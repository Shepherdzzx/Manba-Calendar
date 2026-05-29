package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	parserlib "github.com/Shepherdzzx/manba-alert/parser"
)

type config struct {
	now     string
	pretty  bool
	demo    bool
	input   string
}

type demoCase struct {
	name  string
	input string
}

var demoCases = []demoCase{
	{name: "example1_create_meeting", input: "明天下午三点和张三开会"},
	{name: "example2_create_shopping", input: "后天上午十点去超市买东西"},
	{name: "example3_query_today", input: "帮我看看今天还有什么安排"},
	{name: "example4_query_meetings_next_monday", input: "下周一有什么会议"},
	{name: "example5_delete_gym_tonight", input: "取消今晚七点的健身"},
	{name: "example6_delete_interview_tomorrow_morning", input: "删除明天上午的面试"},
	{name: "example7_create_meeting_next_wednesday", input: "下周三下午两点半开会"},
	{name: "example8_create_movie_friday_evening", input: "星期五晚上八点看电影"},
	{name: "create_event_with_reminder", input: "明天下午三点提醒我和张三开会"},
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	cfg, help, err := parseFlags(args)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}
	if help {
		printUsage(stdout)
		return 0
	}

	now, err := parseNow(cfg.now, cfg.demo)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	if cfg.demo {
		return runDemo(now, cfg.pretty, stdout, stderr)
	}

	input, err := readInput(cfg.input, stdin)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	result := parserlib.New(now).Parse(input)
	if result.Err != nil {
		fmt.Fprintf(stderr, "%s: %s\n", result.Err.Reason, result.Err.Message)
		return 1
	}

	if err := writeJSON(stdout, result.Command, cfg.pretty); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func parseFlags(args []string) (config, bool, error) {
	fs := flag.NewFlagSet("parser-cli", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var cfg config
	var help bool
	fs.StringVar(&cfg.now, "now", "", "override reference time (YYYY-MM-DD or RFC3339)")
	fs.BoolVar(&cfg.pretty, "pretty", false, "pretty-print JSON output")
	fs.BoolVar(&cfg.demo, "demo", false, "run built-in parser demo cases")
	fs.BoolVar(&help, "help", false, "show usage")
	if err := fs.Parse(args); err != nil {
		return config{}, false, err
	}
	remaining := fs.Args()
	if len(remaining) > 1 {
		return config{}, false, fmt.Errorf("too many positional arguments")
	}
	if len(remaining) == 1 {
		cfg.input = remaining[0]
	}
	return cfg, help, nil
}

func parseNow(raw string, demo bool) (time.Time, error) {
	if raw == "" {
		if demo {
			return time.Date(2026, 5, 29, 12, 0, 0, 0, time.UTC), nil
		}
		return time.Now(), nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return time.Date(t.Year(), t.Month(), t.Day(), 12, 0, 0, 0, time.UTC), nil
	}
	return time.Time{}, fmt.Errorf("invalid --now value: %s", raw)
}

func readInput(positional string, stdin io.Reader) (string, error) {
	if strings.TrimSpace(positional) != "" {
		return strings.TrimSpace(positional), nil
	}
	data, err := io.ReadAll(stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read stdin: %w", err)
	}
	input := strings.TrimSpace(string(data))
	if input == "" {
		return "", fmt.Errorf("no input provided")
	}
	return input, nil
}

func runDemo(now time.Time, pretty bool, stdout, stderr io.Writer) int {
	writer := bufio.NewWriter(stdout)
	defer writer.Flush()
	fmt.Fprintln(writer, "=== Demo: Parser Examples ===")
	fmt.Fprintln(writer)
	p := parserlib.New(now)
	for i, tc := range demoCases {
		fmt.Fprintf(writer, "#%d %s\n", i+1, tc.name)
		fmt.Fprintf(writer, "Input: %s\n", tc.input)
		result := p.Parse(tc.input)
		if result.Err != nil {
			fmt.Fprintf(stderr, "%s: %s\n", result.Err.Reason, result.Err.Message)
			return 1
		}
		if err := writeJSON(writer, result.Command, pretty); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		fmt.Fprintln(writer)
		fmt.Fprintln(writer)
	}
	return 0
}

func writeJSON(w io.Writer, v any, pretty bool) error {
	var (
		data []byte
		err  error
	)
	if pretty {
		data, err = json.MarshalIndent(v, "", "  ")
	} else {
		data, err = json.Marshal(v)
	}
	if err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}
	_, err = fmt.Fprintln(w, string(data))
	return err
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  go run ./cmd/parser-cli [flags] \"中文文本\"")
	fmt.Fprintln(w, "  echo \"中文文本\" | go run ./cmd/parser-cli [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --now     override reference time (YYYY-MM-DD or RFC3339)")
	fmt.Fprintln(w, "  --pretty  pretty-print JSON output")
	fmt.Fprintln(w, "  --demo    run built-in demo cases")
	fmt.Fprintln(w, "  --help    show this help message")
}
