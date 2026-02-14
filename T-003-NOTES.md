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

---

## Architect Plan

### Decision 1: Synchronous vs Async Export

- **Option A: Synchronous** — Handler renders HTML → calls Gotenberg → streams PDF bytes back in same request. No export_jobs table, no polling. Pros: Simple, matches acceptance criteria exactly, no async infrastructure needed. Cons: Request blocks until Gotenberg finishes (~1-3s).
- **Option B: Async** — Create export job → background worker → client polls for completion. Pros: Non-blocking, scales better for large PDFs. Cons: Over-engineered for MVP, frontend has no polling logic, requires job queue infrastructure.
- **Choice:** A (Synchronous) — **Rationale:** The acceptance criteria is a single `curl` returning a PDF. Gotenberg conversion takes 1-3s which is acceptable for a synchronous request. The async job table (`export_jobs`) and its SQLC queries are already scaffolded and can be used later if needed. Frontend export modal in T-005 expects a direct download, not polling.

### Decision 2: Gotenberg Client Design

- **Option A: Stdlib client** — New `GotenbergClient` struct in `internal/export/gotenberg.go` using `net/http` + `mime/multipart`. Define a `PDFConverter` interface for testability. Pros: Zero dependencies, full control, small footprint (~40 lines). Cons: Must construct multipart form manually.
- **Option B: Third-party library** (`github.com/dcaraxes/gotenberg-go-client/v8`) — Pre-built client. Pros: Handles multipart construction. Cons: New dependency for ~15 lines of saved code, may lag behind Gotenberg releases.
- **Choice:** A (Stdlib) — **Rationale:** The Gotenberg API is a single multipart POST. The Go stdlib handles this natively. Adding a dependency for something this simple violates YAGNI and the constraint of no new Go dependencies.

### Decision 3: Where to Put Business Logic

- **Option A: ExportService** — New `service/export_service.go` with `ExportPDF(ctx, userID) ([]byte, error)`. Orchestrates: fetch CV → parse content → render HTML → convert to PDF. Depends on interfaces: `repository.CVRepository`, `HTMLRenderer`, `PDFConverter`. Pros: Clean separation, independently testable, follows existing service pattern. Cons: One more file.
- **Option B: Handler-level** — Handler directly calls renderer + Gotenberg client. Pros: Fewer files. Cons: Handler becomes fat, harder to test (must mock 3 things instead of 1), breaks the Handler → Service → Repository pattern.
- **Choice:** A (ExportService) — **Rationale:** Follows the established `Handler → Service → Repository` architecture. Handler tests mock one interface (`ExportServiceInterface`). Service tests mock three granular interfaces. This is the same layering used by `CVService`, `UserService`, `AssetService`.

### Files Changing

| # | File | Action | Lines (est) |
|---|------|--------|-------------|
| 1 | `api/internal/export/gotenberg.go` | **NEW** — `PDFConverter` interface + `GotenbergClient` struct | ~45 |
| 2 | `api/internal/export/gotenberg_test.go` | **NEW** — httptest mock: success, error, timeout | ~55 |
| 3 | `api/internal/service/export_service.go` | **NEW** — `ExportService` + `HTMLRenderer`/`PDFConverter` interfaces | ~50 |
| 4 | `api/internal/service/export_service_test.go` | **NEW** — 4 test cases (happy, CV not found, render error, convert error) | ~60 |
| 5 | `api/internal/handler/export.go` | **MODIFY** — Replace 501 stub with sync PDF streaming, feature flag check | ~30 net |
| 6 | `api/internal/handler/handler.go` | **MODIFY** — Add `exportService` field + `ExportService` in `Dependencies` | ~4 |
| 7 | `api/internal/handler/export_test.go` | **NEW** — 4 test cases (success, invalid format, feature disabled, service error) | ~50 |
| 8 | `api/internal/handler/testutil_test.go` | **MODIFY** — Add `ExportServiceInterface` + `MockExportService` + TestHandler field | ~25 |
| 9 | `api/cmd/server/main.go` | **MODIFY** — Initialize renderer, gotenberg client, export service; wire into handler | ~10 |

**Estimated total:** ~290 lines changed (within 300-line PR target)

### Files NOT Changing

| File | Reason |
|------|--------|
| `domain/asset.go` | `ExportJob`, format constants, `IsValidExportFormat` — all reused as-is |
| `domain/cv.go` | `CV.ParseContent()`, `CVContent` — already complete |
| `export/renderer.go` | HTML renderer from T-002 — used as-is, no modifications needed |
| `export/section.go` | Section parsers — used as-is |
| `repository/repository.go` | `CVRepository` interface already has `GetByUserID` — sufficient |
| `postgres/export.go` | Not needed — synchronous export doesn't use export_jobs table |
| `config/config.go` | `ExportConfig.GotenbergURL` and `Features.EnableExportPDF` already defined |

### Detailed Design

#### 1. Gotenberg Client (`export/gotenberg.go`)

