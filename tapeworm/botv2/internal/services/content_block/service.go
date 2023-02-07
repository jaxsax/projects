package contentblock

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/db"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/errors"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

type Strategy interface {
	IsBlocked(data string) (bool, bool)
	Priority() int
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

type byPriority []Strategy

func (s byPriority) Len() int {
	return len(s)
}

func (s byPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byPriority) Less(i, j int) bool {
	return s[i].Priority() < s[j].Priority()
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
				fn: func(s string) (bool, bool) {
					return strings.Contains(s, st.Content), false
				},
			})
		case "substring_case_insensitive":
			strategyFuncs = append(strategyFuncs, &funcStrategy{
				strategyName: fmt.Sprintf("substring_case_insensitive(%s)", strings.ToLower(st.Content)),
				fn: func(s string) (bool, bool) {
					return strings.Contains(strings.ToLower(s), strings.ToLower(st.Content)), false
				},
			})
		case "substring_exclude":
			strategyFuncs = append(strategyFuncs, &funcStrategy{
				strategyName: fmt.Sprintf("substring_exclude(%s)", st.Content),
				priority:     -100,
				fn: func(s string) (bool, bool) {
					return strings.Contains(strings.ToLower(s), strings.ToLower(st.Content)), true
				},
			})
		}
	}

	s.m.Lock()
	s.strategies = strategyFuncs
	sort.Sort(byPriority(s.strategies))
	s.m.Unlock()

	return nil
}

type funcStrategy struct {
	strategyName string
	priority     int
	fn           func(string) (bool, bool)
}

func (s *funcStrategy) Name() string {
	return s.strategyName
}

func (s *funcStrategy) Priority() int {
	return s.priority
}

func (s *funcStrategy) IsBlocked(data string) (bool, bool) {
	return s.fn(data)
}

func (s *Service) IsAllowed(ctx context.Context, data string, source string) error {
	s.m.RLock()
	defer s.m.RUnlock()

	for _, st := range s.strategies {
		matched, stopChain := st.IsBlocked(data)
		logging.FromContext(ctx).Info(
			"blocklist check",
			"name", st.Name(),
			"source", source,
			"data", data,
			"matched", matched,
			"stopChain", stopChain,
		)

		if matched && stopChain {
			break
		}

		if matched {
			return errors.ErrBlockedContent
		}
	}

	return nil
}
