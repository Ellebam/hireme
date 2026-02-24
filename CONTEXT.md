# HireMe — Project Context

## Current State

**Status:** Task-based workflow active — T-001–T-012, T-016–T-020, T-023–T-026, T-030, T-032–T-034 done, 12 tasks in backlog (3 P1, 9 P2; 10 unblocked, 2 blocked, 4 need splitting when picked up)

### What's Working

**Backend:**
- ✅ Infrastructure: PostgreSQL + Gotenberg via Docker
- ✅ Database: Migrations applied, dev data seeded
- ✅ API endpoints: Health, Users, CVs (full CRUD), Assets
- ✅ Auth bypass: Dev mode uses `dev-user-001`
- ✅ Repository layer: sqlc queries wired to domain types
- ✅ Asset management: Upload, retrieve, delete with local storage
- ✅ Storage abstraction: LocalStorage implemented, R2 placeholder ready
- ✅ API response format: `{ data: ... }` wrapper with proper error handling
- ✅ Go tests: Domain, service, handler, middleware, repository, storage layers tested
- ✅ CI/CD: GitHub Actions with lint, test, build for both API and web

**Frontend:**
- ✅ Three-column editor layout (palette | preview | properties)
- ✅ Dashboard page with CV list
- ✅ 8 section editors (Personal, Summary, Experience, Education, Skills, Languages, Certifications, Projects)
- ✅ Drag & drop section reordering (dnd-kit)
- ✅ Live A4 CV preview with zoom (50-200%)
- ✅ Undo/redo with 50-step history
- ✅ Auto-save with 2s debounce + manual save button
- ✅ Export modal (PDF, DOCX, JSON)
- ✅ Keyboard shortcuts (Ctrl+Z/Y, Ctrl+/-/0)
- ✅ CV templates system (starter + blank templates)
- ✅ Logging system with configurable levels
- ✅ Save error feedback — tooltip on error icon, click-to-retry
- ✅ Dashboard is root page (`/`), `/dashboard` redirects
- ✅ Edit button links to `/editor` (not `/editor/${cv.id}`)
- ✅ Burger menu (MoreVertical) has DropdownMenu with Edit/Delete
- ✅ Delete confirmation dialog on CV card
- ✅ Unified header across all pages (Dashboard, Editor, Templates)
- ✅ shadcn/ui dropdown-menu component
- ✅ 200 passing tests (Vitest)
- ✅ Build verified, full stack integration tested

### Remaining Work
- ✅ Export: PDF (T-003), DOCX (T-004), frontend wiring (T-005) — all working
- ❌ R2 cloud storage — placeholder only, needs AWS SDK implementation → T-013

---

## Tech Stack

| Component | Technology |
|-----------|------------|
| Backend | Go 1.22+, Chi router, sqlc |
| Frontend | Next.js 15, React 19, Tailwind, shadcn/ui, Zustand, dnd-kit |
| Database | PostgreSQL 16 (JSONB for CV content) |
| Export | Gotenberg (HTML → PDF), godocx (CV data → DOCX) |
| Storage | Local filesystem (dev), Cloudflare R2 (prod) |
| Testing | Vitest (frontend), Go testing (backend - planned) |

---

## Frontend Architecture

### Editor Layout
```
+------------------------------------------------------------------+
|                         Top Toolbar                               |
|  [Undo] [Redo] | Zoom [50-200%] | [Export v]                     |
+------------------------------------------------------------------+
|           |                              |                        |
|  Section  |      CV Preview              |    Properties          |
|  Palette  |      (Live A4 Render)        |    Panel               |
|  (240px)  |      (Flex)                  |    (320px)             |
|           |                              |                        |
+------------------------------------------------------------------+
```

### State Management (Zustand)
- **EditorStore**: CV data, sections, undo/redo history, save status
- **UIStore**: Sidebar toggles, modals, preview scale

