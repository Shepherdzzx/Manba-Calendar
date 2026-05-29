package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRun_PositionalInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--now=2026-05-29", "明天下午三点和张三开会"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "create_event") {
		t.Fatalf("expected create_event in output, got %s", stdout.String())
	}
}

func TestRun_StdinInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--now=2026-05-29"}, strings.NewReader("取消今晚七点的健身\n"), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "delete_event") {
		t.Fatalf("expected delete_event in output, got %s", stdout.String())
	}
}

func TestRun_ErrorPath(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"你好世界"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "unknown_intent") {
		t.Fatalf("expected unknown_intent in stderr, got %s", stderr.String())
	}
}

func TestRun_EmptyInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
}

func TestRun_PrettyOutput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--pretty", "--now=2026-05-29", "明天下午三点和张三开会"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "\n") || !strings.Contains(stdout.String(), "  \"intent\"") {
		t.Fatalf("expected pretty JSON output, got %s", stdout.String())
	}
}

func TestRun_DemoMode(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--demo", "--pretty"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d, stderr=%s", exitCode, stderr.String())
	}
	if !strings.Contains(stdout.String(), "=== Demo: Parser Examples ===") || !strings.Contains(stdout.String(), "example1_create_meeting") {
		t.Fatalf("expected demo output, got %s", stdout.String())
	}
}

func TestRun_Help(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--help"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("unexpected exit code: %d", exitCode)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Fatalf("expected usage text, got %s", stdout.String())
	}
}

func TestRun_BadNow(t *testing.T) {
	var stdout, stderr bytes.Buffer
	exitCode := run([]string{"--now=bad-time", "明天下午三点和张三开会"}, strings.NewReader(""), &stdout, &stderr)
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
	if !strings.Contains(stderr.String(), "invalid --now value") {
		t.Fatalf("expected invalid --now value error, got %s", stderr.String())
	}
}
