package telegrampoller

import (
	"unicode/utf16"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func NewReplyToMessage(text string, originalMessage *tgbotapi.Message) tgbotapi.MessageConfig {
	m := tgbotapi.NewMessage(originalMessage.Chat.ID, text)
	m.ReplyToMessageID = originalMessage.MessageID

	return m
}

func ExtractURL(text string, offset, length int) string {
	// https://github.com/go-telegram-bot-api/telegram-bot-api/issues/231
	encodedString := utf16.Encode([]rune(text))
	runeString := utf16.Decode(encodedString[offset : offset+length])

	return string(runeString)
}
