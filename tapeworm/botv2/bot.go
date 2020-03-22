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
	*Logger
	cfg       *Config
	botAPI    *tgbotapi.BotAPI
	linksDB   *LinksDB
	updatesDB *UpdateDB
}

func NewBot(logger *Logger, config *Config) *Bot {
	return &Bot{
		Logger: logger,
		cfg:    config,
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
	b.Log(
		"action", "init",
		"dependency", "telegram",
	)
	bot, err := tgbotapi.NewBotAPI(b.cfg.Token)
	if err != nil {
		return fmt.Errorf("init telegram: %w", err)
	}
	b.Log(
		"action", "init_ok",
		"dependency", "telegram",
	)

	b.botAPI = bot

	b.Log(
		"action", "init",
		"dependency", "postgres",
	)
	db, err := connectDB(&b.cfg.Database)
	if err != nil {
		return fmt.Errorf("connect db: %w", err)
	}
	b.Log(
		"action", "init_ok",
		"dependency", "postgres",
	)

	b.linksDB = NewLinksDB(db)
	b.updatesDB = NewUpdateDB(db)

	return nil
}

func (b *Bot) Run() error {
	if b.botAPI == nil {
		return errors.New("not initialized yet")
	}

	b.Message("listening for messages")

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botAPI.GetUpdatesChan(u)
	if err != nil {
		b.Log("err", "failed to retrieve updates channel")
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
	log := kitlog.WithPrefix(b.Logger,
		"from", message.From.UserName,
		"from_userid", message.From.ID,
		"message_id", message.MessageID,
	)

	err := b.updatesDB.Create(Update{Data: &update})
	if err != nil {
		log.Log(
			"action", "persist_update",
			"err", err,
		)
	}

	log.Log("message", message.Text)

	switch message.Text {
	case "ping":
		reply := tgbotapi.NewMessage(message.Chat.ID, "pong")
		reply.ReplyToMessageID = message.MessageID
		b.botAPI.Send(reply)
	case "!debug":
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
	case "!links":
		all, err := b.linksDB.List()
		if err != nil {
			log.Log("err", err)
		}

		fmt.Printf("%+v\n", all)
	default:
		if message.Entities != nil {
			res := HandleEntities(message.Text, message.Entities)

			linksToAdd := []Link{}
			for _, entity := range res.Parsed {
				linksToAdd = append(linksToAdd, Link{
					Title: entity,
					Link:  entity,
					ExtraData: map[string]interface{}{
						"created_username": message.From.UserName,
					},
					CreatedTS: int64(message.Date),
					CreatedBy: int64(message.From.ID),
				})
			}
			err := b.linksDB.Create(linksToAdd)
			if err != nil {
				log.Log("err", err)
				return
			}

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

			_, err = b.botAPI.Send(reply)
			if err != nil {
				log.Log("err", err)
			}
		}
	}
}
