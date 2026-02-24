package export

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/gomutex/godocx"
	"github.com/gomutex/godocx/docx"
	"github.com/gomutex/godocx/wml/ctypes"
	"github.com/gomutex/godocx/wml/stypes"

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

// docxStyle holds resolved template-specific styling for DOCX generation.
type docxStyle struct {
	primaryColor   string // hex without "#", e.g. "c0392b"
	secondaryColor string
	templateID     string
	nameSizePt     uint64 // full name font size in points
	titleSizePt    uint64 // section title font size
	bodySizePt     uint64 // body text font size
	metaSizePt     uint64 // metadata (dates, locations) font size
	centerHeader   bool   // true for Classic
	capsTitle      bool   // true for Classic + Visionary
	titleBorder    string // "bottom" for Classic/Visionary, "left" for Modern
	sidebarStyle   bool   // true when rendering inside Visionary sidebar (white text)
	rightTabPos    int    // twips — right tab stop for date alignment
}

func (s docxStyle) textColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	return "000000"
}

func (s docxStyle) titleColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	return s.primaryColor
}

func (s docxStyle) metaColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	return s.secondaryColor
}

func (s docxStyle) jobTitleColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	if s.templateID == "modern" {
		return s.primaryColor
	}
	return s.secondaryColor
}

func (s docxStyle) entryTitleColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	if s.templateID == "modern" {
		return s.primaryColor
	}
	return "000000"
}

func (s docxStyle) borderColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	return s.primaryColor
}

func (s docxStyle) descriptionColor() string {
	if s.sidebarStyle {
		return "FFFFFF"
	}
	return s.secondaryColor
}

func (s docxStyle) metaSeparator() string {
	if s.templateID == "modern" {
		return " | "
	}
	return " · "
}

func (s docxStyle) usesAtPrefix() bool {
	return s.templateID == "classic" || s.templateID == ""
}

func resolveDocxStyle(content domain.CVContent) docxStyle {
	primaryColor := defaultPrimaryColor
	secondaryColor := defaultSecondaryColor
	fontSize := defaultFontSize

	if content.Styling != nil {
		if content.Styling.PrimaryColor != "" {
			primaryColor = content.Styling.PrimaryColor
		}
		if content.Styling.SecondaryColor != "" {
			secondaryColor = content.Styling.SecondaryColor
		}
		if content.Styling.FontSize != "" {
			fontSize = content.Styling.FontSize
		}
	}

	var bodySizePt, metaSizePt, titleSizePt, nameSizePt uint64
	switch fontSize {
	case "small":
		bodySizePt, metaSizePt, titleSizePt, nameSizePt = 10, 9, 11, 22
	case "large":
		bodySizePt, metaSizePt, titleSizePt, nameSizePt = 12, 11, 13, 26
	default: // "medium" or unknown
		bodySizePt, metaSizePt, titleSizePt, nameSizePt = 11, 10, 12, 24
	}

	templateID := content.TemplateID
	centerHeader := true
	capsTitle := true
	titleBorder := "bottom"

	switch templateID {
	case "modern":
		centerHeader = false
		capsTitle = false
		titleBorder = "left"
	case "visionary":
		centerHeader = false
	default: // "classic" or unknown
		// defaults are classic-like
	}

	rightTabPos := 9072 // A4 text width with 1" margins (~15.9cm)
	if templateID == "visionary" {
		rightTabPos = 5800 // narrower main cell
	}

	return docxStyle{
		primaryColor:   sanitizeHexColor(primaryColor, strings.TrimPrefix(defaultPrimaryColor, "#")),
		secondaryColor: sanitizeHexColor(secondaryColor, strings.TrimPrefix(defaultSecondaryColor, "#")),
		templateID:     templateID,
		nameSizePt:     nameSizePt,
		titleSizePt:    titleSizePt,
		bodySizePt:     bodySizePt,
		metaSizePt:     metaSizePt,
		centerHeader:   centerHeader,
		capsTitle:      capsTitle,
		titleBorder:    titleBorder,
		rightTabPos:    rightTabPos,
	}
}

