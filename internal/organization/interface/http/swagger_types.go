package http

// DepartmentCreateRequest documents the department creation payload. Keep in sync with dto.CreateDepartmentRequest.
type DepartmentCreateRequest struct {
	FullName          string `json:"full_name"`
	Shortname         string `json:"shortname"`
	BusinessUnitID    *int64 `json:"business_unit_id"`
	GroupDepartmentID *int64 `json:"group_department_id"`
	LeaderID          *int64 `json:"leader_id"`
}

// DepartmentUpdateRequest documents the department update payload. Keep in sync with dto.UpdateDepartmentRequest.
type DepartmentUpdateRequest struct {
	FullName          *string `json:"full_name"`
	Shortname         *string `json:"shortname"`
	BusinessUnitID    *int64  `json:"business_unit_id"`
	GroupDepartmentID *int64  `json:"group_department_id"`
	LeaderID          *int64  `json:"leader_id"`
}
