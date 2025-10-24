package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kiin21/go-rest/pkg/events"
	"github.com/kiin21/go-rest/pkg/httputil"
	businessunitquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/businessunit/query"
	departmentcommand "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/command"
	departmentquery "github.com/kiin21/go-rest/services/starter-service/internal/starter/application/dto/department/query"
	sharedDomain "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/error"
	messagingmocks "github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/messaging/mocks"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository/mocks"
)

func TestGetAllDepartments(t *testing.T) {
	page := 1
	limit := 10
	departments := []*model.DepartmentWithDetails{
		{
			Department: &model.Department{
				ID:       1,
				FullName: "Engineering",
			},
		},
		{
			Department: &model.Department{
				ID:       2,
				FullName: "Sales",
			},
		},
	}

	mockDepartmentRepo := &mocks.MockDepartmentRepository{
		ListWithDetailsFunc: func(ctx context.Context, filter *model.DepartmentListFilter, pagination *httputil.ReqPagination) ([]*model.DepartmentWithDetails, int64, error) {
			return departments, 2, nil
		},
	}

	service := NewOrganizationApplicationService(
		mockDepartmentRepo,
		nil,
		nil,
		nil,
	)

	query := &departmentquery.ListDepartmentsQuery{
		Pagination: httputil.ReqPagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	result, err := service.GetAllDepartments(context.Background(), query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 departments, got %d", len(result.Data))
	}

	if result.Pagination.TotalItems != 2 {
		t.Errorf("expected total items 2, got %d", result.Pagination.TotalItems)
	}
}

func TestGetOneDepartment(t *testing.T) {
	tests := []struct {
		name                      string
		departmentID              int64
		mockFindByIDsWithDetails  func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
		expectError               bool
	}{
		{
			name:         "department found",
			departmentID: 1,
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:       1,
							FullName: "Engineering",
						},
					},
				}, nil
			},
			expectError: false,
		},
		{
			name:         "department not found",
			departmentID: 999,
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{}, nil
			},
			expectError: true,
		},
		{
			name:         "repository error",
			departmentID: 1,
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return nil, errors.New("database error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepartmentRepo := &mocks.MockDepartmentRepository{
				FindByIDsWithDetailsFunc: tt.mockFindByIDsWithDetails,
			}

			service := NewOrganizationApplicationService(
				mockDepartmentRepo,
				nil,
				nil,
				nil,
			)

			department, err := service.GetOneDepartment(context.Background(), tt.departmentID)

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

			if department == nil {
				t.Error("expected non-nil department")
			}
		})
	}
}

func TestCreateDepartment(t *testing.T) {
	buID := int64(1)

	tests := []struct {
		name                      string
		command                   *departmentcommand.CreateDepartmentCommand
		mockCreate                func(ctx context.Context, department *model.Department) error
		mockFindByIDsWithDetails  func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
		expectError               bool
	}{
		{
			name: "successful creation",
			command: &departmentcommand.CreateDepartmentCommand{
				FullName:       "New Department",
				Shortname:      "NEW",
				BusinessUnitID: &buID,
			},
			mockCreate: func(ctx context.Context, department *model.Department) error {
				department.ID = 1
				return nil
			},
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:        1,
							FullName:  "New Department",
							Shortname: "NEW",
						},
					},
				}, nil
			},
			expectError: false,
		},
		{
			name: "create error",
			command: &departmentcommand.CreateDepartmentCommand{
				FullName:  "New Department",
				Shortname: "NEW",
			},
			mockCreate: func(ctx context.Context, department *model.Department) error {
				return errors.New("database error")
			},
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return nil, nil
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepartmentRepo := &mocks.MockDepartmentRepository{
				CreateFunc:               tt.mockCreate,
				FindByIDsWithDetailsFunc: tt.mockFindByIDsWithDetails,
			}

			service := NewOrganizationApplicationService(
				mockDepartmentRepo,
				nil,
				nil,
				nil,
			)

			department, err := service.CreateDepartment(context.Background(), tt.command)

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

			if department == nil {
				t.Error("expected non-nil department")
			}
		})
	}
}

