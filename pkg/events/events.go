package events

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        uuid.UUID       `json:"id"`
	Type      string          `json:"type"`
	Payload   json.RawMessage `json:"payload,omitempty"`
	Timestamp time.Time       `json:"timestamp"`
}

func NewEvent(eventType string, eventPayload interface{}) (*Event, error) {
	payloadBytes, err := json.Marshal(eventPayload)

	if err != nil {
		return nil, fmt.Errorf("failed to marshal event payload: %w", err)
	}

	return &Event{
		ID:        uuid.New(),
		Type:      eventType,
		Timestamp: time.Now().UTC(),
		Payload:   payloadBytes,
	}, nil
}

func (e *Event) ToBytes() ([]byte, error) {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event: %w", err)
	}
	return bytes, nil
}

func (e *Event) UnmarshalPayload(target interface{}) error {
	if len(e.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.Payload, target)
}

func BytesToEvent(data []byte) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, fmt.Errorf("failed to unmarshal event: %w", err)
	}
	return &event, nil
}
