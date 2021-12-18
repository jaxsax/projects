// https://peter.bourgon.org/go-best-practices-2016/#configuration
package botv2

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jaxsax/projects/tapeworm/botv2/enhancers"
	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
	"github.com/jaxsax/projects/tapeworm/botv2/skippedlinks"
	"github.com/jaxsax/projects/tapeworm/botv2/updates"
	"go.uber.org/zap"
)

type Bot struct {
	Logger                 *zap.Logger
	cfg                    *internal.Config
	botAPI                 *tgbotapi.BotAPI
	updatesRepository      updates.Repository
	linksRepository        links.Repository
	skippedLinksRepository skippedlinks.Repository
}

func NewBot(
	logger *zap.Logger,
	config *internal.Config,
	linksRepository links.Repository,
	updatesRepository updates.Repository,
	skippedLinksRepository skippedlinks.Repository,
	botAPI *tgbotapi.BotAPI,
) *Bot {
	return &Bot{
		Logger:                 logger,
		cfg:                    config,
		linksRepository:        linksRepository,
		botAPI:                 botAPI,
		updatesRepository:      updatesRepository,
		skippedLinksRepository: skippedLinksRepository,
	}
}

func (b *Bot) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.botAPI.GetUpdatesChan(u)
	if err != nil {
		return fmt.Errorf("get updates: %w", err)
	}

	b.Logger.Info("listening for messages")
	for update := range updates {
		go b.handleUpdate(update)
	}

	return nil
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	message := update.Message

	start := time.Now()
	ctxLogger := b.Logger.With()

	defer func() {
		// Placing this here until I figure out why using b.Logger.With() causes duplicates
		// in canonical logs
		internal.Emit(
			ctxLogger,
			zap.Int("update_id", update.UpdateID),
			zap.String("from", message.From.UserName),
			zap.Int("from_userid", message.From.ID),
			zap.Int("message_id", message.MessageID),
			zap.Duration("update_duration", time.Since(start)),
		)
	}()

	err := b.updatesRepository.Create(updates.Update{Data: &update})
	if err != nil {
		ctxLogger.Error("save update failed", zap.Error(err))
	}

	ctxLogger.Debug("message received", zap.String("message", message.Text))

	switch message.Text {
	case "ping":
		reply := tgbotapi.NewMessage(message.Chat.ID, "pong")
		reply.ReplyToMessageID = message.MessageID
		b.botAPI.Send(reply)
	case "!links":
		// all := b.linksRepository.List()

		ctxLogger.Info(
			"command received",
			zap.String("cmd", "links"),
		)
	default:
		if message.Entities != nil {
			ctxLogger.Debug(
				"parsing entities",
			)
			res := HandleEntities(message.Text, message.Entities)

			if len(res.Parsed) == 0 {
				return
			}

			linksToAdd := []links.Link{}
			bodyParsed := ""
			for i, entity := range res.Parsed {
				enhancedLink, err := enhancers.EnhanceLink(entity)
				if err != nil {
					ctxLogger.Error(
						"parse link failed",
						zap.Bool("parse_ok", false),
						zap.Error(err),
					)

					skippedLink := skippedlinks.SkippedLink{
						Link:        entity,
						ErrorReason: err.Error(),
					}
					serr := b.skippedLinksRepository.Create(skippedLink)
					if serr != nil {
						ctxLogger.Error(
							"store skipped link failed",
							zap.Bool("store_skipped_ok", false),
							zap.Error(serr),
						)
					}
					continue
				}

				linksToAdd = append(linksToAdd, links.Link{
					Title: enhancedLink.Title,
					Link:  enhancedLink.Link,
					ExtraData: map[string]string{
						"created_username": message.From.UserName,
					},
					CreatedTS: int64(message.Date),
					CreatedBy: int64(message.From.ID),
				})
				bodyParsed += fmt.Sprintf("%v. %v\n", i+1, enhancedLink.Title)
			}
			err := b.linksRepository.CreateMany(linksToAdd)
			if err != nil {
				ctxLogger.Error(
					"store links failed",
					zap.Bool("store_ok", false),
					zap.Error(err),
				)
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
				ctxLogger.Error(
					"send response failed",
					zap.Bool("send_response_ok", false),
					zap.Error(err),
				)
			}
		}
	}
}
