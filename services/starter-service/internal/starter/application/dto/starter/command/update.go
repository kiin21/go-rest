package command

type UpdateStarterCommand struct {
	OriginalDomain string
	Domain         *string
	Name           *string
	Email          *string
	Mobile         *string
	WorkPhone      *string
	JobTitle       *string
	DepartmentID   *int64
	LineManagerID  *int64
}
