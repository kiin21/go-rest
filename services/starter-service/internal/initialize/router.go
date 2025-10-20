package initialize

import (
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/docs"
	initStarter "github.com/kiin21/go-rest/services/starter-service/internal/initialize/starter"
	"github.com/kiin21/go-rest/services/starter-service/internal/middleware"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	persistentMySQL "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/repository/mysql"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(
	db *gorm.DB,
	esClient *elasticsearch.Client,
	isLogger string,
	notificationPublisher domainmessaging.NotificationPublisher,
	kafkaBrokers string,
	kafkaSyncTopic string,
	kafkaConsumerGroup string,
) *gin.Engine {
	router := newGinEngine(isLogger)
	requestURLResolver := httputil.NewRequestURLResolver()

	registerHealthCheck(router)
	registerSwaggerRoutes(router, requestURLResolver)

	v1 := router.Group("/api/v1")

	starterRepo := persistentMySQL.NewStarterRepository(db)
	businessUnitRepo := persistentMySQL.NewBusinessUnitRepository(db)
	departmentRepo := persistentMySQL.NewDepartmentRepository(db)

	orgHandler := initStarter.InitOrganization(
		requestURLResolver,
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		notificationPublisher,
	)

	starterHandler := initStarter.InitStarter(
		esClient,
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		requestURLResolver,
		kafkaBrokers,
		kafkaSyncTopic,
		kafkaConsumerGroup,
	)

	registerAPIRoutes(v1, orgHandler, starterHandler)

	return router
}

func newGinEngine(isLogger string) *gin.Engine {
	var router *gin.Engine
	if isLogger == "debug" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		router = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
	}
	router.Use(middleware.CORS)
	return router
}

func registerHealthCheck(router *gin.Engine) {
	router.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
}

func registerSwaggerRoutes(router *gin.Engine, requestURLResolver httputil.RequestURLResolver) {
	swaggerHandler := ginSwagger.WrapHandler(swaggerFiles.Handler)
	router.GET("/swagger/*any", func(ctx *gin.Context) {
		docs.SwaggerInfo.Schemes = []string{requestURLResolver.Scheme(ctx)}
		if host := requestURLResolver.Host(ctx); host != "" {
			docs.SwaggerInfo.Host = host
		}
		swaggerHandler(ctx)
	})
}

func registerAPIRoutes(
	group *gin.RouterGroup,
	organizationHandler *orgHttp.OrganizationHandler,
	starterHandler *orgHttp.StarterHandler,
) {
	orgHttp.RegisterOrganizationRoutes(group, organizationHandler)
	orgHttp.RegisterStarterRoutes(group, starterHandler)
}
