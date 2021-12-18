package internal

import (
	"fmt"
	"io"

	"gopkg.in/yaml.v2"
)

var (
	// ErrEmptyToken is returned when the token is empty
	ErrEmptyToken = fmt.Errorf("token cannot be empty")

	// ErrEmptyConfig is returned when contents is empty
	ErrEmptyConfig = fmt.Errorf("empty config")

	ErrEmptyPort = fmt.Errorf("port cannot be empty")
)

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type Config struct {
	Token        string   `yaml:"token"`
	Database     DBConfig `yaml:"database"`
	SqliteDBPath string   `yaml:"sqlite_db_path"`
	Port         int      `yaml:"port"`
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

	if cfg.Port == 0 {
		return nil, ErrEmptyPort
	}

	return &cfg, nil
}
