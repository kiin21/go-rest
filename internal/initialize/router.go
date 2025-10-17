package initialize

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	initOrg "github.com/kiin21/go-rest/internal/initialize/organization"
	initStarter "github.com/kiin21/go-rest/internal/initialize/starter"
	"github.com/kiin21/go-rest/internal/middleware"
	orgHttp "github.com/kiin21/go-rest/internal/organization/presentation/http"
	messagingKafka "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker/kafka"
	starterRepository "github.com/kiin21/go-rest/internal/starter/infrastructure/persistence/repository"
	starterHttp "github.com/kiin21/go-rest/internal/starter/presentation/http"
	"github.com/kiin21/go-rest/pkg/httpctx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(
	db *gorm.DB,
	esClient *elasticsearch.Client,
	isLogger string,
	publicBaseURL string,
	kafkaBrokers string,
	kafkaTopic string,
	kafkaConsumerGroup string,
) *gin.Engine {
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

	r.GET("/ping/100", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	// Health check endpoint
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// Swagger UI
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/api/v1")

	requestURLResolver := httpctx.NewRequestURLResolver(publicBaseURL)
	starterRepo := starterRepository.NewMySQLStarterRepository(db)

	// Register the organization routes (get department repo and business unit repo for sharing)
	orgHandler, departmentRepo := initOrg.InitOrganization(db, requestURLResolver, starterRepo)
	orgHttp.RegisterOrganizationRoutes(v1, orgHandler)

	starterHandler, kafkaConsumer := initStarter.InitStarter(
		esClient,
		starterRepo,
		departmentRepo,
		requestURLResolver,
		kafkaBrokers,
		kafkaTopic,
		kafkaConsumerGroup,
	)
	starterHttp.RegisterStarterRoutes(v1, starterHandler)

	setupGracefulShutdown(kafkaConsumer)

	return r
}

// setupGracefulShutdown handles graceful shutdown of Kafka consumer
func setupGracefulShutdown(kafkaConsumer *messagingKafka.Consumer) {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigterm
		log.Println("Shutting down gracefully...")

		// Stop Kafka consumer
		if kafkaConsumer != nil {
			kafkaConsumer.Stop()
		}

		log.Println("Shutdown complete")
		os.Exit(0)
	}()
}