### Key Frontend Files
```
web/src/
├── stores/
│   ├── editor-store.ts    # CV content + history
│   └── ui-store.ts        # UI state
├── lib/
│   ├── api/client.ts      # Typed API client (unwraps responses)
│   ├── logger.ts          # Configurable logging
│   └── templates/         # CV templates (starter, blank)
├── components/editor/
│   ├── EditorLayout.tsx   # Main layout (responsive breakpoints)
│   ├── CVPreview.tsx      # Live preview (template dispatch)
│   ├── SectionPalette.tsx # Add sections
│   ├── PropertiesPanel.tsx# Edit section content
│   └── editors/           # Section-specific editors
├── components/templates/
│   ├── ClassicTemplate.tsx# Traditional CV layout
│   ├── ModernTemplate.tsx # Contemporary with timeline
│   ├── VisionaryTemplate.tsx # Two-column with sidebar
│   └── index.ts           # Template exports
└── hooks/
    ├── useAutoSave.ts     # 2s debounced save
    └── useKeyboardShortcuts.ts
```

### Frontend Debugging
```js
// Enable debug logs in browser console:
localStorage.setItem('LOG_LEVEL', 'debug');
location.reload();
```

---

## Backend Architecture

```
HTTP Request → Handler → Service → Repository → PostgreSQL
                 ↓
              Middleware (Auth, Logging, CORS)
```

- **Handlers**: Parse HTTP, validate input, call services
- **Services**: Business logic, orchestration
- **Repositories**: Data access via sqlc-generated code
- **Domain**: Pure types, no external dependencies

### API Response Format
All responses wrapped in standard format:
```json
// Success
{ "data": { ... } }

// Error
{ "error": { "code": "...", "message": "..." } }
```

### Key Backend Files
```
api/
├── cmd/server/main.go     # Entry point, router setup
├── internal/
│   ├── handler/           # HTTP handlers
│   ├── service/           # Business logic
│   ├── repository/        # Data access
│   ├── domain/            # Types + errors
│   └── middleware/        # Auth, CORS
└── db/
    ├── migrations/        # SQL migrations
    ├── queries/           # sqlc queries
    └── seed/              # Dev seed data
```

---

## Project Structure

```
hireme/
├── api/                    # Go backend
├── web/                    # Next.js frontend
├── docker/                 # Docker configs
├── schemas/                # Shared JSON schemas
├── scripts/                # Dev utilities
├── CONTEXT.md              # This file
├── WORKLOG.md              # Trimmed — historical logs in git history
└── CLAUDE.md               # AI assistant instructions
```

---

## API Endpoints

```
GET  /health                    → {"data": {"status": "healthy"}}
GET  /ready                     → {"data": {"status": "ready"}}
GET  /api/v1/users/me           → Current user profile
PATCH /api/v1/users/me          → Update user profile
GET  /api/v1/cv                 → User's active CV
POST /api/v1/cv                 → Create new CV
PUT  /api/v1/cv/{id}            → Update CV
DELETE /api/v1/cv/{id}          → Delete CV
POST /api/v1/assets             → Upload asset (image)
GET  /api/v1/assets/{id}        → Get asset metadata
DELETE /api/v1/assets/{id}      → Delete asset
POST /api/v1/export/{format}    → Export CV (pdf, docx)
```

---

## Common Commands

```bash
# Full Development (recommended)
task dev               # Start infra + API + web together
task dev:stop          # Stop everything + kill dev ports
task dev:restart       # Full stop and restart

# Partial Development
task dev:api           # Start infra + API only
task dev:web           # Start infra + web only

# Infrastructure
task infra:up          # Start PostgreSQL + Gotenberg
task infra:down        # Stop infrastructure

# Database
task db:migrate        # Run migrations
task db:seed           # Seed dev user + sample CV
task db:psql           # Open psql shell

# Testing
npm test               # Frontend tests (in web/)
task api:test          # Backend tests (when implemented)

# Code Generation
task api:sqlc          # Generate sqlc code
```

