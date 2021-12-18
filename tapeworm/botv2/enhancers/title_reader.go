package enhancers

import (
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

func pageTitle(n *html.Node) string {
	var title string
	if n.Type == html.ElementNode && n.Data == "title" {
		return n.FirstChild.Data
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		title = pageTitle(c)
		if title != "" {
			break
		}
	}
	return title
}

var successiveSpaces = regexp.MustCompile(`\s+`)

func ReadTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	title := pageTitle(doc)
	if title == "" {
		return "", fmt.Errorf("not found")
	}

	title = strings.TrimSpace(title)
	title = successiveSpaces.ReplaceAllLiteralString(title, " ")
	return title, nil
}
