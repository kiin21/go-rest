package domain

import (
	"time"

	shareddomain "github.com/kiin21/go-rest/internal/shared/domain"
)

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
	ParentDepartment    *shareddomain.OrgDepartmentNested
	Subdepartments      []*shareddomain.OrgDepartmentNested
	MembersCount        int
	SubdepartmentsCount int
}

// Type alias for backward compatibility
type DepartmentNested = shareddomain.OrgDepartmentNested

type Leader struct {
	ID       int64
	Domain   string
	Name     string
	Email    string
	JobTitle string
}
