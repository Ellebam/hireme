# T-003: Implement PDF Export via Gotenberg

## Summary
- **Goal:** Wire `POST /api/v1/export/pdf` to render HTML (via T-002 renderer) and convert it to PDF using Gotenberg's Chromium endpoint
- **Acceptance:** `curl -X POST http://localhost:8080/api/v1/export/pdf` returns a valid PDF binary with correct `Content-Type: application/pdf`
- **Branch:** `feat/t-003-pdf-export`

---

## Investigation Checklist
- [x] Trace existing export handler stub
- [x] Trace HTML renderer from T-002
- [x] Trace Gotenberg Docker/config setup
- [x] Map service & handler dependency patterns
- [x] Map test patterns for handler tests
- [x] Identify all missing components
- [x] Check domain types and repository layer

## Findings

### 1. Export Handler (STUB)
**File:** `api/internal/handler/export.go` (50 lines)
- `CreateExport(w, r)` — validates format via `domain.IsValidExportFormat()`, returns **501 Not Implemented**
- `GetExport(w, r)` — returns **501 Not Implemented**
- Route: `POST /api/v1/export/{format}` registered in `cmd/server/main.go:225`
- Response type `ExportResponse` defined (job-oriented: ID, Format, Status, URL) — **this is wrong for synchronous PDF export** which should stream the binary

### 2. HTML Renderer (COMPLETE from T-002)
**File:** `api/internal/export/renderer.go` (404 lines)
- `NewRenderer() (*Renderer, error)` — parses 3 embedded templates (classic, modern, visionary)
- `Render(content domain.CVContent) (string, error)` — takes parsed CVContent, returns self-contained HTML string
- Templates embedded via `//go:embed templates/*.tmpl`
- Section parsers in `export/section.go` (55 lines) — 6 section types
- **Comprehensive tests** in `renderer_test.go` (533 lines)

### 3. Gotenberg Configuration
**Docker (dev):** `docker/docker-compose.infra.yml:29-44`
- Image: `gotenberg/gotenberg:8`
- Port: `3001:3000` (host:container)
- Flags: `--api-timeout=60s`, `--chromium-disable-javascript=true`, `--chromium-allow-list=file:///tmp/.*`
- Health check: `curl -f http://localhost:3000/health`

**Config:** `api/internal/config/config.go:64-66`
- `ExportConfig.GotenbergURL` loaded from `GOTENBERG_URL` env (default: `http://localhost:3001`)

**Feature flags:** `config.go:77-81`
- `EnableExportPDF` (default: true), `EnableExportDOCX` (default: true)

### 4. Domain Types (COMPLETE)
**File:** `api/internal/domain/asset.go:59-97`
- `ExportJob` struct with full lifecycle fields
- Format constants: `pdf`, `docx`, `json`, `yaml`
- Status constants: `pending`, `processing`, `completed`, `failed`
- `IsValidExportFormat(format) bool` validation

**File:** `api/internal/domain/cv.go`
- `CV.Content` is `json.RawMessage`
- `CV.ParseContent() (*CVContent, error)` — parses into structured type the renderer expects
- `CVContent` has `TemplateID`, `Sections`, `Styling` — exactly what `Renderer.Render()` takes

### 5. Database / Repository Layer
- `export_jobs` table exists in migrations (`000001_init.up.sql:84-102`)
- SQLC queries generated (`postgres/queries/export_jobs.sql.go`, 188 lines)
- `ExportJobRepository` interface defined (`repository/repository.go:80-93`)
- **No implementation file** (`postgres/export.go` does not exist)

### 6. Handler Structure Pattern
**File:** `api/internal/handler/handler.go` (33 lines)
- `Handler` struct holds concrete service types: `*service.UserService`, `*service.CVService`, `*service.AssetService`
- No `ExportService` wired yet
- Constructor via `Dependencies` struct → `New(deps)`

**Test pattern:** `handler/testutil_test.go` (637 lines)
- `TestHandler` uses **interfaces** (not concrete types) for mocking
- Interfaces: `UserServiceInterface`, `CVServiceInterface`, `AssetServiceInterface`
- Mocks use func fields (e.g., `GetByUserIDFunc`)
- Helpers: `newAuthenticatedRequest()`, `newAuthenticatedRequestWithParams()`, `parseJSONResponse()`

### 7. Service Pattern
**File:** `api/internal/service/cv_service.go`
- Constructor takes repository interfaces
- `GetByUserID(ctx, userID) (*domain.CV, error)` — retrieves user's active CV

---

## Gaps Found

### BLOCKER — No Gotenberg HTTP Client
There is no code to call Gotenberg's API. Need to create a client that:
- POSTs multipart form to `{GotenbergURL}/forms/chromium/convert/html`
- Sends HTML as a file part named `files` (with filename `index.html`)
- Receives PDF binary response
- Handles errors (timeout, Gotenberg errors)

### BLOCKER — Export Handler is Job-Oriented, Task Requires Synchronous
The current `ExportResponse` (ID, Format, Status, URL) implies an async job queue pattern. But the task acceptance criteria is a **synchronous** curl → PDF response. Two options:
1. **Synchronous (MVP):** Render HTML → call Gotenberg → stream PDF binary back. No export_jobs table needed.
2. **Async (future):** Create job → process in background → poll for result.

