package event

import (
	"encoding/json"
	"fmt"
	"strings"
)

type Event struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

const (
	EventSendMessage = "message"
)

type MessageEvent struct {
	Message string `json:"message"`
	User    string `json:"user"`
}

func (m MessageEvent) String() string {
	return fmt.Sprintf("%s:%s", strings.TrimSpace(m.User), strings.TrimSpace(m.Message))
}
