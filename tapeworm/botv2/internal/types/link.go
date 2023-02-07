package types

import (
	"encoding/json"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Link struct {
	ID           uint64            `json:"id"`
	Link         string            `json:"link"`
	Title        string            `json:"title"`
	CreatedAt    uint64            `json:"created_ts"`
	CreatedByID  uint64            `json:"created_by"`
	ExtraData    map[string]string `json:"extra_data"`
	Domain       string            `json:"domain"`
	Path         string            `json:"path"`
	DeletedAt    *time.Time        `json:"deleted_at,omitempty"`
	Labels       []*Label          `json:"labels"`
	Dimensions   []*Dimension      `json:"dimensions"`
	DimCollected bool              `json:"-"`
}

type Label struct {
	Name string `json:"name"`
}

type Dimension struct {
	Kind DimensionKind   `json:"kind"`
	Data json.RawMessage `json:"data"`
}

type LinkFilter struct {
	LinkWithoutScheme string
	Domain            string
	PageNumber        int
	Limit             int

	TitleSearch  string
	UniqueLink   bool
	DimCollected *bool
}

type TelegramUpdate struct {
	ID   uint64 `json:"id"`
	Data tgbotapi.Update
}
