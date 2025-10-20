package service

import (
	"context"
	"fmt"

	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	repo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

type StarterEnrichmentService struct {
	starterRepo      repo.StarterRepository
	departmentRepo   repo.DepartmentRepository
	businessUnitRepo repo.BusinessUnitRepository
}

func NewStarterEnrichmentService(
	starterRepo repo.StarterRepository,
	departmentRepo repo.DepartmentRepository,
	businessUnitRepo repo.BusinessUnitRepository,
) *StarterEnrichmentService {
	return &StarterEnrichmentService{
		starterRepo:      starterRepo,
		departmentRepo:   departmentRepo,
		businessUnitRepo: businessUnitRepo,
	}
}

func (s *StarterEnrichmentService) EnrichStarters(ctx context.Context, starters []*model.Starter) (*model.EnrichedData, error) {
	enriched := &model.EnrichedData{
		Departments:   make(map[int64]*model.DepartmentNested),
		LineManagers:  make(map[int64]*model.LineManagerNested),
		BusinessUnits: make(map[int64]*model.BusinessUnitNested),
	}

	// Collect unique IDs.
	departmentIDs := s.collectDepartmentIDs(starters)
	lineManagerIDs := s.collectLineManagerIDs(starters)

	if len(departmentIDs) > 0 {
		if err := s.loadDepartments(ctx, departmentIDs, enriched); err != nil {
			return nil, err
		}
	}

	// Load line managers
	if len(lineManagerIDs) > 0 {
		if err := s.loadLineManagers(ctx, lineManagerIDs, enriched); err != nil {
			return nil, err
		}
	}

	return enriched, nil
}

func (s *StarterEnrichmentService) collectDepartmentIDs(starters []*model.Starter) map[int64]bool {
	departmentIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.DepartmentID() != nil {
			departmentIDs[*starter.DepartmentID()] = true
		}
	}
	return departmentIDs
}

func (s *StarterEnrichmentService) collectLineManagerIDs(starters []*model.Starter) map[int64]bool {
	lineManagerIDs := make(map[int64]bool)
	for _, starter := range starters {
		if starter.LineManagerID() != nil {
			lineManagerIDs[*starter.LineManagerID()] = true
		}
	}
	return lineManagerIDs
}

func (s *StarterEnrichmentService) loadDepartments(
	ctx context.Context,
	departmentIDs map[int64]bool,
	enriched *model.EnrichedData,
) error {
	ids := make([]int64, 0, len(departmentIDs))
	for id := range departmentIDs {
		ids = append(ids, id)
	}

	relations, err := s.departmentRepo.FindByIDsWithDetails(ctx, ids)
	if err != nil {
		return fmt.Errorf("failed to load departments with relations: %w", err)
	}

	for _, rel := range relations {
		s.mapDepartmentRelation(rel, enriched)
	}

	return nil
}

func (s *StarterEnrichmentService) mapDepartmentRelation(
	rel *model.DepartmentWithDetails,
	enriched *model.EnrichedData,
) {
	if rel == nil || rel.Department == nil {
		return
	}

	dept := rel.Department
	deptNested := &model.DepartmentNested{
		ID:        dept.ID,
		Name:      dept.FullName,
		Shortname: dept.Shortname,
	}

	if rel.ParentDepartment != nil {
		deptNested.GroupDepartment = &model.GroupDepartmentNested{
			ID:        rel.ParentDepartment.ID,
			Name:      rel.ParentDepartment.FullName,
			Shortname: rel.ParentDepartment.Shortname,
		}
	}

	enriched.Departments[dept.ID] = deptNested

	if rel.BusinessUnit != nil {
		enriched.BusinessUnits[dept.ID] = &model.BusinessUnitNested{
			ID:        rel.BusinessUnit.ID,
			Name:      rel.BusinessUnit.Name,
			Shortname: rel.BusinessUnit.Shortname,
		}
	}
}

func (s *StarterEnrichmentService) loadLineManagers(
	ctx context.Context,
	lineManagerIDs map[int64]bool,
	enriched *model.EnrichedData,
) error {
	ids := make([]int64, 0, len(lineManagerIDs))
	for id := range lineManagerIDs {
		ids = append(ids, id)
	}

	for _, id := range ids {
		manager, err := s.starterRepo.FindByID(ctx, id)
		if err != nil {
			continue // Skip if not found.
		}

		enriched.LineManagers[manager.ID()] = &model.LineManagerNested{
			ID:       manager.ID(),
			Domain:   manager.Domain(),
			Name:     manager.Name(),
			Email:    manager.Email(),
			JobTitle: manager.JobTitle(),
		}
	}

	return nil
}
