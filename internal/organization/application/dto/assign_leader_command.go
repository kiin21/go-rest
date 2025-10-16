package dto

// AssignLeaderCommand represents the command for assigning a leader to a department
type AssignLeaderCommand struct {
	DepartmentID       int64
	LeaderID           *int64  // Used when assigning by ID
	LeaderDomain       *string // Used when assigning by domain
	LeaderIdentifier   interface{}
	LeaderIdentifierType string // "id" or "domain"
}
