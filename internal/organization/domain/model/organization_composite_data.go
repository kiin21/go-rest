package model

import (
	sharedModel "github.com/kiin21/go-rest/internal/shared/domain/model"
)

type EnrichedData struct {
	Departments   map[int64]*sharedModel.DepartmentNested
	LineManagers  map[int64]*sharedModel.LineManagerNested
	BusinessUnits map[int64]*sharedModel.BusinessUnitNested
}

type DepartmentRelation struct {
	Department   *sharedModel.DepartmentNested
	BusinessUnit *sharedModel.BusinessUnitNested
}

type BusinessUnitWithDetails struct {
	*BusinessUnit
	Leader  *sharedModel.LineManagerNested
	Company *Company
}

type DepartmentWithDetails struct {
	*Department
	BusinessUnit     *BusinessUnit
	Leader           *sharedModel.LineManagerNested
	ParentDepartment *sharedModel.OrgDepartmentNested
	Subdepartments   []*sharedModel.OrgDepartmentNested
}
