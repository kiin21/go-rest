package messaging

import "github.com/kiin21/go-rest/pkg/events"

type NotificationProducer interface {
	SendNotification(event *events.Event) error
	Close() error
}
