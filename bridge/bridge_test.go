package bridge

import (
	"encoding/json"
	"testing"
)

func TestParseText(t *testing.T) {
	result, err := ParseText("明天下午三点和张三开会", "2026-05-29T12:00:00Z")
	if err != nil {
		t.Fatalf("ParseText returned error: %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}

	if parsed["intent"] != "create_event" {
		t.Fatalf("unexpected intent: %v", parsed["intent"])
	}
	if parsed["title"] != "和张三开会" {
		t.Fatalf("unexpected title: %v", parsed["title"])
	}
	if parsed["date"] != "2026-05-30" {
		t.Fatalf("unexpected date: %v", parsed["date"])
	}
	if parsed["time"] != "15:00" {
		t.Fatalf("unexpected time: %v", parsed["time"])
	}
}
