package repository

import (
	"context"

	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
)

type StarterRepository interface {
	FindByIDs(ctx context.Context, ids []int64) ([]*model.Starter, error)
	FindByDomain(ctx context.Context, domain string) (*model.Starter, error)
	SearchByKeyword(ctx context.Context, listStarterQuery *starterquery.ListStartersQuery) ([]*model.Starter, int64, error)
	Create(ctx context.Context, starter *model.Starter) error
	Update(ctx context.Context, starter *model.Starter) error
	SoftDelete(ctx context.Context, domain string) (*model.Starter, error)
}

type SearchQueryBuilder func(*starterquery.ListStartersQuery) map[string]interface{}

type StarterSearchRepository interface {
	Search(ctx context.Context, listStarterQuery *starterquery.ListStartersQuery, buildSearchQuery SearchQueryBuilder) ([]int64, int64, error)
	IndexStarter(ctx context.Context, starter *model.StarterESDoc) error
	DeleteFromIndex(ctx context.Context, domain string) error
	BulkIndex(ctx context.Context, starters []*model.StarterESDoc) error
}
