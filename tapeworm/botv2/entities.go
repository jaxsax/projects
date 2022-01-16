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
	return s == "readme.md"
}

func HandleEntities(entities *[]tgbotapi.MessageEntity) *HandleEntitiesResponse {
	// Decide how to do contextual logging, in the caller of this function, we've already
	// defined a context aware logger but passing a logger instance all over the place
	// seems iffy

	entitiesValue := *entities
	urlsToParse := make([]string, 0, len(entitiesValue))
	for _, entity := range entitiesValue {
		if entity.URL == "" {
			continue
		}

		url := entity.URL
		if ignoreURL(url) {
			continue
		}

		urlsToParse = append(urlsToParse, string(url))
	}

	return &HandleEntitiesResponse{
		Parsed: urlsToParse,
	}
}
