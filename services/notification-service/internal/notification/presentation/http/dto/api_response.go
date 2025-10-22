package dto

import "github.com/kiin21/go-rest/pkg/httputil"

// GenericAPIResponse describes the common envelope returned by the API.
type GenericAPIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

// ListNotiListAPIResponse wraps paginated notification data.
type ListNotiListAPIResponse struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    ListNotiListAPIResponseData `json:"data"`
	Error   interface{}                 `json:"error,omitempty"`
}

// ListNotiListAPIResponseData contains notifications and pagination metadata.
type ListNotiListAPIResponseData struct {
	Data       []*ListNotiResponse     `json:"data"`
	Pagination httputil.RespPagination `json:"pagination"`
}
