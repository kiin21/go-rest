package initialize

import (
	"github.com/kiin21/go-rest/pkg/httputil"
	orgAppService "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	orgRepository "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	messagebroker "github.com/kiin21/go-rest/services/starter-service/internal/starter/infrastructure/messagebroker"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
)

func InitOrganization(
	requestURLResolver httputil.RequestURLResolver,
	starterRepo orgRepository.StarterRepository,
	departmentRepo orgRepository.DepartmentRepository,
	businessUnitRepo orgRepository.BusinessUnitRepository,
	notificationPublisher messagebroker.NotificationPublisher,
) *orgHttp.OrganizationHandler {
	organizationService := orgAppService.NewOrganizationApplicationService(
		departmentRepo,
		businessUnitRepo,
		starterRepo,
		notificationPublisher,
	)

	return orgHttp.NewOrganizationHandler(organizationService, requestURLResolver)
}
