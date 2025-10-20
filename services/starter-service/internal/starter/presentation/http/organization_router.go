package http

import (
	"github.com/gin-gonic/gin"
	businessunitdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/businessunit"
	departmentdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/department"
)

// for swagger documents
var (
	_ departmentdto.CreateDepartmentRequest
	_ departmentdto.UpdateDepartmentRequest
	_ departmentdto.AssignLeaderRequest
	_ departmentdto.LeaderInfo
	_ departmentdto.DeleteDepartmentRequest
	_ departmentdto.DepartmentDetailResponse
	_ businessunitdto.BusinessUnitDetailResponse
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	departments := rg.Group("/departments")
	// @Summary List departments
	// @Description Retrieve departments with optional filters
	// @Tags Departments
	// @Produce json
	// @Param business_unit_id query int false "Filter by business unit"
	// @Param page query int false "Page number"
	// @Param limit query int false "Items per page"
	// @Success 200 {array} departmentdto.DepartmentDetailResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Router /departments [get]
	departments.GET("", handler.ListDepartments)

	// @Summary Get department detail
	// @Description Retrieve a department with nested details by ID
	// @Tags Departments
	// @Produce json
	// @Param id path int true "Department ID"
	// @Success 200 {object} departmentdto.DepartmentDetailResponse
	// @Failure 404 {object} responsepkg.APIError
	// @Router /departments/{id} [get]
	departments.GET("/:id", handler.GetDepartmentDetail)

	// @Summary Create department
	// @Description Create a new department
	// @Tags Departments
	// @Accept json
	// @Produce json
	// @Param request body departmentdto.CreateDepartmentRequest true "Department payload"
	// @Success 201 {object} departmentdto.DepartmentDetailResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 500 {object} responsepkg.APIError
	// @Router /departments [post]
	departments.POST("", handler.CreateDepartment)

	// @Summary Update department
	// @Description Update department information
	// @Tags Departments
	// @Accept json
	// @Produce json
	// @Param id path int true "Department ID"
	// @Param request body departmentdto.UpdateDepartmentRequest true "Update payload"
	// @Success 200 {object} departmentdto.DepartmentDetailResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 404 {object} responsepkg.APIError
	// @Router /departments/{id} [patch]
	departments.PATCH("/:id", handler.UpdateDepartment)

	// @Summary Assign department leader
	// @Description Assign or update the leader of a department
	// @Tags Departments
	// @Accept json
	// @Produce json
	// @Param id path int true "Department ID"
	// @Param request body departmentdto.AssignLeaderRequest true "Leader assignment payload"
	// @Success 200 {object} departmentdto.DepartmentDetailResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 404 {object} responsepkg.APIError
	// @Router /departments/{id}/leader [patch]
	departments.PATCH("/:id/leader", handler.AssignLeaderToDepartment)

	// @Summary Delete department
	// @Description Delete a department by ID
	// @Tags Departments
	// @Produce json
	// @Param id path int true "Department ID"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 404 {object} responsepkg.APIError
	// @Router /departments/{id} [delete]
	departments.DELETE("/:id", handler.DeleteDepartment)

	businessUnits := rg.Group("/business-units")
	// @Summary List business units
	// @Description Retrieve business units with pagination
	// @Tags Business Units
	// @Produce json
	// @Param page query int false "Page number"
	// @Param limit query int false "Items per page"
	// @Success 200 {array} businessunitdto.BusinessUnitDetailResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Router /business-units [get]
	businessUnits.GET("", handler.ListBusinessUnits)

	// @Summary Get business unit detail
	// @Description Retrieve a business unit with nested details by ID
	// @Tags Business Units
	// @Produce json
	// @Param id path int true "Business unit ID"
	// @Success 200 {object} businessunitdto.BusinessUnitDetailResponse
	// @Failure 404 {object} responsepkg.APIError
	// @Router /business-units/{id} [get]
	businessUnits.GET("/:id", handler.GetBusinessUnit)
}
