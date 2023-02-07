package algolia

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Algolia struct {
	host string
}

func New(host string) *Algolia {
	return &Algolia{
		host: host,
	}
}

func (a *Algolia) getURL(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	return fmt.Sprintf("%s%s", a.host, path)
}

func (a *Algolia) Search(req *SearchRequest) (*SearchResponse, error) {
	if req.Query == "" {
		return nil, fmt.Errorf("query cannot be empty")
	}

	param := url.Values{}
	param.Set("query", req.Query)

	if !req.Analytics {
		param.Set("analytics", "false")
	}

	if len(req.RestrictSearchableAttributes) > 0 {
		param.Set("restrictSearchableAttributes", strings.Join(req.RestrictSearchableAttributes, ","))
	}

	if req.TypoTolerance != "" {
		param.Set("typoTolerance", req.TypoTolerance)
	}

	resp, err := http.Get(a.getURL("/api/v1/search") + "?" + param.Encode())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var searchResponse SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResponse); err != nil {
		return nil, err
	}

	return &searchResponse, nil
}
