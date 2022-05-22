package enhancers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
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

func EnhanceLinkWithContext(ctx context.Context, link string) (*EnhancedLink, error) {
	url, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	removeUTMParameters(url)

	for _, strategy := range StrategyList {
		lg := logr.FromContextOrDiscard(ctx).WithValues("strategy", strategy.Name(), "url", url)
		if strategy.Accepts(url) {
			lg.Info("accepted")

			e, err := strategy.Provide(url)
			if err != nil {
				lg.Error(err, "strategy failed to provide")
				continue
			}

			e.Title = strings.TrimSpace(e.Title)

			lg.Info("strategy provided", "info", e)
			return e, nil
		}
	}

	return nil, fmt.Errorf("no acceptable strategy")
}
