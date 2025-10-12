package initialize

import (
	orgApplication "github.com/kiin21/go-rest/internal/organization/application"
	orgDomain "github.com/kiin21/go-rest/internal/organization/domain"
	orgRepository "github.com/kiin21/go-rest/internal/organization/infrastructure/persistence/repository"
	orgHttp "github.com/kiin21/go-rest/internal/organization/interface/http"
	"gorm.io/gorm"
)

func InitOrganization(db *gorm.DB) (*orgHttp.OrganizationHandler, orgDomain.DepartmentRepository, orgDomain.BusinessUnitRepository) {
	departmentRepo := orgRepository.NewMySQLDepartmentRepository(db)
	businessUnitRepo := orgRepository.NewMySQLBusinessUnitRepository(db)
	organizationService := orgApplication.NewOrganizationApplicationService(departmentRepo)
	organizationHandler := orgHttp.NewOrganizationHandler(organizationService)

	return organizationHandler, departmentRepo, businessUnitRepo
}
