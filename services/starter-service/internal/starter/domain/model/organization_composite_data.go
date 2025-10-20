package model

type EnrichedData struct {
	Departments   map[int64]*DepartmentNested
	LineManagers  map[int64]*LineManagerNested
	BusinessUnits map[int64]*BusinessUnitNested
}

type DepartmentRelation struct {
	Department   *DepartmentNested
	BusinessUnit *BusinessUnitNested
}

type BusinessUnitWithDetails struct {
	*BusinessUnit
	Leader  *LineManagerNested
	Company *Company
}

type DepartmentWithDetails struct {
	*Department
	BusinessUnit     *BusinessUnit
	Leader           *LineManagerNested
	ParentDepartment *OrgDepartmentNested
	Subdepartments   []*OrgDepartmentNested
}
