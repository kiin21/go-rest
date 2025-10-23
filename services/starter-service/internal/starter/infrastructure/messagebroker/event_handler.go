package messagebroker

import (
	"context"
	"fmt"
	"log"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
)

// EventHandler handles synchronization events between MySQL and Elasticsearch
type EventHandler struct {
	starterRepo       repository.StarterRepository
	searchRepo        repository.StarterSearchRepository
	enrichmentService *domainService.StarterEnrichmentService
}

func NewEventHandler(
	starterRepo repository.StarterRepository,
	searchRepo repository.StarterSearchRepository,
) *EventHandler {
	return &EventHandler{
		starterRepo: starterRepo,
		searchRepo:  searchRepo,
	}
}

// Handle processes incoming sync events
func (h *EventHandler) Handle(ctx context.Context, event *events.Event) error {
	log.Printf("Processing event: type=%s, timestamp=%s", event.Type, event.Timestamp)

	switch event.Type {
	case events.EventTypeStarterIndex, events.EventTypeStarterUpdate, events.EventTypeStarterInsert:
		return h.handleIndexEvent(ctx, event)
	case events.EventTypeStarterDelete:
		return h.handleDeleteEvent(ctx, event)
	default:
		log.Printf("Warning: unknown event type: %s", event.Type)
		return nil
	}
}

// handleIndexEvent handles index/update/insert events
func (h *EventHandler) handleIndexEvent(ctx context.Context, event *events.Event) error {
	// Unmarshal the payload into the struct
	var payload events.IndexStarterPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Now you can access the fields from the unmarshaled payload
	starter, err := h.starterRepo.FindByDomain(ctx, payload.Domain)
	if err != nil {
		return fmt.Errorf("failed to fetch starter from MySQL: %w", err)
	}

	// TODO: refactor
	batchSize := 100
	page := 1

	emptyQuery := &starterquery.ListStartersQuery{
		Pagination: httputil.ReqPagination{
			Page: &page, Limit: &batchSize,
		},
		Keyword:  starter.Domain(),
		SearchBy: "domain",
	}

	starters, _, err := h.starterRepo.SearchByKeyword(ctx, emptyQuery)
	if err != nil {
		return err
	}

	enriched, err := h.enrichmentService.EnrichStarters(ctx, starters)

	esDocs := make([]*model.StarterESDoc, len(starters))
	for i := range starters {
		esDoc := model.NewStarterESDocFromStarter(starters[i], enriched)
		esDocs[i] = esDoc
	}

	// Index to Elasticsearch
	if err := h.searchRepo.IndexStarter(ctx, esDocs[0]); err != nil {
		return fmt.Errorf("failed to index starter to Elasticsearch: %w", err)
	}

	log.Printf("Successfully indexed starter: domain=%s", payload.Domain)
	return nil
}

// handleDeleteEvent handles delete events
func (h *EventHandler) handleDeleteEvent(ctx context.Context, event *events.Event) error {
	var payload events.IndexStarterPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if err := h.searchRepo.DeleteFromIndex(ctx, payload.Domain); err != nil {
		return fmt.Errorf("failed to delete starter from Elasticsearch: %w", err)
	}

	log.Printf("Successfully deleted starter from index: domain=%s", payload.Domain)
	return nil
}
