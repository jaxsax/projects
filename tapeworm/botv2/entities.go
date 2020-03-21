package botv2

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type HandleEntitiesResponse struct {
	Parsed []string
}

func HandleEntities(text string, entities *[]tgbotapi.MessageEntity) *HandleEntitiesResponse {
	// Decide how to do contextual logging, in the caller of this function, we've already
	// defined a context aware logger but passing a logger instance all over the place
	// seems iffy

	entitiesValue := *entities
	urlsToParse := make([]string, 0, len(entitiesValue))
	for _, entity := range entitiesValue {
		urlsToParse = append(urlsToParse, text[entity.Offset:entity.Offset+entity.Length])
	}

	return &HandleEntitiesResponse{
		Parsed: urlsToParse,
	}
}
