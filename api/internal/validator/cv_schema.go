package validator

import (
	"embed"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v5"

	"github.com/ellebam/hireme/api/internal/domain"
)

//go:embed schema/cv-schema.json
var schemaFS embed.FS

// CVValidator validates CV content against the JSON schema
type CVValidator struct {
	schema *jsonschema.Schema
}

// NewCVValidator creates a new CVValidator
func NewCVValidator() (*CVValidator, error) {
	// Read embedded schema
	schemaData, err := schemaFS.ReadFile("schema/cv-schema.json")
	if err != nil {
		return nil, fmt.Errorf("reading embedded schema: %w", err)
	}

	// Compile schema
	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("cv-schema.json", strings.NewReader(string(schemaData))); err != nil {
		return nil, fmt.Errorf("adding schema resource: %w", err)
	}

	schema, err := compiler.Compile("cv-schema.json")
	if err != nil {
		return nil, fmt.Errorf("compiling schema: %w", err)
	}

	return &CVValidator{schema: schema}, nil
}

// Validate validates CV content against the JSON schema
func (v *CVValidator) Validate(content json.RawMessage) error {
	var data interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return domain.NewValidationError("content", fmt.Sprintf("invalid JSON: %s", err.Error()))
	}

	if err := v.schema.Validate(data); err != nil {
		// Extract validation error details
		if validationErr, ok := err.(*jsonschema.ValidationError); ok {
			return domain.NewValidationError("content", formatValidationError(validationErr))
		}
		return domain.NewValidationError("content", err.Error())
	}

	return nil
}

// formatValidationError formats a JSON Schema validation error
func formatValidationError(err *jsonschema.ValidationError) string {
	if len(err.Causes) > 0 {
		// Return the first specific error
		return formatValidationError(err.Causes[0])
	}

	// Format the error message with path
	if err.InstanceLocation != "" {
		return fmt.Sprintf("%s: %s", err.InstanceLocation, err.Message)
	}
	return err.Message
}

// ValidateTemplateID validates that a template ID is valid
func ValidateTemplateID(templateID string) error {
	switch templateID {
	case domain.TemplateClassic, domain.TemplateModern, domain.TemplateMinimal:
		return nil
	default:
		return domain.NewValidationError("templateId", "invalid template ID")
	}
}

// ValidateLocale validates that a locale is valid
func ValidateLocale(locale string) error {
	switch locale {
	case domain.LocaleEN, domain.LocaleDE:
		return nil
	default:
		return domain.NewValidationError("locale", "invalid locale")
	}
}
