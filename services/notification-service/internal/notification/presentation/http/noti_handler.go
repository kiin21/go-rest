package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	notiapp "github.com/kiin21/go-rest/services/notification-service/internal/notification/application"
	"github.com/kiin21/go-rest/services/notification-service/internal/notification/presentation/http/dto"
)

type NotiHandler struct {
	service     *notiapp.NotiApplicationService
	urlResolver httputil.RequestURLResolver
}

func NewNotiHandler(
	service *notiapp.NotiApplicationService,
	urlResolver httputil.RequestURLResolver,
) *NotiHandler {
	return &NotiHandler{
		service:     service,
		urlResolver: urlResolver,
	}
}

// GetList godoc
// @Summary List notifications
// @Description Retrieve notifications with pagination and sorting options
// @Tags Notifications
// @Produce json
// @Param sort_by query string false "Sort field" Enums(from,to,type,timestamp) default(timestamp)
// @Param sort_order query string false "Sort order" Enums(asc,desc) default(desc)
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(20)
// @Success 200 {object} dto.ListNotiListAPIResponse
// @Failure 400 {object} dto.GenericAPIResponse
// @Failure 500 {object} dto.GenericAPIResponse
// @Router /notifications [get]
func (h *NotiHandler) GetList(ctx *gin.Context) {
	httputil.Wrap(h.getList)(ctx)
}

func (h *NotiHandler) getList(ctx *gin.Context) (res interface{}, err error) {
	var req dto.ListNotiRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := notiapp.ListNotificationsQuery{
		Pagination: httputil.ReqPagination{
			Page:  &req.Page,
			Limit: &req.Limit,
		},
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	rawResult, err := h.service.ListNotifications(ctx.Request.Context(), query)
	if err != nil {
		return nil, httputil.NewAPIError(http.StatusInternalServerError, "Failed to list notifications", err.Error())
	}

	responseData := dto.FromDomain(rawResult.Data)

	return &httputil.PaginatedResult[*dto.ListNotiResponse]{
		Data:       responseData,
		Pagination: httputil.CursorPagination(ctx, rawResult.Pagination),
	}, nil
}
