package export

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/ellebam/hireme/api/internal/domain"
)

// testCVContent returns a full CVContent matching the seed data structure.
func testCVContent(templateID string) domain.CVContent {
	endDate := "2019-12"
	return domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    templateID,
		Locale:        "en",
		Title:         "My Professional CV",
		Sections: []domain.CVSection{
			{
				ID: "sec-personal-001", Type: "personal", Order: 0, Visible: true,
				Content: json.RawMessage(`{
					"firstName": "Max",
					"lastName": "Developer",
					"jobTitle": "Senior Software Engineer",
					"email": "max@example.com",
					"phone": "+49 123 456789",
					"location": "Berlin, Germany",
					"links": [
						{"type": "linkedin", "url": "https://linkedin.com/in/maxdev", "label": "LinkedIn"},
						{"type": "github", "url": "https://github.com/maxdev"}
					]
				}`),
			},
			{
				ID: "sec-summary-001", Type: "summary", Order: 1, Visible: true,
				Content: json.RawMessage(`{"text": "Experienced software engineer with 10+ years in building scalable web applications."}`),
			},
			{
				ID: "sec-experience-001", Type: "experience", Order: 2, Visible: true,
				Content: mustMarshal(domain.ExperienceContent{
					Entries: []domain.ExperienceEntry{
						{
							ID: "exp-001", Company: "TechCorp GmbH", Position: "Senior Software Engineer",
							Location: "Berlin, Germany", StartDate: "2020-01", Current: true,
							Description: "Leading development of cloud-native applications.",
							Highlights:  []string{"Architected microservices platform", "Reduced deployment time by 70%"},
						},
						{
							ID: "exp-002", Company: "StartupXYZ", Position: "Full Stack Developer",
							Location: "Munich, Germany", StartDate: "2017-06", EndDate: &endDate,
							Description: "Built and maintained multiple web applications.",
						},
					},
				}),
			},
			{
				ID: "sec-education-001", Type: "education", Order: 3, Visible: true,
				Content: json.RawMessage(`{
					"entries": [{
						"id": "edu-001",
						"institution": "Technical University of Munich",
						"degree": "M.Sc.",
						"field": "Computer Science",
						"location": "Munich, Germany",
						"startDate": "2014-10",
						"endDate": "2017-03",
						"grade": "1.3"
					}]
				}`),
			},
			{
				ID: "sec-skills-001", Type: "skills", Order: 4, Visible: true,
				Content: json.RawMessage(`{
					"categories": [
						{"id": "cat-001", "name": "Programming Languages", "skills": [{"name": "Go"}, {"name": "TypeScript"}, {"name": "Python"}]},
						{"id": "cat-002", "name": "Frameworks", "skills": [{"name": "React"}, {"name": "Next.js"}]}
					]
				}`),
			},
			{
				ID: "sec-languages-001", Type: "languages", Order: 5, Visible: true,
				Content: json.RawMessage(`{
					"entries": [
						{"language": "German", "proficiency": "native"},
						{"language": "English", "proficiency": "fluent"}
					]
				}`),
			},
			{
				ID: "sec-certifications-001", Type: "certifications", Order: 6, Visible: true,
				Content: mustMarshal(domain.CertificationsContent{
					Entries: []domain.CertificationEntry{
						{
							ID: "cert-001", Name: "AWS Solutions Architect", Issuer: "Amazon Web Services",
							Date: "2023-06", CredentialID: "SAA-C03-12345",
						},
					},
				}),
			},
			{
				ID: "sec-projects-001", Type: "projects", Order: 7, Visible: true,
				Content: mustMarshal(domain.ProjectsContent{
					Entries: []domain.ProjectEntry{
						{
							ID: "proj-001", Name: "Open Source CLI Tool", Role: "Creator",
							Description: "A developer productivity tool", StartDate: "2022-03",
							Technologies: []string{"Go", "Cobra"},
						},
					},
				}),
			},
		},
		Styling: &domain.CVStyling{
			PrimaryColor: "#2563eb",
			FontFamily:   "inter",
			FontSize:     "medium",
		},
	}
}

func mustMarshal(v any) json.RawMessage {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}

