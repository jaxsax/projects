package db

import (
	"context"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/go-logr/logr"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

func (q *Store) Search(ctx context.Context, req *types.SearchRequest) (*types.SearchResponse, error) {
	query := bleve.NewMatchQuery(req.FullText)
	sr, err := q.linkIndex.Search(bleve.NewSearchRequest(query))
	if err != nil {
		return nil, err
	}

	itemsFound := make([]types.Link, 0, len(sr.Hits))
	for _, item := range sr.Hits {
		asID, err := strconv.Atoi(item.ID)
		if err != nil {
			return nil, err
		}

		linkItem, err := q.GetLink(ctx, uint64(asID))
		if err != nil {
			return nil, err
		}

		itemsFound = append(itemsFound, *linkItem)
	}

	return &types.SearchResponse{
		Links: itemsFound,
	}, nil
}

func (q *Store) IndexAllItems(ctx context.Context) error {
	b := q.linkIndex.NewBatch()
	logr.FromContextOrDiscard(ctx).Info("listing items")
	items, err := q.ListLinks(ctx)
	if err != nil {
		return err
	}

	for _, item := range items {
		logr.FromContextOrDiscard(ctx).Info("indexing item", "id", item.ID, "title", item.Title)
		err := b.Index(strconv.Itoa(int(item.ID)), item.Title)
		if err != nil {
			return err
		}
	}

	if err := q.linkIndex.Batch(b); err != nil {
		return err
	}

	return nil
}
