package utils

import (
	"encoding/json"
	"net/http"
	"time"
)

// APIResponse represents a standardized API response
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// Meta contains metadata about the response
type Meta struct {
	Page       int `json:"page,omitempty"`
	Limit      int `json:"limit,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"total_pages,omitempty"`
}

// PaginationMeta creates pagination metadata
func PaginationMeta(page, limit, total int) *Meta {
	totalPages := (total + limit - 1) / limit
	return &Meta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}
}

// SuccessResponse creates a successful response
func SuccessResponse(data interface{}, message string) *APIResponse {
	return &APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now(),
	}
}

// SuccessResponseWithMeta creates a successful response with metadata
func SuccessResponseWithMeta(data interface{}, message string, meta *Meta) *APIResponse {
	return &APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

// ErrorResponse creates an error response
func ErrorResponse(message string, err error) *APIResponse {
	response := &APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now(),
	}

	if err != nil {
		// Add error details to the message if provided
		if message == "" {
			response.Error = err.Error()
		} else {
			response.Error = message + ": " + err.Error()
		}
	}

	return response
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(errors map[string]string) *APIResponse {
	return &APIResponse{
		Success:   false,
		Error:     "Validation failed",
		Data:      errors,
		Timestamp: time.Now(),
	}
}

// WriteJSON writes a JSON response with proper headers
func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}
	}
}

// WriteSuccess writes a successful JSON response
func WriteSuccess(w http.ResponseWriter, statusCode int, data interface{}, message string) {
	response := SuccessResponse(data, message)
	WriteJSON(w, statusCode, response)
}

// WriteSuccessWithMeta writes a successful JSON response with metadata
func WriteSuccessWithMeta(w http.ResponseWriter, statusCode int, data interface{}, message string, meta *Meta) {
	response := SuccessResponseWithMeta(data, message, meta)
	WriteJSON(w, statusCode, response)
}

// WriteError writes an error JSON response
func WriteError(w http.ResponseWriter, statusCode int, message string, err error) {
	response := ErrorResponse(message, err)
	WriteJSON(w, statusCode, response)
}

// WriteValidationError writes a validation error JSON response
func WriteValidationError(w http.ResponseWriter, errors map[string]string) {
	response := ValidationErrorResponse(errors)
	WriteJSON(w, http.StatusBadRequest, response)
}

// Common HTTP status responses
func WriteCreated(w http.ResponseWriter, data interface{}, message string) {
	WriteSuccess(w, http.StatusCreated, data, message)
}

func WriteOK(w http.ResponseWriter, data interface{}, message string) {
	WriteSuccess(w, http.StatusOK, data, message)
}

func WriteNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func WriteBadRequest(w http.ResponseWriter, message string, err error) {
	WriteError(w, http.StatusBadRequest, message, err)
}

func WriteUnauthorized(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusUnauthorized, message, nil)
}

func WriteForbidden(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusForbidden, message, nil)
}

func WriteNotFound(w http.ResponseWriter, message string) {
	WriteError(w, http.StatusNotFound, message, nil)
}

func WriteInternalServerError(w http.ResponseWriter, message string, err error) {
	WriteError(w, http.StatusInternalServerError, message, err)
}

func WriteMethodNotAllowed(w http.ResponseWriter, allowedMethods string) {
	w.Header().Set("Allow", allowedMethods)
	WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
}
