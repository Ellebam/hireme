# T-015: Multiple CV Support

## Summary
- **Goal:** Allow users to create, list, and manage multiple CVs — dashboard shows all CVs, editor loads a specific CV by ID, export targets a specific CV
- **Acceptance:** User can create multiple CVs, see them all on the dashboard, edit each independently via `/editor/[id]`, export any specific CV
- **Branch:** `feat/t-015-multiple-cv-support`

---

## Investigation Checklist
- [x] Database schema — unique constraints, is_active, cv_limit
- [x] SQL queries — GetCVByUserID (LIMIT 1), ListCVsByUserID, CountByUserID
- [x] Repository layer — CVRepository interface methods
- [x] Service layer — CVService (GetByUserID, Create with quota check)
- [x] Handler layer — GetCV, CreateCV, UpdateCV, DeleteCV
- [x] Export service + handler — fetchCVContent, ExportPDF/DOCX
- [x] Router/routes — /cv, /cv/{id}, /export/{format}
- [x] Seed data — single seed CV
- [x] API client (frontend) — cvApi methods
- [x] Editor store — single CV in state
- [x] Dashboard page — shows single CV
- [x] Editor page/route — static /editor, no [id] segment
- [x] Templates page — links to /editor?template=X
- [x] Auto-save hook — uses cv.id (already compatible)
- [x] Export modal — uses cv.id (already compatible)
- [x] CV types/interfaces

## Findings

### What's Already Multi-CV Ready

| Component | File | Status |
|-----------|------|--------|
| DB schema | `api/db/migrations/000001_init.up.sql:40-56` | No UNIQUE(user_id) constraint — multiple CVs allowed |
| `ListCVsByUserID` query | `api/db/queries/cvs.sql:11-14` | Returns all CVs for a user |
| `CountByUserID` query | `api/db/queries/cvs.sql:16-17` | For quota enforcement |
| `User.CanCreateCV()` | `api/internal/domain/user.go:42-45` | Checks against `CVLimit` |
| CV quota field | `api/db/migrations/000001_init.up.sql` (users table) | `cv_limit INTEGER DEFAULT 1`, per-tier |
| `ListByUserID` repo | `api/internal/repository/postgres/cv.go:52-63` | Returns `[]*domain.CV` |
| `CountByUserID` repo | `api/internal/repository/postgres/cv.go:65-71` | Used for quota check |
| Update/Delete handlers | `api/internal/handler/cv.go` | Already take CV ID from URL |
| Auto-save hook | `web/src/hooks/useAutoSave.ts:16,28` | Uses `currentCV.id` explicitly |
| Export modal | `web/src/components/editor/modals/ExportModal.tsx:49,81` | Uses `cv.id` explicitly |
| Delete action (dashboard) | `web/src/app/page.tsx:55` | Uses `cv.id` for delete |
| CV TypeScript types | `web/src/types/api.ts` | CV interface already has `id: string` |

### What Enforces Single-CV (Must Change)

#### Backend

| What | File:Line | Problem |
|------|-----------|---------|
| `GetCVByUserID` query | `api/db/queries/cvs.sql:4-9` | `LIMIT 1` — returns only most recent active CV |
| `GetByUserID` repo method | `api/internal/repository/postgres/cv.go:41-50` | Calls single-CV query |
| `CVService.GetByUserID()` | `api/internal/service/cv_service.go` | Returns single CV, used by handler + export |
| `Handler.GetCV()` | `api/internal/handler/cv.go` | `GET /api/v1/cv` → returns one CV |
| `ExportService.fetchCVContent()` | `api/internal/service/export_service.go` | Uses `GetByUserID` — no CV ID param |
| `ExportService.ExportPDF/DOCX()` | `api/internal/service/export_service.go` | Takes `userID` only, not `cvID` |
| `Handler.CreateExport()` | `api/internal/handler/export.go:24-74` | No CV ID in request — exports "the" active CV |
| Export route | `api/cmd/server/main.go:241` | `POST /api/v1/export/{format}` — no CV ID |
| CV routes | `api/cmd/server/main.go:229-233` | `GET /cv` returns single CV, no list endpoint |

#### Frontend

| What | File:Line | Problem |
|------|-----------|---------|
| API client `cvApi.get()` | `web/src/lib/api/client.ts:249` | `GET /api/v1/cv` — no list method, no get-by-id |
| Editor store | `web/src/stores/editor-store.ts:32` | `cv: CV \| null` — single CV in state |
| Dashboard page | `web/src/app/page.tsx:35` | Calls `api.cv.get()` — shows one CV |
| Dashboard hardcoded "1" | `web/src/app/page.tsx:126` | Document count is hardcoded |
| Dashboard edit link | `web/src/app/page.tsx:132,159` | Links to `/editor` — no CV ID |
| Editor route | `web/src/app/editor/page.tsx` | Static `/editor` — no `[id]` dynamic segment |
| Editor CV loading | `web/src/app/editor/page.tsx:35` | `api.cv.get()` — loads "the" CV, auto-creates if missing |
| Templates page links | `web/src/app/templates/page.tsx:118` | Links to `/editor?template=X` — no CV context |
| Create new document link | `web/src/app/page.tsx:179` | Goes to `/editor` — creates/loads single CV |

