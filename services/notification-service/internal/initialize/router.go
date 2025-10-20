package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/notification-service/docs"
	"github.com/kiin21/go-rest/services/notification-service/internal/middleware"
	notihttp "github.com/kiin21/go-rest/services/notification-service/internal/notification/presentation/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(logLevel string, requestURLResolver httputil.RequestURLResolver, handler *notihttp.NotiHandler) *gin.Engine {
	var r *gin.Engine
	// Set the mode based on the environment
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}
	// middlewares
	r.Use(middleware.CORS)
	// Health check endpoint
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	v1 := r.Group("/api/v1")

	// Swagger UI
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", func(ctx *gin.Context) {
		docs.SwaggerInfo.Schemes = []string{requestURLResolver.Scheme(ctx)}
		if host := requestURLResolver.Host(ctx); host != "" {
			docs.SwaggerInfo.Host = host
		}
		swaggerHandler(ctx)
	})

	notihttp.RegisterNotificationRoutes(v1, handler)

	return r
}
