# T-026: DOCX Export Styling Overhaul

## Summary
- **Goal:** Make DOCX output match the finalized CV template designs (Classic, Modern, Visionary) by applying proper typography, colors, spacing, and layout
- **Acceptance:** DOCX exports for each template look recognizably similar to their HTML/PDF counterparts, with template-specific styling applied
- **Branch:** `feat/t-026-docx-styling-overhaul`

---

## Investigation Checklist
- [x] Current DOCX export implementation — what it does and doesn't do
- [x] Frontend template visual designs — what we need to match
- [x] godocx library capabilities — what's possible
- [x] Domain types — CVContent, CVStyling, template IDs
- [x] HTML renderer reference — how templates are differentiated for PDF
- [x] Gap analysis — what's missing

## Findings

### Current DOCX Export (`api/internal/export/docx.go`, 259 lines)
- **One-size-fits-all**: Ignores `content.TemplateID` entirely — identical output for Classic/Modern/Visionary
- **No styling applied**: `CVStyling` (primaryColor, fontFamily, fontSize, lineHeight) is completely unused
- **Minimal formatting**: Only uses `AddHeading()` (levels 0, 1), `AddParagraph()`, `AddText()` with `.Bold(true)`
- **No colors, fonts, sizes, spacing, alignment, or layout** — just basic text structure
- **Well-structured code**: Clean section handlers, sorted by order, invisible sections excluded
- **10 tests** in `docx_test.go` — verify content presence in `word/document.xml` via ZIP extraction

### Frontend Templates (What We Need to Match)

**Classic** (`ClassicTemplate.tsx`) — Clean, traditional
- Single-column, centered header, horizontal separators
- Section titles: UPPERCASE, semibold, primary color, 1px bottom border
- Name: text-2xl bold, centered; Job title: text-base, secondary color
- Contact: centered, middot (·) separated
- Entry layout: bold title + company + date right-aligned
- Default primary: #c0392b (deep red)

**Modern** (`ModernTemplate.tsx`) — Contemporary with visual flourishes
- Single-column, left-aligned, timeline dots + left border accents
- Section titles: text-lg bold, 3px left border in primary color
- Skill badges: rounded pills with light primary background
- Language proficiency: 5-bar visual indicator
- Timeline: vertical line + dot markers for experience/projects

**Visionary** (`VisionaryTemplate.tsx`) — Two-column with colored sidebar
- **Two-column layout**: 250px sidebar (primary color background, white text) + main content
- Sidebar holds: personal info, skills, languages
- Main content holds: summary, experience, education, certifications, projects
- Section titles: UPPERCASE, primary color, 2px bottom border

### godocx v0.1.5 Capabilities (What's Possible)

**Fully Supported:**
- Text: `Color(hex)`, `Size(halfPoints)`, `Bold()`, `Italic()`, `Underline()`
- Text: `Caps()`, `SmallCaps()`, character `Spacing()`
- Paragraph: `Style()`, `Justification()` (left/center/right/both)
- Paragraph: `Spacing`, `Indent`, `PageBreakBefore()`, `KeepNext()`
- Paragraph: `Shading` (background color), `Border`
- Tables: `AddTable()`, `AddRow()`, `AddCell()`, `Style()` (48+ built-in styles)
- Lists: `Numbering(id, level)`, "List Bullet" / "List Number" styles
- Images: `AddPicture(path, width, height)`
- Metadata: CoreProperties (title, creator, etc.)

**Limitations:**
- **Font family**: `RunFonts` type exists but NOT exposed via public `Run` methods — requires lower-level `RunProperty.Fonts` access
- **Table column widths**: No high-level API — must use lower-level cell width properties
- **Hyperlinks**: Partially stubbed/commented out in library source
- **Page margins**: No high-level API — requires low-level `SectionProp` access
- **Multi-column page layout**: Not supported at high level (but tables can simulate)

