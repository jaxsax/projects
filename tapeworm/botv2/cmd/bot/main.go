package main

import (
	"flag"
	"os"
	"path/filepath"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jaxsax/projects/tapeworm/botv2"
	"github.com/jaxsax/projects/tapeworm/botv2/sql"
	_ "github.com/lib/pq"
)

var configPath = flag.String("config_path", "config.yml", "path to config file")

func readConfig() (*botv2.Config, error) {
	fp, err := filepath.Abs(*configPath)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(fp)
	if err != nil {
		return nil, err
	}

	return botv2.ReadConfig(f)
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
		linksRepository = sql.NewLinksRepository(db)
	)

	botAPI, err := botv2.NewTelegramBotAPI(config.Token)
	logErrorAndExit("connect_telegram", err)

	b := botv2.NewBot(
		&botv2.Logger{Logger: logger},
		config,
		linksRepository,
		botAPI,
	)

	err = b.Run()
	logErrorAndExit("run", err)
}
