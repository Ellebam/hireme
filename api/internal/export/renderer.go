package export

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"sort"
	"strings"
	"time"

	"github.com/ellebam/hireme/api/internal/domain"
)

//go:embed templates/*.tmpl
var templateFS embed.FS

// Default styling values
const (
	defaultPrimaryColor   = "#2563eb"
	defaultSecondaryColor = "#64748b"
	defaultFontFamily     = "inter"
	defaultFontSize       = "medium"
	defaultLineHeight     = "normal"
)

// Font family CSS stacks
var fontStacks = map[string]string{
	"inter":        `"Inter", system-ui, -apple-system, "Segoe UI", sans-serif`,
	"roboto":       `"Roboto", "Helvetica Neue", Arial, sans-serif`,
	"opensans":     `"Open Sans", "Helvetica Neue", Arial, sans-serif`,
	"lato":         `"Lato", "Helvetica Neue", Arial, sans-serif`,
	"merriweather": `"Merriweather", Georgia, "Times New Roman", serif`,
}

// Font size to px mapping
var fontSizes = map[string]string{
	"small":  "13px",
	"medium": "14px",
	"large":  "16px",
}

// Line height mapping
var lineHeights = map[string]string{
	"compact": "1.4",
	"normal":  "1.6",
	"relaxed": "1.8",
}

// Proficiency level mapping (for Modern template's bar chart)
var proficiencyLevels = map[string]int{
	"native":       5,
	"fluent":       4,
	"advanced":     3,
	"intermediate": 2,
	"basic":        1,
}

// Section type labels (matching frontend SECTION_LABELS)
var sectionLabels = map[string]string{
	"personal":       "Personal",
	"summary":        "Professional Summary",
	"experience":     "Work Experience",
	"education":      "Education",
	"skills":         "Skills",
	"languages":      "Languages",
	"certifications": "Certifications",
	"projects":       "Projects",
	"awards":         "Awards",
	"publications":   "Publications",
	"references":     "References",
	"custom":         "Custom",
}

// Sidebar section types for Visionary template
var sidebarTypes = map[string]bool{
	"personal":  true,
	"skills":    true,
	"languages": true,
}

// TemplateData is the top-level data passed to templates.
type TemplateData struct {
	Title           string
	Sections        []SectionData
	SidebarSections []SectionData // Visionary only
	MainSections    []SectionData // Visionary only
	HeaderPersonal  *PersonalData // Visionary only
	Styling         StylingData
}

// SectionData wraps a parsed section with typed content.
type SectionData struct {
	Type       string
	Title      string
	Personal   *PersonalData
	Summary    *domain.SummaryContent
	Experience *domain.ExperienceContent
	Education  *domain.EducationContent
	Skills     *domain.SkillsContent
	Languages  *domain.LanguagesContent
}

// PersonalData wraps personal content with pre-computed fields.
type PersonalData struct {
	Content      domain.PersonalContent
	FullName     string
	ContactParts []string
}

// StylingData holds resolved styling values.
type StylingData struct {
	PrimaryColor   string
	SecondaryColor string
	FontStack      template.CSS
	FontSizePx     string
	LineHeightVal  string
}

// Renderer generates self-contained HTML from CV content.
type Renderer struct {
	templates map[string]*template.Template
}

// NewRenderer creates a Renderer with all templates parsed and ready.
func NewRenderer() (*Renderer, error) {
	funcMap := template.FuncMap{
		"dateRange":       formatDateRange,
		"formatDate":      formatDate,
		"degreeField":     degreeField,
		"skillNames":      skillNames,
		"linkLabel":       linkLabel,
		"proficiencyBars": proficiencyBars,
	}

	r := &Renderer{
		templates: make(map[string]*template.Template),
	}

	for _, name := range []string{"classic", "modern", "visionary"} {
		tmpl, err := template.New("").Funcs(funcMap).ParseFS(
			templateFS,
			"templates/base.tmpl",
			fmt.Sprintf("templates/%s.tmpl", name),
		)
		if err != nil {
			return nil, fmt.Errorf("parsing template %s: %w", name, err)
		}
		r.templates[name] = tmpl
	}

	return r, nil
}

// Render generates a self-contained HTML string for the given CV content.
func (r *Renderer) Render(content domain.CVContent) (string, error) {
	tmpl, ok := r.templates[content.TemplateID]
	if !ok {
		return "", fmt.Errorf("unknown template: %q", content.TemplateID)
	}

	data := buildTemplateData(content)

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
		return "", fmt.Errorf("executing template %s: %w", content.TemplateID, err)
	}

	return buf.String(), nil
}

