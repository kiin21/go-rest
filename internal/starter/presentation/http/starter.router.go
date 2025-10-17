package http

import (
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/response"
)

func RegisterStarterRoutes(rg *gin.RouterGroup, handler *StarterHandler) {
	route := rg.Group("/starters")
	route.POST("", response.Wrap(handler.CreateStarter))
	route.GET("", response.Wrap(handler.ListStarters))
	route.GET("/:domain", response.Wrap(handler.Find))
	route.PATCH("/:domain", response.Wrap(handler.UpdateStarter))
	route.DELETE("/:domain", response.Wrap(handler.SoftDeleteStarter))
}
