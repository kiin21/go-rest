package initialize

import (
	orgAppService "github.com/kiin21/go-rest/internal/organization/application/service"
	orgRepository "github.com/kiin21/go-rest/internal/organization/domain/repository"
	orgHttp "github.com/kiin21/go-rest/internal/organization/presentation/http"
	"github.com/kiin21/go-rest/pkg/httpctx"
)

func InitOrganization(
	requestURLResolver httpctx.RequestURLResolver,
	starterRepo orgRepository.StarterRepository,
	departmentRepo orgRepository.DepartmentRepository,
	businessUnitRepo orgRepository.BusinessUnitRepository,
) *orgHttp.OrganizationHandler {
	organizationService := orgAppService.NewOrganizationApplicationService(
		departmentRepo,
		businessUnitRepo,
		starterRepo,
	)

	return orgHttp.NewOrganizationHandler(organizationService, requestURLResolver)
}
