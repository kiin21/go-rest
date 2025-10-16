package domain

import (
	shareddomain "github.com/kiin21/go-rest/internal/shared/domain"
)

// EnrichedData holds enriched data for starters (departments, line managers, business units)
// This is a Domain model to avoid importing Interface DTOs
type EnrichedData struct {
	Departments   map[int64]*shareddomain.DepartmentNested
	LineManagers  map[int64]*shareddomain.LineManagerNested
	BusinessUnits map[int64]*shareddomain.BusinessUnitNested
}

// Type aliases for backward compatibility
type DepartmentNested = shareddomain.DepartmentNested
type GroupDepartmentNested = shareddomain.GroupDepartmentNested
type LineManagerNested = shareddomain.LineManagerNested
type BusinessUnitNested = shareddomain.BusinessUnitNested
