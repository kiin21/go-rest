package initialize

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/notification-service/internal/config"
	initDB "github.com/kiin21/go-rest/services/notification-service/internal/initialize/db"
	initmessagebroker "github.com/kiin21/go-rest/services/notification-service/internal/initialize/event_handler"
	notiapp "github.com/kiin21/go-rest/services/notification-service/internal/notification/application"
	domainmessaging "github.com/kiin21/go-rest/services/notification-service/internal/notification/domain/messaging"
	notinfra "github.com/kiin21/go-rest/services/notification-service/internal/notification/infrastructure/repository"
	notihttp "github.com/kiin21/go-rest/services/notification-service/internal/notification/presentation/http"
)

func Run() (*gin.Engine, string, domainmessaging.NotificationConsumer) {
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

	// 3> Setup repository, service, handler
	collection := client.Collection(collectionName)
	repo := notinfra.NewNotificationMongoRepository(collection)
	service := notiapp.NewNotiApplicationService(repo)
	requestURLResolver := httputil.NewRequestURLResolver()
	handler := notihttp.NewNotiHandler(service, requestURLResolver)

	// 4> Initialize Kafka consumer
	eventHandler := initmessagebroker.InitEventHandler(repo)
	consumer := initmessagebroker.InitGroupConsumer(cfg, eventHandler)

	// 5> Setup router
	r := InitRouter(cfg.LogLevel, requestURLResolver, handler)

	return r, cfg.ServerPort, consumer
}
