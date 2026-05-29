package parser

type Intent string

const (
	IntentCreateEvent Intent = "create_event"
	IntentQueryEvents Intent = "query_events"
	IntentDeleteEvent Intent = "delete_event"
)

const (
	ErrEmptyInput      = "empty_input"
	ErrUnknownIntent   = "unknown_intent"
	ErrMissingTitle    = "missing_title"
	ErrMissingDate     = "missing_date"
	ErrMissingTime     = "missing_time"
	ErrInvalidDate     = "invalid_date"
	ErrInvalidTime     = "invalid_time"
	ErrInvalidIntent   = "invalid_intent"
)

type ParsedCommand struct {
	Intent       Intent `json:"intent"`
	Title        string `json:"title,omitempty"`
	Date         string `json:"date,omitempty"`
	Time         string `json:"time,omitempty"`
	NeedReminder *bool  `json:"need_reminder,omitempty"`
}

type ParseError struct {
	Reason  string
	Message string
}

func (e *ParseError) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

type ParseResult struct {
	Command *ParsedCommand
	Err     *ParseError
}
