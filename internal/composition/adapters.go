package composition

import (
	"context"

	orgApplication "github.com/kiin21/go-rest/internal/organization/application"
	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	starterModel "github.com/kiin21/go-rest/internal/starter/domain/model"
	starterPort "github.com/kiin21/go-rest/internal/starter/domain/port"
)

type starterLeaderLookupAdapter struct {
	repo starterPort.StarterRepository
}

// NewStarterLeaderLookup adapts starter repository to organization leader lookup port.
func NewStarterLeaderLookup(repo starterPort.StarterRepository) orgApplication.LeaderLookup {
	if repo == nil {
		return nil
	}
	return &starterLeaderLookupAdapter{repo: repo}
}

func (s *starterLeaderLookupAdapter) FindStarterIDByDomain(ctx context.Context, domain string) (int64, error) {
	starter, err := s.repo.FindByDomain(ctx, domain)
	if err != nil {
		return 0, err
	}
	if starter == nil {
		return 0, sharedDomain.ErrNotFound
	}
	return starter.ID(), nil
}

func (s *starterLeaderLookupAdapter) FindStarterById(ctx context.Context, id int64) (*aggregate.Starter, error) {
	starter, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if starter == nil {
		return nil, sharedDomain.ErrNotFound
	}
	return starter, nil
}

type organizationLookupAdapter struct {
	departmentRepo orgDomain.DepartmentRepository
}

func NewOrganizationLookup(departmentRepo orgDomain.DepartmentRepository) starterPort.OrganizationLookup {
	if departmentRepo == nil {
		return nil
	}
	return &organizationLookupAdapter{
		departmentRepo: departmentRepo,
	}
}

func (a *organizationLookupAdapter) FindDepartmentRelations(ctx context.Context, ids []int64) ([]*starterModel.DepartmentRelation, error) {
	if a.departmentRepo == nil {
		return nil, nil
	}

	relations, err := a.departmentRepo.FindByIDsWithRelations(ctx, ids)
	if err != nil {
		return nil, err
	}

	results := make([]*starterModel.DepartmentRelation, 0, len(relations))
	for _, rel := range relations {
		if rel == nil || rel.Department == nil {
			continue
		}

		dept := rel.Department
		deptNested := &sharedDomain.DepartmentNested{
			ID:        dept.ID,
			Name:      dept.FullName,
			Shortname: dept.Shortname,
		}

		if rel.ParentDepartment != nil {
			deptNested.GroupDepartment = &sharedDomain.GroupDepartmentNested{
				ID:        rel.ParentDepartment.ID,
				Name:      rel.ParentDepartment.FullName,
				Shortname: rel.ParentDepartment.Shortname,
			}
		}

		var businessUnitNested *sharedDomain.BusinessUnitNested
		if rel.BusinessUnit != nil {
			businessUnitNested = &sharedDomain.BusinessUnitNested{
				ID:        rel.BusinessUnit.ID,
				Name:      rel.BusinessUnit.Name,
				Shortname: rel.BusinessUnit.Shortname,
			}
		}

		results = append(results, &starterModel.DepartmentRelation{
			Department:   deptNested,
			BusinessUnit: businessUnitNested,
		})
	}

	return results, nil
}
