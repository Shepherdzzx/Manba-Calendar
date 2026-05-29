package parser

import (
	"regexp"
	"time"
)

var (
	dateRe = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	timeRe = regexp.MustCompile(`^\d{2}:\d{2}$`)
)

func Validate(cmd ParsedCommand) *ParseError {
	switch cmd.Intent {
	case IntentCreateEvent, IntentQueryEvents, IntentDeleteEvent:
	default:
		return &ParseError{Reason: ErrInvalidIntent, Message: "intent is invalid"}
	}

	if cmd.Date != "" {
		if !dateRe.MatchString(cmd.Date) {
			return &ParseError{Reason: ErrInvalidDate, Message: "date must be in YYYY-MM-DD format"}
		}
		if _, err := time.Parse("2006-01-02", cmd.Date); err != nil {
			return &ParseError{Reason: ErrInvalidDate, Message: "date is not a valid calendar date"}
		}
	}

	if cmd.Time != "" {
		if !timeRe.MatchString(cmd.Time) {
			return &ParseError{Reason: ErrInvalidTime, Message: "time must be in HH:MM format"}
		}
		if _, err := time.Parse("15:04", cmd.Time); err != nil {
			return &ParseError{Reason: ErrInvalidTime, Message: "time is not valid"}
		}
	}

	switch cmd.Intent {
	case IntentCreateEvent:
		if cmd.Title == "" {
			return &ParseError{Reason: ErrMissingTitle, Message: "create_event requires title"}
		}
		if cmd.Date == "" {
			return &ParseError{Reason: ErrMissingDate, Message: "create_event requires date"}
		}
		if cmd.Time == "" {
			return &ParseError{Reason: ErrMissingTime, Message: "create_event requires time"}
		}
	case IntentDeleteEvent:
		if cmd.Title == "" {
			return &ParseError{Reason: ErrMissingTitle, Message: "delete_event requires title"}
		}
	}

	return nil
}
