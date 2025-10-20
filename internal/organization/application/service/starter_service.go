package service

import (
	"context"
	"errors"
	"log"
	"strconv"

	startercommand "github.com/kiin21/go-rest/internal/organization/application/dto/starter/command"
	starterquery "github.com/kiin21/go-rest/internal/organization/application/dto/starter/query"
	"github.com/kiin21/go-rest/internal/organization/domain/model"
	repo "github.com/kiin21/go-rest/internal/organization/domain/repository"
	domainService "github.com/kiin21/go-rest/internal/organization/domain/service"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterApplicationService struct {
	repo              repo.StarterRepository
	domainService     *domainService.StarterDomainService
	enrichmentService *domainService.StarterEnrichmentService
	searchService     *StarterSearchService
}

func NewStarterApplicationService(
	repo repo.StarterRepository,
	domainService *domainService.StarterDomainService,
	enrichmentService *domainService.StarterEnrichmentService,
	searchService *StarterSearchService,
) *StarterApplicationService {
	return &StarterApplicationService{
		repo:              repo,
		domainService:     domainService,
		searchService:     searchService,
		enrichmentService: enrichmentService,
	}
}

func (s *StarterApplicationService) ListStarters(
	ctx context.Context,
	query starterquery.ListStartersQuery,
) (*response.PaginatedResult[*model.Starter], error) {
	sortBy := query.SortBy
	if sortBy == "" {
		sortBy = "id"
	}
	sortOrder := query.SortOrder
	if sortOrder == "" {
		sortOrder = "asc"
	}

	if query.Keyword != "" && s.searchService != nil {
		searchQuery := starterquery.SearchStartersQuery{
			Keyword:        query.Keyword,
			DepartmentID:   query.DepartmentID,
			BusinessUnitID: query.BusinessUnitID,
			Pagination:     query.Pagination,
			SortBy:         sortBy,
			SortOrder:      sortOrder,
		}
		return s.searchService.Search(ctx, searchQuery)
	}

	filter := model.StarterListFilter{
		DepartmentID:   query.DepartmentID,
		BusinessUnitID: query.BusinessUnitID,
		SortBy:         sortBy,
		SortOrder:      sortOrder,
	}

	var (
		starters []*model.Starter
		total    int64
		err      error
	)

	if query.Keyword != "" {
		starters, total, err = s.repo.SearchByKeyword(ctx, query.Keyword, filter, query.Pagination)
	} else {
		starters, total, err = s.repo.List(ctx, filter, query.Pagination)
	}
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

func (s *StarterApplicationService) CreateStarter(ctx context.Context, command startercommand.CreateStarterCommand) (*model.Starter, error) {
	if err := s.domainService.ValidateDomainUniqueness(ctx, command.Domain); err != nil {
		if errors.Is(err, sharedDomain.ErrDomainAlreadyExists) {
			return nil, err
		}
		return nil, err
	}

	starter, err := model.NewStarter(
		command.Domain,
		command.Name,
		command.Email,
		command.Mobile,
		command.WorkPhone,
		command.JobTitle,
		command.DepartmentID,
		command.LineManagerID,
	)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Create(ctx, starter); err != nil {
		return nil, err
	}

	if s.searchService != nil {
		go func() {
			if err := s.searchService.IndexStarter(context.Background(), starter); err != nil {
				log.Printf("Failed to index starter to Elasticsearch: %v", err)
			}
		}()
	}

	return starter, nil
}

func (s *StarterApplicationService) GetStarterByDomain(ctx context.Context, domainName string) (*model.Starter, error) {
	starter, err := s.repo.FindByDomain(ctx, domainName)
	if err != nil {
		return nil, err
	}

	return starter, nil
}

func (s *StarterApplicationService) UpdateStarter(ctx context.Context, command startercommand.UpdateStarterCommand) (*model.Starter, error) {
	starter, err := s.repo.FindByDomain(ctx, command.Domain)
	if err != nil {
		return nil, err
	}

	name := starter.Name()
	if command.Name != nil {
		name = *command.Name
	}

	email := starter.Email()
	if command.Email != nil {
		email = *command.Email
	}

	mobile := starter.Mobile()
	if command.Mobile != nil {
		mobile = *command.Mobile
	}

	workPhone := starter.WorkPhone()
	if command.WorkPhone != nil {
		workPhone = *command.WorkPhone
	}

	jobTitle := starter.JobTitle()
	if command.JobTitle != nil {
		jobTitle = *command.JobTitle
	}

	departmentID := starter.DepartmentID()
	if command.DepartmentID != nil {
		departmentID = command.DepartmentID
	}

	lineManagerID := starter.LineManagerID()
	if command.LineManagerID != nil {
		lineManagerID = command.LineManagerID
	}

	if err := starter.UpdateInfo(name, email, mobile, workPhone, jobTitle, departmentID, lineManagerID); err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, starter); err != nil {
		return nil, err
	}

	if s.searchService != nil {
		go func() {
			if err := s.searchService.IndexStarter(context.Background(), starter); err != nil {
				log.Printf("Failed to index starter to Elasticsearch: %v", err)
			}
		}()
	}

	return starter, nil
}

// SoftDeleteStarter soft deletes a starter by domain
func (s *StarterApplicationService) SoftDeleteStarter(ctx context.Context, domain string) error {
	// Soft delete from MySQL
	if err := s.repo.SoftDelete(ctx, domain); err != nil {
		return err
	}

	// Remove from Elasticsearch index (async, non-blocking)
	if s.searchService != nil {
		go func() {
			if err := s.searchService.DeleteFromIndex(context.Background(), domain); err != nil {
				log.Printf("Failed to delete starter from Elasticsearch: %v", err)
			}
		}()
	}

	return nil
}

func (s *StarterApplicationService) EnrichStarters(ctx context.Context, starters []*model.Starter) (*model.EnrichedData, error) {
	return s.enrichmentService.EnrichStarters(ctx, starters)
}
