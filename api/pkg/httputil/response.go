package httputil

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/ellebam/hireme/api/internal/domain"
)

// Response is the standard API response wrapper
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorBody  `json:"error,omitempty"`
}

// ErrorBody contains error details
type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Field   string `json:"field,omitempty"`
}

// JSON writes a JSON response with the given status code
func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{Data: data}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Log error but can't do much at this point
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// Error writes an error response with the given status code
func Error(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Error: &ErrorBody{
			Code:    http.StatusText(status),
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// ErrorWithCode writes an error response with a custom error code
func ErrorWithCode(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	response := Response{
		Error: &ErrorBody{
			Code:    code,
			Message: message,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// ValidationError writes a validation error response
func ValidationError(w http.ResponseWriter, field, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)

	response := Response{
		Error: &ErrorBody{
			Code:    "validation_error",
			Message: message,
			Field:   field,
		},
	}
	json.NewEncoder(w).Encode(response)
}

// HandleError maps domain errors to HTTP responses
func HandleError(w http.ResponseWriter, err error) {
	// Check for validation errors with field info
	var validationErr *domain.ValidationError
	if errors.As(err, &validationErr) {
		ValidationError(w, validationErr.Field, validationErr.Message)
		return
	}

	// Map domain errors to HTTP status codes
	switch {
	case errors.Is(err, domain.ErrNotFound):
		Error(w, http.StatusNotFound, "resource not found")

	case errors.Is(err, domain.ErrUnauthorized):
		Error(w, http.StatusUnauthorized, "unauthorized")

	case errors.Is(err, domain.ErrForbidden):
		Error(w, http.StatusForbidden, "forbidden")

	case errors.Is(err, domain.ErrValidation):
		Error(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, domain.ErrCVLimitReached):
		ErrorWithCode(w, http.StatusForbidden, "cv_limit_reached", "CV limit reached for your plan")

	case errors.Is(err, domain.ErrStorageLimitReached):
		ErrorWithCode(w, http.StatusForbidden, "storage_limit_reached", "storage limit reached for your plan")

	case errors.Is(err, domain.ErrInvalidFileType):
		ErrorWithCode(w, http.StatusBadRequest, "invalid_file_type", "file type not allowed")

	case errors.Is(err, domain.ErrFileTooLarge):
		ErrorWithCode(w, http.StatusBadRequest, "file_too_large", "file exceeds size limit")

	case errors.Is(err, domain.ErrConflict):
		Error(w, http.StatusConflict, "resource already exists")

	default:
		Error(w, http.StatusInternalServerError, "internal server error")
	}
}

// DecodeJSON decodes JSON from request body into the given struct
func DecodeJSON(r *http.Request, v interface{}) error {
	if r.Body == nil {
		return errors.New("empty request body")
	}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}

// NoContent writes a 204 No Content response
func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// Created writes a 201 Created response with the created resource
func Created(w http.ResponseWriter, data interface{}) {
	JSON(w, http.StatusCreated, data)
}
