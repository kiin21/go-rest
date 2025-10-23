package department

import (
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/query"
)

type ListDepartmentsRequest struct {
	BusinessUnitID *int64 `form:"business_unit_id" binding:"omitempty,gt=0"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	Limit          int    `form:"limit" binding:"omitempty,min=1,max=100"`
}

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

func (r *ListDepartmentsRequest) ToQuery() *query.ListDepartmentsQuery {
	return &query.ListDepartmentsQuery{
		BusinessUnitID: r.BusinessUnitID,
		Pagination: httputil.ReqPagination{
			Page:  &r.Page,
			Limit: &r.Limit,
		},
	}
}
