package types

import (
	"time"
)

type Link struct {
	Link        string            `json:"link"`
	Title       string            `json:"title"`
	CreatedAt   time.Time         `json:"created_at"`
	CreatedByID uint64            `json:"created_by"`
	ExtraData   map[string]string `json:"extra_data"`
	DeletedAt   *time.Time        `json:"deleted_at,omitempty"`
}
