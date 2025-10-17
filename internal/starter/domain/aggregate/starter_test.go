package aggregate

import (
	"testing"
	"time"
)

func TestNewStarter(t *testing.T) {
	tests := []struct {
		name          string
		domain        string
		starterName   string
		email         string
		mobile        string
		workPhone     string
		jobTitle      string
		departmentID  *int64
		lineManagerID *int64
		wantError     bool
		errorMsg      string
	}{
		{
			name:          "Valid starter with all fields",
			domain:        "john.doe",
			starterName:   "John Doe",
			email:         "john.doe@company.com",
			mobile:        "0123456789",
			workPhone:     "0987654321",
			jobTitle:      "Software Engineer",
			departmentID:  int64Ptr(5),
			lineManagerID: int64Ptr(10),
			wantError:     false,
		},
		{
			name:          "Valid starter with minimal fields",
			domain:        "jane.smith",
			starterName:   "Jane Smith",
			email:         "",
			mobile:        "0123456789",
			workPhone:     "",
			jobTitle:      "Senior Engineer",
			departmentID:  nil,
			lineManagerID: nil,
			wantError:     false,
		},
		{
			name:          "Missing domain - should fail",
			domain:        "",
			starterName:   "Test User",
			email:         "test@company.com",
			mobile:        "0123456789",
			workPhone:     "",
			jobTitle:      "Engineer",
			departmentID:  nil,
			lineManagerID: nil,
			wantError:     true,
			errorMsg:      "domain is required",
		},
		{
			name:          "Missing name - should fail",
			domain:        "test.user",
			starterName:   "",
			email:         "test@company.com",
			mobile:        "0123456789",
			workPhone:     "",
			jobTitle:      "Engineer",
			departmentID:  nil,
			lineManagerID: nil,
			wantError:     true,
			errorMsg:      "name is required",
		},
		{
			name:          "Missing mobile - should fail",
			domain:        "test.user",
			starterName:   "Test User",
			email:         "test@company.com",
			mobile:        "",
			workPhone:     "",
			jobTitle:      "Engineer",
			departmentID:  nil,
			lineManagerID: nil,
			wantError:     true,
			errorMsg:      "mobile is required",
		},
		{
			name:          "Missing job title - should fail",
			domain:        "test.user",
			starterName:   "Test User",
			email:         "test@company.com",
			mobile:        "0123456789",
			workPhone:     "",
			jobTitle:      "",
			departmentID:  nil,
			lineManagerID: nil,
			wantError:     true,
			errorMsg:      "job title is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starter, err := NewStarter(
				tt.domain,
				tt.starterName,
				tt.email,
				tt.mobile,
				tt.workPhone,
				tt.jobTitle,
				tt.departmentID,
				tt.lineManagerID,
			)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error %q but got nil", tt.errorMsg)
				}
				if err.Error() != tt.errorMsg {
					t.Fatalf("error message = %v, want %v", err.Error(), tt.errorMsg)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if starter.Domain() != tt.domain {
				t.Errorf("Domain = %v, want %v", starter.Domain(), tt.domain)
			}
			if starter.Email() != tt.email {
				t.Errorf("Email = %v, want %v", starter.Email(), tt.email)
			}
			if starter.Mobile() != tt.mobile {
				t.Errorf("Mobile = %v, want %v", starter.Mobile(), tt.mobile)
			}
			if starter.WorkPhone() != tt.workPhone {
				t.Errorf("WorkPhone = %v, want %v", starter.WorkPhone(), tt.workPhone)
			}
			if starter.JobTitle() != tt.jobTitle {
				t.Errorf("JobTitle = %v, want %v", starter.JobTitle(), tt.jobTitle)
			}
		})
	}
}

func TestRehydrate(t *testing.T) {
	now := time.Now()
	deptID := int64(5)
	managerID := int64(10)

	starter := Rehydrate(
		1,
		"john.doe",
		"John Doe",
		"john@company.com",
		"0123456789",
		"0987654321",
		"Software Engineer",
		&deptID,
		&managerID,
		now,
		now,
	)

	if starter.ID() != 1 {
		t.Errorf("ID = %v, want 1", starter.ID())
	}
	if starter.Domain() != "john.doe" {
		t.Errorf("Domain = %v, want john.doe", starter.Domain())
	}
	if starter.Email() != "john@company.com" {
		t.Errorf("Email = %v, want john@company.com", starter.Email())
	}
	if starter.Mobile() != "0123456789" {
		t.Errorf("Mobile = %v, want 0123456789", starter.Mobile())
	}
	if starter.JobTitle() != "Software Engineer" {
		t.Errorf("JobTitle = %v, want Software Engineer", starter.JobTitle())
	}
	if !equalInt64Ptr(starter.DepartmentID(), &deptID) {
		t.Errorf("DepartmentID = %v, want %v", starter.DepartmentID(), deptID)
	}
	if !equalInt64Ptr(starter.LineManagerID(), &managerID) {
		t.Errorf("LineManagerID = %v, want %v", starter.LineManagerID(), managerID)
	}
}