## Gaps Found

### BLOCKER — Must fix for multi-CV to work

1. **No `GET /api/v1/cvs` list endpoint** — Backend has the query and repo method (`ListByUserID`) but no handler/route
2. **No `GET /api/v1/cvs/{id}` endpoint** — Need to fetch a specific CV by ID (with ownership check)
3. **Export endpoint lacks CV ID** — `POST /api/v1/export/{format}` needs to know which CV to export
4. **Editor route is static** — Must change `/editor` to `/editor/[id]` for direct-linking
5. **Dashboard shows single CV** — Must fetch and display all CVs
6. **API client missing `list()` and `getById()`** — Frontend can't fetch multiple CVs

### Medium — Should fix, important for UX

7. **"Create New" flow undefined** — Where does the user go to create a new CV? Dashboard button? Templates page? Both?
8. **Template selection with multi-CV** — Does template page create a new CV or modify an existing one?
9. **CV title/naming** — Users need a way to title/rename CVs for identification
10. **No "active CV" indicator** — If multiple CVs exist, which is the "primary"?
11. **Seed data only creates one CV** — Should seed 2-3 CVs for development/testing

### Low — Nice to have

12. **CV duplication** — Users might want to clone an existing CV
13. **CV archive/restore** — `is_active` exists but no UI for it
14. **Dashboard sorting/filtering** — If user has many CVs, need search/filter

## Architect Handoff

### Decisions Needed

**Decision 1: API Route Design**
- **Option A (RESTful rename):** `GET /api/v1/cvs` (list), `GET /api/v1/cvs/{id}` (get one), `POST /api/v1/cvs` (create), `PUT /api/v1/cvs/{id}` (update), `DELETE /api/v1/cvs/{id}` (delete)
  - Pro: Clean REST conventions, plural nouns
  - Con: Breaking change — frontend must update all calls, existing `/cv` routes deprecated
- **Option B (Additive):** Keep `GET /api/v1/cv` (backward compat, returns active), ADD `GET /api/v1/cv/list`, `GET /api/v1/cv/{id}`
  - Pro: Non-breaking, gradual migration
  - Con: Inconsistent naming (`/cv` singular + `/cv/list`)
- **Recommendation:** Option A — clean break is fine since there's no production traffic

**Decision 2: Export Route**
- **Option A:** `POST /api/v1/export/{cvId}/{format}` — CV ID in path
- **Option B:** `POST /api/v1/export/{format}?cvId={id}` — CV ID as query param
- **Option C:** `POST /api/v1/cvs/{cvId}/export/{format}` — nested under CV resource
- **Recommendation:** Option C — most RESTful, CV is the parent resource

**Decision 3: "Create New CV" Flow**
- **Option A:** Dashboard "New Document" → creates blank CV immediately → redirects to `/editor/[newId]`
- **Option B:** Dashboard "New Document" → goes to `/templates` → user picks template → creates CV → `/editor/[newId]`
- **Option C:** Dashboard "New Document" → modal with template picker → create → redirect
- **Recommendation:** Option B — reuses existing templates page, already has template cards

**Decision 4: Editor Route Migration**
- Old: `/editor` (static)
- New: `/editor/[id]` (dynamic)
- **Need redirect:** `/editor` → dashboard (or latest CV?)
- **Recommendation:** `/editor` without ID redirects to dashboard; only `/editor/[id]` loads the editor

