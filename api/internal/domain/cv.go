package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// CV represents a curriculum vitae document
type CV struct {
	ID            uuid.UUID
	UserID        string
	Title         string
	SchemaVersion string
	Content       json.RawMessage
	IsActive      bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// CVContent represents the structured content of a CV
// This is used for parsing/validation, the actual storage is JSON
type CVContent struct {
	SchemaVersion string       `json:"schemaVersion"`
	TemplateID    string       `json:"templateId"`
	Locale        string       `json:"locale,omitempty"`
	Title         string       `json:"title,omitempty"`
	Sections      []CVSection  `json:"sections"`
	Styling       *CVStyling   `json:"styling,omitempty"`
}

// CVSection represents a section in the CV
type CVSection struct {
	ID      string          `json:"id"`
	Type    string          `json:"type"`
	Order   int             `json:"order"`
	Visible bool            `json:"visible"`
	Title   string          `json:"title,omitempty"`
	Content json.RawMessage `json:"content"`
}

// CVStyling represents the visual styling options
type CVStyling struct {
	PrimaryColor   string `json:"primaryColor,omitempty"`
	SecondaryColor string `json:"secondaryColor,omitempty"`
	FontFamily     string `json:"fontFamily,omitempty"`
	FontSize       string `json:"fontSize,omitempty"`
	LineHeight     string `json:"lineHeight,omitempty"`
	ShowIcons      *bool  `json:"showIcons,omitempty"`
}

// Section type constants
const (
	SectionTypePersonal       = "personal"
	SectionTypeSummary        = "summary"
	SectionTypeExperience     = "experience"
	SectionTypeEducation      = "education"
	SectionTypeSkills         = "skills"
	SectionTypeLanguages      = "languages"
	SectionTypeCertifications = "certifications"
	SectionTypeProjects       = "projects"
	SectionTypeAwards         = "awards"
	SectionTypePublications   = "publications"
	SectionTypeReferences     = "references"
	SectionTypeCustom         = "custom"
)

// Template ID constants
const (
	TemplateClassic = "classic"
	TemplateModern  = "modern"
	TemplateVisionary = "visionary"
)

// ParseContent parses the CV's raw JSON content into CVContent
func (cv *CV) ParseContent() (*CVContent, error) {
	var content CVContent
	if err := json.Unmarshal(cv.Content, &content); err != nil {
		return nil, err
	}
	return &content, nil
}

// SetContent marshals CVContent back to raw JSON
func (cv *CV) SetContent(content *CVContent) error {
	data, err := json.Marshal(content)
	if err != nil {
		return err
	}
	cv.Content = data
	return nil
}

// PersonalContent represents the personal info section
type PersonalContent struct {
	FirstName       string        `json:"firstName,omitempty"`
	LastName        string        `json:"lastName,omitempty"`
	JobTitle        string        `json:"jobTitle,omitempty"`
	Email           string        `json:"email,omitempty"`
	Phone           string        `json:"phone,omitempty"`
	Location        string        `json:"location,omitempty"`
	PortraitAssetID *string       `json:"portraitAssetId,omitempty"`
	Links           []ProfileLink `json:"links,omitempty"`
}

// ProfileLink represents a social/professional link
type ProfileLink struct {
	Type  string `json:"type"`
	URL   string `json:"url"`
	Label string `json:"label,omitempty"`
}

// ExperienceContent represents work experience section
type ExperienceContent struct {
	Entries []ExperienceEntry `json:"entries"`
}

// ExperienceEntry represents a single work experience
type ExperienceEntry struct {
	ID          string   `json:"id"`
	Company     string   `json:"company"`
	Position    string   `json:"position"`
	Location    string   `json:"location,omitempty"`
	StartDate   string   `json:"startDate"`
	EndDate     *string  `json:"endDate,omitempty"`
	Current     bool     `json:"current,omitempty"`
	Description string   `json:"description,omitempty"`
	Highlights  []string `json:"highlights,omitempty"`
}

// EducationContent represents education section
type EducationContent struct {
	Entries []EducationEntry `json:"entries"`
}

// EducationEntry represents a single education entry
type EducationEntry struct {
	ID          string  `json:"id"`
	Institution string  `json:"institution"`
	Degree      string  `json:"degree"`
	Field       string  `json:"field,omitempty"`
	Location    string  `json:"location,omitempty"`
	StartDate   string  `json:"startDate,omitempty"`
	EndDate     *string `json:"endDate,omitempty"`
	Current     bool    `json:"current,omitempty"`
	Grade       string  `json:"grade,omitempty"`
	Description string  `json:"description,omitempty"`
}

// SkillsContent represents skills section
type SkillsContent struct {
	Categories []SkillCategory `json:"categories"`
}

// SkillCategory represents a category of skills
type SkillCategory struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Skills []Skill `json:"skills"`
}

// Skill represents a single skill
type Skill struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"`
}

// SummaryContent represents the professional summary section
type SummaryContent struct {
	Text string `json:"text"`
}

// LanguagesContent represents the languages section
type LanguagesContent struct {
	Entries []LanguageEntry `json:"entries"`
}

// LanguageEntry represents a single language
type LanguageEntry struct {
	ID            string `json:"id,omitempty"`
	Language      string `json:"language"`
	Proficiency   string `json:"proficiency"`
	Certification string `json:"certification,omitempty"`
}

// CertificationsContent represents the certifications section
type CertificationsContent struct {
	Entries []CertificationEntry `json:"entries"`
}

// CertificationEntry represents a single certification
type CertificationEntry struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Issuer       string  `json:"issuer"`
	Date         string  `json:"date,omitempty"`
	ExpiryDate   *string `json:"expiryDate,omitempty"`
	CredentialID string  `json:"credentialId,omitempty"`
	URL          string  `json:"url,omitempty"`
}

// ProjectsContent represents the projects section
type ProjectsContent struct {
	Entries []ProjectEntry `json:"entries"`
}

// ProjectEntry represents a single project
type ProjectEntry struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Role         string   `json:"role,omitempty"`
	Description  string   `json:"description"`
	URL          string   `json:"url,omitempty"`
	Technologies []string `json:"technologies,omitempty"`
	StartDate    string   `json:"startDate,omitempty"`
	EndDate      *string  `json:"endDate,omitempty"`
}
