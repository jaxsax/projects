package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	kitlog "github.com/go-kit/kit/log"
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

	lw := kitlog.NewSyncWriter(os.Stderr)
	logger := kitlog.NewLogfmtLogger(lw)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	logErrorAndExit := func(action string, err error) {
		if err != nil {
			logger.Log(
				"action", action,
				"err", err,
			)
			os.Exit(1)
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
		logger.Log(
			"action", "starting",
			"component", "bot",
		)

		b := botv2.NewBot(
			&internal.Logger{Logger: logger},
			config,
			linksRepository,
			updatesRepository,
			skippedLinksRepository,
			botAPI,
		)

		err := b.Run()
		logger.Log(
			"action", "ended",
			"component", "bot",
			"err", fmt.Sprintf("%+v", err),
		)
		wg.Done()
	}()

	go func() {
		logger.Log(
			"action", "starting",
			"component", "web",
		)

		webServer := web.NewServer(
			&internal.Logger{Logger: logger},
			config,
			linksRepository,
		)

		err := webServer.Run()

		logger.Log(
			"action", "ended",
			"component", "web",
			"err", fmt.Sprintf("%+v", err),
		)
		wg.Done()
	}()

	wg.Wait()
}
