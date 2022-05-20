package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/httpserver"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

var (
	flagParser        = flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)
	logger            logr.Logger
	httpserverOptions = &httpserver.Options{}
	logOptions        = &loggingOptions{}
)

type loggingOptions struct {
	DevelopmentLog bool `long:"pretty_logs" description:"use the nicer-to-look at development log" env:"PRETTY_LOGS"`
}

func main() {
	if _, err := flagParser.AddGroup("http", "", httpserverOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.AddGroup("logging", "", logOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}

	if err := setupLogger(logOptions); err != nil {
		panic(err)
	}

	done := make(chan struct{}, 1)
	waitSigterm := make(chan os.Signal, 1)
	signal.Notify(waitSigterm, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	httpserver := httpserver.New(httpserverOptions, logger)
	go func() {
		if err := httpserver.Start(); err != nil {
			logger.V(0).Error(err, "start httpserver")
			done <- struct{}{}
		}
	}()

	go func() {
		<-waitSigterm

		logger.V(0).Info("interrupt received")
		done <- struct{}{}
	}()

	<-done
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	if err := httpserver.Stop(shutdownCtx); err != nil {
		logger.V(0).Error(err, "httpserver shutdown")
	}
}

func setupLogger(opts *loggingOptions) error {
	var (
		zapl *zap.Logger
		err  error
	)

	zapl, err = zap.NewProduction()
	if err != nil {
		return err
	}

	if opts.DevelopmentLog {
		zapl, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}

	logger = zapr.NewLogger(zapl)
	return nil
}
