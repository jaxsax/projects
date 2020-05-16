package main

import (
	"flag"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/jaxsax/projects/tapeworm/botv2"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/sql"
	"github.com/jaxsax/projects/tapeworm/botv2/web"
	_ "github.com/lib/pq"
)

var configPath = flag.String("config_path", "config.yml", "path to config file")

func readConfig() (*internal.Config, error) {
	fp, err := filepath.Abs(*configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	return internal.ReadConfig(f)
}

func main() {
	flag.Parse()

	lcfg := zap.NewDevelopmentConfig()
	lcfg.OutputPaths = []string{"stdout"}

	options := []zap.Option{
		internal.Core(),
	}
	logger, err := lcfg.Build(options...)
	if err != nil {
		zap.L().Fatal("error building logger", zap.Error(err))
	}

	logErrorAndExit := func(action string, err error) {
		if err != nil {
			logger.Fatal("init error", zap.String("action", action), zap.Error(err))
		}
	}

	config, err := readConfig()
	logErrorAndExit("read_config", err)

	db, err := botv2.ConnectDB(&config.Database)
	logErrorAndExit("connect_db", err)

	var (
		linksRepository        = sql.NewLinksRepository(db)
		skippedLinksRepository = sql.NewSkippedLinksRepository(db)
		updatesRepository      = sql.NewUpdatesRepository(db)
	)

	botAPI, err := botv2.NewTelegramBotAPI(config.Token)
	logErrorAndExit("connect_telegram", err)

	// https://medium.com/rungo/running-multiple-http-servers-in-go-d15300f4e59f
	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		componentLogger := logger.Named("app.bot")

		componentLogger.Info(
			"starting",
		)

		b := botv2.NewBot(
			componentLogger,
			config,
			linksRepository,
			updatesRepository,
			skippedLinksRepository,
			botAPI,
		)

		err := b.Run()
		componentLogger.Info(
			"stopped",
			zap.String("state", "stopped"),
			zap.Error(err),
		)

		wg.Done()
	}()

	go func() {
		componentLogger := logger.Named("app.botapi")

		componentLogger.Info(
			"starting",
		)

		webServer := web.NewServer(
			componentLogger,
			config,
			linksRepository,
		)

		err := webServer.Run()
		componentLogger.Info(
			"stopped",
			zap.Error(err),
		)
		wg.Done()
	}()

	wg.Wait()
}
