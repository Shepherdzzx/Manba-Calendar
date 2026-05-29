package executor

import (
	"testing"
	"time"
)

func TestInMemoryStore_CreateQueryDelete(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })

	created, err := store.Create(Event{Title: "和张三开会", Date: "2026-05-30", Time: "15:00"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if created.ID == "" || created.Status != "active" {
		t.Fatalf("unexpected created event: %+v", created)
	}

	found, err := store.Query(QueryOptions{Date: "2026-05-30"})
	if err != nil || len(found) != 1 {
		t.Fatalf("unexpected query result: %+v, err=%v", found, err)
	}

	deleted, err := store.Delete(created.ID)
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if deleted.Status != "deleted" || deleted.DeletedAt == nil {
		t.Fatalf("unexpected deleted event: %+v", deleted)
	}

	foundAfterDelete, err := store.Query(QueryOptions{Date: "2026-05-30"})
	if err != nil {
		t.Fatalf("unexpected query error after delete: %v", err)
	}
	if len(foundAfterDelete) != 0 {
		t.Fatalf("expected 0 active events after delete, got %+v", foundAfterDelete)
	}
}

func TestInMemoryStore_AllowsMultipleEventsSameDateTime(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })

	first, err := store.Create(Event{Title: "要写作业", Date: "2026-06-06", Time: "18:00"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	second, err := store.Create(Event{Title: "去超市", Date: "2026-06-06", Time: "18:00"})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if first.ID == second.ID {
		t.Fatalf("expected unique IDs, got same ID %s", first.ID)
	}

	found, err := store.Query(QueryOptions{Date: "2026-06-06", Time: "18:00"})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if len(found) != 2 {
		t.Fatalf("expected 2 events, got %+v", found)
	}
}
