package executor

import (
	"testing"
	"time"

	parserlib "github.com/Shepherdzzx/manba-alert/parser"
)

func TestExecutor_CreateQueryDeleteFlow(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	createResult, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentCreateEvent,
		Title:  "要写作业",
		Date:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if createResult.Kind != ResultCreated || createResult.Event == nil {
		t.Fatalf("unexpected create result: %+v", createResult)
	}

	queryResult, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentQueryEvents,
		Date:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if queryResult.Kind != ResultQueried || len(queryResult.Events) != 1 {
		t.Fatalf("unexpected query result: %+v", queryResult)
	}

	deleteResult, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Title:  "要写作业",
		Date:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected delete error: %v", err)
	}
	if deleteResult.Kind != ResultDeleted || deleteResult.Event == nil {
		t.Fatalf("unexpected delete result: %+v", deleteResult)
	}
}

func TestExecutor_DeleteConflictWhenMultipleMatches(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "开会", Date: "2026-05-30", Time: "15:00"})
	_, _ = store.Create(Event{Title: "开会", Date: "2026-05-31", Time: "16:00"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Title:  "开会",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultConflict || len(result.Events) != 2 {
		t.Fatalf("unexpected conflict result: %+v", result)
	}
	if result.NeedsConfirmation {
		t.Fatalf("expected no confirmation needed for title-only conflict, got %+v", result)
	}
}

func TestExecutor_DeleteDateOnlyWhenSingleMatch(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "要写作业", Date: "2026-05-29"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Date:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultDeleted || result.Event == nil {
		t.Fatalf("unexpected delete result: %+v", result)
	}
}

func TestExecutor_DeleteDateOnlyConflictWhenMultipleMatches(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "要写作业", Date: "2026-05-29"})
	_, _ = store.Create(Event{Title: "去超市", Date: "2026-05-29"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Date:   "2026-05-29",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultConflict || len(result.Events) != 2 {
		t.Fatalf("unexpected conflict result: %+v", result)
	}
	if !result.NeedsConfirmation {
		t.Fatalf("expected confirmation needed for date-only conflict, got %+v", result)
	}
}

func TestExecutor_DeleteTimeOnlyWhenSingleMatch(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "健身", Date: "2026-05-29", Time: "19:00"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Date:   "2026-05-29",
		Time:   "19:00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultDeleted || result.Event == nil {
		t.Fatalf("unexpected delete result: %+v", result)
	}
}

func TestExecutor_DeleteTimeOnlyConflictWhenMultipleMatches(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "健身", Date: "2026-05-29", Time: "19:00"})
	_, _ = store.Create(Event{Title: "开会", Date: "2026-05-29", Time: "19:00"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Date:   "2026-05-29",
		Time:   "19:00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultConflict || len(result.Events) != 2 {
		t.Fatalf("unexpected conflict result: %+v", result)
	}
	if !result.NeedsConfirmation {
		t.Fatalf("expected confirmation needed for date+time conflict, got %+v", result)
	}
}

func TestExecutor_DeleteTitleAndTimeConflictNeedsConfirmation(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "项目会议", Date: "2026-05-29", Time: "19:00"})
	_, _ = store.Create(Event{Title: "周会会议", Date: "2026-05-29", Time: "19:00"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Title:  "会议",
		Date:   "2026-05-29",
		Time:   "19:00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultConflict || len(result.Events) != 2 {
		t.Fatalf("unexpected conflict result: %+v", result)
	}
	if !result.NeedsConfirmation {
		t.Fatalf("expected confirmation needed for title+date+time conflict, got %+v", result)
	}
}

func TestExecutor_DeleteNoMatchDoesNotNeedConfirmation(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	_, _ = store.Create(Event{Title: "健身", Date: "2026-05-29", Time: "19:00"})

	result, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentDeleteEvent,
		Date:   "2026-05-29",
		Time:   "21:00",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Kind != ResultConflict || result.Message != "no matching event found" {
		t.Fatalf("unexpected no-match result: %+v", result)
	}
	if result.NeedsConfirmation {
		t.Fatalf("expected no confirmation for no-match result, got %+v", result)
	}
}

func TestExecutor_AllowsMultipleEventsSameDateTime(t *testing.T) {
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store := NewInMemoryStore(func() time.Time { return now })
	exec := New(store)

	first, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentCreateEvent,
		Title:  "要写作业",
		Date:   "2026-06-06",
		Time:   "18:00",
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	second, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentCreateEvent,
		Title:  "去超市买东西",
		Date:   "2026-06-06",
		Time:   "18:00",
	})
	if err != nil {
		t.Fatalf("unexpected create error: %v", err)
	}
	if first.Event == nil || second.Event == nil || first.Event.ID == second.Event.ID {
		t.Fatalf("expected distinct created events, got first=%+v second=%+v", first, second)
	}

	queryResult, err := exec.Execute(parserlib.ParsedCommand{
		Intent: parserlib.IntentQueryEvents,
		Date:   "2026-06-06",
		Time:   "18:00",
	})
	if err != nil {
		t.Fatalf("unexpected query error: %v", err)
	}
	if queryResult.Kind != ResultQueried || len(queryResult.Events) != 2 {
		t.Fatalf("unexpected query result: %+v", queryResult)
	}
}
