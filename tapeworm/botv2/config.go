package botv2

import (
	"fmt"
	"io"

	kitlog "github.com/go-kit/kit/log"
	"gopkg.in/yaml.v2"
)

var (
	// ErrEmptyToken is returned when the token is empty
	ErrEmptyToken = fmt.Errorf("token cannot be empty")

	// ErrEmptyConfig is returned when contents is empty
	ErrEmptyConfig = fmt.Errorf("empty config")
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

func ReadConfig(r io.Reader) (*Config, error) {
	var cfg Config

	err := yaml.NewDecoder(r).Decode(&cfg)
	if err != nil {
		if err == io.EOF {
			return nil, ErrEmptyConfig
		}
		return nil, err
	}

	if cfg.Token == "" {
		return nil, ErrEmptyToken
	}

	return &cfg, nil
}
