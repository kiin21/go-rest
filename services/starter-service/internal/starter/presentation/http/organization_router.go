package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	org := rg.Group("/organization")
	
	departments := org.Group("/departments")
	departments.GET("", handler.ListDepartments)
	departments.GET("/:id", handler.GetDepartmentDetail)
	departments.POST("", handler.CreateDepartment)
	departments.PATCH("/:id", handler.UpdateDepartment)
	departments.PATCH("/:id/leader", handler.AssignLeaderToDepartment)
	departments.DELETE("/:id", handler.DeleteDepartment)
	
	businessUnits := org.Group("/business-units")
	businessUnits.GET("", handler.ListBusinessUnits)
	businessUnits.GET("/:id", handler.GetBusinessUnit)
}
