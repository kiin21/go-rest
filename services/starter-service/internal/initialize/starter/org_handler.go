package initialize

import (
	orgAppSvc "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/service"
	domainmessaging "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging"
	orgRepo "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	orgHttp "github.com/kiin21/go-rest/services/starter-service/internal/starter/presentation/http"
)

func InitOrganization(
	starterRepo orgRepo.StarterRepository,
	deptRepo orgRepo.DepartmentRepository,
	buRepo orgRepo.BusinessUnitRepository,
	notifPublisher domainmessaging.NotificationProducer,
) *orgHttp.OrganizationHandler {
	organizationService := orgAppSvc.NewOrganizationApplicationService(deptRepo, buRepo, starterRepo, notifPublisher)

	return orgHttp.NewOrganizationHandler(organizationService)
}
