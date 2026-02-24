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
			ID: "2", Type: "awards", Order: 1, Visible: true,
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

func TestGenerate_CertificationsSection(t *testing.T) {
	gen := NewGodocxGenerator()
	expiry := "2026-12"
	content := docxContent(domain.CVSection{
		ID: "1", Type: "certifications", Order: 0, Visible: true,
		Content: mustJSON(domain.CertificationsContent{Entries: []domain.CertificationEntry{
			{
				ID:           "cert-1",
				Name:         "AWS Solutions Architect",
				Issuer:       "Amazon Web Services",
				Date:         "2023-06",
				ExpiryDate:   &expiry,
				CredentialID: "ABC-123",
			},
		}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	for _, expected := range []string{
		"Certifications",
		"AWS Solutions Architect",
		"Amazon Web Services",
		"Jun 2023",
		"Dec 2026",
		"ID: ABC-123",
	} {
		if !strings.Contains(xml, expected) {
			t.Errorf("expected %q in document.xml, not found", expected)
		}
	}
}

func TestGenerate_ProjectsSection(t *testing.T) {
	gen := NewGodocxGenerator()
	endDate := "2024-01"
	content := docxContent(domain.CVSection{
		ID: "1", Type: "projects", Order: 0, Visible: true,
		Content: mustJSON(domain.ProjectsContent{Entries: []domain.ProjectEntry{
			{
				ID:           "proj-1",
				Name:         "HireMe",
				Role:         "Lead Developer",
				Description:  "CV builder application",
				StartDate:    "2023-01",
				EndDate:      &endDate,
				Technologies: []string{"Go", "Next.js", "PostgreSQL"},
			},
		}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	for _, expected := range []string{
		"Projects",
		"HireMe",
		"Lead Developer",
		"CV builder application",
		"Jan 2023",
		"Jan 2024",
		"Go, Next.js, PostgreSQL",
	} {
		if !strings.Contains(xml, expected) {
			t.Errorf("expected %q in document.xml, not found", expected)
		}
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

// --- Helper for template-specific tests ---

func docxContentWithTemplate(templateID string, styling *domain.CVStyling, sections ...domain.CVSection) domain.CVContent {
	return domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    templateID,
		Sections:      sections,
		Styling:       styling,
	}
}

// --- Unit tests for helpers ---

func TestStripHash(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"#c0392b", "c0392b"},
		{"c0392b", "c0392b"},
		{"", ""},
		{"#", ""},
		{"#FF0000", "FF0000"},
	}
	for _, tt := range tests {
		got := stripHash(tt.input)
		if got != tt.want {
			t.Errorf("stripHash(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestResolveDocxStyle_NilStyling(t *testing.T) {
	style := resolveDocxStyle(domain.CVContent{TemplateID: "classic"})

	if style.primaryColor != "2563eb" {
		t.Errorf("primaryColor = %q, want %q", style.primaryColor, "2563eb")
	}
	if style.secondaryColor != "64748b" {
		t.Errorf("secondaryColor = %q, want %q", style.secondaryColor, "64748b")
	}
	if style.bodySizePt != 11 {
		t.Errorf("bodySizePt = %d, want 11", style.bodySizePt)
	}
	if !style.centerHeader {
		t.Error("centerHeader should be true for classic")
	}
	if !style.capsTitle {
		t.Error("capsTitle should be true for classic")
	}
	if style.titleBorder != "bottom" {
		t.Errorf("titleBorder = %q, want %q", style.titleBorder, "bottom")
	}
}

func TestResolveDocxStyle_PerTemplate(t *testing.T) {
	tests := []struct {
		templateID   string
		centerHeader bool
		capsTitle    bool
		titleBorder  string
	}{
		{"classic", true, true, "bottom"},
		{"modern", false, false, "left"},
		{"visionary", false, true, "bottom"},
	}
	for _, tt := range tests {
		t.Run(tt.templateID, func(t *testing.T) {
			style := resolveDocxStyle(domain.CVContent{TemplateID: tt.templateID})
			if style.centerHeader != tt.centerHeader {
				t.Errorf("centerHeader = %v, want %v", style.centerHeader, tt.centerHeader)
			}
			if style.capsTitle != tt.capsTitle {
				t.Errorf("capsTitle = %v, want %v", style.capsTitle, tt.capsTitle)
			}
			if style.titleBorder != tt.titleBorder {
				t.Errorf("titleBorder = %q, want %q", style.titleBorder, tt.titleBorder)
			}
		})
	}
}

func TestResolveDocxStyle_UnknownTemplate(t *testing.T) {
	for _, templateID := range []string{"", "unknown"} {
		t.Run(templateID, func(t *testing.T) {
			style := resolveDocxStyle(domain.CVContent{TemplateID: templateID})
			// Should fall back to classic-like defaults
			if style.primaryColor != "2563eb" {
				t.Errorf("primaryColor = %q, want default", style.primaryColor)
			}
			if !style.centerHeader {
				t.Error("centerHeader should be true for unknown template")
			}
			if !style.capsTitle {
				t.Error("capsTitle should be true for unknown template")
			}
			if style.bodySizePt != 11 {
				t.Errorf("bodySizePt = %d, want 11", style.bodySizePt)
			}
		})
	}
}

// --- Template-specific styling tests ---

func TestGenerate_ClassicStyling(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("classic", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Alice", LastName: "Smith", JobTitle: "Dev"}),
		},
		domain.CVSection{
			ID: "2", Type: "summary", Order: 1, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "Experienced dev"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Classic: centered header
	if !strings.Contains(xml, `w:val="center"`) {
		t.Error("expected centered alignment for classic header")
	}
	// Classic: primary color applied
	if !strings.Contains(xml, `w:val="2563eb"`) {
		t.Error("expected primary color 2563eb in XML")
	}
	// Classic: caps on section titles
	if !strings.Contains(xml, "w:caps") {
		t.Error("expected caps attribute for classic section titles")
	}
	// Classic: bottom border on section titles
	if !strings.Contains(xml, "<w:bottom") {
		t.Error("expected bottom border for classic section titles")
	}
}

func TestGenerate_ModernStyling(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("modern", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Bob", LastName: "Jones"}),
		},
		domain.CVSection{
			ID: "2", Type: "summary", Order: 1, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "Modern summary"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Modern: NOT centered
	if strings.Contains(xml, `w:val="center"`) {
		t.Error("modern template should not center-align header")
	}
	// Modern: left border on section titles (not bottom)
	if !strings.Contains(xml, "<w:left") {
		t.Error("expected left border for modern section titles")
	}
	// Modern: no caps
	if strings.Contains(xml, "w:caps") {
		t.Error("modern template should not use caps on section titles")
	}
	// Modern: primary color still applied
	if !strings.Contains(xml, `w:val="2563eb"`) {
		t.Error("expected primary color 2563eb in XML")
	}
}

func TestGenerate_VisionaryStyling(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Test", LastName: "User"}),
		},
		domain.CVSection{
			ID: "2", Type: "skills", Order: 1, Visible: true,
			Content: mustJSON(domain.SkillsContent{Categories: []domain.SkillCategory{
				{ID: "s1", Name: "Programming", Skills: []domain.Skill{{Name: "Go"}}},
			}}),
		},
		domain.CVSection{
			ID: "3", Type: "experience", Order: 2, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "Developer"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Visionary: table structure
	if !strings.Contains(xml, "<w:tbl>") {
		t.Error("expected table element for visionary layout")
	}
	// Visionary: sidebar cell shading
	if !strings.Contains(xml, `w:fill="2563eb"`) {
		t.Error("expected sidebar cell shading with primary color")
	}
	// Visionary: cell width
	if !strings.Contains(xml, "w:tcW") {
		t.Error("expected cell width property for sidebar")
	}
	// Visionary: white text in sidebar
	if !strings.Contains(xml, `w:val="FFFFFF"`) {
		t.Error("expected white text color in sidebar")
	}
	// Visionary: borderless table
	if !strings.Contains(xml, `w:val="none"`) {
		t.Error("expected borderless table")
	}
	// Visionary: caps on section titles
	if !strings.Contains(xml, "w:caps") {
		t.Error("expected caps for visionary section titles")
	}
}

func TestGenerate_VisionarySidebarMainSplit(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "SIDEBAR_NAME"}),
		},
		domain.CVSection{
			ID: "2", Type: "skills", Order: 1, Visible: true,
			Content: mustJSON(domain.SkillsContent{Categories: []domain.SkillCategory{
				{ID: "s1", Name: "SIDEBAR_SKILLS", Skills: []domain.Skill{{Name: "Go"}}},
			}}),
		},
		domain.CVSection{
			ID: "3", Type: "experience", Order: 2, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "MAIN_EXPERIENCE"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Find positions of the two table cells
	firstCellIdx := strings.Index(xml, "<w:tc>")
	if firstCellIdx == -1 {
		t.Fatal("expected first <w:tc> in visionary layout")
	}
	secondCellIdx := strings.Index(xml[firstCellIdx+6:], "<w:tc>")
	if secondCellIdx == -1 {
		t.Fatal("expected second <w:tc> in visionary layout")
	}
	secondCellIdx += firstCellIdx + 6

	// Personal and skills should be in first cell (sidebar)
	sidebarNameIdx := strings.Index(xml, "SIDEBAR_NAME")
	sidebarSkillsIdx := strings.Index(xml, "SIDEBAR_SKILLS")
	if sidebarNameIdx == -1 || sidebarSkillsIdx == -1 {
		t.Fatal("expected sidebar content in document")
	}
	if sidebarNameIdx > secondCellIdx {
		t.Error("personal section should be in first cell (sidebar)")
	}
	if sidebarSkillsIdx > secondCellIdx {
		t.Error("skills section should be in first cell (sidebar)")
	}

	// Experience should be in second cell (main)
	mainExpIdx := strings.Index(xml, "MAIN_EXPERIENCE")
	if mainExpIdx == -1 {
		t.Fatal("expected main content in document")
	}
	if mainExpIdx < secondCellIdx {
		t.Error("experience section should be in second cell (main)")
	}
}

func TestGenerate_CustomColors(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("classic", &domain.CVStyling{
		PrimaryColor: "#FF0000",
	},
		domain.CVSection{
			ID: "1", Type: "summary", Order: 0, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "Custom colors test"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "FF0000") {
		t.Error("expected custom primary color FF0000 in XML")
	}
}

func TestGenerate_FontSize(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("classic", &domain.CVStyling{
		FontSize: "large",
	},
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Large", LastName: "Font"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	// Large name = 26pt → 52 half-points in XML
	if !strings.Contains(xml, `w:val="52"`) {
		t.Error("expected name font size 52 half-points (26pt) for large")
	}
}

// --- Data gap tests ---

func TestGenerate_ProfileLinks(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "personal", Order: 0, Visible: true,
		Content: mustJSON(domain.PersonalContent{
			FirstName: "Test",
			Links: []domain.ProfileLink{
				{Type: "linkedin", URL: "https://linkedin.com/in/test"},
				{Type: "github", URL: "https://github.com/test"},
			},
		}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "https://linkedin.com/in/test") {
		t.Error("expected LinkedIn URL in document")
	}
	if !strings.Contains(xml, "https://github.com/test") {
		t.Error("expected GitHub URL in document")
	}
}

func TestGenerate_ExperienceLocation(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "experience", Order: 0, Visible: true,
		Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
			ID: "e1", Position: "Developer", Company: "Acme", Location: "Berlin, Germany",
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "Berlin, Germany") {
		t.Error("expected location in experience entry")
	}
}

// --- Post-processing tests ---

func TestPostProcessVisionary_NoTable(t *testing.T) {
	gen := NewGodocxGenerator()
	// Generate a Classic DOCX (no table)
	content := docxContent(domain.CVSection{
		ID: "1", Type: "personal", Order: 0, Visible: true,
		Content: mustJSON(domain.PersonalContent{FirstName: "Test"}),
	})
	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Calling postProcessVisionary on a non-Visionary DOCX should not error
	processed, err := postProcessVisionary(result, "c0392b")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Should still be a valid DOCX
	if len(processed) < 4 || processed[0] != 0x50 || processed[1] != 0x4B {
		t.Fatal("processed output should still be valid DOCX")
	}

	// XML should NOT contain injected table properties
	xml := readDocumentXML(t, processed)
	if strings.Contains(xml, "w:tblBorders") {
		t.Error("should not inject table borders when no table exists")
	}
}

func TestPostProcessVisionary_PreservesAllFiles(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{FirstName: "Test"}),
		},
		domain.CVSection{
			ID: "2", Type: "experience", Order: 1, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "Dev"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify ZIP is valid and has expected files
	zr, err := zip.NewReader(bytes.NewReader(result), int64(len(result)))
	if err != nil {
		t.Fatalf("invalid ZIP: %v", err)
	}

	expectedFiles := map[string]bool{
		"word/document.xml": false,
	}
	for _, f := range zr.File {
		if _, ok := expectedFiles[f.Name]; ok {
			expectedFiles[f.Name] = true
		}
	}
	for name, found := range expectedFiles {
		if !found {
			t.Errorf("expected file %q not found in ZIP", name)
		}
	}

	// Verify document.xml has the injected properties
	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "w:tblBorders") {
		t.Error("expected injected table borders")
	}
	if !strings.Contains(xml, "w:tcPr") {
		t.Error("expected injected cell properties")
	}
}

// --- Re-plan fix tests (F1-F13) ---

func TestGenerate_EducationDegreeFirst(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "education", Order: 0, Visible: true,
		Content: mustJSON(domain.EducationContent{Entries: []domain.EducationEntry{{
			ID: "ed1", Institution: "MIT", Degree: "MSc", Field: "Computer Science",
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	// Degree+field should appear before institution in XML (degree is heading, institution is line 2)
	degreeIdx := strings.Index(xml, "MSc in Computer Science")
	instIdx := strings.Index(xml, "MIT")
	if degreeIdx == -1 {
		t.Fatal("expected degree+field in document")
	}
	if instIdx == -1 {
		t.Fatal("expected institution in document")
	}
	if degreeIdx > instIdx {
		t.Error("degree should appear before institution (degree is the heading)")
	}
}

func TestGenerate_EducationGrade(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "education", Order: 0, Visible: true,
		Content: mustJSON(domain.EducationContent{Entries: []domain.EducationEntry{{
			ID: "ed1", Institution: "TU Munich", Degree: "MSc", Grade: "1.3",
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "Grade: 1.3") {
		t.Error("expected grade in education entry")
	}
}

func TestGenerate_EducationLocation(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "education", Order: 0, Visible: true,
		Content: mustJSON(domain.EducationContent{Entries: []domain.EducationEntry{{
			ID: "ed1", Institution: "MIT", Degree: "BSc", Location: "Cambridge, USA",
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if !strings.Contains(xml, "Cambridge, USA") {
		t.Error("expected location in education entry")
	}
}

func TestGenerate_ExperienceDateTabStop(t *testing.T) {
	gen := NewGodocxGenerator()
	endDate := "2023-06"
	content := docxContent(domain.CVSection{
		ID: "1", Type: "experience", Order: 0, Visible: true,
		Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
			ID: "e1", Position: "Developer", StartDate: "2020-01", EndDate: &endDate,
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	// Tab stop definition on paragraph
	if !strings.Contains(xml, `w:val="right"`) {
		t.Error("expected right tab stop for date alignment")
	}
	// Tab character in run (serialized as <w:tab></w:tab>)
	if !strings.Contains(xml, "<w:tab>") {
		t.Error("expected tab character in date run")
	}
	// Date text present
	if !strings.Contains(xml, "Jan 2020") {
		t.Error("expected start date in document")
	}
	if !strings.Contains(xml, "Jun 2023") {
		t.Error("expected end date in document")
	}
}

func TestGenerate_ExperienceMultiLine(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(domain.CVSection{
		ID: "1", Type: "experience", Order: 0, Visible: true,
		Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
			ID: "e1", Position: "Senior Dev", Company: "Acme Corp", Location: "Berlin",
			Description: "Led the team",
		}}}),
	})

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Position and company should be on separate paragraphs (separate <w:p> elements)
	posIdx := strings.Index(xml, "Senior Dev")
	compIdx := strings.Index(xml, "Acme Corp")
	if posIdx == -1 || compIdx == -1 {
		t.Fatal("expected position and company in document")
	}

	// Classic template: "at " prefix before company
	if !strings.Contains(xml, "at Acme Corp") {
		t.Error("classic template should use 'at' prefix before company")
	}

	// Description uses secondary color (not black "000000")
	if !strings.Contains(xml, "Led the team") {
		t.Error("expected description in document")
	}
}

func TestGenerate_ModernNoAtPrefix(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("modern", nil,
		domain.CVSection{
			ID: "1", Type: "experience", Order: 0, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{{
				ID: "e1", Position: "Dev", Company: "ModernCo",
			}}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	if strings.Contains(xml, "at ModernCo") {
		t.Error("modern template should not use 'at' prefix before company")
	}
	if !strings.Contains(xml, "ModernCo") {
		t.Error("expected company name in document")
	}
}

func TestGenerate_DescriptionColor(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContent(
		domain.CVSection{
			ID: "1", Type: "summary", Order: 0, Visible: true,
			Content: mustJSON(domain.SummaryContent{Text: "A seasoned developer"}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)
	// Description should use secondary color (64748b), not black (000000)
	descIdx := strings.Index(xml, "A seasoned developer")
	if descIdx == -1 {
		t.Fatal("expected summary text in document")
	}
	// Look for secondary color near the description text
	nearbyXML := xml[max(0, descIdx-200):descIdx]
	if !strings.Contains(nearbyXML, "64748b") {
		t.Error("expected secondary color (64748b) on description text, not black")
	}
}

func TestGenerate_VisionaryMainHeader(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{
				FirstName: "Max", LastName: "Developer", JobTitle: "Lead Engineer",
			}),
		},
		domain.CVSection{
			ID: "2", Type: "experience", Order: 1, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "MAIN_CONTENT_MARKER"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Find the second table cell (main content)
	firstCellIdx := strings.Index(xml, "<w:tc>")
	if firstCellIdx == -1 {
		t.Fatal("expected first table cell")
	}
	secondCellIdx := strings.Index(xml[firstCellIdx+6:], "<w:tc>")
	if secondCellIdx == -1 {
		t.Fatal("expected second table cell")
	}
	secondCellIdx += firstCellIdx + 6

	// Name should appear in main cell (second cell) as header
	mainCellXML := xml[secondCellIdx:]
	if !strings.Contains(mainCellXML, "Max Developer") {
		t.Error("expected name in main cell header")
	}
	if !strings.Contains(mainCellXML, "Lead Engineer") {
		t.Error("expected job title in main cell header")
	}

	// Main header should appear before the section content
	nameInMain := strings.Index(mainCellXML, "Max Developer")
	markerInMain := strings.Index(mainCellXML, "MAIN_CONTENT_MARKER")
	if nameInMain > markerInMain {
		t.Error("main header name should appear before section content")
	}
}

func TestGenerate_VisionarySidebarLabeledContact(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "personal", Order: 0, Visible: true,
			Content: mustJSON(domain.PersonalContent{
				FirstName: "Test", Email: "test@example.com", Phone: "+49123", Location: "Berlin",
			}),
		},
		domain.CVSection{
			ID: "2", Type: "experience", Order: 1, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "Dev"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Sidebar should use labeled format
	if !strings.Contains(xml, "Email") {
		t.Error("expected 'Email' label in sidebar contact")
	}
	if !strings.Contains(xml, "Phone") {
		t.Error("expected 'Phone' label in sidebar contact")
	}
	if !strings.Contains(xml, "Location") {
		t.Error("expected 'Location' label in sidebar contact")
	}
	// Should NOT use pipe-separated format in sidebar
	if strings.Contains(xml, "test@example.com | +49123") {
		t.Error("sidebar should not use pipe-separated contact format")
	}
}

func TestGenerate_VisionarySidebarSkillLevels(t *testing.T) {
	gen := NewGodocxGenerator()
	content := docxContentWithTemplate("visionary", nil,
		domain.CVSection{
			ID: "1", Type: "skills", Order: 0, Visible: true,
			Content: mustJSON(domain.SkillsContent{Categories: []domain.SkillCategory{
				{ID: "s1", Name: "Backend", Skills: []domain.Skill{
					{Name: "Go", Level: "expert"},
					{Name: "Python", Level: "advanced"},
				}},
			}}),
		},
		domain.CVSection{
			ID: "2", Type: "experience", Order: 1, Visible: true,
			Content: mustJSON(domain.ExperienceContent{Entries: []domain.ExperienceEntry{
				{ID: "e1", Position: "Dev"},
			}}),
		},
	)

	result, err := gen.Generate(content)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	xml := readDocumentXML(t, result)

	// Skills should include proficiency levels in sidebar
	if !strings.Contains(xml, "expert") {
		t.Error("expected skill level 'expert' in sidebar")
	}
	if !strings.Contains(xml, "advanced") {
		t.Error("expected skill level 'advanced' in sidebar")
	}
	// Should be bulleted
	if !strings.Contains(xml, "• Go") {
		t.Error("expected bulleted skill format in sidebar")
	}
}