**Decision 5: Task Splitting**
This is sized M but should be split. Recommended sub-tasks:
- **T-015a:** Backend — add list/get-by-id endpoints, update export route
- **T-015b:** Frontend — dynamic editor route, dashboard multi-CV display
- **T-015c:** "Create new CV" flow (dashboard → templates → editor)
- **T-015d:** Tests + seed data update

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `api/db/queries/cvs.sql` | Add `GetCVByID` query (with user ownership check) | #2 |
| `api/internal/repository/postgres/cv.go` | Add `GetByID` with ownership check | #2 |
| `api/internal/service/cv_service.go` | Add `GetByID`, `ListByUserID` service methods | #1, #2 |
| `api/internal/handler/cv.go` | Add `ListCVs`, `GetCVByID` handlers; rename routes | #1, #2 |
| `api/internal/handler/export.go` | Accept CV ID from path/query | #3 |
| `api/internal/service/export_service.go` | Add `cvID` param to `ExportPDF`, `ExportDOCX`, `fetchCVContent` | #3 |
| `api/cmd/server/main.go` | Update routes: `/cvs`, `/cvs/{id}`, `/cvs/{id}/export/{format}` | #1, #2, #3 |
| `web/src/lib/api/client.ts` | Add `cvApi.list()`, `cvApi.getById(id)`, update export call | #6 |
| `web/src/app/editor/[id]/page.tsx` | NEW — dynamic editor route, loads CV by ID from params | #4 |
| `web/src/app/editor/page.tsx` | Convert to redirect to `/` | #4 |
| `web/src/app/page.tsx` | Fetch all CVs, display list, update links to `/editor/[id]` | #5 |
| `web/src/app/templates/page.tsx` | Create new CV on template select, redirect to `/editor/[id]` | #8 |
| `web/src/stores/editor-store.ts` | Minor — ensure `setCV` works with any CV (already does) | — |
| `api/db/seed/dev_seed.sql` | Add 2-3 sample CVs with different templates | #11 |

### Constraints
- PR size < 300 lines per sub-task — split into 3-4 PRs
- Backward compat not required (no production users)
- Quota enforcement already works — no schema migration needed
- `is_active` field exists but deferred for later (no archive UI in this task)

### Recommended Next Agent
**@architect** — Multiple design decisions needed (API routes, create flow, route migration, task splitting)

## Test Plan

