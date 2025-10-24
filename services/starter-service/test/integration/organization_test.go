package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepartment_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("List Departments", func(t *testing.T) {
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/departments?limit=10", nil)
		AssertSuccess(t, w)

		data := ExtractPaginatedData(t, w.Body.Bytes())
		assert.Greater(t, len(data), 0)

		// Check first department structure
		firstDept := data[0]
		assert.NotNil(t, firstDept["id"])
		assert.NotNil(t, firstDept["full_name"])
		assert.NotNil(t, firstDept["shortname"])
	})

	t.Run("Get Department Detail", func(t *testing.T) {
		// Use existing department from migration (id=9)
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/departments/9", nil)
		AssertSuccess(t, w)

		data := ExtractSingleData(t, w.Body.Bytes())
		assert.Equal(t, float64(9), data["id"])
		assert.NotEmpty(t, data["full_name"])
	})

	t.Run("Create Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		payload := map[string]interface{}{
			"full_name":        "Test Department",
			"shortname":        "TD",
			"business_unit_id": 1, // Use existing BU from migration
		}

		w := MakeRequest(t, env, http.MethodPost, "/api/v1/organization/departments", payload)
		
		if w.Code != http.StatusOK {
			t.Logf("Response status: %d", w.Code)
			t.Logf("Response body: %s", w.Body.String())
		}
		
		AssertSuccess(t, w)

		data := ExtractSingleData(t, w.Body.Bytes())
		assert.Equal(t, "Test Department", data["full_name"])
		assert.Equal(t, "TD", data["shortname"])
	})

	t.Run("Update Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a department first
		dept := createDepartment(t, env, "Update Dept", "UD", 1)
		deptID := int64(dept["id"].(float64))

		// Update the department
		payload := map[string]interface{}{
			"full_name":        "Updated Department",
			"shortname":        "UPD",
			"business_unit_id": 1,
		}

		w := MakeRequest(t, env, http.MethodPatch, fmt.Sprintf("/api/v1/organization/departments/%d", deptID), payload)
		AssertSuccess(t, w)

		data := ExtractSingleData(t, w.Body.Bytes())
		assert.Equal(t, "Updated Department", data["full_name"])
		assert.Equal(t, "UPD", data["shortname"])
	})

	t.Run("Assign Leader to Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a department
		dept := createDepartment(t, env, "Leader Dept", "LD", 1)
		deptID := int64(dept["id"].(float64))

		// Create a starter to be the leader
		leader := createStarter(t, env, "leader1", "Leader One", "leader1@vng.com.vn")
		leaderID := int64(leader["id"].(float64))

		// Assign leader
		payload := map[string]interface{}{
			"leader_id": leaderID,
		}

		w := MakeRequest(t, env, http.MethodPatch, fmt.Sprintf("/api/v1/organization/departments/%d/leader", deptID), payload)
		AssertSuccess(t, w)

		data := ExtractSingleData(t, w.Body.Bytes())
		if leaderData, ok := data["leader"].(map[string]interface{}); ok {
			assert.Equal(t, "leader1", leaderData["domain"])
		}
	})

	t.Run("Delete Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a department
		dept := createDepartment(t, env, "Delete Dept", "DD", 1)
		deptID := int64(dept["id"].(float64))

		// Delete the department
		w := MakeRequest(t, env, http.MethodDelete, fmt.Sprintf("/api/v1/organization/departments/%d", deptID), nil)
		AssertSuccess(t, w)

		// Verify department is deleted
		w = MakeRequest(t, env, http.MethodGet, fmt.Sprintf("/api/v1/organization/departments/%d", deptID), nil)
		AssertError(t, w)
	})

	t.Run("Create Department - Validation Errors", func(t *testing.T) {
		testCases := []struct {
			name    string
			payload map[string]interface{}
		}{
			{
				name: "Missing fullName",
				payload: map[string]interface{}{
					"shortname":        "TD",
					"business_unit_id": 1,
				},
			},
			{
				name: "Missing shortname",
				payload: map[string]interface{}{
					"full_name":        "Test Dept",
					"business_unit_id": 1,
				},
			},
			{
				name: "Invalid businessUnitId",
				payload: map[string]interface{}{
					"full_name":        "Test Dept",
					"shortname":        "TD",
					"business_unit_id": 9999,
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/organization/departments", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				env.Router.ServeHTTP(w, req)

				assert.NotEqual(t, http.StatusOK, w.Code)
			})
		}
	})
}

