package department

// UpdateDepartmentRequest represents the HTTP request for updating a department
// All fields are optional (pointers) to support partial updates
type UpdateDepartmentRequest struct {
	FullName          *string `json:"full_name" binding:"omitempty,min=3,max=255"`
	Shortname         *string `json:"shortname" binding:"omitempty,min=2,max=100"`
	BusinessUnitID    *int64  `json:"business_unit_id" binding:"omitempty,gt=0"`
	GroupDepartmentID *int64  `json:"group_department_id" binding:"omitempty,gt=0"`
	LeaderID          *int64  `json:"leader_id" binding:"omitempty,gt=0"`
}
