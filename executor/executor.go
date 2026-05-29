package executor

import (
	"fmt"

	parserlib "github.com/Shepherdzzx/manba-alert/parser"
)

type Executor struct {
	store Store
}

func New(store Store) *Executor {
	return &Executor{store: store}
}

func (e *Executor) Execute(cmd parserlib.ParsedCommand) (ExecuteResult, error) {
	switch cmd.Intent {
	case parserlib.IntentCreateEvent:
		event, err := e.store.Create(Event{
			Title:        cmd.Title,
			Date:         cmd.Date,
			Time:         cmd.Time,
			NeedReminder: cmd.NeedReminder,
		})
		if err != nil {
			return ExecuteResult{}, err
		}
		return ExecuteResult{Kind: ResultCreated, Event: &event, Message: "event created"}, nil
	case parserlib.IntentQueryEvents:
		events, err := e.store.Query(QueryOptions{Title: cmd.Title, Date: cmd.Date, Time: cmd.Time})
		if err != nil {
			return ExecuteResult{}, err
		}
		return ExecuteResult{Kind: ResultQueried, Events: events, Message: fmt.Sprintf("%d event(s) found", len(events))}, nil
	case parserlib.IntentDeleteEvent:
		var (
			matches []Event
			err     error
		)
		if cmd.Title != "" {
			matches, err = e.store.Query(QueryOptions{Title: cmd.Title, Date: cmd.Date, Time: cmd.Time})
		} else {
			matches, err = e.store.Query(QueryOptions{Date: cmd.Date, Time: cmd.Time})
		}
		if err != nil {
			return ExecuteResult{}, err
		}
		if len(matches) == 0 {
			return ExecuteResult{Kind: ResultConflict, Message: "no matching event found"}, nil
		}
		if len(matches) > 1 {
			return ExecuteResult{
				Kind:              ResultConflict,
				Events:            matches,
				Message:           "multiple matching events found",
				NeedsConfirmation: needsDeleteConfirmation(cmd),
			}, nil
		}
		deleted, err := e.store.Delete(matches[0].ID)
		if err != nil {
			return ExecuteResult{}, err
		}
		return ExecuteResult{Kind: ResultDeleted, Event: &deleted, Message: "event deleted"}, nil
	default:
		return ExecuteResult{}, fmt.Errorf("unsupported intent: %s", cmd.Intent)
	}
}

func needsDeleteConfirmation(cmd parserlib.ParsedCommand) bool {
	return cmd.Intent == parserlib.IntentDeleteEvent && (cmd.Date != "" || cmd.Time != "")
}
