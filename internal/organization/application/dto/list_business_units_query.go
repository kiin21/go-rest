package dto

import "github.com/kiin21/go-rest/pkg/response"

// ListBusinessUnitsQuery represents the application-level query for listing business units.
type ListBusinessUnitsQuery struct {
	Pagination response.ReqPagination
}
