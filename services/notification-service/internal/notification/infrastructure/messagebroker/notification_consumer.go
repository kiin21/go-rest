package messagebroker

import (
	"context"
	"errors"
	"log"

	"github.com/kiin21/go-rest/pkg/events"
	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	domainmessaging "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/messaging"
)

// NotificationService defines the contract required by the consumer for persisting notifications.
type NotificationService interface {
	StoreNotification(context.Context, *events.LeaderAssignmentNotification) error
}

type kafkaNotificationConsumer struct {
	consumer *sharedKafka.EventConsumer
}

func NewKafkaNotificationConsumer(
	brokers []string,
	groupID string,
	topic string,
	service NotificationService,
) (domainmessaging.NotificationConsumer, error) {
	if service == nil {
		return nil, errors.New("notification service is required")
	}
	if len(brokers) == 0 {
		return nil, errors.New("notification consumer requires at least one broker")
	}
	if groupID == "" {
		return nil, errors.New("notification consumer requires a group id")
	}
	if topic == "" {
		return nil, errors.New("notification consumer requires a topic")
	}

	handler := func(ctx context.Context, event *events.Event) error {
		if event == nil {
			return nil
		}
		if event.Type != events.EventTypeNotificationLeaderAssignment {
			log.Printf("Warning: skipping unsupported notification event type %s", event.Type)
			return nil
		}

		var payload events.LeaderAssignmentNotification
		if err := event.DecodePayload(&payload); err != nil {
			return err
		}
		if payload.Timestamp.IsZero() {
			payload.Timestamp = event.Timestamp
		}

		return service.StoreNotification(ctx, &payload)
	}

	eventConsumer, err := sharedKafka.NewEventConsumer(
		brokers,
		groupID,
		[]string{topic},
		handler,
	)
	if err != nil {
		return nil, err
	}

	return &kafkaNotificationConsumer{consumer: eventConsumer}, nil
}

func (c *kafkaNotificationConsumer) Start() {
	if c == nil || c.consumer == nil {
		return
	}

	c.consumer.Start()
}

func (c *kafkaNotificationConsumer) Stop() {
	if c == nil || c.consumer == nil {
		return
	}

	c.consumer.Stop()
}
