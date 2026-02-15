package export

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/ellebam/hireme/api/internal/domain"
)

func docxContent(sections ...domain.CVSection) domain.CVContent {
	return domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections:      sections,
	}
}

func mustJSON(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

func readDocumentXML(t *testing.T, docxBytes []byte) string {
	t.Helper()
	zr, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		t.Fatalf("opening zip: %v", err)
	}
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				t.Fatalf("opening document.xml: %v", err)
			}
			defer func() {
				if err := rc.Close(); err != nil {
					t.Fatalf("closing document.xml: %v", err)
				}
			}()
			var buf bytes.Buffer
			if _, err := buf.ReadFrom(rc); err != nil {
				t.Fatalf("reading document.xml: %v", err)
			}
			return buf.String()
		}
	}
	t.Fatal("word/document.xml not found in DOCX")
	return ""
}

func TestGenerate_ValidDOCX(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "personal", Order: 0, Visible: true,
		Content: mustJSON(domain.PersonalContent{FirstName: "Jane", LastName: "Doe"}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check ZIP magic bytes (PK\x03\x04)
	if len(result) < 4 || result[0] != 0x50 || result[1] != 0x4B || result[2] != 0x03 || result[3] != 0x04 {
		t.Fatalf("output does not start with ZIP magic bytes, got %x", result[:4])
	}

	// Verify it's a valid ZIP with word/document.xml
	zr, err := zip.NewReader(bytes.NewReader(result), int64(len(result)))
	if err != nil {
		t.Fatalf("invalid ZIP: %v", err)
	}
	found := false
	for _, f := range zr.File {
		if f.Name == "word/document.xml" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("word/document.xml not found in DOCX")
	}
}

func TestGenerate_PersonalSection(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "personal", Order: 0, Visible: true,
		Content: mustJSON(domain.PersonalContent{
			FirstName: "Alice",
			LastName:  "Smith",
			JobTitle:  "Software Engineer",
			Email:     "alice@example.com",
			Phone:     "+1234567890",
			Location:  "Berlin",
		}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	for _, expected := range []string{"Alice Smith", "Software Engineer", "alice@example.com", "+1234567890", "Berlin"} {
		if !strings.Contains(xml, expected) {
			t.Errorf("expected %q in document.xml, not found", expected)
		}
	}
}

func TestGenerate_AllSectionTypes(t *testing.T) {
	gen := NewGodocxGenerator()
	endDate := "2023-06"
	content := docxContent(
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Bob", LastName: "Jones"}),
		},
		domain.CVSection{
			ID: "2", Type: "summary", Order: 1, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "Experienced developer"}),
		},
		domain.CVSection{
			ID: "3", Type: "experience", Order: 2, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
				ID: "e1", Company: "Acme", Position: "Lead", StartDate: "2020-01", EndDate: &endDate,
				Description: "Led team", Highlights: []string{"Grew revenue 2x", "Hired 5 engineers"},
			}}}),
		},
		domain.CVSection{
			ID: "4", Type: "education", Order: 3, Visible: true,
			Content: mustJSON(domain.EducationContent{Entries: []domain.EducationEntry{{
				ID: "ed1", Institution: "MIT", Degree: "BSc", Field: "CS",
			}}}),
		},
		domain.CVSection{
			ID: "5", Type: "skills", Order: 4, Visible: true,
			Content: mustJSON(domain.SkillsContent{Categories: []domain.SkillCategory{{
				ID: "s1", Name: "Programming", Skills: []domain.Skill{{Name: "Go"}, {Name: "TypeScript"}},
			}}}),
		},
		domain.CVSection{
			ID: "6", Type: "languages", Order: 5, Visible: true,
			Content: mustJSON(domain.LanguagesContent{Entries: []domain.LanguageEntry{{
				Language: "English", Proficiency: "Native",
			}}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	for _, expected := range []string{
		"Bob Jones", "Experienced developer",
		"Lead", "Acme", "Led team",
		"Grew revenue 2x", "Hired 5 engineers",
		"MIT", "BSc", "CS",
		"Programming", "Go", "TypeScript",
		"English", "Native",
	} {
		if !strings.Contains(xml, expected) {
			t.Errorf("expected %q in document.xml, not found", expected)
		}
	}
}

func TestGenerate_ExperienceHighlights(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "experience", Order: 0, Visible: true,
		Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
			ID: "e1", Position: "Dev", Highlights: []string{"Led team of 5", "Increased revenue"},
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "• Led team of 5") {
		t.Error("expected bullet '• Led team of 5' in document.xml")
	}
	if !strings.Contains(xml, "• Increased revenue") {
		t.Error("expected bullet '• Increased revenue' in document.xml")
	}
}

func TestGenerate_EmptySections(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent()

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still be valid DOCX
	if len(result) < 4 || result[0] != 0x50 || result[1] != 0x4B {
		t.Fatal("empty sections should still produce valid DOCX")
	}
}

func TestGenerate_UnknownSectionSkipped(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Test"}),
		},
		domain.CVSection{
			ID: "2", Type: "certifications", Order: 1, Visible: true,
			Content: json.RawMessage(`{"entries":[]}`),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "Test") {
		t.Error("expected personal section content")
	}
}

func TestGenerate_InvisibleSectionsExcluded(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Visible"}),
		},
		domain.CVSection{
			ID: "2", Type: "summary", Order: 1, Visible: false,
			Content: mustJSON(domain.SummaryContent{Text: "Hidden summary text"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "Visible") {
		t.Error("expected visible section content")
	}
	if strings.Contains(xml, "Hidden summary text") {
		t.Error("invisible section content should not appear in DOCX")
	}
}

func TestGenerate_SectionOrdering(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(
		domain.CVSection{
			ID: "1", Type: "summary", Order: 2, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "SUMMARY_SECOND"}),
		},
		domain.CVSection{
			ID: "2", Type: "experience", Order: 1, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
				ID: "e1", Position: "EXPERIENCE_FIRST",
			}}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	expIdx := strings.Index(xml, "EXPERIENCE_FIRST")
	sumIdx := strings.Index(xml, "SUMMARY_SECOND")
	if expIdx == -1 || sumIdx == -1 {
		t.Fatal("expected both sections in document.xml")
	}
	if expIdx > sumIdx {
		t.Error("experience (order=1) should appear before summary (order=2)")
	}
}
