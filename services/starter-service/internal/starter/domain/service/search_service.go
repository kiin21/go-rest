package service

import (
	"context"
	"log"
	"strconv"
	"strings"

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
	syncProducer  messaging.SyncProducer
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
		StarterID: starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
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
		StarterID: starter.ID(),
		Domain:    starter.Domain(),
		Name:      starter.Name(),
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

func (s *StarterSearchService) buildSearchQuery(q *starterquery.ListStartersQuery) map[string]interface{} {
	if q == nil {
		return nil
	}
	kw := strings.TrimSpace(q.Keyword)
	if kw == "" {
		return nil
	}

	field := s.mapSearchByToFieldName(q.SearchBy)

	var must []any
	if field == "" {
		must = append(must, map[string]any{
			"multi_match": map[string]any{
				"query": kw,
				"type":  "bool_prefix",
				"fields": []string{
					"domain^3",
					"name^2",
					"department_name",
					"business_unit_name",
					"full_text",
					"search_tokens",
				},
			},
		})
	} else if field == "domain" {
		must = append(must, map[string]any{
			"wildcard": map[string]any{
				"domain": map[string]any{"value": "*" + kw + "*"},
			},
		})
	} else {
		must = append(must, map[string]any{
			"match_phrase_prefix": map[string]any{
				field: map[string]any{"query": kw},
			},
		})
	}

	es := map[string]any{
		"query": map[string]any{
			"bool": map[string]any{
				"must": must,
			},
		},
	}

	if sort := s.buildSortClause(q.SortBy, q.SortOrder); len(sort) > 0 {
		es["sort"] = sort
	}
	return es
}

func (s *StarterSearchService) mapSearchByToFieldName(searchBy string) string {
	switch strings.ToLower(strings.TrimSpace(searchBy)) {
	case "domain":
		return "domain"
	case "fullname", "name":
		return "name"
	case "dept_name":
		return "department_name"
	case "bu_name":
		return "business_unit_name"
	default:
		return ""
	}
}

func (s *StarterSearchService) buildSortClause(sortBy, sortOrder string) []interface{} {
	field := s.mapSortFieldToESField(sortBy)
	if field == "" {
		return nil
	}
	order := strings.ToLower(strings.TrimSpace(sortOrder))
	if order != "desc" {
		order = "asc"
	}
	return []interface{}{
		map[string]interface{}{
			field: map[string]interface{}{"order": order},
		},
	}
}

func (s *StarterSearchService) mapSortFieldToESField(sortBy string) string {
	switch strings.ToLower(strings.TrimSpace(sortBy)) {
	case "id":
		return "id"
	case "domain":
		return "domain"
	case "name", "fullname":
		return "name"
	case "dept_name":
		return "department_name"
	case "bu_name":
		return "business_unit_name"
	default:
		return "id"
	}
}
