// https://peter.bourgon.org/go-best-practices-2016/#configuration
package bot_v2

import (
	"errors"
	"fmt"

	kitlog "github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Bot struct {
	logger *Logger
	cfg    *Config
	botApi *tgbotapi.BotAPI
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

func (b *Bot) Init() error {
	b.logger.Message("init connection to telegram")

	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return fmt.Errorf("init: %w", err)
	}

	b.botApi = bot

	return nil
}

func (b *Bot) Run() error {
	if b.botApi == nil {
		return errors.New("not initialized yet")
	}

	b.logger.Message("listening for messages")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botApi.GetUpdatesChan(u)
	if err != nil {
		b.logger.Log("err", "failed to retrieve updates channel")
	}

	for update := range updates {
		b.handleUpdate(update)
	}

	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message
	log := kitlog.WithPrefix(b.logger,
		"from", message.From.UserName,
		"from_userid", message.From.ID,
		"message_id", message.MessageID,
	)

	log.Log("message", message.Text)

	if message.Text == "ping" {
		reply := tgbotapi.NewMessage(message.Chat.ID, "pong")
		reply.ReplyToMessageID = message.MessageID
		b.botApi.Send(reply)
	}

	if message.Entities != nil {
		res := HandleEntities(message.Text, message.Entities)

		bodyParsed := ""
		for i, url := range res.Parsed {
			bodyParsed += fmt.Sprintf("%v. %v\n", i+1, url)
		}
		body := fmt.Sprintf(`<b>Links parsed</b>

%v
`, bodyParsed)

		reply := tgbotapi.NewMessage(message.Chat.ID, body)
		reply.ParseMode = "HTML"
		reply.DisableNotification = true
		reply.DisableWebPagePreview = true
		reply.ReplyToMessageID = message.MessageID

		_, err := b.botApi.Send(reply)
		if err != nil {
			log.Log("err", err)
		}
	}
}
