package httputil

type ReqPagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type RespPagination struct {
	Limit      int     `json:"limit"`
	TotalItems int64   `json:"total_items,omitempty"`
	Prev       *string `json:"prev"`
	Next       *string `json:"next"`
}

type PaginatedResult[T any] struct {
	Data       []T            `json:"data"`
	Pagination RespPagination `json:"pagination"`
}
