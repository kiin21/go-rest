package events

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	EventTypeStarterInsert = "starter.insert"
	EventTypeStarterUpdate = "starter.update"
	EventTypeStarterDelete = "starter.delete"
	EventTypeStarterIndex  = "starter.index"
)

type Event struct {
	Type      string            `json:"type"`
	Domain    string            `json:"domain,omitempty"`
	Key       string            `json:"key,omitempty"`
	Payload   json.RawMessage   `json:"payload,omitempty"`
	Metadata  map[string]string `json:"metadata,omitempty"`
	Timestamp time.Time         `json:"timestamp"`
	Retries   int               `json:"retries,omitempty"`
}

func NewEvent(eventType string, opts ...EventOption) (*Event, error) {
	if eventType == "" {
		return nil, errors.New("event type is required")
	}

	e := &Event{
		Type:      eventType,
		Timestamp: time.Now(),
	}
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(e); err != nil {
			return nil, err
		}
	}
	return e, nil
}

type EventOption func(*Event) error

func WithDomain(domain string) EventOption {
	return func(e *Event) error {
		e.Domain = domain
		return nil
	}
}

func WithKey(key string) EventOption {
	return func(e *Event) error {
		e.Key = key
		return nil
	}
}

func WithMetadata(metadata map[string]string) EventOption {
	return func(e *Event) error {
		if len(metadata) == 0 {
			return nil
		}
		if e.Metadata == nil {
			e.Metadata = make(map[string]string, len(metadata))
		}
		for k, v := range metadata {
			e.Metadata[k] = v
		}
		return nil
	}
}

func WithPayload(payload interface{}) EventOption {
	return func(e *Event) error {
		if payload == nil {
			e.Payload = nil
			return nil
		}
		bytes, err := json.Marshal(payload)
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}
		e.Payload = bytes
		return nil
	}
}

func WithTimestamp(ts time.Time) EventOption {
	return func(e *Event) error {
		if !ts.IsZero() {
			e.Timestamp = ts
		}
		return nil
	}
}

func WithRetries(count int) EventOption {
	return func(e *Event) error {
		if count < 0 {
			return fmt.Errorf("retry count cannot be negative: %d", count)
		}
		e.Retries = count
		return nil
	}
}

func (e *Event) DecodePayload(out interface{}) error {
	if out == nil {
		return errors.New("decode target must not be nil")
	}
	if len(e.Payload) == 0 {
		return nil
	}
	return json.Unmarshal(e.Payload, out)
}
