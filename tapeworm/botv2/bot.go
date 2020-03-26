// https://peter.bourgon.org/go-best-practices-2016/#configuration
package botv2

import (
	"fmt"

	kitlog "github.com/go-kit/kit/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jaxsax/projects/tapeworm/botv2/enhancers"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
	"github.com/jaxsax/projects/tapeworm/botv2/updates"
)

type Bot struct {
	*internal.Logger
	cfg               *internal.Config
	botAPI            *tgbotapi.BotAPI
	updatesRepository updates.Repository
	linksRepository   links.Repository
}

func NewBot(
	logger *internal.Logger,
	config *internal.Config,
	linksRepository links.Repository,
	updatesRepository updates.Repository,
	botAPI *tgbotapi.BotAPI,
) *Bot {
	return &Bot{
		Logger:            logger,
		cfg:               config,
		linksRepository:   linksRepository,
		botAPI:            botAPI,
		updatesRepository: updatesRepository,
	}
}

func (b *Bot) Run() error {
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

	err := b.updatesRepository.Create(updates.Update{Data: &update})
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
	case "!links":
		all := b.linksRepository.List()
		log.Log("links", fmt.Sprintf("%+v\n", all))
	default:
		if message.Entities != nil {
			res := HandleEntities(message.Text, message.Entities)

			if len(res.Parsed) == 0 {
				return
			}

			linksToAdd := []links.Link{}
			bodyParsed := ""
			for i, entity := range res.Parsed {
				enhancedLink, err := enhancers.EnhanceLink(entity)
				if err != nil {
					log.Log(
						"action", "parse_link",
						"err", err,
						"url", entity,
					)
					continue
				}

				linksToAdd = append(linksToAdd, links.Link{
					Title: enhancedLink.Title,
					Link:  enhancedLink.Link,
					ExtraData: map[string]interface{}{
						"created_username": message.From.UserName,
					},
					CreatedTS: int64(message.Date),
					CreatedBy: int64(message.From.ID),
				})
				bodyParsed += fmt.Sprintf("%v. %v\n", i+1, enhancedLink.Title)
			}
			err := b.linksRepository.CreateMany(linksToAdd)
			if err != nil {
				log.Log("err", err)
				return
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
