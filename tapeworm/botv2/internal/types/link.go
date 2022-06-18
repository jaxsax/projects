package types

import (
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Link struct {
	ID          uint64            `json:"id"`
	Link        string            `json:"link"`
	Title       string            `json:"title"`
	CreatedAt   uint64            `json:"created_ts"`
	CreatedByID uint64            `json:"created_by"`
	ExtraData   map[string]string `json:"extra_data"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty"`
	Domain      string            `json:"domain"`
}

type LinkFilter struct {
	LinkWithoutScheme string
}

type TelegramUpdate struct {
	ID   uint64 `json:"id"`
	Data tgbotapi.Update
}

type SearchRequest struct {
	FullText string
}

type SearchResponse struct {
	Links []Link `json:"links"`
}
