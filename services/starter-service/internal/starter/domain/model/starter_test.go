package model

import (
	"testing"
	"time"
)

func TestNewStarter(t *testing.T) {
	deptID := int64(1)
	lineManagerID := int64(2)

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
		expectError   bool
	}{
		{
			name:          "valid starter",
			domain:        "testdomain",
			starterName:   "Test User",
			email:         "test@vng.com.vn",
			mobile:        "0123456789",
			workPhone:     "0987654321",
			jobTitle:      "Developer",
			departmentID:  &deptID,
			lineManagerID: &lineManagerID,
			expectError:   false,
		},
		{
			name:          "valid starter without optional fields",
			domain:        "testdomain",
			starterName:   "Test User",
			email:         "test@vng.com.vn",
			mobile:        "0123456789",
			workPhone:     "",
			jobTitle:      "Developer",
			departmentID:  nil,
			lineManagerID: nil,
			expectError:   false,
		},
		{
			name:        "empty domain",
			domain:      "",
			starterName: "Test User",
			email:       "test@vng.com.vn",
			mobile:      "0123456789",
			jobTitle:    "Developer",
			expectError: true,
		},
		{
			name:        "empty name",
			domain:      "testdomain",
			starterName: "",
			email:       "test@vng.com.vn",
			mobile:      "0123456789",
			jobTitle:    "Developer",
			expectError: true,
		},
		{
			name:        "empty mobile",
			domain:      "testdomain",
			starterName: "Test User",
			email:       "test@vng.com.vn",
			mobile:      "",
			jobTitle:    "Developer",
			expectError: true,
		},
		{
			name:        "empty job title",
			domain:      "testdomain",
			starterName: "Test User",
			email:       "test@vng.com.vn",
			mobile:      "0123456789",
			jobTitle:    "",
			expectError: true,
		},
		{
			name:        "invalid email",
			domain:      "testdomain",
			starterName: "Test User",
			email:       "test@gmail.com",
			mobile:      "0123456789",
			jobTitle:    "Developer",
			expectError: true,
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

			if starter.Domain != tt.domain {
				t.Errorf("expected domain %s, got %s", tt.domain, starter.Domain)
			}
			if starter.Name != tt.starterName {
				t.Errorf("expected name %s, got %s", tt.starterName, starter.Name)
			}
			if starter.GetEmail() != tt.email {
				t.Errorf("expected email %s, got %s", tt.email, starter.GetEmail())
			}
			if starter.Mobile != tt.mobile {
				t.Errorf("expected mobile %s, got %s", tt.mobile, starter.Mobile)
			}
			if starter.WorkPhone != tt.workPhone {
				t.Errorf("expected work phone %s, got %s", tt.workPhone, starter.WorkPhone)
			}
			if starter.JobTitle != tt.jobTitle {
				t.Errorf("expected job title %s, got %s", tt.jobTitle, starter.JobTitle)
			}
		})
	}
}

