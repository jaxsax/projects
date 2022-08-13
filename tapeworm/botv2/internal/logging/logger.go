package logging

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/google/wire"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger logr.Logger

type Options struct {
	DevelopmentLog bool `long:"pretty_logs" description:"use the nicer-to-look at development log" env:"PRETTY_LOGS"`
}

func Info(msg string, keysAndValues ...interface{}) {
	logger.WithCallDepth(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	logger.WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func V(level int) logr.Logger {
	return logger.V(level)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return logger.WithValues(keysAndValues...)
}

func WithContext(ctx context.Context) context.Context {
	return logr.NewContext(ctx, logger)
}

func FromContext(ctx context.Context) logr.Logger {
	l, err := logr.FromContext(ctx)
	if err == nil {
		return l
	}

	return logger
}

func New(opts *Options) (logr.Logger, error) {
	var (
		zapl *zap.Logger
		err  error
	)

	zapl, err = zap.NewProduction()
	if err != nil {
		return logr.Logger{}, err
	}

	if opts.DevelopmentLog {
		devConfig := zap.NewDevelopmentConfig()
		devConfig.Level = zap.NewAtomicLevelAt(zapcore.Level(-2))
		zapl, err = devConfig.Build()
		if err != nil {
			return logr.Logger{}, err
		}
	}

	logrSink := zapr.NewLogger(zapl)
	logger = logrSink
	return logrSink, nil
}

var ProviderSet = wire.NewSet(
	New,
)
