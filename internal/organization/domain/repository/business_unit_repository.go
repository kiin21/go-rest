package repository

import (
	"context"

	"github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
)

type BusinessUnitRepository interface {
	FindByIDs(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error)
	List(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnit, int64, error)
	FindByIDWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error)
	ListWithDetails(ctx context.Context, pg response.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error)
}
