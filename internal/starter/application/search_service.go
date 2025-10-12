package application

import (
	"context"
	"strconv"
	"time"

	"github.com/kiin21/go-rest/internal/shared/infrastructure"
	appDto "github.com/kiin21/go-rest/internal/starter/application/dto"
	"github.com/kiin21/go-rest/internal/starter/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

// StarterSearchService handles search operations using Elasticsearch
// Separated from StarterApplicationService for clear responsibility
type StarterSearchService struct {
	searchRepo    domain.StarterSearchRepository
	repo          domain.StarterRepository // For fetching full data if needed
	kafkaProducer *infrastructure.KafkaProducer
}

// NewStarterSearchService creates a new search service
func NewStarterSearchService(
	searchRepo domain.StarterSearchRepository,
	repo domain.StarterRepository,
	kafkaProducer *infrastructure.KafkaProducer,
) *StarterSearchService {
	return &StarterSearchService{
		searchRepo:    searchRepo,
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

// Search performs search using Elasticsearch
func (s *StarterSearchService) Search(
	ctx context.Context,
	query appDto.SearchStartersQuery,
) (*response.PaginatedResult[*domain.Starter], error) {

	filter := domain.ListFilter{
		DepartmentID:   query.DepartmentID,
		BusinessUnitID: query.BusinessUnitID,
	}

	// Use Elasticsearch for search
	starters, total, err := s.searchRepo.Search(ctx, query.Query, filter, query.Pagination)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Pagination.Limit
	if int(total)%query.Pagination.Limit > 0 {
		totalPages++
	}

	var prev, next *string
	if query.Pagination.Page > 1 {
		value := strconv.Itoa(query.Pagination.Page - 1)
		prev = &value
	}
	if query.Pagination.Page < totalPages {
		value := strconv.Itoa(query.Pagination.Page + 1)
		next = &value
	}

	return &response.PaginatedResult[*domain.Starter]{
		Data: starters,
		Pagination: response.RespPagination{
			Limit:      query.Pagination.Limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

// IndexStarter publishes a starter indexing event to Kafka (call after Create/Update)
func (s *StarterSearchService) IndexStarter(ctx context.Context, starter *domain.Starter) error {
	if s.kafkaProducer == nil {
		// Kafka not configured, skip sync
		return nil
	}

	event := &infrastructure.SyncEvent{
		Type:      "index",
		Domain:    starter.Domain(),
		Data:      nil, // Don't send full object, consumer will fetch from DB
		Timestamp: time.Now(),
		Retries:   0,
	}

	return s.kafkaProducer.PublishSyncEvent(ctx, event)
}

// DeleteFromIndex publishes a starter deletion event to Kafka (call after SoftDelete)
func (s *StarterSearchService) DeleteFromIndex(ctx context.Context, domain string) error {
	if s.kafkaProducer == nil {
		// Kafka not configured, skip sync
		return nil
	}

	event := &infrastructure.SyncEvent{
		Type:      "delete",
		Domain:    domain,
		Data:      nil, // Only need domain for deletion
		Timestamp: time.Now(),
		Retries:   0,
	}

	return s.kafkaProducer.PublishSyncEvent(ctx, event)
}

// ReindexAll reindexes all starters (for initial setup or data migration)
func (s *StarterSearchService) ReindexAll(ctx context.Context) error {
	// Get all starters from database in batches
	const batchSize = 100
	page := 1
	totalIndexed := 0

	for {
		filter := domain.ListFilter{}
		pagination := response.ReqPagination{Page: page, Limit: batchSize}

		starters, total, err := s.repo.List(ctx, filter, pagination)
		if err != nil {
			return err
		}

		if len(starters) == 0 {
			break
		}

		// Bulk index to Elasticsearch
		if err := s.searchRepo.BulkIndex(ctx, starters); err != nil {
			return err
		}

		totalIndexed += len(starters)

		// Check if we've processed all records
		if int64(totalIndexed) >= total {
			break
		}

		page++
	}

	return nil
}
