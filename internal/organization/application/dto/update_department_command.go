package dto

// UpdateDepartmentCommand represents the command to update a department
// Using pointers to distinguish between "field not sent" vs "field sent as null/empty"
type UpdateDepartmentCommand struct {
	FullName          *string
	Shortname         *string
	BusinessUnitID    *int64
	GroupDepartmentID *int64
	LeaderID          *int64
}
