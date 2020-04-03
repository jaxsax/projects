package enhancers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type EnhancedLink struct {
	Original string
	Link     string
	Title    string
}

func EnhanceLink(link string) (*EnhancedLink, error) {
	url, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	res, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(body)
	switch contentType {
	case "text/xml; charset=utf-8":
		fallthrough
	case "text/html; charset=utf-8":
		title, err := ReadTitle(bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("read title: %w", err)
		}

		return &EnhancedLink{
			Original: link,
			Link:     url.String(),
			Title:    title,
		}, nil
	}

	return nil, fmt.Errorf("not imeplemented")
}
