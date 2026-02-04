package validator

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/ellebam/hireme/api/internal/domain"
)

func TestCVValidator_Validate_Valid(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	// Minimal valid CV JSON according to schema
	validCV := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": []
	}`

	err = validator.Validate(json.RawMessage(validCV))
	if err != nil {
		t.Errorf("expected valid CV to pass validation, got error: %v", err)
	}
}

func TestCVValidator_Validate_ValidWithSections(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	// Valid CV with a personal section
	validCV := `{
		"schemaVersion": "1.0.0",
		"templateId": "modern",
		"locale": "en",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440000",
				"type": "personal",
				"order": 0,
				"content": {
					"firstName": "John",
					"lastName": "Doe"
				}
			}
		]
	}`

	err = validator.Validate(json.RawMessage(validCV))
	if err != nil {
		t.Errorf("expected valid CV with sections to pass validation, got error: %v", err)
	}
}

func TestCVValidator_Validate_MissingRequired(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		content string
	}{
		{
			name: "missing templateId",
			content: `{
				"schemaVersion": "1.0.0",
				"sections": []
			}`,
		},
		{
			name: "missing schemaVersion",
			content: `{
				"templateId": "classic",
				"sections": []
			}`,
		},
		{
			name: "missing sections",
			content: `{
				"schemaVersion": "1.0.0",
				"templateId": "classic"
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(json.RawMessage(tt.content))
			if err == nil {
				t.Error("expected validation error for missing required field, got nil")
			}

			var validationErr *domain.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
		})
	}
}

func TestCVValidator_Validate_InvalidJSON(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "malformed JSON",
			content: `{broken`,
		},
		{
			name:    "incomplete JSON",
			content: `{"schemaVersion": `,
		},
		{
			name:    "empty content",
			content: ``,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(json.RawMessage(tt.content))
			if err == nil {
				t.Error("expected error for invalid JSON, got nil")
			}

			var validationErr *domain.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
		})
	}
}

func TestCVValidator_Validate_InvalidTemplateId(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	invalidCV := `{
		"schemaVersion": "1.0.0",
		"templateId": "unknown",
		"sections": []
	}`

	err = validator.Validate(json.RawMessage(invalidCV))
	if err == nil {
		t.Error("expected validation error for invalid templateId, got nil")
	}
}

func TestCVValidator_Validate_InvalidLocale(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	invalidCV := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"locale": "fr",
		"sections": []
	}`

	err = validator.Validate(json.RawMessage(invalidCV))
	if err == nil {
		t.Error("expected validation error for invalid locale, got nil")
	}
}

func TestValidateTemplateID_Valid(t *testing.T) {
	validTemplates := []string{"classic", "modern", "minimal"}

	for _, template := range validTemplates {
		t.Run(template, func(t *testing.T) {
			err := ValidateTemplateID(template)
			if err != nil {
				t.Errorf("expected no error for valid template %q, got %v", template, err)
			}
		})
	}
}

func TestValidateTemplateID_Invalid(t *testing.T) {
	invalidTemplates := []string{"unknown", "fancy", "", "CLASSIC", "Modern"}

	for _, template := range invalidTemplates {
		t.Run(template, func(t *testing.T) {
			err := ValidateTemplateID(template)
			if err == nil {
				t.Errorf("expected error for invalid template %q, got nil", template)
			}

			var validationErr *domain.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
			if validationErr.Field != "templateId" {
				t.Errorf("expected field 'templateId', got %q", validationErr.Field)
			}
		})
	}
}

func TestValidateLocale_Valid(t *testing.T) {
	validLocales := []string{"en", "de"}

	for _, locale := range validLocales {
		t.Run(locale, func(t *testing.T) {
			err := ValidateLocale(locale)
			if err != nil {
				t.Errorf("expected no error for valid locale %q, got %v", locale, err)
			}
		})
	}
}

func TestValidateLocale_Invalid(t *testing.T) {
	invalidLocales := []string{"fr", "es", "", "EN", "De", "english"}

	for _, locale := range invalidLocales {
		t.Run(locale, func(t *testing.T) {
			err := ValidateLocale(locale)
			if err == nil {
				t.Errorf("expected error for invalid locale %q, got nil", locale)
			}

			var validationErr *domain.ValidationError
			if !errors.As(err, &validationErr) {
				t.Errorf("expected ValidationError, got %T", err)
			}
			if validationErr.Field != "locale" {
				t.Errorf("expected field 'locale', got %q", validationErr.Field)
			}
		})
	}
}
