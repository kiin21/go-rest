package initialize

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/docs"
	"github.com/kiin21/go-rest/services/starter-service/internal/middleware"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(
	logLevel string,
	requestURLResolver *httputil.RequestURLResolver,
	orgHandler *orgHttp.OrganizationHandler,
	starterHandler *orgHttp.StarterHandler,
) *gin.Engine {
	var router *gin.Engine
	if logLevel == "debug" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		router = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	}

	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)

	router.Use(middleware.CORS)

	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.GET("/swagger/*any", func(ctx *gin.Context) {
		docs.SwaggerInfo.Schemes = []string{requestURLResolver.Scheme(ctx)}
		if host := requestURLResolver.Host(ctx); host != "" {
			docs.SwaggerInfo.Host = host
		}
		swaggerHandler(ctx)
	})

	v1 := router.Group("/api/v1")

	orgHttp.RegisterOrganizationRoutes(v1, orgHandler)
	orgHttp.RegisterStarterRoutes(v1, starterHandler)

	return router
}
