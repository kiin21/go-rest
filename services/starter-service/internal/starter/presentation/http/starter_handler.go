package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	startercommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/command"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	shareddto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/shared"
	starterdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/starter"
)

var (
	_ shareddto.GenericAPIResponse
	_ starterdto.StarterListAPIResponse
	_ starterdto.StarterAPIResponse
	_ starterdto.StarterDeleteAPIResponse
)

type StarterHandler struct {
	service     *service.StarterApplicationService
	urlResolver httputil.RequestURLResolver
}

func NewStarterHandler(
	service *service.StarterApplicationService,
	urlResolver httputil.RequestURLResolver,
) *StarterHandler {
	return &StarterHandler{
		service:     service,
		urlResolver: urlResolver,
	}
}

// ListStarters godoc
// @Summary List starters
// @Description Retrieve starters with optional search and pagination
// @Tags Starters
// @Produce json
// @Param q query string false "Keyword search"
// @Param search_by query string false "Search field" Enums(fullname,domain,dept_name,bu_name)
// @Param sort_by query string false "Sort field" Enums(id,domain,created_at) default(id)
// @Param sort_order query string false "Sort order" Enums(asc,desc) default(asc)
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(20)
// @Success 200 {object} starterdto.StarterListAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /starters [get]
func (sh *StarterHandler) ListStarters(ctx *gin.Context) {
	httputil.Wrap(sh.listStarters)(ctx)
}

func (sh *StarterHandler) listStarters(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.ListStartersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	var keyword string
	if req.Query != nil {
		keyword = *req.Query
	}

	query := starterquery.ListStartersQuery{
		Pagination: httputil.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Keyword:   keyword,
		SearchBy:  req.SearchBy,
		SortBy:    req.SortBy,
		SortOrder: req.SortOrder,
	}

	rawResult, err := sh.service.ListStarters(ctx, query)
	if err != nil {
		return nil, err
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, rawResult.Data)
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	responseData := starterdto.FromStartersEnriched(rawResult.Data, enrichedDTO)

	return &httputil.PaginatedResult[*starterdto.StarterResponse]{
		Data:       responseData,
		Pagination: httputil.DecoratePagination(ctx, sh.urlResolver, rawResult.Pagination),
	}, nil
}

// CreateStarter godoc
// @Summary Create starter
// @Description Create a new starter record
// @Tags Starters
// @Accept json
// @Produce json
// @Param request body starterdto.CreateStarterRequest true "Starter payload"
// @Success 200 {object} starterdto.StarterAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 409 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /starters [post]
func (sh *StarterHandler) CreateStarter(ctx *gin.Context) {
	httputil.Wrap(sh.createStarter)(ctx)
}

func (sh *StarterHandler) createStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.CreateStarterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	command := startercommand.CreateStarterCommand{
		Domain:        req.Domain,
		Name:          req.Name,
		Email:         req.Email,
		Mobile:        req.Mobile,
		WorkPhone:     req.WorkPhone,
		JobTitle:      req.JobTitle,
		DepartmentID:  req.DepartmentID,
		LineManagerID: req.LineManagerID,
	}

	starter, err := sh.service.CreateStarter(ctx, command)
	if err != nil {
		switch {
		case errors.Is(err, sharedDomain.ErrDomainAlreadyExists), errors.Is(err, sharedDomain.ErrDuplicateEntry):
			return nil, httputil.NewAPIError(http.StatusConflict, "Domain already exists", err.Error())
		case errors.Is(err, sharedDomain.ErrValidation), errors.Is(err, sharedDomain.ErrInvalidInput):
			return nil, httputil.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		default:
			return nil, err
		}
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, httputil.NewAPIError(http.StatusInternalServerError, "Failed to enrich starter", err.Error())
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// Find godoc
// @Summary Get starter detail
// @Description Retrieve a starter by domain
// @Tags Starters
// @Produce json
// @Param domain path string true "Starter domain"
// @Success 200 {object} starterdto.StarterAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /starters/{domain} [get]
func (sh *StarterHandler) Find(ctx *gin.Context) {
	httputil.Wrap(sh.find)(ctx)
}

func (sh *StarterHandler) find(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.FindStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	starter, err := sh.service.GetStarterByDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// UpdateStarter godoc
// @Summary Update starter
// @Description Partially update a starter by domain
// @Tags Starters
// @Accept json
// @Produce json
// @Param domain path string true "Starter domain"
// @Param request body starterdto.UpdateStarterRequest true "Update payload"
// @Success 200 {object} starterdto.StarterAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /starters/{domain} [patch]
func (sh *StarterHandler) UpdateStarter(ctx *gin.Context) {
	httputil.Wrap(sh.updateStarter)(ctx)
}

func (sh *StarterHandler) updateStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.UpdateStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	command := startercommand.UpdateStarterCommand{
		Domain:        req.Domain,
		Name:          req.Name,
		Email:         req.Email,
		Mobile:        req.Mobile,
		WorkPhone:     req.WorkPhone,
		JobTitle:      req.JobTitle,
		DepartmentID:  req.DepartmentID,
		LineManagerID: req.LineManagerID,
	}

	starter, err := sh.service.UpdateStarter(ctx, command)
	if err != nil {
		switch {
		case errors.Is(err, sharedDomain.ErrValidation), errors.Is(err, sharedDomain.ErrInvalidInput):
			return nil, httputil.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		case errors.Is(err, sharedDomain.ErrNotFound):
			return nil, httputil.NewAPIError(http.StatusNotFound, "Starter not found", err.Error())
		default:
			return nil, err
		}
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, err
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// SoftDeleteStarter godoc
// @Summary Delete starter
// @Description Soft delete a starter by domain
// @Tags Starters
// @Produce json
// @Param domain path string true "Starter domain"
// @Success 200 {object} starterdto.StarterDeleteAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /starters/{domain} [delete]
func (sh *StarterHandler) SoftDeleteStarter(ctx *gin.Context) {
	httputil.Wrap(sh.softDeleteStarter)(ctx)
}

func (sh *StarterHandler) softDeleteStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.DeleteStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	if err := sh.service.SoftDeleteStarter(ctx, req.Domain); err != nil {
		return nil, err
	}

	return gin.H{
		"message": "Starter soft deleted successfully",
		"domain":  req.Domain,
	}, nil
}
