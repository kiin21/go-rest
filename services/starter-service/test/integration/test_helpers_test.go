package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Data       []map[string]interface{} `json:"data"`
	Pagination map[string]interface{}   `json:"pagination"`
}

// SingleResponse represents a single item API response
type SingleResponse struct {
	Data map[string]interface{} `json:"data"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Error   map[string]interface{} `json:"error,omitempty"`
}

// ExtractPaginatedData extracts the data array from a paginated response
func ExtractPaginatedData(t *testing.T, responseBody []byte) []map[string]interface{} {
	var outerResponse map[string]interface{}
	err := json.Unmarshal(responseBody, &outerResponse)
	require.NoError(t, err)

	dataField, ok := outerResponse["data"].(map[string]interface{})
	require.True(t, ok, "Expected data to be a map")

	dataArray, ok := dataField["data"].([]interface{})
	require.True(t, ok, "Expected data.data to be an array")

	// Convert to []map[string]interface{}
	result := make([]map[string]interface{}, len(dataArray))
	for i, item := range dataArray {
		result[i] = item.(map[string]interface{})
	}

	return result
}

// ExtractSingleData extracts a single item from the response
func ExtractSingleData(t *testing.T, responseBody []byte) map[string]interface{} {
	var outerResponse map[string]interface{}
	err := json.Unmarshal(responseBody, &outerResponse)
	require.NoError(t, err)

	data, ok := outerResponse["data"].(map[string]interface{})
	require.True(t, ok, "Expected data to be a map")

	return data
}

// MakeRequest is a helper to make HTTP requests in tests
func MakeRequest(t *testing.T, env *TestEnv, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(t, err)
		reqBody = bytes.NewBuffer(jsonBody)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	env.Router.ServeHTTP(w, req)

	return w
}

// AssertSuccess checks if the response is successful
func AssertSuccess(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code != http.StatusOK && w.Code != http.StatusCreated {
		t.Logf("Response body: %s", w.Body.String())
		t.Fatalf("Expected success status, got %d", w.Code)
	}
}

// AssertError checks if the response is an error
func AssertError(t *testing.T, w *httptest.ResponseRecorder) {
	if w.Code == http.StatusOK || w.Code == http.StatusCreated {
		t.Fatalf("Expected error status, got %d", w.Code)
	}
}
