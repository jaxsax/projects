package enhancers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/go-logr/logr"
)

type EnhancedLink struct {
	Original string
	Link     string
	Title    string
}

type Strategy interface {
	Name() string
	Accepts(u *url.URL) bool
	Provide(u *url.URL) (*EnhancedLink, error)
}

var StrategyList = []Strategy{
	&OEmbedStrategy{},
	&DefaultStrategy{},
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func removeUTMParameters(url *url.URL) {
	values := url.Query()
	for k := range values {
		if strings.HasPrefix(k, "utm_") {
			values.Del(k)
		}
	}

	url.RawQuery = values.Encode()
}

func EnhanceLink(link string) (*EnhancedLink, error) {
	return EnhanceLinkWithContext(context.Background(), link)
}

var remapHostMap = map[string]string{
	"www.reddit.com": "old.reddit.com",
	"reddit.com":     "old.reddit.com",
}

var successiveSpaces = regexp.MustCompile(`\s+`)

func EnhanceLinkWithContext(ctx context.Context, link string) (*EnhancedLink, error) {
	providedURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if providedURL.Scheme == "" {
		providedURL.Scheme = "http"
	}

	removeUTMParameters(providedURL)

	urlToRetrieveFrom, err := url.Parse(providedURL.String())
	if err != nil {
		return nil, err
	}

	rm, hasRemap := remapHostMap[urlToRetrieveFrom.Host]
	if hasRemap {
		urlToRetrieveFrom.Host = rm
	}

	for _, strategy := range StrategyList {
		lg := logr.FromContextOrDiscard(ctx).WithValues("strategy", strategy.Name(), "url", urlToRetrieveFrom)
		if strategy.Accepts(urlToRetrieveFrom) {
			lg.Info("accepted")

			e, err := strategy.Provide(urlToRetrieveFrom)
			if err != nil {
				lg.Error(err, "strategy failed to provide")
				continue
			}

			title := strings.TrimSpace(e.Title)
			title = successiveSpaces.ReplaceAllLiteralString(title, " ")

			e.Title = title
			e.Link = providedURL.String()

			lg.Info("strategy provided", "info", e)
			return e, nil
		}
	}

	return nil, fmt.Errorf("no acceptable strategy")
}
