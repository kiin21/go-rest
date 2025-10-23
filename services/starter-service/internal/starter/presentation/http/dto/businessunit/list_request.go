package businessunit

import (
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/businessunit/query"
)

type ListBusinessUnitsRequest struct {
	Page  int `form:"page" binding:"omitempty,min=1"`
	Limit int `form:"limit" binding:"omitempty,min=1,max=100"`
}

func (r *ListBusinessUnitsRequest) SetDefaults() {
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

func (r *ListBusinessUnitsRequest) ToQuery() *query.ListBusinessUnitsQuery {
	return &query.ListBusinessUnitsQuery{
		Pagination: httputil.ReqPagination{
			Page:  &r.Page,
			Limit: &r.Limit,
		},
	}
}
