package dimcollector

import (
	"context"
	"fmt"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/dimension"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Service struct {
	collectors []dimension.Collector
	store      *db.Store

	shutdown         chan struct{}
	shutdownComplete chan struct{}
}

func New(
	store *db.Store,
	hnCollector *dimension.HNCollector,
) *Service {
	return &Service{
		collectors: []dimension.Collector{
			hnCollector,
		},
		store:            store,
		shutdown:         make(chan struct{}),
		shutdownComplete: make(chan struct{}),
	}
}

func (s *Service) Start(ctx context.Context) error {
	ticker := time.NewTicker(time.Minute)
	go func() {
		for {
			select {
			case <-s.shutdown:
				goto shutdown
			case <-ticker.C:
			}

			ctx := context.Background()
			t := false
			links, err := s.store.ListLinksWithFilter(ctx, &types.LinkFilter{
				Limit:        100,
				DimCollected: &t,
			})

			if err != nil {
				logging.FromContext(ctx).Error(err, "error listing links")
				continue
			}

			for _, link := range links {
				if err := s.PopulateDimensions(ctx, link); err != nil {
					logging.FromContext(ctx).Error(err, "error populating dims", "link", link)
					continue
				}
			}
		}

	shutdown:
		s.shutdownComplete <- struct{}{}
	}()

	return nil
}

func (s *Service) Stop(ctx context.Context) error {
	close(s.shutdown)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.shutdownComplete:
	}

	return nil
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
