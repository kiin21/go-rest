package dto

import (
	"time"

	"github.com/kiin21/go-rest/internal/organization/domain"
)

type DepartmentDetailResponse struct {
	ID                  int64               `json:"id"`
	FullName            string              `json:"full_name"`
	Shortname           string              `json:"shortname"`
	BusinessUnit        *BusinessUnitNested `json:"business_unit,omitempty"`
	Leader              *LeaderNested       `json:"leader,omitempty"`
	ParentDepartment    *DepartmentNested   `json:"parent_department,omitempty"`
	Subdepartments      []*DepartmentNested `json:"subdepartments,omitempty"`
	MembersCount        int                 `json:"members_count"`
	SubdepartmentsCount int                 `json:"subdepartments_count"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

type DepartmentNested struct {
	ID        int64  `json:"id"`
	FullName  string `json:"full_name"`
	Shortname string `json:"shortname"`
}

func FromDomainWithDetails(dept *domain.DepartmentWithDetails) *DepartmentDetailResponse {
	response := &DepartmentDetailResponse{
		ID:                  dept.ID,
		FullName:            dept.FullName,
		Shortname:           dept.Shortname,
		MembersCount:        dept.MembersCount,
		SubdepartmentsCount: dept.SubdepartmentsCount,
		CreatedAt:           dept.CreatedAt,
		UpdatedAt:           dept.UpdatedAt,
	}

	if dept.BusinessUnit != nil {
		response.BusinessUnit = &BusinessUnitNested{
			ID:        dept.BusinessUnit.ID,
			Name:      dept.BusinessUnit.Name,
			Shortname: dept.BusinessUnit.Shortname,
		}
	}

	if dept.Leader != nil {
		response.Leader = &LeaderNested{
			ID:       dept.Leader.ID,
			Domain:   dept.Leader.Domain,
			Name:     dept.Leader.Name,
			Email:    dept.Leader.Email,
			JobTitle: dept.Leader.JobTitle,
		}
	}

	if dept.ParentDepartment != nil {
		response.ParentDepartment = &DepartmentNested{
			ID:        dept.ParentDepartment.ID,
			FullName:  dept.ParentDepartment.FullName,
			Shortname: dept.ParentDepartment.Shortname,
		}
	}

	if len(dept.Subdepartments) > 0 {
		response.Subdepartments = make([]*DepartmentNested, len(dept.Subdepartments))
		for i, sd := range dept.Subdepartments {
			response.Subdepartments[i] = &DepartmentNested{
				ID:        sd.ID,
				FullName:  sd.FullName,
				Shortname: sd.Shortname,
			}
		}
	}

	return response
}

// FromDomainsWithDetails converts slice of domain.DepartmentWithDetails to DTOs
func FromDomainsWithDetails(depts []*domain.DepartmentWithDetails) []*DepartmentDetailResponse {
	responses := make([]*DepartmentDetailResponse, len(depts))
	for i, dept := range depts {
		responses[i] = FromDomainWithDetails(dept)
	}
	return responses
}