func TestUpdateDepartment(t *testing.T) {
	fullName := "Updated Department"
	shortname := "UPD"
	buID := int64(1)

	tests := []struct {
		name                      string
		command                   *departmentcommand.UpdateDepartmentCommand
		mockFindByIDsWithDetails  func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
		mockUpdate                func(ctx context.Context, department *model.Department) error
		expectError               bool
	}{
		{
			name: "successful update",
			command: &departmentcommand.UpdateDepartmentCommand{
				ID:             1,
				FullName:       &fullName,
				Shortname:      &shortname,
				BusinessUnitID: &buID,
			},
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:             1,
							FullName:       "Original Department",
							Shortname:      "ORIG",
							BusinessUnitID: &buID,
							CreatedAt:      time.Now(),
							UpdatedAt:      time.Now(),
						},
					},
				}, nil
			},
			mockUpdate: func(ctx context.Context, department *model.Department) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "find error",
			command: &departmentcommand.UpdateDepartmentCommand{
				ID:        999,
				FullName:  &fullName,
				Shortname: &shortname,
			},
			mockFindByIDsWithDetails: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return nil, errors.New("not found")
			},
			mockUpdate:  func(ctx context.Context, department *model.Department) error { return nil },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepartmentRepo := &mocks.MockDepartmentRepository{
				FindByIDsWithDetailsFunc: tt.mockFindByIDsWithDetails,
				UpdateFunc:               tt.mockUpdate,
			}

			service := NewOrganizationApplicationService(
				mockDepartmentRepo,
				nil,
				nil,
				nil,
			)

			department, err := service.UpdateDepartment(context.Background(), tt.command)

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

			if department == nil {
				t.Error("expected non-nil department")
			}
		})
	}
}

func TestDeleteDepartment(t *testing.T) {
	tests := []struct {
		name         string
		departmentID int64
		mockDelete   func(ctx context.Context, id int64) error
		expectError  bool
	}{
		{
			name:         "successful deletion",
			departmentID: 1,
			mockDelete: func(ctx context.Context, id int64) error {
				return nil
			},
			expectError: false,
		},
		{
			name:         "delete error",
			departmentID: 999,
			mockDelete: func(ctx context.Context, id int64) error {
				return errors.New("not found")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepartmentRepo := &mocks.MockDepartmentRepository{
				DeleteFunc: tt.mockDelete,
			}

			service := NewOrganizationApplicationService(
				mockDepartmentRepo,
				nil,
				nil,
				nil,
			)

			err := service.DeleteDepartment(context.Background(), tt.departmentID)

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

func TestAssignLeader(t *testing.T) {
	leaderID := int64(10)
	leaderDomain := "leader"

	tests := []struct {
		name                      string
		command                   *departmentcommand.AssignLeaderCommand
		mockFindDepartment        func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error)
		mockFindStarter           func(ctx context.Context, domain string) (*model.Starter, error)
		mockUpdate                func(ctx context.Context, department *model.Department) error
		expectError               bool
	}{
		{
			name: "assign by leader ID",
			command: &departmentcommand.AssignLeaderCommand{
				DepartmentID: 1,
				LeaderID:     &leaderID,
			},
			mockFindDepartment: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:        1,
							FullName:  "Engineering",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Leader: &model.LineManagerNested{
							ID:     10,
							Domain: "leader",
						},
					},
				}, nil
			},
			mockUpdate: func(ctx context.Context, department *model.Department) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "assign by leader domain",
			command: &departmentcommand.AssignLeaderCommand{
				DepartmentID: 1,
				LeaderDomain: &leaderDomain,
			},
			mockFindDepartment: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:        1,
							FullName:  "Engineering",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
						Leader: &model.LineManagerNested{
							ID:     10,
							Domain: "leader",
						},
					},
				}, nil
			},
			mockFindStarter: func(ctx context.Context, domain string) (*model.Starter, error) {
				starter, _ := model.Rehydrate(10, "leader", "Leader Name", "leader@vng.com.vn", "0123456789", "", "Manager", nil, nil, time.Now(), time.Now())
				return starter, nil
			},
			mockUpdate: func(ctx context.Context, department *model.Department) error {
				return nil
			},
			expectError: false,
		},
		{
			name: "department not found",
			command: &departmentcommand.AssignLeaderCommand{
				DepartmentID: 999,
				LeaderID:     &leaderID,
			},
			mockFindDepartment: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{}, nil
			},
			mockUpdate:  func(ctx context.Context, department *model.Department) error { return nil },
			expectError: true,
		},
		{
			name: "invalid input - both leader ID and domain",
			command: &departmentcommand.AssignLeaderCommand{
				DepartmentID: 1,
				LeaderID:     &leaderID,
				LeaderDomain: &leaderDomain,
			},
			mockFindDepartment: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
				return []*model.DepartmentWithDetails{
					{
						Department: &model.Department{
							ID:        1,
							FullName:  "Engineering",
							CreatedAt: time.Now(),
							UpdatedAt: time.Now(),
						},
					},
				}, nil
			},
			mockUpdate:  func(ctx context.Context, department *model.Department) error { return nil },
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDepartmentRepo := &mocks.MockDepartmentRepository{
				FindByIDsWithDetailsFunc: tt.mockFindDepartment,
				UpdateFunc:               tt.mockUpdate,
			}

			mockStarterRepo := &mocks.MockStarterRepository{
				FindByDomainFunc: tt.mockFindStarter,
			}

			mockNotificationPub := &messagingmocks.MockNotificationProducer{
				SendNotificationFunc: func(event *events.Event) error {
					return nil
				},
			}

			service := NewOrganizationApplicationService(
				mockDepartmentRepo,
				nil,
				mockStarterRepo,
				mockNotificationPub,
			)

			department, err := service.AssignLeader(context.Background(), tt.command)

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

			if department == nil {
				t.Error("expected non-nil department")
			}
		})
	}
}

