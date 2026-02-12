# T-002: Create HTML Generation for CV Export

## Summary
- **Goal:** Build a Go service that takes CV content + templateId and produces a self-contained HTML string with inline CSS, one function per template (classic, modern, visionary)
- **Acceptance:** Unit tests pass, HTML renders correctly in browser
- **Branch:** `feat/t-002-html-generation`

---

## Investigation Checklist
- [x] Trace CV domain types (Go) — what does CV content look like?
- [x] Trace existing export handler/service (the 501 stub)
- [x] Trace frontend templates (Classic, Modern, Visionary) — what HTML/CSS do they produce?
- [x] Check HTML template designs in `data/` folder
- [x] Identify JSON schema for CV content
- [x] Map Go project structure — where should the new service live?
- [x] Check Gotenberg API requirements — what HTML format does it expect?
- [x] Identify existing test patterns in Go codebase

---

## Findings

### 1. CV Domain Types (Go)

**File:** `api/internal/domain/cv.go`

Key structures:
- **`CV`** (L11-20): Main entity — `ID`, `UserID`, `Title`, `SchemaVersion`, `Content` (json.RawMessage), `IsActive`
  - `ParseContent()` (L76) decodes Content into `CVContent`
- **`CVContent`** (L24-31): `SchemaVersion`, `TemplateID`, `Locale`, `Title`, `Sections` ([]CVSection), `Styling` (*CVStyling)
- **`CVSection`** (L34-41): `ID`, `Type`, `Order`, `Visible`, `Title`, `Content` (json.RawMessage)
- **`CVStyling`** (L44-51): `PrimaryColor`, `SecondaryColor`, `FontFamily`, `FontSize`, `LineHeight`, `ShowIcons`
- **Section content types** (L96-167): `PersonalContent`, `ExperienceContent` (→`ExperienceEntry`), `EducationContent` (→`EducationEntry`), `SkillsContent` (→`SkillCategory`→`Skill`), `LanguagesContent` — no dedicated type yet, languages are likely bare json.RawMessage
- **Template constants** (L71-73): `TemplateClassic = "classic"`, `TemplateModern = "modern"`, `TemplateVisionary = "visionary"`
- **Section type constants** (L54-67): personal, summary, experience, education, skills, languages, certifications, projects, awards, publications, references, custom

**Note:** `CVSection.Content` is `json.RawMessage` — needs type-specific unmarshaling per section type.

### 2. Existing Export Stub

**File:** `api/internal/handler/export.go` (L1-51)

- `CreateExport()` — extracts `format` from URL param, validates with `domain.IsValidExportFormat()`, returns **501 Not Implemented**
- `GetExport()` — also 501 stub
- Response type: `ExportResponse` (ID, Format, Status, URL, Error, CreatedAt)
- Route: `r.Post("/export/{format}", h.CreateExport)` in `api/cmd/server/main.go:225`

**File:** `api/internal/domain/asset.go` (L59-97)

- `ExportJob` struct with ID, UserID, CVID, Format, Status, ResultPath, ErrorMessage, timestamps
- Valid formats: pdf, docx, json, yaml
- Statuses: pending, processing, completed, failed

**File:** `api/internal/repository/repository.go` (L81-93)

- `ExportJobRepository` interface defined but **not implemented** — GetByID, ListPending, Create, UpdateStatus
- SQL queries exist: `api/db/queries/export_jobs.sql` (L1-38)

### 3. Frontend Templates (React)

Three templates producing different HTML/CSS:

| Template | File | Lines | Layout |
|----------|------|-------|--------|
| Classic | `web/src/components/templates/ClassicTemplate.tsx` | 429 | Single column, centered, serif accents, bottom-border sections |
| Modern | `web/src/components/templates/ModernTemplate.tsx` | 423 | Single column, left-border sections, timeline dots for experience, badge-style skills, bar-chart languages |
| Visionary | `web/src/components/templates/VisionaryTemplate.tsx` | 521 | Two-column (250px sidebar + flex main), sidebar has personal/skills/languages on colored background |

**Styling approach:** Tailwind CSS utility classes + inline styles for dynamic colors (primaryColor, secondaryColor).

**Default colors:** primaryColor `#2563eb`, secondaryColor `#64748b`

