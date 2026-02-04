package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestJSON_Success(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := map[string]string{"name": "test"}
	JSON(recorder, http.StatusOK, data)

	result := recorder.Result()
	defer result.Body.Close()

	// Check status code
	if result.StatusCode != http.StatusOK {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusOK)
	}

	// Check content type
	contentType := result.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	// Parse response body
	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Check data is wrapped
	if response.Data == nil {
		t.Error("expected Data field to be set")
	}
	if response.Error != nil {
		t.Error("expected Error field to be nil")
	}
}

func TestJSON_WithStruct(t *testing.T) {
	recorder := httptest.NewRecorder()

	type TestData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := TestData{ID: 1, Name: "test"}
	JSON(recorder, http.StatusCreated, data)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusCreated {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusCreated)
	}

	body, _ := io.ReadAll(result.Body)
	if !strings.Contains(string(body), `"id":1`) {
		t.Error("expected response to contain id field")
	}
	if !strings.Contains(string(body), `"name":"test"`) {
		t.Error("expected response to contain name field")
	}
}

func TestError_NotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	Error(recorder, http.StatusNotFound, "resource not found")

	result := recorder.Result()
	defer result.Body.Close()

	// Check status code
	if result.StatusCode != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusNotFound)
	}

	// Check content type
	contentType := result.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Content-Type = %q, want %q", contentType, "application/json")
	}

	// Parse response body
	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error == nil {
		t.Fatal("expected Error field to be set")
	}
	if response.Error.Code != "Not Found" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "Not Found")
	}
	if response.Error.Message != "resource not found" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "resource not found")
	}
	if response.Data != nil {
		t.Error("expected Data field to be nil")
	}
}

func TestError_BadRequest(t *testing.T) {
	recorder := httptest.NewRecorder()

	Error(recorder, http.StatusBadRequest, "invalid input")

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusBadRequest)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "Bad Request" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "Bad Request")
	}
}

func TestHandleError_ErrNotFound(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, domain.ErrNotFound)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusNotFound {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusNotFound)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error == nil {
		t.Fatal("expected Error field to be set")
	}
}

func TestHandleError_ErrForbidden(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, domain.ErrForbidden)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusForbidden)
	}
}

func TestHandleError_ErrUnauthorized(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, domain.ErrUnauthorized)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusUnauthorized {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusUnauthorized)
	}
}

func TestHandleError_ValidationError(t *testing.T) {
	recorder := httptest.NewRecorder()

	validationErr := domain.NewValidationError("email", "invalid email format")
	HandleError(recorder, validationErr)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusBadRequest)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error == nil {
		t.Fatal("expected Error field to be set")
	}
	if response.Error.Field != "email" {
		t.Errorf("Error.Field = %q, want %q", response.Error.Field, "email")
	}
	if response.Error.Code != "validation_error" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "validation_error")
	}
}

func TestHandleError_ErrCVLimitReached(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, domain.ErrCVLimitReached)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusForbidden)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "cv_limit_reached" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "cv_limit_reached")
	}
}

func TestHandleError_ErrStorageLimitReached(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, domain.ErrStorageLimitReached)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusForbidden)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "storage_limit_reached" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "storage_limit_reached")
	}
}

func TestHandleError_UnknownError(t *testing.T) {
	recorder := httptest.NewRecorder()

	HandleError(recorder, errors.New("some random error"))

	result := recorder.Result()
	defer result.Body.Close()

	// Unknown errors should return 500
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusInternalServerError)
	}
}

func TestDecodeJSON_Valid(t *testing.T) {
	type TestInput struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	body := `{"name": "John", "email": "john@example.com"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var input TestInput
	err := DecodeJSON(req, &input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if input.Name != "John" {
		t.Errorf("Name = %q, want %q", input.Name, "John")
	}
	if input.Email != "john@example.com" {
		t.Errorf("Email = %q, want %q", input.Email, "john@example.com")
	}
}

func TestDecodeJSON_Invalid(t *testing.T) {
	type TestInput struct {
		Name string `json:"name"`
	}

	tests := []struct {
		name string
		body string
	}{
		{
			name: "malformed JSON",
			body: `{broken`,
		},
		{
			name: "incomplete JSON",
			body: `{"name": `,
		},
		{
			name: "wrong type",
			body: `{"name": 123}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))

			var input TestInput
			err := DecodeJSON(req, &input)
			if err == nil {
				t.Error("expected error for invalid JSON, got nil")
			}
		})
	}
}

func TestDecodeJSON_EmptyBody(t *testing.T) {
	type TestInput struct {
		Name string `json:"name"`
	}

	req := httptest.NewRequest(http.MethodPost, "/", nil)

	var input TestInput
	err := DecodeJSON(req, &input)
	if err == nil {
		t.Error("expected error for empty body, got nil")
	}
}

func TestDecodeJSON_UnknownFields(t *testing.T) {
	type TestInput struct {
		Name string `json:"name"`
	}

	// JSON with unknown field should fail because DisallowUnknownFields is set
	body := `{"name": "John", "unknown": "field"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))

	var input TestInput
	err := DecodeJSON(req, &input)
	if err == nil {
		t.Error("expected error for unknown fields, got nil")
	}
}

func TestNoContent(t *testing.T) {
	recorder := httptest.NewRecorder()

	NoContent(recorder)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusNoContent {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusNoContent)
	}

	body, _ := io.ReadAll(result.Body)
	if len(body) != 0 {
		t.Errorf("expected empty body, got %q", string(body))
	}
}

func TestCreated(t *testing.T) {
	recorder := httptest.NewRecorder()

	data := map[string]string{"id": "123"}
	Created(recorder, data)

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusCreated {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusCreated)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Data == nil {
		t.Error("expected Data field to be set")
	}
}

func TestValidationError_Response(t *testing.T) {
	recorder := httptest.NewRecorder()

	ValidationError(recorder, "email", "invalid email format")

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusBadRequest {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusBadRequest)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error == nil {
		t.Fatal("expected Error field to be set")
	}
	if response.Error.Code != "validation_error" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "validation_error")
	}
	if response.Error.Field != "email" {
		t.Errorf("Error.Field = %q, want %q", response.Error.Field, "email")
	}
	if response.Error.Message != "invalid email format" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "invalid email format")
	}
}

func TestErrorWithCode(t *testing.T) {
	recorder := httptest.NewRecorder()

	ErrorWithCode(recorder, http.StatusForbidden, "custom_code", "custom message")

	result := recorder.Result()
	defer result.Body.Close()

	if result.StatusCode != http.StatusForbidden {
		t.Errorf("status code = %d, want %d", result.StatusCode, http.StatusForbidden)
	}

	var response Response
	if err := json.NewDecoder(result.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Error.Code != "custom_code" {
		t.Errorf("Error.Code = %q, want %q", response.Error.Code, "custom_code")
	}
	if response.Error.Message != "custom message" {
		t.Errorf("Error.Message = %q, want %q", response.Error.Message, "custom message")
	}
}