func TestBusinessUnit_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("List Business Units", func(t *testing.T) {
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/business-units?limit=10", nil)
		AssertSuccess(t, w)

		data := ExtractPaginatedData(t, w.Body.Bytes())
		assert.Greater(t, len(data), 0)

		// Check first business unit structure
		firstBU := data[0]
		assert.NotNil(t, firstBU["id"])
		assert.NotNil(t, firstBU["name"])
		assert.NotNil(t, firstBU["shortname"])
	})

	t.Run("Get Business Unit Detail", func(t *testing.T) {
		// Use existing business unit from migration (id=1)
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/business-units/1", nil)
		AssertSuccess(t, w)

		data := ExtractSingleData(t, w.Body.Bytes())
		assert.Equal(t, float64(1), data["id"])
		assert.NotEmpty(t, data["name"])
	})

	t.Run("Get Non-existent Business Unit", func(t *testing.T) {
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/business-units/9999", nil)
		AssertError(t, w)
	})
}

func TestDepartment_WithFilters_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("List Departments by Business Unit", func(t *testing.T) {
		// List departments filtered by business unit id=1
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/departments?businessUnitId=1&limit=10", nil)
		AssertSuccess(t, w)

		data := ExtractPaginatedData(t, w.Body.Bytes())
		// Should have departments belonging to business unit 1
		for _, dept := range data {
			if businessUnit, ok := dept["business_unit"].(map[string]interface{}); ok {
				if businessUnit["id"] != nil {
					assert.Equal(t, float64(1), businessUnit["id"])
				}
			}
		}
	})

	t.Run("List Departments with Pagination", func(t *testing.T) {
		// Request with small limit
		w := MakeRequest(t, env, http.MethodGet, "/api/v1/organization/departments?limit=2", nil)
		AssertSuccess(t, w)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		// Extract pagination data
		dataField := response["data"].(map[string]interface{})
		dataArray := dataField["data"].([]interface{})
		assert.LessOrEqual(t, len(dataArray), 2)

		// Check pagination info
		pagination := dataField["pagination"].(map[string]interface{})
		assert.NotNil(t, pagination["next"])
	})
}

// Helper function to create a department
func createDepartment(t *testing.T, env *TestEnv, fullName, shortname string, businessUnitID int64) map[string]interface{} {
	payload := map[string]interface{}{
		"full_name":        fullName,
		"shortname":        shortname,
		"business_unit_id": businessUnitID,
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/organization/departments", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Failed to create department: %s", w.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response["data"].(map[string]interface{})
}

func TestDepartment_ComplexScenarios_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("Department with Multiple Starters and Leader", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create department
		dept := createDepartment(t, env, "Engineering", "ENG", 1)
		deptID := int64(dept["id"].(float64))

		// Create leader
		leader := createStarter(t, env, "techlead", "Tech Lead", "techlead@vng.com.vn")
		leaderID := int64(leader["id"].(float64))

		// Assign leader to department
		payload := map[string]interface{}{
			"leader_id": leaderID,
		}
		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/api/v1/organization/departments/%d/leader", deptID), bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		env.Router.ServeHTTP(w, req)
		require.Equal(t, http.StatusOK, w.Code)

		// Create multiple team members in the department
		for i := 1; i <= 3; i++ {
			memberPayload := map[string]interface{}{
				"domain":          fmt.Sprintf("engineer%d", i),
				"name":            fmt.Sprintf("Engineer %d", i),
				"email":           fmt.Sprintf("engineer%d@vng.com.vn", i),
				"mobile":          "+84901234567",
				"work_phone":      "1234567",
				"job_title":       "Software Engineer",
				"department_id":   deptID,
				"line_manager_id": leaderID,
			}

			body, _ := json.Marshal(memberPayload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			env.Router.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code)
		}

		// Get department detail - should show leader and member count
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/organization/departments/%d", deptID), nil)
		w = httptest.NewRecorder()
		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Engineering", data["fullName"])

		// Verify leader is set
		if leaderData, ok := data["leader"].(map[string]interface{}); ok {
			assert.Equal(t, "techlead", leaderData["domain"])
		}

		// List starters in this department
		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/starters?departmentId=%d&limit=10", deptID), nil)
		w = httptest.NewRecorder()
		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		err = json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		starters := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(starters), 3) // At least 3 engineers
	})

	t.Run("Hierarchical Department Structure", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create parent department
		parentDept := createDepartment(t, env, "Parent Dept", "PD", 1)
		parentDeptID := int64(parentDept["id"].(float64))

		// Create child department with groupDepartmentId
		childPayload := map[string]interface{}{
			"full_name":            "Child Dept",
			"shortname":            "CD",
			"business_unit_id":     1,
			"group_department_id":  parentDeptID,
		}

		body, _ := json.Marshal(childPayload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/organization/departments", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Child Dept", data["fullName"])
		
		// Verify parent department reference
		if groupDept, ok := data["groupDepartment"].(map[string]interface{}); ok {
			assert.Equal(t, float64(parentDeptID), groupDept["id"])
		}
	})
}

