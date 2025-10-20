package initialize

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/notification-service/internal/config"
	initDB "github.com/kiin21/go-rest/services/notification-service/internal/initialize/db"
	initmessagebroker "github.com/kiin21/go-rest/services/notification-service/internal/initialize/messagebroker"
	notiapp "github.com/kiin21/go-rest/services/notification-service/internal/notification/application"
	notinfra "github.com/kiin21/go-rest/services/notification-service/internal/notification/infrastructure/repository"
	notihttp "github.com/kiin21/go-rest/services/notification-service/internal/notification/presentation/http"
)

func Run() (*gin.Engine, string) {
	// 1> Read config -> environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2> Initialize database connection
	client, err := initDB.InitDB(&cfg)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	collectionName := cfg.MongoCollection

	collection := client.Collection(collectionName)
	repo := notinfra.NewNotificationMongoRepository(collection)
	service := notiapp.NewNotiApplicationService(repo)
	requestURLResolver := httputil.NewRequestURLResolver()
	handler := notihttp.NewNotiHandler(service, requestURLResolver)

	startNotificationConsumer(cfg, service)

	r := InitRouter(
		cfg.LogLevel,
		requestURLResolver,
		handler,
	)

	return r, cfg.ServerPort
}

func startNotificationConsumer(cfg config.Config, service *notiapp.NotiApplicationService) {
	consumer := initmessagebroker.InitNotificationConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaConsumerGroup,
		cfg.KafkaTopicNotifications,
		service,
	)
	if consumer == nil {
		return
	}

	go consumer.Start()
	log.Printf("Kafka notification consumer started for topic %s", cfg.KafkaTopicNotifications)
}
