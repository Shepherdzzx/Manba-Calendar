package executor

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type jsonStoreData struct {
	Events map[string]Event `json:"events"`
}

type JSONFileStore struct {
	mu      sync.RWMutex
	path    string
	events  map[string]Event
	nowFunc func() time.Time
}

func NewJSONFileStore(path string, nowFunc func() time.Time) (*JSONFileStore, error) {
	if nowFunc == nil {
		nowFunc = time.Now
	}
	store := &JSONFileStore{
		path:    path,
		events:  make(map[string]Event),
		nowFunc: nowFunc,
	}
	if err := store.load(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *JSONFileStore) load() error {
	data, err := os.ReadFile(s.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if len(data) == 0 {
		return nil
	}
	var payload jsonStoreData
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("failed to decode store file: %w", err)
	}
	if payload.Events == nil {
		payload.Events = make(map[string]Event)
	}
	s.events = payload.Events
	return nil
}

func (s *JSONFileStore) save() error {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}
	payload := jsonStoreData{Events: s.events}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return err
	}
	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmpPath, s.path)
}

func (s *JSONFileStore) Create(e Event) (Event, error) {
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
	if err := s.save(); err != nil {
		delete(s.events, id)
		return Event{}, err
	}
	return e, nil
}

func (s *JSONFileStore) Query(opts QueryOptions) ([]Event, error) {
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
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}

func (s *JSONFileStore) Delete(id string) (Event, error) {
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
	if err := s.save(); err != nil {
		return Event{}, err
	}
	return e, nil
}

func (s *JSONFileStore) ListAll() ([]Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]Event, 0, len(s.events))
	for _, e := range s.events {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].CreatedAt.Before(out[j].CreatedAt)
	})
	return out, nil
}
