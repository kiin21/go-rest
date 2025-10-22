package http

import (
	"github.com/gin-gonic/gin"
	businessunitdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/businessunit"
	departmentdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/department"
	shareddto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/shared"
)

// for swagger documents
var (
	_ departmentdto.CreateDepartmentRequest
	_ departmentdto.UpdateDepartmentRequest
	_ departmentdto.AssignLeaderRequest
	_ departmentdto.LeaderInfo
	_ departmentdto.DeleteDepartmentRequest
	_ departmentdto.DepartmentDetailResponse
	_ departmentdto.DepartmentListAPIResponse
	_ departmentdto.DepartmentDetailAPIResponse
	_ departmentdto.DepartmentDeleteAPIResponse
	_ businessunitdto.BusinessUnitDetailResponse
	_ businessunitdto.BusinessUnitListAPIResponse
	_ businessunitdto.BusinessUnitDetailAPIResponse
	_ shareddto.GenericAPIResponse
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	departments := rg.Group("/departments")
	departments.GET("", handler.ListDepartments)
	departments.GET("/:id", handler.GetDepartmentDetail)
	departments.POST("", handler.CreateDepartment)
	departments.PATCH("/:id", handler.UpdateDepartment)
	departments.PATCH("/:id/leader", handler.AssignLeaderToDepartment)
	departments.DELETE("/:id", handler.DeleteDepartment)
	businessUnits := rg.Group("/business-units")
	businessUnits.GET("", handler.ListBusinessUnits)
	businessUnits.GET("/:id", handler.GetBusinessUnit)
}
