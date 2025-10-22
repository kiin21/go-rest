package starter

import "github.com/kiin21/go-rest/pkg/httputil"

// StarterListAPIResponse wraps paginated starter data.
type StarterListAPIResponse struct {
	Code    int                        `json:"code"`
	Message string                     `json:"message"`
	Data    StarterListAPIResponseData `json:"data"`
	Error   interface{}                `json:"error,omitempty"`
}

// StarterListAPIResponseData contains starters and pagination metadata.
type StarterListAPIResponseData struct {
	Data       []*StarterResponse      `json:"data"`
	Pagination httputil.RespPagination `json:"pagination"`
}

// StarterAPIResponse wraps a single starter payload.
type StarterAPIResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *StarterResponse `json:"data"`
	Error   interface{}      `json:"error,omitempty"`
}

// StarterDeleteAPIResponse contains the response after deleting a starter.
type StarterDeleteAPIResponse struct {
	Code    int                       `json:"code"`
	Message string                    `json:"message"`
	Data    StarterDeleteResponseData `json:"data"`
	Error   interface{}               `json:"error,omitempty"`
}

// StarterDeleteResponseData holds delete confirmation details.
type StarterDeleteResponseData struct {
	Domain  string `json:"domain"`
	Message string `json:"message"`
}
