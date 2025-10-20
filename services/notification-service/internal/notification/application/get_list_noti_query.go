package application

import "github.com/kiin21/go-rest/pkg/httputil"

type ListNotificationsQuery struct {
	Pagination httputil.ReqPagination
	SortBy     string
	SortOrder  string
}
