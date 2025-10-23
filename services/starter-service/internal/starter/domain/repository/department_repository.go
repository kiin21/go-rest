package repository

import (
	"context"

	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
)

type DepartmentRepository interface {
	ListWithDetails(ctx context.Context, filter *model.DepartmentListFilter, pg *httputil.ReqPagination) ([]*model.DepartmentWithDetails, int64, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Department, error)
	SearchByKeyword(ctx context.Context, keyword string) ([]*model.Department, int64, error)
	FindByIDsWithDetails(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
	Create(ctx context.Context, department *model.Department) error
	Update(ctx context.Context, department *model.Department) error
	Delete(ctx context.Context, id int64) error
}
