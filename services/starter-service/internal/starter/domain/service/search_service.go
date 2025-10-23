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
	searchRepo    repository.StarterSearchRepository
	repo          repository.StarterRepository
	kafkaProducer messaging.NotificationProducer
}

func NewStarterSearchService(
	searchRepo repository.StarterSearchRepository,
	repo repository.StarterRepository,
	kafkaProducer messaging.NotificationProducer,
) *StarterSearchService {
	return &StarterSearchService{
		searchRepo:    searchRepo,
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

// IndexStarter publishes a starter indexing event to Kafka (call after Create/Update)
func (s *StarterSearchService) IndexStarter(ctx context.Context, starter *model.Starter) error {
	if s.kafkaProducer == nil {
		return nil
	}

	// Create an event payload
	payload := events.IndexStarterPayload{
		StarterID: starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
	}

	event, err := events.NewEvent(events.EventTypeStarterIndex, payload)
	if err != nil {
		log.Printf("Failed to create starter index event: %v", err)
		return err
	}

	if err := s.kafkaProducer.SendNotification(event); err != nil {
		log.Printf("Failed to send starter index event: %v", err)
		return err
	}

	return nil
}

// DeleteFromIndex publishes a starter deletion event to Kafka (call after SoftDelete)
func (s *StarterSearchService) DeleteFromIndex(ctx context.Context, starter *model.Starter) error {
	if s.kafkaProducer == nil {
		return nil
	}

	// Create an event payload
	payload := events.IndexStarterPayload{
		StarterID: starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
	}

	event, err := events.NewEvent(events.EventTypeStarterDelete, payload)
	if err != nil {
		log.Printf("Failed to create starter delete event: %v", err)
		return err
	}

	if err := s.kafkaProducer.SendNotification(event); err != nil {
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
	starterIds, total, err := s.searchRepo.Search(ctx, query, s.buildSearchQuery)
	if err != nil {
		return nil, err
	}

	starters, err := s.repo.FindByIDs(ctx, starterIds)

	// TODO: refactor this kind of pagination
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

func (s *StarterSearchService) buildSearchQuery(*starterquery.ListStartersQuery) map[string]interface{} {
	// TODO: implement
	return nil
}

// mapSearchByToFieldName maps the SearchBy parameter to Elasticsearch field name
func (s *StarterSearchService) mapSearchByToFieldName(searchBy string) string {
	// TODO: implement
	return ""
}

// buildSortClause builds the sort clause for Elasticsearch
func (s *StarterSearchService) buildSortClause(sortBy, sortOrder string) []interface{} {
	// TODO: implement
	return nil
}

// mapSortFieldToESField maps domain sort fields to Elasticsearch fields
func (s *StarterSearchService) mapSortFieldToESField(sortBy string) string {
	// TODO: implement
	return ""
}
