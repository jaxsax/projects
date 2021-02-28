package search

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"

	"github.com/expectedsh/go-sonic/sonic"
	"github.com/jaxsax/projects/tapeworm/botv2/links"
)

type LinkSearcher interface {
	Search(query string, limit int, offset int) ([]links.Link, error)
	Ingest([]links.Link) error
}

type collectionBucket struct {
	Collection string
	Bucket     string
}

type SonicLinkSearcher struct {
	s                 sonic.Searchable
	conf              internal.SonicConfig
	sonicMu           sync.Mutex
	linksRepository   links.Repository
	byLinkTitleBucket collectionBucket
	logger            *zap.Logger
}

var _ LinkSearcher = &SonicLinkSearcher{}

func NewSonicLinkSearcher(
	logger *zap.Logger,
	sonicSearcher sonic.Searchable, conf internal.SonicConfig, linksRepository links.Repository) *SonicLinkSearcher {
	return &SonicLinkSearcher{
		s:               sonicSearcher,
		conf:            conf,
		logger:          logger,
		linksRepository: linksRepository,
		byLinkTitleBucket: collectionBucket{
			Collection: "links",
			Bucket:     "title",
		},
	}
}

// Attempts to reconnect when necessary, giving up after 10 attempts
func (searcher *SonicLinkSearcher) reconnectIfNecessary() error {
	// s might be nil because this is the first connection, so we initialize the connection
	searcher.sonicMu.Lock()
	defer searcher.sonicMu.Unlock()

	var (
		maxAttemptsAllowed       = 10
		currentAttempts          = 0
		lastErr            error = nil
	)

	increaseAttempt := func() {
		currentAttempts++
		time.Sleep(100 * time.Millisecond)
	}
	for {
		if currentAttempts >= maxAttemptsAllowed {
			return fmt.Errorf("max attempts hit: %w", lastErr)
		}

		if searcher.s == nil {
			searchable, err := sonic.NewSearch(searcher.conf.Host, searcher.conf.Port, searcher.conf.Password)
			if err != nil {
				lastErr = err
				increaseAttempt()
				continue
			}

			searcher.s = searchable
		}

		// try to .Ping()
		// on ping error, try to reconnect
		// query: write tcp 172.18.0.14:48638->172.18.0.3:1491: write: broken pipe

		err := searcher.s.Ping()
		searcher.logger.Info(
			"trying to ping sonic",
			zap.Int("attempts", currentAttempts),
			zap.Int("maxAttempts", maxAttemptsAllowed),
			zap.Error(err),
		)
		if err != nil {
			lastErr = err
			searcher.s = nil
			increaseAttempt()
			continue
		}

		return nil
	}
}

func (searcher *SonicLinkSearcher) Search(query string, limit int, offset int) ([]links.Link, error) {
	err := searcher.reconnectIfNecessary()
	if err != nil {
		return []links.Link{}, err
	}

	rs, err := searcher.s.Query(
		searcher.byLinkTitleBucket.Collection,
		searcher.byLinkTitleBucket.Bucket,
		query,
		limit,
		offset,
		sonic.LangNone,
	)
	if err != nil {
		return []links.Link{}, fmt.Errorf("query: %w", err)
	}

	linkIDs := make([]int64, 0, len(rs))

	for _, obj := range rs {
		p := strings.Split(obj, ":")
		if len(p) != 2 {
			continue
		}

		if p[0] != "id" {
			continue
		}

		id, err := strconv.ParseInt(p[1], 10, 64)
		if err != nil {
			continue
		}

		linkIDs = append(linkIDs, id)
	}

	matchingLinks, err := searcher.linksRepository.ListMatchingIDs(linkIDs)
	if err != nil {
		return []links.Link{}, fmt.Errorf("list: %w", err)
	}

	return matchingLinks, nil
}

func (searcher *SonicLinkSearcher) Ingest([]links.Link) error {
	return fmt.Errorf("not implemented yet")
}