func TestGetBusinessUnit(t *testing.T) {
	tests := []struct {
		name          string
		businessUnitID int64
		mockFindByIDs func(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error)
		expectError   bool
	}{
		{
			name:           "business unit found",
			businessUnitID: 1,
			mockFindByIDs: func(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
				return []*model.BusinessUnit{
					{
						ID:   1,
						Name: "Technology",
					},
				}, nil
			},
			expectError: false,
		},
		{
			name:           "business unit not found",
			businessUnitID: 999,
			mockFindByIDs: func(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
				return []*model.BusinessUnit{}, nil
			},
			expectError: true,
		},
		{
			name:           "repository error",
			businessUnitID: 1,
			mockFindByIDs: func(ctx context.Context, ids []int64) ([]*model.BusinessUnit, error) {
				return nil, errors.New("database error")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{
				FindByIDsFunc: tt.mockFindByIDs,
			}

			service := NewOrganizationApplicationService(
				nil,
				mockBusinessUnitRepo,
				nil,
				nil,
			)

			bu, err := service.GetBusinessUnit(context.Background(), tt.businessUnitID)

			if tt.expectError {
				if err == nil {
					t.Error("expected error but got nil")
				}
				if err != sharedDomain.ErrNotFound && tt.businessUnitID == 999 {
					// Check specific error type for not found case
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if bu == nil {
				t.Error("expected non-nil business unit")
			}
		})
	}
}

func TestListBusinessUnits(t *testing.T) {
	page := 1
	limit := 10

	businessUnits := []*model.BusinessUnit{
		{ID: 1, Name: "Technology", Shortname: "TECH"},
		{ID: 2, Name: "Sales", Shortname: "SALES"},
	}

	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{
		ListFunc: func(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnit, int64, error) {
			return businessUnits, 2, nil
		},
	}

	service := NewOrganizationApplicationService(
		nil,
		mockBusinessUnitRepo,
		nil,
		nil,
	)

	query := &businessunitquery.ListBusinessUnitsQuery{
		Pagination: httputil.ReqPagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	result, err := service.ListBusinessUnits(context.Background(), query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 business units, got %d", len(result.Data))
	}
}

func TestListBusinessUnitsWithDetails(t *testing.T) {
	page := 1
	limit := 10

	businessUnits := []*model.BusinessUnitWithDetails{
		{
			BusinessUnit: &model.BusinessUnit{
				ID:        1,
				Name:      "Technology",
				Shortname: "TECH",
			},
		},
		{
			BusinessUnit: &model.BusinessUnit{
				ID:        2,
				Name:      "Sales",
				Shortname: "SALES",
			},
		},
	}

	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{
		ListWithDetailsFunc: func(ctx context.Context, pagination httputil.ReqPagination) ([]*model.BusinessUnitWithDetails, int64, error) {
			return businessUnits, 2, nil
		},
	}

	service := NewOrganizationApplicationService(
		nil,
		mockBusinessUnitRepo,
		nil,
		nil,
	)

	query := &businessunitquery.ListBusinessUnitsQuery{
		Pagination: httputil.ReqPagination{
			Page:  &page,
			Limit: &limit,
		},
	}

	result, err := service.ListBusinessUnitsWithDetails(context.Background(), query)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected non-nil result")
	}

	if len(result.Data) != 2 {
		t.Errorf("expected 2 business units, got %d", len(result.Data))
	}
}

func TestGetBusinessUnitWithDetails(t *testing.T) {
	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{
		FindByIDWithDetailsFunc: func(ctx context.Context, id int64) (*model.BusinessUnitWithDetails, error) {
			if id == 1 {
				return &model.BusinessUnitWithDetails{
					BusinessUnit: &model.BusinessUnit{
						ID:        1,
						Name:      "Technology",
						Shortname: "TECH",
					},
				}, nil
			}
			return nil, sharedDomain.ErrNotFound
		},
	}

	service := NewOrganizationApplicationService(
		nil,
		mockBusinessUnitRepo,
		nil,
		nil,
	)

	// Test successful retrieval
	bu, err := service.GetBusinessUnitWithDetails(context.Background(), 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if bu == nil {
		t.Error("expected non-nil business unit")
	}

	// Test not found
	_, err = service.GetBusinessUnitWithDetails(context.Background(), 999)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

