package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	budto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/businessunit"
	departmentdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/department"
)

type OrganizationHandler struct {
	orgSvc *service.OrganizationApplicationService
}

func NewOrganizationHandler(
	orgSvc *service.OrganizationApplicationService,
) *OrganizationHandler {
	return &OrganizationHandler{
		orgSvc: orgSvc,
	}
}

// ListDepartments godoc
// @Summary List departments
// @Description Retrieve departments with optional business unit filter and pagination
// @Tags Departments
// @Accept json
// @Produce json
// @Param business_unit_id query int false "Filter by business unit ID" minimum(1)
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Page size" default(10) minimum(1) maximum(100)
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments [get]
func (h *OrganizationHandler) ListDepartments(ctx *gin.Context) {
	httputil.Wrap(h.listDepartments)(ctx)
}

func (h *OrganizationHandler) listDepartments(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.ListDepartmentsRequest

	if err := httputil.ValidateQuery(ctx, &req); err != nil {
		return nil, err
	}
	req.SetDefaults()

	result, err := h.orgSvc.GetAllDepartments(ctx, req.ToQuery())
	if err != nil {
		return nil, err
	}

	data := departmentdto.FromDomainsWithDetails(result.Data)
	pagination := httputil.CursorPagination(ctx, result.Pagination)
	return &httputil.PaginatedResult[*departmentdto.DepartmentDetailResponse]{
		Data:       data,
		Pagination: pagination,
	}, nil
}

// GetDepartmentDetail godoc
// @Summary Get department detail
// @Description Retrieve a department with nested details by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 404 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments/{id} [get]
func (h *OrganizationHandler) GetDepartmentDetail(ctx *gin.Context) {
	httputil.Wrap(h.getDepartmentDetail)(ctx)
}

func (h *OrganizationHandler) getDepartmentDetail(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.GetOneDepartment(ctx, uriReq.ID)
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// ListBusinessUnits godoc
// @Summary List business units
// @Description Retrieve business units with pagination
// @Tags Business Units
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1) minimum(1)
// @Param limit query int false "Page size" default(10) minimum(1) maximum(100)
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/business-units [get]
func (h *OrganizationHandler) ListBusinessUnits(ctx *gin.Context) {
	httputil.Wrap(h.listBusinessUnits)(ctx)
}

func (h *OrganizationHandler) listBusinessUnits(ctx *gin.Context) (res interface{}, err error) {
	var req budto.ListBusinessUnitsRequest

	if err := httputil.ValidateQuery(ctx, &req); err != nil {
		return nil, err
	}
	req.SetDefaults()

	result, err := h.orgSvc.ListBusinessUnitsWithDetails(ctx, req.ToQuery())
	if err != nil {
		return nil, err
	}

	data := budto.FromBusinessUnitsWithDetails(result.Data)
	pagination := httputil.CursorPagination(ctx, result.Pagination)

	return &httputil.PaginatedResult[*budto.BusinessUnitDetailResponse]{
		Data:       data,
		Pagination: pagination,
	}, nil
}

// GetBusinessUnit godoc
// @Summary Get business unit detail
// @Description Retrieve a business unit with nested details by ID
// @Tags Business Units
// @Accept json
// @Produce json
// @Param id path int true "Business unit ID" minimum(1)
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 404 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/business-units/{id} [get]
func (h *OrganizationHandler) GetBusinessUnit(ctx *gin.Context) {
	httputil.Wrap(h.getBusinessUnit)(ctx)
}

func (h *OrganizationHandler) getBusinessUnit(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		ID int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	unit, err := h.orgSvc.GetBusinessUnitWithDetails(ctx, uriReq.ID)
	if err != nil {
		return nil, err
	}

	return budto.FromBusinessUnitWithDetails(unit), nil
}

// CreateDepartment godoc
// @Summary Create department
// @Description Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body department.CreateDepartmentRequest true "Department payload"
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments [post]
func (h *OrganizationHandler) CreateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.createDepartment)(ctx)
}

func (h *OrganizationHandler) createDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.CreateDepartmentRequest
	if err := httputil.ValidateBody(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.CreateDepartment(ctx, req.ToCommand())
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
// @Param request body department.UpdateDepartmentRequest true "Update payload"
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 404 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments/{id} [patch]
func (h *OrganizationHandler) UpdateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.updateDepartment)(ctx)
}

func (h *OrganizationHandler) updateDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		DeptId int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}
	var req departmentdto.UpdateDepartmentRequest
	if err := httputil.ValidateBody(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.UpdateDepartment(ctx, req.ToCommand(uriReq.DeptId))
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
// @Param request body department.AssignLeaderRequest true "Leader assignment payload"
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 404 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments/{id}/leader [patch]
func (h *OrganizationHandler) AssignLeaderToDepartment(ctx *gin.Context) {
	httputil.Wrap(h.assignLeaderToDepartment)(ctx)
}

func (h *OrganizationHandler) assignLeaderToDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		DeptId int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}
	var req departmentdto.AssignLeaderRequest
	if err := httputil.ValidateBody(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.AssignLeader(ctx, req.ToCommand(uriReq.DeptId))
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// DeleteDepartment godoc
// @Summary Delete department
// @Description Delete a department by ID
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID" minimum(1)
// @Success 200 {object} httputil.APIResponse
// @Failure 400 {object} httputil.APIResponse
// @Failure 404 {object} httputil.APIResponse
// @Failure 500 {object} httputil.APIResponse
// @Router /organization/departments/{id} [delete]
func (h *OrganizationHandler) DeleteDepartment(ctx *gin.Context) {
	httputil.Wrap(h.deleteDepartment)(ctx)
}

func (h *OrganizationHandler) deleteDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		Id int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	err = h.orgSvc.DeleteDepartment(ctx, uriReq.Id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"message": "DepartmentName deleted successfully",
		"id":      uriReq.Id,
	}, nil
}
