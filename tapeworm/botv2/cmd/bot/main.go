package main

import (
	sql1 "database/sql"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"

	"github.com/jaxsax/projects/tapeworm/botv2"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
	"github.com/jaxsax/projects/tapeworm/botv2/skippedlinks"
	"github.com/jaxsax/projects/tapeworm/botv2/updates"
	"github.com/jaxsax/projects/tapeworm/botv2/web"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

var (
	configPath    = flag.String("config_path", "config.yml", "path to config file")
	sqlitedbpath  = flag.String("sqlite-db-path", "database.db", "path to sqlite .db file")
	port          = flag.Int("port", 8080, "port for webapplication")
	token         = flag.String("telegram-token", "", "telegram token")
	staticDirPath = flag.String("static-dir", "", "path to static dir")
)

func readConfig() (*internal.Config, error) {
	fp, err := filepath.Abs(*configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fp)
	if err != nil {
		if os.IsNotExist(err) {
			// Create something with sane defaults
			var config internal.Config
			config.SqliteDBPath = *sqlitedbpath
			config.Token = *token
			config.Port = *port

			return &config, nil
		}
		return nil, err
	}

	cf, err := internal.ReadConfig(f)
	return cf, err
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

	if config.SqliteDBPath == "" {
		logErrorAndExit("connect_sqlite", errors.New("sqlite_db_path cannot be empty"))
	}
	sqliteDB, err := sql1.Open("sqlite3", config.SqliteDBPath)
	logErrorAndExit("connect_sqlite", err)

	var (
		linksRepository        = links.NewSqliteRepository(sqliteDB)
		skippedLinksRepository = skippedlinks.NewSqliteRepository(sqliteDB)
		updatesRepository      = updates.NewSqliteRepository(sqliteDB)
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
			*staticDirPath,
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
