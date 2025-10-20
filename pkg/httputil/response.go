package httputil

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func SuccessResponse(c *gin.Context, data interface{}) {
	c.PureJSON(200, APIResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	c.PureJSON(code, APIResponse{
		Code:    code,
		Message: message,
		Error:   err,
	})
}

type HandlerFunc func(ctx *gin.Context) (res interface{}, err error)

func Wrap(handler HandlerFunc) func(c *gin.Context) {
	return func(ctx *gin.Context) {
		res, err := handler(ctx)
		if err != nil {
			var apiErr *APIError
			if errors.As(err, &apiErr) {
				ErrorResponse(ctx, apiErr.StatusCode, apiErr.Message, apiErr.Err)
			}
			return
		}
		SuccessResponse(ctx, res)
	}
}

func DecoratePagination(ctx *gin.Context, resolver RequestURLResolver, pagination RespPagination) RespPagination {
	updated := pagination
	if pagination.Prev != nil {
		prevURL := AbsoluteURLWithPage(ctx, resolver, *pagination.Prev)
		updated.Prev = &prevURL
	}
	if pagination.Next != nil {
		nextURL := AbsoluteURLWithPage(ctx, resolver, *pagination.Next)
		updated.Next = &nextURL
	}
	return updated
}

func AbsoluteURLWithPage(ctx *gin.Context, resolver RequestURLResolver, page string) string {
	query := ctx.Request.URL.Query()
	query.Set("page", page)
	return resolver.AbsoluteURL(ctx, ctx.Request.URL.Path, query)
}