// buildTemplateData converts CVContent into the TemplateData used by templates.
func buildTemplateData(content domain.CVContent) TemplateData {
	styling := resolveStyling(content.Styling)

	// Build title from personal section if available
	title := content.Title
	if title == "" {
		title = "CV"
	}

	// Parse visible sections, sorted by order
	type orderedSection struct {
		order int
		data  SectionData
	}
	var ordered []orderedSection
	for _, sec := range content.Sections {
		if !sec.Visible {
			continue
		}
		ordered = append(ordered, orderedSection{order: sec.Order, data: parseSectionData(sec)})
	}
	sort.Slice(ordered, func(i, j int) bool {
		return ordered[i].order < ordered[j].order
	})
	allSections := make([]SectionData, len(ordered))
	for i, o := range ordered {
		allSections[i] = o.data
	}

	// Split for Visionary template
	var sidebar, main []SectionData
	var headerPersonal *PersonalData
	for _, sd := range allSections {
		if sidebarTypes[sd.Type] {
			sidebar = append(sidebar, sd)
			if sd.Type == "personal" && sd.Personal != nil {
				headerPersonal = sd.Personal
			}
		} else {
			main = append(main, sd)
		}
	}

	return TemplateData{
		Title:           title,
		Sections:        allSections,
		SidebarSections: sidebar,
		MainSections:    main,
		HeaderPersonal:  headerPersonal,
		Styling:         styling,
	}
}

func parseSectionData(sec domain.CVSection) SectionData {
	title := sec.Title
	if title == "" {
		if label, ok := sectionLabels[sec.Type]; ok {
			title = label
		} else {
			title = sec.Type
		}
	}

	sd := SectionData{
		Type:  sec.Type,
		Title: title,
	}

	switch sec.Type {
	case domain.SectionTypePersonal:
		p := ParsePersonal(sec.Content)
		fullName := strings.TrimSpace(p.FirstName + " " + p.LastName)
		var contactParts []string
		if p.Email != "" {
			contactParts = append(contactParts, p.Email)
		}
		if p.Phone != "" {
			contactParts = append(contactParts, p.Phone)
		}
		if p.Location != "" {
			contactParts = append(contactParts, p.Location)
		}
		sd.Personal = &PersonalData{
			Content:      p,
			FullName:     fullName,
			ContactParts: contactParts,
		}
	case domain.SectionTypeSummary:
		s := ParseSummary(sec.Content)
		sd.Summary = &s
	case domain.SectionTypeExperience:
		e := ParseExperience(sec.Content)
		sd.Experience = &e
	case domain.SectionTypeEducation:
		e := ParseEducation(sec.Content)
		sd.Education = &e
	case domain.SectionTypeSkills:
		s := ParseSkills(sec.Content)
		sd.Skills = &s
	case domain.SectionTypeLanguages:
		l := ParseLanguages(sec.Content)
		sd.Languages = &l
	}

	return sd
}

func resolveStyling(s *domain.CVStyling) StylingData {
	primary := defaultPrimaryColor
	secondary := defaultSecondaryColor
	family := defaultFontFamily
	size := defaultFontSize
	lh := defaultLineHeight

	if s != nil {
		if s.PrimaryColor != "" {
			primary = s.PrimaryColor
		}
		if s.SecondaryColor != "" {
			secondary = s.SecondaryColor
		}
		if s.FontFamily != "" {
			family = s.FontFamily
		}
		if s.FontSize != "" {
			size = s.FontSize
		}
		if s.LineHeight != "" {
			lh = s.LineHeight
		}
	}

	fontStack := fontStacks[family]
	if fontStack == "" {
		fontStack = fontStacks[defaultFontFamily]
	}

	fontSize := fontSizes[size]
	if fontSize == "" {
		fontSize = fontSizes[defaultFontSize]
	}

	lineHeight := lineHeights[lh]
	if lineHeight == "" {
		lineHeight = lineHeights[defaultLineHeight]
	}

	return StylingData{
		PrimaryColor:   primary,
		SecondaryColor: secondary,
		FontStack:      template.CSS(fontStack),
		FontSizePx:     fontSize,
		LineHeightVal:  lineHeight,
	}
}

// Template functions

func formatDate(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	parts := strings.Split(dateStr, "-")
	if len(parts) < 2 {
		return parts[0] // year only
	}
	t, err := time.Parse("2006-01", dateStr)
	if err != nil {
		return dateStr
	}
	return t.Format("Jan 2006")
}

func formatDateRange(startDate string, endDate *string, current bool) string {
	start := formatDate(startDate)
	var end string
	if current {
		end = "Present"
	} else if endDate != nil {
		end = formatDate(*endDate)
	}
	parts := []string{}
	if start != "" {
		parts = append(parts, start)
	}
	if end != "" {
		parts = append(parts, end)
	}
	return strings.Join(parts, " - ")
}

func degreeField(degree, field string) string {
	parts := []string{}
	if degree != "" {
		parts = append(parts, degree)
	}
	if field != "" {
		parts = append(parts, "in "+field)
	}
	return strings.Join(parts, " ")
}

func skillNames(skills []domain.Skill) string {
	names := make([]string, len(skills))
	for i, s := range skills {
		names[i] = s.Name
	}
	return strings.Join(names, ", ")
}

func linkLabel(link domain.ProfileLink) string {
	if link.Label != "" {
		return link.Label
	}
	if link.URL != "" {
		return link.URL
	}
	return link.Type
}

func proficiencyBars(proficiency, primaryColor string) []template.CSS {
	level := proficiencyLevels[strings.ToLower(proficiency)]
	bars := make([]template.CSS, 5)
	for i := range 5 {
		if i < level {
			bars[i] = template.CSS(primaryColor)
		} else {
			bars[i] = template.CSS(primaryColor + "20")
		}
	}
	return bars
}
