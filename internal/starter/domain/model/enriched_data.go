package model

import (
	sharedDomain "github.com/kiin21/go-rest/internal/shared/domain"
)

// EnrichedData holds enriched data for starters (departments, line managers, business units).
type EnrichedData struct {
	Departments   map[int64]*sharedDomain.DepartmentNested
	LineManagers  map[int64]*sharedDomain.LineManagerNested
	BusinessUnits map[int64]*sharedDomain.BusinessUnitNested
}

// DepartmentRelation contains department and related business unit data required for enrichment.
type DepartmentRelation struct {
	Department   *sharedDomain.DepartmentNested
	BusinessUnit *sharedDomain.BusinessUnitNested
}
