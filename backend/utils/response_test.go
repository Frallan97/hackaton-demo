package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSuccessResponse(t *testing.T) {
	response := SuccessResponse("test data", "test message")

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Data != "test data" {
		t.Errorf("Expected Data to be 'test data', got %v", response.Data)
	}

	if response.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %v", response.Message)
	}

	if response.Error != "" {
		t.Error("Expected Error to be empty")
	}

	if response.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}
}

func TestErrorResponse(t *testing.T) {
	err := &testError{message: "test error"}
	response := ErrorResponse("test error message", err)

	if response.Success {
		t.Error("Expected Success to be false")
	}

	expectedError := "test error message: test error"
	if response.Error != expectedError {
		t.Errorf("Expected Error to be '%s', got %v", expectedError, response.Error)
	}

	if response.Data != nil {
		t.Error("Expected Data to be nil")
	}

	if response.Timestamp.IsZero() {
		t.Error("Expected Timestamp to be set")
	}
}

func TestValidationErrorResponse(t *testing.T) {
	errors := map[string]string{
		"field1": "Field 1 is required",
		"field2": "Field 2 is invalid",
	}

	response := ValidationErrorResponse(errors)

	if response.Success {
		t.Error("Expected Success to be false")
	}

	if response.Error != "Validation failed" {
		t.Errorf("Expected Error to be 'Validation failed', got %v", response.Error)
	}

	if response.Data == nil {
		t.Error("Expected Data to contain validation errors")
	}

	// Check if validation errors are properly set
	if response.Data.(map[string]string)["field1"] != "Field 1 is required" {
		t.Error("Expected field1 error to be set correctly")
	}
}

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	WriteJSON(w, http.StatusOK, data)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("Expected Content-Type to be application/json, got %s", w.Header().Get("Content-Type"))
	}

	var response map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response["key"] != "value" {
		t.Errorf("Expected response to contain key=value, got %v", response)
	}
}

func TestWriteSuccess(t *testing.T) {
	w := httptest.NewRecorder()
	data := "test data"

	WriteSuccess(w, http.StatusCreated, data, "test message")

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Error("Expected Success to be true")
	}

	if response.Data != "test data" {
		t.Errorf("Expected Data to be 'test data', got %v", response.Data)
	}

	if response.Message != "test message" {
		t.Errorf("Expected Message to be 'test message', got %v", response.Message)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	err := &testError{message: "test error"}

	WriteError(w, http.StatusBadRequest, "test error message", err)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Success {
		t.Error("Expected Success to be false")
	}

	expectedError := "test error message: test error"
	if response.Error != expectedError {
		t.Errorf("Expected Error to be '%s', got %v", expectedError, response.Error)
	}
}

func TestPaginationMeta(t *testing.T) {
	meta := PaginationMeta(2, 10, 25)

	if meta.Page != 2 {
		t.Errorf("Expected Page to be 2, got %d", meta.Page)
	}

	if meta.Limit != 10 {
		t.Errorf("Expected Limit to be 10, got %d", meta.Limit)
	}

	if meta.Total != 25 {
		t.Errorf("Expected Total to be 25, got %d", meta.Total)
	}

	if meta.TotalPages != 3 {
		t.Errorf("Expected TotalPages to be 3, got %d", meta.TotalPages)
	}
}

func TestWriteSuccessWithMeta(t *testing.T) {
	w := httptest.NewRecorder()
	data := "test data"
	meta := &Meta{Page: 1, Limit: 10, Total: 100}

	WriteSuccessWithMeta(w, http.StatusOK, data, "test message", meta)

	var response APIResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta == nil {
		t.Error("Expected Meta to be set")
	}

	if response.Meta.Page != 1 {
		t.Errorf("Expected Meta.Page to be 1, got %d", response.Meta.Page)
	}
}

// Helper type for testing
type testError struct {
	message string
}

func (e *testError) Error() string {
	return e.message
}
