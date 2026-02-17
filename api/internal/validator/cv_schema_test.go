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

func TestCVValidator_Validate_AllValidTemplateIds(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	validTemplates := []string{"classic", "modern", "visionary"}

	for _, tmpl := range validTemplates {
		t.Run(tmpl, func(t *testing.T) {
			cv := `{
				"schemaVersion": "1.0.0",
				"templateId": "` + tmpl + `",
				"sections": []
			}`
			err := validator.Validate(json.RawMessage(cv))
			if err != nil {
				t.Errorf("expected templateId %q to pass schema validation, got: %v", tmpl, err)
			}
		})
	}
}

func TestCVValidator_Validate_InvalidTemplateId(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	invalidTemplates := []string{"unknown", "minimal", "blank", "CLASSIC"}

	for _, tmpl := range invalidTemplates {
		t.Run(tmpl, func(t *testing.T) {
			cv := `{
				"schemaVersion": "1.0.0",
				"templateId": "` + tmpl + `",
				"sections": []
			}`
			err := validator.Validate(json.RawMessage(cv))
			if err == nil {
				t.Errorf("expected templateId %q to fail schema validation, got nil", tmpl)
			}
		})
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
	validTemplates := []string{"classic", "modern", "visionary"}

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
	invalidTemplates := []string{"unknown", "fancy", "", "CLASSIC", "Modern", "minimal"}

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

func TestCVValidator_Validate_EmptyEntriesSection(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	// A certifications section with empty entries should validate
	// (this was previously broken due to oneOf collision)
	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "certifications",
				"order": 0,
				"content": { "entries": [] }
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err != nil {
		t.Errorf("expected empty certifications entries to pass validation, got: %v", err)
	}
}

func TestCVValidator_Validate_MultipleEmptyEntriesSections(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	// Multiple entry-based sections all with empty entries
	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "certifications",
				"order": 0,
				"content": { "entries": [] }
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440002",
				"type": "projects",
				"order": 1,
				"content": { "entries": [] }
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440003",
				"type": "awards",
				"order": 2,
				"content": { "entries": [] }
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err != nil {
		t.Errorf("expected multiple empty-entries sections to pass validation, got: %v", err)
	}
}

func TestCVValidator_Validate_PopulatedCertificationsSection(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "certifications",
				"order": 0,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440010",
							"name": "AWS Solutions Architect",
							"issuer": "Amazon Web Services",
							"date": "2023-06"
						}
					]
				}
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err != nil {
		t.Errorf("expected populated certifications to pass validation, got: %v", err)
	}
}

func TestCVValidator_Validate_PopulatedLanguagesSection(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "languages",
				"order": 0,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440020",
							"language": "German",
							"proficiency": "native"
						},
						{
							"id": "550e8400-e29b-41d4-a716-446655440021",
							"language": "English",
							"proficiency": "fluent",
							"certification": "TOEFL 120"
						}
					]
				}
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err != nil {
		t.Errorf("expected populated languages to pass validation, got: %v", err)
	}
}

func TestCVValidator_Validate_LanguageEntryWithoutId(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "languages",
				"order": 0,
				"content": {
					"entries": [
						{
							"language": "German",
							"proficiency": "native"
						}
					]
				}
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err == nil {
		t.Error("expected validation error for language entry without id, got nil")
	}
}

func TestCVValidator_Validate_LanguageEntryExtraProperty(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "languages",
				"order": 0,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440020",
							"language": "German",
							"proficiency": "native",
							"extra": "should-fail"
						}
					]
				}
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err == nil {
		t.Error("expected validation error for language entry with extra property, got nil")
	}
}

func TestCVValidator_Validate_FullRealisticCV(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "modern",
		"locale": "en",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "personal",
				"order": 0,
				"content": {
					"firstName": "John",
					"lastName": "Doe",
					"jobTitle": "Software Engineer",
					"email": "john@example.com"
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440002",
				"type": "summary",
				"order": 1,
				"content": {
					"text": "Experienced engineer."
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440003",
				"type": "experience",
				"order": 2,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440010",
							"company": "TechCorp",
							"position": "Senior Engineer",
							"startDate": "2020-01"
						}
					]
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440004",
				"type": "education",
				"order": 3,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440011",
							"institution": "MIT",
							"degree": "M.Sc."
						}
					]
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440005",
				"type": "skills",
				"order": 4,
				"content": {
					"categories": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440012",
							"name": "Languages",
							"skills": [{"name": "Go", "level": "expert"}]
						}
					]
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440006",
				"type": "languages",
				"order": 5,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440013",
							"language": "English",
							"proficiency": "native"
						}
					]
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440007",
				"type": "certifications",
				"order": 6,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440014",
							"name": "AWS SA",
							"issuer": "Amazon"
						}
					]
				}
			},
			{
				"id": "550e8400-e29b-41d4-a716-446655440008",
				"type": "projects",
				"order": 7,
				"content": {
					"entries": [
						{
							"id": "550e8400-e29b-41d4-a716-446655440015",
							"name": "CLI Tool",
							"description": "A command-line tool."
						}
					]
				}
			}
		],
		"styling": {
			"primaryColor": "#2563eb",
			"fontFamily": "inter",
			"fontSize": "medium"
		}
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err != nil {
		t.Errorf("expected full realistic CV to pass validation, got: %v", err)
	}
}

func TestCVValidator_Validate_WrongContentForSectionType(t *testing.T) {
	validator, err := NewCVValidator()
	if err != nil {
		t.Fatalf("failed to create validator: %v", err)
	}

	// Certifications section with summary content (has "text" instead of "entries")
	cv := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"sections": [
			{
				"id": "550e8400-e29b-41d4-a716-446655440001",
				"type": "certifications",
				"order": 0,
				"content": { "text": "This is summary content, not certifications" }
			}
		]
	}`

	err = validator.Validate(json.RawMessage(cv))
	if err == nil {
		t.Error("expected validation error for wrong content type in certifications section, got nil")
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
