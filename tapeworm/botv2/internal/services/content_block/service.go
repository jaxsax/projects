package contentblock

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/errors"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

type Strategy interface {
	IsBlocked(data string) bool
	Name() string
}

type Service struct {
	db *db.Store

	m          sync.RWMutex
	strategies []Strategy
}

func New(db *db.Store) *Service {
	return &Service{
		db:         db,
		strategies: make([]Strategy, 0),
	}
}

func (s *Service) Start() error {
	if err := s.refresh(); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(time.Minute)
		for range ticker.C {
			if err := s.refresh(); err != nil {
				logging.FromContext(context.Background()).Error(err, "reload content blocklist strategy failed")
			}
		}
	}()

	return nil
}

func (s *Service) refresh() error {
	strategies, err := s.db.ListBLocklistStrategies(context.Background())
	if err != nil {
		return err
	}

	strategyFuncs := make([]Strategy, 0)
	for _, st := range strategies {
		st := st
		logging.FromContext(context.Background()).Info("loading strategy func", "param", st)
		switch st.Strategy {
		case "substring":
			strategyFuncs = append(strategyFuncs, &funcStrategy{
				strategyName: fmt.Sprintf("substring(%s)", st.Content),
				fn: func(s string) bool {
					return strings.Contains(s, st.Content)
				},
			})
		case "substring_case_insensitive":
			strategyFuncs = append(strategyFuncs, &funcStrategy{
				strategyName: fmt.Sprintf("substring_case_insensitive(%s)", strings.ToLower(st.Content)),
				fn: func(s string) bool {
					return strings.Contains(strings.ToLower(s), strings.ToLower(st.Content))
				},
			})
		}
	}

	s.m.Lock()
	s.strategies = strategyFuncs
	s.m.Unlock()

	return nil
}

type funcStrategy struct {
	strategyName string
	fn           func(string) bool
}

func (s *funcStrategy) Name() string {
	return s.strategyName
}

func (s *funcStrategy) IsBlocked(data string) bool {
	return s.fn(data)
}

func (s *Service) IsAllowed(ctx context.Context, data string, source string) error {
	s.m.RLock()
	defer s.m.RUnlock()

	for _, st := range s.strategies {
		if st.IsBlocked(data) {
			logging.FromContext(ctx).Info(
				"content blocked by strategy",
				"name", st.Name(),
				"source", source,
				"data", data,
			)

			return errors.ErrBlockedContent
		}
	}

	return nil
}
