package query

import "github.com/kiin21/go-rest/pkg/httputil"

type ListStartersQuery struct {
	Pagination httputil.ReqPagination

	SearchBy string

	Keyword string

	SortBy    string
	SortOrder string
}
