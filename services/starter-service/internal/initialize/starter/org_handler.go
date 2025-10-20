package initialize

import (
	"github.com/kiin21/go-rest/pkg/httputil"
	orgAppService "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	orgRepository "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
)

func InitOrganization(
	requestURLResolver httputil.RequestURLResolver,
	starterRepo orgRepository.StarterRepository,
	departmentRepo orgRepository.DepartmentRepository,
	businessUnitRepo orgRepository.BusinessUnitRepository,
	notificationPublisher domainmessaging.NotificationPublisher,
) *orgHttp.OrganizationHandler {
	organizationService := orgAppService.NewOrganizationApplicationService(
		departmentRepo,
		businessUnitRepo,
		starterRepo,
		notificationPublisher,
	)

	return orgHttp.NewOrganizationHandler(organizationService, requestURLResolver)
}
