package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	application "github.com/kiin21/go-rest/internal/organization/application"
	applicationDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	presentationDto "github.com/kiin21/go-rest/internal/organization/interface/http/dto"
	"github.com/kiin21/go-rest/pkg/response"
)

type OrganizationHandler struct {
	organizationService *application.OrganizationApplicationService
}

// Constructor
func NewOrganizationHandler(
	organizationService *application.OrganizationApplicationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		organizationService: organizationService,
	}
}

// [GET]: /api/v1/departments?business_unit_id=1&page=1&limit=10&include_subdepartments=false
// ListDepartments godoc
// @Summary List departments
// @Description Returns a paginated list of departments with optional business unit filter.
// @Tags Departments
// @Accept json
// @Produce json
// @Param business_unit_id query int false "Filter by business unit ID"
// @Param include_subdepartments query bool false "Include nested subdepartments"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 20)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments [get]
func (h *OrganizationHandler) ListDepartments(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.ListDepartmentsRequest

	// Auto validation with binding tag
	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	// Convert Presentation DTO → Application Query
	query := applicationDto.ListDepartmentsQuery{
		Pagination: response.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		BusinessUnitID:        req.BusinessUnitID,
		IncludeSubdepartments: req.IncludeSubdepartments,
	}

	// Call service
	result, err := h.organizationService.GetAllDepartments(ctx, query)
	if err != nil {
		return nil, err
	}

	// Convert domain entities with details to response DTOs
	responseData := presentationDto.FromDomainsWithDetails(result.Data)
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

	return &response.PaginatedResult[*presentationDto.DepartmentDetailResponse]{
		Data:       responseData,
		Pagination: result.Pagination,
	}, nil
}

// [POST]: /api/v1/departments
// CreateDepartment godoc
// @Summary Create department
// @Description Creates a new department within a business unit.
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body DepartmentCreateRequest true "Department payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments [post]
func (h *OrganizationHandler) CreateDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.CreateDepartmentRequest

	// Auto validation with binding tag
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Convert Presentation DTO → Application Command
	cmd := applicationDto.CreateDepartmentCommand{
		FullName:          req.FullName,
		Shortname:         req.Shortname,
		BusinessUnitID:    req.BusinessUnitID,
		GroupDepartmentID: req.GroupDepartmentID,
		LeaderID:          req.LeaderID,
	}

	// Call service
	result, err := h.organizationService.CreateDepartment(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	return presentationDto.FromDepartment(result), nil
}

// [PATCH]: /api/v1/departments/:id/leader
// AssignLeaderForDepartment godoc
// @Summary Update department
// @Description Updates department details by identifier.
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body DepartmentUpdateRequest true "Department attributes to update"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments/{id} [patch]
func (h *OrganizationHandler) AssignLeaderForDepartment(ctx *gin.Context) (res interface{}, err error) {
	// Parse ID from URI with auto validation
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	// Parse request body with auto validation
	var req presentationDto.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Convert Presentation DTO → Application Command
	cmd := applicationDto.UpdateDepartmentCommand{
		FullName:          req.FullName,
		Shortname:         req.Shortname,
		BusinessUnitID:    req.BusinessUnitID,
		GroupDepartmentID: req.GroupDepartmentID,
		LeaderID:          req.LeaderID,
	}

	// Call service
	result, err := h.organizationService.UpdateDepartment(ctx, uriReq.ID, cmd)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	return presentationDto.FromDepartment(result), nil
}
