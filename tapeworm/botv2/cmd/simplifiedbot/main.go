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
	"github.com/jaxsax/projects/tapeworm/botv2/internal/managed"
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

	lifecycleManager := managed.New()
	lifecycleManager.Add(app.HTTPServer, "http_server")
	lifecycleManager.Add(app.TelegramSource, "telegram_source")

	go func() {
		if err := lifecycleManager.Start(context.Background()); err != nil {
			app.Logger.V(0).Error(err, "start lifecycle failed")
			os.Exit(1)
		}
	}()

	go func() {
		<-waitSigterm

		app.Logger.V(0).Info("interrupt received")
		done <- struct{}{}
	}()

	<-done
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := lifecycleManager.Stop(shutdownCtx); err != nil {
		app.Logger.V(0).Error(err, "stop lifecyle failed")
		os.Exit(1)
	}
}
