package updates

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

// Update is the domain model
type Update struct {
	ID   int64
	Data *tgbotapi.Update
}

// Repository provides access to an updates store
type Repository interface {
	Create(update Update) error
}
