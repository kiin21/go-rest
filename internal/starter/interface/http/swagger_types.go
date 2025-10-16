package http

// StarterCreateRequest documents the starter creation payload. Keep in sync with dto.CreateStarterRequest.
type StarterCreateRequest struct {
	Domain        string `json:"domain"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	WorkPhone     string `json:"work_phone"`
	JobTitle      string `json:"job_title"`
	DepartmentID  *int64 `json:"department_id"`
	LineManagerID *int64 `json:"line_manager_id"`
}

// StarterUpdateRequest documents the starter update payload. Keep in sync with dto.UpdateStarterRequest.
type StarterUpdateRequest struct {
	Name          *string `json:"name"`
	Email         *string `json:"email"`
	Mobile        *string `json:"mobile"`
	WorkPhone     *string `json:"work_phone"`
	JobTitle      *string `json:"job_title"`
	DepartmentID  *int64  `json:"department_id"`
	LineManagerID *int64  `json:"line_manager_id"`
}
