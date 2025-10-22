package businessunit

import "github.com/kiin21/go-rest/pkg/httputil"

// BusinessUnitListAPIResponse wraps paginated business unit data.
type BusinessUnitListAPIResponse struct {
	Code    int                             `json:"code"`
	Message string                          `json:"message"`
	Data    BusinessUnitListAPIResponseData `json:"data"`
	Error   interface{}                     `json:"error,omitempty"`
}

// BusinessUnitListAPIResponseData contains business units and pagination metadata.
type BusinessUnitListAPIResponseData struct {
	Data       []*BusinessUnitDetailResponse `json:"data"`
	Pagination httputil.RespPagination       `json:"pagination"`
}

// BusinessUnitDetailAPIResponse wraps a single business unit payload.
type BusinessUnitDetailAPIResponse struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    *BusinessUnitDetailResponse `json:"data"`
	Error   interface{}                 `json:"error,omitempty"`
}