// docWriter abstracts paragraph creation so section renderers work with
// both RootDoc (single-column) and Cell (Visionary table).
type docWriter struct {
	addPara      func(text string) *docx.Paragraph
	addEmptyPara func() *docx.Paragraph
}

func writerFromDoc(doc *docx.RootDoc) docWriter {
	return docWriter{
		addPara:      doc.AddParagraph,
		addEmptyPara: doc.AddEmptyParagraph,
	}
}

func writerFromCell(cell *docx.Cell) docWriter {
	return docWriter{
		addPara:      cell.AddParagraph,
		addEmptyPara: cell.AddEmptyPara,
	}
}

var hexColorRe = regexp.MustCompile(`^[0-9a-fA-F]{3,8}$`)

// sanitizeHexColor strips a leading "#" and validates the value is a hex color.
// Returns the fallback if the color is invalid.
func sanitizeHexColor(color, fallback string) string {
	color = strings.TrimPrefix(color, "#")
	if hexColorRe.MatchString(color) {
		return color
	}
	return fallback
}

func ensureParaProps(para *docx.Paragraph) {
	ct := para.GetCT()
	if ct.Property == nil {
		ct.Property = &ctypes.ParagraphProp{}
	}
}

func setSpacing(para *docx.Paragraph, before, after uint64) {
	ensureParaProps(para)
	para.GetCT().Property.Spacing = &ctypes.Spacing{
		Before: &before,
		After:  &after,
	}
}

func addDateOnRight(para *docx.Paragraph, style docxStyle, dateRange string) {
	if dateRange == "" {
		return
	}
	ensureParaProps(para)
	para.GetCT().Property.Tabs = ctypes.Tabs{
		Tab: []ctypes.Tab{{
			Val:      stypes.CustTabStopRight,
			Position: style.rightTabPos,
		}},
	}
	halfPts := style.metaSizePt * 2
	dateColor := style.metaColor()
	para.GetCT().Children = append(para.GetCT().Children, ctypes.ParagraphChild{
		Run: &ctypes.Run{
			Property: &ctypes.RunProperty{
				Size:  &ctypes.FontSize{Value: halfPts},
				Color: &ctypes.Color{Val: dateColor},
			},
			Children: []ctypes.RunChild{
				{Tab: &ctypes.Empty{}},
				{Text: &ctypes.Text{Text: dateRange}},
			},
		},
	})
}

func addStyledTitle(w docWriter, style docxStyle, title string) {
	para := w.addEmptyPara()
	run := para.AddText(title).Bold(true).Size(style.titleSizePt).Color(style.titleColor())
	if style.capsTitle {
		run.Caps(true)
	}

	setSpacing(para, 240, 80)

	props := para.GetCT().Property // guaranteed non-nil after setSpacing
	borderColor := style.borderColor()
	switch style.titleBorder {
	case "bottom":
		spaceVal := "1"
		props.Border = &ctypes.ParaBorder{
			Bottom: &ctypes.Border{
				Val:   stypes.BorderStyleSingle,
				Color: &borderColor,
				Space: &spaceVal,
			},
		}
	case "left":
		spaceVal := "4"
		props.Border = &ctypes.ParaBorder{
			Left: &ctypes.Border{
				Val:   stypes.BorderStyleSingle,
				Color: &borderColor,
				Space: &spaceVal,
			},
		}
	}
}

type orderedSection struct {
	order   int
	section domain.CVSection
}

func findPersonalSection(ordered []orderedSection) *domain.CVSection {
	for _, entry := range ordered {
		if entry.section.Type == domain.SectionTypePersonal {
			sec := entry.section
			return &sec
		}
	}
	return nil
}

func sectionTitle(sec domain.CVSection) string {
	title := sec.Title
	if title == "" {
		if label, ok := sectionLabels[sec.Type]; ok {
			title = label
		}
	}
	return title
}

func renderSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	switch sec.Type {
	case domain.SectionTypePersonal:
		addPersonalSection(w, style, sec)
	case domain.SectionTypeSummary:
		addSummarySection(w, style, title, sec)
	case domain.SectionTypeExperience:
		addExperienceSection(w, style, title, sec)
	case domain.SectionTypeEducation:
		addEducationSection(w, style, title, sec)
	case domain.SectionTypeSkills:
		addSkillsSection(w, style, title, sec)
	case domain.SectionTypeLanguages:
		addLanguagesSection(w, style, title, sec)
	case domain.SectionTypeCertifications:
		addCertificationsSection(w, style, title, sec)
	case domain.SectionTypeProjects:
		addProjectsSection(w, style, title, sec)
	}
}

// Generate produces a DOCX document from structured CV content.
func (g *GodocxGenerator) Generate(content domain.CVContent) ([]byte, error) {
	doc, err := godocx.NewDocument()
	if err != nil {
		return nil, fmt.Errorf("creating document: %w", err)
	}

	style := resolveDocxStyle(content)

	// Collect visible sections sorted by order
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

	if style.templateID == "visionary" {
		table := doc.AddTable()
		row := table.AddRow()
		sidebarCell := row.AddCell()
		mainCell := row.AddCell()

		sidebarDocStyle := style
		sidebarDocStyle.sidebarStyle = true

		// Visionary: add name + job title to main cell header
		personalSec := findPersonalSection(ordered)
		if personalSec != nil {
			personal := ParsePersonal(personalSec.Content)
			fullName := strings.TrimSpace(personal.FirstName + " " + personal.LastName)
			if fullName != "" {
				namePara := mainCell.AddEmptyPara()
				namePara.AddText(fullName).Bold(true).Size(style.nameSizePt).Color(style.titleColor())
				setSpacing(namePara, 0, 40)
			}
			if personal.JobTitle != "" {
				titlePara := mainCell.AddEmptyPara()
				titlePara.AddText(personal.JobTitle).Size(style.bodySizePt).Color(style.jobTitleColor())
				setSpacing(titlePara, 0, 80)
			}
		}

		for _, entry := range ordered {
			sec := entry.section
			title := sectionTitle(sec)
			if sidebarTypes[sec.Type] {
				renderSection(writerFromCell(sidebarCell), sidebarDocStyle, title, sec)
			} else {
				renderSection(writerFromCell(mainCell), style, title, sec)
			}
		}
	} else {
		w := writerFromDoc(doc)
		for _, entry := range ordered {
			sec := entry.section
			title := sectionTitle(sec)
			renderSection(w, style, title, sec)
		}
	}

	var buf bytes.Buffer
	if err := doc.Write(&buf); err != nil {
		return nil, fmt.Errorf("writing document: %w", err)
	}

	docxBytes := buf.Bytes()
	if style.templateID == "visionary" {
		docxBytes, err = postProcessVisionary(docxBytes, style.primaryColor)
		if err != nil {
			return nil, fmt.Errorf("post-processing visionary: %w", err)
		}
	}

	return docxBytes, nil
}

