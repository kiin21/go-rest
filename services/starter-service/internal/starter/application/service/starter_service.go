package service

import (
	"context"
	"errors"
	"fmt"
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
	starterRepo       repo.StarterRepository
	searchRepo        repo.StarterSearchRepository
	domainService     *domainService.StarterDomainService
	enrichmentService *domainService.StarterEnrichmentService
	searchService     *domainService.StarterSearchService
}

func NewStarterApplicationService(
	starterRepo repo.StarterRepository,
	searchRepo repo.StarterSearchRepository,
	domainService *domainService.StarterDomainService,
	enrichmentService *domainService.StarterEnrichmentService,
	searchService *domainService.StarterSearchService,
) *StarterApplicationService {
	return &StarterApplicationService{
		starterRepo:       starterRepo,
		searchRepo:        searchRepo,
		domainService:     domainService,
		searchService:     searchService,
		enrichmentService: enrichmentService,
	}
}

func (s *StarterApplicationService) ListStarters(
	ctx context.Context,
	query *starterquery.ListStartersQuery,
) (*httputil.PaginatedResult[*model.Starter], error) {
	// Use Elasticsearch if keyword exists and search service is available
	if query.Keyword != "" && s.searchService != nil {
		log.Println("Using Elasticsearch for search")
		fmt.Println("QUERY: ", query.SearchBy)
		return s.searchService.Search(ctx, query)
	}

	// Fallback to MySQL
	log.Printf("Using MySQL for search: keyword=%s, by=%s", query.Keyword, query.SearchBy)
	return s.listFromMySQL(ctx, query)
}

func (s *StarterApplicationService) CreateStarter(
	ctx context.Context,
	command *startercommand.CreateStarterCommand,
) (*model.Starter, error) {
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

	if err := s.starterRepo.Create(ctx, starter); err != nil {
		return nil, err
	}

	s.asyncIndexStarter(starter)
	return starter, nil
}

func (s *StarterApplicationService) GetStarterByDomain(
	ctx context.Context,
	domainName string,
) (*model.Starter, error) {
	return s.starterRepo.FindByDomain(ctx, domainName)
}

func (s *StarterApplicationService) UpdateStarter(
	ctx context.Context,
	command *startercommand.UpdateStarterCommand,
) (*model.Starter, error) {
	starter, err := s.starterRepo.FindByDomain(ctx, command.OriginalDomain)
	if err != nil {
		return nil, err
	}

	domain, name, email, mobile, workPhone, jobTitle, departmentID, lineManagerID := s.applyUpdates(starter, command)

	if err := starter.UpdateInfo(domain, name, email, mobile, workPhone, jobTitle, departmentID, lineManagerID); err != nil {
		return nil, err
	}

	if err := s.starterRepo.Update(ctx, starter); err != nil {
		return nil, err
	}

	s.asyncIndexStarter(starter)
	return starter, nil
}

func (s *StarterApplicationService) SoftDeleteStarter(
	ctx context.Context,
	domain string,
) error {
	entity, err := s.starterRepo.SoftDelete(ctx, domain)
	if err != nil {
		return err
	}

	s.asyncDeleteFromIndex(entity)
	return nil
}

func (s *StarterApplicationService) ReindexAll(ctx context.Context) error {
	const batchSize = 100
	var totalIndexed int

	for page := 1; ; page++ {
		query := &starterquery.ListStartersQuery{
			Pagination: httputil.ReqPagination{
				Page:  &page,
				Limit: func() *int { v := batchSize; return &v }(),
			},
		}

		starters, total, err := s.starterRepo.SearchByKeyword(ctx, query)
		if err != nil {
			return fmt.Errorf("failed to fetch starters page %d: %w", page, err)
		}

		if len(starters) == 0 {
			break
		}

		enriched, err := s.enrichmentService.EnrichStarters(ctx, starters)
		if err != nil {
			return fmt.Errorf("failed to enrich starters batch: %w", err)
		}

		esDocs := make([]*model.StarterESDoc, len(starters))
		for i, starter := range starters {
			esDocs[i] = model.NewStarterESDocFromStarter(starter, enriched)
		}

		if err := s.searchRepo.BulkIndex(ctx, esDocs); err != nil {
			return fmt.Errorf("failed to bulk index batch: %w", err)
		}

		totalIndexed += len(starters)
		log.Printf("Reindexed %d/%d starters", totalIndexed, total)

		if int64(totalIndexed) >= total {
			break
		}
	}

	log.Printf("Reindexing completed: %d starters indexed", totalIndexed)
	return nil
}

func (s *StarterApplicationService) asyncIndexStarter(starter *model.Starter) {
	if s.searchService == nil {
		return
	}

	go func() {
		if err := s.searchService.IndexStarter(context.Background(), starter); err != nil {
			log.Printf("Failed to index starter to Elasticsearch: %v", err)
		}
	}()
}

func (s *StarterApplicationService) asyncDeleteFromIndex(entity *model.Starter) {
	if s.searchService == nil {
		return
	}

	go func() {
		if err := s.searchService.DeleteFromIndex(context.Background(), entity); err != nil {
			log.Printf("Failed to delete starter from Elasticsearch: %v", err)
		}
	}()
}

func (s *StarterApplicationService) listFromMySQL(
	ctx context.Context,
	query *starterquery.ListStartersQuery,
) (*httputil.PaginatedResult[*model.Starter], error) {
	starters, total, err := s.starterRepo.SearchByKeyword(ctx, query)
	if err != nil {
		return nil, err
	}

	limit := query.Pagination.GetLimit()
	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	var prev, next *string
	currentPage := query.Pagination.GetPage()
	if currentPage > 1 {
		value := strconv.Itoa(currentPage - 1)
		prev = &value
	}
	if currentPage < totalPages {
		value := strconv.Itoa(currentPage + 1)
		next = &value
	}

	return &httputil.PaginatedResult[*model.Starter]{
		Data: starters,
		Pagination: httputil.RespPagination{
			Limit:      limit,
			TotalItems: total,
			Prev:       prev,
			Next:       next,
		},
	}, nil
}

func (s *StarterApplicationService) applyUpdates(
	starter *model.Starter,
	command *startercommand.UpdateStarterCommand,
) (domain, name, email, mobile, workPhone, jobTitle string, departmentID, lineManagerID *int64) {
	domain = starter.Domain
	if command.Domain != nil {
		domain = *command.Domain
	}

	name = starter.Name
	if command.Name != nil {
		name = *command.Name
	}

	email = starter.GetEmail()
	if command.Email != nil {
		email = *command.Email
	}

	mobile = starter.Mobile
	if command.Mobile != nil {
		mobile = *command.Mobile
	}

	workPhone = starter.WorkPhone
	if command.WorkPhone != nil {
		workPhone = *command.WorkPhone
	}

	jobTitle = starter.JobTitle
	if command.JobTitle != nil {
		jobTitle = *command.JobTitle
	}

	departmentID = starter.DepartmentID
	if command.DepartmentID != nil {
		departmentID = command.DepartmentID
	}

	lineManagerID = starter.LineManagerID
	if command.LineManagerID != nil {
		lineManagerID = command.LineManagerID
	}

	return
}
