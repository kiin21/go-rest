package domain

import (
	"context"

	"github.com/kiin21/go-rest/pkg/response"
)

// DepartmentRepository describes persistence operations for the Department aggregate.
type DepartmentRepository interface {
	// List retrieves departments with optional filtering and pagination
	List(ctx context.Context, filter DepartmentListFilter, pg response.ReqPagination) ([]*Department, int64, error)

	// ListWithDetails retrieves departments with all related data (leader, business_unit, parent, subdepartments, counts)
	ListWithDetails(ctx context.Context, filter DepartmentListFilter, pg response.ReqPagination) ([]*DepartmentWithDetails, int64, error)

	// FindByID retrieves a department by its ID
	FindByID(ctx context.Context, id int64) (*Department, error)

	// FindByIDsWithRelations retrieves departments with group_department and business_unit in one query
	FindByIDsWithRelations(ctx context.Context, ids []int64) ([]*DepartmentWithRelations, error)

	// Create inserts a new department
	Create(ctx context.Context, department *Department) error

	// Update updates an existing department
	Update(ctx context.Context, department *Department) error
}

// DepartmentWithRelations holds department with its relations loaded
type DepartmentWithRelations struct {
	Department      *Department
	GroupDepartment *Department
	BusinessUnit    *BusinessUnit
}

// DepartmentListFilter contains filtering options for listing departments
type DepartmentListFilter struct {
	BusinessUnitID        *int64
	IncludeSubdepartments bool
}