func TestStarter_UpdateInfo(t *testing.T) {
	tests := []struct {
		name          string
		starter       *Starter
		newName       string
		email         string
		mobile        string
		workPhone     string
		jobTitle      string
		departmentID  *int64
		lineManagerID *int64
	}{
		{
			name:          "Valid update - all fields",
			starter:       mustCreateStarter("john.doe", "John Doe", "old@company.com", "1111111111", "", "Old Title", nil, nil),
			newName:       "John Updated",
			email:         "new@company.com",
			mobile:        "2222222222",
			workPhone:     "3333333333",
			jobTitle:      "New Title",
			departmentID:  nil,
			lineManagerID: nil,
		},
		{
			name:          "Valid update - minimal fields",
			starter:       mustCreateStarter("john.doe", "John Doe", "old@company.com", "1111111111", "", "Old Title", nil, nil),
			newName:       "Jane Doe",
			email:         "",
			mobile:        "2222222222",
			workPhone:     "",
			jobTitle:      "New Title",
			departmentID:  nil,
			lineManagerID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldUpdatedAt := tt.starter.UpdatedAt()
			time.Sleep(1 * time.Millisecond) // Ensure time difference

			tt.starter.UpdateInfo(tt.newName, tt.email, tt.mobile, tt.workPhone, tt.jobTitle, tt.departmentID, tt.lineManagerID)

			if tt.starter.Name() != tt.newName {
				t.Errorf("Name = %v, want %v", tt.starter.Name(), tt.newName)
			}
			if tt.starter.Email() != tt.email {
				t.Errorf("Email = %v, want %v", tt.starter.Email(), tt.email)
			}
			if tt.starter.Mobile() != tt.mobile {
				t.Errorf("Mobile = %v, want %v", tt.starter.Mobile(), tt.mobile)
			}
			if tt.starter.WorkPhone() != tt.workPhone {
				t.Errorf("WorkPhone = %v, want %v", tt.starter.WorkPhone(), tt.workPhone)
			}
			if tt.starter.JobTitle() != tt.jobTitle {
				t.Errorf("JobTitle = %v, want %v", tt.starter.JobTitle(), tt.jobTitle)
			}
			if !tt.starter.UpdatedAt().After(oldUpdatedAt) {
				t.Errorf("UpdatedAt should be updated")
			}
		})
	}
}

func TestStarter_AssignToDepartment(t *testing.T) {
	starter := mustCreateStarter("john.doe", "John Doe", "john@company.com", "0123456789", "", "Engineer", nil, nil)
	oldUpdatedAt := starter.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	deptID := int64(5)
	starter.AssignToDepartment(deptID)

	if !equalInt64Ptr(starter.DepartmentID(), &deptID) {
		t.Errorf("DepartmentID = %v, want %v", starter.DepartmentID(), deptID)
	}
	if !starter.UpdatedAt().After(oldUpdatedAt) {
		t.Errorf("UpdatedAt should be updated")
	}
}

func TestStarter_AssignLineManager(t *testing.T) {
	starter := mustCreateStarter("john.doe", "John Doe", "john@company.com", "0123456789", "", "Engineer", nil, nil)
	oldUpdatedAt := starter.UpdatedAt()
	time.Sleep(1 * time.Millisecond)

	tests := []struct {
		name      string
		managerID *int64
	}{
		{
			name:      "Assign manager",
			managerID: int64Ptr(10),
		},
		{
			name:      "Remove manager (set to nil)",
			managerID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starter.AssignLineManager(tt.managerID)

			if !equalInt64Ptr(starter.LineManagerID(), tt.managerID) {
				t.Errorf("LineManagerID = %v, want %v", starter.LineManagerID(), tt.managerID)
			}
			if !starter.UpdatedAt().After(oldUpdatedAt) {
				t.Errorf("UpdatedAt should be updated")
			}
		})
	}
}

func TestStarter_CanBeDeleted(t *testing.T) {
	tests := []struct {
		name                 string
		hasSubordinates      bool
		isDepartmentLeader   bool
		isBusinessUnitLeader bool
		expectedCanBeDeleted bool
	}{
		{
			name:                 "Can be deleted - no dependencies",
			hasSubordinates:      false,
			isDepartmentLeader:   false,
			isBusinessUnitLeader: false,
			expectedCanBeDeleted: true,
		},
		{
			name:                 "Cannot be deleted - has subordinates",
			hasSubordinates:      true,
			isDepartmentLeader:   false,
			isBusinessUnitLeader: false,
			expectedCanBeDeleted: false,
		},
		{
			name:                 "Cannot be deleted - is department leader",
			hasSubordinates:      false,
			isDepartmentLeader:   true,
			isBusinessUnitLeader: false,
			expectedCanBeDeleted: false,
		},
		{
			name:                 "Cannot be deleted - is business unit leader",
			hasSubordinates:      false,
			isDepartmentLeader:   false,
			isBusinessUnitLeader: true,
			expectedCanBeDeleted: false,
		},
		{
			name:                 "Cannot be deleted - multiple dependencies",
			hasSubordinates:      true,
			isDepartmentLeader:   true,
			isBusinessUnitLeader: true,
			expectedCanBeDeleted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			starter := mustCreateStarter("john.doe", "John Doe", "john@company.com", "0123456789", "", "Engineer", nil, nil)
			canDelete := starter.CanBeDeleted(tt.hasSubordinates, tt.isDepartmentLeader, tt.isBusinessUnitLeader)

			if canDelete != tt.expectedCanBeDeleted {
				t.Errorf("CanBeDeleted = %v, want %v", canDelete, tt.expectedCanBeDeleted)
			}
		})
	}
}

// Helper functions.
func int64Ptr(i int64) *int64 {
	return &i
}

func equalInt64Ptr(a, b *int64) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func mustCreateStarter(domain, name, email, mobile, workPhone, jobTitle string, deptID, managerID *int64) *Starter {
	starter, err := NewStarter(domain, name, email, mobile, workPhone, jobTitle, deptID, managerID)
	if err != nil {
		panic(err)
	}
	return starter
}
