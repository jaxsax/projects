package telegrampoller

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Options struct {
	Token                string `long:"telegram_token" description:"telegram bot token" env:"TELEGRAM_BOT_TOKEN"`
	UpdateRequestTimeout int    `long:"telegram_update_request_timeout" description:"how long to keep updates channel open" env:"TELEGRAM_UPDATE_REQUEST_TIMEOUT" default:"30"`
}

type TelegramPoller struct {
	options *Options
	logger  logr.Logger

	botapi *tgbotapi.BotAPI
	done   chan struct{}
}

func New(opt *Options, logger logr.Logger) *TelegramPoller {
	return &TelegramPoller{
		options: opt,
		logger:  logger,
		done:    make(chan struct{}, 1),
	}
}

func (p *TelegramPoller) Start() error {
	tgbotapi.SetLogger(&botLogger{p.logger.WithName("tgbotapi")})
	api, err := tgbotapi.NewBotAPI(p.options.Token)
	if err != nil {
		return err
	}

	api.Debug = true
	p.botapi = api

	u := tgbotapi.NewUpdate(0)
	u.Timeout = p.options.UpdateRequestTimeout

	updatesChan := api.GetUpdatesChan(u)
	for update := range updatesChan {
		p.handle(update)
	}

	p.done <- struct{}{}

	return nil
}

func (p *TelegramPoller) handle(update tgbotapi.Update) {
	l := p.logger.WithValues("message_id", update.Message.MessageID)

	l.Info("telegram message", "update", update)
}

func (p *TelegramPoller) Stop(ctx context.Context) error {
	p.botapi.StopReceivingUpdates()

	p.logger.V(0).Info("waiting for telegram updates to drain")
	select {
	case <-ctx.Done():
		return fmt.Errorf("stop timed out")
	case <-p.done:
		return nil
	}
}
