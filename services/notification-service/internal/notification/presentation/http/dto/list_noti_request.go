package dto

type ListNotiRequest struct {
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=from to type timestamp"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`

	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (r *ListNotiRequest) SetDefaults() {
	if r.Page == 0 {
		r.Page = 1
	}
	if r.Limit == 0 {
		r.Limit = 20
	}
	if r.SortBy == "" {
		r.SortBy = "timestamp"
	}
	if r.SortOrder == "" {
		r.SortOrder = "desc"
	}
}
