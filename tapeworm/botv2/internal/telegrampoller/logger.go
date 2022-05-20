package telegrampoller

import (
	"fmt"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type botLogger struct {
	logger logr.Logger
}

var _ tgbotapi.BotLogger = &botLogger{}

func (l *botLogger) Println(v ...interface{}) {
	if len(v) == 0 {
		return
	}

	v0 := v[0]
	switch vt := v0.(type) {
	case error:
		l.logger.V(0).Error(vt, "")
	default:
		l.logger.V(1).Info(fmt.Sprintf("%s", v0))
	}
}

func (l *botLogger) Printf(format string, v ...interface{}) {
	l.logger.V(1).Info(fmt.Sprintf(format, v...))
}
