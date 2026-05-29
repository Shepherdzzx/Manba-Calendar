package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	executorlib "github.com/Shepherdzzx/manba-alert/executor"
	parserlib "github.com/Shepherdzzx/manba-alert/parser"
)

type config struct {
	now    string
	pretty bool
	json   bool
	store  string
	input  string
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

var interactiveStdin = os.Stdin

var isInteractive = func() bool {
	info, err := interactiveStdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
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

	now, err := parseNow(cfg.now)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	input, err := readInput(cfg.input, stdin)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 2
	}

	parseResult := parserlib.New(now).Parse(input)
	if parseResult.Err != nil {
		fmt.Fprintf(stderr, "%s: %s\n", parseResult.Err.Reason, parseResult.Err.Message)
		return 1
	}

	store, err := openStore(cfg.store, func() time.Time { return now })
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	exec := executorlib.New(store)
	result, err := exec.Execute(*parseResult.Command)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if result.Kind == executorlib.ResultConflict && result.NeedsConfirmation && !cfg.json && isInteractive() {
		return confirmAndDelete(store, result, cfg.pretty, stdout, stderr)
	}

	if err := writeJSON(stdout, result, cfg.pretty); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func parseFlags(args []string) (config, bool, error) {
	fs := flag.NewFlagSet("executor-cli", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	var cfg config
	var help bool
	fs.StringVar(&cfg.now, "now", "", "override reference time (YYYY-MM-DD or RFC3339)")
	fs.BoolVar(&cfg.pretty, "pretty", false, "pretty-print JSON output")
	fs.BoolVar(&cfg.json, "json", false, "output raw JSON only (disable interactive confirmation)")
	fs.StringVar(&cfg.store, "store", defaultStorePath(), "path to JSON store file")
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

func parseNow(raw string) (time.Time, error) {
	if raw == "" {
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

func confirmAndDelete(store executorlib.Store, result executorlib.ExecuteResult, pretty bool, stdout, stderr io.Writer) int {
	fmt.Fprintln(stderr, "找到多个匹配事件:")
	for i, event := range result.Events {
		if event.Time != "" {
			fmt.Fprintf(stderr, "  %d. %s (%s %s)\n", i+1, event.Title, event.Date, event.Time)
			continue
		}
		fmt.Fprintf(stderr, "  %d. %s (%s)\n", i+1, event.Title, event.Date)
	}
	fmt.Fprint(stderr, "请输入要删除的编号，或输入 c 取消: ")

	reader := bufio.NewReader(interactiveStdin)
	selection, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(stderr, "failed to read confirmation input: %v\n", err)
		return 1
	}
	selection = strings.TrimSpace(selection)
	if strings.EqualFold(selection, "c") {
		if err := writeJSON(stdout, executorlib.ExecuteResult{Kind: executorlib.ResultConflict, Message: "delete cancelled"}, pretty); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
		return 0
	}

	index, err := strconv.Atoi(selection)
	if err != nil || index < 1 || index > len(result.Events) {
		fmt.Fprintln(stderr, "invalid selection")
		return 1
	}

	deleted, err := store.Delete(result.Events[index-1].ID)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	if err := writeJSON(stdout, executorlib.ExecuteResult{Kind: executorlib.ResultDeleted, Event: &deleted, Message: "event deleted"}, pretty); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}

func defaultStorePath() string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return ".manba-alert-events.json"
	}
	return home + "/.manba-alert/events.json"
}

func openStore(path string, nowFunc func() time.Time) (executorlib.Store, error) {
	return executorlib.NewJSONFileStore(path, nowFunc)
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  go run ./cmd/executor-cli [flags] \"中文文本\"")
	fmt.Fprintln(w, "  echo \"中文文本\" | go run ./cmd/executor-cli [flags]")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "Flags:")
	fmt.Fprintln(w, "  --now     override reference time (YYYY-MM-DD or RFC3339)")
	fmt.Fprintln(w, "  --pretty  pretty-print JSON output")
	fmt.Fprintln(w, "  --json    output raw JSON only (disable interactive confirmation)")
	fmt.Fprintln(w, "  --store   path to JSON store file")
	fmt.Fprintln(w, "  --help    show this help message")
}
