package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/model"
	"github.com/kiin21/go-rest/services/starter-service/internal/starter/domain/repository/mocks"
)

func TestEnrichStarters(t *testing.T) {
	deptID := int64(1)
	lineManagerID := int64(2)

	starter1, _ := model.Rehydrate(
		1,
		"user1",
		"User One",
		"user1@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		&deptID,
		&lineManagerID,
		time.Now(),
		time.Now(),
	)

	starter2, _ := model.Rehydrate(
		3,
		"user2",
		"User Two",
		"user2@vng.com.vn",
		"0987654321",
		"",
		"Senior Developer",
		&deptID,
		nil,
		time.Now(),
		time.Now(),
	)

	starters := []*model.Starter{starter1, starter2}

	mockDepartmentRepo := &mocks.MockDepartmentRepository{
		FindByIDsWithDetailsFunc: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
			return []*model.DepartmentWithDetails{
				{
					Department: &model.Department{
						ID:        1,
						FullName:  "Engineering",
						Shortname: "ENG",
					},
					BusinessUnit: &model.BusinessUnit{
						ID:        1,
						Name:      "Technology",
						Shortname: "TECH",
					},
					ParentDepartment: nil,
				},
			}, nil
		},
	}

	mockStarterRepo := &mocks.MockStarterRepository{
		FindByIDsFunc: func(ctx context.Context, ids []int64) ([]*model.Starter, error) {
			lineManager, _ := model.Rehydrate(
				2,
				"manager",
				"Line Manager",
				"manager@vng.com.vn",
				"1111111111",
				"",
				"Manager",
				nil,
				nil,
				time.Now(),
				time.Now(),
			)
			return []*model.Starter{lineManager}, nil
		},
	}

	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{}

	service := NewStarterEnrichmentService(mockStarterRepo, mockDepartmentRepo, mockBusinessUnitRepo)

	enriched, err := service.EnrichStarters(context.Background(), starters)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if enriched == nil {
		t.Fatal("expected non-nil enriched data")
	}

	if len(enriched.Departments) != 1 {
		t.Errorf("expected 1 department, got %d", len(enriched.Departments))
	}

	if dept, ok := enriched.Departments[1]; ok {
		if dept.Name != "Engineering" {
			t.Errorf("expected department name Engineering, got %s", dept.Name)
		}
	} else {
		t.Error("expected department with ID 1")
	}

	if len(enriched.LineManagers) != 1 {
		t.Errorf("expected 1 line manager, got %d", len(enriched.LineManagers))
	}

	if manager, ok := enriched.LineManagers[2]; ok {
		if manager.Name != "Line Manager" {
			t.Errorf("expected line manager name 'Line Manager', got %s", manager.Name)
		}
	} else {
		t.Error("expected line manager with ID 2")
	}

	if len(enriched.BusinessUnits) != 1 {
		t.Errorf("expected 1 business unit, got %d", len(enriched.BusinessUnits))
	}
}

func TestEnrichStartersEmptyStarters(t *testing.T) {
	mockDepartmentRepo := &mocks.MockDepartmentRepository{}
	mockStarterRepo := &mocks.MockStarterRepository{}
	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{}

	service := NewStarterEnrichmentService(mockStarterRepo, mockDepartmentRepo, mockBusinessUnitRepo)

	enriched, err := service.EnrichStarters(context.Background(), []*model.Starter{})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if enriched == nil {
		t.Fatal("expected non-nil enriched data")
	}

	if len(enriched.Departments) != 0 {
		t.Errorf("expected 0 departments, got %d", len(enriched.Departments))
	}

	if len(enriched.LineManagers) != 0 {
		t.Errorf("expected 0 line managers, got %d", len(enriched.LineManagers))
	}
}

func TestEnrichStartersDepartmentError(t *testing.T) {
	deptID := int64(1)
	starter, _ := model.Rehydrate(
		1,
		"user1",
		"User One",
		"user1@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		&deptID,
		nil,
		time.Now(),
		time.Now(),
	)

	mockDepartmentRepo := &mocks.MockDepartmentRepository{
		FindByIDsWithDetailsFunc: func(ctx context.Context, ids []int64) ([]*model.DepartmentWithDetails, error) {
			return nil, errors.New("department repository error")
		},
	}

	mockStarterRepo := &mocks.MockStarterRepository{}
	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{}

	service := NewStarterEnrichmentService(mockStarterRepo, mockDepartmentRepo, mockBusinessUnitRepo)

	_, err := service.EnrichStarters(context.Background(), []*model.Starter{starter})

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestEnrichStartersLineManagerError(t *testing.T) {
	lineManagerID := int64(2)
	starter, _ := model.Rehydrate(
		1,
		"user1",
		"User One",
		"user1@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		&lineManagerID,
		time.Now(),
		time.Now(),
	)

	mockDepartmentRepo := &mocks.MockDepartmentRepository{}
	mockStarterRepo := &mocks.MockStarterRepository{
		FindByIDsFunc: func(ctx context.Context, ids []int64) ([]*model.Starter, error) {
			return nil, errors.New("starter repository error")
		},
	}
	mockBusinessUnitRepo := &mocks.MockBusinessUnitRepository{}

	service := NewStarterEnrichmentService(mockStarterRepo, mockDepartmentRepo, mockBusinessUnitRepo)

	_, err := service.EnrichStarters(context.Background(), []*model.Starter{starter})

	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestCollectDepartmentIDs(t *testing.T) {
	deptID1 := int64(1)
	deptID2 := int64(2)

	starter1, _ := model.NewStarter("user1", "User One", "user1@vng.com.vn", "0123456789", "", "Developer", &deptID1, nil)
	starter2, _ := model.NewStarter("user2", "User Two", "user2@vng.com.vn", "0987654321", "", "Developer", &deptID2, nil)
	starter3, _ := model.NewStarter("user3", "User Three", "user3@vng.com.vn", "1111111111", "", "Developer", nil, nil)

	service := &StarterEnrichmentService{}

	deptIDs := service.collectDepartmentIDs([]*model.Starter{starter1, starter2, starter3})

	if len(deptIDs) != 2 {
		t.Errorf("expected 2 department IDs, got %d", len(deptIDs))
	}

	if !deptIDs[1] || !deptIDs[2] {
		t.Error("expected department IDs 1 and 2 to be present")
	}
}

func TestCollectLineManagerIDs(t *testing.T) {
	lineManagerID1 := int64(10)
	lineManagerID2 := int64(20)

	starter1, _ := model.NewStarter("user1", "User One", "user1@vng.com.vn", "0123456789", "", "Developer", nil, &lineManagerID1)
	starter2, _ := model.NewStarter("user2", "User Two", "user2@vng.com.vn", "0987654321", "", "Developer", nil, &lineManagerID2)
	starter3, _ := model.NewStarter("user3", "User Three", "user3@vng.com.vn", "1111111111", "", "Developer", nil, nil)

	service := &StarterEnrichmentService{}

	managerIDs := service.collectLineManagerIDs([]*model.Starter{starter1, starter2, starter3})

	if len(managerIDs) != 2 {
		t.Errorf("expected 2 line manager IDs, got %d", len(managerIDs))
	}

	if !managerIDs[10] || !managerIDs[20] {
		t.Error("expected line manager IDs 10 and 20 to be present")
	}
}

