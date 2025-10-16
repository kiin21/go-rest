package domain

import (
	"context"

	"github.com/kiin21/go-rest/pkg/response"
)

type BusinessUnitRepository interface {
	FindByID(ctx context.Context, id int64) (*BusinessUnit, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*BusinessUnit, error)
	List(ctx context.Context, pg response.ReqPagination) ([]*BusinessUnit, int64, error)
	FindByIDWithDetails(ctx context.Context, id int64) (*BusinessUnitWithDetails, error)
	ListWithDetails(ctx context.Context, pg response.ReqPagination) ([]*BusinessUnitWithDetails, int64, error)
}
