package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	businessunitquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	businessunitdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/businessunit"
	departmentdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/department"
	shareddto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/shared"
)

var (
	_ shareddto.GenericAPIResponse
	_ departmentdto.DepartmentListAPIResponse
	_ departmentdto.DepartmentDetailAPIResponse
	_ departmentdto.DepartmentDeleteAPIResponse
)

type OrganizationHandler struct {
	orgService  *service.OrganizationApplicationService
	urlResolver httputil.RequestURLResolver
}

func NewOrganizationHandler(
	orgService *service.OrganizationApplicationService,
	urlResolver httputil.RequestURLResolver,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgService:  orgService,
		urlResolver: urlResolver,
	}
}

// ListDepartments godoc
// @Summary List departments
// @Description Retrieve departments with optional business unit filter and pagination
// @Tags Departments
// @Produce json
// @Param business_unit_id query int false "Filter by business unit ID" minimum(1)
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(10)
// @Success 200 {object} departmentdto.DepartmentListAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments [get]
func (h *OrganizationHandler) ListDepartments(ctx *gin.Context) {
	httputil.Wrap(h.listDepartments)(ctx)
}

func (h *OrganizationHandler) listDepartments(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.ListDepartmentsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := departmentquery.ListDepartmentsQuery{
		Pagination: httputil.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
		BusinessUnitID: req.BusinessUnitID,
	}

	result, err := h.orgService.GetAllDepartments(ctx, query)
	if err != nil {
		return nil, err
	}

	responseData := departmentdto.FromDomainsWithDetails(result.Data)
	pagination := httputil.DecoratePagination(ctx, h.urlResolver, result.Pagination)
	return &httputil.PaginatedResult[*departmentdto.DepartmentDetailResponse]{
		Data:       responseData,
		Pagination: pagination,
	}, nil
}

// GetDepartmentDetail godoc
// @Summary Get department detail
// @Description Retrieve a department with nested details by ID
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Success 200 {object} departmentdto.DepartmentDetailAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments/{id} [get]
func (h *OrganizationHandler) GetDepartmentDetail(ctx *gin.Context) {
	httputil.Wrap(h.getDepartmentDetail)(ctx)
}

func (h *OrganizationHandler) getDepartmentDetail(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.GetDepartmentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	query := &departmentquery.GetDepartmentQuery{
		ID: req.ID,
	}

	result, err := h.orgService.GetOneDepartment(ctx, *query)
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// ListBusinessUnits godoc
// @Summary List business units
// @Description Retrieve business units with pagination
// @Tags Business Units
// @Produce json
// @Param page query int false "Page number" minimum(1) default(1)
// @Param limit query int false "Page size" minimum(1) maximum(100) default(10)
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /business-units [get]
func (h *OrganizationHandler) ListBusinessUnits(ctx *gin.Context) {
	httputil.Wrap(h.listBusinessUnits)(ctx)
}

func (h *OrganizationHandler) listBusinessUnits(ctx *gin.Context) (res interface{}, err error) {
	var req businessunitdto.ListBusinessUnitsRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}
	req.SetDefaults()

	query := businessunitquery.ListBusinessUnitsQuery{
		Pagination: httputil.ReqPagination{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}

	result, err := h.orgService.ListBusinessUnitsWithDetails(ctx, query)
	if err != nil {
		return nil, err
	}

	responseData := businessunitdto.FromBusinessUnitsWithDetails(result.Data)
	pagination := httputil.DecoratePagination(ctx, h.urlResolver, result.Pagination)

	return &httputil.PaginatedResult[*businessunitdto.BusinessUnitDetailResponse]{
		Data:       responseData,
		Pagination: pagination,
	}, nil
}

// GetBusinessUnit godoc
// @Summary Get business unit detail
// @Description Retrieve a business unit with nested details by ID
// @Tags Business Units
// @Produce json
// @Param id path int true "Business unit ID" minimum(1)
// @Success 200 {object} businessunitdto.BusinessUnitDetailAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /business-units/{id} [get]
func (h *OrganizationHandler) GetBusinessUnit(ctx *gin.Context) {
	httputil.Wrap(h.getBusinessUnit)(ctx)
}

func (h *OrganizationHandler) getBusinessUnit(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid business unit ID", err.Error())
	}

	unit, err := h.orgService.GetBusinessUnitWithDetails(ctx, uriReq.ID)
	if err != nil {
		return nil, err
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
// @Success 200 {object} departmentdto.DepartmentDetailAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments [post]
func (h *OrganizationHandler) CreateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.createDepartment)(ctx)
}

func (h *OrganizationHandler) createDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.CreateDepartmentRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
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
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// UpdateDepartment godoc
// @Summary Update department
// @Description Update department information by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Param request body departmentdto.UpdateDepartmentRequest true "Update payload"
// @Success 200 {object} departmentdto.DepartmentDetailAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments/{id} [patch]
func (h *OrganizationHandler) UpdateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.updateDepartment)(ctx)
}

func (h *OrganizationHandler) updateDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	var req departmentdto.UpdateDepartmentRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	cmd := departmentcommand.UpdateDepartmentCommand{
		ID:                uriReq.ID,
		FullName:          req.FullName,
		Shortname:         req.Shortname,
		BusinessUnitID:    req.BusinessUnitID,
		GroupDepartmentID: req.GroupDepartmentID,
		LeaderID:          req.LeaderID,
	}

	result, err := h.orgService.UpdateDepartment(ctx, cmd)
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// AssignLeaderToDepartment godoc
// @Summary Assign department leader
// @Description Assign or update the leader of a department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Param request body departmentdto.AssignLeaderRequest true "Leader assignment payload"
// @Success 200 {object} departmentdto.DepartmentDetailAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments/{id}/leader [patch]
func (h *OrganizationHandler) AssignLeaderToDepartment(ctx *gin.Context) {
	httputil.Wrap(h.assignLeaderToDepartment)(ctx)
}

func (h *OrganizationHandler) assignLeaderToDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := ctx.ShouldBindUri(&uriReq); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid department ID", err.Error())
	}

	var req departmentdto.AssignLeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid request body", err.Error())
	}

	if err := req.Validate(); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid leader assignment", err.Error())
	}

	identifier, identifierType, err := req.GetLeaderIdentifier()
	if err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid leader identifier", err.Error())
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
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// DeleteDepartment godoc
// @Summary Delete department
// @Description Delete a department by ID
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Success 200 {object} departmentdto.DepartmentDeleteAPIResponse
// @Failure 400 {object} shareddto.GenericAPIResponse
// @Failure 404 {object} shareddto.GenericAPIResponse
// @Failure 500 {object} shareddto.GenericAPIResponse
// @Router /departments/{id} [delete]
func (h *OrganizationHandler) DeleteDepartment(ctx *gin.Context) {
	httputil.Wrap(h.deleteDepartment)(ctx)
}

func (h *OrganizationHandler) deleteDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.DeleteDepartmentRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		return nil, httputil.NewAPIError(http.StatusBadRequest, "Invalid query parameters", err.Error())
	}

	query := &departmentcommand.DeleteDepartmentCommand{
		ID: req.ID,
	}

	err = h.orgService.DeleteDepartment(ctx, *query)
	if err != nil {
		return nil, err
	}

	return gin.H{
		"message": "Department deleted successfully",
		"id":      req.ID,
	}, nil
}
