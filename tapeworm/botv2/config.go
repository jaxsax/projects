package botv2

import (
	"errors"
	"os"
	"path/filepath"

	kitlog "github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"
)

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Token    string        `yaml:"token"`
	Database DBConfig      `yaml:"database"`
	Logger   kitlog.Logger `yaml:"-"`
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

	fp, err := filepath.Abs(configPath)
	handleError(log, err)

	f, err := os.Open(fp)

	handleError(log, err)
	defer func() {
		handleError(log, f.Close())
	}()

	var cfg Config
	cfg.Logger = log

	decoder := yaml.NewDecoder(f)
	handleError(log, decoder.Decode(&cfg))

	if cfg.Token == "" {
		handleError(log, errors.New("invalid token"))
	}

	return &cfg
}
