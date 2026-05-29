package parser

import (
	"strings"
	"time"
)

type Parser struct {
	now time.Time
}

func New(now time.Time) Parser {
	return Parser{now: now}
}

func (p Parser) Parse(input string) ParseResult {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return ParseResult{Err: &ParseError{Reason: ErrEmptyInput, Message: "input text is empty"}}
	}

	intent := DetectIntent(trimmed)
	dt := NewDateTimeParser(p.now).Parse(trimmed)
	if intent == "" {
		if strings.Contains(trimmed, "开会") || strings.Contains(trimmed, "看电影") || strings.Contains(trimmed, "买东西") {
			intent = IntentCreateEvent
		} else if dt.Date != "" || dt.Time != "" {
			intent = IntentCreateEvent
		} else {
			return ParseResult{Err: &ParseError{Reason: ErrUnknownIntent, Message: "unable to detect intent"}}
		}
	}

	title := ExtractTitle(trimmed, dt)

	cmd := ParsedCommand{
		Intent: intent,
		Title:  title,
		Date:   dt.Date,
		Time:   dt.Time,
	}

	if strings.Contains(trimmed, "提醒") && intent == IntentCreateEvent {
		v := true
		cmd.NeedReminder = &v
	}

	if err := Validate(cmd); err != nil {
		return ParseResult{Err: err}
	}

	return ParseResult{Command: &cmd}
}
