package service

import (
	"context"
	"fmt"

	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	"github.com/kiin21/go-rest/internal/starter/domain/model"
	"github.com/kiin21/go-rest/internal/starter/domain/port"
)

// StarterEnrichmentService enriches starter data with related entities.
type StarterEnrichmentService struct {
	starterRepo        port.StarterRepository
	organizationLookup port.OrganizationLookup
}

// NewStarterEnrichmentService creates enrichment service.
func NewStarterEnrichmentService(
	starterRepo port.StarterRepository,
	organizationLookup port.OrganizationLookup,
) *StarterEnrichmentService {
	return &StarterEnrichmentService{
		starterRepo:        starterRepo,
		organizationLookup: organizationLookup,
	}
}

// EnrichStarters loads and enriches starter data with departments, line managers, and business units.
func (s *StarterEnrichmentService) EnrichStarters(ctx context.Context, starters []*aggregate.Starter) (*model.EnrichedData, error) {
	enriched := &model.EnrichedData{
		Departments:   make(map[int64]*sharedDomain.DepartmentNested),
		LineManagers:  make(map[int64]*sharedDomain.LineManagerNested),
		BusinessUnits: make(map[int64]*sharedDomain.BusinessUnitNested),
	}

	// Collect unique IDs.
	departmentIDs := s.collectDepartmentIDs(starters)
	lineManagerIDs := s.collectLineManagerIDs(starters)

	// Load departments in batch.
	if len(departmentIDs) > 0 {
		if err := s.loadDepartments(ctx, departmentIDs, enriched); err != nil {
			return nil, err
		}
	}

	// Load line managers in batch.
	if len(lineManagerIDs) > 0 {
		if err := s.loadLineManagers(ctx, lineManagerIDs, enriched); err != nil {
			return nil, err
		}
	}

	return enriched, nil
}

// collectDepartmentIDs extracts unique department IDs from starters.
func (s *StarterEnrichmentService) collectDepartmentIDs(starters []*aggregate.Starter) map[int64]bool {
	departmentIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.DepartmentID() != nil {
			departmentIDs[*starter.DepartmentID()] = true
		}
	}
	return departmentIDs
}

// collectLineManagerIDs extracts unique line manager IDs from starters.
func (s *StarterEnrichmentService) collectLineManagerIDs(starters []*aggregate.Starter) map[int64]bool {
	lineManagerIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.LineManagerID() != nil {
			lineManagerIDs[*starter.LineManagerID()] = true
		}
	}
	return lineManagerIDs
}

// loadDepartments fetches departments with relations (group_department, business_unit).
func (s *StarterEnrichmentService) loadDepartments(
	ctx context.Context,
	departmentIDs map[int64]bool,
	enriched *model.EnrichedData,
) error {
	if s.organizationLookup == nil {
		return nil
	}

	// Convert map to slice.
	ids := make([]int64, 0, len(departmentIDs))
	for id := range departmentIDs {
		ids = append(ids, id)
	}

	// Batch load departments with relations.
	relations, err := s.organizationLookup.FindDepartmentRelations(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to load departments with relations: %w", err)
	}

	// Map results to enriched data.
	for _, rel := range relations {
		s.mapDepartmentRelation(rel, enriched)
	}

	return nil
}

// mapDepartmentRelation maps department relation to enriched data structure.
func (s *StarterEnrichmentService) mapDepartmentRelation(
	rel *model.DepartmentRelation,
	enriched *model.EnrichedData,
) {
	if rel == nil || rel.Department == nil {
		return
	}

	// Store department in enriched data.
	enriched.Departments[rel.Department.ID] = rel.Department

	// Store business unit if exists.
	if rel.BusinessUnit != nil {
		enriched.BusinessUnits[rel.Department.ID] = rel.BusinessUnit
	}
}

// loadLineManagers fetches line manager information.
func (s *StarterEnrichmentService) loadLineManagers(
	ctx context.Context,
	lineManagerIDs map[int64]bool,
	enriched *model.EnrichedData,
) error {
	// Convert map to slice.
	ids := make([]int64, 0, len(lineManagerIDs))
	for id := range lineManagerIDs {
		ids = append(ids, id)
	}

	// Batch load starters (line managers are also starters).
	for _, id := range ids {
		manager, err := s.starterRepo.FindByID(ctx, id)
		if err != nil {
			continue // Skip if not found.
		}

		enriched.LineManagers[manager.ID()] = &sharedDomain.LineManagerNested{
			ID:       manager.ID(),
			Domain:   manager.Domain(),
			Name:     manager.Name(),
			Email:    manager.Email(),
			JobTitle: manager.JobTitle(),
		}
	}

	return nil
}
