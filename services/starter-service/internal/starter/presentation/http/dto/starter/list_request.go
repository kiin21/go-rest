package starter

import (
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
)

type ListStartersRequest struct {
	Query    *string `form:"q"`
	SearchBy string  `form:"search_by" binding:"omitempty,oneof=fullname domain dept_name bu_name"`

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

func (r *ListStartersRequest) ToQuery() *query.ListStartersQuery {
	var keyword string
	if r.Query != nil {
		keyword = *r.Query
	}

	return &query.ListStartersQuery{
		Pagination: httputil.ReqPagination{
			Page:  &r.Page,
			Limit: &r.Limit,
		},
		Keyword:   keyword,
		SearchBy:  r.SearchBy,
		SortBy:    r.SortBy,
		SortOrder: r.SortOrder,
	}
}
