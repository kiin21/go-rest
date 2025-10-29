package service

import (
	"context"
	"log"
	"strconv"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type StarterSearchService struct {
	searchRepo   repository.StarterSearchRepository
	repo         repository.StarterRepository
	syncProducer messaging.SyncProducer
}

func NewStarterSearchService(
	searchRepo repository.StarterSearchRepository,
	repo repository.StarterRepository,
	syncProducer messaging.SyncProducer,
) *StarterSearchService {
	return &StarterSearchService{
		searchRepo:   searchRepo,
		repo:         repo,
		syncProducer: syncProducer,
	}
}

func (s *StarterSearchService) IndexStarter(ctx context.Context, starter *model.Starter) error {
	if s.syncProducer == nil {
		return nil
	}

	// Create an event payload
	payload := events.IndexStarterPayload{
		StarterID: starter.ID,
		Domain:    starter.Domain,
		Name:      starter.Name,
	}

	event, err := events.NewEvent(events.EventTypeStarterIndex, payload)
	if err != nil {
		log.Printf("Failed to create starter index event: %v", err)
		return err
	}

	if err := s.syncProducer.SendSyncEvent(event); err != nil {
		log.Printf("Failed to send starter index event: %v", err)
		return err
	}

	return nil
}

func (s *StarterSearchService) DeleteFromIndex(ctx context.Context, starter *model.Starter) error {
	if s.syncProducer == nil {
		return nil
	}

	// Create an event payload
	payload := events.IndexStarterPayload{
		StarterID: starter.ID,
		Domain:    starter.Domain,
		Name:      starter.Name,
	}

	event, err := events.NewEvent(events.EventTypeStarterDelete, payload)
	if err != nil {
		log.Printf("Failed to create starter delete event: %v", err)
		return err
	}

	if err := s.syncProducer.SendSyncEvent(event); err != nil {
		log.Printf("Failed to send starter delete event: %v", err)
		return err
	}

	return nil
}

func (s *StarterSearchService) Search(
	ctx context.Context,
	query *starterquery.ListStartersQuery,
) (*httputil.PaginatedResult[*model.Starter], error) {
	// Elasticsearch search with query builder
	starterIds, total, err := s.searchRepo.Search(ctx, query)
	if err != nil {
		return nil, err
	}

	starters, err := s.repo.FindByIDs(ctx, starterIds)

	totalPages := int(total) / query.Pagination.GetLimit()
	if int(total)%(query.Pagination.GetLimit()) > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.GetPage() > 1 {
		value := strconv.Itoa(query.Pagination.GetPage() - 1)
		prev = &value
	}
	if query.Pagination.GetPage() < totalPages {
		value := strconv.Itoa(query.Pagination.GetPage() + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*model.Starter]{
		Data: starters,
		Pagination: httputil.RespPagination{
			Limit:      query.Pagination.GetLimit(),
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}
