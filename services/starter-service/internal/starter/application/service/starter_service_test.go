package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kiin21/go-rest/pkg/httputil"
	startercommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/command"
	starterquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/starter/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository/mocks"
	domainService "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/service"
)

func TestCreateStarter(t *testing.T) {
	tests := []struct {
		name                string
		command             *startercommand.CreateStarterCommand
		mockFindByDomain    func(ctx context.Context, domain string) (*model.Starter, error)
		mockCreate          func(ctx context.Context, starter *model.Starter) error
		expectError         bool
		expectErrorContains string
	}{
		{
			name: "successful creation",
			command: &startercommand.CreateStarterCommand{
				Domain:   "newuser",
				Name:     "New User",
				Email:    "newuser@vng.com.vn",
				Mobile:   "0123456789",
				JobTitle: "Developer",
			},
			mockFindByDomain: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("resource not found")
			},
			mockCreate: func(ctx context.Context, starter *model.Starter) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "domain already exists",
			command: &startercommand.CreateStarterCommand{
				Domain:   "existinguser",
				Name:     "Existing User",
				Email:    "existing@vng.com.vn",
				Mobile:   "0123456789",
				JobTitle: "Developer",
			},
			mockFindByDomain: func(ctx context.Context, domain string) (*model.Starter, error) {
				starter, _ := model.NewStarter(
					"existinguser",
					"Existing User",
					"existing@vng.com.vn",
					"0123456789",
					"",
					"Developer",
					nil,
					nil,
				)
				return starter, nil
			},
			mockCreate:          func(ctx context.Context, starter *model.Starter) error { return nil },
			expectError:         true,
			expectErrorContains: "domain already exists",
		},
		{
			name: "repository create error",
			command: &startercommand.CreateStarterCommand{
				Domain:   "newuser",
				Name:     "New User",
				Email:    "newuser@vng.com.vn",
				Mobile:   "0123456789",
				JobTitle: "Developer",
			},
			mockFindByDomain: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("resource not found")
			},
			mockCreate: func(ctx context.Context, starter *model.Starter) error {
				return errors.New("database error")
			},
			expectError:         true,
			expectErrorContains: "database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStarterRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFindByDomain,
				CreateFunc:       tt.mockCreate,
			}

			mockSearchRepo := &mocks.MockStarterSearchRepository{}
			domainSvc := domainService.NewStarterDomainService(mockStarterRepo)

			service := NewStarterApplicationService(
				mockStarterRepo,
				mockSearchRepo,
				domainSvc,
				nil,
				nil,
			)

			starter, err := service.CreateStarter(context.Background(), tt.command)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
					return
				}
				if tt.expectErrorContains != "" && err.Error() != tt.expectErrorContains {
					t.Errorf("expected error containing '%s', got '%s'", tt.expectErrorContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if starter == nil {
				t.Error("expected non-nil starter")
				return
			}

			if starter.Domain != tt.command.Domain {
				t.Errorf("expected domain %s, got %s", tt.command.Domain, starter.Domain)
			}
		})
	}
}

func TestGetStarterByDomain(t *testing.T) {
	mockStarter, _ := model.NewStarter(
		"testuser",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	tests := []struct {
		name             string
		domain           string
		mockFindByDomain func(ctx context.Context, domain string) (*model.Starter, error)
		expectError      bool
	}{
		{
			name:   "starter found",
			domain: "testuser",
			mockFindByDomain: func(ctx context.Context, domain string) (*model.Starter, error) {
				return mockStarter, nil
			},
			expectError: false,
		},
		{
			name:   "starter not found",
			domain: "nonexistent",
			mockFindByDomain: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, sharedDomain.ErrNotFound
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStarterRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFindByDomain,
			}

			service := NewStarterApplicationService(
				mockStarterRepo,
				nil,
				nil,
				nil,
				nil,
			)

			starter, err := service.GetStarterByDomain(context.Background(), tt.domain)

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

			if starter == nil {
				t.Error("expected non-nil starter")
			}
		})
	}
}

