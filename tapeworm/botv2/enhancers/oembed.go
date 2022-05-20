package enhancers

import (
	"fmt"
	"io"
	"net/url"

	"golang.org/x/net/html"
)

// https://oembed.com/#section4
type OEmbedStrategy struct{}

func (s *OEmbedStrategy) Name() string {
	return "oembed"
}

// <link rel="alternate" type="application/json+oembed" href="https://www.youtube.com/oembed?format=json&amp;url=https%3A%2F%2Fwww.youtube.com%2Fwatch%3Fv%3DwUB8l1Fz0mA" title="HOW TO DECIDE by Annie Duke | Core Message">
func (s *OEmbedStrategy) Accepts(u *url.URL) bool {
	return true
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

func (s *OEmbedStrategy) Provide(u *url.URL) (*EnhancedLink, error) {
	resp, err := httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	tk := html.NewTokenizer(resp.Body)
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
	return nil, fmt.Errorf("no data")
}
