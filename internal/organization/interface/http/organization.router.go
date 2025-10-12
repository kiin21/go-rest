package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/response"
)

func RegisterOrganizationRoutes(rg *gin.RouterGroup, handler *OrganizationHandler) {
	departments := rg.Group("/departments")
	departments.GET("", response.Wrap(handler.ListDepartments))
	departments.POST("", response.Wrap(handler.CreateDepartment))
	departments.PATCH("/:id/leader", response.Wrap(handler.AssignLeaderForDepartment))
}
