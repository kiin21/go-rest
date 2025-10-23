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

// ListDepartments GET /api/v1/organization/departments
func (h *OrganizationHandler) ListDepartments(ctx *gin.Context) {
	httputil.Wrap(h.listDepartments)(ctx)
}

func (h *OrganizationHandler) listDepartments(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.ListDepartmentsRequest

	if err := httputil.ValidateReq(ctx, &req); err != nil {
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

// GetDepartmentDetail GET /api/v1/organization/departments/:id
func (h *OrganizationHandler) GetDepartmentDetail(ctx *gin.Context) {
	httputil.Wrap(h.getDepartmentDetail)(ctx)
}

func (h *OrganizationHandler) getDepartmentDetail(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		id int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.GetOneDepartment(ctx, uriReq.id)
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// ListBusinessUnits GET /api/v1/organization/business-units
func (h *OrganizationHandler) ListBusinessUnits(ctx *gin.Context) {
	httputil.Wrap(h.listBusinessUnits)(ctx)
}

func (h *OrganizationHandler) listBusinessUnits(ctx *gin.Context) (res interface{}, err error) {
	var req budto.ListBusinessUnitsRequest

	if err := httputil.ValidateReq(ctx, &req); err != nil {
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

// GetBusinessUnit GET /api/v1/organization/business-units/:id
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

// CreateDepartment POST /api/v1/organization/departments
func (h *OrganizationHandler) CreateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.createDepartment)(ctx)
}

func (h *OrganizationHandler) createDepartment(ctx *gin.Context) (res interface{}, err error) {
	var req departmentdto.CreateDepartmentRequest
	if err := httputil.ValidateReq(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.CreateDepartment(ctx, req.ToCommand())
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// UpdateDepartment PATCH /api/v1/organization/departments/:id
func (h *OrganizationHandler) UpdateDepartment(ctx *gin.Context) {
	httputil.Wrap(h.updateDepartment)(ctx)
}

func (h *OrganizationHandler) updateDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		deptId int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}
	var req departmentdto.UpdateDepartmentRequest
	if err := httputil.ValidateReq(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.UpdateDepartment(ctx, req.ToCommand(uriReq.deptId))
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// AssignLeaderToDepartment PATCH /api/v1/organization/departments/:id/leader
func (h *OrganizationHandler) AssignLeaderToDepartment(ctx *gin.Context) {
	httputil.Wrap(h.assignLeaderToDepartment)(ctx)
}

func (h *OrganizationHandler) assignLeaderToDepartment(ctx *gin.Context) (res interface{}, err error) {
	var uriReq struct {
		deptId int64 `uri:"id" binding:"required,min=1"`
	}
	if err := httputil.ValidateURI(ctx, &uriReq); err != nil {
		return nil, err
	}
	var req departmentdto.AssignLeaderRequest
	if err := httputil.ValidateReq(ctx, &req); err != nil {
		return nil, err
	}

	result, err := h.orgSvc.AssignLeader(ctx, req.ToCommand(uriReq.deptId))
	if err != nil {
		return nil, err
	}

	return departmentdto.FromDomainWithDetails(result), nil
}

// DeleteDepartment DELETE /api/v1/organization/departments/:id
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
