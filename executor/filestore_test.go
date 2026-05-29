package executor

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestJSONFileStore_PersistsAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)

	store, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	created, err := store.Create(Event{Title: "要写作业", Date: "2026-05-29"})
	if err != nil {
		t.Fatalf("unexpected store create error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected store file to exist: %v", err)
	}

	reopened, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected reopen error: %v", err)
	}
	found, err := reopened.Query(QueryOptions{Date: "2026-05-29"})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if len(found) != 1 || found[0].ID != created.ID {
		t.Fatalf("unexpected persisted events: %+v", found)
	}
}

func TestJSONFileStore_DeletePersists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)

	store, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	created, err := store.Create(Event{Title: "开会", Date: "2026-05-30", Time: "15:00"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	_, err = store.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}

	reopened, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected reopen error: %v", err)
	}
	found, err := reopened.Query(QueryOptions{Date: "2026-05-30"})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if len(found) != 0 {
		t.Fatalf("expected deleted event to stay filtered out, got %+v", found)
	}
	all, err := reopened.ListAll()
	if err != nil {
		t.Fatalf("unexpected list error: %v", err)
	}
	if len(all) != 1 || all[0].Status != "deleted" || all[0].DeletedAt == nil {
		t.Fatalf("unexpected all-events view: %+v", all)
	}
}

func TestJSONFileStore_AllowsMultipleEventsSameDateAcrossReopen(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)

	store, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	first, err := store.Create(Event{Title: "要写作业", Date: "2026-06-06"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	second, err := store.Create(Event{Title: "去超市买东西", Date: "2026-06-06"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if first.ID == second.ID {
		t.Fatalf("expected unique IDs, got same ID %s", first.ID)
	}

	reopened, err := NewJSONFileStore(path, func() time.Time { return now })
	if err != nil {
		t.Fatalf("unexpected reopen error: %v", err)
	}
	third, err := reopened.Create(Event{Title: "看电影", Date: "2026-06-06"})
	if err != nil {
		t.Fatalf("unexpected create error after reopen: %v", err)
	}
	if third.ID == first.ID || third.ID == second.ID {
		t.Fatalf("expected unique ID after reopen, got duplicate %s", third.ID)
	}
	found, err := reopened.Query(QueryOptions{Date: "2026-06-06"})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if len(found) != 3 {
		t.Fatalf("expected 3 events, got %+v", found)
	}
}
