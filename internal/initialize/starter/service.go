package initialize

import (
	"context"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
	"github.com/kiin21/go-rest/internal/shared/infrastructure"
	"github.com/kiin21/go-rest/internal/starter/application"
	starterDomain "github.com/kiin21/go-rest/internal/starter/domain"
	"github.com/kiin21/go-rest/internal/starter/infrastructure/persistence/repository"
	searchRepo "github.com/kiin21/go-rest/internal/starter/infrastructure/search/repository"
	starterHttp "github.com/kiin21/go-rest/internal/starter/interface/http"
	"gorm.io/gorm"
)

func InitStarter(
	db *gorm.DB,
	esClient *elasticsearch.Client,
	departmentRepo orgDomain.DepartmentRepository,
	businessUnitRepo orgDomain.BusinessUnitRepository,
	kafkaBrokers string,
	kafkaTopic string,
	kafkaConsumerGroup string,
) (*starterHttp.StarterHandler, *application.StarterSearchService, *infrastructure.KafkaConsumer) {
	starterRepo := repository.NewMySQLStarterRepository(db)

	// Initialize Kafka Producer only if config is provided
	var kafkaProducer *infrastructure.KafkaProducer
	var kafkaConsumer *infrastructure.KafkaConsumer

	if kafkaBrokers != "" && kafkaTopic != "" {
		brokers := strings.Split(kafkaBrokers, ",")
		var err error
		kafkaProducer, err = infrastructure.NewKafkaProducer(brokers, kafkaTopic)
		if err != nil {
			log.Printf("Warning: Failed to create Kafka producer: %v", err)
			log.Printf("Sync events will be disabled")
			kafkaProducer = nil
		} else {
			log.Printf("✅ Kafka producer connected successfully")
		}
	} else {
		log.Printf("Warning: Kafka configuration not provided, sync events disabled")
	}

	// Initialize Elasticsearch repository (for search) - only if esClient is available
	var starterSearchService *application.StarterSearchService
	var searchRepository starterDomain.StarterSearchRepository
	if esClient != nil {
		log.Println("🔍 Initializing Elasticsearch index...")
		// Ensure index exists before using Elasticsearch
		indexManager := searchRepo.NewIndexManager(esClient)
		if err := indexManager.CreateIndex(context.Background()); err != nil {
			log.Printf("Warning: Failed to create Elasticsearch index: %v", err)
			log.Printf("Elasticsearch search will be disabled")
		} else {
			log.Printf("✅ Elasticsearch index ready")
			searchRepository = searchRepo.NewElasticsearchStarterRepository(esClient)
			starterSearchService = application.NewStarterSearchService(searchRepository, starterRepo, kafkaProducer)

			// Check if index is empty and auto-reindex from MySQL
			log.Println("🔍 Checking if index is empty...")
			isEmpty, err := indexManager.IsIndexEmpty(context.Background())
			if err != nil {
				log.Printf("Warning: Failed to check if index is empty: %v", err)
			} else if isEmpty {
				log.Println("📦 Elasticsearch index is empty, starting auto-reindex from MySQL...")
				if err := starterSearchService.ReindexAll(context.Background()); err != nil {
					log.Printf("⚠️  Auto-reindex failed: %v", err)
				} else {
					log.Println("✅ Auto-reindex completed successfully")
				}
			} else {
				log.Println("ℹ️  Elasticsearch index already contains data, skipping auto-reindex")
			}
		}
	} else {
		log.Println("⚠️  Elasticsearch client is nil, skipping index initialization")
	}

	// Initialize Kafka Consumer with event handler only if producer exists
	if kafkaProducer != nil && searchRepository != nil {
		eventHandler := func(ctx context.Context, event *infrastructure.SyncEvent) error {
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
		kafkaConsumer, err = infrastructure.NewKafkaConsumer(brokers, kafkaConsumerGroup, kafkaTopic, eventHandler)
		if err != nil {
			log.Printf("Warning: Failed to create Kafka consumer: %v", err)
		} else {
			// Start Kafka consumer
			kafkaConsumer.Start()
			log.Printf("✅ Kafka consumer started successfully")
		}
	}

	// Create Domain Services
	// Domain Service encapsulates business rules that don't belong to a single entity
	starterDomainService := starterDomain.NewStarterDomainService(starterRepo)

	// Enrichment Service is a Domain Service that coordinates between multiple aggregates
	enrichmentService := starterDomain.NewStarterEnrichmentService(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
	)

	starterService := application.NewStarterApplicationService(
		starterRepo,
		starterDomainService,
		starterSearchService,
		enrichmentService,
	)

	// Initialize handler with services
	starterHandler := starterHttp.NewStarterHandler(starterService, starterSearchService)

	return starterHandler, starterSearchService, kafkaConsumer
}
