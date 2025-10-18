package repository

import (
	"context"

	model "github.com/kiin21/go-rest/internal/organization/domain/model"
	"github.com/kiin21/go-rest/pkg/response"
)

type StarterRepository interface {
	FindByID(ctx context.Context, id int64) (*model.Starter, error)
	FindByDomain(ctx context.Context, domain string) (*model.Starter, error)
	List(ctx context.Context, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error)
	SearchByKeyword(ctx context.Context, keyword string, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error)
	Create(ctx context.Context, starter *model.Starter) error
	Update(ctx context.Context, starter *model.Starter) error
	SoftDelete(ctx context.Context, domain string) error
}

type StarterSearchRepository interface {
	Search(ctx context.Context, query string, filter model.StarterListFilter, pg response.ReqPagination) ([]*model.Starter, int64, error)
	IndexStarter(ctx context.Context, starter *model.Starter) error
	DeleteFromIndex(ctx context.Context, domain string) error
	BulkIndex(ctx context.Context, starters []*model.Starter) error
}
