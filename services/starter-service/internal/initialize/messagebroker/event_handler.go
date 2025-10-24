package messagebroker

import (
	"context"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	domainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
)

type EventHandler struct {
	starterRepo       domainRepo.StarterRepository
	starterSearchRepo domainRepo.StarterSearchRepository
	enrichmentService *domainService.StarterEnrichmentService
}

func NewEventHandler(
	starterRepo domainRepo.StarterRepository,
	starterSearchRepo domainRepo.StarterSearchRepository,
	enrichmentService *domainService.StarterEnrichmentService,
) *EventHandler {
	return &EventHandler{
		starterRepo:       starterRepo,
		starterSearchRepo: starterSearchRepo,
		enrichmentService: enrichmentService,
	}
}

func (h *EventHandler) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (h *EventHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *EventHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := session.Context()

		// Parse event
		event, err := events.BytesToEvent(msg.Value)
		if err != nil {
			log.Printf("Failed to parse event: %v", err)
			session.MarkMessage(msg, "")
			continue
		}

		log.Printf("Processing event: Type=%s, ID=%s", event.Type, event.ID)

	var payload events.IndexStarterPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return fmt.Errorf("failed to unmarshal IndexStarterPayload: %w", err)
	}

	switch event.Type {
		case events.EventTypeStarterInsert, events.EventTypeStarterUpdate:
			{
				esDocs, err := h.fetchAndEnrichStarter(ctx, payload.Domain)
				if err != nil {
					return fmt.Errorf("failed to fetch and enrich starter: %w", err)
				}
				return h.starterSearchRepo.IndexStarter(ctx, esDocs)
			}
		case events.EventTypeStarterDelete:
			{
				return h.starterSearchRepo.DeleteFromIndex(ctx, payload.Domain)
			}

		default:
			log.Printf("Unknown event type: %s", event.Type)
		}

		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *EventHandler) fetchAndEnrichStarter(
	ctx context.Context,
	domain string,
) (*model.StarterESDoc, error) {
	// Fetch starter from MySQL
	starter, err := h.starterRepo.FindByDomain(ctx, domain)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch starter from MySQL: %w", err)
	}

	// Enrich single starter
	enriched, err := h.enrichmentService.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, fmt.Errorf("failed to enrich starter: %w", err)
	}

	// Convert to ES document
	esDoc := model.NewStarterESDocFromStarter(starter, enriched)
	return esDoc, nil
}

func (h *EventHandler) enrichStartersToESDocs(
	ctx context.Context,
	starters []*model.Starter,
) ([]*model.StarterESDoc, error) {
	// Enrich all starters in batch
	enriched, err := h.enrichmentService.EnrichStarters(ctx, starters)
	if err != nil {
		return nil, err
	}

	// Convert to ES documents
	esDocs := make([]*model.StarterESDoc, len(starters))
	for i := range starters {
		esDocs[i] = model.NewStarterESDocFromStarter(starters[i], enriched)
	}

	return esDocs, nil
}
