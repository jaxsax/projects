package botv2

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewTelegramBotAPI(token string) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return bot, nil
}
