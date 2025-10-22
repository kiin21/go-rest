package http

import (
	"github.com/gin-gonic/gin"
	shareddto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/shared"
	starterdto "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http/dto/starter"
)

// for swagger documents
var (
	_ starterdto.CreateStarterRequest
	_ starterdto.UpdateStarterRequest
	_ starterdto.StarterResponse
	_ starterdto.StarterListAPIResponse
	_ starterdto.StarterAPIResponse
	_ starterdto.StarterDeleteAPIResponse
	_ shareddto.GenericAPIResponse
)

func RegisterStarterRoutes(rg *gin.RouterGroup, handler *StarterHandler) {
	route := rg.Group("/starters")

	route.POST("", handler.CreateStarter)

	route.GET("", handler.ListStarters)

	route.GET("/:domain", handler.Find)

	route.PATCH("/:domain", handler.UpdateStarter)

	route.DELETE("/:domain", handler.SoftDeleteStarter)
}
