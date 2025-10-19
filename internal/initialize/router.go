package initialize

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/docs"
	initOrg "github.com/kiin21/go-rest/internal/initialize/organization"
	"github.com/kiin21/go-rest/internal/middleware"
	orgInfra "github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/repository/mysql"
	orgHttp "github.com/kiin21/go-rest/internal/organization/presentation/http"
	"github.com/kiin21/go-rest/pkg/httpctx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, esClient *elasticsearch.Client, isLogger string, kafkaBrokers string, kafkaTopic string, kafkaConsumerGroup string) *gin.Engine {
	var r *gin.Engine
	// Set the mode based on the environment
	if isLogger == "debug" {
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

	requestURLResolver := httpctx.NewRequestURLResolver()

	// Swagger UI
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	r.GET("/swagger/*any", func(ctx *gin.Context) {
		docs.SwaggerInfo.Schemes = []string{requestURLResolver.Scheme(ctx)}
		if host := requestURLResolver.Host(ctx); host != "" {
			docs.SwaggerInfo.Host = host
		}
		swaggerHandler(ctx)
	})

	starterRepo := orgInfra.NewStarterRepository(db)
	businessUnitRepo := orgInfra.NewBusinessUnitRepository(db)
	departmentRepo := orgInfra.NewDepartmentRepository(db)

	orgHandler := initOrg.InitOrganization(requestURLResolver, starterRepo, departmentRepo, businessUnitRepo)
	orgHttp.RegisterOrganizationRoutes(v1, orgHandler)

	starterHandler := initOrg.InitStarter(esClient, starterRepo, departmentRepo, businessUnitRepo, requestURLResolver, kafkaBrokers, kafkaTopic, kafkaConsumerGroup)
	orgHttp.RegisterStarterRoutes(v1, starterHandler)

	return r
}
