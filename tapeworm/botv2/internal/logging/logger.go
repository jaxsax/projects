package logging

import "github.com/go-logr/logr"

var Logger logr.Logger

func Info(msg string, keysAndValues ...interface{}) {
	Logger.WithCallDepth(1).Info(msg, keysAndValues...)
}

func Error(err error, msg string, keysAndValues ...interface{}) {
	Logger.WithCallDepth(1).Error(err, msg, keysAndValues...)
}

func V(level int) logr.Logger {
	return Logger.V(level)
}

func WithValues(keysAndValues ...interface{}) logr.Logger {
	return Logger.WithValues(keysAndValues...)
}
