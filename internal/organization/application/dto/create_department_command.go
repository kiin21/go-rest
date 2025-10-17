package dto

type CreateDepartmentCommand struct {
	FullName          string
	Shortname         string
	BusinessUnitID    *int64
	GroupDepartmentID *int64
	LeaderID          *int64
}
