package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	startercommand "github.com/kiin21/go-rest/internal/organization/application/dto/starter/command"
	starterquery "github.com/kiin21/go-rest/internal/organization/application/dto/starter/query"
	"github.com/kiin21/go-rest/internal/organization/application/service"
	"github.com/kiin21/go-rest/internal/organization/domain/model"
	starterdto "github.com/kiin21/go-rest/internal/organization/presentation/http/dto/starter"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/pkg/httpctx"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterHandler struct {
	service     *service.StarterApplicationService
	urlResolver httpctx.RequestURLResolver
}

func NewStarterHandler(
	service *service.StarterApplicationService,
	urlResolver httpctx.RequestURLResolver,
) *StarterHandler {
	return &StarterHandler{
		service:     service,
		urlResolver: urlResolver,
	}
}

// ListStarters godoc
// @Summary List starters
// @Description Retrieve starters with optional filters and pagination
// @Tags Starters
// @Produce json
// @Param business_unit_id query int false "Filter by business unit"
// @Param department_id query int false "Filter by department"
// @Param q query string false "Keyword search"
// @Param sort_by query string false "Sort field"
// @Param sort_order query string false "Sort order"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /starters [get]
func (sh *StarterHandler) ListStarters(ctx *gin.Context) {
	response.Wrap(sh.listStarters)(ctx)
}

func (sh *StarterHandler) listStarters(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.ListStartersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	var keyword string
	if req.Query != nil {
		keyword = *req.Query
	}

	query := starterquery.ListStartersQuery{
		Pagination: response.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		Keyword:        keyword,
		BusinessUnitID: req.BusinessUnitID,
		DepartmentID:   req.DepartmentID,
		SortBy:         req.SortBy,
		SortOrder:      req.SortOrder,
	}

	rawResult, err := sh.service.ListStarters(ctx, query)
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to list starters")
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, rawResult.Data)
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to enrich starters")
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	responseData := starterdto.FromStartersEnriched(rawResult.Data, enrichedDTO)

	return &response.PaginatedResult[*starterdto.StarterResponse]{
		Data:       responseData,
		Pagination: decoratePagination(ctx, sh.urlResolver, rawResult.Pagination),
	}, nil
}

// CreateStarter godoc
// @Summary Create a starter
// @Description Create a new starter record
// @Tags Starters
// @Accept json
// @Produce json
// @Param request body starterdto.CreateStarterRequest true "Starter payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters [post]
func (sh *StarterHandler) CreateStarter(ctx *gin.Context) {
	response.Wrap(sh.createStarter)(ctx)
}

func (sh *StarterHandler) createStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.CreateStarterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
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
			return nil, response.NewAPIError(http.StatusConflict, "Domain already exists", err.Error())
		case errors.Is(err, sharedDomain.ErrValidation), errors.Is(err, sharedDomain.ErrInvalidInput):
			return nil, response.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		default:
			return nil, mapServiceError(err, "", "Failed to create starter")
		}
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to enrich starter", err.Error())
	}

	enrichedDTO := starterdto.FromDomainEnrichment(enrichedDomain)
	return starterdto.FromDomainEnriched(starter, enrichedDTO), nil
}

// Find godoc
// @Summary Get starter detail
// @Description Get a starter by domain
// @Tags Starters
// @Produce json
// @Param domain path string true "Starter domain"
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /starters/{domain} [get]
func (sh *StarterHandler) Find(ctx *gin.Context) {
	response.Wrap(sh.find)(ctx)
}

func (sh *StarterHandler) find(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.FindStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	starter, err := sh.service.GetStarterByDomain(ctx, req.Domain)
	if err != nil {
		return nil, mapServiceError(err, "Starter not found", "Failed to fetch starter")
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to enrich starter")
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
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /starters/{domain} [patch]
func (sh *StarterHandler) UpdateStarter(ctx *gin.Context) {
	response.Wrap(sh.updateStarter)(ctx)
}

func (sh *StarterHandler) updateStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.UpdateStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
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
			return nil, response.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		case errors.Is(err, sharedDomain.ErrNotFound):
			return nil, response.NewAPIError(http.StatusNotFound, "Starter not found", err.Error())
		default:
			return nil, mapServiceError(err, "", "Failed to update starter")
		}
	}

	enrichedDomain, err := sh.service.EnrichStarters(ctx, []*model.Starter{starter})
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to enrich starter")
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
// @Success 200 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /starters/{domain} [delete]
func (sh *StarterHandler) SoftDeleteStarter(ctx *gin.Context) {
	response.Wrap(sh.softDeleteStarter)(ctx)
}

func (sh *StarterHandler) softDeleteStarter(ctx *gin.Context) (res interface{}, err error) {
	var req starterdto.DeleteStarterRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	if err := sh.service.SoftDeleteStarter(ctx, req.Domain); err != nil {
		return nil, mapServiceError(err, "Starter not found", "Failed to delete starter")
	}

	return gin.H{
		"message": "Starter soft deleted successfully",
		"domain":  req.Domain,
	}, nil
}
