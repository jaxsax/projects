package enhancers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/errors"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

type EnhancedLink struct {
	Original string
	Link     string
	Title    string
}

type Strategy interface {
	Name() string
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

var remapHostMap = map[string]string{
	"www.reddit.com": "old.reddit.com",
	"reddit.com":     "old.reddit.com",
}

var successiveSpaces = regexp.MustCompile(`\s+`)

func EnhanceLinkWithContext(ctx context.Context, link string, p *db.Store) (*EnhancedLink, error) {
	if !strings.HasPrefix(link, "https://") && !strings.HasPrefix(link, "http://") {
		link = fmt.Sprintf("http://%s", link)
	}

	providedURL, err := url.Parse(link)
	if err != nil {
		return nil, err
	}

	removeUTMParameters(providedURL)

	if providedURL.Hostname() == "" {
		return nil, errors.ErrInvalidDomain
	}

	urlToRetrieveFrom, err := url.Parse(providedURL.String())
	if err != nil {
		return nil, err
	}

	rm, hasRemap := remapHostMap[urlToRetrieveFrom.Host]
	if hasRemap {
		urlToRetrieveFrom.Host = rm
	}

	for _, strategy := range StrategyList {
		lg := logging.FromContext(ctx).WithValues("strategy", strategy.Name(), "url", urlToRetrieveFrom)

		lg.Info("trying strategy")

		e, err := strategy.Provide(urlToRetrieveFrom)
		if err != nil {
			lg.Error(err, "strategy failed to provide")
			continue
		}

		if e == nil {
			lg.Info("nil object")
			continue
		}

		title := strings.TrimSpace(e.Title)
		title = successiveSpaces.ReplaceAllLiteralString(title, " ")

		e.Title = title
		e.Link = providedURL.String()

		lg.Info("strategy provided", "info", e)
		return e, nil
	}

	return nil, fmt.Errorf("no acceptable strategy")
}
