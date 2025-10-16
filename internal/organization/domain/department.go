package domain

import "time"

type Department struct {
	ID                int64
	GroupDepartmentID *int64
	FullName          string
	Shortname         string
	BusinessUnitID    *int64
	LeaderID          *int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         *time.Time
}

type DepartmentWithDetails struct {
	*Department
	BusinessUnit        *BusinessUnit
	Leader              *Leader
	ParentDepartment    *DepartmentNested
	Subdepartments      []*DepartmentNested
	MembersCount        int
	SubdepartmentsCount int
}

type DepartmentNested struct {
	ID           int64
	FullName     string
	Shortname    string
	MembersCount int
}

type Leader struct {
	ID       int64
	Domain   string
	Name     string
	Email    string
	JobTitle string
}

