package dto

type AssignLeaderCommand struct {
	DepartmentID         int64
	LeaderID             *int64
	LeaderDomain         *string
	LeaderIdentifier     interface{}
	LeaderIdentifierType string // "id" or "domain"
}
