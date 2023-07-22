package event

import (
	"encoding/json"
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
	From    string `json:"from"`
}
