package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/httpserver"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/telegrampoller"
)

type App struct {
	HTTPServer     *httpserver.Server
	TelegramSource *telegrampoller.TelegramPoller
	Logger         logr.Logger
}

func main() {
	app, err := initialize()
	if err != nil {
		log.Printf("panic initializing: %v", err)
		os.Exit(1)
	}

	done := make(chan struct{}, 1)
	waitSigterm := make(chan os.Signal, 1)
	signal.Notify(waitSigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.HTTPServer.Start(); err != nil {
			app.Logger.V(0).Error(err, "start httpserver")
			done <- struct{}{}
		}
	}()

	go func() {
		if err := app.TelegramSource.Start(); err != nil {
			app.Logger.V(0).Error(err, "start telegram poller")
			done <- struct{}{}
		}
	}()

	go func() {
		<-waitSigterm

		app.Logger.V(0).Info("interrupt received")
		done <- struct{}{}
	}()

	<-done
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := app.TelegramSource.Stop(shutdownCtx); err != nil {
		app.Logger.V(0).Error(err, "telegram poller shutdown")
	}

	if err := app.HTTPServer.Stop(shutdownCtx); err != nil {
		app.Logger.V(0).Error(err, "httpserver shutdown")
	}
}
