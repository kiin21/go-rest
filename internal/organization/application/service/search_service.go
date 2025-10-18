package service

import (
	"context"
	"strconv"
	"time"

	starterquery "github.com/kiin21/go-rest/internal/organization/application/dto/starter/query"
	model "github.com/kiin21/go-rest/internal/organization/domain/model"
	repository "github.com/kiin21/go-rest/internal/organization/domain/repository"
	messaging "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
	"github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker/kafka"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterSearchService struct {
	searchRepo    repository.StarterSearchRepository
	repo          repository.StarterRepository
	kafkaProducer *kafka.Producer
}

func NewStarterSearchService(
	searchRepo repository.StarterSearchRepository,
	repo repository.StarterRepository,
	kafkaProducer *kafka.Producer,
) *StarterSearchService {
	return &StarterSearchService{
		searchRepo:    searchRepo,
		repo:          repo,
		kafkaProducer: kafkaProducer,
	}
}

func (s *StarterSearchService) Search(
	ctx context.Context,
	query starterquery.SearchStartersQuery,
) (*response.PaginatedResult[*model.Starter], error) {

	filter := model.StarterListFilter{
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

	return &response.PaginatedResult[*model.Starter]{
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
func (s *StarterSearchService) IndexStarter(ctx context.Context, starter *model.Starter) error {
	if s.kafkaProducer == nil {
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
		return nil
	}

	event := &messaging.SyncEvent{
		Type:      "delete",
		Domain:    domain,
		Data:      nil,
		Timestamp: time.Now(),
		Retries:   0,
	}

	return s.kafkaProducer.PublishSyncEvent(event)
}

// ReindexAll reindexes all starters (for initial setup or data migration)
func (s *StarterSearchService) ReindexAll(ctx context.Context) error {
	const batchSize = 100
	page := 1
	totalIndexed := 0

	for {
		filter := model.StarterListFilter{}
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

		if int64(totalIndexed) >= total {
			break
		}

		page++
	}

	return nil
}
