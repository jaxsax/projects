package main

import (
	"flag"
	"os"
	"path/filepath"

	kitlog "github.com/go-kit/kit/log"
	"github.com/jaxsax/projects/tapeworm/botv2"
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

	b := botv2.NewBot(&botv2.Logger{Logger: logger}, config)
	err = b.Init()
	logErrorAndExit("init", err)

	err = b.Run()
	logErrorAndExit("run", err)
}
