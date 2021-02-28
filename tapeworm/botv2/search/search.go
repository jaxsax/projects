package search

import (
	"fmt"
	"strconv"
	"strings"

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
	linksRepository   links.Repository
	byLinkTitleBucket collectionBucket
}

var _ LinkSearcher = &SonicLinkSearcher{}

func NewSonicLinkSearcher(searchable sonic.Searchable, linksRepository links.Repository) *SonicLinkSearcher {
	return &SonicLinkSearcher{
		s:               searchable,
		linksRepository: linksRepository,
		byLinkTitleBucket: collectionBucket{
			Collection: "links",
			Bucket:     "title",
		},
	}
}

func (searcher *SonicLinkSearcher) Search(query string, limit int, offset int) ([]links.Link, error) {
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
