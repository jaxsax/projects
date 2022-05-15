package botv2

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/scylladb/go-set/strset"
)

type HandleEntitiesResponse struct {
	Parsed []string
}

func ignoreURL(s string) bool {
	s = strings.ToLower(s)
	return s == "readme.md" || s == ""
}

func HandleEntities(msg string, entities []tgbotapi.MessageEntity) *HandleEntitiesResponse {
	runeText := []rune(msg)
	uniqueUrls := strset.New()
	for _, entity := range entities {
		if !entity.IsTextLink() && !entity.IsURL() {
			continue
		}

		url := string(runeText[entity.Offset : entity.Offset+entity.Length])
		if ignoreURL(url) {
			continue
		}

		uniqueUrls.Add(url)
	}

	return &HandleEntitiesResponse{
		Parsed: uniqueUrls.List(),
	}
}
