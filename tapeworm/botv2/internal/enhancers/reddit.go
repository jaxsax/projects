package enhancers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/pkg/errors"
)

type Reddit struct{}

var _ Strategy = (*Reddit)(nil)

func (r *Reddit) Name() string {
	return "reddit"
}

var (
	// https://old.reddit.com/r/programming/comments/17ooxwe/interruptions_cost_23_minutes_15_seconds_right/?share_id=OPTVwmKJ9uA465C9H1AlX
	commentThreadPath = regexp.MustCompile(`r/(?P<subreddit>.*)/comments/(?P<postId>.*)/(?P<postslug>.*)/$`)

	// https://www.reddit.com/r/programming/s/iJI7LKOI3f
	shareThreadPath = regexp.MustCompile(`r/(?P<subreddit>.*)/s/(?P<shareId>.*)$`)
)

func (r *Reddit) Provide(u *url.URL) (*EnhancedLink, error) {
	if !strings.Contains(u.Hostname(), "reddit.com") {
		return nil, fmt.Errorf("not reddit domain")
	}

	commentMatches := commentThreadPath.FindStringSubmatch(u.Path)
	shareMatches := shareThreadPath.FindStringSubmatch(u.Path)
	logging.FromContext(context.TODO()).Info("debug", "path", u.Path, "comment", commentMatches, "shares", shareMatches)

	if len(commentMatches) > 0 {
		redditRemap := map[string]string{
			"www.reddit.com": "old.reddit.com",
			"reddit.com":     "old.reddit.com",
		}

		rm, hasRemap := redditRemap[u.Host]
		if hasRemap {
			u.Host = rm
		}

		logging.FromContext(context.TODO()).Info("making request", "url", u.String())

		req, err := http.NewRequest(http.MethodGet, u.String(), nil)
		if err != nil {
			return nil, errors.Wrap(err, "build request for comment thread")
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "send http request")
		}

		title, err := ReadTitle(resp.Body)
		if err != nil {
			return nil, errors.Wrap(err, "read title")
		}

		return &EnhancedLink{
			Original: u.String(),
			Link:     u.String(),
			Title:    title,
		}, nil
	}

	if len(shareMatches) > 0 {
		req, err := http.NewRequest(http.MethodHead, u.String(), nil)
		if err != nil {
			return nil, errors.Wrap(err, "build request for share match")
		}

		resp, err := httpClientNoRedirect.Do(req)
		if err != nil {
			return nil, errors.Wrap(err, "send http request")
		}

		logging.FromContext(context.TODO()).Info("share url", "url", u.String())

		newURL, err := resp.Location()
		if err != nil {
			if errors.Is(err, http.ErrNoLocation) {
				return nil, fmt.Errorf("no location redirect")
			}

			return nil, errors.Wrap(err, "get location")
		}

		removeUTMParameters(newURL)
		el, err := r.Provide(newURL)
		if err != nil {
			return nil, errors.Wrap(err, "comment match")
		}

		return el, nil
	}

	return nil, fmt.Errorf("no match")
}