### Backend Tests
- List CVs returns all user CVs (empty, 1, multiple)
- Get CV by ID returns correct CV (own CV, other user's CV → 403, not found → 404)
- Create CV respects quota (at limit → error, under limit → success)
- Export with CV ID works (own CV, other user's → 403, not found → 404)
- Seed data includes multiple CVs

### Frontend Tests
- Dashboard renders multiple CV cards
- Dashboard "New Document" navigates correctly
- Editor loads specific CV by ID from route params
- Editor handles invalid/missing CV ID (redirect to dashboard)
- Template selection creates new CV and redirects to `/editor/[id]`
- Auto-save works correctly when switching between CVs
- Export targets the currently loaded CV

### QA-Recommended Additional Tests

**Existing coverage reviewed:** 12 CV handler tests (`cv_test.go`), 9 export handler tests (`export_test.go`), 11 CV service tests (`cv_service_test.go`), 11 export service tests (`export_service_test.go`), 5 API client CV tests (`client.test.ts`), 46 editor store tests (`editor-store.test.ts`), 11 auto-save tests (`useAutoSave.test.ts`), 4 dashboard tests (`DashboardPage.test.tsx`), 4 editor page tests (`EditorPage.test.tsx`), 13 template factory tests (`templates.test.ts`).

#### Backend

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 1 | GetCVByID — invalid UUID → 400 | Handler | High | Add `TestGetCVByID_InvalidID` in `cv_test.go` alongside the existing `TestUpdateCV_InvalidID` / `TestDeleteCV_InvalidID` pattern. Pass a non-UUID string as `{id}`, assert 400 with "invalid CV ID" message. Without this, malformed URLs could panic on `uuid.Parse`. |
| 2 | ListCVs — response is JSON array | Handler | High | Add `TestListCVs_Success_Multiple` in `cv_test.go`. Mock `ListByUserIDFunc` returning 3 CVs. Assert response is a JSON array (not `{ data: [...] }` or `null`), length matches, and each element has correct fields. Confirms the new response shape clients depend on. |
| 3 | ListCVs — empty list returns `[]` | Handler | High | Add `TestListCVs_Success_Empty` in `cv_test.go`. Mock `ListByUserIDFunc` returning empty slice. Assert response is `[]` not `null`. Frontend `CV[]` destructuring breaks on `null`. |
| 4 | ListCVs — service error propagation | Handler | Medium | Add `TestListCVs_ServiceError` in `cv_test.go`. Mock `ListByUserIDFunc` returning `domain.ErrInternal`. Assert 500 response. Matches the error-propagation pattern in existing GetCV/Create tests. |
| 5 | Export handler — invalid CV UUID → 400 | Handler | High | Add `TestCreateExport_InvalidCVID` in `export_test.go`. With the new route `/cv/{id}/export/{format}`, pass a non-UUID as `{id}`. Assert 400. Currently only format validation is tested; the new `{id}` param needs the same treatment. |
| 6 | Export service — ownership check | Service | High | Add `TestExportPDF_WrongUser` and `TestExportDOCX_WrongUser` in `export_service_test.go`. After signature changes to `(ctx, cvID, userID)`, mock `GetByIDFunc` returning a CV with different `UserID`. Assert `ErrForbidden`. This is the most critical new security boundary. |
| 7 | CVService.ListByUserID — delegates to repo | Service | Medium | Add `TestCVService_ListByUserID_Success` and `TestCVService_ListByUserID_Empty` in `cv_service_test.go`. Verify the service method correctly delegates to `cvRepo.ListByUserID` and returns the result. Simple passthrough but confirms wiring. |
| 8 | ListCVs — ordering (updated_at DESC) | Handler | Medium | In `TestListCVs_Success_Multiple`, additionally assert that CVs are returned in `updated_at` descending order. The SQL query `ListCVsByUserID` in `cvs.sql:11-14` should ORDER BY `updated_at DESC` — verify the handler preserves that order. |
| 9 | Export handler — missing CV ID (old route removed) | Handler | Low | Verify that the old `POST /export/{format}` route no longer matches (returns 404 or 405). Prevents clients from accidentally using the deprecated endpoint and getting silent failures. |

#### Frontend

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 10 | `cvApi.list()` returns array | API Client | High | Update or add test in `client.test.ts` → `cvApi` describe block. Mock fetch returning `{ data: [cv1, cv2] }`. Call `cvApi.list()`, assert result is `CV[]` with length 2. The old `cvApi.get()` test returns a single object — the new `list()` must unwrap an array. |
| 11 | `cvApi.get(id)` sends correct URL | API Client | High | Add test in `client.test.ts` → `cvApi` describe block. Call `cvApi.get('some-uuid')`, assert fetch was called with `/api/v1/cv/some-uuid`. Verifies the ID is interpolated into the path, not sent as a query param. |
| 12 | `exportApi.export` URL includes CV ID | API Client | High | Update existing export tests in `client.test.ts` → `exportApi` describe block. After the URL change to `/api/v1/cv/{cvId}/export/{format}`, verify fetch is called with the new URL pattern. All 6 existing export tests need URL assertions updated. |
| 13 | Dashboard — empty state links to `/templates` | Page | Medium | Update `DashboardPage.test.tsx`. The empty state currently shows when API returns 404. After change, empty state should show when `list()` returns `[]`. Verify "Create New Document" link points to `/templates` (not `/editor`). |
| 14 | Dashboard — document count is dynamic | Page | Medium | Add test in `DashboardPage.test.tsx`. Mock `cvApi.list()` returning 3 CVs. Assert the document count displays "3" (not hardcoded "1"). Catches the regression from `page.tsx:126`. |
| 15 | Dashboard — delete removes CV from list | Page | Medium | Add test in `DashboardPage.test.tsx`. Render with 2 CVs, delete one, assert only 1 CV card remains. Existing delete test only covers single-CV scenario. |
| 16 | Dashboard — edit links include CV ID | Page | High | Add test in `DashboardPage.test.tsx`. Mock `cvApi.list()` returning a CV with known ID. Assert edit link `href` is `/editor/{cv.id}` (not `/editor`). This is the core navigation change. |
| 17 | Editor page — `useParams` provides CV ID | Page | High | Create `web/src/app/__tests__/EditorIdPage.test.tsx` for the new `editor/[id]/page.tsx`. Mock `useParams` returning `{ id: 'some-uuid' }`. Mock `cvApi.get('some-uuid')`. Assert `EditorLayout` renders with correct CV data. |
| 18 | Editor page — 404 redirects to dashboard | Page | High | In the new `EditorIdPage.test.tsx`. Mock `cvApi.get(id)` throwing 404. Assert `router.replace('/')` is called. Without this, users see a broken editor on invalid URLs. |
| 19 | Editor page — no auto-create on missing CV | Page | Medium | In `EditorIdPage.test.tsx`. Verify that unlike the old editor page, the new dynamic route does NOT call `cvApi.create()` when the CV is not found. Auto-create is removed — creation happens only via templates page. Regression test for the old behavior. |
| 20 | Editor redirect — `/editor` goes to `/` | Page | Medium | Add test for the new `editor/page.tsx` redirect component. Render it, assert `router.replace('/')` is called. Ensures bookmarks to old `/editor` URL don't break. |
| 21 | Templates page — creates CV and redirects | Page | High | Create `web/src/app/__tests__/TemplatesPage.test.tsx`. Mock `cvApi.create()` returning a CV with known ID. Click a template card, assert `cvApi.create` was called with correct template content, then assert `router.push('/editor/{newId}')`. This is the new "Create CV" entry point. |
| 22 | Templates page — quota error display | Page | Medium | In `TemplatesPage.test.tsx`. Mock `cvApi.create()` throwing `ErrCVLimitReached` (403). Click a template, assert error message is shown to user (not a silent failure or unhandled rejection). |

---

## Architect Plan

### Decisions

**Decision 1: API Route Design**
- Option A (RESTful plural `/cvs`): Rename all routes to `/cvs`, `/cvs/{id}` — clean REST but changes POST/PUT/DELETE too (larger diff, more test churn)
- Option B (Keep singular `/cv`, repurpose GET): `GET /cv` returns list (was single), add `GET /cv/{id}` — POST/PUT/DELETE routes unchanged
- **Choice:** Option B — **Rationale:** Minimizes diff. PUT `/cv/{id}`, DELETE `/cv/{id}`, and POST `/cv` stay identical. Only GET changes behavior (returns array instead of single object) and we add one new GET route. Frontend must update the GET call anyway.

**Final API surface:**
```
GET    /api/v1/cv              → List all user's CVs (was: get single)
GET    /api/v1/cv/{id}         → Get specific CV by ID (NEW)
POST   /api/v1/cv              → Create new CV (unchanged)
PUT    /api/v1/cv/{id}         → Update CV (unchanged)
DELETE /api/v1/cv/{id}         → Delete CV (unchanged)
POST   /api/v1/cv/{id}/export/{format}  → Export specific CV (was: /export/{format})
```

**Decision 2: Export Route**
- Option A: `POST /api/v1/cv/{id}/export/{format}` — nested under CV resource
- Option B: `POST /api/v1/export/{format}` with CV ID in body — minimal route change
- **Choice:** Option A — **Rationale:** Makes the CV ownership explicit in the URL. The old `/export/{format}` route becomes dead (removed). This is more RESTful and the handler already gets `{id}` from chi params, matching the update/delete pattern.

**Decision 3: "Create New CV" Flow**
- Option A: Dashboard → `/templates` → pick template → create CV → `/editor/[id]`
- Option B: Dashboard → create blank CV immediately → `/editor/[id]`
- **Choice:** Option A — **Rationale:** Reuses existing templates page. The templates page becomes a "create CV with template" page. This avoids creating CVs users didn't intend (no more auto-create in editor). Templates page must become a client component to call the API.

**Decision 4: Editor Route Migration**
- `/editor/[id]` becomes the real editor (dynamic route)
- `/editor` (no ID) redirects to `/` (dashboard)
- Header nav: remove "Editor" link (meaningless without a specific CV), keep "New CV" button → `/templates`
- **Choice:** Redirect to dashboard — **Rationale:** Simple, predictable. Users go through dashboard to pick a CV.

**Decision 5: Task Splitting**
- **PR 1 (T-015a — Backend):** Routes, handlers, service, export refactor, all backend tests (~200 lines)
- **PR 2 (T-015b — Frontend):** Dynamic editor route, dashboard multi-CV, API client, templates flow, header nav, seed data (~250 lines)
- **Rationale:** Backend must land first so frontend can test against it. Two PRs keeps each under 300 lines.

### Deferred (Out of Scope)

| Item | Why deferred |
|------|-------------|
| CV duplication | Nice-to-have, not required for multi-CV |
| `is_active` archive/restore UI | Exists in DB but no UI needed yet |
| Dashboard sorting/filtering | Only needed when user has many CVs |
| CV rename in dashboard | Users can rename in editor (title field already editable in content) |
| "Active CV" indicator | Deferred — all CVs are active for now |

### Files NOT Changing

| File | Why |
|------|-----|
| `web/src/stores/editor-store.ts` | Already works with any CV — `setCV(cv)` takes any CV object |
| `web/src/hooks/useAutoSave.ts` | Already uses `currentCV.id` — works correctly with any loaded CV |
| `web/src/components/editor/modals/ExportModal.tsx` | Already uses `cv.id` — just needs the API client URL to change |
| `web/src/types/api.ts` | CV interface already has `id`, `title`, etc. |
| `web/src/types/cv.ts` | No changes needed |
| `api/db/migrations/` | No schema changes — DB already supports multiple CVs |
| `api/internal/repository/repository.go` | Interface already has `GetByID`, `ListByUserID`, `CountByUserID` |
| `api/internal/repository/postgres/cv.go` | Already has `GetByID`, `ListByUserID` — no changes needed |
| `api/internal/domain/cv.go` | No changes needed |
| `api/internal/domain/errors.go` | Already has `ErrForbidden`, `ErrNotFound`, `ErrCVLimitReached` |
| `api/db/queries/cvs.sql` | `GetCVByID` query already exists (line 1-2), `ListCVsByUserID` exists |

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Chi routing conflict between `GET /cv` and `GET /cv/{id}` | Low | High | Chi handles static vs dynamic routes correctly — static routes take priority. Both `GET /cv` (list) and `GET /cv/{id}` coexist fine. |
| Frontend cache serves old `/editor` page after route change | Medium | Medium | Clear `.next` cache after changing route structure. Already a known issue (MEMORY.md). |
| Export tests break due to signature change | Certain | Low | Both handler mocks (`ExportServiceInterface`) and service mocks (`mockCVRepo`) must be updated. Plan includes specific mock changes. |
| Templates page redirect after CV creation fails | Low | Medium | Use `router.push()` after `await api.cv.create()`. Handle errors with try/catch and show error state. |
| Old bookmarks to `/editor` break | Low | Low | Redirect catches this — `/editor` → dashboard. |

### Consequences

**What becomes easier:**
- Users can manage multiple CVs for different job applications
- Each CV is independently editable with its own URL
- Export targets a specific CV — no ambiguity
- Foundation for future features: CV duplication, archiving, comparison

**What becomes harder:**
- "Quick edit" requires 2 clicks (dashboard → CV) instead of 1 (direct to editor)
- Frontend must always know which CV ID to operate on — no more "just get the CV"

---

### PR 1: T-015a — Backend (Routes + Export Refactor)

#### Code Adjustments

**1. `api/internal/handler/cv.go`** — Add `ListCVs` and `GetCVByID` handlers

```go
// ListCVs returns all CVs for the authenticated user
func (h *Handler) ListCVs(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := middleware.MustGetUserID(ctx)

    cvs, err := h.cvService.ListByUserID(ctx, userID)
    if err != nil {
        httputil.HandleError(w, err)
        return
    }

    responses := make([]CVResponse, len(cvs))
    for i, cv := range cvs {
        responses[i] = CVResponse{
            ID:            cv.ID.String(),
            Title:         cv.Title,
            SchemaVersion: cv.SchemaVersion,
            Content:       cv.Content,
            CreatedAt:     cv.CreatedAt.Format("2006-01-02T15:04:05Z"),
            UpdatedAt:     cv.UpdatedAt.Format("2006-01-02T15:04:05Z"),
        }
    }

    httputil.JSON(w, http.StatusOK, responses)
}

// GetCVByID returns a specific CV by ID
func (h *Handler) GetCVByID(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    userID := middleware.MustGetUserID(ctx)

    idParam := chi.URLParam(r, "id")
    cvID, err := uuid.Parse(idParam)
    if err != nil {
        httputil.Error(w, http.StatusBadRequest, "invalid CV ID")
        return
    }

    cv, err := h.cvService.GetByID(ctx, cvID, userID)
    if err != nil {
        httputil.HandleError(w, err)
        return
    }

    response := CVResponse{...}
    httputil.JSON(w, http.StatusOK, response)
}
```

- Remove old `GetCV` handler (replaced by `ListCVs`)

**2. `api/internal/service/cv_service.go`** — Add `ListByUserID`

```go
// ListByUserID returns all CVs for a user
func (s *CVService) ListByUserID(ctx context.Context, userID string) ([]*domain.CV, error) {
    return s.cvRepo.ListByUserID(ctx, userID)
}
```

- `GetByID` already exists with ownership check (lines 36-48) — no change needed

**3. `api/internal/service/export_service.go`** — Change signature from `userID` to `cvID + userID`

- `fetchCVContent(ctx, cvID uuid.UUID, userID string)` — use `cvRepo.GetByID(ctx, cvID)` then check `cv.UserID == userID`
- `renderHTML(ctx, cvID, userID)` — pass through
- `ExportPDF(ctx, cvID, userID)` — pass through
- `ExportDOCX(ctx, cvID, userID)` — pass through

**4. `api/internal/handler/export.go`** — Parse CV ID from URL param `{id}`

- Handler reads `chi.URLParam(r, "id")`, parses UUID, passes to `exportService.ExportPDF(ctx, cvID, userID)`
- Include CV title in Content-Disposition filename: `attachment; filename="{cv.Title}.pdf"`

**5. `api/cmd/server/main.go`** — Update routes (lines 229-242)

```go
// CV routes
r.Get("/cv", h.ListCVs)                          // Changed: list all
r.Get("/cv/{id}", h.GetCVByID)                   // NEW
r.Post("/cv", h.CreateCV)                         // Unchanged
r.Put("/cv/{id}", h.UpdateCV)                     // Unchanged
r.Delete("/cv/{id}", h.DeleteCV)                  // Unchanged

// Export routes — nested under CV
r.Post("/cv/{id}/export/{format}", h.CreateExport)  // Changed: CV ID in path
```

- Remove old `r.Post("/export/{format}", ...)` and `r.Get("/export/{id}", ...)`

**6. Backend Test Updates**

`api/internal/handler/testutil_test.go`:
- Add `ListByUserIDFunc` and `GetByIDFunc` to `MockCVService` + `CVServiceInterface`
- Change `ExportServiceInterface` signatures: `ExportPDF(ctx, cvID uuid.UUID, userID string)`, `ExportDOCX(ctx, cvID uuid.UUID, userID string)`
- Update `MockExportService` to match
- Add `ListCVs` and `GetCVByID` to `TestHandler`
- Update `CreateExport` in `TestHandler` to parse `{id}` from chi params

`api/internal/handler/cv_test.go`:
- Replace `TestGetCV_*` with `TestListCVs_*` (success with 0, 1, multiple CVs) and `TestGetCVByID_*` (success, not found, forbidden)
- Keep `TestCreateCV_*`, `TestUpdateCV_*`, `TestDeleteCV_*` unchanged

`api/internal/handler/export_test.go`:
- Update mock signatures to `(ctx, cvID uuid.UUID, userID string)`
- Pass CV ID in chi route params

`api/internal/service/export_service_test.go`:
- Change `mockCVRepo.GetByUserIDFunc` → `mockCVRepo.GetByIDFunc`
- Pass `cvID uuid.UUID` to `ExportPDF`/`ExportDOCX` instead of `userID` string
- Add ownership check test (wrong user → ErrForbidden)

`api/internal/service/cv_service_test.go`:
- Add `TestCVService_ListByUserID_*` (empty, multiple)
- Existing `GetByID` tests already cover success, not found, wrong user

#### Verification Gates (PR 1)
1. `cd api && go build ./...` — compiles
2. `cd api && go test ./...` — all tests pass
3. `cd api && go vet ./...` — no issues

---

### PR 2: T-015b — Frontend (Editor Route + Dashboard + Templates)

#### Code Adjustments

**1. `web/src/lib/api/client.ts`** — Update `cvApi` and `exportApi`

```typescript
export const cvApi = {
  /** List all user's CVs */
  list: () => client.get<CV[]>('/api/v1/cv'),

  /** Get a specific CV by ID */
  get: (id: string) => client.get<CV>(`/api/v1/cv/${id}`),

  /** Create a new CV */
  create: (data: CreateCVRequest) => client.post<CV>('/api/v1/cv', data),

  /** Update an existing CV */
  update: (id: string, data: UpdateCVRequest) =>
    client.put<CV>(`/api/v1/cv/${id}`, data),

  /** Delete a CV */
  delete: (id: string) => client.delete(`/api/v1/cv/${id}`),
};

export const exportApi = {
  export: async (cvId: string, format: ExportFormat): Promise<Blob> => {
    const url = `${API_BASE}/api/v1/cv/${cvId}/export/${format}`;
    // ... same fetch logic, just URL changed
  },
};
```

**2. `web/src/app/editor/[id]/page.tsx`** — NEW dynamic route

- Copy structure from existing `editor/page.tsx`
- Use `useParams()` to get `id` from URL
- Load CV via `api.cv.get(id)` instead of `api.cv.get()`
- Remove auto-create logic (no more "create if not found" — that happens on templates page)
- If 404 → redirect to dashboard
- If 403 → show error "CV not found or access denied"
- Keep template query param support (`?template=X`) for applying template on load

**3. `web/src/app/editor/page.tsx`** — Convert to redirect

```typescript
'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function EditorRedirectPage() {
  const router = useRouter();
  useEffect(() => { router.replace('/'); }, [router]);
  return null;
}
```

**4. `web/src/app/page.tsx`** — Dashboard multi-CV

- Change `api.cv.get()` → `api.cv.list()`
- State: `const [cvs, setCVs] = useState<CV[]>([])`
- Render `DocumentList` for each CV with `cv.id` in links
- Update edit links: `href={`/editor/${cv.id}`}`
- Update document count: `{cvs.length}` instead of hardcoded `1`
- "Create New Document" → link to `/templates`
- Delete updates `cvs` state by filtering out deleted ID
- Empty state links to `/templates` instead of `/editor`
- Row numbering: `String(index + 1).padStart(2, '0')`

**5. `web/src/app/templates/page.tsx`** — Create CV on template select

- Convert to client component (`'use client'`)
- On template click: call `api.cv.create({ title: 'Untitled CV', content: templateContent(templateId) })` → then `router.push(`/editor/${newCv.id}`)`
- Need to import template content generators
- Show loading state while creating
- Handle errors (quota reached → show message)

**6. `web/src/components/layout/Header.tsx`** — Update navigation

```typescript
const navItems = [
  { href: '/', label: 'Dashboard' },
  { href: '/templates', label: 'Templates' },
  // Remove { href: '/editor', label: 'Editor' }
];
```

- "New CV" button in top-right: change `href="/editor"` → `href="/templates"`
- `isActive` logic for `/editor` routes: `pathname.startsWith('/editor/')` (still highlights when in editor)

Actually, keep "Editor" in nav but handle it differently:
- If user is on `/editor/[id]`, "Editor" highlights
- "Editor" link in nav should go to dashboard (since we can't know which CV to edit)
- Better: just remove "Editor" from the nav items list. The editor is reached through the dashboard.

**7. `api/db/seed/dev_seed.sql`** — Add second sample CV

Add a second CV with `classic` template and different content, using a different UUID. Keep existing CV unchanged for backward compat with existing dev databases.

#### Verification Gates (PR 2)
1. `cd web && npx tsc --noEmit` — type check passes
2. `cd web && npm test` — all tests pass
3. `cd web && npx next build` — build succeeds
4. Manual E2E: dashboard shows multiple CVs, click opens correct editor, create new via templates works

---

### Engineer Execution Steps

#### PR 1: Backend

1. **Update `export_service.go`** — Change `fetchCVContent`, `renderHTML`, `ExportPDF`, `ExportDOCX` signatures to accept `(ctx, cvID uuid.UUID, userID string)`. Use `cvRepo.GetByID` + ownership check.
2. **Update `cv_service.go`** — Add `ListByUserID` method.
3. **Update `cv.go` handler** — Remove `GetCV`. Add `ListCVs` and `GetCVByID`.
4. **Update `export.go` handler** — Parse `{id}` from URL params, pass to export service.
5. **Update `main.go` routes** — Wire new handlers, remove old export route.
6. **Update test mocks** — `testutil_test.go`: add `ListByUserID`/`GetByID` to mock CV service, update export mock signatures. `export_service_test.go`: update `mockCVRepo` to use `GetByID`.
7. **Update handler tests** — `cv_test.go`: replace GetCV tests with ListCVs/GetCVByID tests. `export_test.go`: update to pass CV ID.
8. **Update service tests** — `export_service_test.go`: change to use `GetByID`, add ownership test.
9. **Run `go test ./...`** — verify all pass.
10. **Run `go build ./...` && `go vet ./...`** — verify clean.

#### PR 2: Frontend

1. **Update `client.ts`** — Replace `cvApi.get()` with `cvApi.list()` and `cvApi.get(id)`. Update `exportApi.export` URL.
2. **Create `editor/[id]/page.tsx`** — Dynamic route, load CV by ID from params.
3. **Convert `editor/page.tsx`** — Redirect to `/`.
4. **Update `page.tsx` (dashboard)** — Fetch list, render multiple rows, update all links.
5. **Update `templates/page.tsx`** — Client component, create CV on click, redirect to editor.
6. **Update `Header.tsx`** — Remove "Editor" nav link, update "New CV" button to `/templates`.
7. **Update seed data** — Add second CV.
8. **Run `npx tsc --noEmit`** — type check.
9. **Run `npm test`** — all tests pass.
10. **Run `npx next build`** — verify build.
11. **Manual test** — Start full stack, verify multi-CV flow end-to-end.

## Local QA Review (2026-02-24)

**Verdict**: FAIL

### Static Checks
| Check | Result | Details |
|-------|--------|---------|
| Go tests | PASS | All packages pass |
| Frontend unit tests | PASS | 205/205 passed (19 files) |
| TypeScript | PASS | Zero errors |
| Production build | PASS | Pre-existing warnings only |

### E2E Browser Tests
| Test | Result | Details |
|------|--------|---------|
| Dashboard loads | FAIL | TypeError crash |
| Editor loads | FAIL | Preview stuck on "Loading..." |
| Section editing | SKIP | Blocked by editor failure |
| Template switching | SKIP | Blocked by editor failure |
| Data persistence | SKIP | Blocked |
| Navigation | PASS | Nav links and routing work |

### Findings Requiring Attention

**(blocking) Dashboard crashes — `GET /cv` returns array, frontend expects single object**
- `TypeError: Cannot read properties of undefined (reading 'sections')` at `src/app/page.tsx:141`
- Backend `GET /api/v1/cv` changed to `ListCVs` (returns `{data: [...]}`), but frontend `cvApi.get()` (`web/src/lib/api/client.ts:249`) still expects `{data: {...}}` (single CV)
- Dashboard shows error overlay instead of CV list

**(blocking) Editor preview fails to load CV data**
- Editor shell renders but preview shows "Loading..." permanently, console: `[Store] CV content is missing`
- Same root cause: `cvApi.get()` returns array, editor store can't parse it as a single CV

**(blocking) Export endpoint URL mismatch**
- Frontend `exportApi.export()` (`client.ts:288`) calls `POST /api/v1/export/{format}` (old URL)
- Backend moved to `POST /api/v1/cv/{id}/export/{format}` (new URL)
- Export would 404 if attempted

**Resolution**: These are all expected for a "backend-only first PR" — the frontend updates are planned for PR 2 (T-015b). However, this PR cannot be merged alone without breaking the app. Options:
1. Ship both PRs together (backend + frontend) as a single PR
2. Add backward-compatible shim: keep old `GET /cv` returning single CV alongside new `ListCVs`
3. Merge backend first behind a feature flag (not implemented)
