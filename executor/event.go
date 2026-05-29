package executor

import "time"

type Event struct {
	ID           string     `json:"id"`
	Title        string     `json:"title"`
	Date         string     `json:"date"`
	Time         string     `json:"time,omitempty"`
	NeedReminder *bool      `json:"need_reminder,omitempty"`
	Status       string     `json:"status"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}

type QueryOptions struct {
	Title string
	Date  string
	Time  string
}

type ResultKind string

const (
	ResultCreated  ResultKind = "created"
	ResultQueried  ResultKind = "queried"
	ResultDeleted  ResultKind = "deleted"
	ResultConflict ResultKind = "conflict"
)

type ExecuteResult struct {
	Kind              ResultKind `json:"kind"`
	Event             *Event     `json:"event,omitempty"`
	Events            []Event    `json:"events,omitempty"`
	Message           string     `json:"message,omitempty"`
	NeedsConfirmation bool       `json:"needs_confirmation,omitempty"`
}
