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
	Accepts(*url.URL) bool
	Provide(*url.URL) (*EnhancedLink, error)
}

var StrategyList = []Strategy{
	&Youtube{},
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
		if strategy.Accepts(url) {
			logr.FromContextOrDiscard(ctx).V(0).Info(
				"using strategy",
				"strategy_name", strategy.Name(),
				"url", url,
				"host", url.Hostname(),
			)
			return strategy.Provide(url)
		}
	}

	return nil, fmt.Errorf("no acceptable strategy")
}
