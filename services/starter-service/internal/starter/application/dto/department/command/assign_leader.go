package command

type AssignLeaderCommand struct {
	DepartmentID int64
	LeaderID     *int64
	LeaderDomain *string
}
