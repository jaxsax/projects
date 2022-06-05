package main

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

var (
	flagParser = flags.NewParser(nil, flags.HelpFlag|flags.PassDoubleDash)
	dbOptions  = &db.Options{}
	logOptions = &loggingOptions{}
)

type loggingOptions struct {
	DevelopmentLog bool `long:"pretty_logs" description:"use the nicer-to-look at development log" env:"PRETTY_LOGS"`
}

func main() {
	if _, err := flagParser.AddGroup("logging", "", logOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.AddGroup("db", "", dbOptions); err != nil {
		panic(err)
	}

	if _, err := flagParser.Parse(); err != nil {
		panic(err)
	}

	logger, err := setupLogger(logOptions)
	if err != nil {
		panic(err)
	}

	store, err := db.Setup(dbOptions)
	if err != nil {
		panic(err)
	}

	err = store.IndexAllItems(logr.NewContext(context.Background(), *logger))
	if err != nil {
		panic(err)
	}
}

func setupLogger(opts *loggingOptions) (*logr.Logger, error) {
	var (
		zapl *zap.Logger
		err  error
	)

	zapl, err = zap.NewProduction()
	if err != nil {
		return nil, err
	}

	if opts.DevelopmentLog {
		zapl, err = zap.NewDevelopment()
		if err != nil {
			return nil, err
		}
	}

	l := zapr.NewLogger(zapl)
	return &l, nil
}
