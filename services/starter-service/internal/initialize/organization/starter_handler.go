package initialize

import (
	"context"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	sharedKafka "github.com/kiin21/go-rest/pkg/kafka"
	orgApp "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	orgRepository "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	orgService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
	orgInfraSearch "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/search/repository"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
)

func InitStarter(
	esClient *elasticsearch.Client,
	starterRepo orgRepository.StarterRepository,
	departmentRepo orgRepository.DepartmentRepository,
	businessUnitRepo orgRepository.BusinessUnitRepository,
	requestURLResolver httputil.RequestURLResolver,
	kafkaBrokers string,
	kafkaSyncTopic string,
	kafkaConsumerGroup string,
) *orgHttp.StarterHandler {
	brokers := parseKafkaBrokers(kafkaBrokers)

	// Initialize Kafka Producer for search sync events
	var kafkaProducer *sharedKafka.Producer
	var kafkaConsumer *sharedKafka.EventConsumer
	if len(brokers) > 0 && kafkaSyncTopic != "" {
		var err error
		kafkaProducer, err = sharedKafka.NewProducerWithTopic(brokers, kafkaSyncTopic)
		if err != nil {
			log.Printf("Warning: failed to create Kafka producer for sync events: %v", err)
			log.Printf("Sync events will be disabled")
			kafkaProducer = nil
		} else {
			log.Printf("Kafka producer for sync events connected successfully")
		}
	} else {
		log.Printf("Warning: Kafka configuration for sync events not provided, sync events disabled")
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
	if kafkaProducer != nil && searchRepository != nil && kafkaConsumerGroup != "" && kafkaSyncTopic != "" && len(brokers) > 0 {
		eventHandler := func(ctx context.Context, event *events.Event) error {
			switch event.Type {
			case events.EventTypeStarterIndex, events.EventTypeStarterUpdate, events.EventTypeStarterInsert:
				starter, err := starterRepo.FindByDomain(ctx, event.Domain)
				if err != nil {
					return err
				}
				return searchRepository.IndexStarter(ctx, starter)
			case events.EventTypeStarterDelete:
				return searchRepository.DeleteFromIndex(ctx, event.Domain)
			}
			return nil
		}

		var err error
		kafkaConsumer, err = sharedKafka.NewEventConsumer(brokers, kafkaConsumerGroup, []string{kafkaSyncTopic}, eventHandler)
		if err != nil {
			log.Printf("Warning: failed to create Kafka consumer: %v", err)
		} else {
			go kafkaConsumer.Start()
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
