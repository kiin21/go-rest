package event_handler

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/repository"
)

// EventHandler handle events from Kafka
type EventHandler struct {
	repo domainRepo.NotificationRepository
}

func NewEventHandler(repo domainRepo.NotificationRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

func (h *EventHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *EventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf(
			"✅ Received message from topic[%s], partition[%d], offset[%d]: %s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Value),
		)

		// Parse event
		event, err := events.BytesToEvent(msg.Value)
		if err != nil {
			log.Printf("Failed to parse event: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		log.Printf("Processing event: Type=%s, ID=%s, Key=%s", event.Type, event.ID, event.Key)

		// Get context from session
		ctx := session.Context()

		switch event.Type {
		case events.EventTypeNotificationLeaderAssignment:
			if err := h.handleLeaderAssignment(ctx, event); err != nil {
				log.Printf("Failed to handle leader assignment: %v", err)
			}
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
func (h *EventHandler) handleLeaderAssignment(ctx context.Context, event *events.Event) error {
	// BytesToEvent event payload
	var payload events.LeaderAssignmentEventPayload

	if err := event.UnmarshalData(&payload); err != nil {
		log.Printf("Failed to unmarshal leader assignment event: %v", err)
		return err
	}

	// Create model
	notification := &model.Notification{
		ID:          event.ID.String(),
		FromStarter: payload.FromStarter,
		ToStarter:   payload.ToStarter,
		Message:     payload.Message,
		Type:        event.Type,
		Timestamp:   event.Timestamp,
	}

	// Save to db
	if err := h.repo.Create(ctx, notification); err != nil {
		log.Printf("Failed to create notification: %v", err)
		return err
	}

	log.Printf("Leader assignment notification created for user %s", notification.ID)
	return nil
}
