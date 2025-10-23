package initialize

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	starterApp "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	starterDomainRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	starterDomainSvc "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
	starterInfraSearch "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/search/repository"
	starterHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
)

func InitStarter(
	starterRepo starterDomainRepo.StarterRepository,
	departmentRepo starterDomainRepo.DepartmentRepository,
	businessUnitRepo starterDomainRepo.BusinessUnitRepository,
	esClient *elasticsearch.Client,
	syncProducer messaging.SyncProducer,
) (*starterHttp.StarterHandler, starterDomainRepo.StarterSearchRepository, *starterDomainSvc.StarterEnrichmentService) {
	var (
		starterSearchRepo    starterDomainRepo.StarterSearchRepository
		starterSearchService *starterDomainSvc.StarterSearchService
	)

	// Initialize Elasticsearch repository
	if esClient != nil {
		log.Println("Initializing Elasticsearch index...")

		indexManager, err := starterInfraSearch.NewIndexManager(esClient)
		if err != nil {
			log.Printf("Warning: failed to create index manager: %v", err)
			log.Printf("Elasticsearch search will be disabled")
		} else {
			if err := indexManager.CreateIndex(context.Background()); err != nil {
				log.Printf("Warning: failed to create Elasticsearch index: %v", err)
				log.Printf("Elasticsearch search will be disabled")
			} else {
				log.Printf("Elasticsearch index ready")

				starterSearchRepo = starterInfraSearch.NewElasticsearchStarterRepository(esClient)
				starterSearchService = starterDomainSvc.NewStarterSearchService(
					starterSearchRepo,
					starterRepo,
					syncProducer,
				)

				log.Println("Checking if index is empty...")
				isEmpty, err := indexManager.IsIndexEmpty(context.Background())
				if err != nil {
					log.Printf("Warning: failed to check if index is empty: %v", err)
				} else if isEmpty {
					log.Println("Elasticsearch index is empty, starting auto-reindex from MySQL...")
					// Auto-reindex will be triggered below
				} else {
					log.Println("Elasticsearch index already contains data, skipping auto-reindex")
				}
			}
		}
	} else {
		log.Println("Elasticsearch client is nil, skipping index initialization")
	}

	// Initialize Domain Services
	starterDomainService := starterDomainSvc.NewStarterDomainService(starterRepo)

	starterEnrichmentService := starterDomainSvc.NewStarterEnrichmentService(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
	)

	// Initialize Application Service
	starterAppService := starterApp.NewStarterApplicationService(
		starterRepo,
		starterSearchRepo,
		starterDomainService,
		starterEnrichmentService,
		starterSearchService,
	)

	// Auto-reindex on startup if ES is enabled
	if starterSearchRepo != nil {
		log.Println("Starting auto-reindex...")
		if err := starterAppService.ReindexAll(context.Background()); err != nil {
			log.Printf("Auto-reindex failed: %v", err)
		} else {
			log.Println("Auto-reindex completed successfully")
		}
	}

	starterHandler := starterHttp.NewStarterHandler(starterAppService, starterEnrichmentService)

	return starterHandler, starterSearchRepo, starterEnrichmentService
}
