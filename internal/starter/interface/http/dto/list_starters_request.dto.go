package dto

type ListStartersRequest struct {
	BusinessUnitID *int64  `form:"business_unit_id" binding:"omitempty,gt=0"`
	DepartmentID   *int64  `form:"department_id" binding:"omitempty,gt=0"`
	Query          *string `form:"q"`

	SortBy    string `form:"sort_by" binding:"omitempty,oneof=id domain created_at"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`

	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (r *ListStartersRequest) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 20
	}
	if r.SortBy == "" {
		r.SortBy = "id"
	}
	if r.SortOrder == "" {
		r.SortOrder = "asc"
	}
}