---

## Task Board

### Active
| ID | Task | Branch | Blocked By | Status |
|----|------|--------|------------|--------|
| — | — | — | — | — |

### Backlog — P1 (Should)
| ID | Task | Blocked By | Size |
|----|------|------------|------|
| T-015 | Multiple CV support | — | M |
| T-035 | Export save location — user-configurable download path with feature flag | — | S |
| T-036 | Fix modal animation jank — editor dialogs snap to bottom-right before centering | — | XS–S |

### Backlog — P2 (Features)
| ID | Task | Blocked By | Size |
|----|------|------------|------|
| T-031 | True responsive editor — collapsible sidebars + mobile layout | — | M |
| T-021 | Export markdown + bundle download (PDF+DOCX+JSON+MD zip) | — | S |
| T-022 | JSON round-trip import | — | S |
| T-027 | Markdown import (parse exported markdown back into CV) | T-021 | S–M |
| T-028 | CV import from uploaded DOCX/PDF — app-based extraction | — | M |
| T-029 | Consultant profile CV type — new section structure + template | — | M+ |
| T-013 | R2 cloud storage implementation | — | M |
| T-014 | OAuth authentication (Google OIDC) | — | M |
| T-037 | Unified export design alignment — PDF + DOCX match editor templates | T-035 | M |

