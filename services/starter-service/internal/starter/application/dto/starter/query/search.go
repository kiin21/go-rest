package query

import "github.com/kiin21/go-rest/pkg/httputil"

type SearchStartersQuery struct {
	Keyword string

	SearchBy string

	Pagination httputil.ReqPagination

	SortBy    string
	SortOrder string
}
