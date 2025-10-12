package dto

import (
	"time"

	"github.com/kiin21/go-rest/internal/organization/domain"
)

// DepartmentDetailResponse represents the full department response with all nested data
type DepartmentDetailResponse struct {
	ID                  int64               `json:"id"`
	GroupDepartmentID   *int64              `json:"group_department_id"`
	FullName            string              `json:"full_name"`
	Shortname           string              `json:"shortname"`
	BusinessUnit        *BusinessUnitNested `json:"business_unit,omitempty"`
	Leader              *LeaderNested       `json:"leader,omitempty"`
	ParentDepartment    *DepartmentNested   `json:"parent_department,omitempty"`
	MembersCount        int                 `json:"members_count"`
	SubdepartmentsCount int                 `json:"subdepartments_count"`
	Subdepartments      []*DepartmentNested `json:"subdepartments,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

// BusinessUnitNested represents nested business unit data
type BusinessUnitNested struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"shortname,omitempty"`
}

// LeaderNested represents nested leader (starter) data
type LeaderNested struct {
	ID       int64  `json:"id"`
	Domain   string `json:"domain"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	JobTitle string `json:"job_title"`
}

// DepartmentNested represents simplified department (for parent/subdepartments)
type DepartmentNested struct {
	ID           int64  `json:"id"`
	FullName     string `json:"full_name"`
	Shortname    string `json:"shortname"`
	MembersCount int    `json:"members_count"`
}

// FromDomainWithDetails converts domain.DepartmentWithDetails to DTO
func FromDomainWithDetails(dept *domain.DepartmentWithDetails) *DepartmentDetailResponse {
	response := &DepartmentDetailResponse{
		ID:                  dept.ID,
		GroupDepartmentID:   dept.GroupDepartmentID,
		FullName:            dept.FullName,
		Shortname:           dept.Shortname,
		MembersCount:        dept.MembersCount,
		SubdepartmentsCount: dept.SubdepartmentsCount,
		CreatedAt:           dept.CreatedAt,
		UpdatedAt:           dept.UpdatedAt,
	}

	// Map business unit
	if dept.BusinessUnit != nil {
		response.BusinessUnit = &BusinessUnitNested{
			ID:        dept.BusinessUnit.ID,
			Name:      dept.BusinessUnit.Name,
			Shortname: dept.BusinessUnit.Shortname,
		}
	}

	// Map leader
	if dept.Leader != nil {
		response.Leader = &LeaderNested{
			ID:       dept.Leader.ID,
			Domain:   dept.Leader.Domain,
			Name:     dept.Leader.Name,
			Email:    dept.Leader.Email,
			JobTitle: dept.Leader.JobTitle,
		}
	}

	// Map parent department
	if dept.ParentDepartment != nil {
		response.ParentDepartment = &DepartmentNested{
			ID:           dept.ParentDepartment.ID,
			FullName:     dept.ParentDepartment.FullName,
			Shortname:    dept.ParentDepartment.Shortname,
			MembersCount: dept.ParentDepartment.MembersCount,
		}
	}

	// Map subdepartments
	if len(dept.Subdepartments) > 0 {
		response.Subdepartments = make([]*DepartmentNested, len(dept.Subdepartments))
		for i, sd := range dept.Subdepartments {
			response.Subdepartments[i] = &DepartmentNested{
				ID:           sd.ID,
				FullName:     sd.FullName,
				Shortname:    sd.Shortname,
				MembersCount: sd.MembersCount,
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
