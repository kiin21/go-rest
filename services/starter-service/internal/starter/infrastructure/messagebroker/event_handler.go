package messagebroker

import (
	"context"
	"fmt"
	"log"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

// SyncEventHandler handles synchronization events between MySQL and Elasticsearch
type SyncEventHandler struct {
	starterRepo repository.StarterRepository
	searchRepo  repository.StarterSearchRepository
}

// NewSyncEventHandler creates a new sync event handler
func NewSyncEventHandler(
	starterRepo repository.StarterRepository,
	searchRepo repository.StarterSearchRepository,
) *SyncEventHandler {
	return &SyncEventHandler{
		starterRepo: starterRepo,
		searchRepo:  searchRepo,
	}
}

// Handle processes incoming sync events
func (h *SyncEventHandler) Handle(ctx context.Context, event *events.Event) error {
	log.Printf("Processing event: type=%s, domain=%s, timestamp=%s",
		event.Type, event.Domain, event.Timestamp)

	switch event.Type {
	case events.EventTypeStarterIndex, events.EventTypeStarterUpdate, events.EventTypeStarterInsert:
		return h.handleIndexEvent(ctx, event)
	case events.EventTypeStarterDelete:
		return h.handleDeleteEvent(ctx, event)
	default:
		log.Printf("Warning: unknown event type: %s", event.Type)
		return nil // Don't fail on unknown events
	}
}

// handleIndexEvent handles index/update/insert events
func (h *SyncEventHandler) handleIndexEvent(ctx context.Context, event *events.Event) error {
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
func (h *SyncEventHandler) handleDeleteEvent(ctx context.Context, event *events.Event) error {
	if err := h.searchRepo.DeleteFromIndex(ctx, event.Domain); err != nil {
		return fmt.Errorf("failed to delete starter from Elasticsearch: %w", err)
	}

	log.Printf("Successfully deleted starter from index: domain=%s", event.Domain)
	return nil
}
