package initialize

import (
	"context"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	orgApp "github.com/kiin21/go-rest/internal/organization/application/service"
	orgRepository "github.com/kiin21/go-rest/internal/organization/domain/repository"
	orgService "github.com/kiin21/go-rest/internal/organization/domain/service"
	orgInfraSearch "github.com/kiin21/go-rest/internal/organization/infrastructure/search/repository"
	orgHttp "github.com/kiin21/go-rest/internal/organization/presentation/http"
	messaging "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
	messagingKafka "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker/kafka"
	"github.com/kiin21/go-rest/pkg/httpctx"
)

func InitStarter(
	esClient *elasticsearch.Client,
	starterRepo orgRepository.StarterRepository,
	departmentRepo orgRepository.DepartmentRepository,
	businessUnitRepo orgRepository.BusinessUnitRepository,
	requestURLResolver httpctx.RequestURLResolver,
	kafkaBrokers string,
	kafkaTopic string,
	kafkaConsumerGroup string,
) *orgHttp.StarterHandler {
	// Initialize Kafka Producer
	var kafkaProducer *messagingKafka.Producer
	var kafkaConsumer *messagingKafka.Consumer

	if kafkaBrokers != "" && kafkaTopic != "" {
		brokers := strings.Split(kafkaBrokers, ",")
		var err error
		kafkaProducer, err = messagingKafka.NewProducer(brokers, kafkaTopic)
		if err != nil {
			log.Printf("Warning: failed to create Kafka producer: %v", err)
			log.Printf("Sync events will be disabled")
			kafkaProducer = nil
		} else {
			log.Printf("Kafka producer connected successfully")
		}
	} else {
		log.Printf("Warning: Kafka configuration not provided, sync events disabled")
	}

	// Initialize Elasticsearch repository
	var starterSearchService *orgApp.StarterSearchService
	var searchRepository orgRepository.StarterSearchRepository
	if esClient != nil {
		log.Println("Initializing Elasticsearch index...")
		indexManager := orgInfraSearch.NewIndexManager(esClient)
		if err := indexManager.CreateIndex(context.Background()); err != nil {
			log.Printf("Warning: failed to create Elasticsearch index: %v", err)
			log.Printf("Elasticsearch search will be disabled")
		} else {
			log.Printf("Elasticsearch index ready")
			searchRepository = orgInfraSearch.NewElasticsearchStarterRepository(esClient)
			starterSearchService = orgApp.NewStarterSearchService(searchRepository, starterRepo, kafkaProducer)

			log.Println("Checking if index is empty...")
			isEmpty, err := indexManager.IsIndexEmpty(context.Background())
			if err != nil {
				log.Printf("Warning: failed to check if index is empty: %v", err)
			} else if isEmpty {
				log.Println("Elasticsearch index is empty, starting auto-reindex from MySQL...")
				if err := starterSearchService.ReindexAll(context.Background()); err != nil {
					log.Printf("Auto-reindex failed: %v", err)
				} else {
					log.Println("Auto-reindex completed successfully")
				}
			} else {
				log.Println("Elasticsearch index already contains data, skipping auto-reindex")
			}
		}
	} else {
		log.Println("Elasticsearch client is nil, skipping index initialization")
	}

	// Initialize Kafka Consumer
	if kafkaProducer != nil && searchRepository != nil {
		eventHandler := func(ctx context.Context, event *messaging.SyncEvent) error {
			switch event.Type {
			case "index":
				starter, err := starterRepo.FindByDomain(ctx, event.Domain)
				if err != nil {
					return err
				}
				return searchRepository.IndexStarter(ctx, starter)
			case "delete":
				return searchRepository.DeleteFromIndex(ctx, event.Domain)
			}
			return nil
		}

		// Create Kafka consumer.
		brokers := strings.Split(kafkaBrokers, ",")
		var err error
		kafkaConsumer, err = messagingKafka.NewConsumer(brokers, kafkaConsumerGroup, kafkaTopic, eventHandler)
		if err != nil {
			log.Printf("Warning: failed to create Kafka consumer: %v", err)
		} else {
			// Start Kafka consumer.
			kafkaConsumer.Start()
			log.Printf("Kafka consumer started successfully")
		}
	}

	// Domain Services.
	starterDomainService := orgService.NewStarterDomainService(starterRepo)

	enrichmentService := orgService.NewStarterEnrichmentService(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
	)

	starterService := orgApp.NewStarterApplicationService(
		starterRepo,
		starterDomainService,
		enrichmentService,
		starterSearchService,
	)

	return orgHttp.NewStarterHandler(starterService, requestURLResolver)
}
