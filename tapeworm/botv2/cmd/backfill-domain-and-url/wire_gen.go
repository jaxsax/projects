// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jaxsax/projects/tapeworm/botv2/internal/config"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

// Injectors from wire.go:

func initialize() (*App, error) {
	options := config.ProvideDB()
	store, err := db.Setup(options)
	if err != nil {
		return nil, err
	}
	loggingOptions := config.ProvideLogging()
	logger, err := logging.New(loggingOptions)
	if err != nil {
		return nil, err
	}
	app := &App{
		store:     store,
		dbOptions: options,
		logger:    logger,
	}
	return app, nil
}
