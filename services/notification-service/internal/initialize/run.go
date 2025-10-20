package initialize

import (
	"context"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	"github.com/kiin21/go-rest/services/notification-service/internal/config"
	initDB "github.com/kiin21/go-rest/services/notification-service/internal/initialize/db"
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
	brokers := parseKafkaBrokers(cfg.KafkaBrokers)
	if len(brokers) == 0 || cfg.KafkaTopicNotifications == "" || cfg.KafkaConsumerGroup == "" {
		log.Printf("Warning: Kafka notification consumer disabled due to incomplete configuration")
		return
	}

	handler := func(ctx context.Context, event *events.Event) error {
		if event.Type != events.EventTypeNotificationLeaderAssignment {
			log.Printf("Warning: skipping unsupported notification event type %s", event.Type)
			return nil
		}

		var payload events.LeaderAssignmentNotification
		if err := event.DecodePayload(&payload); err != nil {
			return err
		}
		if payload.Timestamp.IsZero() {
			payload.Timestamp = event.Timestamp
		}

		return service.StoreNotification(ctx, &payload)
	}

	consumer, err := sharedKafka.NewEventConsumer(
		brokers,
		cfg.KafkaConsumerGroup,
		[]string{cfg.KafkaTopicNotifications},
		handler,
	)
	if err != nil {
		log.Printf("Warning: failed to initialize Kafka consumer: %v", err)
		return
	}

	go consumer.Start()
	log.Printf("Kafka notification consumer started for topic %s", cfg.KafkaTopicNotifications)
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
