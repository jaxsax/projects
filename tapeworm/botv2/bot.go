// https://peter.bourgon.org/go-best-practices-2016/#configuration
package botv2

import (
	"errors"
	"fmt"

	kitlog "github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jmoiron/sqlx"
)

type Bot struct {
	logger  *Logger
	cfg     *Config
	botAPI  *tgbotapi.BotAPI
	linksDB *LinksDB
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

func connectDB(conf *DBConfig) (*sqlx.DB, error) {
	connString := fmt.Sprintf(
		"user=%s password=%s host=%s dbname=%s sslmode=disable",
		conf.User, conf.Password, conf.Host, conf.Name,
	)
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (b *Bot) Init() error {
	b.logger.Log(
		"action", "init",
		"dependency", "telegram",
	)
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return fmt.Errorf("init telegram: %w", err)
	}
	b.logger.Log(
		"action", "init_ok",
		"dependency", "telegram",
	)

	b.botAPI = bot

	b.logger.Log(
		"action", "init",
		"dependency", "postgres",
	)
	db, err := connectDB(&b.cfg.Database)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	b.logger.Log(
		"action", "init_ok",
		"dependency", "postgres",
	)

	b.linksDB = NewLinksDB(db)

	return nil
}

func (b *Bot) Run() error {
	if b.botAPI == nil {
		return errors.New("not initialized yet")
	}

	b.logger.Message("listening for messages")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botAPI.GetUpdatesChan(u)
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
		b.botAPI.Send(reply)
	} else if message.Text == "!debug" {
		body := fmt.Sprintf(`
Database:
	%+v
`, b.linksDB.db.Stats())[1:]

		reply := tgbotapi.NewMessage(message.Chat.ID, body)
		reply.DisableNotification = true
		reply.DisableWebPagePreview = true
		reply.ReplyToMessageID = message.MessageID

		_, err := b.botAPI.Send(reply)
		if err != nil {
			log.Log("err", err)
		}
	}

	if message.Entities != nil {
		res := HandleEntities(message.Text, message.Entities)

		bodyParsed := ""
		for i, url := range res.Parsed {
			bodyParsed += fmt.Sprintf("%v. %v\n", i+1, url)
		}
		body := fmt.Sprintf(`
<b>Links parsed</b>

%v
`, bodyParsed)[1:]

		reply := tgbotapi.NewMessage(message.Chat.ID, body)
		reply.ParseMode = "HTML"
		reply.DisableNotification = true
		reply.DisableWebPagePreview = true
		reply.ReplyToMessageID = message.MessageID

		_, err := b.botAPI.Send(reply)
		if err != nil {
			log.Log("err", err)
		}
	}
}
