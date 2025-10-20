package dto

type ListNotiPagination struct {
	Limit      int     `json:"limit"`
	TotalItems int64   `json:"total_items,omitempty"`
	Prev       *string `json:"prev,omitempty"`
	Next       *string `json:"next,omitempty"`
}

type ListNotiListAPIResponseData struct {
	Data       []*ListNotiResponse `json:"data"`
	Pagination ListNotiPagination  `json:"pagination"`
}

type ListNotiListAPIResponse struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Data    ListNotiListAPIResponseData `json:"data"`
	Error   interface{}                 `json:"error,omitempty"`
}

type GenericAPIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}
