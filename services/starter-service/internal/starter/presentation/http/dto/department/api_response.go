package department

import (
	"github.com/kiin21/go-rest/pkg/httputil"
)

// DepartmentListAPIResponse wraps paginated department data.
type DepartmentListAPIResponse struct {
	Code    int                           `json:"code"`
	Message string                        `json:"message"`
	Data    DepartmentListAPIResponseData `json:"data"`
	Error   interface{}                   `json:"error,omitempty"`
}

// DepartmentListAPIResponseData contains departments and pagination metadata.
type DepartmentListAPIResponseData struct {
	Data       []*DepartmentDetailResponse `json:"data"`
	Pagination httputil.RespPagination     `json:"pagination"`
}

// DepartmentDetailAPIResponse wraps a single department payload.
type DepartmentDetailAPIResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    *DepartmentDetailResponse `json:"data"`
	Error   interface{}               `json:"error,omitempty"`
}

// DepartmentDeleteAPIResponse contains the response after deleting a department.
type DepartmentDeleteAPIResponse struct {
	Code    int                          `json:"code"`
	Message string                       `json:"message"`
	Data    DepartmentDeleteResponseData `json:"data"`
	Error   interface{}                  `json:"error,omitempty"`
}

// DepartmentDeleteResponseData holds delete confirmation details.
type DepartmentDeleteResponseData struct {
	ID      int64  `json:"id"`
	Message string `json:"message"`
}
