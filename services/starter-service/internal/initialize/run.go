package initialize

import (
	"log"
	"strings"

	"github.com/kiin21/go-rest/services/starter-service/internal/config"
	initDB "github.com/kiin21/go-rest/services/starter-service/internal/initialize/db"
	initES "github.com/kiin21/go-rest/services/starter-service/internal/initialize/elasticsearch"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
)

func Run() (*gin.Engine, string) {
	// 1> Read config -> environment variables
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// 2> Initialize database connection
	db, err := initDB.InitDB(&cfg)
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

	// 4> Initialize router with Kafka config
	r := InitRouter(
		db,
		esClient,
		cfg.LogLevel,
		cfg.KafkaBrokers,
		cfg.KafkaTopicSyncEvents,
		cfg.KafkaTopicNotifications,
		cfg.KafkaConsumerGroup,
	)

	return r, cfg.ServerPort
}
