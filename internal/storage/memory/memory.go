package memory

import (
	"sync"

	"github.com/blankInPajamas/SentinelGo/internal/logs"
	"github.com/blankInPajamas/SentinelGo/internal/storage"
)

type InMemoryStorage struct {
	mu sync.RWMutex
	events []logs.Event
}

func New() *InMemoryStorage {
	return &InMemoryStorage{
		events: make([]logs.Event, 0),
	}
}

func (s *InMemoryStorage) Save(event logs.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events = append(s.events, event)
	return nil
}

func (s *InMemoryStorage) Query(filter storage.QueryFilter) ([]logs.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var results []logs.Event

	for _, e := range s.events {
		if !filter.StartTime.IsZero() && e.Timestamp.Before(filter.StartTime) {
			continue
		}

		if filter.SourceIP != "" && e.SourceIP != filter.SourceIP {
			continue
		}

		if filter.User != "" && e.User != filter.User {
			continue
		}

		if filter.Outcome != "" && e.Outcome != filter.Outcome {
			continue
		}
		if filter.EventType != "" && e.EventType != filter.EventType {
			continue
		}

		results = append(results, e)
	}

	if filter.Offset > 0 && filter.Offset < len(results) {
		results = results[filter.Offset:]
	} else if filter.Offset >= len(results) {
		results = []logs.Event{}
	}

	if filter.Limit > 0 && len(results) > filter.Limit {
		results = results[:filter.Limit]
	}

	return results, nil
}

func (s *InMemoryStorage) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.events)
}
	
func (s * InMemoryStorage) Close() error {
	return nil
}