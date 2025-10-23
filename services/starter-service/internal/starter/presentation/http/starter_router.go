package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterStarterRoutes(rg *gin.RouterGroup, handler *StarterHandler) {
	route := rg.Group("/starters")
	route.POST("", handler.CreateStarter)
	route.GET("", handler.ListStarters)
	route.GET("/:domain", handler.Find)
	route.PATCH("/:domain", handler.UpdateStarter)
	route.DELETE("/:domain", handler.SoftDeleteStarter)
}
