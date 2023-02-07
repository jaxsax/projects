package dimcollector

import (
	"context"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/dimension"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Service struct {
	collectors []dimension.Collector
	store      *db.Store
}

func New(
	store *db.Store,
	hnCollector *dimension.HNCollector,
) *Service {
	return &Service{
		collectors: []dimension.Collector{
			hnCollector,
		},
		store: store,
	}
}

// todo: omit if the link is already processed according to db records since this can be run in a reconcile method
// not sure if this should be done here though
func (s *Service) PopulateDimensions(ctx context.Context, link *types.Link) error {
	dims := make([]*types.Dimension, 0)
	for _, c := range s.collectors {
		dim, err := c.Collect(link)
		if err != nil {
			return fmt.Errorf("collect dimension: %w", err)
		}

		dims = append(dims, dim...)
	}

	if err := s.store.UpdateLinkDimensions(ctx, link, dims); err != nil {
		return err
	}

	return nil
}