func TestRehydrate(t *testing.T) {
	id := int64(1)
	domain := "testdomain"
	name := "Test User"
	email := "test@vng.com.vn"
	mobile := "0123456789"
	workPhone := "0987654321"
	jobTitle := "Developer"
	deptID := int64(1)
	lineManagerID := int64(2)
	createdAt := time.Now()
	updatedAt := time.Now()

	starter, err := Rehydrate(
		id,
		domain,
		name,
		email,
		mobile,
		workPhone,
		jobTitle,
		&deptID,
		&lineManagerID,
		createdAt,
		updatedAt,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if starter.ID != id {
		t.Errorf("expected id %d, got %d", id, starter.ID)
	}
	if starter.Domain != domain {
		t.Errorf("expected domain %s, got %s", domain, starter.Domain)
	}
	if starter.Name != name {
		t.Errorf("expected name %s, got %s", name, starter.Name)
	}
	if starter.GetEmail() != email {
		t.Errorf("expected email %s, got %s", email, starter.GetEmail())
	}
	if starter.Mobile != mobile {
		t.Errorf("expected mobile %s, got %s", mobile, starter.Mobile)
	}
	if starter.WorkPhone != workPhone {
		t.Errorf("expected work phone %s, got %s", workPhone, starter.WorkPhone)
	}
	if starter.JobTitle != jobTitle {
		t.Errorf("expected job title %s, got %s", jobTitle, starter.JobTitle)
	}
	if starter.DepartmentID == nil || *starter.DepartmentID != deptID {
		t.Errorf("expected department id %d, got %v", deptID, starter.DepartmentID)
	}
	if starter.LineManagerID == nil || *starter.LineManagerID != lineManagerID {
		t.Errorf("expected line manager id %d, got %v", lineManagerID, starter.LineManagerID)
	}
}

func TestRehydrateInvalidEmail(t *testing.T) {
	_, err := Rehydrate(
		1,
		"testdomain",
		"Test User",
		"invalid@gmail.com",
		"0123456789",
		"0987654321",
		"Developer",
		nil,
		nil,
		time.Now(),
		time.Now(),
	)

	if err == nil {
		t.Error("expected error for invalid email but got nil")
	}
}

func TestStarterUpdateInfo(t *testing.T) {
	deptID := int64(1)
	lineManagerID := int64(2)

	starter, err := NewStarter(
		"testdomain",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"0987654321",
		"Developer",
		&deptID,
		&lineManagerID,
	)
	if err != nil {
		t.Fatalf("failed to create starter: %v", err)
	}

	newName := "Updated User"
	newEmail := "updated@vng.com.vn"
	newMobile := "1111111111"
	newWorkPhone := "2222222222"
	newJobTitle := "Senior Developer"
	newDeptID := int64(3)
	newLineManagerID := int64(4)

	err = starter.UpdateInfo(
		newName,
		newEmail,
		newMobile,
		newWorkPhone,
		newJobTitle,
		&newDeptID,
		&newLineManagerID,
	)
	if err != nil {
		t.Fatalf("unexpected error during update: %v", err)
	}

	if starter.Name != newName {
		t.Errorf("expected name %s, got %s", newName, starter.Name)
	}
	if starter.GetEmail() != newEmail {
		t.Errorf("expected email %s, got %s", newEmail, starter.GetEmail())
	}
	if starter.Mobile != newMobile {
		t.Errorf("expected mobile %s, got %s", newMobile, starter.Mobile)
	}
	if starter.WorkPhone != newWorkPhone {
		t.Errorf("expected work phone %s, got %s", newWorkPhone, starter.WorkPhone)
	}
	if starter.JobTitle != newJobTitle {
		t.Errorf("expected job title %s, got %s", newJobTitle, starter.JobTitle)
	}
	if starter.DepartmentID == nil || *starter.DepartmentID != newDeptID {
		t.Errorf("expected department id %d, got %v", newDeptID, starter.DepartmentID)
	}
	if starter.LineManagerID == nil || *starter.LineManagerID != newLineManagerID {
		t.Errorf("expected line manager id %d, got %v", newLineManagerID, starter.LineManagerID)
	}
}

func TestStarterUpdateInfoInvalidEmail(t *testing.T) {
	starter, _ := NewStarter(
		"testdomain",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	err := starter.UpdateInfo(
		"Updated User",
		"invalid@gmail.com",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	if err == nil {
		t.Error("expected error for invalid email but got nil")
	}
}

func TestNewStarterESDocFromStarter(t *testing.T) {
	deptID := int64(1)
	starter, _ := NewStarter(
		"testdomain",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		&deptID,
		nil,
	)

	enrichedData := &EnrichedData{
		Departments: map[int64]*DepartmentNested{
			1: {
				ID:        1,
				Name:      "Engineering",
				Shortname: "ENG",
			},
		},
		BusinessUnits: map[int64]*BusinessUnitNested{
			1: {
				ID:        1,
				Name:      "Technology",
				Shortname: "TECH",
			},
		},
		LineManagers: map[int64]*LineManagerNested{},
	}

	esDoc := NewStarterESDocFromStarter(starter, enrichedData)

	if esDoc == nil {
		t.Fatal("expected non-nil ES doc")
	}

	if esDoc.Domain() != starter.Domain {
		t.Errorf("expected domain %s, got %s", starter.Domain, esDoc.Domain())
	}
	if esDoc.Name() != starter.Name {
		t.Errorf("expected name %s, got %s", starter.Name, esDoc.Name())
	}
	if esDoc.DepartmentName() != "Engineering" {
		t.Errorf("expected department name Engineering, got %s", esDoc.DepartmentName())
	}
	if esDoc.BusinessUnitName() != "Technology" {
		t.Errorf("expected business unit name Technology, got %s", esDoc.BusinessUnitName())
	}
}

func TestNewStarterESDocFromStarterNilStarter(t *testing.T) {
	esDoc := NewStarterESDocFromStarter(nil, nil)
	if esDoc != nil {
		t.Error("expected nil ES doc for nil starter")
	}
}

func TestNewStarterESDocFromStarterNilEnrichedData(t *testing.T) {
	starter, _ := NewStarter(
		"testdomain",
		"Test User",
		"test@vng.com.vn",
		"0123456789",
		"",
		"Developer",
		nil,
		nil,
	)

	esDoc := NewStarterESDocFromStarter(starter, nil)

	if esDoc == nil {
		t.Fatal("expected non-nil ES doc")
	}

	if esDoc.DepartmentName() != "" {
		t.Errorf("expected empty department name, got %s", esDoc.DepartmentName())
	}
	if esDoc.BusinessUnitName() != "" {
		t.Errorf("expected empty business unit name, got %s", esDoc.BusinessUnitName())
	}
}

