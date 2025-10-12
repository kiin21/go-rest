package dto

// CreateDepartmentCommand represents the command to create a new department
type CreateDepartmentCommand struct {
	FullName          string
	Shortname         string
	BusinessUnitID    *int64
	GroupDepartmentID *int64
	LeaderID          *int64
}
