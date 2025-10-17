package application

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/kiin21/go-rest/internal/shared/domain"
	appDto "github.com/kiin21/go-rest/internal/starter/application/dto"
	starterAggregate "github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	starterPort "github.com/kiin21/go-rest/internal/starter/domain/port"
	starterService "github.com/kiin21/go-rest/internal/starter/domain/service"
	"github.com/kiin21/go-rest/internal/starter/presentation/http/dto"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterApplicationService struct {
	repo              starterPort.StarterRepository
	domainService     *starterService.StarterDomainService
	searchService     *StarterSearchService
	enrichmentService *starterService.StarterEnrichmentService
}

func NewStarterApplicationService(
	repo starterPort.StarterRepository,
	domainService *starterService.StarterDomainService,
	searchService *StarterSearchService,
	enrichmentService *starterService.StarterEnrichmentService,
) *StarterApplicationService {
	return &StarterApplicationService{
		repo:              repo,
		domainService:     domainService,
		searchService:     searchService,
		enrichmentService: enrichmentService,
	}
}

func (s *StarterApplicationService) GetAllStarters(
	ctx context.Context,
	query appDto.ListStartersQuery,
) (*response.PaginatedResult[*starterAggregate.Starter], error) {
	// If keyword exists and Elasticsearch is available → Use Elasticsearch
	if query.Keyword != "" && s.searchService != nil {
		searchQuery := appDto.SearchStartersQuery{
			Keyword:        query.Keyword,
			DepartmentID:   query.DepartmentID,
			BusinessUnitID: query.BusinessUnitID,
			Pagination:     query.Pagination,
		}
		return s.searchService.Search(ctx, searchQuery)
	}

	// Fallback to MySQL (keyword search or list)
	filter := starterPort.ListFilter{
		DepartmentID:   query.DepartmentID,
		BusinessUnitID: query.BusinessUnitID,
	}

	var starters []*starterAggregate.Starter
	var total int64
	var err error

	if query.Keyword != "" {
		// Elasticsearch not available, use MySQL LIKE search
		starters, total, err = s.repo.SearchByKeyword(ctx, query.Keyword, filter, query.Pagination)
	} else {
		// No keyword, simple list
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

// CreateStarter creates a new starter and syncs to Elasticsearch
func (s *StarterApplicationService) CreateStarter(ctx context.Context, command appDto.CreateStarterCommand) (*starterAggregate.Starter, error) {
	if err := s.domainService.ValidateDomainUniqueness(ctx, command.Domain); err != nil {
		if errors.Is(err, domain.ErrDomainAlreadyExists) {
			return nil, fmt.Errorf("domain '%s' already exists: %w", command.Domain, domain.ErrDomainAlreadyExists)
		}
		return nil, err
	}

	// Create new starter entity
	starter, err := starterAggregate.NewStarter(
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

	// Persist to MySQL
	if err := s.repo.Create(ctx, starter); err != nil {
		return nil, err
	}

	// Sync to Elasticsearch
	if s.searchService != nil {
		go func() {
			if err := s.searchService.IndexStarter(context.Background(), starter); err != nil {
				log.Printf("Failed to index starter to Elasticsearch: %v", err)
			}
		}()
	}

	return starter, nil
}

// GetStarterByDomain returns a single starter by domain (username)
func (s *StarterApplicationService) GetStarterByDomain(ctx context.Context, domainName string) (*starterAggregate.Starter, error) {
	starter, err := s.repo.FindByDomain(ctx, domainName)
	if err != nil {
		return nil, err
	}

	return starter, nil
}

// UpdateStarter updates an existing starter (supports partial updates)
func (s *StarterApplicationService) UpdateStarter(ctx context.Context, command appDto.UpdateStarterCommand) (*starterAggregate.Starter, error) {
	// Find existing starter
	starter, err := s.repo.FindByDomain(ctx, command.Domain)
	if err != nil {
		return nil, err
	}

	// Update basic info only if provided (partial update)
	// Build current values, only override if new value is provided
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

	// Update starter info (will validate required fields)
	starter.UpdateInfo(name, email, mobile, workPhone, jobTitle, departmentID, lineManagerID)

	// Persist changes
	if err := s.repo.Update(ctx, starter); err != nil {
		return nil, err
	}

	// Sync to Elasticsearch (async, non-blocking)
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

// EnrichStarters delegates enrichment to the domain service and converts to presentation DTO
func (s *StarterApplicationService) EnrichStarters(ctx context.Context, starters []*starterAggregate.Starter) (*dto.StarterEnrichedData, error) {
	// Get domain enriched data
	domainEnriched, err := s.enrichmentService.EnrichStarters(ctx, starters)
	if err != nil {
		return nil, err
	}

	// Convert to presentation DTO
	interfaceEnriched := &dto.StarterEnrichedData{
		Departments:   make(map[int64]*dto.DepartmentNested),
		LineManagers:  make(map[int64]*dto.LineManagerNested),
		BusinessUnits: make(map[int64]*dto.BusinessUnitNested),
	}

	if domainEnriched == nil {
		return interfaceEnriched, nil
	}

	// Map departments
	for id, dept := range domainEnriched.Departments {
		deptDTO := &dto.DepartmentNested{
			ID:        dept.ID,
			Name:      dept.Name,
			Shortname: dept.Shortname,
		}
		if dept.GroupDepartment != nil {
			deptDTO.GroupDepartment = &dto.GroupDepartmentNested{
				ID:        dept.GroupDepartment.ID,
				Name:      dept.GroupDepartment.Name,
				Shortname: dept.GroupDepartment.Shortname,
			}
		}
		interfaceEnriched.Departments[id] = deptDTO
	}

	// Map line managers
	for id, manager := range domainEnriched.LineManagers {
		interfaceEnriched.LineManagers[id] = &dto.LineManagerNested{
			ID:       manager.ID,
			Domain:   manager.Domain,
			Name:     manager.Name,
			Email:    manager.Email,
			JobTitle: manager.JobTitle,
		}
	}

	// Map business units
	for id, bu := range domainEnriched.BusinessUnits {
		interfaceEnriched.BusinessUnits[id] = &dto.BusinessUnitNested{
			ID:        bu.ID,
			Name:      bu.Name,
			Shortname: bu.Shortname,
		}
	}

	return interfaceEnriched, nil
}
