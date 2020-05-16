package internal

import (
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Core() zap.Option {
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		core := &CannonicalLog{
			EmptyCore:   c,
			WrappedCore: c,
		}
		return core
	})
}

func Emit(log *zap.Logger, fields ...zap.Field) error {
	c, ok := log.Core().(*CannonicalLog)
	if !ok {
		return errors.New("unknown logger type")
	}
	if len(c.Fields)+len(fields) == 0 {
		return nil
	}

	concatenatedFields := make([]zap.Field, 0, len(c.Fields)+len(fields))
	for _, f := range c.Fields {
		concatenatedFields = append(concatenatedFields, f)
	}

	for _, f := range fields {
		concatenatedFields = append(concatenatedFields, f)
	}

	if err := c.EmptyCore.Write(zapcore.Entry{
		Time:    time.Now(),
		Message: "cannonical_log_line",
	}, concatenatedFields); err != nil {
		return err
	}

	c.Reset()

	if err := c.EmptyCore.Sync(); err != nil {
		return err
	}

	return nil
}
