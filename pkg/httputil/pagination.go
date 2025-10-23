package httputil

import "github.com/gin-gonic/gin"

type ReqPagination struct {
	Page  *int `json:"page" form:"page"`
	Limit *int `json:"limit" form:"limit"`
}

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

func (p *ReqPagination) GetPage() int {
	if p.Page == nil || *p.Page < 1 {
		return DefaultPage
	}
	return *p.Page
}

func (p *ReqPagination) GetLimit() int {
	if p.Limit == nil || *p.Limit < 1 {
		return DefaultLimit
	}

	limit := *p.Limit
	if limit > MaxLimit {
		return MaxLimit
	}

	return limit
}

func (p *ReqPagination) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
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

func CursorPagination(ctx *gin.Context, pagination RespPagination) RespPagination {
	updated := pagination
	if pagination.Prev != nil {
		prevURL := AbsoluteURLWithPage(ctx, *pagination.Prev)
		updated.Prev = &prevURL
	}
	if pagination.Next != nil {
		nextURL := AbsoluteURLWithPage(ctx, *pagination.Next)
		updated.Next = &nextURL
	}
	return updated
}

func AbsoluteURLWithPage(ctx *gin.Context, page string) string {
	query := ctx.Request.URL.Query()
	query.Set("page", page)

	resolver := NewRequestURLResolver()
	return resolver.AbsoluteURL(ctx, ctx.Request.URL.Path, query)
}
