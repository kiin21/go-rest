package messaging

import "github.com/kiin21/go-rest/pkg/events"

// NotificationProducer sends notification events (leader assignment, etc.)
type NotificationProducer interface {
	SendNotification(event *events.Event) error
	Close() error
}

// SyncProducer sends sync events for Elasticsearch indexing
type SyncProducer interface {
	SendSyncEvent(event *events.Event) error
	Close() error
}
