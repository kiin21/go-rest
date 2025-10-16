package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/response"
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	departments := rg.Group("/departments")
	departments.GET("", response.Wrap(handler.ListDepartments))
	departments.GET("/:id", response.Wrap(handler.GetDepartmentDetail))
	departments.POST("", response.Wrap(handler.CreateDepartment))
	departments.PATCH("/:id", response.Wrap(handler.UpdateDepartment))
	departments.PATCH("/:id/leader", response.Wrap(handler.AssignLeaderToDepartment))

	businessUnits := rg.Group("/business-units")
	businessUnits.GET("", response.Wrap(handler.ListBusinessUnits))
	businessUnits.GET("/:id", response.Wrap(handler.GetBusinessUnit))
}
