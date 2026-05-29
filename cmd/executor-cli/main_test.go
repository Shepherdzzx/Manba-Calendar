package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	executorlib "github.com/Shepherdzzx/manba-alert/executor"
)

func TestRun_CreateCommand(t *testing.T) {
	var stdout, stderr bytes.Buffer
	storePath := filepath.Join(t.TempDir(), "events.json")
	exitCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "今天要写作业"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "created") || !strings.Contains(stdout.String(), "要写作业") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRun_PersistenceAcrossInvocations(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	var createOut, createErr bytes.Buffer
	createCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "今天要写作业"}, strings.NewReader(""), &createOut, &createErr)
	if createCode != 0 {
		t.Fatalf("unexpected create exit code: %d, stderr=%s", createCode, createErr.String())
	}

	var queryOut, queryErr bytes.Buffer
	queryCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "帮我看看今天还有什么安排"}, strings.NewReader(""), &queryOut, &queryErr)
	if queryCode != 0 {
		t.Fatalf("unexpected query exit code: %d, stderr=%s", queryCode, queryErr.String())
	}
	if !strings.Contains(queryOut.String(), "queried") || !strings.Contains(queryOut.String(), "要写作业") {
		t.Fatalf("expected persisted event in query result, got %s", queryOut.String())
	}
}

func TestRun_InteractiveDeleteConfirm(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store, err := executorlib.NewJSONFileStore(storePath, func() time.Time { return now })
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "要写作业", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed first event: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "去超市", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed second event: %v", err)
	}

	oldInteractive := isInteractive
	oldStdin := interactiveStdin
	defer func() {
		isInteractive = oldInteractive
		interactiveStdin = oldStdin
	}()
	isInteractive = func() bool { return true }
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	interactiveStdin = r
	if _, err := w.WriteString("1\n"); err != nil {
		t.Fatalf("failed to write selection: %v", err)
	}
	_ = w.Close()

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "删除今天的事"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "deleted") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "要写作业") && !strings.Contains(stdout.String(), "去超市") {
		t.Fatalf("expected deleted event details in output, got %s", stdout.String())
	}
	if !strings.Contains(stderr.String(), "找到多个匹配事件") {
		t.Fatalf("expected confirmation prompt, got stderr=%s", stderr.String())
	}
}

func TestRun_InteractiveDeleteCancel(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store, err := executorlib.NewJSONFileStore(storePath, func() time.Time { return now })
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "要写作业", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed first event: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "去超市", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed second event: %v", err)
	}

	oldInteractive := isInteractive
	oldStdin := interactiveStdin
	defer func() {
		isInteractive = oldInteractive
		interactiveStdin = oldStdin
	}()
	isInteractive = func() bool { return true }
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	interactiveStdin = r
	if _, err := w.WriteString("c\n"); err != nil {
		t.Fatalf("failed to write selection: %v", err)
	}
	_ = w.Close()

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "删除今天的事"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "delete cancelled") || !strings.Contains(stdout.String(), "conflict") {
		t.Fatalf("unexpected output: %s", stdout.String())
	}
}

func TestRun_InteractiveDeleteInvalidSelection(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store, err := executorlib.NewJSONFileStore(storePath, func() time.Time { return now })
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "要写作业", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed first event: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "去超市", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed second event: %v", err)
	}

	oldInteractive := isInteractive
	oldStdin := interactiveStdin
	defer func() {
		isInteractive = oldInteractive
		interactiveStdin = oldStdin
	}()
	isInteractive = func() bool { return true }
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	interactiveStdin = r
	if _, err := w.WriteString("99\n"); err != nil {
		t.Fatalf("failed to write selection: %v", err)
	}
	_ = w.Close()

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "删除今天的事"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stderr.String(), "invalid selection") {
		t.Fatalf("expected invalid selection error, got stderr=%s", stderr.String())
	}
}

func TestRun_DeleteJSONFlagKeepsRawConflictOutput(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store, err := executorlib.NewJSONFileStore(storePath, func() time.Time { return now })
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "要写作业", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed first event: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "去超市", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed second event: %v", err)
	}

	oldInteractive := isInteractive
	defer func() { isInteractive = oldInteractive }()
	isInteractive = func() bool { return true }

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--json", "--store=" + storePath, "--now=2026-05-29", "删除今天的事"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "needs_confirmation") || !strings.Contains(stdout.String(), "multiple matching events found") {
		t.Fatalf("expected raw conflict output, got %s", stdout.String())
	}
	if strings.Contains(stderr.String(), "找到多个匹配事件") {
		t.Fatalf("did not expect interactive prompt, got stderr=%s", stderr.String())
	}
}

func TestRun_DeleteConflictFallsBackToRawJSONWhenNotInteractive(t *testing.T) {
	storePath := filepath.Join(t.TempDir(), "events.json")
	now := time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC)
	store, err := executorlib.NewJSONFileStore(storePath, func() time.Time { return now })
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "要写作业", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed first event: %v", err)
	}
	if _, err := store.Create(executorlib.Event{Title: "去超市", Date: "2026-05-29"}); err != nil {
		t.Fatalf("failed to seed second event: %v", err)
	}

	oldInteractive := isInteractive
	defer func() { isInteractive = oldInteractive }()
	isInteractive = func() bool { return false }

	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--store=" + storePath, "--now=2026-05-29", "删除今天的事"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "needs_confirmation") || !strings.Contains(stdout.String(), "multiple matching events found") {
		t.Fatalf("expected raw conflict output, got %s", stdout.String())
	}
	if stderr.Len() != 0 {
		t.Fatalf("expected no prompt output, got stderr=%s", stderr.String())
	}
}

func TestRun_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--help"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected usage output, got %s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "--json") {
		t.Fatalf("expected --json in usage output, got %s", stdout.String())
	}
}
