package port

import (
	"context"

	"github.com/kiin21/go-rest/internal/starter/domain/aggregate"
	"github.com/kiin21/go-rest/internal/starter/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
)

// StarterRepository describes persistence operations for the Starter aggregate.
type StarterRepository interface {
	FindByID(ctx context.Context, id int64) (*aggregate.Starter, error)
	FindByDomain(ctx context.Context, domain string) (*aggregate.Starter, error)
	List(ctx context.Context, filter ListFilter, pg response.ReqPagination) ([]*aggregate.Starter, int64, error)
	SearchByKeyword(ctx context.Context, keyword string, filter ListFilter, pg response.ReqPagination) ([]*aggregate.Starter, int64, error)
	Create(ctx context.Context, starter *aggregate.Starter) error
	Update(ctx context.Context, starter *aggregate.Starter) error
	SoftDelete(ctx context.Context, domain string) error
}

// StarterSearchRepository encapsulates search/read model operations.
type StarterSearchRepository interface {
	Search(ctx context.Context, query string, filter ListFilter, pg response.ReqPagination) ([]*aggregate.Starter, int64, error)
	IndexStarter(ctx context.Context, starter *aggregate.Starter) error
	DeleteFromIndex(ctx context.Context, domain string) error
	BulkIndex(ctx context.Context, starters []*aggregate.Starter) error
}

// OrganizationLookup exposes the minimal data needed from the organization bounded context.
type OrganizationLookup interface {
	FindDepartmentRelations(ctx context.Context, ids []int64) ([]*model.DepartmentRelation, error)
}

// ListFilter contains filtering options for listing starters.
type ListFilter struct {
	DepartmentID   *int64
	BusinessUnitID *int64
	LineManagerID  *int64
}
