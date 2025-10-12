package domain

import (
	"context"

	"github.com/kiin21/go-rest/pkg/response"
)

// StarterRepository describes persistence operations for the Starter aggregate.
type StarterRepository interface {
	// Business queries
	FindByID(ctx context.Context, id int64) (*Starter, error)
	FindByDomain(ctx context.Context, domain string) (*Starter, error)
	List(ctx context.Context, filter ListFilter, pg response.ReqPagination) ([]*Starter, int64, error)
	SearchByKeyword(ctx context.Context, keyword string, filter ListFilter, pg response.ReqPagination) ([]*Starter, int64, error)

	// Write operations
	Create(ctx context.Context, starter *Starter) error
	Update(ctx context.Context, starter *Starter) error
	SoftDelete(ctx context.Context, domain string) error
} // ListFilter contains filtering options for listing starters
type ListFilter struct {
	DepartmentID   *int64
	BusinessUnitID *int64
	LineManagerID  *int64
}
