package initialize

import (
	"log"
	"net/http"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	"github.com/kiin21/go-rest/services/starter-service/docs"
	initOrg "github.com/kiin21/go-rest/services/starter-service/internal/initialize/organization"
	"github.com/kiin21/go-rest/services/starter-service/internal/middleware"
	messageBroker "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
	orgInfra "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/repository/mysql"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func InitRouter(db *gorm.DB, esClient *elasticsearch.Client, isLogger string, kafkaBrokers string, kafkaSyncTopic string, kafkaNotificationTopic string, kafkaConsumerGroup string) *gin.Engine {
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

	requestURLResolver := httputil.NewRequestURLResolver()

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

	notificationPublisher := buildNotificationPublisher(kafkaBrokers, kafkaNotificationTopic)
	orgHandler := initOrg.InitOrganization(requestURLResolver, starterRepo, departmentRepo, businessUnitRepo, notificationPublisher)
	orgHttp.RegisterOrganizationRoutes(v1, orgHandler)

	starterHandler := initOrg.InitStarter(esClient, starterRepo, departmentRepo, businessUnitRepo, requestURLResolver, kafkaBrokers, kafkaSyncTopic, kafkaConsumerGroup)
	orgHttp.RegisterStarterRoutes(v1, starterHandler)

	return r
}

func buildNotificationPublisher(kafkaBrokers, kafkaNotificationTopic string) messageBroker.NotificationPublisher {
	brokers := parseKafkaBrokers(kafkaBrokers)
	if len(brokers) == 0 || kafkaNotificationTopic == "" {
		log.Printf("Warning: Kafka notification configuration incomplete, leader notifications disabled")
		return nil
	}

	producer, err := sharedKafka.NewProducer(brokers)
	if err != nil {
		log.Printf("Warning: failed to initialize Kafka notification producer: %v", err)
		return nil
	}

	log.Printf("Kafka notification producer initialised for topic %s", kafkaNotificationTopic)
	return messageBroker.NewKafkaNotificationPublisher(producer, kafkaNotificationTopic)
}

func parseKafkaBrokers(raw string) []string {
	if raw == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	brokers := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			brokers = append(brokers, trimmed)
		}
	}
	return brokers
}
