package application

import (
	"context"
	"strconv"
	"time"

	messaging "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
	"github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker/kafka"
	appDto "github.com/kiin21/go-rest/internal/starter/application/dto"
	starterAggregate "github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	starterPort "github.com/kiin21/go-rest/internal/starter/domain/port"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterSearchService struct {
	searchRepo    starterPort.StarterSearchRepository
	repo          starterPort.StarterRepository
	kafkaProducer *kafka.Producer
}

func NewStarterSearchService(
	searchRepo starterPort.StarterSearchRepository,
	repo starterPort.StarterRepository,
	kafkaProducer *kafka.Producer,
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
) (*response.PaginatedResult[*starterAggregate.Starter], error) {

	filter := starterPort.ListFilter{
		DepartmentID:   query.DepartmentID,
		BusinessUnitID: query.BusinessUnitID,
	}

	// Elasticsearch
	starters, total, err := s.searchRepo.Search(ctx, query.Keyword, filter, query.Pagination)
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

	return &response.PaginatedResult[*starterAggregate.Starter]{
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
func (s *StarterSearchService) IndexStarter(ctx context.Context, starter *starterAggregate.Starter) error {
	if s.kafkaProducer == nil {
		// Kafka not configured, skip sync
		return nil
	}

	event := &messaging.SyncEvent{
		Type:      "index",
		Domain:    starter.Domain(),
		Data:      nil,
		Timestamp: time.Now(),
		Retries:   0,
	}

	return s.kafkaProducer.PublishSyncEvent(event)
}

// DeleteFromIndex publishes a starter deletion event to Kafka (call after SoftDelete)
func (s *StarterSearchService) DeleteFromIndex(ctx context.Context, domain string) error {
	if s.kafkaProducer == nil {
		// Kafka not configured, skip sync
		return nil
	}

	event := &messaging.SyncEvent{
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
		filter := starterPort.ListFilter{}
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
