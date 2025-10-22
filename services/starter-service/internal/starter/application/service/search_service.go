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

func (s *StarterSearchService) Search(
	ctx context.Context,
	query starterquery.ListStartersQuery,
) (*httputil.PaginatedResult[*model.Starter], error) {
	// Elasticsearch
	starters, total, err := s.searchRepo.Search(ctx, query)
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

	// Create an event payload
	payload := events.IndexStarterPayload{
		StarterID: starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
	}

	event := events.NewEvent(events.EventTypeStarterIndex, payload)
	event.Key = starter.Domain()

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

	event := events.NewEvent(events.EventTypeStarterDelete, payload)
	event.Key = starter.Domain() // Use domain as partition key

	if err := s.kafkaProducer.SendNotification(event); err != nil {
		log.Printf("Failed to send starter delete event: %v", err)
		return err
	}

	return nil
}

// ReindexAll reindex all starters (for initial setup or data migration)
func (s *StarterSearchService) ReindexAll(ctx context.Context) error {
	const batchSize = 100
	page := 1
	totalIndexed := 0

	for {
		emptyQuery := starterquery.ListStartersQuery{
			Pagination: httputil.ReqPagination{
				Page: page, Limit: batchSize,
			},
		}

		starters, total, err := s.repo.SearchByKeyword(ctx, emptyQuery)
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
		log.Printf("Reindexed %d/%d starters", totalIndexed, total)

		if int64(totalIndexed) >= total {
			break
		}

		page++
	}

	log.Printf("✅ Reindexing completed: %d starters indexed", totalIndexed)
	return nil
}
