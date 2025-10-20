package query

import "github.com/kiin21/go-rest/pkg/httputil"

type ListStartersQuery struct {
	Pagination httputil.ReqPagination

	BusinessUnitID *int64
	DepartmentID   *int64

	Keyword string

	SortBy    string
	SortOrder string
}
