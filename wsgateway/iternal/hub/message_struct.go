package hub

import (
	"encoding/json"
	"time"
)

type ClientMessage struct {
	Type    string          `json:"type"` // "chat_message", "join_room", "typing"
	Payload json.RawMessage `json:"payload"`
}

type MessageToBroadcast struct {
	From   string
	Text   string
	RoomID string
	time   time.Time
}
