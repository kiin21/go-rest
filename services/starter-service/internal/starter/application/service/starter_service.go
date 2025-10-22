package service

import (
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/kiin21/go-rest/pkg/httputil"
	startercommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/command"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
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
) (*httputil.PaginatedResult[*model.Starter], error) {

	// Use Elasticsearch if a keyword is provided and Elasticsearch is enabled
	if query.Keyword != "" && s.searchService != nil {
		log.Println("Using Elasticsearch for search")
		return s.searchService.Search(ctx, query)
	}

	// Use MySQL if a keyword is provided
	// TODO: refactor
	var (
		starters []*model.Starter
		total    int64
		err      error
	)
	log.Printf("Using MySQL for search: keyword=%s", query.Keyword)
	if query.Keyword != "" {
		starters, total, err = s.repo.SearchByKeyword(ctx, query)
	} else {
		starters, total, err = s.repo.SearchByKeyword(ctx, query)
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
	entity, err := s.repo.SoftDelete(ctx, domain)
	if err != nil {
		return err
	}

	// Remove from Elasticsearch index (async, non-blocking)
	if s.searchService != nil {
		go func() {
			if err := s.searchService.DeleteFromIndex(context.Background(), entity); err != nil {
				log.Printf("Failed to delete starter from Elasticsearch: %v", err)
			}
		}()
	}

	return nil
}

func (s *StarterApplicationService) EnrichStarters(ctx context.Context, starters []*model.Starter) (*model.EnrichedData, error) {
	return s.enrichmentService.EnrichStarters(ctx, starters)
}
