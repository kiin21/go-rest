package command

type UpdateDepartmentCommand struct {
	ID                int64
	FullName          *string
	Shortname         *string
	BusinessUnitID    *int64
	GroupDepartmentID *int64
	LeaderID          *int64
}
