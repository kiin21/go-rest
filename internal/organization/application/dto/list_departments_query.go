package dto

import "github.com/kiin21/go-rest/pkg/response"

// ListDepartmentsQuery represents the application-level query for listing departments
type ListDepartmentsQuery struct {
	Pagination            response.ReqPagination
	BusinessUnitID        *int64
}
