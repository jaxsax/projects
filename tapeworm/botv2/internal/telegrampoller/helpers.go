package telegrampoller

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func NewReplyToMessage(text string, originalMessage *tgbotapi.Message) tgbotapi.MessageConfig {
	m := tgbotapi.NewMessage(originalMessage.Chat.ID, text)
	m.ReplyToMessageID = originalMessage.MessageID

	return m
}
