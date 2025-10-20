package department

import (
	"time"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/internal/organization/presentation/http/dto/shared"
)

type DepartmentDetailResponse struct {
	ID               int64                      `json:"id"`
	FullName         string                     `json:"full_name"`
	Shortname        string                     `json:"shortname"`
	BusinessUnit     *shared.BusinessUnitNested `json:"business_unit,omitempty"`
	Leader           *shared.LineManagerNested  `json:"leader,omitempty"`
	ParentDepartment *shared.DepartmentNested   `json:"parent_department,omitempty"`
	Subdepartments   []*shared.DepartmentNested `json:"sub_departments,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

// DepartmentNested represents a nested department object in a response.

func FromDomainWithDetails(dept *model.DepartmentWithDetails) *DepartmentDetailResponse {
	response := &DepartmentDetailResponse{
		ID:        dept.ID,
		FullName:  dept.FullName,
		Shortname: dept.Shortname,
		CreatedAt: dept.CreatedAt,
		UpdatedAt: dept.UpdatedAt,
	}

	if dept.BusinessUnit != nil {
		response.BusinessUnit = &shared.BusinessUnitNested{
			ID:        dept.BusinessUnit.ID,
			Name:      dept.BusinessUnit.Name,
			Shortname: dept.BusinessUnit.Shortname,
		}
	}

	if dept.Leader != nil {
		response.Leader = &shared.LineManagerNested{
			ID:       dept.Leader.ID,
			Domain:   dept.Leader.Domain,
			Name:     dept.Leader.Name,
			Email:    dept.Leader.Email,
			JobTitle: dept.Leader.JobTitle,
		}
	}

	if dept.ParentDepartment != nil {
		response.ParentDepartment = &shared.DepartmentNested{
			ID:        dept.ParentDepartment.ID,
			Name:      dept.ParentDepartment.FullName,
			Shortname: dept.ParentDepartment.Shortname,
		}
	}

	if len(dept.Subdepartments) > 0 {
		response.Subdepartments = make([]*shared.DepartmentNested, len(dept.Subdepartments))
		for i, sd := range dept.Subdepartments {
			response.Subdepartments[i] = &shared.DepartmentNested{
				ID:        sd.ID,
				Name:      sd.FullName,
				Shortname: sd.Shortname,
			}
		}
	}

	return response
}

// FromDomainsWithDetails converts slice of domain.DepartmentWithDetails to DTOs
func FromDomainsWithDetails(depts []*model.DepartmentWithDetails) []*DepartmentDetailResponse {
	responses := make([]*DepartmentDetailResponse, len(depts))
	for i, dept := range depts {
		responses[i] = FromDomainWithDetails(dept)
	}
	return responses
}