**Key visual differences:**
- **Classic**: Section titles with bottom border, centered personal header, comma-separated skills, middot-separated languages
- **Modern**: Section titles with 3px left border, timeline with dots for experience, rounded-full pill badges for skills, 5-bar progress chart for languages
- **Visionary**: Colored sidebar (250px, bg=primaryColor) with white text, main content with 2px bottom-border section titles, bullet-dot highlights

### 4. JSON Schema

**Canonical (embedded in Go validator):** `api/internal/validator/schema/cv-schema.json`
- templateId enum: `["classic", "modern", "visionary"]` ✅
- Required: `schemaVersion`, `templateId`, `sections`
- 12 section types with type-specific content schemas
- Styling: primaryColor, secondaryColor, fontFamily, fontSize, lineHeight, showIcons

**Stale copy (project root):** `schemas/cv-schema.json`
- Has `"minimal"` instead of `"visionary"` — OUT OF DATE ⚠️

### 5. Gotenberg Configuration

**Docker:** `docker/docker-compose.infra.yml` (L29-44)
- Image: `gotenberg/gotenberg:8`
- Port: `3001:3000`
- **JavaScript disabled** (`--chromium-disable-javascript=true`)
- File allow-list: `file:///tmp/.*`
- API timeout: 60s

**Go config:** `api/internal/config/config.go` (L64-66, L131-132)
- `ExportConfig.GotenbergURL` defaults to `http://localhost:3001`

**Implication:** HTML must be fully self-contained — no external JS, no CDN CSS. All styles must be inline or in a `<style>` block.

### 6. Go Project Patterns

**Service pattern** (`api/internal/service/cv_service.go`):
- Constructor with dependency injection: `NewCVService(repos, validator)`
- Context as first param
- Returns domain types + errors
- Registered in `main.go` L81-83

**Handler pattern** (`api/internal/handler/cv.go`):
- Uses `middleware.MustGetUserID(ctx)` for auth
- `httputil.DecodeJSON()`, `httputil.JSON()`, `httputil.HandleError()`
- Request/response structs with JSON tags

**Test pattern** (`api/internal/handler/testutil_test.go`):
- Mock services with function pointers
- `httptest.NewRecorder()` for handler tests
- Test data helpers: `createTestUser()`, `createTestCV()`

**Dependency injection** (`api/internal/handler/handler.go`):
- `Dependencies` struct holds all services
- Handler receives deps in `New()`

### 7. Seed Data (Sample CV)

**File:** `api/db/seed/dev_seed.sql` (L36-192)
- Full CV with sections: personal, summary, experience (2 entries), education (1), skills (2 categories), languages (2)
- templateId: "modern", styling with primaryColor "#2563eb"
- Good test fixture reference

---

## Gaps Found

### BLOCKER
1. **No `LanguageEntry` Go type** — Experience, Education, Skills have dedicated content types in `domain/cv.go`, but languages content is not explicitly typed. Need to check if it's handled via generic json.RawMessage unmarshaling or if a type needs to be added.

### Medium
2. **Stale `schemas/cv-schema.json`** — Root-level schema has `"minimal"` instead of `"visionary"`. Embedded copy in `api/internal/validator/schema/` is correct. Should be synced.
3. **No summary/languages content types in Go** — `SummaryContent` and `LanguagesContent` may need to be defined (or confirmed as already parseable from json.RawMessage).