**Recommendation:** Synchronous for T-003. The job infrastructure can be used later.

### BLOCKER — Handler Needs ExportService or Direct Gotenberg Access
The `Handler` struct doesn't have an export service. Need to:
1. Create an `ExportService` that depends on: Renderer, Gotenberg client, CVService (or CV repo)
2. Add `ExportService` to `Handler.Dependencies`
3. Or keep it simpler: handler calls renderer + gotenberg client directly

### Medium — No ExportJobRepository Implementation
`postgres/export.go` doesn't exist. The SQLC code is generated, the interface is defined, but the adapter isn't built. **Not needed for synchronous PDF export** — can be deferred.

### Medium — Feature Flag Not Checked
`config.Features.EnableExportPDF` is defined but the handler doesn't check it. Should gate the endpoint.

### Low — `ExportResponse` Struct May Need Rework
If we go synchronous, the `ExportResponse` struct in `export.go` becomes unused for PDF. It might still be useful for JSON export or future async jobs.

---

## Architect Handoff

### Decisions Needed

1. **Synchronous vs Async Export** (Quick Decision)
   - **Option A (Recommended): Synchronous** — Handler renders HTML → calls Gotenberg → streams PDF bytes back. Simple, matches acceptance criteria. Export job table deferred.
   - **Option B: Async** — Create export job → process → poll. Over-engineered for MVP, frontend doesn't support polling yet.

2. **Gotenberg Client Design** (Quick Decision)
   - **Option A (Recommended): Simple function/struct** in `internal/export/gotenberg.go` — `ConvertHTMLToPDF(ctx, gotenbergURL, html string) ([]byte, error)`. Uses `net/http` + `mime/multipart`. Interface for testability.
   - **Option B: Third-party library** — e.g., `github.com/nickvdyck/gotenberg-go-client`. Adds dependency, may be overkill.

3. **Where to Put Business Logic** (Quick Decision)
   - **Option A (Recommended): ExportService** — New `service/export_service.go` with `ExportPDF(ctx, userID) ([]byte, error)`. Fetches CV, parses content, renders HTML, calls Gotenberg. Clean separation.
   - **Option B: Handler-level** — Handler directly calls renderer + Gotenberg client. Simpler but less testable.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `api/internal/export/gotenberg.go` | **NEW** — Gotenberg HTTP client | BLOCKER: no client |
| `api/internal/export/gotenberg_test.go` | **NEW** — Client tests (mock HTTP) | BLOCKER: no tests |
| `api/internal/service/export_service.go` | **NEW** — ExportService with ExportPDF method | BLOCKER: no service |
| `api/internal/service/export_service_test.go` | **NEW** — Service tests (mock renderer + client) | BLOCKER: no tests |
| `api/internal/handler/export.go` | **MODIFY** — Replace 501 stub with synchronous PDF streaming | BLOCKER: stub |
| `api/internal/handler/handler.go` | **MODIFY** — Add ExportService to Dependencies | BLOCKER: not wired |
| `api/internal/handler/export_test.go` | **NEW** — Handler tests | BLOCKER: no tests |
| `api/internal/handler/testutil_test.go` | **MODIFY** — Add ExportServiceInterface + MockExportService | BLOCKER: mock needed |
| `api/cmd/server/main.go` | **MODIFY** — Wire ExportService in dependency setup | BLOCKER: not wired |

### Tests to Add

1. **Gotenberg client tests** (`export/gotenberg_test.go`)
   - Mock HTTP server returns PDF bytes → verify client passes HTML correctly
   - Mock HTTP server returns error → verify error propagation
   - Timeout test

2. **ExportService tests** (`service/export_service_test.go`)
   - Happy path: mock CV repo + renderer + Gotenberg client → returns PDF bytes
   - CV not found → returns domain.ErrNotFound
   - Renderer error → returns error
   - Gotenberg error → returns error

3. **Handler tests** (`handler/export_test.go`)
   - `POST /api/v1/export/pdf` → 200 with `Content-Type: application/pdf`
   - Invalid format → 400
   - Feature flag disabled → 501 or 403
   - Service error → appropriate HTTP error

### Constraints
- PR target: < 300 lines changed
- No new Go dependencies needed (stdlib `net/http` + `mime/multipart` suffice)
- Gotenberg must be running for integration tests (can be skipped in CI if needed)
- Handler tests should use mocked ExportService (no real Gotenberg)

### Recommended Next Agent
**@architect** — to finalize the 3 design decisions above, then **@engineer** for implementation.

---

## Test Plan

### Unit Tests (must pass without Gotenberg)
- Gotenberg client: httptest mock server
- ExportService: mock renderer + mock Gotenberg client + mock CV repo
- Handler: mock ExportService

### Integration Test (requires Gotenberg running)
- Start infra → call `POST /api/v1/export/pdf` → verify response is valid PDF (check `%PDF-` magic bytes)
- Can be tagged `//go:build integration` to skip in CI
