package query

import "github.com/kiin21/go-rest/pkg/httputil"

type SearchStartersQuery struct {
	Keyword string

	DepartmentID   *int64
	BusinessUnitID *int64

	Pagination httputil.ReqPagination

	SortBy    string
	SortOrder string
}
