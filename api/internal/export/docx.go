package export

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"

	"github.com/ellebam/hireme/api/internal/domain"
)

// DOCXGenerator generates DOCX documents from structured CV content.
type DOCXGenerator interface {
	Generate(content domain.CVContent) ([]byte, error)
}

// GodocxGenerator implements DOCXGenerator using the gomutex/godocx library.
type GodocxGenerator struct{}

// NewGodocxGenerator creates a new GodocxGenerator.
func NewGodocxGenerator() *GodocxGenerator {
	return &GodocxGenerator{}
}

// Generate produces a DOCX document from structured CV content.
func (g *GodocxGenerator) Generate(content domain.CVContent) ([]byte, error) {
	doc, err := godocx.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("creating document: %w", err)
	}

	// Collect visible sections sorted by order
	type orderedSection struct {
		order   int
		section domain.CVSection
	}
	var ordered []orderedSection
	for _, sec := range content.Sections {
		if !sec.Visible {
			continue
		}
		ordered = append(ordered, orderedSection{order: sec.Order, section: sec})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].order < ordered[j].order
	})

	for _, entry := range ordered {
		sec := entry.section
		title := sec.Title
		if title == "" {
			if label, ok := sectionLabels[sec.Type]; ok {
				title = label
			}
		}

		switch sec.Type {
		case domain.SectionTypePersonal:
			addPersonalSection(doc, sec)
		case domain.SectionTypeSummary:
			addSummarySection(doc, title, sec)
		case domain.SectionTypeExperience:
			addExperienceSection(doc, title, sec)
		case domain.SectionTypeEducation:
			addEducationSection(doc, title, sec)
		case domain.SectionTypeSkills:
			addSkillsSection(doc, title, sec)
		case domain.SectionTypeLanguages:
			addLanguagesSection(doc, title, sec)
		default:
			// Unknown section types silently skipped
		}
	}

	var buf bytes.Buffer
	if err := doc.Write(&buf); err != nil {
		return nil, fmt.Errorf("writing document: %w", err)
	}

	return buf.Bytes(), nil
}

func addPersonalSection(doc *docx.RootDoc, sec domain.CVSection) {
	personal := ParsePersonal(sec.Content)

	fullName := strings.TrimSpace(personal.FirstName + " " + personal.LastName)
	if fullName != "" {
		doc.AddHeading(fullName, 0) //nolint:errcheck
	}

	if personal.JobTitle != "" {
		doc.AddParagraph(personal.JobTitle)
	}

	var contactParts []string
	if personal.Email != "" {
		contactParts = append(contactParts, personal.Email)
	}
	if personal.Phone != "" {
		contactParts = append(contactParts, personal.Phone)
	}
	if personal.Location != "" {
		contactParts = append(contactParts, personal.Location)
	}
	if len(contactParts) > 0 {
		doc.AddParagraph(strings.Join(contactParts, " | "))
	}
}

func addSummarySection(doc *docx.RootDoc, title string, sec domain.CVSection) {
	summary := ParseSummary(sec.Content)
	doc.AddHeading(title, 1) //nolint:errcheck
	if summary.Text != "" {
		doc.AddParagraph(summary.Text)
	}
}

func addExperienceSection(doc *docx.RootDoc, title string, sec domain.CVSection) {
	experience := ParseExperience(sec.Content)
	doc.AddHeading(title, 1) //nolint:errcheck

	for _, entry := range experience.Entries {
		para := doc.AddEmptyParagraph()
		para.AddText(entry.Position).Bold(true)

		dateRange := formatDateRange(entry.StartDate, entry.EndDate, entry.Current)
		if entry.Company != "" || dateRange != "" {
			details := []string{}
			if entry.Company != "" {
				details = append(details, entry.Company)
			}
			if dateRange != "" {
				details = append(details, dateRange)
			}
			para.AddText(" at " + strings.Join(details, " | "))
		}

		if entry.Description != "" {
			doc.AddParagraph(entry.Description)
		}

		for _, highlight := range entry.Highlights {
			doc.AddParagraph("• " + highlight)
		}
	}
}

func addEducationSection(doc *docx.RootDoc, title string, sec domain.CVSection) {
	education := ParseEducation(sec.Content)
	doc.AddHeading(title, 1) //nolint:errcheck

	for _, entry := range education.Entries {
		para := doc.AddEmptyParagraph()
		para.AddText(entry.Institution).Bold(true)

		df := degreeField(entry.Degree, entry.Field)
		dateRange := formatDateRange(entry.StartDate, entry.EndDate, entry.Current)
		var suffix []string
		if df != "" {
			suffix = append(suffix, df)
		}
		if dateRange != "" {
			suffix = append(suffix, dateRange)
		}
		if len(suffix) > 0 {
			para.AddText(" — " + strings.Join(suffix, " | "))
		}

		if entry.Description != "" {
			doc.AddParagraph(entry.Description)
		}
	}
}

func addSkillsSection(doc *docx.RootDoc, title string, sec domain.CVSection) {
	skills := ParseSkills(sec.Content)
	doc.AddHeading(title, 1) //nolint:errcheck

	for _, category := range skills.Categories {
		para := doc.AddEmptyParagraph()
		para.AddText(category.Name + ": ").Bold(true)
		para.AddText(skillNames(category.Skills))
	}
}

func addLanguagesSection(doc *docx.RootDoc, title string, sec domain.CVSection) {
	languages := ParseLanguages(sec.Content)
	doc.AddHeading(title, 1) //nolint:errcheck

	for _, entry := range languages.Entries {
		text := entry.Language
		if entry.Proficiency != "" {
			text += " — " + entry.Proficiency
		}
		doc.AddParagraph(text)
	}
}