### Low
5. **`blank` template in frontend but not in Go** — Frontend has `TemplateId = 'blank'` but Go only defines classic/modern/visionary. Not relevant for export (blank isn't exportable), but worth noting.

---

## Architect Handoff

### Decisions Needed

**D1: HTML generation approach** (Quick Decision)
- **Option A: Go `html/template`** — Use Go's standard template engine with `.tmpl` files. Pros: familiar, easy to maintain, natural Go pattern. Cons: template syntax can be verbose for complex HTML.
- **Option B: Go string builder** — Construct HTML programmatically with `strings.Builder` or `fmt.Sprintf`. Pros: full control, easy to test individual parts. Cons: harder to read, mixing logic and markup.
- **Option C: Embedded HTML templates with Go `text/template`** — Similar to A but allows more flexible string output. Pros: clean separation. Cons: no auto-escaping (but we control all data).
- **Recommendation: Option A** — `html/template` with embedded `.tmpl` files via `//go:embed`. Matches existing patterns (validator uses `//go:embed` for schema). Auto-escaping protects against XSS in CV data.

**D2: CSS approach** (Quick Decision)
- **Option A: Inline styles on every element** — `style="..."` on each HTML tag. Pros: maximum portability, works everywhere. Cons: verbose, hard to maintain.
- **Option B: `<style>` block in `<head>`** — CSS classes in a style block. Pros: cleaner HTML, reusable classes. Cons: some PDF renderers strip `<style>`.
- **Option C: Hybrid** — `<style>` for base layout + inline for dynamic colors. Pros: clean + portable. Cons: slightly more complex.
- **Recommendation: Option C** — Gotenberg uses Chromium which fully supports `<style>` blocks. Use classes for layout/typography, inline only for dynamic colors (primaryColor, secondaryColor).

**D3: Where to put the new code** (Quick Decision)
- New service: `api/internal/service/export_service.go`
- Template files: `api/internal/service/templates/*.tmpl` (embedded via `//go:embed`)
- Or dedicated package: `api/internal/export/` with `renderer.go` + `templates/`
- **Recommendation: `api/internal/export/`** — keeps template rendering isolated from service orchestration. ExportService in `service/` calls into `export.Renderer`.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `api/internal/export/renderer.go` | **NEW** — HTML renderer with `Render(content CVContent) (string, error)` | Core task |
| `api/internal/export/renderer_test.go` | **NEW** — Unit tests per template | Core task |
| `api/internal/export/templates/classic.tmpl` | **NEW** — Classic HTML template | Core task |
| `api/internal/export/templates/modern.tmpl` | **NEW** — Modern HTML template | Core task |
| `api/internal/export/templates/visionary.tmpl` | **NEW** — Visionary HTML template | Core task |
| `api/internal/export/templates/base.tmpl` | **NEW** — Shared HTML boilerplate (head, fonts, reset CSS) | Core task |
| `api/internal/domain/cv.go` | **MODIFY** — Add missing content types if needed (LanguageEntry, SummaryContent) | Gap #1, #4 |
| `schemas/cv-schema.json` | **MODIFY** — Sync templateId enum to match embedded copy | Gap #2 |

### Constraints
- HTML must be **self-contained** — no external resources (Gotenberg has JS disabled, file-only allow-list)
- Fonts: either use web-safe fonts or embed font data via base64 in `<style>`
- Output should visually match React template rendering as closely as possible
- PR target: < 300 lines (may be tight — consider splitting templates into separate commits or accepting slightly over)

### Recommended Next Agent
**@architect** — Several design decisions needed (D1-D4) before implementation. After architect decides, **@engineer** implements.

---

## Test Plan

### Unit Tests (`api/internal/export/renderer_test.go`)

1. **Per-template rendering** — For each template (classic, modern, visionary):
   - Render a full CV (all section types) → verify HTML contains expected elements
   - Render minimal CV (personal only) → verify no errors, basic structure present
   - Verify correct CSS classes/styles for that template's visual identity

2. **Section-specific tests:**
   - Personal section: name, job title, contact info, links
   - Summary section: text content
   - Experience section: entries with position, company, dates, highlights
   - Education section: entries with degree, institution, dates, grade
   - Skills section: categories with skill names
   - Languages section: entries with language name and proficiency

3. **Edge cases:**
   - Empty sections (visible but no content)
   - Hidden sections (`visible: false`) should be omitted
   - Missing optional fields (no phone, no links, no highlights)
   - Custom styling (non-default colors, fonts)
   - Special characters in CV data (HTML escaping)

4. **HTML validity:**
   - Output is valid HTML5 (contains `<!DOCTYPE html>`, `<html>`, `<head>`, `<body>`)
   - Self-contained (no external resource references)
   - CSS variables / inline styles use correct color values from styling

### Test Data
Use seed data structure from `api/db/seed/dev_seed.sql` as reference fixture. Create a `testdata/` directory with sample `CVContent` JSON.

### Integration Test (optional, for T-003)
Open generated HTML in browser and visually verify — can be manual for T-002, automated in T-003 when Gotenberg is wired up.

---

## Architect Plan

### Phase 2: Research & Compatibility

**Go `html/template` with `embed`** — The validator already uses `//go:embed` for the schema JSON (`api/internal/validator/cv_schema.go:14`). This is a proven pattern in the codebase. `html/template` is stdlib, no dependency additions needed.

**Gotenberg Chromium compatibility** — Gotenberg 8 uses Chromium with JavaScript disabled. `<style>` blocks in `<head>` are fully supported. No need for inline-only CSS. External resources (CDN fonts, external CSS) are blocked by the file-only allow-list.

**Font strategy** — Using web-safe font stacks (system fonts). The styling `fontFamily` field maps to CSS stacks: `inter` → system-ui stack, `merriweather` → serif stack, etc. No base64 font embedding needed for MVP — Gotenberg's Chromium has standard system fonts.

### Phase 3: Decisions

**Decision 1: HTML generation approach**
- Option A: `html/template` with embedded `.tmpl` files — Pros: auto-escaping (XSS protection on CV data), clean separation of markup and logic, familiar Go pattern, `.tmpl` files match `//go:embed` pattern already used. Cons: template syntax can be verbose for conditionals.
- Option B: `strings.Builder` / programmatic — Pros: full control, easy to test fragments. Cons: mixing markup and logic, no auto-escaping, harder to read/maintain.
- **Choice: A** — **Rationale:** Auto-escaping is important since CV data is user-input. The validator already proves `//go:embed` works in this project. Template files are easier to review and maintain than string builders.

**Decision 2: CSS approach**
- Option A: All inline styles — Pros: maximum portability. Cons: extremely verbose, hard to maintain, duplicated across templates.
- Option B: `<style>` block only — Pros: clean, reusable classes. Cons: some email clients strip styles (irrelevant here).
- Option C: Hybrid — `<style>` for base layout/typography + inline for dynamic colors — Pros: clean + portable. Cons: slightly more complex.
- **Choice: C** — **Rationale:** Gotenberg uses Chromium which fully supports `<style>`. Use classes for layout, inline only for dynamic values (`primaryColor`, `secondaryColor`). This keeps templates readable while supporting user-customizable colors.

**Decision 3: Package structure**
- Option A: `api/internal/service/export_service.go` + templates in same package — Pros: simpler. Cons: mixes template rendering with service orchestration.
- Option B: `api/internal/export/` package with `renderer.go` + `templates/` — Pros: clean separation, export package owns HTML generation, service layer orchestrates. Cons: one more package.
- **Choice: B** — **Rationale:** Keeps rendering isolated. `export.Renderer` is a pure function (CVContent → HTML string) with no service dependencies. The export service (T-003) will call into it. This also keeps the `templates/` directory co-located with the renderer.

**Decision 4: Missing Go types (Gap #1, #3)**
- `SummaryContent` — Schema requires `{ text: string }`. Go domain lacks this type. Need to add.
- `LanguagesContent` / `LanguageEntry` — Schema requires `{ entries: [{ language, proficiency, certification? }] }`. Go domain lacks these types. Need to add.
- **Choice:** Add `SummaryContent`, `LanguagesContent`, and `LanguageEntry` to `api/internal/domain/cv.go`. Keep them consistent with existing patterns (`ExperienceContent` → `ExperienceEntry` pattern).

**Decision 5: Template scope — which section types to render**
- The frontend renders 6 section types: personal, summary, experience, education, skills, languages.
- The schema defines 12 types total, but 6 are not yet implemented in the frontend (certifications, projects, awards, publications, references, custom).
- **Choice:** Match the frontend — render the same 6 types. Unimplemented types get a generic fallback (section title only). This avoids scope creep and keeps T-002 as a size S task. Future tasks (T-011, T-012) will add certifications and projects to both frontend and export.

**Decision 6: Stale root schema (Gap #2)**
- **Choice:** Out of scope for T-002. It's a documentation/sync issue, not a runtime issue. The authoritative schema is in `api/internal/validator/schema/`. File a note for a future cleanup.

### Phase 4: Implementation Plan

#### Target Structure

```
api/internal/export/
├── renderer.go           # Render(CVContent) → (string, error)
├── renderer_test.go      # Unit tests
├── section.go            # Section content parsing helpers
├── templates/
│   ├── base.tmpl         # Shared HTML boilerplate (doctype, head, reset CSS, font stacks)
│   ├── classic.tmpl      # Classic template
│   ├── modern.tmpl       # Modern template
│   └── visionary.tmpl    # Visionary template
```

#### Files Changing

| File | Action | What Changes |
|------|--------|-------------|
| `api/internal/domain/cv.go` | MODIFY | Add `SummaryContent`, `LanguagesContent`, `LanguageEntry` types |
| `api/internal/export/renderer.go` | NEW | `Renderer` struct, `Render(CVContent) (string, error)`, template loading via `//go:embed`, section content parsing |
| `api/internal/export/section.go` | NEW | Helper functions to parse `json.RawMessage` → typed content per section type |
| `api/internal/export/renderer_test.go` | NEW | Unit tests for all 3 templates |
| `api/internal/export/templates/base.tmpl` | NEW | HTML5 boilerplate, CSS reset, font-family map, shared utility classes |
| `api/internal/export/templates/classic.tmpl` | NEW | Classic template: centered header, bottom-border sections, comma-separated skills |
| `api/internal/export/templates/modern.tmpl` | NEW | Modern template: left-border sections, timeline dots, pill badges, bar-chart languages |
| `api/internal/export/templates/visionary.tmpl` | NEW | Visionary template: two-column layout (sidebar + main), colored sidebar |

#### Files NOT Changing

| File | Reason |
|------|--------|
| `api/internal/handler/export.go` | T-003 will wire this up; T-002 only builds the renderer |
| `api/internal/handler/handler.go` | No new service injection needed yet |
| `api/cmd/server/main.go` | No routing changes — renderer is not exposed via HTTP yet |
| `schemas/cv-schema.json` | Out of scope (stale copy, not used at runtime) |
| `web/src/components/templates/*` | Read-only reference; no frontend changes |

#### Domain Type Additions (`api/internal/domain/cv.go`)

```go
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
    Language      string `json:"language"`
    Proficiency   string `json:"proficiency"`
    Certification string `json:"certification,omitempty"`
}
```

#### Renderer API (`api/internal/export/renderer.go`)

```go
package export

type Renderer struct {
    templates map[string]*template.Template // keyed by template ID
}

func NewRenderer() (*Renderer, error) // loads & parses embedded templates
func (r *Renderer) Render(content domain.CVContent) (string, error) // dispatches to correct template
```

Key design points:
- `NewRenderer()` is called once at startup (or in tests); parses all templates
- `Render()` selects template by `content.TemplateID`, executes into a `bytes.Buffer`, returns string
- Template data struct wraps `CVContent` with pre-parsed section content (typed, not `json.RawMessage`)
- Each `.tmpl` file `{{template "base" .}}` defines a base layout, then template-specific blocks

#### Template Data Model

The renderer pre-processes `CVContent` into a `TemplateData` struct before passing to templates:

```go
type TemplateData struct {
    Personal   *domain.PersonalContent
    Summary    *domain.SummaryContent
    Experience *domain.ExperienceContent
    Education  *domain.EducationContent
    Skills     *domain.SkillsContent
    Languages  *domain.LanguagesContent
    Sections   []SectionData  // ordered, visible-only, with parsed content
    Styling    StylingData    // with defaults applied
}

type SectionData struct {
    Type    string
    Title   string
    Content interface{} // typed content
}

type StylingData struct {
    PrimaryColor   string // default: #2563eb
    SecondaryColor string // default: #64748b
    FontFamily     string // CSS font-stack string
}
```

#### CSS Strategy

**Base template (`base.tmpl`)** provides:
- HTML5 doctype, `<html>`, `<head>`, `<body>`
- CSS reset (margin, padding, box-sizing)
- Font-family mapping: `inter` → `"Inter", system-ui, -apple-system, sans-serif`, etc.
- Typography scale: base sizes for `h1`, `h2`, `p`, `li`
- Print-friendly defaults: `@page { margin: 0; }`, no backgrounds stripped

**Per-template CSS** is in each `.tmpl` file's `<style>` block:
- Layout-specific rules (e.g., Visionary's two-column flex)
- Section title styling differences (Classic: bottom-border, Modern: left-border, Visionary: bottom-border + uppercase)
- Component styling (Modern's timeline dots, skill pills, language bars)

**Dynamic colors** are injected via inline `style=""` attributes in the HTML, using Go template variables like `{{.Styling.PrimaryColor}}`.

#### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| HTML doesn't match React templates visually | Medium | Medium | Use React templates as direct reference; visual diff in browser during QA |
| Template parsing errors at runtime | Low | High | `NewRenderer()` fails fast at startup; unit tests verify all templates parse |
| Missing section content type causes panic | Low | High | `section.go` returns zero-value structs for nil content; tests cover empty sections |
| Template output too large (> 300 lines PR) | Medium | Low | Templates share base; keep CSS minimal; if needed, split into 2 PRs (renderer + templates) |

#### Engineer Execution Steps

1. **Add missing domain types** — Edit `api/internal/domain/cv.go`:
   - Add `SummaryContent`, `LanguagesContent`, `LanguageEntry` after `SkillsContent` (L167)
   - Keep consistent with existing patterns

2. **Create `api/internal/export/` directory** and files:
   - `section.go` — Section content parsing helpers:
     - `ParsePersonal(json.RawMessage) *domain.PersonalContent`
     - `ParseSummary(json.RawMessage) *domain.SummaryContent`
     - `ParseExperience(json.RawMessage) *domain.ExperienceContent`
     - `ParseEducation(json.RawMessage) *domain.EducationContent`
     - `ParseSkills(json.RawMessage) *domain.SkillsContent`
     - `ParseLanguages(json.RawMessage) *domain.LanguagesContent`
     - Each returns zero-value struct on parse failure (never nil)

3. **Create `templates/base.tmpl`** — Shared HTML boilerplate:
   - `<!DOCTYPE html>`, `<html lang="en">`, `<head>`, `<body>`
   - CSS reset, font-family map, typography base
   - Defines `{{block "content" .}}{{end}}` for template-specific content
   - A4 dimensions: `width: 210mm` for print/PDF

4. **Create `templates/classic.tmpl`** — Translate React `ClassicTemplate.tsx`:
   - Centered personal header with name, job title, contact middot-separated
   - Section titles: uppercase, bottom-border with primaryColor
   - Experience: position at company · location | date range; description; highlight bullets
   - Education: degree in field at institution · location | date range; grade
   - Skills: category name: skill1, skill2, skill3 (comma-separated)
   - Languages: language (proficiency) middot-separated inline

5. **Create `templates/modern.tmpl`** — Translate React `ModernTemplate.tsx`:
   - Left-aligned personal header, job title in primaryColor
   - Section titles: left-border (3px) with primaryColor
   - Experience: timeline with dots (circles on vertical line), company in primaryColor
   - Skills: rounded pill badges with primaryColor background (15% opacity)
   - Languages: 5-bar chart (filled vs empty based on proficiency level)

6. **Create `templates/visionary.tmpl`** — Translate React `VisionaryTemplate.tsx`:
   - Two-column layout: sidebar (250px, bg=primaryColor, white text) + main content
   - Sidebar: personal info, skills (bullet dots), languages (text)
   - Main: header strip with name + job title, then sections
   - Section titles: uppercase, bottom-border (2px) with primaryColor

7. **Create `renderer.go`** — Main renderer:
   - `//go:embed templates/*.tmpl` for all template files
   - `NewRenderer()` parses base + each template
   - `Render(CVContent)` builds `TemplateData` from CVContent, executes template
   - Template functions: `formatDate(string) string`, `joinStrings([]string, string) string`, proficiency-to-level mapping

8. **Create `renderer_test.go`** — Unit tests:
   - Test fixture: `CVContent` struct matching seed data structure
   - Per-template tests (classic, modern, visionary):
     - Full CV render → verify HTML contains name, job title, company, skills, etc.
     - Verify HTML structure: `<!DOCTYPE html>`, `<html>`, `<head>`, `<body>`
     - Verify self-contained (no external URLs in `<link>` or `<script>`)
     - Verify dynamic colors appear in output
   - Edge case tests:
     - Hidden sections (visible=false) omitted
     - Empty sections (no entries) don't crash
     - Missing optional fields (no phone, no links)
     - Special characters in data are HTML-escaped
     - Unknown template ID returns error
   - Minimal CV test (personal section only)

9. **Run verification gates:**
   - `cd api && go build ./...` — compiles
   - `cd api && go test ./internal/export/...` — all tests pass
   - `cd api && go vet ./...` — no vet issues

#### Verification Gates

| Gate | Command | Must Pass |
|------|---------|-----------|
| Compile | `cd api && go build ./...` | Yes |
| Unit tests | `cd api && go test ./internal/export/... -v` | Yes |
| Existing tests | `cd api && go test ./...` | Yes (no regressions) |
| Vet | `cd api && go vet ./...` | Yes |

#### Consequences

**What becomes easier:**
- T-003 (PDF export) and T-004 (DOCX export) can call `export.Renderer.Render()` directly — they just need to POST the HTML to Gotenberg
- Adding new section types (T-011 certifications, T-012 projects) just means adding a new `Parse*` function and template blocks

**What becomes harder:**
- Nothing significant. The export package is isolated and has no dependencies on the service layer.
