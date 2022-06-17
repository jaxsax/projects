package enhancers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DefaultStrategy struct{}

func (s *DefaultStrategy) Name() string {
	return "default"
}

func (s *DefaultStrategy) Accepts(u *url.URL) bool {
	return true
}

func (s *DefaultStrategy) Provide(url *url.URL) (*EnhancedLink, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
	)

	res, err := httpClient.Do(req)
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
			Link:  url.String(),
			Title: title,
		}, nil
	case "application/pdf":
		return &EnhancedLink{
			Link:  url.String(),
			Title: url.String(),
		}, nil
	}

	return nil, fmt.Errorf("unimplemented type %v", contentType)
}