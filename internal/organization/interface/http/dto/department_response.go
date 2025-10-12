package dto

import (
	"time"

	"github.com/kiin21/go-rest/internal/organization/domain"
)

// DepartmentResponse represents a department in API responses
type DepartmentResponse struct {
	ID                int64      `json:"id"`
	GroupDepartmentID *int64     `json:"group_department_id,omitempty"`
	FullName          string     `json:"full_name"`
	Shortname         string     `json:"shortname"`
	BusinessUnitID    *int64     `json:"business_unit_id,omitempty"`
	LeaderID          *int64     `json:"leader_id,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"`
}

// FromDepartment converts domain entity to response DTO
func FromDepartment(dept *domain.Department) *DepartmentResponse {
	if dept == nil {
		return nil
	}

	return &DepartmentResponse{
		ID:                dept.ID,
		GroupDepartmentID: dept.GroupDepartmentID,
		FullName:          dept.FullName,
		Shortname:         dept.Shortname,
		BusinessUnitID:    dept.BusinessUnitID,
		LeaderID:          dept.LeaderID,
		CreatedAt:         dept.CreatedAt,
		UpdatedAt:         dept.UpdatedAt,
		DeletedAt:         dept.DeletedAt,
	}
}

// FromDepartments converts multiple domain entities to response DTOs
func FromDepartments(depts []*domain.Department) []*DepartmentResponse {
	responses := make([]*DepartmentResponse, len(depts))
	for i, dept := range depts {
		responses[i] = FromDepartment(dept)
	}
	return responses
}
