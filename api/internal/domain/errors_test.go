package domain

import (
	"errors"
	"testing"
)

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}

	got := err.Error()
	want := "invalid email format"

	if got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestValidationError_Unwrap(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid",
	}

	unwrapped := err.Unwrap()
	if unwrapped != ErrValidation {
		t.Errorf("Unwrap() = %v, want %v", unwrapped, ErrValidation)
	}
}

func TestValidationError_ErrorsIs(t *testing.T) {
	err := &ValidationError{
		Field:   "email",
		Message: "invalid",
	}

	if !errors.Is(err, ErrValidation) {
		t.Error("expected errors.Is(ValidationError, ErrValidation) to be true")
	}
}

func TestValidationError_ErrorsAs(t *testing.T) {
	originalErr := &ValidationError{
		Field:   "email",
		Message: "invalid email format",
	}

	// Wrap the error
	wrappedErr := errors.New("wrapped: " + originalErr.Error())
	_ = wrappedErr // We'll test the direct case

	var validationErr *ValidationError
	if !errors.As(originalErr, &validationErr) {
		t.Fatal("expected errors.As to succeed")
	}

	if validationErr.Field != "email" {
		t.Errorf("Field = %q, want %q", validationErr.Field, "email")
	}
	if validationErr.Message != "invalid email format" {
		t.Errorf("Message = %q, want %q", validationErr.Message, "invalid email format")
	}
}

func TestNewValidationError(t *testing.T) {
	tests := []struct {
		name    string
		field   string
		message string
	}{
		{
			name:    "email field",
			field:   "email",
			message: "invalid email format",
		},
		{
			name:    "required field",
			field:   "name",
			message: "field is required",
		},
		{
			name:    "empty field name",
			field:   "",
			message: "general validation error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := NewValidationError(tt.field, tt.message)

			if err.Field != tt.field {
				t.Errorf("Field = %q, want %q", err.Field, tt.field)
			}
			if err.Message != tt.message {
				t.Errorf("Message = %q, want %q", err.Message, tt.message)
			}
			if err.Error() != tt.message {
				t.Errorf("Error() = %q, want %q", err.Error(), tt.message)
			}
			if !errors.Is(err, ErrValidation) {
				t.Error("expected errors.Is(err, ErrValidation) to be true")
			}
		})
	}
}

func TestDomainErrors_AreDistinct(t *testing.T) {
	domainErrors := []error{
		ErrNotFound,
		ErrUnauthorized,
		ErrForbidden,
		ErrValidation,
		ErrCVLimitReached,
		ErrStorageLimitReached,
		ErrInvalidFileType,
		ErrFileTooLarge,
		ErrConflict,
		ErrInternal,
	}

	// Each error should only match itself
	for i, err1 := range domainErrors {
		for j, err2 := range domainErrors {
			if i == j {
				if !errors.Is(err1, err2) {
					t.Errorf("expected errors.Is(%v, %v) to be true", err1, err2)
				}
			} else {
				if errors.Is(err1, err2) {
					t.Errorf("expected errors.Is(%v, %v) to be false", err1, err2)
				}
			}
		}
	}
}
