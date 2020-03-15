// https://peter.bourgon.org/go-best-practices-2016/#configuration
package bot_v2

import (
	"fmt"
	"log"

	kitlog "github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	logger *Logger
	cfg    *Config
}

type Logger struct {
	kitlog.Logger
}

func (l *Logger) Message(s string) {
	l.Log("msg", s)
}

func NewBot(configPath string) *Bot {
	if configPath == "" {
		configPath = "config.yml"
	}
	cfg := readConfig(configPath)
	return &Bot{
		logger: &Logger{cfg.Logger},
		cfg:    cfg,
	}
}

func (b *Bot) Run() error {
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	b.logger.Message("init")

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		b.logger.Log("err", "failed to retrieve updates channel")
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		b.logger.Log("update", fmt.Sprintf("%+v", update))
	}

	return nil
}
