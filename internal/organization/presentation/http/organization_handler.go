package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	businessunitquery "github.com/kiin21/go-rest/internal/organization/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/internal/organization/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/internal/organization/application/dto/department/query"
	"github.com/kiin21/go-rest/internal/organization/application/service"
	businessunitdto "github.com/kiin21/go-rest/internal/organization/presentation/http/dto/businessunit"
	departmentdto "github.com/kiin21/go-rest/internal/organization/presentation/http/dto/department"
	"github.com/kiin21/go-rest/pkg/httpctx"
	"github.com/kiin21/go-rest/pkg/response"
)

type OrganizationHandler struct {
	orgService  *service.OrganizationApplicationService
	urlResolver httpctx.RequestURLResolver
}

func NewOrganizationHandler(
	orgService *service.OrganizationApplicationService,
	urlResolver httpctx.RequestURLResolver,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgService:  orgService,
		urlResolver: urlResolver,
	}
}

// ListDepartments godoc
// @Summary List departments
// @Description Retrieve departments with optional filters
// @Tags Departments
// @Produce json
// @Param business_unit_id query int false "Filter by business unit"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /departments [get]
func (h *OrganizationHandler) ListDepartments(ctx *gin.Context) {
	response.Wrap(h.listDepartments)(ctx)
}

func (h *OrganizationHandler) listDepartments(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.ListDepartmentsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := departmentquery.ListDepartmentsQuery{
		Pagination: response.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		BusinessUnitID: req.BusinessUnitID,
	}

	result, err := h.orgService.GetAllDepartments(ctx, query)
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to list departments")
	}

	responseData := departmentdto.FromDomainsWithDetails(result.Data)
	pagination := decoratePagination(ctx, h.urlResolver, result.Pagination)

	return &response.PaginatedResult[*departmentdto.DepartmentDetailResponse]{
		Data:       responseData,
		Pagination: pagination,
	}, nil
}

// GetDepartmentDetail godoc
// @Summary Get department detail
// @Description Retrieve a department with nested details by ID
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /departments/{id} [get]
func (h *OrganizationHandler) GetDepartmentDetail(ctx *gin.Context) {
	response.Wrap(h.getDepartmentDetail)(ctx)
}

func (h *OrganizationHandler) getDepartmentDetail(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.GetDepartmentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	query := &departmentquery.GetDepartmentQuery{
		ID: req.ID,
	}

	result, err := h.orgService.GetOneDepartment(ctx, *query)
	if err != nil {
		return nil, mapServiceError(err, "Department not found", "Failed to fetch department")
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// ListBusinessUnits godoc
// @Summary List business units
// @Description Retrieve business units with pagination
// @Tags Business Units
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /business-units [get]
func (h *OrganizationHandler) ListBusinessUnits(ctx *gin.Context) {
	response.Wrap(h.listBusinessUnits)(ctx)
}

func (h *OrganizationHandler) listBusinessUnits(ctx *gin.Context) (res interface{}, err error) {
	var req businessunitdto.ListBusinessUnitsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := businessunitquery.ListBusinessUnitsQuery{
		Pagination: response.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}

	result, err := h.orgService.ListBusinessUnitsWithDetails(ctx, query)
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to list business units")
	}

	responseData := businessunitdto.FromBusinessUnitsWithDetails(result.Data)
	pagination := decoratePagination(ctx, h.urlResolver, result.Pagination)

	return &response.PaginatedResult[*businessunitdto.BusinessUnitDetailResponse]{
		Data:       responseData,
		Pagination: pagination,
	}, nil
}

// GetBusinessUnit godoc
// @Summary Get business unit detail
// @Description Retrieve a business unit with nested details by ID
// @Tags Business Units
// @Produce json
// @Param id path int true "Business unit ID"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /business-units/{id} [get]
func (h *OrganizationHandler) GetBusinessUnit(ctx *gin.Context) {
	response.Wrap(h.getBusinessUnit)(ctx)
}

func (h *OrganizationHandler) getBusinessUnit(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid business unit ID", err.Error())
	}

	unit, err := h.orgService.GetBusinessUnitWithDetails(ctx, uriReq.ID)
	if err != nil {
		return nil, mapServiceError(err, "Business unit not found", "Failed to fetch business unit")
	}

	return businessunitdto.FromBusinessUnitWithDetails(unit), nil
}

// CreateDepartment godoc
// @Summary Create department
// @Description Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body departmentdto.CreateDepartmentRequest true "Department payload"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /departments [post]
func (h *OrganizationHandler) CreateDepartment(ctx *gin.Context) {
	response.Wrap(h.createDepartment)(ctx)
}

func (h *OrganizationHandler) createDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.CreateDepartmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	cmd := departmentcommand.CreateDepartmentCommand{
		FullName:          req.FullName,
		Shortname:         req.Shortname,
		BusinessUnitID:    req.BusinessUnitID,
		GroupDepartmentID: req.GroupDepartmentID,
		LeaderID:          req.LeaderID,
	}

	result, err := h.orgService.CreateDepartment(ctx, cmd)
	if err != nil {
		return nil, mapServiceError(err, "", "Failed to create department")
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// UpdateDepartment godoc
// @Summary Update department
// @Description Update department information
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body departmentdto.UpdateDepartmentRequest true "Update payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /departments/{id} [patch]
func (h *OrganizationHandler) UpdateDepartment(ctx *gin.Context) {
	response.Wrap(h.updateDepartment)(ctx)
}

func (h *OrganizationHandler) updateDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	var req departmentdto.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	cmd := departmentcommand.UpdateDepartmentCommand{
		FullName:          req.FullName,
		Shortname:         req.Shortname,
		BusinessUnitID:    req.BusinessUnitID,
		GroupDepartmentID: req.GroupDepartmentID,
		LeaderID:          req.LeaderID,
	}

	result, err := h.orgService.UpdateDepartment(ctx, uriReq.ID, cmd)
	if err != nil {
		return nil, mapServiceError(err, "Department not found", "Failed to update department")
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// AssignLeaderToDepartment godoc
// @Summary Assign department leader
// @Description Assign or update the leader of a department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param request body departmentdto.AssignLeaderRequest true "Leader assignment payload"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 404 {object} response.APIResponse
// @Router /departments/{id}/leader [patch]
func (h *OrganizationHandler) AssignLeaderToDepartment(ctx *gin.Context) {
	response.Wrap(h.assignLeaderToDepartment)(ctx)
}

func (h *OrganizationHandler) assignLeaderToDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	var req departmentdto.AssignLeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := req.Validate(); err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid leader assignment", err.Error())
	}

	identifier, identifierType, err := req.GetLeaderIdentifier()
	if err != nil {
		return nil, response.NewAPIError(http.StatusBadRequest, "Invalid leader identifier", err.Error())
	}

	cmd := departmentcommand.AssignLeaderCommand{
		DepartmentID:         uriReq.ID,
		LeaderIdentifier:     identifier,
		LeaderIdentifierType: identifierType,
	}

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

	result, err := h.orgService.AssignLeader(ctx, cmd)
	if err != nil {
		return nil, mapServiceError(err, "Department or leader not found", "Failed to assign leader")
	}

	return departmentdto.FromDomainWithDetails(result), nil
}
