package telegrampoller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jaxsax/projects/tapeworm/botv2/enhancers"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Options struct {
	Token                string `long:"telegram_token" description:"telegram bot token" env:"TELEGRAM_BOT_TOKEN"`
	UpdateRequestTimeout int    `long:"telegram_update_request_timeout" description:"how long to keep updates channel open" env:"TELEGRAM_UPDATE_REQUEST_TIMEOUT" default:"30"`
}

type TelegramPoller struct {
	options *Options
	logger  logr.Logger

	botapi *tgbotapi.BotAPI
	store  *db.Store
	done   chan struct{}
}

func New(opt *Options, store *db.Store, logger logr.Logger) *TelegramPoller {
	return &TelegramPoller{
		options: opt,
		store:   store,
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
	l := p.logger.WithValues(
		"from_chat_id", update.FromChat().ID,
		"messsage_id", update.Message.MessageID,
	)
	ctx := logr.NewContext(context.Background(), l)

	l.Info("telegram message received", "update", update)

	if update.Message != nil {
		p.handleMessage(ctx, update.Message)
	}
}

func (p *TelegramPoller) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if len(message.Entities) == 0 {
		return
	}

	var (
		processedLinkGroup []*processLinkResponse
		anyLinks           bool
	)
	for _, entity := range message.Entities {
		if !entity.IsURL() {
			continue
		}

		anyLinks = true

		req := &processLinkRequest{
			URL: message.Text[entity.Offset : entity.Offset+entity.Length],
		}

		resp, err := p.linkProcessor(ctx, req)
		if err != nil {
			p.replyWithError(ctx, err, "failed to process link", message)
			return
		}

		processedLinkGroup = append(processedLinkGroup, resp)
	}

	if len(processedLinkGroup) == 0 && anyLinks {
		p.replyWithError(ctx, fmt.Errorf("no links found"), "no link found", message)
		return
	}

	// Persist to storage
	processedLinks := make([]*types.Link, 0, len(processedLinkGroup))
	for _, link := range processedLinkGroup {
		lt := &types.Link{
			Link:        link.Link,
			Title:       link.Title,
			CreatedAt:   time.Now(),
			CreatedByID: uint64(message.From.ID),
			ExtraData:   map[string]string{},
		}

		processedLinks = append(processedLinks, lt)
	}

	if err := p.store.CreateLinks(ctx, processedLinks); err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "failed to persist links to storage")
		p.replyWithError(ctx, err, "failed to store links", message)
		return
	}

	// Reply to the user
	processedMessage := processedLinkGroupMessageBody(processedLinkGroup)
	m := NewReplyToMessage(processedMessage, message)
	m.DisableNotification = true

	_, err := p.botapi.Send(m)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "failed to send processed message")
		return
	}
}

func processedLinkGroupMessageBody(g []*processLinkResponse) string {
	var sb strings.Builder

	for i, r := range g {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, r.Title))
	}

	return sb.String()
}

func (p *TelegramPoller) replyWithError(
	ctx context.Context,
	processErr error,
	clientFacingMessage string,
	originMessage *tgbotapi.Message,
) {
	// todo: mask server sided errors
	m := tgbotapi.NewMessage(originMessage.Chat.ID, clientFacingMessage)
	m.ReplyToMessageID = originMessage.MessageID

	logr.FromContextOrDiscard(ctx).Error(processErr, clientFacingMessage)
	_, err := p.botapi.Send(m)
	if err != nil {
		logr.FromContextOrDiscard(ctx).Error(err, "failed to reply with error")
		return
	}
}

type processLinkRequest struct {
	URL string
}

type processLinkResponse struct {
	Title string
	Link  string
}

func (p *TelegramPoller) linkProcessor(ctx context.Context, req *processLinkRequest) (*processLinkResponse, error) {
	l, err := enhancers.EnhanceLinkWithContext(ctx, req.URL)
	if err != nil {
		return nil, err
	}

	return &processLinkResponse{
		Title: l.Title,
		Link:  l.Link,
	}, nil
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