func addPersonalSection(w docWriter, style docxStyle, sec domain.CVSection) {
	personal := ParsePersonal(sec.Content)

	fullName := strings.TrimSpace(personal.FirstName + " " + personal.LastName)
	if fullName != "" {
		para := w.addEmptyPara()
		para.AddText(fullName).Bold(true).Size(style.nameSizePt).Color(style.textColor())
		if style.centerHeader {
			para.Justification(stypes.JustificationCenter)
		}
		setSpacing(para, 0, 40)
	}

	if personal.JobTitle != "" {
		para := w.addEmptyPara()
		para.AddText(personal.JobTitle).Size(style.bodySizePt).Color(style.jobTitleColor())
		if style.centerHeader {
			para.Justification(stypes.JustificationCenter)
		}
		setSpacing(para, 0, 40)
	}

	if style.sidebarStyle {
		// Sidebar: labeled fields on separate lines
		if personal.Email != "" {
			para := w.addEmptyPara()
			para.AddText("Email  ").Bold(true).Size(style.metaSizePt).Color(style.metaColor())
			para.AddText(personal.Email).Size(style.metaSizePt).Color(style.metaColor())
		}
		if personal.Phone != "" {
			para := w.addEmptyPara()
			para.AddText("Phone  ").Bold(true).Size(style.metaSizePt).Color(style.metaColor())
			para.AddText(personal.Phone).Size(style.metaSizePt).Color(style.metaColor())
		}
		if personal.Location != "" {
			para := w.addEmptyPara()
			para.AddText("Location  ").Bold(true).Size(style.metaSizePt).Color(style.metaColor())
			para.AddText(personal.Location).Size(style.metaSizePt).Color(style.metaColor())
		}
	} else {
		// Standard: pipe-separated contact line
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
			para := w.addEmptyPara()
			para.AddText(strings.Join(contactParts, " | ")).Size(style.metaSizePt).Color(style.metaColor())
			if style.centerHeader {
				para.Justification(stypes.JustificationCenter)
			}
			setSpacing(para, 0, 40)
		}
	}

	// Profile links
	for _, link := range personal.Links {
		if link.URL == "" {
			continue
		}
		text := link.URL
		if link.Label != "" {
			text = link.Label + ": " + link.URL
		} else if link.Type != "" {
			text = link.Type + ": " + link.URL
		}
		para := w.addEmptyPara()
		para.AddText(text).Size(style.metaSizePt).Color(style.metaColor())
		if style.centerHeader {
			para.Justification(stypes.JustificationCenter)
		}
	}
}

func addSummarySection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	summary := ParseSummary(sec.Content)
	addStyledTitle(w, style, title)
	if summary.Text != "" {
		para := w.addEmptyPara()
		para.AddText(summary.Text).Size(style.bodySizePt).Color(style.descriptionColor())
		setSpacing(para, 0, 40)
	}
}

func addExperienceSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	experience := ParseExperience(sec.Content)
	addStyledTitle(w, style, title)

	for _, entry := range experience.Entries {
		// Line 1: Position (bold) + right-aligned date
		positionPara := w.addEmptyPara()
		positionPara.AddText(entry.Position).Bold(true).Size(style.bodySizePt).Color(style.entryTitleColor())
		dateRange := formatDateRange(entry.StartDate, entry.EndDate, entry.Current)
		addDateOnRight(positionPara, style, dateRange)
		setSpacing(positionPara, 120, 0)

		// Line 2: [prefix] Company (primary) + separator + Location (meta)
		if entry.Company != "" || entry.Location != "" {
			compPara := w.addEmptyPara()
			if entry.Company != "" {
				prefix := ""
				if style.usesAtPrefix() {
					prefix = "at "
				}
				compPara.AddText(prefix + entry.Company).Size(style.bodySizePt).Color(style.titleColor())
			}
			if entry.Location != "" {
				sep := ""
				if entry.Company != "" {
					sep = style.metaSeparator()
				}
				compPara.AddText(sep + entry.Location).Size(style.metaSizePt).Color(style.metaColor())
			}
			setSpacing(compPara, 0, 0)
		}

		// Description
		if entry.Description != "" {
			descPara := w.addEmptyPara()
			descPara.AddText(entry.Description).Size(style.bodySizePt).Color(style.descriptionColor())
			setSpacing(descPara, 0, 40)
		}

		// Highlights
		for _, highlight := range entry.Highlights {
			hlPara := w.addEmptyPara()
			hlPara.AddText("• " + highlight).Size(style.bodySizePt).Color(style.descriptionColor())
		}
	}
}

func addEducationSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	education := ParseEducation(sec.Content)
	addStyledTitle(w, style, title)

	for _, entry := range education.Entries {
		// Line 1: Degree+Field (bold) + right-aligned date; fall back to Institution if no degree
		df := degreeField(entry.Degree, entry.Field)
		heading := df
		if heading == "" {
			heading = entry.Institution
		}

		headingPara := w.addEmptyPara()
		headingPara.AddText(heading).Bold(true).Size(style.bodySizePt).Color(style.entryTitleColor())
		dateRange := formatDateRange(entry.StartDate, entry.EndDate, entry.Current)
		addDateOnRight(headingPara, style, dateRange)
		setSpacing(headingPara, 120, 0)

		// Line 2: Institution (primary) + separator + Location (meta)
		// Only when degree was shown as heading
		if df != "" && (entry.Institution != "" || entry.Location != "") {
			instPara := w.addEmptyPara()
			if entry.Institution != "" {
				prefix := ""
				if style.usesAtPrefix() {
					prefix = "at "
				}
				instPara.AddText(prefix + entry.Institution).Size(style.bodySizePt).Color(style.titleColor())
			}
			if entry.Location != "" {
				sep := ""
				if entry.Institution != "" {
					sep = style.metaSeparator()
				}
				instPara.AddText(sep + entry.Location).Size(style.metaSizePt).Color(style.metaColor())
			}
			setSpacing(instPara, 0, 0)
		} else if df == "" && entry.Location != "" {
			// Institution was the heading — show location on its own line
			locPara := w.addEmptyPara()
			locPara.AddText(entry.Location).Size(style.metaSizePt).Color(style.metaColor())
			setSpacing(locPara, 0, 0)
		}

		// Grade
		if entry.Grade != "" {
			gradePara := w.addEmptyPara()
			gradePara.AddText("Grade: " + entry.Grade).Size(style.metaSizePt).Color(style.metaColor())
		}

		// Description
		if entry.Description != "" {
			descPara := w.addEmptyPara()
			descPara.AddText(entry.Description).Size(style.bodySizePt).Color(style.descriptionColor())
			setSpacing(descPara, 0, 40)
		}
	}
}

func addSkillsSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	skills := ParseSkills(sec.Content)
	addStyledTitle(w, style, title)

	if style.sidebarStyle {
		// Sidebar: bulleted list with proficiency levels
		for _, category := range skills.Categories {
			catPara := w.addEmptyPara()
			catPara.AddText(category.Name).Bold(true).Size(style.metaSizePt).Color(style.metaColor())
			for _, skill := range category.Skills {
				skillPara := w.addEmptyPara()
				text := "• " + skill.Name
				if skill.Level != "" {
					text += "  " + skill.Level
				}
				skillPara.AddText(text).Size(style.bodySizePt).Color(style.textColor())
			}
		}
	} else {
		// Standard: comma-separated
		for _, category := range skills.Categories {
			para := w.addEmptyPara()
			para.AddText(category.Name + ": ").Bold(true).Size(style.bodySizePt).Color(style.textColor())
			para.AddText(skillNames(category.Skills)).Size(style.bodySizePt).Color(style.textColor())
		}
	}
}

func addLanguagesSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	languages := ParseLanguages(sec.Content)
	addStyledTitle(w, style, title)

	for _, entry := range languages.Entries {
		para := w.addEmptyPara()
		text := entry.Language
		if entry.Proficiency != "" {
			text += " — " + entry.Proficiency
		}
		para.AddText(text).Size(style.bodySizePt).Color(style.textColor())
	}
}

func addCertificationsSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	certs := ParseCertifications(sec.Content)
	addStyledTitle(w, style, title)

	for _, entry := range certs.Entries {
		para := w.addEmptyPara()
		para.AddText(entry.Name).Bold(true).Size(style.bodySizePt).Color(style.entryTitleColor())

		var details []string
		if entry.Issuer != "" {
			details = append(details, entry.Issuer)
		}
		if entry.Date != "" {
			details = append(details, certDateRange(entry.Date, entry.ExpiryDate))
		}
		if len(details) > 0 {
			para.AddText(" — " + strings.Join(details, " | ")).Size(style.metaSizePt).Color(style.metaColor())
		}
		setSpacing(para, 120, 0)

		if entry.CredentialID != "" {
			credPara := w.addEmptyPara()
			credPara.AddText("ID: " + entry.CredentialID).Size(style.metaSizePt).Color(style.metaColor())
		}
	}
}

