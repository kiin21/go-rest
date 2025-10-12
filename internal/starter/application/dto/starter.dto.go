package dto

import "github.com/kiin21/go-rest/pkg/response"

// ============================================
// QUERIES (Read operations)
// ============================================

// ListStartersQuery represents parameters for listing starters
type ListStartersQuery struct {
	// Pagination
	Pagination response.ReqPagination

	// Filters
	BusinessUnitID *int64
	DepartmentID   *int64

	// Search
	Keyword string

	// Sorting
	SortBy    string
	SortOrder string
}

// SearchStartersQuery represents parameters for intelligent searching starters
// Priority: domain > name > business unit name
type SearchStartersQuery struct {
	// Search query string - intelligent search
	// Priority: domain > name > business unit name
	Query string

	// Filters
	DepartmentID   *int64
	BusinessUnitID *int64

	// Pagination
	Pagination response.ReqPagination
}

// ============================================
// COMMANDS (Write operations)
// ============================================

// CreateStarterCommand represents the command to create a new starter
type CreateStarterCommand struct {
	Domain        string
	Name          string // Required for creation
	Email         string
	Mobile        string
	WorkPhone     string
	JobTitle      string
	DepartmentID  *int64
	LineManagerID *int64
}

// UpdateStarterCommand represents the command to update a starter
// Use pointers to support partial updates (nil = don't update, value = update)
type UpdateStarterCommand struct {
	Domain        string
	Name          *string
	Email         *string
	Mobile        *string
	WorkPhone     *string
	JobTitle      *string
	DepartmentID  *int64
	LineManagerID *int64
}

// AssignToDepartmentCommand represents the command to assign a starter to a department
type AssignToDepartmentCommand struct {
	Domain       string
	DepartmentID int64
}

// RemoveStarterCommand represents the command to remove a starter
type RemoveStarterCommand struct {
	Domain       string
	NewManagerID *int64 // Optional: for delegating subordinates
}

// ============================================
// DTOs (Data Transfer Objects for responses)
// ============================================

// StarterDTO represents a starter for external consumption
type StarterDTO struct {
	ID            int64  `json:"id"`
	Domain        string `json:"domain"`
	Email         string `json:"email"`
	Mobile        string `json:"mobile"`
	WorkPhone     string `json:"work_phone"`
	JobTitle      string `json:"job_title"`
	DepartmentID  *int64 `json:"department_id"`
	LineManagerID *int64 `json:"line_manager_id"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}
