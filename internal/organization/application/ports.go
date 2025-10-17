package application

import (
	"context"

	"github.com/kiin21/go-rest/internal/starter/domain/aggregate"
)

type LeaderLookup interface {
	FindStarterIDByDomain(ctx context.Context, domain string) (int64, error)
	FindStarterById(ctx context.Context, id int64) (*aggregate.Starter, error)
}
