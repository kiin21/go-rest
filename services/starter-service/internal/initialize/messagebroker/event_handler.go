package messagebroker

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type EventHandler struct {
	repo domainRepo.StarterSearchRepository
}

func NewEventHandler(repo domainRepo.StarterSearchRepository) *EventHandler {
	return &EventHandler{repo: repo}
}

func (h *EventHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *EventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Printf(
			"Received message from topic[%s], partition[%d], offset[%d]: %s\n",
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

		var payload events.IndexStarterPayload

		if err := event.UnmarshalData(&payload); err != nil {
			log.Printf("Failed to unmarshal leader assignment event: %v", err)
			return err
		}

		// TODO: switch to handle base on event type
		starter, err := model.NewStarter(payload.Domain, payload.Name, "", "", "", "", nil, nil)
		if err != nil {
			log.Printf("Failed to create starter: %v from message", err)
			return err
		}
		ctx := session.Context()

		switch event.Type {
		case events.EventTypeStarterInsert, events.EventTypeStarterUpdate:
			err := h.repo.IndexStarter(ctx, starter)
			if err != nil {
				return err
			}
		case events.EventTypeStarterDelete:
			err := h.repo.DeleteFromIndex(ctx, starter.Domain())
			if err != nil {
				return err
			}
		default:
			log.Printf("Unknown event type: %s", event.Type)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}
