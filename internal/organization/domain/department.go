package domain

import "time"

// Department represents the department aggregate in the organization domain.
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

// DepartmentWithDetails includes department with all related data
type DepartmentWithDetails struct {
	*Department
	BusinessUnit        *BusinessUnit
	Leader              *Leader
	ParentDepartment    *DepartmentNested
	MembersCount        int
	SubdepartmentsCount int
	Subdepartments      []*DepartmentNested
}

// DepartmentNested represents a simplified department (for parent/subdepartments)
type DepartmentNested struct {
	ID           int64
	FullName     string
	Shortname    string
	MembersCount int
}

// Leader represents a department leader (simplified starter)
type Leader struct {
	ID       int64
	Domain   string
	Name     string
	Email    string
	JobTitle string
}

// UpdateInfo updates the basic department information
func (d *Department) UpdateInfo(fullName, shortname *string) {
	if fullName != nil {
		d.FullName = *fullName
	}
	if shortname != nil {
		d.Shortname = *shortname
	}
}

// AssignToBusinessUnit assigns the department to a business unit
func (d *Department) AssignToBusinessUnit(businessUnitID *int64) {
	d.BusinessUnitID = businessUnitID
}

// AssignToGroupDepartment assigns the department to a parent department
func (d *Department) AssignToGroupDepartment(groupDepartmentID *int64) {
	d.GroupDepartmentID = groupDepartmentID
}

// AssignLeader assigns a leader to the department
func (d *Department) AssignLeader(leaderID *int64) {
	d.LeaderID = leaderID
}
