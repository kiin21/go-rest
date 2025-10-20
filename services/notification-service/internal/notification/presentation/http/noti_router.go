package http

import "github.com/gin-gonic/gin"

func RegisterNotificationRoutes(rg *gin.RouterGroup, handler *NotiHandler) {
	route := rg.Group("/notifications")

	route.GET("", handler.GetList)
}