### Done
| ID | Task | PR |
|----|------|----|
| T-001 | Verify template switching persists to backend | feat/t-001-template-persist |
| T-002 | Create HTML generation for CV export | feat/t-002-html-generation |
| T-003 | Implement PDF export via Gotenberg | feat/t-003-pdf-export |
| T-004 | Implement DOCX export via godocx | feat/t-004-docx-export |
| T-005 | Wire frontend export modal to real API | (shipped in MVP) |
| T-006 | P0 smoke tests — section editors | test/t-006-007-smoke-tests |
| T-007 | P0 smoke tests — page components | test/t-006-007-smoke-tests |
| T-008 | P1 interaction tests — editors | test/t-008-010-interaction-ux-tests |
| T-009 | P1 interaction tests — export error paths | test/t-008-010-interaction-ux-tests |
| T-010 | P2 UX regression tests | test/t-008-010-interaction-ux-tests |
| T-011 | Add Certifications section | feat/t-011-012-certifications-projects |
| T-012 | Add Projects section | feat/t-011-012-certifications-projects |
| T-016 | Fix save bug after template switch | fix/t-016-save-after-template-switch |
| T-023 | README overhaul — content + design | docs/t-023-readme-overhaul |
| T-024 | Design direction — define visual identity, palette, typography | (decided: Editorial Craft prototype) |
| T-025 | Implement editorial design overhaul | feat/t-025-design-overhaul (#40) |
| T-030 | Extract design system from Editorial Craft prototype | feat/t-030-design-system |
| T-017 | Responsive layout + clipping fixes | feat/t-017-responsive-layout (#41) |
| T-018 | Section click UX — single-click opens properties | feat/t-018-020-ux-fixes (#42) |
| T-019 | Replace date pickers with MonthYearPicker | feat/t-018-020-ux-fixes (#42) |
| T-020 | Tag input for technologies field | feat/t-018-020-ux-fixes (#42) |
| T-033 | Editor pane design system alignment | feat/t-033-editor-design-alignment (#43) |
| T-034 | Unified header across all pages (including editor) | feat/t-034-unified-header (#44) |
| T-032 | Section delete button in palette sidebar | feat/t-032-section-delete-button (#45) |
| T-026 | DOCX export styling overhaul — match actual CV templates | feat/t-026-docx-styling-overhaul (#46) |

---

## Task Details

### P2 — Features

**T-021: Export markdown + bundle download** — S
- Markdown export: client-side `cvContentToMarkdown()` serializer, same pattern as JSON export
- Bundle: `jszip` to package PDF + DOCX + JSON + MD into one zip download
- All export calls already exist; bundle calls them in parallel via `Promise.all`
- Files: `ExportModal.tsx`, new `markdown-serializer.ts`, new `bundle-export.ts`
- Note: adds `jszip` runtime dependency

**T-022: JSON round-trip import** — S
- File input → read JSON → validate against CV schema → hydrate editor via `setContent()`
- Backend validation already exists (`validator.CVValidator`)
- Files: new import UI (dashboard or editor toolbar), `editor-store.ts`

**T-027: Markdown import** — S–M (blocked by T-021)
- Parse previously exported markdown back into `CVContent` structure
- Must handle user modifications to the exported markdown
- Needs custom markdown-to-domain-model parser — fragile, design carefully
- Files: new `markdown-parser.ts`, import UI

**T-028: CV import from uploaded DOCX/PDF** — M (needs splitting when picked up)
- Upload DOCX/PDF → extract structured CV data → create new CV
- App-based extraction: Go libraries for DOCX parsing, PDF text extraction
- Backend: new `POST /api/v1/cv/import` endpoint + document parser service
- No LLM — deterministic parsing with Go libraries (e.g. `unidoc` for PDF, `docx` for DOCX)
- Mapping from raw extracted sections → `CVContent` JSON

**T-029: Consultant profile CV type** — M+ (needs splitting when picked up)
- Entirely new CV class: skills-first, projects-first, certifications-first structure
- Requires: `cvType` discriminant in `CVContent`, different validation rules, new section palette behavior, new template(s)
- Schema migration strategy needed for existing CVs
- Split into sub-tasks when picked up: schema design, backend validation, frontend palette, template(s)

**T-013: R2 cloud storage implementation** — M
- Implement `R2Storage` methods using AWS SDK
- Tests: storage interface tests with mock S3 client

**T-014: OAuth authentication (Google OIDC)** — M (needs splitting when picked up)
- Full auth flow: login, callback, JWT, middleware
- Tests: middleware tests, auth flow tests

**T-031: True responsive editor — collapsible sidebars + mobile layout** — M
- Current T-017 responsive work only shrinks the editor pane; sidebars stay fixed width
- Need: collapsible sidebars at breakpoints, mobile layout where sidebars become drawers/sheets
- Touch-friendly interactions for drag & drop on mobile
- Files: `EditorLayout.tsx`, `SectionPalette.tsx`, `PropertiesPanel.tsx`

**T-015: Multiple CV support** — M (needs splitting when picked up)
- Dashboard shows list, editor route `/editor/[id]`, create new CV
- Tests: routing tests, dashboard list test

**T-035: Export save location** — S
- Let users choose where exported files are saved (custom download path)
- Build with feature flag for potential paid version gating
- Enables AI agent QA workflows — agents can download exports to known paths for browser-based side-by-side comparison
- Files: `ExportModal.tsx`, possibly new settings/preferences UI

**T-036: Fix modal animation jank** — XS–S
- Editor modals (Experience, Education, etc.) briefly appear in bottom-right corner before jumping to center
- Reproducible across all editors that open an additional component (dialog/modal pattern)
- Likely CSS/animation timing issue in dialog mount or radix dialog positioning
- Files: editor dialog components, possibly shared dialog/modal wrapper

**T-037: Unified export design alignment** — M (blocked by T-035)
- Make PDF and DOCX exports visually match the editor template designs more closely
- Builds on T-026 (DOCX styling) work; extends to PDF export as well
- Blocked by T-035 (export save location) to enable side-by-side QA comparison workflow
- Files: `api/internal/export/docx.go`, `api/internal/export/pdf.go`, HTML templates

### Dependency Graph
```
INDEPENDENT (can start anytime):
  T-015, T-035, T-036 (P1)
  T-031 (P2)
  T-021, T-022, T-028
  T-013, T-014

DEPENDENCY CHAINS:
  T-021 (markdown export) → T-027 (markdown import)
  T-035 (export save location) → T-037 (unified export design alignment)

T-029 (consultant profile): independent but needs splitting when picked up
```

### Recommended Sequencing
**Phase 1 — P1 priorities:**
1. **T-036** (modal animation fix) — quick UX bug fix
2. **T-035** (export save location) — enables QA workflows, unlocks T-037
3. **T-015** (multiple CV support) — core feature, needs splitting

**Phase 2 — P2 UX + features:**
4. **T-037** (unified export design) — PDF + DOCX match editors
5. **T-031** (responsive/mobile) — bigger UX effort
6. **T-021** (markdown + bundle) — new capability, unlocks T-027
7. **T-022** (JSON import) — completes export/import story

**Phase 3 — Remaining features + infra:**
8. **T-027** (markdown import)
9. **T-028** (CV upload/extract)
10. **T-013/T-014** — R2, OAuth
11. **T-029** (consultant profile) — split into sub-tasks first

---

## Task Workflow

Every task follows this lifecycle:

```
Backlog → Active (in-progress) → QA Review → Tech Lead Approval → Done
```

### Steps
1. **Investigate** — `/investigate <task>` — PM traces code, documents findings, prepares handoff in notes file
2. **Plan** — `/plan <task>` — Architect researches, makes decisions, writes implementation plan in notes file
3. **QA Plan** — `/qa-plan <task>` — QA audits coverage gaps, appends test recommendations to notes file
4. **Implement** — `/implement <task>` — Engineer implements code + tests per plan, runs verification gates
5. **Local QA** — `/local-qa` — Full local QA gate (unit tests, type check, build, E2E browser testing)
6. **If QA finds issues** — `/plan <task>` (re-plan mode) → `/implement <task>` (re-implementation mode)
7. **Ship** — `/ship` — Commit, push, open PR
8. **QA Review** — Tech lead + CodeRabbit review the PR
9. **Close** — `/close-task <task>` — PM moves task to Done, cleans up notes

### Task Size Rules
- **XS** (< 1 hour): Config changes, verifications, tiny fixes
- **S** (1-3 hours): Single feature, single editor, test suite
- **M** (3-6 hours): Feature with backend+frontend; must be split if larger

### Branch & PR Convention
- Branch: `feat/t-003-pdf-export`, `fix/t-001-template-persist`, `test/t-006-editor-smoke`
- PR title: `[T-003] Implement PDF export via Gotenberg`
- Target: < 300 lines changed per PR

---

## Stakeholder Feedback

New feedback: say `stakeholder feedback: [description]` → PM agent creates/updates tasks in the Task Board.

---

## Coding Guidelines

### Variable Naming
- Use descriptive variable names — no single-letter or cryptic abbreviations
- Good: `personal`, `styling`, `downloadLink`, `bytesPerUnit`, `lineHeightKey`, `widthVal`
- Bad: `p`, `s`, `a`, `k`, `lh`, `w`
- Loop indices (`i`, `j`) are acceptable when the context is obvious

---

## Key Patterns

### Auth Bypass (Development)
Set `AUTH_BYPASS_ENABLED=true` in `.env` — all requests use `dev-user-001`

### sqlc Workflow
1. Write SQL in `api/db/queries/*.sql`
2. Run `task api:sqlc`
3. Use generated `queries.Queries` in repositories

### Repository Error Mapping
```go
// pgx.ErrNoRows → domain.ErrNotFound
user, err := r.q.GetUserByID(ctx, id)
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domain.ErrNotFound
}
```

---

## Development Setup

```bash
task setup             # Full setup (infra, migrate, seed, deps)
task dev               # Start everything (API :8080, Web :3000)
```

## Testing

```bash
# Frontend
cd web && npm test

# API Health check
curl http://localhost:8080/health

# Get CV (auth bypass must be enabled)
curl http://localhost:8080/api/v1/cv
```
