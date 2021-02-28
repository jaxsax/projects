package internal

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

var (
	// ErrEmptyToken is returned when the token is empty
	ErrEmptyToken = fmt.Errorf("token cannot be empty")

	// ErrEmptyConfig is returned when contents is empty
	ErrEmptyConfig = fmt.Errorf("empty config")

	// ErrInvalidPort is returned when an invalid port number is returned
	ErrInvalidPort = fmt.Errorf("port cannot be empty, WEB_PORT undefined")
)

type DBConfig struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
}

type SonicConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
}

type Config struct {
	Token        string      `yaml:"token"`
	Database     DBConfig    `yaml:"database"`
	SqliteDBPath string      `yaml:"sqlite_db_path"`
	Sonic        SonicConfig `yaml:"sonic"`
	Port         int         `yaml:"-"`
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

	port := os.Getenv("WEB_PORT")
	if port == "" {
		return nil, ErrInvalidPort
	}

	p, err := strconv.ParseUint(port, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	cfg.Port = int(p)

	if cfg.Sonic.Port == 0 {
		cfg.Sonic.Port = 1491
	}

	return &cfg, nil
}
