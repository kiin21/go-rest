package repository

import (
	"context"

	"github.com/kiin21/go-rest/pkg/httputil"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
)

type BusinessUnitRepository interface {
	FindByIDs(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error)
	List(ctx context.Context, pg httputil.ReqPagination) ([]*model.BusinessUnit, int64, error)
	FindByIDWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error)
	ListWithDetails(ctx context.Context, pg httputil.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error)
}
