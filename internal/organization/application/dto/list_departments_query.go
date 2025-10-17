package dto

import "github.com/kiin21/go-rest/pkg/response"

type ListDepartmentsQuery struct {
	Pagination     response.ReqPagination
	BusinessUnitID *int64
}
