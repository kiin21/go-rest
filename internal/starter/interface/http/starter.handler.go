package http

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
	"github.com/kiin21/go-rest/internal/starter/application"
	applicationDto "github.com/kiin21/go-rest/internal/starter/application/dto"
	"github.com/kiin21/go-rest/internal/starter/domain"
	presentationDto "github.com/kiin21/go-rest/internal/starter/interface/http/dto"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterHandler struct {
	service       *application.StarterApplicationService
	searchService *application.StarterSearchService
}

func NewStarterHandler(
	service *application.StarterApplicationService,
	searchService *application.StarterSearchService,
) *StarterHandler {
	return &StarterHandler{
		service:       service,
		searchService: searchService,
	}
}

// [GET]: /api/v1/starters?q=duydh2&business_unit_id=1&department_id=1&sort_by=domain&sort_order=asc&page=1&limit=3
// ListStarters godoc
// @Summary List starters
// @Description Returns a paginated list of starters with optional filters and sorting.
// @Tags Starters
// @Accept json
// @Produce json
// @Param q query string false "Keyword to search by domain, name, or email"
// @Param business_unit_id query int false "Filter by business unit ID"
// @Param department_id query int false "Filter by department ID"
// @Param sort_by query string false "Sort field" Enums(id,name,domain,business_unit_id,department_id,created_at,updated_at)
// @Param sort_order query string false "Sort order" Enums(asc,desc)
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters [get]
func (sh *StarterHandler) ListStarters(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.ListStartersRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	var keyword string
	if req.Query != nil {
		keyword = *req.Query
	}

	query := applicationDto.ListStartersQuery{
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

	result, err := sh.service.GetAllStarters(ctx, query)
	if err != nil {
		return nil, err
	}

	// Use 2 service calls to avoid N+1 queries
	// seperate concern in service, make function `GetAllStarters` reusable
	enrichedData, err := sh.service.EnrichStarters(ctx, result.Data)
	if err != nil {
		return nil, err
	}

	// Convert domain entities to enriched response DTOs
	responseData := presentationDto.FromStartersEnriched(result.Data, enrichedData)

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	q := ctx.Request.URL.Query()
	if result.Pagination.Prev != nil {
		q.Set("page", fmt.Sprint(req.Page-1))

		u := &url.URL{
			Scheme:   scheme,
			Host:     ctx.Request.Host,
			Path:     ctx.Request.URL.Path,
			RawQuery: q.Encode(),
		}
		prevURL := u.String()
		result.Pagination.Prev = &prevURL
	}
	if result.Pagination.Next != nil {
		q.Set("page", fmt.Sprint(req.Page+1))

		u := &url.URL{
			Scheme:   scheme,
			Host:     ctx.Request.Host,
			Path:     ctx.Request.URL.Path,
			RawQuery: q.Encode(),
		}
		nextURL := u.String()
		result.Pagination.Next = &nextURL
	}

	return &response.PaginatedResult[*presentationDto.StarterResponse]{
		Data:       responseData,
		Pagination: result.Pagination,
	}, nil
}

// CreateStarter [POST]: /api/v1/starters
// CreateStarter godoc
// @Summary Create a starter
// @Description Creates a new starter record.
// @Tags Starters
// @Accept json
// @Produce json
// @Param request body StarterCreateRequest true "Starter payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters [post]
func (sh *StarterHandler) CreateStarter(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.CreateStarterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	command := applicationDto.CreateStarterCommand{
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
		// Check if it's a domain already exists error (from Domain Service)
		if errors.Is(err, sharedDomain.ErrDomainAlreadyExists) {
			return nil, response.NewAPIError(http.StatusConflict, "Domain already exists", err.Error())
		}
		// Check if it's a duplicate entry error (fallback from Repository)
		if errors.Is(err, sharedDomain.ErrDuplicateEntry) {
			return nil, response.NewAPIError(http.StatusConflict, "Domain already exists", err.Error())
		}
		// Check for validation errors
		if errors.Is(err, sharedDomain.ErrValidation) || errors.Is(err, sharedDomain.ErrInvalidInput) {
			return nil, response.NewAPIError(http.StatusBadRequest, "Validation failed", err.Error())
		}
		// Return 500 for other unexpected errors
		return nil, response.NewAPIError(http.StatusInternalServerError, "Failed to create starter", err.Error())
	}

	// Enrich with related data
	enrichedData, err := sh.service.EnrichStarters(ctx, []*domain.Starter{starter})
	if err != nil {
		// Log error but continue with basic data
	}

	return presentationDto.FromDomainEnriched(starter, enrichedData), nil
}

// [GET] /api/v1/starters/:domain
// Find godoc
// @Summary Get starter detail
// @Description Retrieves starter information by domain.
// @Tags Starters
// @Accept json
// @Produce json
// @Param domain path string true "Starter domain"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters/{domain} [get]
func (sh *StarterHandler) Find(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.FindStarterRequest

	// Auto validation with binding tag
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	// Call service
	result, err := sh.service.GetStarterByDomain(ctx, req.Domain)
	if err != nil {
		return nil, err
	}

	// Enrich with related data
	enrichedData, err := sh.service.EnrichStarters(ctx, []*domain.Starter{result})
	if err != nil {
		// Log error but continue with basic data
	}

	// Convert domain entity to enriched response DTO
	enrichedResponse := presentationDto.FromDomainEnriched(result, enrichedData)
	return enrichedResponse, nil
}

// [PATCH] /api/v1/starters/:domain
// UpdateStarter godoc
// @Summary Update starter
// @Description Partially updates starter information by domain.
// @Tags Starters
// @Accept json
// @Produce json
// @Param domain path string true "Starter domain"
// @Param request body StarterUpdateRequest true "Starter attributes to update"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters/{domain} [patch]
func (sh *StarterHandler) UpdateStarter(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.UpdateStarterRequest

	// Auto validation with binding tag (URI)
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	// Auto validation with binding tag (JSON body)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Convert Presentation DTO → Application Command
	command := applicationDto.UpdateStarterCommand{
		Domain:        req.Domain,
		Name:          req.Name,
		Email:         req.Email,
		Mobile:        req.Mobile,
		WorkPhone:     req.WorkPhone,
		JobTitle:      req.JobTitle,
		DepartmentID:  req.DepartmentID,
		LineManagerID: req.LineManagerID,
	}

	// Call service
	result, err := sh.service.UpdateStarter(ctx, command)
	if err != nil {
		return nil, err
	}

	// Enrich with related data
	enrichedData, err := sh.service.EnrichStarters(ctx, []*domain.Starter{result})
	if err != nil {
		// Log error but continue with basic data
	}

	// Convert domain entity to enriched response DTO
	return presentationDto.FromDomainEnriched(result, enrichedData), nil
}

// [DELETE] /api/v1/starters/:domain
// SoftDeleteStarter godoc
// @Summary Soft delete starter
// @Description Soft deletes a starter by domain.
// @Tags Starters
// @Accept json
// @Produce json
// @Param domain path string true "Starter domain"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /starters/{domain} [delete]
func (sh *StarterHandler) SoftDeleteStarter(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.DeleteStarterRequest

	// Auto validation with binding tag
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid domain parameter", err.Error())
	}

	// Call service to soft delete
	if err := sh.service.SoftDeleteStarter(ctx, req.Domain); err != nil {
		return nil, err
	}

	// Return success message
	return gin.H{
		"message": "Starter soft deleted successfully",
		"domain":  req.Domain,
	}, nil
}
