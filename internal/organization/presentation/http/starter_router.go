package http

import (
	"github.com/gin-gonic/gin"
	starterdto "github.com/kiin21/go-rest/internal/organization/presentation/http/dto/starter"
)

// for swagger documents
var (
	_ starterdto.CreateStarterRequest
	_ starterdto.UpdateStarterRequest
	_ starterdto.StarterResponse
)

func RegisterStarterRoutes(rg *gin.RouterGroup, handler *StarterHandler) {
	route := rg.Group("/starters")

	// @Summary Create a starter
	// @Description Create a new starter record
	// @Tags Starters
	// @Accept json
	// @Produce json
	// @Param request body starterdto.CreateStarterRequest true "Starter payload"
	// @Success 201 {object} starterdto.StarterResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 409 {object} responsepkg.APIError
	// @Failure 500 {object} responsepkg.APIError
	// @Router /starters [post]
	route.POST("", handler.CreateStarter)

	// @Summary List starters
	// @Description Retrieve starters with optional filters and pagination
	// @Tags Starters
	// @Produce json
	// @Param business_unit_id query int false "Filter by business unit"
	// @Param department_id query int false "Filter by department"
	// @Param q query string false "Keyword search"
	// @Param sort_by query string false "Sort field"
	// @Param sort_order query string false "Sort order"
	// @Param page query int false "Page number"
	// @Param limit query int false "Items per page"
	// @Success 200 {array} starterdto.StarterResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Router /starters [get]
	route.GET("", handler.ListStarters)

	// @Summary Get starter detail
	// @Description Get a starter by domain
	// @Tags Starters
	// @Produce json
	// @Param domain path string true "Starter domain"
	// @Success 200 {object} starterdto.StarterResponse
	// @Failure 404 {object} responsepkg.APIError
	// @Router /starters/{domain} [get]
	route.GET("/:domain", handler.Find)

	// @Summary Update starter
	// @Description Partially update a starter by domain
	// @Tags Starters
	// @Accept json
	// @Produce json
	// @Param domain path string true "Starter domain"
	// @Param request body starterdto.UpdateStarterRequest true "Update payload"
	// @Success 200 {object} starterdto.StarterResponse
	// @Failure 400 {object} responsepkg.APIError
	// @Failure 404 {object} responsepkg.APIError
	// @Router /starters/{domain} [patch]
	route.PATCH("/:domain", handler.UpdateStarter)

	// @Summary Delete starter
	// @Description Soft delete a starter by domain
	// @Tags Starters
	// @Produce json
	// @Param domain path string true "Starter domain"
	// @Success 200 {object} map[string]interface{}
	// @Failure 404 {object} responsepkg.APIError
	// @Router /starters/{domain} [delete]
	route.DELETE("/:domain", handler.SoftDeleteStarter)
}
