package enhancers

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type DefaultStrategy struct{}

func (s *DefaultStrategy) Name() string {
	return "default"
}

func (s *DefaultStrategy) Provide(url *url.URL) (*EnhancedLink, error) {
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36",
	)

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = res.Body.Close()
	}()

	bodyReader := io.LimitReader(res.Body, 25*1024*1024)
	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("invalid status code: %d", res.StatusCode)
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
