package dto

import (
	"time"
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
