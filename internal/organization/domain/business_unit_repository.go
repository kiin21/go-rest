package domain

import "context"

// BusinessUnitRepository defines operations for business unit persistence
type BusinessUnitRepository interface {
	FindByID(ctx context.Context, id int64) (*BusinessUnit, error)
	FindByIDs(ctx context.Context, ids []int64) ([]*BusinessUnit, error)
}
