package bridge

import (
	"encoding/json"
	"fmt"
	"time"

	parserlib "github.com/Shepherdzzx/manba-alert/parser"
)

// ParseText parses a Chinese natural language command and returns a JSON command.
func ParseText(input string, nowRFC3339 string) (string, error) {
	now, err := time.Parse(time.RFC3339, nowRFC3339)
	if err != nil {
		now = time.Now()
	}

	result := parserlib.New(now).Parse(input)
	if result.Err != nil {
		return "", fmt.Errorf("%s: %s", result.Err.Reason, result.Err.Message)
	}

	data, err := json.Marshal(result.Command)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
