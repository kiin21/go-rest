package query

import "github.com/kiin21/go-rest/pkg/response"

type ListStartersQuery struct {
	Pagination response.ReqPagination

	BusinessUnitID *int64
	DepartmentID   *int64

	Keyword string

	SortBy    string
	SortOrder string
}
