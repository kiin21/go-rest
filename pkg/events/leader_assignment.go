package events

import (
	"errors"
	"time"
)

const EventTypeNotificationLeaderAssignment = "notification.leader_assignment"

type LeaderAssignmentNotification struct {
	FromStarter string    `json:"from_starter"`
	ToStarter   string    `json:"to_starter"`
	Message     string    `json:"message"`
	Type        string    `json:"type"`
	Timestamp   time.Time `json:"timestamp"`
}

func NewLeaderAssignmentEvent(notification *LeaderAssignmentNotification) (*Event, error) {
	if notification == nil {
		return nil, errors.New("leader assignment notification is required")
	}

	payload := *notification
	if payload.Timestamp.IsZero() {
		payload.Timestamp = time.Now()
	}

	return NewEvent(
		EventTypeNotificationLeaderAssignment,
		WithDomain(payload.ToStarter),
		WithKey(payload.ToStarter),
		WithTimestamp(payload.Timestamp),
		WithPayload(payload),
	)
}
