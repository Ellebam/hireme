package domain

import (
	"encoding/json"
	"testing"
)

func TestCV_ParseContent_Valid(t *testing.T) {
	validJSON := `{
		"schemaVersion": "1.0.0",
		"templateId": "classic",
		"locale": "en",
		"sections": []
	}`

	cv := &CV{
		Content: json.RawMessage(validJSON),
	}

	content, err := cv.ParseContent()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if content.SchemaVersion != "1.0.0" {
		t.Errorf("expected schemaVersion '1.0.0', got %q", content.SchemaVersion)
	}
	if content.TemplateID != "classic" {
		t.Errorf("expected templateId 'classic', got %q", content.TemplateID)
	}
	if content.Locale != "en" {
		t.Errorf("expected locale 'en', got %q", content.Locale)
	}
}

func TestCV_ParseContent_Invalid(t *testing.T) {
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
			name:    "invalid syntax",
			content: `{schemaVersion: "1.0.0"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cv := &CV{
				Content: json.RawMessage(tt.content),
			}

			_, err := cv.ParseContent()
			if err == nil {
				t.Error("expected error for invalid JSON, got nil")
			}
		})
	}
}

func TestCV_SetContent(t *testing.T) {
	cv := &CV{}

	content := &CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "modern",
		Locale:        "de",
		Title:         "My CV",
		Sections:      []CVSection{},
	}

	err := cv.SetContent(content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the content was set correctly by parsing it back
	parsed, err := cv.ParseContent()
	if err != nil {
		t.Fatalf("failed to parse set content: %v", err)
	}

	if parsed.SchemaVersion != content.SchemaVersion {
		t.Errorf("schemaVersion: got %q, want %q", parsed.SchemaVersion, content.SchemaVersion)
	}
	if parsed.TemplateID != content.TemplateID {
		t.Errorf("templateId: got %q, want %q", parsed.TemplateID, content.TemplateID)
	}
	if parsed.Locale != content.Locale {
		t.Errorf("locale: got %q, want %q", parsed.Locale, content.Locale)
	}
	if parsed.Title != content.Title {
		t.Errorf("title: got %q, want %q", parsed.Title, content.Title)
	}
}

func TestCV_SetContent_WithSections(t *testing.T) {
	cv := &CV{}

	content := &CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections: []CVSection{
			{
				ID:      "section-1",
				Type:    SectionTypePersonal,
				Order:   0,
				Visible: true,
				Content: json.RawMessage(`{"firstName": "John"}`),
			},
		},
	}

	err := cv.SetContent(content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsed, err := cv.ParseContent()
	if err != nil {
		t.Fatalf("failed to parse set content: %v", err)
	}

	if len(parsed.Sections) != 1 {
		t.Fatalf("expected 1 section, got %d", len(parsed.Sections))
	}

	section := parsed.Sections[0]
	if section.ID != "section-1" {
		t.Errorf("section ID: got %q, want %q", section.ID, "section-1")
	}
	if section.Type != SectionTypePersonal {
		t.Errorf("section Type: got %q, want %q", section.Type, SectionTypePersonal)
	}
}

func TestCV_SetContent_WithStyling(t *testing.T) {
	cv := &CV{}
	showIcons := true

	content := &CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "visionary",
		Sections:      []CVSection{},
		Styling: &CVStyling{
			PrimaryColor: "#2563eb",
			FontFamily:   "inter",
			ShowIcons:    &showIcons,
		},
	}

	err := cv.SetContent(content)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	parsed, err := cv.ParseContent()
	if err != nil {
		t.Fatalf("failed to parse set content: %v", err)
	}

	if parsed.Styling == nil {
		t.Fatal("expected styling to be set, got nil")
	}
	if parsed.Styling.PrimaryColor != "#2563eb" {
		t.Errorf("primaryColor: got %q, want %q", parsed.Styling.PrimaryColor, "#2563eb")
	}
	if parsed.Styling.FontFamily != "inter" {
		t.Errorf("fontFamily: got %q, want %q", parsed.Styling.FontFamily, "inter")
	}
}
