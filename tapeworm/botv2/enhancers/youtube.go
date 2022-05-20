package enhancers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/scylladb/go-set/strset"
)

type Youtube struct{}

type youtubeoEmbedResponse struct {
	Title           string `json:"title"`
	AuthorName      string `json:"author_name"`
	AuthorURL       string `json:"author_url"`
	Type            string `json:"video"`
	Height          int    `json:"height"`
	Width           int    `json:"width"`
	Version         string `json:"version"`
	ProviderName    string `json:"provider_name"`
	ProviderURL     string `json:"provider_url"`
	ThumbnailHeight int    `json:"thumbnail_height"`
	ThumbnailWeight int    `json:"thumbnail_weight"`
	ThumbnailURL    string `json:"thumbnail_url"`
	HTML            string `json:"html"`
}

var ytHostSet = strset.New(
	"youtube.com",
	"www.youtube.com",
)

func (s *Youtube) Name() string {
	return "youtube"
}

func (s *Youtube) Accepts(u *url.URL) bool {
	return ytHostSet.Has(u.Hostname()) && u.Path == "/watch" && u.Query().Get("v") != ""
}

func (s *Youtube) Provide(u *url.URL) (*EnhancedLink, error) {
	videoID := u.Query().Get("v")
	if videoID == "" {
		return nil, fmt.Errorf("missing video id")
	}

	req, err := http.NewRequest(http.MethodGet, "https://www.youtube.com/oembed", nil)
	if err != nil {
		return nil, err
	}

	q := url.Values{}
	q.Set("format", "json")
	q.Set("url", u.String())
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	var embedResponse youtubeoEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&embedResponse); err != nil {
		return nil, err
	}

	return &EnhancedLink{
		Link:  u.String(),
		Title: embedResponse.Title,
	}, nil
}
