package dto

// ListDepartmentsRequest represents the HTTP request for listing departments
type ListDepartmentsRequest struct {
	BusinessUnitID        *int64 `form:"business_unit_id" binding:"omitempty,gt=0"`
	Page                  int    `form:"page" binding:"omitempty,min=1"`
	Limit                 int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

// SetDefaults sets default values for the request
func (r *ListDepartmentsRequest) SetDefaults() {
	if r.Page <= 0 {
		r.Page = 1
	}
	if r.Limit <= 0 {
		r.Limit = 10
	}
	if r.Limit > 100 {
		r.Limit = 100
	}
}
