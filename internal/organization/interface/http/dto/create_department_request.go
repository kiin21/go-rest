package dto

// CreateDepartmentRequest represents the HTTP request for creating a new department
type CreateDepartmentRequest struct {
	FullName          string `json:"full_name" binding:"required,min=3,max=255"`
	Shortname         string `json:"shortname" binding:"required,min=2,max=100"`
	BusinessUnitID    *int64 `json:"business_unit_id" binding:"omitempty,gt=0"`
	GroupDepartmentID *int64 `json:"group_department_id" binding:"omitempty,gt=0"`
	LeaderID          *int64 `json:"leader_id" binding:"omitempty,gt=0"`
}
