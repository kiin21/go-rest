package mocks

import (
	"context"

	"github.com/kiin21/go-rest/pkg/httputil"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
)

// MockStarterRepository is a mock implementation of StarterRepository
type MockStarterRepository struct {
	CreateFunc           func(ctx context.Context, starter *model.Starter) error
	UpdateFunc           func(ctx context.Context, starter *model.Starter) error
	SoftDeleteFunc       func(ctx context.Context, domain string) (*model.Starter, error)
	FindByDomainFunc     func(ctx context.Context, domain string) (*model.Starter, error)
	FindByIDsFunc        func(ctx context.Context, ids []int64) ([]*model.Starter, error)
	SearchByKeywordFunc  func(ctx context.Context, query *starterquery.ListStartersQuery) ([]*model.Starter, int64, error)
}

func (m *MockStarterRepository) Create(ctx context.Context, starter *model.Starter) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, starter)
	}
	return nil
}

func (m *MockStarterRepository) Update(ctx context.Context, starter *model.Starter) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, starter)
	}
	return nil
}

func (m *MockStarterRepository) SoftDelete(ctx context.Context, domain string) (*model.Starter, error) {
	if m.SoftDeleteFunc != nil {
		return m.SoftDeleteFunc(ctx, domain)
	}
	return nil, nil
}

func (m *MockStarterRepository) FindByDomain(ctx context.Context, domain string) (*model.Starter, error) {
	if m.FindByDomainFunc != nil {
		return m.FindByDomainFunc(ctx, domain)
	}
	return nil, nil
}

func (m *MockStarterRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Starter, error) {
	if m.FindByIDsFunc != nil {
		return m.FindByIDsFunc(ctx, ids)
	}
	return nil, nil
}

func (m *MockStarterRepository) SearchByKeyword(ctx context.Context, query *starterquery.ListStartersQuery) ([]*model.Starter, int64, error) {
	if m.SearchByKeywordFunc != nil {
		return m.SearchByKeywordFunc(ctx, query)
	}
	return nil, 0, nil
}

// MockStarterSearchRepository is a mock implementation of StarterSearchRepository
type MockStarterSearchRepository struct {
	SearchFunc           func(ctx context.Context, query *starterquery.ListStartersQuery, buildSearchQuery repository.SearchQueryBuilder) ([]int64, int64, error)
	IndexStarterFunc     func(ctx context.Context, doc *model.StarterESDoc) error
	BulkIndexFunc        func(ctx context.Context, docs []*model.StarterESDoc) error
	DeleteFromIndexFunc  func(ctx context.Context, domain string) error
}

func (m *MockStarterSearchRepository) Search(ctx context.Context, query *starterquery.ListStartersQuery, buildSearchQuery repository.SearchQueryBuilder) ([]int64, int64, error) {
	if m.SearchFunc != nil {
		return m.SearchFunc(ctx, query, buildSearchQuery)
	}
	return nil, 0, nil
}

func (m *MockStarterSearchRepository) IndexStarter(ctx context.Context, doc *model.StarterESDoc) error {
	if m.IndexStarterFunc != nil {
		return m.IndexStarterFunc(ctx, doc)
	}
	return nil
}

func (m *MockStarterSearchRepository) BulkIndex(ctx context.Context, docs []*model.StarterESDoc) error {
	if m.BulkIndexFunc != nil {
		return m.BulkIndexFunc(ctx, docs)
	}
	return nil
}

func (m *MockStarterSearchRepository) DeleteFromIndex(ctx context.Context, domain string) error {
	if m.DeleteFromIndexFunc != nil {
		return m.DeleteFromIndexFunc(ctx, domain)
	}
	return nil
}

// MockDepartmentRepository is a mock implementation of DepartmentRepository
type MockDepartmentRepository struct {
	CreateFunc                func(ctx context.Context, department *model.Department) error
	UpdateFunc                func(ctx context.Context, department *model.Department) error
	DeleteFunc                func(ctx context.Context, id int64) error
	FindByIDsFunc             func(ctx context.Context, ids []int64) ([]*model.Department, error)
	FindByIDsWithDetailsFunc  func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
	ListWithDetailsFunc       func(ctx context.Context, filter *model.DepartmentListFilter, pagination *httputil.ReqPagination) ([]*model.DepartmentWithDetails, int64, error)
	SearchByKeywordFunc       func(ctx context.Context, keyword string) ([]*model.Department, int64, error)
}

func (m *MockDepartmentRepository) Create(ctx context.Context, department *model.Department) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, department)
	}
	return nil
}

func (m *MockDepartmentRepository) Update(ctx context.Context, department *model.Department) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, department)
	}
	return nil
}

func (m *MockDepartmentRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return nil
}

func (m *MockDepartmentRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.Department, error) {
	if m.FindByIDsFunc != nil {
		return m.FindByIDsFunc(ctx, ids)
	}
	return nil, nil
}

func (m *MockDepartmentRepository) FindByIDsWithDetails(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
	if m.FindByIDsWithDetailsFunc != nil {
		return m.FindByIDsWithDetailsFunc(ctx, ids)
	}
	return nil, nil
}

func (m *MockDepartmentRepository) ListWithDetails(ctx context.Context, filter *model.DepartmentListFilter, pagination *httputil.ReqPagination) ([]*model.DepartmentWithDetails, int64, error) {
	if m.ListWithDetailsFunc != nil {
		return m.ListWithDetailsFunc(ctx, filter, pagination)
	}
	return nil, 0, nil
}

func (m *MockDepartmentRepository) SearchByKeyword(ctx context.Context, keyword string) ([]*model.Department, int64, error) {
	if m.SearchByKeywordFunc != nil {
		return m.SearchByKeywordFunc(ctx, keyword)
	}
	return nil, 0, nil
}

// MockBusinessUnitRepository is a mock implementation of BusinessUnitRepository
type MockBusinessUnitRepository struct {
	FindByIDsFunc            func(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error)
	FindByIDWithDetailsFunc  func(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error)
	ListFunc                 func(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnit, int64, error)
	ListWithDetailsFunc      func(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error)
}

func (m *MockBusinessUnitRepository) FindByIDs(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
	if m.FindByIDsFunc != nil {
		return m.FindByIDsFunc(ctx, ids)
	}
	return nil, nil
}

func (m *MockBusinessUnitRepository) FindByIDWithDetails(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
	if m.FindByIDWithDetailsFunc != nil {
		return m.FindByIDWithDetailsFunc(ctx, id)
	}
	return nil, nil
}

func (m *MockBusinessUnitRepository) List(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnit, int64, error) {
	if m.ListFunc != nil {
		return m.ListFunc(ctx, pagination)
	}
	return nil, 0, nil
}

func (m *MockBusinessUnitRepository) ListWithDetails(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error) {
	if m.ListWithDetailsFunc != nil {
		return m.ListWithDetailsFunc(ctx, pagination)
	}
	return nil, 0, nil
}

