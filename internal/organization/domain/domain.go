package domain

import (
	"time"

	shareddomain "github.com/kiin21/go-rest/internal/shared/domain"
)

type BusinessUnit struct {
	ID        int64
	Name      string
	Shortname string
	CompanyID int64
	LeaderID  *int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Company struct {
	ID   int64
	Name string
}

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

type BusinessUnitWithDetails struct {
	*BusinessUnit
	Leader  *LineManager
	Company *Company
}

type DepartmentWithDetails struct {
	*Department
	BusinessUnit     *BusinessUnit
	Leader           *LineManager
	ParentDepartment *shareddomain.OrgDepartmentNested
	Subdepartments   []*shareddomain.OrgDepartmentNested
	MembersCount     int
}

type DepartmentNested = shareddomain.OrgDepartmentNested

type LineManager struct {
	ID       int64
	Domain   string
	Name     string
	Email    string
	JobTitle string
}
