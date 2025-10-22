package initialize

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/kiin21/go-rest/pkg/httputil"
	orgApp "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
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
	notifProducer messaging.NotificationProducer,
) (*orgHttp.StarterHandler, orgRepository.StarterSearchRepository) {
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
			starterSearchService = orgApp.NewStarterSearchService(searchRepository, starterRepo, notifProducer)

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

	// Initialize Domain Services
	starterDomainService := orgService.NewStarterDomainService(starterRepo)

	enrichmentService := orgService.NewStarterEnrichmentService(
		starterRepo,
		departmentRepo,
		businessUnitRepo,
	)

	// Initialize Application Service
	starterService := orgApp.NewStarterApplicationService(
		starterRepo,
		starterDomainService,
		enrichmentService,
		starterSearchService,
	)

	return orgHttp.NewStarterHandler(starterService, requestURLResolver), searchRepository
}
