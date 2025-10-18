package command

type CreateStarterCommand struct {
	Domain        string
	Name          string
	Email         string
	Mobile        string
	WorkPhone     string
	JobTitle      string
	DepartmentID  *int64
	LineManagerID *int64
}
