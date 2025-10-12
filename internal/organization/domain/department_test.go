package domain

import (
	"testing"
	"time"
)

func TestDepartment_UpdateInfo(t *testing.T) {
	tests := []struct {
		name              string
		initialDepartment *Department
		fullName          *string
		shortname         *string
		expectedFullName  string
		expectedShortname string
	}{
		{
			name: "Update both full_name and shortname",
			initialDepartment: &Department{
				ID:        1,
				FullName:  "Old Department Name",
				Shortname: "ODN",
			},
			fullName:          stringPtr("New Department Name"),
			shortname:         stringPtr("NDN"),
			expectedFullName:  "New Department Name",
			expectedShortname: "NDN",
		},
		{
			name: "Update only full_name",
			initialDepartment: &Department{
				ID:        1,
				FullName:  "Old Department Name",
				Shortname: "ODN",
			},
			fullName:          stringPtr("Updated Name"),
			shortname:         nil,
			expectedFullName:  "Updated Name",
			expectedShortname: "ODN",
		},
		{
			name: "Update only shortname",
			initialDepartment: &Department{
				ID:        1,
				FullName:  "Old Department Name",
				Shortname: "ODN",
			},
			fullName:          nil,
			shortname:         stringPtr("USN"),
			expectedFullName:  "Old Department Name",
			expectedShortname: "USN",
		},
		{
			name: "Update with nil values - no changes",
			initialDepartment: &Department{
				ID:        1,
				FullName:  "Old Department Name",
				Shortname: "ODN",
			},
			fullName:          nil,
			shortname:         nil,
			expectedFullName:  "Old Department Name",
			expectedShortname: "ODN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept := tt.initialDepartment
			dept.UpdateInfo(tt.fullName, tt.shortname)

			if dept.FullName != tt.expectedFullName {
				t.Errorf("FullName = %v, want %v", dept.FullName, tt.expectedFullName)
			}
			if dept.Shortname != tt.expectedShortname {
				t.Errorf("Shortname = %v, want %v", dept.Shortname, tt.expectedShortname)
			}
		})
	}
}

func TestDepartment_AssignToBusinessUnit(t *testing.T) {
	tests := []struct {
		name               string
		initialDepartment  *Department
		businessUnitID     *int64
		expectedBusinessID *int64
	}{
		{
			name: "Assign to business unit",
			initialDepartment: &Department{
				ID:             1,
				BusinessUnitID: nil,
			},
			businessUnitID:     int64Ptr(5),
			expectedBusinessID: int64Ptr(5),
		},
		{
			name: "Reassign to different business unit",
			initialDepartment: &Department{
				ID:             1,
				BusinessUnitID: int64Ptr(3),
			},
			businessUnitID:     int64Ptr(7),
			expectedBusinessID: int64Ptr(7),
		},
		{
			name: "Remove business unit assignment",
			initialDepartment: &Department{
				ID:             1,
				BusinessUnitID: int64Ptr(5),
			},
			businessUnitID:     nil,
			expectedBusinessID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept := tt.initialDepartment
			dept.AssignToBusinessUnit(tt.businessUnitID)

			if !equalInt64Ptr(dept.BusinessUnitID, tt.expectedBusinessID) {
				t.Errorf("BusinessUnitID = %v, want %v", dept.BusinessUnitID, tt.expectedBusinessID)
			}
		})
	}
}

func TestDepartment_AssignToGroupDepartment(t *testing.T) {
	tests := []struct {
		name                string
		initialDepartment   *Department
		groupDepartmentID   *int64
		expectedGroupDeptID *int64
	}{
		{
			name: "Assign to parent department",
			initialDepartment: &Department{
				ID:                1,
				GroupDepartmentID: nil,
			},
			groupDepartmentID:   int64Ptr(10),
			expectedGroupDeptID: int64Ptr(10),
		},
		{
			name: "Remove parent department",
			initialDepartment: &Department{
				ID:                1,
				GroupDepartmentID: int64Ptr(10),
			},
			groupDepartmentID:   nil,
			expectedGroupDeptID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept := tt.initialDepartment
			dept.AssignToGroupDepartment(tt.groupDepartmentID)

			if !equalInt64Ptr(dept.GroupDepartmentID, tt.expectedGroupDeptID) {
				t.Errorf("GroupDepartmentID = %v, want %v", dept.GroupDepartmentID, tt.expectedGroupDeptID)
			}
		})
	}
}

