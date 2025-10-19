package repository

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
)

type DepartmentRepository interface {
	ListWithDetails(ctx context.Context, filter model.DepartmentListFilter, pg response.ReqPagination) ([]*model.DepartmentWithDetails, int64, error)

	FindByIDsWithDetails(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)

	Create(ctx context.Context, department *model.Department) error

	Update(ctx context.Context, department *model.Department) error
}
