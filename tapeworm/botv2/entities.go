package botv2

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/scylladb/go-set/strset"
)

type HandleEntitiesResponse struct {
	Parsed []string
}

func ignoreURL(s string) bool {
	s = strings.ToLower(s)
	return s == "readme.md"
}

func HandleEntities(msg string, entities *[]tgbotapi.MessageEntity) *HandleEntitiesResponse {
	runeText := []rune(msg)
	entitiesValue := *entities
	uniqueUrls := strset.New()
	for _, entity := range entitiesValue {
		if !entity.IsTextLink() && !entity.IsUrl() {
			continue
		}

		url := string(runeText[entity.Offset : entity.Offset+entity.Length])
		if url == "" {
			continue
		}

		if ignoreURL(url) {
			continue
		}

		uniqueUrls.Add(url)
	}

	return &HandleEntitiesResponse{
		Parsed: uniqueUrls.List(),
	}
}
