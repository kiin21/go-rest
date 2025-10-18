package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/httpctx"
	"github.com/kiin21/go-rest/pkg/response"
)

func decoratePagination(ctx *gin.Context, resolver httpctx.RequestURLResolver, pagination response.RespPagination) response.RespPagination {
	updated := pagination
	if pagination.Prev != nil {
		prevURL := absoluteURLWithPage(ctx, resolver, *pagination.Prev)
		updated.Prev = &prevURL
	}
	if pagination.Next != nil {
		nextURL := absoluteURLWithPage(ctx, resolver, *pagination.Next)
		updated.Next = &nextURL
	}
	return updated
}

func absoluteURLWithPage(ctx *gin.Context, resolver httpctx.RequestURLResolver, page string) string {
	query := ctx.Request.URL.Query()
	query.Set("page", page)
	return resolver.AbsoluteURL(ctx, ctx.Request.URL.Path, query)
}

func mapServiceError(err error, notFoundMsg, genericMsg string) *response.APIError {
	if err == nil {
		return nil
	}

	var apiErr *response.APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}

	if notFoundMsg != "" && errors.Is(err, sharedDomain.ErrNotFound) {
		return response.NewAPIError(http.StatusNotFound, notFoundMsg, err.Error())
	}

	message := genericMsg
	if message == "" {
		message = "Internal server error"
	}
	return response.NewAPIError(http.StatusInternalServerError, message, err.Error())
}