```go
// PDFConverter converts HTML to PDF.
type PDFConverter interface {
    ConvertHTMLToPDF(ctx context.Context, html string) ([]byte, error)
}

// GotenbergClient calls Gotenberg's Chromium HTML-to-PDF endpoint.
type GotenbergClient struct {
    url        string        // e.g. "http://localhost:3001"
    httpClient *http.Client
}

func NewGotenbergClient(url string) *GotenbergClient
func (g *GotenbergClient) ConvertHTMLToPDF(ctx context.Context, html string) ([]byte, error)
```

Implementation notes:
- POST multipart form to `{url}/forms/chromium/convert/html`
- File field: `files` with filename `index.html` and content-type `text/html`
- Check response status: 200 → return body, non-200 → read error body and wrap in `fmt.Errorf`
- Use 30s timeout via `http.Client.Timeout` (Gotenberg itself has 60s timeout)
- Return `([]byte, error)` — caller writes to response

#### 2. Export Service (`service/export_service.go`)

```go
// HTMLRenderer renders CV content to HTML string.
type HTMLRenderer interface {
    Render(content domain.CVContent) (string, error)
}

// ExportService handles document export operations.
type ExportService struct {
    cvRepo   repository.CVRepository
    renderer HTMLRenderer
    pdf      export.PDFConverter
}

func NewExportService(cvRepo repository.CVRepository, renderer HTMLRenderer, pdf export.PDFConverter) *ExportService
func (s *ExportService) ExportPDF(ctx context.Context, userID string) ([]byte, error)
```

`ExportPDF` flow:
1. `s.cvRepo.GetByUserID(ctx, userID)` → `*domain.CV` (propagates `ErrNotFound`)
2. `cv.ParseContent()` → `*domain.CVContent`
3. `s.renderer.Render(*content)` → HTML string
4. `s.pdf.ConvertHTMLToPDF(ctx, html)` → PDF bytes
5. Return `(pdfBytes, nil)`

#### 3. Handler Changes (`handler/export.go`)

`CreateExport` becomes:
1. Get format from `chi.URLParam(r, "format")`
2. Validate format with `domain.IsValidExportFormat(format)`
3. Switch on format:
   - `"pdf"`: check `h.config.Features.EnableExportPDF` → if disabled, return 501. Get userID → call `h.exportService.ExportPDF(ctx, userID)` → set `Content-Type: application/pdf`, `Content-Disposition: attachment; filename="export.pdf"` → write bytes.
   - Other formats: return 501 Not Implemented
4. On error: use `httputil.HandleError(w, err)` for domain errors, `httputil.Error(w, 500, ...)` for others.

`GetExport` stays 501 (deferred to async job feature).

`ExportResponse` struct stays (unused but harmless; may be used by future JSON/YAML export).

#### 4. Handler Wiring (`handler/handler.go`)

Add to `Handler`:
```go
exportService *service.ExportService
```

Add to `Dependencies`:
```go
ExportService *service.ExportService
```

#### 5. Test Handler Additions (`handler/testutil_test.go`)

```go
type ExportServiceInterface interface {
    ExportPDF(ctx context.Context, userID string) ([]byte, error)
}

type MockExportService struct {
    ExportPDFFunc func(ctx context.Context, userID string) ([]byte, error)
}
```

Add `exportService ExportServiceInterface` to `TestHandler` struct.
Add `CreateExport` method to `TestHandler`.
Update `NewTestHandler` to accept export service (4th param).

#### 6. Main Wiring (`cmd/server/main.go`)

After existing service initialization:
```go
renderer, err := export.NewRenderer()
gotenbergClient := export.NewGotenbergClient(cfg.Export.GotenbergURL)
exportSvc := service.NewExportService(cvRepo, renderer, gotenbergClient)
```

Add `ExportService: exportSvc` to `handler.Dependencies`.

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Gotenberg not running → 500 on export | Medium | Low | Error message includes "gotenberg" context; dev env starts it via `task infra:up` |
| Large CV HTML → Gotenberg timeout | Low | Low | 30s client timeout, 60s Gotenberg timeout; CVs are small documents |
| Memory spike from large PDF in memory | Low | Low | Typical CV PDF is <500KB; `[]byte` is fine for MVP |
| Breaking TestHandler constructor signature | Medium | Medium | All existing tests pass `nil` for 4th param; update call sites |

### Consequences

**Becomes easier:**
- T-004 (DOCX export) — add `ExportDOCX` to ExportService, `ConvertHTMLToDOCX` to a new converter, extend handler switch
- T-005 (frontend wiring) — backend endpoint is ready, frontend just calls `POST /api/v1/export/pdf`
- Future async export — ExportService can be wrapped by a job worker without changing the core logic

**Becomes harder:**
- Nothing. The synchronous approach is strictly simpler than async.

### Engineer Execution Steps

