package enhancers

import (
	"fmt"
	"io"

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

func ReadTitle(body io.Reader) (string, error) {
	doc, err := html.Parse(body)
	if err != nil {
		return "", err
	}

	title := pageTitle(doc)
	if title == "" {
		return "", fmt.Errorf("not found")
	}
	return title, nil
}
