package initialize

import (
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	initDB "github.com/kiin21/go-rest/services/starter-service/internal/initialize/db"
	initES "github.com/kiin21/go-rest/services/starter-service/internal/initialize/elasticsearch"
	initBroker "github.com/kiin21/go-rest/services/starter-service/internal/initialize/messagebroker"
	initStarter "github.com/kiin21/go-rest/services/starter-service/internal/initialize/starter"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	persistentMySQL "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/persistence/repository/mysql"
)

func Run() (*gin.Engine, string, domainmessaging.NotificationProducer, domainmessaging.StarterConsumer) {
	// 1> Read config -> environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2> Initialize database connection
	db, err := initDB.InitMySQL(cfg.DBURI)
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	// 3> Initialize Elasticsearch client (optional)
	var esClient *elasticsearch.Client
	if cfg.ElasticsearchAddresses != "" {
		addresses := strings.Split(cfg.ElasticsearchAddresses, ",")
		esConfig := initES.Config{
			Addresses: addresses,
			Username:  cfg.ElasticsearchUsername,
			Password:  cfg.ElasticsearchPassword,
		}

		esClient, err = initES.InitESClient(esConfig)
		if err != nil {
			log.Printf("Warning: Could not initialize Elasticsearch: %v", err)
			esClient = nil
		}
	}

	// 4> Initialize Kafka notification producer
	notificationProducer := initBroker.InitProducer(cfg)

	// 5> Prepare shared dependencies
	requestURLResolver := httputil.NewRequestURLResolver()
	starterRepo := persistentMySQL.NewStarterRepository(db)
	businessUnitRepo := persistentMySQL.NewBusinessUnitRepository(db)
	departmentRepo := persistentMySQL.NewDepartmentRepository(db)

	orgHandler := initStarter.InitOrganization(
		requestURLResolver,
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		notificationProducer,
	)

	starterHandler, searchRepo := initStarter.InitStarter(
		esClient,
		starterRepo,
		departmentRepo,
		businessUnitRepo,
		requestURLResolver,
		notificationProducer,
	)

	eventHandler := initBroker.InitEventHandler(searchRepo)

	consumer := initBroker.InitGroupConsumer(cfg, eventHandler)

	// 6> Initialize router
	r := InitRouter(
		cfg.LogLevel,
		requestURLResolver,
		orgHandler,
		starterHandler,
	)

	return r, cfg.ServerPort, notificationProducer, consumer
}
