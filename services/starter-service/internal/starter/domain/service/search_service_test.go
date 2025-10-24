package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging/mocks"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository"
	repomocks "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository/mocks"
)

func TestMapSearchByToFieldName(t *testing.T) {
	service := &StarterSearchService{}

	tests := []struct {
		input    string
		expected string
	}{
		{"domain", "domain"},
		{"DOMAIN", "domain"},
		{"fullname", "name"},
		{"name", "name"},
		{"dept_name", "department_name"},
		{"bu_name", "business_unit_name"},
		{"", ""},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.mapSearchByToFieldName(tt.input)
			if result != tt.expected {
				t.Errorf("mapSearchByToFieldName(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestMapSortFieldToESField(t *testing.T) {
	service := &StarterSearchService{}

	tests := []struct {
		input    string
		expected string
	}{
		{"id", "id"},
		{"domain", "domain"},
		{"name", "name"},
		{"fullname", "name"},
		{"dept_name", "department_name"},
		{"bu_name", "business_unit_name"},
		{"", "id"},
		{"unknown", "id"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := service.mapSortFieldToESField(tt.input)
			if result != tt.expected {
				t.Errorf("mapSortFieldToESField(%s) = %s; want %s", tt.input, result, tt.expected)
			}
		})
	}
}

func TestBuildSortClause(t *testing.T) {
	service := &StarterSearchService{}

	tests := []struct {
		name      string
		sortBy    string
		sortOrder string
		hasResult bool
	}{
		{"ascending sort", "name", "asc", true},
		{"descending sort", "name", "desc", true},
		{"default asc when invalid order", "name", "invalid", true},
		{"empty sortBy defaults to id", "", "asc", true}, // mapSortFieldToESField returns "id" for empty
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.buildSortClause(tt.sortBy, tt.sortOrder)
			if tt.hasResult {
				if result == nil {
					t.Error("expected non-nil result")
				}
			} else {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			}
		})
	}
}

func TestBuildSearchQuery(t *testing.T) {
	service := &StarterSearchService{}

	tests := []struct {
		name     string
		query    *starterquery.ListStartersQuery
		isNil    bool
	}{
		{
			name:  "nil query",
			query: nil,
			isNil: true,
		},
		{
			name: "empty keyword",
			query: &starterquery.ListStartersQuery{
				Keyword: "",
			},
			isNil: true,
		},
		{
			name: "keyword with no searchBy",
			query: &starterquery.ListStartersQuery{
				Keyword: "test",
			},
			isNil: false,
		},
		{
			name: "keyword with domain searchBy",
			query: &starterquery.ListStartersQuery{
				Keyword:  "test",
				SearchBy: "domain",
			},
			isNil: false,
		},
		{
			name: "keyword with name searchBy",
			query: &starterquery.ListStartersQuery{
				Keyword:  "test",
				SearchBy: "name",
			},
			isNil: false,
		},
		{
			name: "with sort",
			query: &starterquery.ListStartersQuery{
				Keyword:   "test",
				SortBy:    "name",
				SortOrder: "desc",
			},
			isNil: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.buildSearchQuery(tt.query)
			if tt.isNil {
				if result != nil {
					t.Errorf("expected nil result, got %v", result)
				}
			} else {
				if result == nil {
					t.Error("expected non-nil result")
				}
			}
		})
	}
}

func TestIndexStarter(t *testing.T) {
	starter, _ := model.NewStarter(
		"testuser",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	t.Run("nil sync producer", func(t *testing.T) {
		service := &StarterSearchService{
			syncProducer: nil,
		}
		err := service.IndexStarter(context.Background(), starter)
		if err != nil {
			t.Errorf("expected no error with nil sync producer, got: %v", err)
		}
	})

	t.Run("successful index", func(t *testing.T) {
		mockProducer := &mocks.MockSyncProducer{
			SendSyncEventFunc: func(event *events.Event) error {
				return nil
			},
		}
		service := NewStarterSearchService(nil, nil, mockProducer)
		err := service.IndexStarter(context.Background(), starter)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("send event error", func(t *testing.T) {
		mockProducer := &mocks.MockSyncProducer{
			SendSyncEventFunc: func(event *events.Event) error {
				return errors.New("kafka error")
			},
		}
		service := NewStarterSearchService(nil, nil, mockProducer)
		err := service.IndexStarter(context.Background(), starter)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestDeleteFromIndex(t *testing.T) {
	starter, _ := model.NewStarter(
		"testuser",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	t.Run("nil sync producer", func(t *testing.T) {
		service := &StarterSearchService{
			syncProducer: nil,
		}
		err := service.DeleteFromIndex(context.Background(), starter)
		if err != nil {
			t.Errorf("expected no error with nil sync producer, got: %v", err)
		}
	})

	t.Run("successful delete", func(t *testing.T) {
		mockProducer := &mocks.MockSyncProducer{
			SendSyncEventFunc: func(event *events.Event) error {
				return nil
			},
		}
		service := NewStarterSearchService(nil, nil, mockProducer)
		err := service.DeleteFromIndex(context.Background(), starter)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("send event error", func(t *testing.T) {
		mockProducer := &mocks.MockSyncProducer{
			SendSyncEventFunc: func(event *events.Event) error {
				return errors.New("kafka error")
			},
		}
		service := NewStarterSearchService(nil, nil, mockProducer)
		err := service.DeleteFromIndex(context.Background(), starter)
		if err == nil {
			t.Error("expected error but got nil")
		}
	})
}

func TestSearch(t *testing.T) {
	page := 1
	limit := 10

	mockStarter, _ := model.Rehydrate(
		1,
		"testuser",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	tests := []struct {
		name             string
		query            *starterquery.ListStartersQuery
		mockSearchRepo   *repomocks.MockStarterSearchRepository
		mockStarterRepo  *repomocks.MockStarterRepository
		expectError      bool
	}{
		{
			name: "successful search",
			query: &starterquery.ListStartersQuery{
				Keyword: "test",
				Pagination: httputil.ReqPagination{
					Page:  &page,
					Limit: &limit,
				},
			},
			mockSearchRepo: &repomocks.MockStarterSearchRepository{
				SearchFunc: func(ctx context.Context, query *starterquery.ListStartersQuery, buildSearchQuery repository.SearchQueryBuilder) ([]int64, int64, error) {
					return []int64{1}, 1, nil
				},
			},
			mockStarterRepo: &repomocks.MockStarterRepository{
				FindByIDsFunc: func(ctx context.Context, ids []int64) ([]*model.Starter, error) {
					return []*model.Starter{mockStarter}, nil
				},
			},
			expectError: false,
		},
		{
			name: "search repository error",
			query: &starterquery.ListStartersQuery{
				Keyword: "test",
				Pagination: httputil.ReqPagination{
					Page:  &page,
					Limit: &limit,
				},
			},
			mockSearchRepo: &repomocks.MockStarterSearchRepository{
				SearchFunc: func(ctx context.Context, query *starterquery.ListStartersQuery, buildSearchQuery repository.SearchQueryBuilder) ([]int64, int64, error) {
					return nil, 0, errors.New("elasticsearch error")
				},
			},
			mockStarterRepo: &repomocks.MockStarterRepository{},
			expectError:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewStarterSearchService(
				tt.mockSearchRepo,
				tt.mockStarterRepo,
				nil,
			)

			result, err := service.Search(context.Background(), tt.query)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected non-nil result")
			}
		})
	}
}

