package domain

import (
	"context"

	"github.com/kiin21/go-rest/pkg/response"
)

type DepartmentRepository interface {
	ListWithDetails(ctx context.Context, filter DepartmentListFilter, pg response.ReqPagination) ([]*DepartmentWithDetails, int64, error)

	FindByIDsWithRelations(ctx context.Context, ids []int64) ([]*DepartmentWithDetails, error)

	Create(ctx context.Context, department *Department) error

	Update(ctx context.Context, department *Department) error
}

type DepartmentListFilter struct {
	BusinessUnitID *int64
}
