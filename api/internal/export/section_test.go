package export

import (
	"encoding/json"
	"testing"
)

func TestParsePersonal_Valid(t *testing.T) {
	raw := json.RawMessage(`{
		"firstName": "Max",
		"lastName": "Developer",
		"jobTitle": "Engineer",
		"email": "max@example.com",
		"phone": "+49 123",
		"location": "Berlin",
		"links": [{"type": "github", "url": "https://github.com/max", "label": "GitHub"}]
	}`)
	c := ParsePersonal(raw)
	if c.FirstName != "Max" {
		t.Errorf("FirstName: got %q, want %q", c.FirstName, "Max")
	}
	if c.LastName != "Developer" {
		t.Errorf("LastName: got %q, want %q", c.LastName, "Developer")
	}
	if c.Email != "max@example.com" {
		t.Errorf("Email: got %q, want %q", c.Email, "max@example.com")
	}
	if len(c.Links) != 1 {
		t.Fatalf("Links: got %d, want 1", len(c.Links))
	}
	if c.Links[0].Label != "GitHub" {
		t.Errorf("Link label: got %q, want %q", c.Links[0].Label, "GitHub")
	}
}

func TestParsePersonal_Malformed(t *testing.T) {
	c := ParsePersonal(json.RawMessage(`{broken`))
	if c.FirstName != "" {
		t.Errorf("expected zero-value, got FirstName=%q", c.FirstName)
	}
}

func TestParsePersonal_Empty(t *testing.T) {
	c := ParsePersonal(nil)
	if c.FirstName != "" {
		t.Errorf("expected zero-value, got FirstName=%q", c.FirstName)
	}
}

func TestParseSummary_Valid(t *testing.T) {
	raw := json.RawMessage(`{"text": "Experienced engineer"}`)
	c := ParseSummary(raw)
	if c.Text != "Experienced engineer" {
		t.Errorf("Text: got %q, want %q", c.Text, "Experienced engineer")
	}
}

func TestParseSummary_Malformed(t *testing.T) {
	c := ParseSummary(json.RawMessage(`{broken`))
	if c.Text != "" {
		t.Errorf("expected zero-value, got Text=%q", c.Text)
	}
}

func TestParseExperience_Valid(t *testing.T) {
	raw := json.RawMessage(`{
		"entries": [{
			"id": "exp-1",
			"company": "ACME",
			"position": "Dev",
			"startDate": "2020-01",
			"current": true,
			"highlights": ["Built things"]
		}]
	}`)
	c := ParseExperience(raw)
	if len(c.Entries) != 1 {
		t.Fatalf("Entries: got %d, want 1", len(c.Entries))
	}
	if c.Entries[0].Company != "ACME" {
		t.Errorf("Company: got %q, want %q", c.Entries[0].Company, "ACME")
	}
	if !c.Entries[0].Current {
		t.Error("Current: got false, want true")
	}
}

func TestParseExperience_Malformed(t *testing.T) {
	c := ParseExperience(json.RawMessage(`{broken`))
	if len(c.Entries) != 0 {
		t.Errorf("expected zero entries, got %d", len(c.Entries))
	}
}

func TestParseEducation_Valid(t *testing.T) {
	raw := json.RawMessage(`{
		"entries": [{
			"id": "edu-1",
			"institution": "MIT",
			"degree": "B.Sc.",
			"field": "CS",
			"grade": "3.9"
		}]
	}`)
	c := ParseEducation(raw)
	if len(c.Entries) != 1 {
		t.Fatalf("Entries: got %d, want 1", len(c.Entries))
	}
	if c.Entries[0].Institution != "MIT" {
		t.Errorf("Institution: got %q", c.Entries[0].Institution)
	}
	if c.Entries[0].Grade != "3.9" {
		t.Errorf("Grade: got %q", c.Entries[0].Grade)
	}
}

func TestParseSkills_Valid(t *testing.T) {
	raw := json.RawMessage(`{
		"categories": [{
			"id": "cat-1",
			"name": "Languages",
			"skills": [{"name": "Go"}, {"name": "TS", "level": "expert"}]
		}]
	}`)
	c := ParseSkills(raw)
	if len(c.Categories) != 1 {
		t.Fatalf("Categories: got %d, want 1", len(c.Categories))
	}
	if len(c.Categories[0].Skills) != 2 {
		t.Fatalf("Skills: got %d, want 2", len(c.Categories[0].Skills))
	}
}

func TestParseLanguages_Valid(t *testing.T) {
	raw := json.RawMessage(`{
		"entries": [
			{"language": "German", "proficiency": "native"},
			{"language": "English", "proficiency": "fluent", "certification": "TOEFL"}
		]
	}`)
	c := ParseLanguages(raw)
	if len(c.Entries) != 2 {
		t.Fatalf("Entries: got %d, want 2", len(c.Entries))
	}
	if c.Entries[0].Language != "German" {
		t.Errorf("Language: got %q", c.Entries[0].Language)
	}
	if c.Entries[1].Certification != "TOEFL" {
		t.Errorf("Certification: got %q", c.Entries[1].Certification)
	}
}

func TestParseLanguages_Malformed(t *testing.T) {
	c := ParseLanguages(json.RawMessage(`{broken`))
	if len(c.Entries) != 0 {
		t.Errorf("expected zero entries, got %d", len(c.Entries))
	}
}
