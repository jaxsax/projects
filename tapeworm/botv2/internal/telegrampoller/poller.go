package telegrampoller

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/go-logr/logr"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/enhancers"
	ierrors "github.com/jaxsax/projects/tapeworm/botv2/internal/errors"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	contentblock "github.com/jaxsax/projects/tapeworm/botv2/internal/services/content_block"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Options struct {
	Token                string `long:"telegram_token" description:"telegram bot token" env:"TELEGRAM_BOT_TOKEN"`
	UpdateRequestTimeout int    `long:"telegram_update_request_timeout" description:"how long to keep updates channel open" env:"TELEGRAM_UPDATE_REQUEST_TIMEOUT" default:"30"`
	Enable               bool   `long:"enable_telegram_bot" description:"Whether to enable bot message polling" env:"ENABLE_TELEGRAM_BOT"`
}

type TelegramPoller struct {
	options *Options
	logger  logr.Logger

	botapi              *tgbotapi.BotAPI
	store               *db.Store
	contentBlockChecker *contentblock.Service
	done                chan struct{}
}

func New(opt *Options, store *db.Store, logger logr.Logger, bc *contentblock.Service) *TelegramPoller {
	return &TelegramPoller{
		options:             opt,
		store:               store,
		logger:              logger,
		contentBlockChecker: bc,
		done:                make(chan struct{}, 1),
	}
}

func (p *TelegramPoller) Start(ctx context.Context) error {
	if err := p.contentBlockChecker.Start(); err != nil {
		return err
	}

	tgbotapi.SetLogger(&botLogger{p.logger.WithName("tgbotapi")})
	api, err := tgbotapi.NewBotAPI(p.options.Token)
	if err != nil {
		return err
	}

	api.Debug = true
	p.botapi = api

	u := tgbotapi.NewUpdate(0)
	u.Timeout = p.options.UpdateRequestTimeout

	if p.options.Enable {
		updatesChan := api.GetUpdatesChan(u)
		for update := range updatesChan {
			p.handle(update)
		}
	}

	p.done <- struct{}{}

	return nil
}

func (p *TelegramPoller) handle(update tgbotapi.Update) {
	l := p.logger.WithValues(
		"update_id", update.UpdateID,
	)
	ctx := logr.NewContext(context.Background(), l)

	l.Info("telegram message received", "update", update)

	err := p.store.CreateTelegramUpdate(ctx, &types.TelegramUpdate{
		Data: update,
	})
	if err != nil {
		l.Error(err, "failed to save update")
	}

	if update.Message != nil {
		p.handleMessage(ctx, update.Message)
	}
}

func (p *TelegramPoller) handleMessage(ctx context.Context, message *tgbotapi.Message) {
	if len(message.Entities) == 0 {
		return
	}

	err := p.contentBlockChecker.IsAllowed(ctx, message.Text, "telegram_source_msg_text")
	if err != nil {
		p.replyWithErrorV2(ctx, err, message)
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
			URL: ExtractURL(message.Text, entity.Offset, entity.Length),
		}

		resp, err := p.linkProcessor(ctx, req)
		if err != nil {
			p.replyWithErrorV2(ctx, err, message)
			return
		}

		processedLinkGroup = append(processedLinkGroup, resp)
	}

	if !anyLinks {
		return
	}

	if len(processedLinkGroup) == 0 {
		p.replyWithError(ctx, fmt.Errorf("no links found"), "no link found", message)
		return
	}

	// Persist to storage
	processedLinks := make([]*types.Link, 0, len(processedLinkGroup))
	for _, link := range processedLinkGroup {
		path := link.URL.EscapedPath()
		if link.URL.RawQuery != "" {
			path += "?" + link.URL.RawQuery
		}

		if link.URL.Fragment != "" {
			path += "#" + link.URL.EscapedFragment()
		}

		logging.FromContext(ctx).Info("link", "item", link)

		lt := &types.Link{
			Link:        link.URL.String(),
			Title:       link.Title,
			CreatedAt:   uint64(time.Now().Unix()),
			CreatedByID: uint64(message.From.ID),
			ExtraData:   map[string]string{},
			Domain:      link.URL.Hostname(),
			Path:        path,
		}
		processedLinks = append(processedLinks, lt)
	}

	if err := p.store.UpsertLinks(ctx, processedLinks); err != nil {
		logging.FromContext(ctx).Error(err, "failed to persist links to storage")
		p.replyWithError(ctx, err, "failed to store links", message)
		return
	}

	// Reply to the user
	processedMessage := processedLinkGroupMessageBody(processedLinkGroup)
	m := NewReplyToMessage(processedMessage, message)
	m.DisableNotification = true

	_, err = p.botapi.Send(m)
	if err != nil {
		logging.FromContext(ctx).Error(err, "failed to send processed message")
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

func (p *TelegramPoller) replyWithErrorV2(
	ctx context.Context,
	err error,
	originMessage *tgbotapi.Message,
) {
	if err == nil {
		return
	}

	var chattable tgbotapi.Chattable

	var ufe ierrors.UserFacingError
	if errors.As(err, &ufe) {
		if ufe.UserResponse() != "" {
			m := tgbotapi.NewMessage(originMessage.Chat.ID, ufe.UserResponse())
			m.ReplyToMessageID = originMessage.MessageID
			chattable = m
		}

		if ufe.StickerResponse() != "" {
			m := tgbotapi.NewSticker(originMessage.Chat.ID, tgbotapi.FileID(ufe.StickerResponse()))
			m.ReplyToMessageID = originMessage.MessageID
			chattable = m
		}
	}

	if chattable == nil {
		// No specific error types are matched, return generic error
		m := tgbotapi.NewMessage(originMessage.Chat.ID, "Internal server error")
		m.ReplyToMessageID = originMessage.MessageID
		chattable = m
	}

	logging.FromContext(ctx).Info("replying to user with error", "internal_err", err)
	_, err = p.botapi.Send(chattable)
	if err != nil {
		logging.FromContext(ctx).Error(err, "failed to send error reply")
		return
	}
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

	logging.FromContext(ctx).Error(processErr, clientFacingMessage)
	_, err := p.botapi.Send(m)
	if err != nil {
		logging.FromContext(ctx).Error(err, "failed to reply with error")
		return
	}
}

type processLinkRequest struct {
	URL string
}

type processLinkResponse struct {
	Title string
	URL   *url.URL
}

func (p *TelegramPoller) linkProcessor(ctx context.Context, req *processLinkRequest) (*processLinkResponse, error) {
	l, err := enhancers.EnhanceLinkWithContext(ctx, req.URL, p.store)
	if err != nil {
		return nil, err
	}

	lurl, err := url.Parse(l.Link)
	if err != nil {
		return nil, err
	}

	if l.Title == "" {
		return nil, fmt.Errorf("title is empty")
	}

	if err := p.contentBlockChecker.IsAllowed(ctx, l.Title, "link_processor_title_check"); err != nil {
		return nil, err
	}

	return &processLinkResponse{
		Title: l.Title,
		URL:   lurl,
	}, nil
}

func (p *TelegramPoller) Stop(ctx context.Context) error {
	if !p.options.Enable {
		return nil
	}

	p.botapi.StopReceivingUpdates()

	p.logger.V(0).Info("waiting for telegram updates to drain")
	select {
	case <-ctx.Done():
		return fmt.Errorf("stop timed out")
	case <-p.done:
		return nil
	}
}
