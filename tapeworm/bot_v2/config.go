package bot_v2

import (
	"errors"
	"os"

	kitlog "github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token  string        `yaml:"token"`
	Logger kitlog.Logger `yaml:"-"`
}

func handleError(log kitlog.Logger, err error) {
	if err != nil {
		log.Log("err", err, "msg", "failed to parse config")
		os.Exit(1)
	}
}

func initLogger() kitlog.Logger {
	w := kitlog.NewSyncWriter(os.Stderr)
	logger := kitlog.NewLogfmtLogger(w)
	logger = kitlog.With(logger, "ts", kitlog.DefaultTimestampUTC, "caller", kitlog.DefaultCaller)

	return logger
}

func readConfig(configPath string) *Config {
	log := initLogger()

	f, err := os.Open(configPath)

	handleError(log, err)
	defer func() {
		handleError(log, f.Close())
	}()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	handleError(log, decoder.Decode(&cfg))

	if cfg.Token == "" {
		handleError(log, errors.New("invalid token"))
	}

	cfg.Logger = log

	return &cfg
}
