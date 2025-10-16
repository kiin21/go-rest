package domain

import (
	"context"
	"fmt"

	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
)

// StarterEnrichmentService handles enrichment of starter data with related entities
// This is a Domain Service because it coordinates between multiple aggregates (Starter, Department, BusinessUnit)
type StarterEnrichmentService struct {
	starterRepo      StarterRepository
	departmentRepo   orgDomain.DepartmentRepository
	businessUnitRepo orgDomain.BusinessUnitRepository
}

// NewStarterEnrichmentService creates enrichment service
func NewStarterEnrichmentService(
	starterRepo StarterRepository,
	departmentRepo orgDomain.DepartmentRepository,
	businessUnitRepo orgDomain.BusinessUnitRepository,
) *StarterEnrichmentService {
	return &StarterEnrichmentService{
		starterRepo:      starterRepo,
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
	}
}

// EnrichStarters loads and enriches starter data with departments, line managers, and business units
// This prevents N+1 queries by batch loading related entities
func (s *StarterEnrichmentService) EnrichStarters(ctx context.Context, starters []*Starter) (*EnrichedData, error) {
	enriched := &EnrichedData{
		Departments:   make(map[int64]*DepartmentNested),
		LineManagers:  make(map[int64]*LineManagerNested),
		BusinessUnits: make(map[int64]*BusinessUnitNested),
	}

	// Collect unique IDs
	departmentIDs := s.collectDepartmentIDs(starters)
	lineManagerIDs := s.collectLineManagerIDs(starters)

	// Load departments in batch
	if len(departmentIDs) > 0 {
		if err := s.loadDepartments(ctx, departmentIDs, enriched); err != nil {
			return nil, err
		}
	}

	// Load line managers in batch
	if len(lineManagerIDs) > 0 {
		if err := s.loadLineManagers(ctx, lineManagerIDs, enriched); err != nil {
			return nil, err
		}
	}

	return enriched, nil
}

// collectDepartmentIDs extracts unique department IDs from starters
func (s *StarterEnrichmentService) collectDepartmentIDs(starters []*Starter) map[int64]bool {
	departmentIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.DepartmentID() != nil {
			departmentIDs[*starter.DepartmentID()] = true
		}
	}
	return departmentIDs
}

// collectLineManagerIDs extracts unique line manager IDs from starters
func (s *StarterEnrichmentService) collectLineManagerIDs(starters []*Starter) map[int64]bool {
	lineManagerIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.LineManagerID() != nil {
			lineManagerIDs[*starter.LineManagerID()] = true
		}
	}
	return lineManagerIDs
}

// loadDepartments fetches departments with relations (group_department, business_unit)
func (s *StarterEnrichmentService) loadDepartments(
	ctx context.Context,
	departmentIDs map[int64]bool,
	enriched *EnrichedData,
) error {
	// Convert map to slice
	ids := make([]int64, 0, len(departmentIDs))
	for id := range departmentIDs {
		ids = append(ids, id)
	}

	// Batch load departments with relations
	relations, err := s.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to load departments with relations: %w", err)
	}

	// Map results to enriched data
	for _, rel := range relations {
		s.mapDepartmentRelation(rel, enriched)
	}

	return nil
}

// mapDepartmentRelation maps department relation to enriched data structure
func (s *StarterEnrichmentService) mapDepartmentRelation(
	rel *orgDomain.DepartmentWithDetails,
	enriched *EnrichedData,
) {
	dept := rel.Department

	// Create department nested structure
	deptNested := &DepartmentNested{
		ID:        dept.ID,
		Name:      dept.FullName,
		Shortname: dept.Shortname,
	}

	// Add group department if exists
	if rel.ParentDepartment != nil {
		deptNested.GroupDepartment = &GroupDepartmentNested{
			ID:        rel.ParentDepartment.ID,
			Name:      rel.ParentDepartment.FullName,
			Shortname: rel.ParentDepartment.Shortname,
		}
	}

	// Store department in enriched data
	enriched.Departments[dept.ID] = deptNested

	// Store business unit if exists
	if rel.BusinessUnit != nil {
		enriched.BusinessUnits[rel.Department.ID] = &BusinessUnitNested{
			ID:        rel.BusinessUnit.ID,
			Name:      rel.BusinessUnit.Name,
			Shortname: rel.BusinessUnit.Shortname,
		}
	}
}

// loadLineManagers fetches line manager information
func (s *StarterEnrichmentService) loadLineManagers(
	ctx context.Context,
	lineManagerIDs map[int64]bool,
	enriched *EnrichedData,
) error {
	// Convert map to slice
	ids := make([]int64, 0, len(lineManagerIDs))
	for id := range lineManagerIDs {
		ids = append(ids, id)
	}

	// Batch load starters (line managers are also starters)
	for _, id := range ids {
		manager, err := s.starterRepo.FindByID(ctx, id)
		if err != nil {
			continue // Skip if not found
		}

		enriched.LineManagers[manager.ID()] = &LineManagerNested{
			ID:       manager.ID(),
			Domain:   manager.Domain(),
			Name:     manager.Name(),
			Email:    manager.Email(),
			JobTitle: manager.JobTitle(),
		}
	}

	return nil
}
