package messagebroker

import (
	"context"
	"fmt"
	"log"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

// EventHandler handles synchronization events between MySQL and Elasticsearch
type EventHandler struct {
	starterRepo repository.StarterRepository
	searchRepo  repository.StarterSearchRepository
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
	log.Printf("Processing event: type=%s, domain=%s, timestamp=%s",
		event.Type, event.Domain, event.Timestamp)

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
	// Fetch the latest data from MySQL
	starter, err := h.starterRepo.FindByDomain(ctx, event.Domain)
	if err != nil {
		return fmt.Errorf("failed to fetch starter from MySQL: %w", err)
	}

	// Index to Elasticsearch
	if err := h.searchRepo.IndexStarter(ctx, starter); err != nil {
		return fmt.Errorf("failed to index starter to Elasticsearch: %w", err)
	}

	log.Printf("Successfully indexed starter: domain=%s", event.Domain)
	return nil
}

// handleDeleteEvent handles delete events
func (h *EventHandler) handleDeleteEvent(ctx context.Context, event *events.Event) error {
	var payload events.IndexStarterPayload

	if err := event.UnmarshalData(&payload); err != nil {
		log.Printf("Failed to unmarshal leader assignment event: %v", err)
		return err
	}

	if err := h.searchRepo.DeleteFromIndex(ctx, payload.Domain); err != nil {
		return fmt.Errorf("failed to delete starter from Elasticsearch: %w", err)
	}

	log.Printf("Successfully deleted starter from index: domain=%s", event.Domain)
	return nil
}
