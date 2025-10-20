package query

import "github.com/kiin21/go-rest/pkg/httputil"

type ListDepartmentsQuery struct {
	Pagination     httputil.ReqPagination
	BusinessUnitID *int64
}
