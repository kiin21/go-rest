package initialize

import (
	"github.com/kiin21/go-rest/internal/composition"
	orgApplication "github.com/kiin21/go-rest/internal/organization/application"
	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
	orgRepository "github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/repository"
	orgHttp "github.com/kiin21/go-rest/internal/organization/presentation/http"
	"github.com/kiin21/go-rest/internal/starter/domain/port"
	"github.com/kiin21/go-rest/pkg/httpctx"
	"gorm.io/gorm"
)

func InitOrganization(
	db *gorm.DB,
	requestURLResolver httpctx.RequestURLResolver,
	starterRepo port.StarterRepository,
) (*orgHttp.OrganizationHandler, orgDomain.DepartmentRepository) {
	departmentRepo := orgRepository.NewMySQLDepartmentRepository(db)
	businessUnitRepo := orgRepository.NewMySQLBusinessUnitRepository(db)

	leaderLookup := composition.NewStarterLeaderLookup(starterRepo)

	organizationService := orgApplication.NewOrganizationApplicationService(departmentRepo, businessUnitRepo, leaderLookup)
	organizationHandler := orgHttp.NewOrganizationHandler(organizationService, requestURLResolver)

	return organizationHandler, departmentRepo
}
