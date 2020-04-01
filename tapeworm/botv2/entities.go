package botv2

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HandleEntitiesResponse struct {
	Parsed []string
}

func ignoreURL(s string) bool {
	s = strings.ToLower(s)
	if s == "readme.md" {
		return true
	}
	return false
}

func HandleEntities(text string, entities *[]tgbotapi.MessageEntity) *HandleEntitiesResponse {
	// Decide how to do contextual logging, in the caller of this function, we've already
	// defined a context aware logger but passing a logger instance all over the place
	// seems iffy

	runeText := []rune(text)
	entitiesValue := *entities
	urlsToParse := make([]string, 0, len(entitiesValue))
	for _, entity := range entitiesValue {
		if !entity.IsUrl() {
			continue
		}

		url := string(runeText[entity.Offset : entity.Offset+entity.Length])
		if ignoreURL(url) {
			continue
		}

		urlsToParse = append(urlsToParse, string(url))
	}

	return &HandleEntitiesResponse{
		Parsed: urlsToParse,
	}
}
