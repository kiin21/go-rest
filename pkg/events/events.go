package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID       `json:"id"`
	Type      string          `json:"type"`
	Key       string          `json:"key,omitempty"`
	Domain    string          `json:"domain,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
	Payload   json.RawMessage `json:"payload,omitempty"` // Change to json.RawMessage
}

func NewEvent(eventType string, eventPayload interface{}) *Event {
	payloadBytes, _ := json.Marshal(eventPayload)
	return &Event{
		ID:        uuid.New(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payloadBytes,
	}
}

func (e *Event) ToBytes() ([]byte, error) {
	return json.Marshal(e)
}

// UnmarshalData unmarshal the event payload into the provided target
func (e *Event) UnmarshalData(target interface{}) error {
	if len(e.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.Payload, target)
}

func BytesToEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
