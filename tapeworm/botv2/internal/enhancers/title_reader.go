package enhancers

import (
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ReadTitle(body io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return "", err
	}

	var title = ""
	doc.Find("title").Each(func(i int, s *goquery.Selection) {
		if title == "" {
			title = s.Text()
		}
	})

	doc.Find("meta").Each(func(i int, s *goquery.Selection) {
		propertyName, ok := s.Attr("property")
		if !ok {
			return
		}

		if propertyName != "og:title" {
			return
		}

		propertyValue, ok := s.Attr("content")
		if !ok {
			return
		}

		if title == "" {
			title = propertyValue
		}
	})

	title = strings.TrimSpace(title)

	return title, nil
}
