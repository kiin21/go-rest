package initialize

import (
	"context"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kiin21/go-rest/internal/composition"
	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
	messaging "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker"
	messagingKafka "github.com/kiin21/go-rest/internal/shared/infrastructure/messagebroker/kafka"
	"github.com/kiin21/go-rest/internal/starter/application"
	starterPort "github.com/kiin21/go-rest/internal/starter/domain/port"
	starterService "github.com/kiin21/go-rest/internal/starter/domain/service"
	searchRepo "github.com/kiin21/go-rest/internal/starter/infrastructure/search/repository"
	starterHttp "github.com/kiin21/go-rest/internal/starter/presentation/http"
	"github.com/kiin21/go-rest/pkg/httpctx"
)

func InitStarter(
	esClient *elasticsearch.Client,
	starterRepo starterPort.StarterRepository,
	departmentRepo orgDomain.DepartmentRepository,
	requestURLResolver httpctx.RequestURLResolver,
	kafkaBrokers string,
	kafkaTopic string,
	kafkaConsumerGroup string,
) (*starterHttp.StarterHandler, *messagingKafka.Consumer) {

	// Initialize Kafka Producer only if config is provided
	var kafkaProducer *messagingKafka.Producer
	var kafkaConsumer *messagingKafka.Consumer

	if kafkaBrokers != "" && kafkaTopic != "" {
		brokers := strings.Split(kafkaBrokers, ",")
		var err error
		kafkaProducer, err = messagingKafka.NewProducer(brokers, kafkaTopic)
		if err != nil {
			log.Printf("Warning: Failed to create Kafka producer: %v", err)
			log.Printf("Sync events will be disabled")
			kafkaProducer = nil
		} else {
			log.Printf("Kafka producer connected successfully")
		}
	} else {
		log.Printf("Warning: Kafka configuration not provided, sync events disabled")
	}

	// Initialize Elasticsearch repository (for search) - only if esClient is available
	var starterSearchService *application.StarterSearchService
	var searchRepository starterPort.StarterSearchRepository
	if esClient != nil {
		log.Println("🔍 Initializing Elasticsearch index...")
		// Ensure index exists before using Elasticsearch
		indexManager := searchRepo.NewIndexManager(esClient)
		if err := indexManager.CreateIndex(context.Background()); err != nil {
			log.Printf("Warning: Failed to create Elasticsearch index: %v", err)
			log.Printf("Elasticsearch search will be disabled")
		} else {
			log.Printf("Elasticsearch index ready")
			searchRepository = searchRepo.NewElasticsearchStarterRepository(esClient)
			starterSearchService = application.NewStarterSearchService(searchRepository, starterRepo, kafkaProducer)

			// Check if index is empty and auto-reindex from MySQL
			log.Println("Checking if index is empty...")
			isEmpty, err := indexManager.IsIndexEmpty(context.Background())
			if err != nil {
				log.Printf("Warning: Failed to check if index is empty: %v", err)
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

	// Initialize Kafka Consumer with event handler only if producer exists
	if kafkaProducer != nil && searchRepository != nil {
		eventHandler := func(ctx context.Context, event *messaging.SyncEvent) error {
			switch event.Type {
			case "index":
				// Get starter from database using domain
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

		brokers := strings.Split(kafkaBrokers, ",")
		var err error
		kafkaConsumer, err = messagingKafka.NewConsumer(brokers, kafkaConsumerGroup, kafkaTopic, eventHandler)
		if err != nil {
			log.Printf("Warning: Failed to create Kafka consumer: %v", err)
		} else {
			// Start Kafka consumer
			kafkaConsumer.Start()
			log.Printf("Kafka consumer started successfully")
		}
	}

	// omain Services
	starterDomainService := starterService.NewStarterDomainService(starterRepo)

	// Service that coordinates between multiple aggregates
	organizationLookup := composition.NewOrganizationLookup(departmentRepo)
	enrichmentService := starterService.NewStarterEnrichmentService(
		starterRepo,
		organizationLookup,
	)

	_starterService := application.NewStarterApplicationService(
		starterRepo,
		starterDomainService,
		starterSearchService,
		enrichmentService,
	)

	starterHandler := starterHttp.NewStarterHandler(_starterService, starterSearchService, requestURLResolver)

	return starterHandler, kafkaConsumer
}
