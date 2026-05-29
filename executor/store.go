package executor

import (
	"fmt"
	"sync"
	"time"
)

type Store interface {
	Create(Event) (Event, error)
	Query(QueryOptions) ([]Event, error)
	Delete(id string) (Event, error)
	ListAll() ([]Event, error)
}

type InMemoryStore struct {
	mu      sync.RWMutex
	events  map[string]Event
	nowFunc func() time.Time
}

func NewInMemoryStore(nowFunc func() time.Time) *InMemoryStore {
	if nowFunc == nil {
		nowFunc = time.Now
	}
	return &InMemoryStore{
		events:  make(map[string]Event),
		nowFunc: nowFunc,
	}
}

func (s *InMemoryStore) Create(e Event) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	id, err := generateEventID()
	if err != nil {
		return Event{}, err
	}
	e.ID = id
	e.Status = "active"
	e.CreatedAt = s.nowFunc()
	s.events[id] = e
	return e, nil
}

func (s *InMemoryStore) Query(opts QueryOptions) ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []Event
	for _, e := range s.events {
		if e.Status != "active" {
			continue
		}
		if opts.Date != "" && e.Date != opts.Date {
			continue
		}
		if opts.Time != "" && e.Time != opts.Time {
			continue
		}
		if opts.Title != "" && !contains(e.Title, opts.Title) {
			continue
		}
		out = append(out, e)
	}
	return out, nil
}

func (s *InMemoryStore) Delete(id string) (Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.events[id]
	if !ok {
		return Event{}, fmt.Errorf("event not found")
	}
	now := s.nowFunc()
	e.Status = "deleted"
	e.DeletedAt = &now
	s.events[id] = e
	return e, nil
}

func (s *InMemoryStore) ListAll() ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, 0, len(s.events))
	for _, e := range s.events {
		out = append(out, e)
	}
	return out, nil
}

func contains(haystack, needle string) bool {
	return needle == "" || (haystack != "" && len(needle) <= len(haystack) && (haystack == needle || containsFold(haystack, needle)))
}

func containsFold(haystack, needle string) bool {
	return len(needle) == 0 || (len(haystack) >= len(needle) && indexFold(haystack, needle) >= 0)
}

func indexFold(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
