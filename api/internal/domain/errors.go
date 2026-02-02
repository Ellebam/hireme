package domain

import "errors"

// Domain errors - these are mapped to HTTP status codes by handlers
var (
	// ErrNotFound indicates the requested resource was not found
	ErrNotFound = errors.New("not found")

	// ErrUnauthorized indicates the request lacks valid authentication
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden indicates the user doesn't have permission for this action
	ErrForbidden = errors.New("forbidden")

	// ErrValidation indicates the input failed validation
	ErrValidation = errors.New("validation error")

	// ErrCVLimitReached indicates the user has reached their CV limit
	ErrCVLimitReached = errors.New("cv limit reached")

	// ErrStorageLimitReached indicates the user has reached their storage limit
	ErrStorageLimitReached = errors.New("storage limit reached")

	// ErrInvalidFileType indicates an unsupported file type was uploaded
	ErrInvalidFileType = errors.New("invalid file type")

	// ErrFileTooLarge indicates the uploaded file exceeds size limits
	ErrFileTooLarge = errors.New("file too large")

	// ErrConflict indicates a conflict with existing data
	ErrConflict = errors.New("conflict")

	// ErrInternal indicates an unexpected internal error
	ErrInternal = errors.New("internal error")
)

// ValidationError wraps ErrValidation with specific field information
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

// NewValidationError creates a new validation error for a specific field
func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}