### HTML Renderer Reference (`renderer.go`)
- Already has separate `.tmpl` files per template: `classic.tmpl`, `modern.tmpl`, `visionary.tmpl`
- Defines styling constants: `fontStacks`, `fontSizes`, `lineHeights`, `proficiencyLevels`
- Knows about sidebar section types for Visionary: personal, skills, languages
- Default primary color: `#2563eb` (different from frontend's `#c0392b` — frontend wins, user-customizable)

## Gaps Found

### BLOCKER
1. **No template dispatch** — DOCX generator ignores `content.TemplateID`, produces identical output for all templates
2. **No color application** — `CVStyling.PrimaryColor` / `SecondaryColor` not used; all text is default black
3. **No font sizing** — All text uses default Word sizes; no hierarchy between name, section titles, body text

### Medium
4. **No paragraph spacing control** — No spacing between sections, entries, or after headings
5. **No paragraph alignment** — Classic template needs centered header; all templates need proper alignment
6. **No section title styling** — Currently just `AddHeading(title, 1)` with no color/caps/borders
7. **Missing location field** — Experience entries don't include `entry.Location`
8. **Missing profile links** — Personal section doesn't render `Links` (LinkedIn, GitHub, etc.)
9. **Visionary two-column layout** — Requires table-based simulation; godocx tables can do this but it's complex
10. **No font family control** — godocx doesn't expose font family via public API; may need workaround

### Low
11. **Modern timeline visuals** — Cannot render timeline dots/lines in DOCX; best-effort with left borders or indent
12. **Skill badges** — Cannot render rounded pill shapes; best-effort with comma-separated or table cells
13. **Language proficiency bars** — Cannot render visual bars; text-based proficiency display is fine
14. **Highlight bullets** — Currently manual `• ` prefix; should use proper Word bullet list style

## Architect Handoff

### Decisions Needed

**D1: Template-specific rendering scope** (Quick Decision)
- **Option A**: Full template-specific DOCX (3 separate rendering paths) — most faithful, highest effort
- **Option B**: Shared layout with template-specific styling (colors, fonts, spacing vary per template) — balanced
- **Option C**: Single polished layout that applies CVStyling (ignore template differences) — least effort
- **Recommendation**: **Option B** for Classic and Modern (both single-column, differ mainly in styling). **Option C fallback** for Visionary (two-column table layout is complex and fragile). This gives 80% of the visual improvement for 50% of the effort.

**D2: Visionary template approach** (Quick Decision)
- **Option A**: Table-based two-column layout in DOCX (sidebar + main) — faithful but complex, may be fragile
- **Option B**: Single-column layout with sidebar sections grouped at top or bottom — simple, clean
- **Option C**: Skip Visionary-specific layout, use same as Classic/Modern with Visionary colors
- **Recommendation**: **Option B or C** — DOCX two-column tables are notoriously fragile across Word versions. The DOCX is a text-focused format; PDF handles visual fidelity.

**D3: Font family handling** (Quick Decision)
- godocx doesn't expose font family via public Run methods
- **Option A**: Access `RunProperty.Fonts` directly (lower-level, may break on library updates)
- **Option B**: Skip font family, use default Calibri; apply size/color/weight only
- **Option C**: Fork or contribute to godocx to add `Fonts()` method
- **Recommendation**: **Option A** — access the internal property. One helper function wraps it. Low risk for a v0.1.5 library.

**D4: Styling function signature** (Quick Decision)
- Current: `Generate(content domain.CVContent) ([]byte, error)` — styling is already in `content.Styling`
- No interface change needed — `content.Styling` and `content.TemplateID` are already available
- **Recommendation**: No change to interface. Extract styling at top of `Generate()` and pass to section renderers.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `api/internal/export/docx.go` | Major rewrite — add template dispatch, styled section renderers, color/font/spacing application | #1-8, #10-14 |
| `api/internal/export/docx_test.go` | Update tests for styled output; add per-template test cases | All |
| `api/internal/domain/cv.go` | None — types already have everything needed | — |
| `api/internal/export/renderer.go` | May share constants (`fontSizes`, `proficiencyLevels`, `sidebarTypes`) — extract if needed | — |

### Tests to Add/Update
- **Per-template tests**: Generate DOCX for Classic, Modern, Visionary and verify template-specific XML attributes
- **Styling application tests**: Verify primary color appears in document XML, font sizes are set
- **Section title styling**: Verify uppercase, color, spacing in output
- **Profile links**: Test that links are rendered (even as text if hyperlinks aren't supported)
- **Location field**: Test location appears in experience entries
- Existing 10 tests need updating — they verify content presence but not styling

### Constraints
- PR target: < 300 lines changed (may need to split into 2 PRs: "Classic+Modern styling" and "Visionary layout")
- godocx v0.1.5 — don't upgrade unless necessary
- Backward compatible — existing DOCX export tests should still pass (content structure unchanged)
- DOCX is inherently limited vs HTML/PDF — "best effort" is acceptable for complex visual elements

### Recommended Next Agent
**@architect** — Need decisions on D1-D3 before implementation can begin. The template dispatch strategy and Visionary layout approach are design decisions that affect the entire implementation.

---

## Architect Plan

### Decisions

**Decision 1: Template-specific rendering scope**
- Option A: Full template-specific DOCX (3 separate rendering paths) — Pros: most faithful. Cons: 3× the code, high effort, most DOCX-impossible elements are visual (timelines, badges)
- Option B: Shared section renderers + template-specific styling config — Pros: DRY, 80% visual fidelity, manageable diff. Cons: templates look similar in structure
- Option C: Single polished layout ignoring template differences — Pros: minimal code. Cons: no visual distinction between templates
- **Choice:** B with structural exception for Visionary — **Rationale:** Classic and Modern are both single-column; they share section renderers parameterized by a `docxStyle` config. Visionary needs a structural difference (two-column table) on top of styling. The renderers output to a `docWriter` interface so the same code works for both doc-root and table-cell targets.

**Decision 2: Visionary template — two-column table layout via XML post-processing**
- Option A: Fork godocx to add `GetCT()` on Cell/Table/Row — Pros: type-safe API access. Cons: fork maintenance burden, upstream merge not guaranteed.
- Option B: Post-process DOCX XML after generation — Pros: zero dependencies, proven pattern (test infra already unzips+reads XML), full control over any XML property. Cons: string-based XML manipulation (mitigated by predictable godocx output).
- Option C: Paragraph-level shading only (no cell-level) — Pros: no post-processing. Cons: gaps between paragraphs, no cell width control, default table borders.
- Option D: Custom table style via `RootDoc.DocStyles` — Pros: uses godocx's serialization. Cons: OOXML table style conditional formatting is complex, doesn't handle cell widths, arcane.
- **Choice:** B — **Rationale:** The user's principle is "export what we show." Post-processing the DOCX XML gives us cell-level background color, fixed column widths, and borderless table — all without forking godocx. The approach uses only stdlib (`archive/zip`, `strings`, `bytes`). Our test infrastructure already proves the pattern: `readDocumentXML()` unzips and reads `word/document.xml`. We simply extend this to also write. The XML structure is predictable because we control the godocx generation. A dedicated `postProcessVisionary()` function keeps this isolated from the main rendering logic.

**Decision 3: Font family handling**
- Option A: Access `RunProperty.Fonts` via internal struct — Pros: full font control. Cons: `Run.ct` is **private** — CANNOT access from outside the `docx` package.
- Option B: Use Word default (Calibri) — Pros: works, professional font, zero hacks. Cons: ignores user's fontFamily preference.
- **Choice:** B — **Rationale:** `Run.ct` is unexported, making Option A impossible without modifying the library. Calibri is a widely-available, professional sans-serif font that works well for CVs. The fontFamily setting continues to apply in PDF/HTML exports where we have full CSS control. This is an acceptable DOCX limitation.

**Decision 4: Styling function signature**
- **Choice:** No change to `DOCXGenerator` interface — `content.TemplateID` and `content.Styling` are already in `CVContent`. Extract a `docxStyle` struct at the top of `Generate()` and thread it to section renderers as a parameter.

**Decision 5: Visionary post-processing strategy**
- godocx generates the table structure (rows, cells, paragraphs with styled runs)
- A `postProcessVisionary()` function then unzips the DOCX, modifies `word/document.xml` to inject cell/table properties, and re-zips
- **What gets injected:**
  - `<w:tblBorders>` inside `<w:tblPr>` — all borders set to `none`
  - `<w:tcPr>` on first `<w:tc>` (sidebar cell) — width (`<w:tcW>`) + shading (`<w:shd>`)
  - `<w:tcPr>` on second `<w:tc>` (main cell) — borders none (explicit)
- **XML injection approach:** Targeted `strings.Index`-based insertion. The XML is predictable because we control the godocx output. We find specific elements (`<w:tblPr>`, first/second `<w:tc>`) and inject properties at known positions.
- **Isolation:** Post-processing runs ONLY for Visionary template, after all godocx operations complete. Classic/Modern are unaffected.

### Implementation Design

#### Core Type: `docxStyle`

```go
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
    sidebarStyle   bool   // true when rendering inside Visionary sidebar (white text on colored bg)
}
```

#### Font Size Scale (derived from CVStyling.FontSize)

| CVStyling.FontSize | body | meta | title | name |
|---------------------|------|------|-------|------|
| small               | 10   | 9    | 11    | 22   |
| medium (default)    | 11   | 10   | 12    | 24   |
| large               | 12   | 11   | 13    | 26   |

#### Template Style Matrix

| Property | Classic | Modern | Visionary (main) | Visionary (sidebar) |
|----------|---------|--------|------------------|---------------------|
| Header alignment | Center | Left | Left | Left |
| Job title color | secondary | primary | secondary | FFFFFF (white) |
| Section title case | UPPERCASE (Caps) | Normal | UPPERCASE (Caps) | UPPERCASE (Caps) |
| Section title border | Bottom, primary | Left, primary | Bottom, primary | Bottom, white/60 |
| Section title color | primary | primary | primary | FFFFFF |
| Entry title color | default (black) | primary | default (black) | FFFFFF |
| Body/meta text color | secondary | secondary | secondary | FFFFFF (80% opacity N/A in DOCX — use full white) |
| Background | none | none | none | Cell shading = primary color |

#### `docWriter` Abstraction

Section renderers currently take `*docx.RootDoc`. For Visionary, sections render into table cells. Both `RootDoc` and `Cell` can create paragraphs, but their method names differ (`AddEmptyParagraph` vs `AddEmptyPara`). Solution: a thin adapter.

```go
type docWriter struct {
    addPara      func(text string) *docx.Paragraph
    addEmptyPara func() *docx.Paragraph
}

func writerFromDoc(doc *docx.RootDoc) docWriter {
    return docWriter{addPara: doc.AddParagraph, addEmptyPara: doc.AddEmptyParagraph}
}

func writerFromCell(cell *docx.Cell) docWriter {
    return docWriter{addPara: cell.AddParagraph, addEmptyPara: cell.AddEmptyPara}
}
```

All section renderers change signature: `(doc *docx.RootDoc, ...)` → `(w docWriter, ...)`.

#### Visionary Two-Column: Generate + Post-Process

**Step 1: godocx generates the table structure**

```go
// In Generate(), when templateID == "visionary":
table := doc.AddTable()
row := table.AddRow()

sidebarCell := row.AddCell()
mainCell := row.AddCell()

sidebarStyle := style
sidebarStyle.sidebarStyle = true  // switches text colors to white

for _, entry := range ordered {
    if sidebarTypes[entry.section.Type] {
        renderSection(writerFromCell(sidebarCell), sidebarStyle, entry)
    } else {
        renderSection(writerFromCell(mainCell), style, entry)
    }
}
```

Sections use `sidebarTypes` map (already defined in `renderer.go`): personal, skills, languages → sidebar; everything else → main.

**Step 2: Post-process to inject cell/table properties**

After `doc.Write(&buf)`, call `postProcessVisionary()`:

```go
docxBytes := buf.Bytes()
if style.templateID == "visionary" {
    docxBytes, err = postProcessVisionary(docxBytes, style.primaryColor)
    if err != nil {
        return nil, fmt.Errorf("post-processing visionary: %w", err)
    }
}
return docxBytes, nil
```

The `postProcessVisionary()` function:
1. Opens the DOCX as ZIP (`archive/zip`)
2. Reads all files into memory (DOCX is small, typically <100KB)
3. In `word/document.xml`:
   - Finds `<w:tblPr>` → injects `<w:tblBorders>` with all sides `val="none"`
   - Finds first `<w:tc>` → injects `<w:tcPr>` with `<w:tcW w:w="3000" w:type="dxa"/>` (sidebar width ~2.08") and `<w:shd w:val="clear" w:fill="{primaryColor}"/>`
4. Re-zips all files into new buffer

**XML injection targets:**
```xml
<!-- Before: -->
<w:tblPr></w:tblPr>

<!-- After: -->
<w:tblPr>
  <w:tblBorders>
    <w:top w:val="none"/><w:left w:val="none"/>
    <w:bottom w:val="none"/><w:right w:val="none"/>
    <w:insideH w:val="none"/><w:insideV w:val="none"/>
  </w:tblBorders>
</w:tblPr>

<!-- Before: -->
<w:tc><w:p>...

<!-- After: -->
<w:tc>
  <w:tcPr>
    <w:tcW w:w="3000" w:type="dxa"/>
    <w:shd w:val="clear" w:fill="c0392b"/>
  </w:tcPr>
  <w:p>...
```

#### Section Title Rendering

Replace `doc.AddHeading(title, 1)` with manual styled paragraphs. This is necessary because `AddHeading()` applies built-in Word heading styles that override custom formatting (color, caps, borders).

```
Para: w.addEmptyPara()
  → Run: AddText(title).Bold(true).Color(primary or FFFFFF).Size(titleSizePt)
  → Run: .Caps(true) if capsTitle
  → Para props via GetCT(): Border (bottom or left), Spacing.Before = 240 twips (12pt)
```

#### Paragraph Spacing Strategy

Access via `para.GetCT().Property.Spacing`:

| Element | Before (twips) | After (twips) |
|---------|----------------|---------------|
| Section title | 240 (12pt) | 80 (4pt) |
| Entry header (bold line) | 120 (6pt) | 0 |
| Body paragraph | 0 | 40 (2pt) |
| Bullet/highlight | 0 | 0 |

#### Data Gaps Fixed

- **Location field** (Gap #7): Add `entry.Location` to experience entries after company name
- **Profile links** (Gap #8): Render as text lines below contact info: `"LinkedIn: url"`, `"GitHub: url"`

### Code Adjustments

| File | Change | Lines (est.) |
|------|--------|-------------|
| `api/internal/export/docx.go` | Add `docxStyle` struct + `resolveDocxStyle()`. Add `docWriter` abstraction + constructors. Add helpers: `stripHash()`, `getParaProps()`, `addStyledTitle()`. Add `postProcessVisionary()` + ZIP helpers. Modify all 8 section renderers to accept `(w docWriter, style docxStyle, ...)`. Replace `AddHeading()` with manual styled paragraphs. Add Visionary table layout in `Generate()`. Add location + links rendering. New imports for `ctypes`, `stypes`, `archive/zip`, `io`. | ~+230 net |
| `api/internal/export/docx_test.go` | Update `docxContent()` to accept template ID + optional styling. Add per-template tests (Classic, Modern, Visionary) verifying template-specific XML: color hex, caps, border types, table structure + cell shading for Visionary. Add tests for location + links. Existing 10 tests remain. | ~+100 net |

**Estimated total diff: ~+330 net lines** — slightly over 300 target. Split option: Phase 1 PR (styling for all templates, single-column) + Phase 2 PR (Visionary table layout + post-processing).

### Files Explicitly NOT Changing

| File | Reason |
|------|--------|
| `api/internal/domain/cv.go` | Types already have `TemplateID`, `CVStyling`, `Location`, `Links` — no schema changes needed |
| `api/internal/export/renderer.go` | HTML renderer is independent. Constants like `sectionLabels`, `sidebarTypes` are already shared (same package). |
| `api/internal/export/section.go` | Parse functions unchanged — JSON → domain types, no styling |
| `api/internal/service/export_service.go` | Service layer unchanged — calls `docx.Generate(content)`, no signature change |
| `api/internal/handler/export.go` | Handler unchanged — delegates to service |
| `api/go.mod` | No dependency changes — `archive/zip` is stdlib, godocx stays at v0.1.5 |

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| `para.GetCT().Property` returns nil for fresh paragraphs | Medium | Crash | Always nil-check; `getParaProps()` helper handles this. |
| XML injection targets wrong element | Low | Broken DOCX | The XML structure is predictable (we generate it). Tests verify both content AND injected properties. Also: only Visionary triggers post-processing; Classic/Modern are unaffected. |
| Post-processed ZIP loses file metadata (compression, timestamps) | Low | Larger file | Copy original `zip.FileHeader` when re-zipping. Acceptable since DOCX is small (<100KB). |
| Visionary table renders differently in LibreOffice vs Word | Medium | Visual | Test in both. Use standard OOXML properties only (`w:shd`, `w:tcW`, `w:tblBorders`). |
| Cell background doesn't extend to full height in single-row table | Low | Visual | Word auto-extends cells in a single-row table to page height by default. Verify. |
| Existing tests break due to `AddHeading` → manual paragraph switch | Medium | Test failure | Content strings unchanged in XML. `readDocumentXML` checks text content, not XML structure. |
| godocx `Run.Size()` half-point conversion | Low | Wrong sizes | Verified: pass points directly, library handles internally. |

### Consequences

**What becomes easier:**
- DOCX exports match what the user sees in the editor
- User's color customizations carry through to DOCX
- The `docWriter` abstraction makes it easy to render sections into any container
- Post-processing pattern is reusable if we ever need to set other OOXML properties godocx doesn't expose

**What becomes harder:**
- Visionary has a post-processing step that manipulates raw XML — must be kept in sync if godocx output format changes (low risk: pinned at v0.1.5)

### Engineer Execution Steps

**Phase 1: Styling infrastructure**

1. **Add `docxStyle` struct and `resolveDocxStyle()`** at top of `docx.go`. Resolve colors (strip `#`), font sizes (small/medium/large → pt scale), template config (centerHeader, capsTitle, titleBorder, sidebarStyle).

2. **Add `docWriter` struct** with `writerFromDoc()` and `writerFromCell()` constructors.

3. **Add helper functions:**
   - `stripHash(color string) string`
   - `getParaProps(para *docx.Paragraph) *ctypes.ParagraphProp` — nil-safe access via `GetCT()`
   - `addStyledTitle(w docWriter, style docxStyle, title string)` — section title with template styling
   - `setSpacing(para *docx.Paragraph, before, after uint64)` — paragraph spacing helper

4. **Add new imports:** `"github.com/gomutex/godocx/wml/ctypes"`, `"github.com/gomutex/godocx/wml/stypes"`, `"archive/zip"`, `"io"`

**Phase 2: Section renderer refactor**

5. **Change all 8 section renderer signatures** from `(doc *docx.RootDoc, ...)` to `(w docWriter, style docxStyle, ...)`. Replace `doc.AddParagraph` → `w.addPara`, `doc.AddEmptyParagraph` → `w.addEmptyPara`.

6. **Add styling to each renderer:**
   - Replace `doc.AddHeading(title, 1)` with `addStyledTitle(w, style, title)`
   - Replace `doc.AddHeading(fullName, 0)` in personal section with manual sized+colored paragraph
   - Add `.Size()` and `.Color()` calls on all `AddText()` chains
   - Add spacing via `getParaProps()` + `setSpacing()`
   - Personal section: add profile links, center alignment for Classic
   - Experience section: add `entry.Location`
   - Sidebar-aware: when `style.sidebarStyle == true`, use white text colors

**Phase 3: Visionary table layout + post-processing**

7. **In `Generate()`**, add template dispatch:
   - Classic/Modern: `w := writerFromDoc(doc)` → render sections normally → write to buffer → return
   - Visionary: create table → `row.AddCell()` x2 → split sections via `sidebarTypes` → render sidebar sections into first cell (with sidebarStyle), main sections into second cell → write to buffer → post-process → return

8. **Implement `postProcessVisionary(docxBytes []byte, primaryColor string) ([]byte, error)`:**
   - Open ZIP from bytes (`archive/zip`)
   - Read all files into a map
   - Parse `word/document.xml` as string
   - Call `injectTableProperties(xml, primaryColor)` — finds `<w:tblPr>` and injects borderless borders
   - Call `injectSidebarCellProperties(xml, primaryColor)` — finds first `<w:tc>` and injects `<w:tcPr>` with width + shading
   - Write modified files to new ZIP buffer
   - Return new buffer

9. **Implement XML injection helpers:**
   - `injectTableProperties(xml, primaryColor)` — insert `<w:tblBorders>` with all `val="none"` inside `<w:tblPr>`
   - `injectSidebarCellProperties(xml, primaryColor)` — insert `<w:tcPr><w:tcW w:w="3000" w:type="dxa"/><w:shd w:val="clear" w:fill="{color}"/></w:tcPr>` after first `<w:tc>`
   - Both use `strings.Index` to find insertion points. Defensive: return input unchanged if pattern not found.

**Phase 4: Tests**

10. **Update `docxContent()`** to accept template ID parameter
11. **Add template-specific tests:**
    - `TestGenerate_ClassicStyling` — centered header, primary color, caps, bottom border in XML
    - `TestGenerate_ModernStyling` — left-aligned, primary color, no caps, left border in XML
    - `TestGenerate_VisionaryStyling` — table structure (`<w:tbl>`), cell shading (`w:fill="{color}"`), white text (`w:val="FFFFFF"`), sidebar width, caps
    - `TestGenerate_CustomColors` — custom `CVStyling{PrimaryColor: "#FF0000"}`, verify `"FF0000"` in XML
    - `TestGenerate_ProfileLinks` — links rendered as text
    - `TestGenerate_ExperienceLocation` — location appears
12. **Run existing 10 tests** — verify they still pass (content strings unchanged)

**Phase 5: Verification**

13. `cd api && go build ./...`
14. `cd api && go test ./internal/export/...`
15. `cd api && go vet ./...`
16. Manual: generate DOCX for each template, open in Word/LibreOffice, visual sanity check

---

## Test Plan

### Existing Coverage (10 tests in `docx_test.go`)

All verify content presence in `word/document.xml` via ZIP extraction. None verify styling.

| Test | What it covers |
|------|---------------|
| `TestGenerate_ValidDOCX` | ZIP magic bytes, valid structure |
| `TestGenerate_PersonalSection` | Name, job title, contact fields present |
| `TestGenerate_AllSectionTypes` | All 6 main sections produce content |
| `TestGenerate_ExperienceHighlights` | Bullet prefix `•` present |
| `TestGenerate_EmptySections` | Empty input → valid DOCX |
| `TestGenerate_UnknownSectionSkipped` | Unknown type silently skipped |
| `TestGenerate_CertificationsSection` | Cert fields + date formatting |
| `TestGenerate_ProjectsSection` | Project fields + technologies |
| `TestGenerate_InvisibleSectionsExcluded` | Visible=false excluded |
| `TestGenerate_SectionOrdering` | Order field respected |

**Key note:** `docxContent()` hardcodes `TemplateID: "classic"` and never sets `Styling`. When this helper is updated to accept parameters, existing tests must remain unbroken (they should still compile and pass with no styling assertions).

### QA-Recommended Additional Tests

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 1 | `resolveDocxStyle` nil styling | Backend export | High | In `docx_test.go`: call `resolveDocxStyle()` with `CVContent{TemplateID: "classic", Styling: nil}`. Assert defaults: primaryColor = `"2563eb"` (matching `renderer.go` default), bodySizePt = 11, centerHeader = true, capsTitle = true. This is the path ALL existing tests will exercise after the refactor — must not panic. |
| 2 | `resolveDocxStyle` per-template config | Backend export | High | Table-driven test in `docx_test.go` with subtests for each template ID (`classic`, `modern`, `visionary`). For each, assert the template-specific fields: centerHeader (true/false), capsTitle (true/false), titleBorder ("bottom"/"left"). This locks in the Template Style Matrix from the plan. |
| 3 | `resolveDocxStyle` unknown template fallback | Backend export | Medium | Call with `TemplateID: ""` and `TemplateID: "unknown"`. Assert it falls back to sensible defaults (Classic-like) rather than panicking or producing zero-valued style. |
| 4 | `postProcessVisionary` produces valid ZIP | Backend export | High | In `docx_test.go`: generate a Visionary DOCX, then call `postProcessVisionary()` on the output. Verify: (a) result is still valid ZIP with `word/document.xml`, (b) `w:fill="{primaryColor}"` is present, (c) `w:tcW` with width value is present, (d) all `<w:tblBorders>` children have `val="none"`. This tests the post-processing function in isolation. |
| 5 | `postProcessVisionary` defensive — no table | Backend export | Medium | Call `postProcessVisionary()` on a Classic DOCX (which has no `<w:tbl>`). Assert it returns the input unchanged and does not error. This guards against the function being accidentally called for non-Visionary templates. |
| 6 | Font size in XML | Backend export | High | Generate a DOCX with `CVStyling{FontSize: "large"}`. Read XML and verify name font size appears as expected half-point value (e.g., size 26pt → `w:val="52"` in XML). Also verify body text uses the large-scale body size. This confirms the Size Scale table is correctly wired. |
| 7 | Visionary sidebar/main split | Backend export | High | Generate a Visionary DOCX with personal + skills + experience + education sections. Read XML. Verify: (a) personal and skills content appears inside the FIRST `<w:tc>` element, (b) experience and education content appears inside the SECOND `<w:tc>` element. Use `strings.Index` to locate `<w:tc>` boundaries and check content positions relative to them. |
| 8 | Visionary white text in sidebar | Backend export | Medium | Generate a Visionary DOCX with a personal section. Read XML. Verify that within the first `<w:tc>`, run color values contain `"FFFFFF"` (white). This confirms the `sidebarStyle` flag switches text colors. |
| 9 | Classic centered header | Backend export | Medium | Generate a Classic DOCX with personal section. Read XML. Verify `<w:jc w:val="center"/>` appears in the personal section's paragraphs (name, job title, contact). Confirms D1 alignment behavior. |
| 10 | Modern left border on section title | Backend export | Medium | Generate a Modern DOCX with a summary section. Read XML. Verify `<w:pBdr>` contains `<w:left` (not `<w:bottom`) for section title. Confirms the Modern template uses left border accent. |
| 11 | `stripHash` edge cases | Backend export | Medium | Table-driven unit test for `stripHash()`: `"#c0392b"` → `"c0392b"`, `"c0392b"` → `"c0392b"` (no-op), `""` → `""`, `"#"` → `""`. Small but this helper is called on every color value. |
| 12 | Post-process ZIP preserves all files | Backend export | Medium | Generate a Visionary DOCX, run `postProcessVisionary()`. List all files in the output ZIP and verify the same set of files exists as in the input ZIP (word/document.xml, word/styles.xml, [Content_Types].xml, _rels/.rels, etc.). Ensures re-zipping doesn't drop files. |
| 13 | Spacing values on section title | Backend export | Low | Generate any template with a summary section. Read XML. Verify `<w:spacing` with `w:before="240"` appears on the section title paragraph. Locks in the Paragraph Spacing Strategy from the plan. |
| 14 | Caps attribute on Classic/Visionary titles | Backend export | Low | Generate a Classic DOCX with a summary section. Read XML. Verify `<w:caps` appears within the section title run properties. Generate Modern — verify `<w:caps` does NOT appear. Locks in the capsTitle flag behavior. |

---

## QA Visual Comparison (2026-02-23)

Screenshots: `qa-t026-editor.png` (Modern), `qa-t026-visionary.png` (Visionary), `qa-t026-dashboard.png` (dashboard context), `Screenshot 2026-02-23 at 23.25.16.png` (Classic DOCX vs editor side-by-side)

### Classic Template — Flaws Found (from side-by-side DOCX vs Editor)

**What works correctly:**
- Name large, bold, centered ✅
- Job title centered with color ✅
- Contact pipe-separated, centered ✅
- Profile links rendered (labeled format: "linkedin: url") ✅
- Section titles UPPERCASE, primary color, bottom border ✅
- Skills as "Category: skill1, skill2" ✅ (matches Classic editor)
- Languages as "Language — proficiency" ✅

**Experience Section (critical — layout structure wrong):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| C1 | Position + date not two-column | "Senior Software Engineer" left, "Jan 2020 - Present" right-aligned on same line | Everything on one paragraph: "Senior Software Engineer at TechCorp GmbH \| Berlin, Germany \| Jan 2020 - Present" — date inline |
| C2 | Company not on separate line | "at TechCorp GmbH" (primary color) on line 2, "· Berlin, Germany" (secondary) after it | Company inline on same paragraph as position with " at " |
| C3 | Description color wrong | Description in secondary/muted color | Description uses `textColor()` → black instead of secondary |

**Note:** Classic template correctly uses "at" before company/institution (matches editor). The issue is purely layout — it should be on a separate line, not inline.

**Education Section (critical — field order inverted):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| C4 | Degree should be primary bold heading | "M.Sc. in Computer Science" bold on line 1 | `entry.Institution` is bold heading — wrong field first |
| C5 | Institution should be second line | "at Technical University of Munich" (primary color) on line 2 | Institution is the bold heading; degree appended with " — " |
| C6 | Education location missing | "· Munich, Germany" shown after institution | `entry.Location` not rendered |
| C7 | Education date not right-aligned | "Oct 2014 - Mar 2017" right-aligned on heading line | Date inline with " \| " separator |
| C8 | Grade not rendered | "Grade: 1.3" shown below institution line | `entry.Grade` field exists in domain but never used |

### Modern Template — Flaws Found

**Experience Section (critical — layout structure is wrong):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| M1 | Position + date not two-column | "Senior Software Engineer" left, "Jan 2020 - Present" right-aligned on same line | Everything on one line: "Senior Software Engineer at TechCorp GmbH \| Berlin, Germany \| Jan 2020 - Present" — date is inline, not right-aligned |
| M2 | Company not on separate line | "TechCorp GmbH" (blue/primary) on line 2, "| Berlin, Germany" (secondary) after it | Company appended inline with " at " before it on the position line |
| M3 | Spurious "at" conjunction | No "at" — position and company are on separate lines | " at " inserted between position and company |
| M4 | Description color wrong | Description text in secondary/muted color | Description uses `textColor()` → black. Should use `metaColor()` or secondary |

**Education Section (critical — field order is inverted):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| M5 | Degree should be primary bold heading | "M.Sc. in Computer Science" bold on line 1 | `entry.Institution` is bold heading on line 1 — wrong field |
| M6 | Institution should be second line in primary color | "Technical University of Munich" (blue/primary) on line 2 | Institution is the bold heading; degree appended with " — " |
| M7 | Education location missing | "| Munich, Germany" shown after institution | `entry.Location` not rendered in `addEducationSection` |
| M8 | Education date not right-aligned | "Oct 2014 - Mar 2017" right-aligned on heading line | Date appended inline with " \| " separator |
| M9 | Grade not rendered | "Grade: 1.3" shown below institution | `entry.Grade` field exists in domain but `addEducationSection` never uses it |

**Minor / Acceptable Limitations:**

| # | Flaw | Severity | Notes |
|---|------|----------|-------|
| M10 | Profile links on separate lines | Minor | Editor shows both URLs inline on one line; DOCX renders each as separate paragraph |
| M11 | Skill pills → comma list | Acceptable | DOCX can't render rounded pill badges; comma-separated is fine |
| M12 | Timeline dots/line not rendered | Acceptable | DOCX limitation — visual-only elements |
| M13 | Bullet dots not colored | Minor | Editor shows blue dots; DOCX uses plain text "• " |

### Visionary Template — Flaws Found

**Sidebar Personal Section (critical — wrong layout format):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| V1 | Contact fields should be labeled | "Email: max@example.com", "Phone: +49 123 456789", "Location: Berlin, Germany" — each on own line with label | All joined pipe-separated: "email \| phone \| location" — no labels, no separate lines |

**Main Content Header (critical — missing entirely):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| V2 | Name + title missing in main content | "Max Developer" (large, primary/blue) + "Senior Software Engineer" shown as header in the main (right) column | Personal section goes only to sidebar via `sidebarTypes` — main content has no name/title header |

**Sidebar Skills (medium — data loss):**

| # | Flaw | Editor Shows | DOCX Produces |
|---|------|-------------|---------------|
| V3 | Skill proficiency levels dropped | "• Go" with "expert" tag, "• TypeScript" with "advanced" tag, etc. | `skillNames()` only joins `.Name`, drops `.Level` entirely |
| V4 | Skills as bulleted list, not comma-separated | Each skill on own line with bullet: "• Go expert" | "Category: Go, TypeScript, Python" flat comma list |

**Experience Section (same structural flaws as Modern):**

| # | Flaw | Same as |
|---|------|---------|
| V5 | Position + date not two-column | M1 |
| V6 | Company not on separate line in primary color | M2 |
| V7 | Spurious "at" conjunction | M3 |
| V8 | Description color wrong | M4 |

**Education Section (same structural flaws as Modern):**

| # | Flaw | Same as |
|---|------|---------|
| V9 | Degree/institution order inverted | M5 + M6 |
| V10 | Education location missing | M7 |
| V11 | Education date not right-aligned | M8 |
| V12 | Grade not rendered | M9 |

### Summary: All Templates — Flaws by Priority

**P0 — Must Fix (structure/data wrong, affects all 3 templates):**

| ID | Flaw | Templates | Code Location |
|----|------|-----------|---------------|
| F1 | Experience: position + date should be two-column (date right-aligned) | All | `addExperienceSection` — need tab stop or two-paragraph approach |
| F2 | Experience: company should be separate line below position, in primary color | All | `addExperienceSection` — currently inline with " at " |
| F3 | Experience: "at" valid for Classic but not Modern/Visionary | Modern, Visionary | `addExperienceSection` — needs template-aware conjunction |
| F4 | Education: degree should be bold heading (line 1), not institution | All | `addEducationSection` — fields are inverted |
| F5 | Education: institution should be second line in primary color | All | `addEducationSection` — currently bold heading |
| F6 | Education: location not rendered | All | `addEducationSection` — `entry.Location` unused |
| F7 | Education: grade not rendered | All | `addEducationSection` — `entry.Grade` unused |
| F8 | Visionary: main content missing name+title header | Visionary | `Generate()` — personal only goes to sidebar, main has no header |
| F9 | Visionary: sidebar personal should use labeled fields | Visionary | `addPersonalSection` — needs sidebar-aware contact layout |

**P1 — Should Fix (visual fidelity):**

| ID | Flaw | Templates | Code Location |
|----|------|-----------|---------------|
| F10 | Description text should use secondary color, not black | All | `addExperienceSection` — `textColor()` returns black |
| F11 | Education date should be right-aligned | All | `addEducationSection` — same tab stop issue as F1 |
| F12 | Skill proficiency levels dropped | Visionary (sidebar) | `addSkillsSection` — `skillNames()` drops `.Level` |
| F13 | Skills should be bulleted list in Visionary sidebar | Visionary | `addSkillsSection` — needs sidebar-aware rendering |

**P2 — Nice to Have:**

| ID | Flaw | Templates |
|----|------|-----------|
| F14 | Profile links could be on one line | Modern |
| F15 | Bullet dots could use primary color | All |

---

## Architect Re-Plan (2026-02-23)

### QA Findings Classification

**Blocking (must fix):** F1–F13 — structural layout mismatches, missing fields, wrong field order, missing data
**Non-blocking (P2):** F14, F15 — cosmetic improvements, deferred

### Root Cause Analysis

The original architect plan focused on template-specific **styling** (colors, borders, caps, spacing, alignment) and the Visionary two-column layout. It did NOT address the **entry layout structure** within sections — the plan assumed the existing entry rendering (position+company on one line, institution as primary heading) was correct. The QA visual comparison revealed this was wrong: all HTML templates use a multi-line entry structure with right-aligned dates, which the DOCX never matched.

Specifically:
1. **Experience entries** — all 3 HTML templates put Position + Date on line 1, Company + Location on line 2. The DOCX puts everything on one paragraph.
2. **Education entries** — all 3 HTML templates put Degree on line 1 (bold), Institution on line 2 (primary color). The DOCX has these inverted.
3. **Missing fields** — Grade and education Location were never wired into the DOCX renderer.
4. **Visionary gaps** — the original plan covered sidebar/main split and post-processing, but missed the main header duplication and sidebar-specific formatting (labeled contact, bulleted skills with levels).

### Decisions

**Decision 1: Right-aligned dates — tab stops via CT manipulation**
- Option A: Right-aligned tab stop on paragraph + tab character in run via `Paragraph.GetCT().Children` — Pros: matches editor exactly, standard OOXML feature. Cons: requires low-level CT manipulation since `Run` doesn't expose `GetCT()`.
- Option B: Date on separate right-aligned paragraph — Pros: simple. Cons: pushes date to next line, doesn't match editor.
- Option C: Inline date with " — " separator — Pros: simplest. Cons: doesn't match editor at all.
- **Choice:** A — **Rationale:** godocx supports `ParagraphProp.Tabs` with `CustTabStopRight` (confirmed in library source). Tab characters can be inserted by appending `ParagraphChild{Run: &ctypes.Run{Children: [{Tab: &Empty{}}, {Text: &Text{}}]}}` directly to `para.GetCT().Children`. This is the same CT access pattern we already use for spacing and borders.

**Decision 2: Entry layout approach — multi-paragraph per entry**
- All entries restructured to:
  - Para 1: Title (bold) + tab + Date (right-aligned, meta color)
  - Para 2: Subtitle (primary color) + separator + Location (meta color)
  - Para 3+: Description, Grade, Highlights
- Template-specific behavior:
  - Classic: "at " before subtitle, " · " separator
  - Modern: no prefix, " | " separator
  - Visionary: no prefix, " · " separator

**Decision 3: Visionary main header — render before section loop**
- In `Generate()`, for visionary template, extract personal info and render name + job title into the main cell BEFORE iterating sections. The personal section itself still goes to sidebar (via sidebarTypes).
- This mirrors the HTML template's `HeaderPersonal` block.

**Decision 4: Sidebar-aware rendering — conditional branches in addPersonalSection and addSkillsSection**
- `addPersonalSection`: when `sidebarStyle`, render contact as labeled "Email  value" lines instead of pipe-separated
- `addSkillsSection`: when `sidebarStyle`, render each skill as "• name  level" bulleted list instead of "Category: name1, name2"

### Code Adjustments

| File | Change | Fixes |
|------|--------|-------|
| `api/internal/export/docx.go` | Restructure `addExperienceSection` + `addEducationSection` entry layout; add tab stop helper; add Visionary main header in `Generate()`; add sidebar branches in `addPersonalSection` + `addSkillsSection`; fix description color; add `rightTabPos` to docxStyle + methods | F1–F13 |
| `api/internal/export/docx_test.go` | Update content assertions that check for " at " pattern; add tests for right-aligned dates, degree-first layout, grade, education location, Visionary main header | F1–F13 |

### Files NOT Changing (preserved from original implementation)

| File | Reason |
|------|--------|
| `docxStyle` struct + `resolveDocxStyle()` | Working correctly — only adding `rightTabPos` field |
| `docWriter` abstraction | Working correctly |
| `addStyledTitle()` | Working correctly |
| `postProcessVisionary()` + XML injection | Working correctly |
| `addSummarySection()` | Working correctly (F10 color fix is trivial — change `textColor()` → `metaColor()`) |
| `addCertificationsSection()` | No flaws found in QA — keep as-is |
| `addProjectsSection()` | No flaws found in QA — keep as-is |
| All domain types, service layer, handlers | Not in scope |

### Implementation Details

#### 1. Add `rightTabPos` to docxStyle

```go
type docxStyle struct {
    // ... existing fields ...
    rightTabPos int // twips — right tab stop for date alignment
}
```

In `resolveDocxStyle()`:
- Classic/Modern (doc-level paragraphs): `rightTabPos = 9072` (~15.9cm, A4 text width with 1" margins)
- Visionary: `rightTabPos = 5800` (main cell is ~60% of text width after 3000 twip sidebar)

#### 2. Add template-aware helpers to docxStyle

```go
func (s docxStyle) metaSeparator() string {
    if s.templateID == "modern" { return " | " }
    return " · "
}

func (s docxStyle) usesAtPrefix() bool {
    return s.templateID == "classic" || s.templateID == ""
}

func (s docxStyle) descriptionColor() string {
    if s.sidebarStyle { return "FFFFFF" }
    return s.secondaryColor // F10: was textColor() (black), should be secondary
}
```

#### 3. Add tab stop + date helper

```go
func addDateOnRight(para *docx.Paragraph, style docxStyle, dateRange string) {
    if dateRange == "" { return }
    ensureParaProps(para)
    para.GetCT().Property.Tabs = ctypes.Tabs{
        Tab: []ctypes.Tab{{
            Val:      stypes.CustTabStopRight,
            Position: style.rightTabPos,
        }},
    }
    // Append run with tab char + date text via CT
    dateColor := style.metaColor()
    para.GetCT().Children = append(para.GetCT().Children, ctypes.ParagraphChild{
        Run: &ctypes.Run{
            Property: &ctypes.RunProperty{
                // Size and Color need correct ctypes fields — verify during implementation
            },
            Children: []ctypes.RunChild{
                {Tab: &ctypes.Empty{}},
                {Text: &ctypes.Text{Text: dateRange}},
            },
        },
    })
}
```

**Note:** The exact `RunProperty` fields for size/color must be verified from `ctypes.RunProperty` during implementation. The pattern for `Size` is likely `FontSize` or similar; for `Color` it's `ctypes.Color`.

#### 4. Restructure addExperienceSection (F1–F4)

```
For each entry:
  Para 1: Position (bold, entryTitleColor) + [tab] + Date (right-aligned, metaColor)
          → spacing: before=120, after=0
  Para 2: [prefix] + Company (primaryColor) + separator + Location (metaColor)
          → prefix = "at " for Classic, "" for Modern/Visionary
          → separator = " · " for Classic/Visionary, " | " for Modern
          → spacing: before=0, after=0
  Para 3 (if description): Description text (descriptionColor = secondary)
          → spacing: before=0, after=40
  Para 4+ (highlights): "• " + text (descriptionColor)
```

#### 5. Restructure addEducationSection (F4–F9)

```
For each entry:
  Para 1: degreeField (bold, entryTitleColor) + [tab] + Date (right-aligned, metaColor)
          → spacing: before=120, after=0
  Para 2: [prefix] + Institution (primaryColor) + separator + Location (metaColor)
          → prefix = "at " for Classic, "" for Modern/Visionary
          → spacing: before=0, after=0
  Para 3 (if grade): "Grade: " + grade (metaColor)
  Para 4 (if description): Description text (descriptionColor)
```

#### 6. Add Visionary main header (F8)

In `Generate()`, after creating the table but before the section loop:
```go
if style.templateID == "visionary" {
    // ... existing table setup ...

    // Add header to main cell (F8)
    personal := findPersonalSection(ordered)
    if personal != nil {
        p := ParsePersonal(personal.Content)
        fullName := strings.TrimSpace(p.FirstName + " " + p.LastName)
        if fullName != "" {
            namePara := mainCell.AddParagraph(fullName)
            // or mainCell.AddEmptyPara() + AddText(...)
            // Style: large, bold, primary color
        }
        if p.JobTitle != "" {
            // Job title para: secondary color
        }
    }

    // ... existing section loop ...
}
```

Add helper `findPersonalSection()` that returns the personal section from ordered list.

#### 7. Sidebar-aware addPersonalSection (F9)

Add branch when `style.sidebarStyle`:
```
Name (bold, white, large)
Job title (white)
[blank line or spacing]
"Email  " + email  (each on own paragraph, meta color = white)
"Phone  " + phone
"Location  " + location
[links as before]
```

#### 8. Sidebar-aware addSkillsSection (F12, F13)

Add branch when `style.sidebarStyle`:
```
For each category:
  Category name (bold, smaller, metaColor)
  For each skill:
    "• " + skill.Name + "  " + skill.Level  (on own paragraph)
```

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Tab character not rendering correctly in Word | Low | Wrong layout | The `ParagraphChild{Run: &ctypes.Run{Children: [{Tab}, {Text}]}}` pattern matches OOXML spec. If it fails, fall back to "—" separator. |
| RunProperty field names wrong for CT-created runs | Medium | Compile error | Verify exact `ctypes.RunProperty` fields during implementation. Pattern is the same as what `Run.Size()` / `Run.Color()` set internally. |
| Right tab position wrong for Visionary cell | Low | Date misaligned | The 5800 twip value is approximate. Verify visually; adjust if needed. |
| Existing tests break due to layout restructure | Medium | Test failure | Tests check for content strings like "Acme", "Lead", "at". The " at " pattern will change for Modern tests. Update assertions to match new layout. |
| Description color change affects test assertions | Low | Test failure | Existing tests don't check color values for descriptions. New styling tests already exist. |

### Engineer Execution Steps

1. **Read `ctypes.RunProperty`** to find exact field names for Size and Color — needed for the `addDateOnRight` helper

2. **Add `rightTabPos` to `docxStyle`** + set values in `resolveDocxStyle()`:
   - Classic/Modern: 9072
   - Visionary: 5800

3. **Add helper methods** to `docxStyle`: `metaSeparator()`, `usesAtPrefix()`, `descriptionColor()`

4. **Implement `addDateOnRight()`** — tab stop + CT run with tab char + date text

5. **Restructure `addExperienceSection()`**:
   - Para 1: Position + date via `addDateOnRight()`
   - Para 2: [prefix] + Company (primaryColor) + separator + Location (metaColor)
   - Para 3: Description (descriptionColor)
   - Paras 4+: Highlights

6. **Restructure `addEducationSection()`**:
   - Para 1: degreeField (bold) + date via `addDateOnRight()`
   - Para 2: [prefix] + Institution (primaryColor) + separator + Location (metaColor)
   - Para 3: Grade (if present)
   - Para 4: Description

7. **Fix `addSummarySection()`** — change `textColor()` → `descriptionColor()` for body text

8. **Add Visionary main header** in `Generate()`:
   - Extract personal section from ordered list
   - Render name + job title into mainCell before section loop

9. **Add sidebar branch in `addPersonalSection()`** — labeled "Email  value" format

10. **Add sidebar branch in `addSkillsSection()`** — bulleted "• name  level" format

11. **Update tests**:
    - Tests checking for " at " in experience need updating (only Classic keeps "at")
    - Add test for education degree-first ordering
    - Add test for education grade + location
    - Add test for Visionary main header content
    - Add test for tab stop presence in XML (`w:tab`)
    - Add test for description color (secondary, not black)

12. **Verify**: `cd api && go build ./... && go test ./internal/export/... && go vet ./...`

### Verification Gates

1. `go build ./...` — compiles
2. `go vet ./...` — no issues
3. `go test ./internal/export/...` — all tests pass (existing + new)
4. Manual: generate DOCX for Classic, Modern, Visionary → open in Word/LibreOffice → visual check

### Non-blocking Findings (deferred)

| ID | Flaw | Recommendation |
|----|------|----------------|
| F14 | Profile links on one line (Modern) | Low impact — defer to polish pass |
| F15 | Bullet dots with primary color | Would require CT manipulation for colored bullet symbols — defer |
