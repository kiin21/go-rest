package http

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/internal/organization/application"
	applicationDto "github.com/kiin21/go-rest/internal/organization/application/dto"
	presentationDto "github.com/kiin21/go-rest/internal/organization/interface/http/dto"
	"github.com/kiin21/go-rest/pkg/response"
)

type OrganizationHandler struct {
	departmentService   *application.DepartmentApplicationService
	businessUnitService *application.BusinessUnitApplicationService
}

// Constructor
func NewOrganizationHandler(
	departmentService *application.DepartmentApplicationService,
	businessUnitService *application.BusinessUnitApplicationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		departmentService:   departmentService,
		businessUnitService: businessUnitService,
	}
}

// [GET]: /api/v1/departments?business_unit_id=1&page=1&limit=10
// ListDepartments godoc
// @Summary List departments
// @Description Returns a paginated list of departments, optionally filtered by business unit.
// @Tags Departments
// @Accept json
// @Produce json
// @Param business_unit_id query int false "Filter by business unit ID"
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
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
		BusinessUnitID: req.BusinessUnitID,
	}

	// Call service
	result, err := h.departmentService.GetAllDepartments(ctx, query)
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

// [GET]: /api/v1/departments/:id
// GetDepartmentDetail godoc
// @Summary Get department detail
// @Description Retrieves department information by identifier.
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments/{id} [get]
func (h *OrganizationHandler) GetDepartmentDetail(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.GetDepartmentRequest
	fmt.Println("2345353sdgdfsgdsfgterwtr4: ", ctx.Request.RequestURI)
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	query := &applicationDto.GetDepartmentQuery{
		ID: &req.ID,
	}

	results, svcErr := h.departmentService.GetOneDepartment(ctx, *query)
	if svcErr != nil {
		return nil, svcErr
	}

	// Expecting one result for the given ID
	if len(results) == 0 {
		return nil, response.NewAPIError(http.StatusNotFound, "department not found", nil)
	}

	// Convert to response DTO
	return presentationDto.FromDomainWithDetails(results[0]), nil
}

// [GET]: /api/v1/business-units?page=1&limit=10
// ListBusinessUnits godoc
// @Summary List business units
// @Description Returns a paginated list of business units.
// @Tags Business Units
// @Accept json
// @Produce json
// @Param page query int false "Page number (default 1)"
// @Param limit query int false "Items per page (default 10)"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /business-units [get]
func (h *OrganizationHandler) ListBusinessUnits(ctx *gin.Context) (res interface{}, err error) {
	var req presentationDto.ListBusinessUnitsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := applicationDto.ListBusinessUnitsQuery{
		Pagination: response.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}

	result, err := h.businessUnitService.ListBusinessUnitsWithDetails(ctx, query)
	if err != nil {
		return nil, err
	}

	responseData := presentationDto.FromBusinessUnitsWithDetails(result.Data)
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

	return &response.PaginatedResult[*presentationDto.BusinessUnitDetailResponse]{
		Data:       responseData,
		Pagination: result.Pagination,
	}, nil
}

// [GET]: /api/v1/business-units/:id
// GetBusinessUnit godoc
// @Summary Get business unit detail
// @Description Retrieves business unit information by identifier.
// @Tags Business Units
// @Accept json
// @Produce json
// @Param id path int true "Business unit ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /business-units/{id} [get]
func (h *OrganizationHandler) GetBusinessUnit(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid business unit ID", err.Error())
	}

	unit, err := h.businessUnitService.GetBusinessUnitWithDetails(ctx, uriReq.ID)
	if err != nil {
		return nil, err
	}

	return presentationDto.FromBusinessUnitWithDetails(unit), nil
}

// [POST]: /api/v1/departments
// CreateDepartment godoc
// @Summary Create department
// @Description Creates a new department within a business unit.
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body DepartmentCreateRequest true "Department payload"
// @Success 200 {object} response.APIResponse{data=dto.DepartmentDetailResponse}
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
	result, err := h.departmentService.CreateDepartment(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	return presentationDto.FromDomainWithDetails(result), nil
}

// [PATCH]: /api/v1/departments/:id
// UpdateDepartment godoc
// @Summary Update department
// @Description Updates department attributes.
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body dto.UpdateDepartmentRequest true "Department attributes to update"
// @Success 200 {object} response.APIResponse{data=dto.DepartmentDetailResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments/{id} [patch]
func (h *OrganizationHandler) UpdateDepartment(ctx *gin.Context) (res interface{}, err error) {
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
	result, err := h.departmentService.UpdateDepartment(ctx, uriReq.ID, cmd)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	return presentationDto.FromDomainWithDetails(result), nil
}

// [PATCH]: /api/v1/departments/:id/leader
// AssignLeaderToDepartment godoc
// @Summary Assign leader to department
// @Description Assigns a leader to a department using either ID or domain (but not both).
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body dto.AssignLeaderRequest true "Leader assignment payload"
// @Success 200 {object} response.APIResponse{data=dto.DepartmentDetailResponse}
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments/{id}/leader [patch]
func (h *OrganizationHandler) AssignLeaderToDepartment(ctx *gin.Context) (res interface{}, err error) {
	// Parse ID from URI with auto validation
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	// Parse request body with auto validation
	var req presentationDto.AssignLeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	// Validate the request structure (either ID or domain, but not both)
	if err := req.Validate(); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid leader assignment", err.Error())
	}

	// Get leader identifier and type
	identifier, identifierType, err := req.GetLeaderIdentifier()
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid leader identifier", err.Error())
	}

	// Convert Presentation DTO → Application Command
	cmd := applicationDto.AssignLeaderCommand{
		DepartmentID:         uriReq.ID,
		LeaderIdentifier:     identifier,
		LeaderIdentifierType: identifierType,
	}

	// Set the appropriate field based on identifier type
	switch identifierType {
	case "id":
		if id, ok := identifier.(int64); ok {
			cmd.LeaderID = &id
		}
	case "domain":
		if domain, ok := identifier.(string); ok {
			cmd.LeaderDomain = &domain
		}
	}

	// Call service
	result, err := h.departmentService.AssignLeader(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// Convert domain entity to response DTO
	return presentationDto.FromDomainWithDetails(result), nil
}
