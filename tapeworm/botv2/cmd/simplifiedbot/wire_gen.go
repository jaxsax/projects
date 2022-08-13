// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/jaxsax/projects/tapeworm/botv2/internal/config"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/httpserver"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/telegrampoller"
)

// Injectors from wire.go:

func initialize() (*App, error) {
	options := config.ProvideHTTP()
	dbOptions := config.ProvideDB()
	store, err := db.Setup(dbOptions)
	if err != nil {
		return nil, err
	}
	loggingOptions := config.ProvideLogging()
	logger, err := logging.New(loggingOptions)
	if err != nil {
		return nil, err
	}
	server := httpserver.New(options, store, logger)
	telegrampollerOptions := config.ProvideTelegram()
	telegramPoller := telegrampoller.New(telegrampollerOptions, store, logger)
	app := &App{
		HTTPServer:     server,
		TelegramSource: telegramPoller,
		Logger:         logger,
	}
	return app, nil
}