func TestUpdateStarter(t *testing.T) {
	existingStarter, _ := model.Rehydrate(
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

	newName := "Updated User"
	newEmail := "updated@vng.com.vn"

	tests := []struct {
		name        string
		command     *startercommand.UpdateStarterCommand
		mockFind    func(ctx context.Context, domain string) (*model.Starter, error)
		mockUpdate  func(ctx context.Context, starter *model.Starter) error
		expectError bool
	}{
		{
			name: "successful update",
			command: &startercommand.UpdateStarterCommand{
				Domain: "testuser",
				Name:   &newName,
				Email:  &newEmail,
			},
			mockFind: func(ctx context.Context, domain string) (*model.Starter, error) {
				return existingStarter, nil
			},
			mockUpdate: func(ctx context.Context, starter *model.Starter) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "starter not found",
			command: &startercommand.UpdateStarterCommand{
				Domain: "nonexistent",
				Name:   &newName,
			},
			mockFind: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, sharedDomain.ErrNotFound
			},
			mockUpdate:  func(ctx context.Context, starter *model.Starter) error { return nil },
			expectError: true,
		},
		{
			name: "update error",
			command: &startercommand.UpdateStarterCommand{
				Domain: "testuser",
				Name:   &newName,
			},
			mockFind: func(ctx context.Context, domain string) (*model.Starter, error) {
				return existingStarter, nil
			},
			mockUpdate: func(ctx context.Context, starter *model.Starter) error {
				return errors.New("database error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStarterRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFind,
				UpdateFunc:       tt.mockUpdate,
			}

			service := NewStarterApplicationService(
				mockStarterRepo,
				nil,
				nil,
				nil,
				nil,
			)

			starter, err := service.UpdateStarter(context.Background(), tt.command)

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

			if starter == nil {
				t.Error("expected non-nil starter")
			}
		})
	}
}

func TestSoftDeleteStarter(t *testing.T) {
	deletedStarter, _ := model.NewStarter(
		"deleteuser",
		"Delete User",
		"delete@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	tests := []struct {
		name           string
		domain         string
		mockSoftDelete func(ctx context.Context, domain string) (*model.Starter, error)
		expectError    bool
	}{
		{
			name:   "successful deletion",
			domain: "deleteuser",
			mockSoftDelete: func(ctx context.Context, domain string) (*model.Starter, error) {
				return deletedStarter, nil
			},
			expectError: false,
		},
		{
			name:   "deletion error",
			domain: "erroruser",
			mockSoftDelete: func(ctx context.Context, domain string) (*model.Starter, error) {
				return nil, errors.New("database error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStarterRepo := &mocks.MockStarterRepository{
				SoftDeleteFunc: tt.mockSoftDelete,
			}

			service := NewStarterApplicationService(
				mockStarterRepo,
				nil,
				nil,
				nil,
				nil,
			)

			err := service.SoftDeleteStarter(context.Background(), tt.domain)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestListStartersFromMySQL(t *testing.T) {
	starters := []*model.Starter{}
	starter1, _ := model.NewStarter("user1", "User One", "user1@vng.com.vn", "0123456789", "", "Developer", nil, nil)
	starter2, _ := model.NewStarter("user2", "User Two", "user2@vng.com.vn", "0987654321", "", "Developer", nil, nil)
	starters = append(starters, starter1, starter2)

	page := 1
	limit := 10

	mockStarterRepo := &mocks.MockStarterRepository{
		SearchByKeywordFunc: func(ctx context.Context, query *starterquery.ListStartersQuery) ([]*model.Starter, int64, error) {
			return starters, 2, nil
		},
	}

	service := NewStarterApplicationService(
		mockStarterRepo,
		nil,
		nil,
		nil,
		nil,
	)

	query := &starterquery.ListStartersQuery{
		Pagination: httputil.ReqPagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	result, err := service.ListStarters(context.Background(), query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 starters, got %d", len(result.Data))
	}

	if result.Pagination.TotalItems != 2 {
		t.Errorf("expected total items 2, got %d", result.Pagination.TotalItems)
	}
}

func TestApplyUpdates(t *testing.T) {
	existingStarter, _ := model.Rehydrate(
		1,
		"testuser",
		"Original Name",
		"original@vng.com.vn",
		"0123456789",
		"111111",
		"Developer",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	service := &StarterApplicationService{}

	newName := "Updated Name"
	newEmail := "updated@vng.com.vn"

	command := &startercommand.UpdateStarterCommand{
		Domain: "testuser",
		Name:   &newName,
		Email:  &newEmail,
	}

	name, email, mobile, workPhone, jobTitle, deptID, lineManagerID := service.applyUpdates(existingStarter, command)

	if name != newName {
		t.Errorf("expected name %s, got %s", newName, name)
	}

	if email != newEmail {
		t.Errorf("expected email %s, got %s", newEmail, email)
	}

	if mobile != existingStarter.Mobile {
		t.Errorf("expected mobile %s, got %s", existingStarter.Mobile, mobile)
	}

	if workPhone != existingStarter.WorkPhone {
		t.Errorf("expected work phone %s, got %s", existingStarter.WorkPhone, workPhone)
	}

	if jobTitle != existingStarter.JobTitle {
		t.Errorf("expected job title %s, got %s", existingStarter.JobTitle, jobTitle)
	}

	if deptID != existingStarter.DepartmentID {
		t.Error("expected department ID to remain unchanged")
	}

	if lineManagerID != existingStarter.LineManagerID {
		t.Error("expected line manager ID to remain unchanged")
	}
}

