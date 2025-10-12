package domain

import (
	"context"

	"github.com/kiin21/go-rest/pkg/response"
)

type StarterSearchRepository interface {
	// Search operations
	Search(ctx context.Context, query string, filter ListFilter, pg response.ReqPagination) ([]*Starter, int64, error)

	// Index management
	IndexStarter(ctx context.Context, starter *Starter) error
	DeleteFromIndex(ctx context.Context, domain string) error

	// Bulk operations for initial indexing or reindexing
	BulkIndex(ctx context.Context, starters []*Starter) error
}
