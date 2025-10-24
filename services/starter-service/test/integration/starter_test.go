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

func TestStarter_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("Create Starter", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		payload := map[string]interface{}{
			"domain":     "testuser",
			"name":       "Test User",
			"email":      "testuser@vng.com.vn",
			"mobile":     "+84901234567",
			"work_phone": "1234567",
			"job_title":  "Software Engineer",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Logf("Response body: %s", w.Body.String())
		}
		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		require.True(t, ok, "Expected data field in response")
		assert.Equal(t, "testuser", data["domain"])
		assert.Equal(t, "Test User", data["name"])
		assert.Equal(t, "testuser@vng.com.vn", data["email"])
	})

	t.Run("List Starters", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a test starter first
		createStarter(t, env, "listuser1", "List User 1", "listuser1@vng.com.vn")
		createStarter(t, env, "listuser2", "List User 2", "listuser2@vng.com.vn")

		// List starters
		req := httptest.NewRequest(http.MethodGet, "/api/v1/starters?limit=10", nil)
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 2)
	})

	t.Run("Get Starter by Domain", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a test starter
		createStarter(t, env, "getuser", "Get User", "getuser@vng.com.vn")

		// Get starter
		req := httptest.NewRequest(http.MethodGet, "/api/v1/starters/getuser", nil)
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "getuser", data["domain"])
		assert.Equal(t, "Get User", data["name"])
	})

	t.Run("Update Starter", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a test starter
		createStarter(t, env, "updateuser", "Update User", "updateuser@vng.com.vn")

		// Update starter
		payload := map[string]interface{}{
			"domain":     "updateuser",
			"name":       "Updated User",
			"email":      "updateuser@vng.com.vn",
			"mobile":     "+84909999999",
			"work_phone": "9999999",
			"job_title":  "Senior Software Engineer",
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/starters/updateuser", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "Updated User", data["name"])
		assert.Equal(t, "Senior Software Engineer", data["jobTitle"])
	})

	t.Run("Delete Starter", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create a test starter
		createStarter(t, env, "deleteuser", "Delete User", "deleteuser@vng.com.vn")

		// Delete starter
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/starters/deleteuser", nil)
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Verify starter is deleted (soft delete)
		req = httptest.NewRequest(http.MethodGet, "/api/v1/starters/deleteuser", nil)
		w = httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		// Should return error (404 or similar)
		assert.NotEqual(t, http.StatusOK, w.Code)
	})

	t.Run("Create Starter - Validation Errors", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		testCases := []struct {
			name    string
			payload map[string]interface{}
		}{
			{
				name: "Missing domain",
				payload: map[string]interface{}{
					"name":       "Test User",
					"email":      "test@vng.com.vn",
					"mobile":     "+84901234567",
					"work_phone": "1234567",
					"job_title":  "Engineer",
				},
			},
			{
				name: "Invalid email",
				payload: map[string]interface{}{
					"domain":     "testuser",
					"name":       "Test User",
					"email":      "invalid-email",
					"mobile":     "+84901234567",
					"work_phone": "1234567",
					"job_title":  "Engineer",
				},
			},
			{
				name: "Missing required fields",
				payload: map[string]interface{}{
					"domain": "testuser",
				},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				body, _ := json.Marshal(tc.payload)
				req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				env.Router.ServeHTTP(w, req)

				assert.NotEqual(t, http.StatusOK, w.Code)
			})
		}
	})

	t.Run("Get Non-existent Starter", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		req := httptest.NewRequest(http.MethodGet, "/api/v1/starters/nonexistent", nil)
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

// Helper function to create a starter
func createStarter(t *testing.T, env *TestEnv, domain, name, email string) map[string]interface{} {
	payload := map[string]interface{}{
		"domain":     domain,
		"name":       name,
		"email":      email,
		"mobile":     "+84901234567",
		"work_phone": "1234567",
		"job_title":  "Software Engineer",
	}

	body, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code, "Failed to create starter: %s", w.Body.String())

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	return response["data"].(map[string]interface{})
}

func TestStarter_WithDepartmentAndManager_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	env := SetupTestEnvironment(t)
	defer env.Cleanup()

	t.Run("Create Starter with Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Use existing department from migration (id=9)
		departmentID := int64(9)

		payload := map[string]interface{}{
			"domain":        "deptuser",
			"name":          "Department User",
			"email":         "deptuser@vng.com.vn",
			"mobile":        "+84901234567",
			"work_phone":    "1234567",
			"job_title":     "Software Engineer",
			"department_id": departmentID,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "deptuser", data["domain"])
		assert.NotNil(t, data["department"])
	})

	t.Run("Create Starter with Line Manager", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		// Create manager first
		manager := createStarter(t, env, "manager1", "Manager One", "manager1@vng.com.vn")
		managerID := int64(manager["id"].(float64))

		// Create employee with manager
		payload := map[string]interface{}{
			"domain":          "employee1",
			"name":            "Employee One",
			"email":           "employee1@vng.com.vn",
			"mobile":          "+84901234567",
			"work_phone":      "1234567",
			"job_title":       "Junior Engineer",
			"line_manager_id": managerID,
		}

		body, _ := json.Marshal(payload)
		req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].(map[string]interface{})
		assert.Equal(t, "employee1", data["domain"])
		// Line manager info should be enriched
		if lineManager, ok := data["lineManager"].(map[string]interface{}); ok {
			assert.Equal(t, "manager1", lineManager["domain"])
		}
	})

	t.Run("List Starters with Filter by Department", func(t *testing.T) {
		CleanupDatabase(t, env.DB)

		departmentID := int64(9)

		// Create starters in same department
		for i := 1; i <= 3; i++ {
			payload := map[string]interface{}{
				"domain":        fmt.Sprintf("filteruser%d", i),
				"name":          fmt.Sprintf("Filter User %d", i),
				"email":         fmt.Sprintf("filteruser%d@vng.com.vn", i),
				"mobile":        "+84901234567",
				"work_phone":    "1234567",
				"job_title":     "Engineer",
				"department_id": departmentID,
			}

			body, _ := json.Marshal(payload)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/starters", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			env.Router.ServeHTTP(w, req)
			require.Equal(t, http.StatusOK, w.Code)
		}

		// List starters filtered by department
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/starters?departmentId=%d&limit=10", departmentID), nil)
		w := httptest.NewRecorder()

		env.Router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err)

		data := response["data"].([]interface{})
		assert.GreaterOrEqual(t, len(data), 3)
	})
}