func TestDepartment_AssignLeader(t *testing.T) {
	tests := []struct {
		name              string
		initialDepartment *Department
		leaderID          *int64
		expectedLeaderID  *int64
	}{
		{
			name: "Assign leader",
			initialDepartment: &Department{
				ID:       1,
				LeaderID: nil,
			},
			leaderID:         int64Ptr(25),
			expectedLeaderID: int64Ptr(25),
		},
		{
			name: "Reassign to different leader",
			initialDepartment: &Department{
				ID:       1,
				LeaderID: int64Ptr(20),
			},
			leaderID:         int64Ptr(30),
			expectedLeaderID: int64Ptr(30),
		},
		{
			name: "Remove leader",
			initialDepartment: &Department{
				ID:       1,
				LeaderID: int64Ptr(25),
			},
			leaderID:         nil,
			expectedLeaderID: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dept := tt.initialDepartment
			dept.AssignLeader(tt.leaderID)

			if !equalInt64Ptr(dept.LeaderID, tt.expectedLeaderID) {
				t.Errorf("LeaderID = %v, want %v", dept.LeaderID, tt.expectedLeaderID)
			}
		})
	}
}

func TestDepartment_ComplexUpdate(t *testing.T) {
	// Test updating multiple fields in sequence
	dept := &Department{
		ID:                1,
		FullName:          "Original Department",
		Shortname:         "OD",
		BusinessUnitID:    nil,
		GroupDepartmentID: nil,
		LeaderID:          nil,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Update info
	dept.UpdateInfo(stringPtr("Updated Department"), stringPtr("UD"))
	if dept.FullName != "Updated Department" || dept.Shortname != "UD" {
		t.Errorf("UpdateInfo failed: FullName=%s, Shortname=%s", dept.FullName, dept.Shortname)
	}

	// Assign to business unit
	dept.AssignToBusinessUnit(int64Ptr(5))
	if !equalInt64Ptr(dept.BusinessUnitID, int64Ptr(5)) {
		t.Errorf("AssignToBusinessUnit failed: BusinessUnitID=%v", dept.BusinessUnitID)
	}

	// Assign to parent department
	dept.AssignToGroupDepartment(int64Ptr(10))
	if !equalInt64Ptr(dept.GroupDepartmentID, int64Ptr(10)) {
		t.Errorf("AssignToGroupDepartment failed: GroupDepartmentID=%v", dept.GroupDepartmentID)
	}

	// Assign leader
	dept.AssignLeader(int64Ptr(25))
	if !equalInt64Ptr(dept.LeaderID, int64Ptr(25)) {
		t.Errorf("AssignLeader failed: LeaderID=%v", dept.LeaderID)
	}

	// Verify all fields
	if dept.FullName != "Updated Department" {
		t.Errorf("FullName = %v, want 'Updated Department'", dept.FullName)
	}
	if dept.Shortname != "UD" {
		t.Errorf("Shortname = %v, want 'UD'", dept.Shortname)
	}
	if !equalInt64Ptr(dept.BusinessUnitID, int64Ptr(5)) {
		t.Errorf("BusinessUnitID = %v, want 5", dept.BusinessUnitID)
	}
	if !equalInt64Ptr(dept.GroupDepartmentID, int64Ptr(10)) {
		t.Errorf("GroupDepartmentID = %v, want 10", dept.GroupDepartmentID)
	}
	if !equalInt64Ptr(dept.LeaderID, int64Ptr(25)) {
		t.Errorf("LeaderID = %v, want 25", dept.LeaderID)
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

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
