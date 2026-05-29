package parser

import (
	"strings"
	"time"
)

var fallbackCreateActivities = []string{
	"开会",
	"看电影",
	"买东西",
	"写作业",
	"打游戏",
	"复习",
	"上课",
	"学习",
	"考试",
	"写代码",
	"工作",
	"健身",
	"跑步",
	"打球",
	"看书",
	"看剧",
	"散步",
	"吃饭",
	"写",
	"去",
	"见",
	"买",
}

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
		for _, act := range fallbackCreateActivities {
			if strings.Contains(trimmed, act) {
				intent = IntentCreateEvent
				break
			}
		}
		if intent == "" {
			return ParseResult{Err: &ParseError{Reason: ErrUnknownIntent, Message: "unable to detect intent"}}
		}
	}

	title := ExtractTitle(trimmed, intent, dt)

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
