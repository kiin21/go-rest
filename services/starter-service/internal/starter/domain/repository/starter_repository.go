package repository

import (
	"context"

	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
)

type StarterRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Starter, error)
	FindByDomain(ctx context.Context, domain string) (*model.Starter, error)
	SearchByKeyword(ctx context.Context, listStarterQuery starterquery.ListStartersQuery) ([]*model.Starter, int64, error)
	Create(ctx context.Context, starter *model.Starter) error
	Update(ctx context.Context, starter *model.Starter) error
	SoftDelete(ctx context.Context, domain string) (*model.Starter, error)
}

type StarterSearchRepository interface {
	Search(ctx context.Context, listStarterQuery starterquery.ListStartersQuery) ([]*model.Starter, int64, error)
	IndexStarter(ctx context.Context, starter *model.Starter) error
	DeleteFromIndex(ctx context.Context, domain string) error
	BulkIndex(ctx context.Context, starters []*model.Starter) error
}