func TestNewRenderer_Success(t *testing.T) {
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("NewRenderer() error: %v", err)
	}
	if r == nil {
		t.Fatal("NewRenderer() returned nil")
	}
	if len(r.templates) != 3 {
		t.Errorf("expected 3 templates, got %d", len(r.templates))
	}
}

func TestRender_UnknownTemplate(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{TemplateID: "nonexistent"}
	_, err := r.Render(content)
	if err == nil {
		t.Fatal("expected error for unknown template")
	}
	if !strings.Contains(err.Error(), "nonexistent") {
		t.Errorf("error should mention template name, got: %v", err)
	}
}

func TestFormatDate(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{"2020-01", "Jan 2020"},
		{"2017-06", "Jun 2017"},
		{"2020", "2020"},
		{"", ""},
		{"invalid", "invalid"},
	}
	for _, tt := range tests {
		got := formatDate(tt.input)
		if got != tt.want {
			t.Errorf("formatDate(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatDateRange(t *testing.T) {
	end := "2019-12"
	tests := []struct {
		name    string
		start   string
		end     *string
		current bool
		want    string
	}{
		{"current", "2020-01", nil, true, "Jan 2020 - Present"},
		{"range", "2017-06", &end, false, "Jun 2017 - Dec 2019"},
		{"start only", "2020-01", nil, false, "Jan 2020"},
		{"empty", "", nil, false, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDateRange(tt.start, tt.end, tt.current)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}

func TestProficiencyBars(t *testing.T) {
	tests := []struct {
		proficiency string
		wantFilled  int
	}{
		{"native", 5},
		{"fluent", 4},
		{"advanced", 3},
		{"intermediate", 2},
		{"basic", 1},
		{"unknown", 0},
	}
	for _, tt := range tests {
		t.Run(tt.proficiency, func(t *testing.T) {
			bars := proficiencyBars(tt.proficiency, "#2563eb")
			filled := 0
			for _, b := range bars {
				if string(b) == "#2563eb" {
					filled++
				}
			}
			if filled != tt.wantFilled {
				t.Errorf("proficiency %q: got %d filled bars, want %d", tt.proficiency, filled, tt.wantFilled)
			}
			if len(bars) != 5 {
				t.Errorf("expected 5 bars, got %d", len(bars))
			}
		})
	}
}

// Per-template rendering tests

func TestRender_Classic_FullCV(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("classic")
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertHTMLStructure(t, html)
	assertSelfContained(t, html)
	assertContains(t, html, "Max Developer")
	assertContains(t, html, "Senior Software Engineer")
	assertContains(t, html, "max@example.com")
	assertContains(t, html, "TechCorp GmbH")
	assertContains(t, html, "StartupXYZ")
	assertContains(t, html, "Technical University of Munich")
	assertContains(t, html, "M.Sc.")
	assertContains(t, html, "Go")
	assertContains(t, html, "TypeScript")
	assertContains(t, html, "German")
	assertContains(t, html, "English")
	assertContains(t, html, "#2563eb") // dynamic color
	// Certifications
	assertContains(t, html, "AWS Solutions Architect")
	assertContains(t, html, "Amazon Web Services")
	assertContains(t, html, "SAA-C03-12345")
	// Projects
	assertContains(t, html, "Open Source CLI Tool")
	assertContains(t, html, "Creator")
	assertContains(t, html, "A developer productivity tool")
	assertContains(t, html, "Go, Cobra")
	// Classic-specific: bottom-border section titles
	assertContains(t, html, "classic-section-title")
}

func TestRender_Modern_FullCV(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("modern")
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertHTMLStructure(t, html)
	assertSelfContained(t, html)
	assertContains(t, html, "Max Developer")
	assertContains(t, html, "Senior Software Engineer")
	assertContains(t, html, "TechCorp GmbH")
	assertContains(t, html, "Go")
	assertContains(t, html, "German")
	// Certifications & Projects
	assertContains(t, html, "AWS Solutions Architect")
	assertContains(t, html, "Open Source CLI Tool")
	// Modern-specific: timeline, skill pills, language bars
	assertContains(t, html, "modern-timeline")
	assertContains(t, html, "modern-skill-pill")
	assertContains(t, html, "modern-lang-bar")
}

func TestRender_Visionary_FullCV(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("visionary")
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertHTMLStructure(t, html)
	assertSelfContained(t, html)
	assertContains(t, html, "Max Developer")
	assertContains(t, html, "Senior Software Engineer")
	assertContains(t, html, "TechCorp GmbH")
	// Certifications & Projects (in main area)
	assertContains(t, html, "AWS Solutions Architect")
	assertContains(t, html, "Open Source CLI Tool")
	// Visionary-specific: sidebar + main layout
	assertContains(t, html, "visionary-sidebar")
	assertContains(t, html, "visionary-main")
	assertContains(t, html, "visionary-layout")
}

func TestRender_Visionary_SidebarRouting(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("visionary")
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	// Sidebar should contain personal info, skills, languages
	sidebarIdx := strings.Index(html, `class="visionary-sidebar"`)
	mainIdx := strings.Index(html, `class="visionary-main"`)
	if sidebarIdx < 0 || mainIdx < 0 {
		t.Fatal("missing sidebar or main sections")
	}
	sidebar := html[sidebarIdx:mainIdx]
	main := html[mainIdx:]
	// Personal data in sidebar
	if !strings.Contains(sidebar, "max@example.com") {
		t.Error("sidebar should contain email")
	}
	if !strings.Contains(sidebar, "Go") {
		t.Error("sidebar should contain skills")
	}
	if !strings.Contains(sidebar, "German") {
		t.Error("sidebar should contain languages")
	}
	// Experience in main
	if !strings.Contains(main, "TechCorp GmbH") {
		t.Error("main should contain experience")
	}
	if !strings.Contains(main, "Technical University of Munich") {
		t.Error("main should contain education")
	}
}

func TestRender_HiddenSections(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("classic")
	// Hide the summary section
	content.Sections[1].Visible = false
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if strings.Contains(html, "Experienced software engineer with 10+") {
		t.Error("hidden summary should not appear in output")
	}
	// Other sections should still be present
	assertContains(t, html, "Max Developer")
	assertContains(t, html, "TechCorp GmbH")
}

func TestRender_MinimalCV(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections: []domain.CVSection{
			{
				ID: "personal", Type: "personal", Order: 0, Visible: true,
				Content: json.RawMessage(`{"firstName": "Jane", "lastName": "Doe"}`),
			},
		},
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertHTMLStructure(t, html)
	assertContains(t, html, "Jane Doe")
}

func TestRender_EmptySections(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "modern",
		Sections: []domain.CVSection{
			{ID: "exp", Type: "experience", Order: 0, Visible: true, Content: json.RawMessage(`{"entries": []}`)},
			{ID: "edu", Type: "education", Order: 1, Visible: true, Content: json.RawMessage(`{"entries": []}`)},
			{ID: "skills", Type: "skills", Order: 2, Visible: true, Content: json.RawMessage(`{"categories": []}`)},
			{ID: "langs", Type: "languages", Order: 3, Visible: true, Content: json.RawMessage(`{"entries": []}`)},
		},
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertHTMLStructure(t, html)
}

func TestRender_SpecialCharacters(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections: []domain.CVSection{
			{
				ID: "personal", Type: "personal", Order: 0, Visible: true,
				Content: json.RawMessage(`{"firstName": "<script>alert('xss')</script>", "lastName": "O'Brien & Co."}`),
			},
		},
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	if strings.Contains(html, "<script>") {
		t.Error("HTML should escape script tags")
	}
	if !strings.Contains(html, "&lt;script&gt;") {
		t.Error("script tag should be HTML-escaped")
	}
	assertContains(t, html, "O&#39;Brien &amp; Co.")
}

func TestRender_CustomStyling(t *testing.T) {
	r, _ := NewRenderer()
	content := testCVContent("classic")
	content.Styling = &domain.CVStyling{
		PrimaryColor:   "#ff0000",
		SecondaryColor: "#00ff00",
		FontFamily:     "merriweather",
		FontSize:       "large",
		LineHeight:     "relaxed",
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertContains(t, html, "#ff0000")
	assertContains(t, html, "#00ff00")
	assertContains(t, html, "Merriweather")
	assertContains(t, html, "16px")
	assertContains(t, html, "1.8")
}

func TestRender_DefaultStyling(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections: []domain.CVSection{
			{ID: "summary", Type: "summary", Order: 0, Visible: true, Content: json.RawMessage(`{"text": "Hello"}`)},
		},
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	assertContains(t, html, defaultPrimaryColor)
	assertContains(t, html, "Inter")
}

func TestRender_SectionOrdering(t *testing.T) {
	r, _ := NewRenderer()
	content := domain.CVContent{
		SchemaVersion: "1.0.0",
		TemplateID:    "classic",
		Sections: []domain.CVSection{
			{ID: "skills", Type: "skills", Order: 5, Visible: true, Content: json.RawMessage(`{"categories": [{"id":"c1","name":"Cat","skills":[{"name":"SKILL_MARKER"}]}]}`)},
			{ID: "personal", Type: "personal", Order: 0, Visible: true, Content: json.RawMessage(`{"firstName": "PERSONAL_MARKER"}`)},
			{ID: "exp", Type: "experience", Order: 3, Visible: true, Content: json.RawMessage(`{"entries": [{"id":"e1","company":"EXP_MARKER","position":"Dev","startDate":"2020-01"}]}`)},
		},
	}
	html, err := r.Render(content)
	if err != nil {
		t.Fatalf("Render error: %v", err)
	}
	personalIdx := strings.Index(html, "PERSONAL_MARKER")
	expIdx := strings.Index(html, "EXP_MARKER")
	skillIdx := strings.Index(html, "SKILL_MARKER")
	if personalIdx < 0 || expIdx < 0 || skillIdx < 0 {
		t.Fatal("missing markers in output")
	}
	if personalIdx > expIdx || expIdx > skillIdx {
		t.Errorf("sections not in order: personal=%d, exp=%d, skills=%d", personalIdx, expIdx, skillIdx)
	}
}

func TestRender_CrossTemplateConsistency(t *testing.T) {
	r, _ := NewRenderer()
	dataPoints := []string{"Max Developer", "TechCorp GmbH", "Go", "German", "Technical University of Munich", "AWS Solutions Architect", "Open Source CLI Tool"}
	for _, tmplID := range []string{"classic", "modern", "visionary"} {
		t.Run(tmplID, func(t *testing.T) {
			content := testCVContent(tmplID)
			html, err := r.Render(content)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}
			for _, dp := range dataPoints {
				if !strings.Contains(html, dp) {
					t.Errorf("template %s missing data point %q", tmplID, dp)
				}
			}
		})
	}
}

func TestRender_FontFamilyMapping(t *testing.T) {
	r, _ := NewRenderer()
	tests := []struct {
		family   string
		contains string
	}{
		{"inter", "Inter"},
		{"merriweather", "Merriweather"},
		{"roboto", "Roboto"},
		{"lato", "Lato"},
		{"opensans", "Open Sans"},
		{"unknown", "Inter"}, // fallback to default
		{"", "Inter"},        // empty = default
	}
	for _, tt := range tests {
		t.Run(tt.family, func(t *testing.T) {
			content := domain.CVContent{
				SchemaVersion: "1.0.0",
				TemplateID:    "classic",
				Sections:      []domain.CVSection{},
				Styling:       &domain.CVStyling{FontFamily: tt.family},
			}
			html, err := r.Render(content)
			if err != nil {
				t.Fatalf("Render error: %v", err)
			}
			if !strings.Contains(html, tt.contains) {
				t.Errorf("font %q should produce CSS containing %q", tt.family, tt.contains)
			}
		})
	}
}

// Helper assertions

func assertHTMLStructure(t *testing.T, html string) {
	t.Helper()
	checks := []string{"<!DOCTYPE html>", "<html", "<head>", "<body>", "</html>"}
	for _, c := range checks {
		if !strings.Contains(html, c) {
			t.Errorf("HTML missing %q", c)
		}
	}
}

func assertSelfContained(t *testing.T, html string) {
	t.Helper()
	// No external resource references
	forbidden := []string{
		`<link rel="stylesheet" href="http`,
		`<script src="http`,
		`url(http`,
	}
	for _, f := range forbidden {
		if strings.Contains(html, f) {
			t.Errorf("HTML contains external resource reference: %q", f)
		}
	}
}

func assertContains(t *testing.T, html, substr string) {
	t.Helper()
	if !strings.Contains(html, substr) {
		t.Errorf("HTML should contain %q", substr)
	}
}