1. **Create `api/internal/export/gotenberg.go`** — `PDFConverter` interface + `GotenbergClient` implementation
2. **Create `api/internal/export/gotenberg_test.go`** — 3 tests using `httptest.NewServer`
3. **Create `api/internal/service/export_service.go`** — `HTMLRenderer` interface (consumed here), `ExportService` struct, `ExportPDF` method
4. **Create `api/internal/service/export_service_test.go`** — 4 test cases with mocks
5. **Modify `api/internal/handler/handler.go`** — Add `exportService` to Handler + Dependencies
6. **Modify `api/internal/handler/testutil_test.go`** — Add `ExportServiceInterface`, `MockExportService`, update `TestHandler` + `NewTestHandler`
7. **Modify `api/internal/handler/export.go`** — Replace 501 stub with synchronous PDF export logic + feature flag check
8. **Create `api/internal/handler/export_test.go`** — 4 handler test cases
9. **Modify `api/cmd/server/main.go`** — Wire renderer, gotenberg client, export service into handler dependencies
10. **Run `go test ./...`** from `api/` — all tests must pass
11. **Run `go vet ./...`** — no issues
12. **Run `go build ./cmd/server`** — compiles cleanly

### Verification Gates

- [ ] `cd api && go test ./...` — all pass (unit tests, no Gotenberg needed)
- [ ] `cd api && go vet ./...` — no issues
- [ ] `cd api && go build ./cmd/server` — compiles
- [ ] Manual integration test: `task infra:up && task api:dev`, then `curl -s -o /tmp/test.pdf -w "%{http_code}" -X POST http://localhost:8080/api/v1/export/pdf` returns 200 and `file /tmp/test.pdf` shows "PDF document"

### QA-Recommended Additional Tests

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 1 | Multipart form construction | Backend Gotenberg client | High | In `export/gotenberg_test.go`, add a test that inspects the multipart request body sent to the mock `httptest.Server`. Assert: (a) field name is `files`, (b) filename is `index.html`, (c) content-type is `text/html`, (d) body matches the input HTML string. The existing plan tests response handling but not the *request* construction — a malformed multipart would silently produce a Gotenberg 400. |
| 2 | Context cancellation propagation | Backend Gotenberg client | High | In `export/gotenberg_test.go`, add a test where the context is cancelled before the mock server responds. Assert that `ConvertHTMLToPDF` returns a `context.Canceled` (or `context.DeadlineExceeded`) error, not a generic HTTP error. Ensures the `http.NewRequestWithContext` path works. |
| 3 | Content-Type and Content-Disposition headers | Backend handler | High | In `handler/export_test.go`, for the success case, explicitly assert: `Content-Type` is `application/pdf`, `Content-Disposition` is `attachment; filename="export.pdf"`. The plan mentions checking `Content-Type` but not `Content-Disposition` — both are in the architect's design and must be tested. |
| 4 | PDF body bytes match service output | Backend handler | High | In `handler/export_test.go` success case, assert that `rr.Body.Bytes()` equals the exact byte slice returned by the mock `ExportPDFFunc`. This confirms the handler streams bytes directly without wrapping them in JSON (unlike all other handlers which use `httputil.JSON`). |
| 5 | Valid non-PDF format returns 501 | Backend handler | Medium | In `handler/export_test.go`, add a test with format `docx` (a valid format per `IsValidExportFormat`, but not yet implemented). Assert 501 Not Implemented. This is distinct from the "invalid format → 400" test — it covers the handler's format switch/fallthrough for unimplemented formats. |
| 6 | ParseContent error propagation | Backend ExportService | Medium | In `service/export_service_test.go`, add a test where `cv.Content` is invalid JSON (e.g. `json.RawMessage("not json")`). Assert that `ExportPDF` returns an error wrapping the JSON parse failure. This tests the `cv.ParseContent()` step between the CV repo call and the renderer call — the current plan mocks the CV repo response but doesn't test a corrupt `Content` field. |
| 7 | Empty HTML from renderer | Backend ExportService | Medium | In `service/export_service_test.go`, add a test where the mock renderer returns `("", nil)` — technically success but empty HTML. Assert the service either: (a) propagates the empty string to the PDF converter (which should handle it), or (b) returns a specific error. This documents the intended behavior for edge cases. |
| 8 | Table-driven format validation in handler | Backend handler | Medium | In `handler/export_test.go`, consolidate the invalid-format test into a table-driven subtest covering: `""` (empty), `"html"`, `"exe"`, `"PDF"` (wrong case). Uses `t.Run` with `newAuthenticatedRequestWithParams` varying the `format` param. Ensures `IsValidExportFormat` is called correctly for all edge cases at the handler boundary. |
| 9 | Gotenberg non-200 error message | Backend Gotenberg client | Medium | In `export/gotenberg_test.go`, add a test where the mock server returns HTTP 400 with a body like `"invalid file format"`. Assert that the returned error message includes both the status code and the response body text — otherwise debugging Gotenberg issues will be opaque. |
| 10 | TestHandler constructor backward compat | Backend handler test infra | Medium | In `handler/testutil_test.go`, after adding the `exportService` field, verify that existing test files (`cv_test.go`, `user_test.go`, `asset_test.go`, `health_test.go`) all pass with `nil` for the new export service parameter. Add a comment or verify in a dedicated test that `NewTestHandler(userSvc, cvSvc, assetSvc, nil)` does not panic. |
