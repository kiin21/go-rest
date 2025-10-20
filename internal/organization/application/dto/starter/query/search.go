package query

import "github.com/kiin21/go-rest/pkg/response"

type SearchStartersQuery struct {
	Keyword string

	DepartmentID   *int64
	BusinessUnitID *int64

	Pagination response.ReqPagination

	SortBy    string
	SortOrder string
}
