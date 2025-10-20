package service

import (
	"context"
	"strconv"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type StarterSearchService struct {
	searchRepo    repository.StarterSearchRepository
	repo          repository.StarterRepository
	kafkaProducer *sharedKafka.Producer
}

func NewStarterSearchService(
	searchRepo repository.StarterSearchRepository,
	repo repository.StarterRepository,
	kafkaProducer *sharedKafka.Producer,
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
) (*httputil.PaginatedResult[*model.Starter], error) {

	sortBy := query.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := query.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}

	filter := model.StarterListFilter{
		DepartmentID:   query.DepartmentID,
		BusinessUnitID: query.BusinessUnitID,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
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

	return &httputil.PaginatedResult[*model.Starter]{
		Data: starters,
		Pagination: httputil.RespPagination{
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

	event, err := events.NewEvent(
		events.EventTypeStarterIndex,
		events.WithDomain(starter.Domain()),
	)
	if err != nil {
		return err
	}

	_, _, err = s.kafkaProducer.PublishEvent(event)
	return err
}

// DeleteFromIndex publishes a starter deletion event to Kafka (call after SoftDelete)
func (s *StarterSearchService) DeleteFromIndex(ctx context.Context, domain string) error {
	if s.kafkaProducer == nil {
		return nil
	}

	event, err := events.NewEvent(
		events.EventTypeStarterDelete,
		events.WithDomain(domain),
	)
	if err != nil {
		return err
	}

	_, _, err = s.kafkaProducer.PublishEvent(event)
	return err
}

// ReindexAll reindexes all starters (for initial setup or data migration)
func (s *StarterSearchService) ReindexAll(ctx context.Context) error {
	const batchSize = 100
	page := 1
	totalIndexed := 0

	for {
		filter := model.StarterListFilter{}
		pagination := httputil.ReqPagination{Page: page, Limit: batchSize}

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
