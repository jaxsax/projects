package enhancers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

// https://oembed.com/#section4
type OEmbedStrategy struct {
}

func (s *OEmbedStrategy) Name() string {
	return "oembed"
}

func (s *OEmbedStrategy) Provide(u *url.URL) (*EnhancedLink, error) {
	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    respStr := string(respBody)
    if !strings.Contains(respStr, "application/json+oembed") {
        return nil, fmt.Errorf("response did not contain oembed tags")
    }

	tk := html.NewTokenizer(bytes.NewBuffer(respBody))
	for {
		ntk := tk.Next()
		switch {
		case ntk == html.ErrorToken:
			goto done
		case ntk == html.StartTagToken:
			t := tk.Token()
			if t.Data != "link" {
				continue
			}

			data := getOEmbedData(t.Attr)
			if data != nil {
				return &EnhancedLink{
					Link:  u.String(),
					Title: data.Title,
				}, nil
			}
		}
	}

done:
	return nil, nil
}

type oembedData struct {
	Type  string
	Title string
}

func getOEmbedData(attrs []html.Attribute) *oembedData {
	var data oembedData
	for _, attr := range attrs {
		if attr.Key == "type" && attr.Val == "application/json+oembed" {
			data.Type = attr.Val
		}

		if attr.Key == "title" {
			data.Title = attr.Val
		}
	}

	if data.Type == "" || data.Title == "" {
		return nil
	}

	return &data
}
