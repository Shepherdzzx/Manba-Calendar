package executor

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateEventID() (string, error) {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to generate event id: %w", err)
	}
	return "evt-" + hex.EncodeToString(buf), nil
}