func addProjectsSection(w docWriter, style docxStyle, title string, sec domain.CVSection) {
	projects := ParseProjects(sec.Content)
	addStyledTitle(w, style, title)

	for _, entry := range projects.Entries {
		para := w.addEmptyPara()
		para.AddText(entry.Name).Bold(true).Size(style.bodySizePt).Color(style.entryTitleColor())

		var details []string
		if entry.Role != "" {
			details = append(details, entry.Role)
		}
		dateRange := formatDateRange(entry.StartDate, entry.EndDate, false)
		if dateRange != "" {
			details = append(details, dateRange)
		}
		if len(details) > 0 {
			para.AddText(" — " + strings.Join(details, " | ")).Size(style.metaSizePt).Color(style.metaColor())
		}
		setSpacing(para, 120, 0)

		if entry.Description != "" {
			descPara := w.addEmptyPara()
			descPara.AddText(entry.Description).Size(style.bodySizePt).Color(style.textColor())
			setSpacing(descPara, 0, 40)
		}

		if len(entry.Technologies) > 0 {
			techPara := w.addEmptyPara()
			techPara.AddText("Technologies: ").Bold(true).Size(style.metaSizePt).Color(style.metaColor())
			techPara.AddText(strings.Join(entry.Technologies, ", ")).Size(style.metaSizePt).Color(style.metaColor())
		}
	}
}

// postProcessVisionary modifies the generated DOCX XML to add Visionary-specific
// table properties: borderless table and sidebar cell shading with fixed width.
func postProcessVisionary(docxBytes []byte, primaryColor string) ([]byte, error) {
	zipReader, err := zip.NewReader(bytes.NewReader(docxBytes), int64(len(docxBytes)))
	if err != nil {
		return nil, fmt.Errorf("opening zip: %w", err)
	}

	type zipEntry struct {
		header *zip.FileHeader
		data   []byte
	}
	var entries []zipEntry
	for _, file := range zipReader.File {
		rc, err := file.Open()
		if err != nil {
			return nil, fmt.Errorf("opening %s: %w", file.Name, err)
		}
		data, err := io.ReadAll(rc)
		if closeErr := rc.Close(); closeErr != nil {
			return nil, fmt.Errorf("closing %s: %w", file.Name, closeErr)
		}
		if err != nil {
			return nil, fmt.Errorf("reading %s: %w", file.Name, err)
		}

		if file.Name == "word/document.xml" {
			xmlStr := string(data)
			xmlStr = injectTableBorders(xmlStr)
			xmlStr = injectSidebarCellProps(xmlStr, primaryColor)
			data = []byte(xmlStr)
		}

		entries = append(entries, zipEntry{header: &file.FileHeader, data: data})
	}

	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)
	for _, entry := range entries {
		writer, err := zipWriter.CreateHeader(entry.header)
		if err != nil {
			return nil, fmt.Errorf("creating zip entry %s: %w", entry.header.Name, err)
		}
		if _, err := writer.Write(entry.data); err != nil {
			return nil, fmt.Errorf("writing zip entry %s: %w", entry.header.Name, err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("closing zip: %w", err)
	}

	return buf.Bytes(), nil
}

// injectTableBorders adds borderless borders to the first table in the XML.
func injectTableBorders(xml string) string {
	target := "</w:tblPr>"
	idx := strings.Index(xml, target)
	if idx == -1 {
		return xml
	}

	borders := `<w:tblBorders>` +
		`<w:top w:val="none"/>` +
		`<w:left w:val="none"/>` +
		`<w:bottom w:val="none"/>` +
		`<w:right w:val="none"/>` +
		`<w:insideH w:val="none"/>` +
		`<w:insideV w:val="none"/>` +
		`</w:tblBorders>`

	return xml[:idx] + borders + xml[idx:]
}

// injectSidebarCellProps adds width and background shading to the first table cell.
func injectSidebarCellProps(xml string, primaryColor string) string {
	target := "<w:tc>"
	idx := strings.Index(xml, target)
	if idx == -1 {
		return xml
	}

	tcPr := `<w:tcPr>` +
		`<w:tcW w:w="3000" w:type="dxa"/>` +
		`<w:shd w:val="clear" w:color="auto" w:fill="` + primaryColor + `"/>` +
		`</w:tcPr>`

	insertPos := idx + len(target)
	return xml[:insertPos] + tcPr + xml[insertPos:]
}
