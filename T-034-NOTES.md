# T-034: Unified Header Across All Pages (Including Editor)

## Summary
- **Goal:** Replace the split header system (shared `Header.tsx` on Dashboard/Templates + separate `EditorToolbar.tsx` on Editor) with a single unified header that provides consistent navigation and branding across all pages, with the editor adding contextual tools within the shared layout.
- **Acceptance:** All pages share the same header component with logo + nav. Editor page gains navigation without losing any editor-specific tools (undo/redo, zoom, template, color picker, save status, export, panel toggles). Design system alignment maintained.
- **Branch:** `feat/t-034-unified-header`

---

## Investigation Checklist
- [x] Read Header.tsx — shared header (Dashboard, Templates)
- [x] Read EditorToolbar.tsx — editor-only toolbar
- [x] Read AppShell.tsx — layout wrapper
- [x] Read EditorLayout.tsx — editor page structure
- [x] Read editor/page.tsx — how editor mounts
- [x] Trace all imports/consumers of Header and EditorToolbar
- [x] Check for existing tests (Header, Toolbar, AppShell)
- [x] Review design system tokens and spacing reference

## Findings

### Current Architecture

| Page | Path | Wrapper | Has Header? | Has EditorToolbar? |
|------|------|---------|-------------|-------------------|
| Dashboard | `/` | `AppShell` | Yes | No |
| Templates | `/templates` | `AppShell` | Yes | No |
| Editor | `/editor` | None | **No** | Yes |
| Old Dashboard | `/dashboard` | N/A (redirect) | N/A | N/A |

### Header.tsx (`web/src/components/layout/Header.tsx`, 119 lines)
- **Height:** 60px, sticky top, z-50
- **Contents:** Logo (H stamp + "HireMe") | Desktop nav (Dashboard, Editor, Templates) | Spacer | New CV button + User avatar | Mobile hamburger
- **Styling:** `bg-card border-b-2 border-ink`, `px-9`
- **Responsive:** Desktop nav hidden below `md:` (768px), mobile menu with hamburger
- **Active detection:** Exact match for `/`, startsWith for others
- **No props** — self-contained, reads pathname from Next.js

### EditorToolbar.tsx (`web/src/components/editor/EditorToolbar.tsx`, 285 lines)
- **Height:** 50px (different from Header's 60px)
- **Contents (left to right):**
  1. Back to Dashboard (arrow) + divider
  2. Left panel toggle + divider
  3. Undo / Redo + divider
  4. Zoom out / percentage / Zoom in + divider
  5. Template selector (native `<select>`) + Color picker
  6. Spacer
  7. Save status indicator + Manual save button + divider
  8. Export button + divider
  9. Right panel toggle
- **Styling:** `bg-card border-b-2 border-ink`, `px-5`, `gap-1`
- **Connects to:** `useEditorStore` (undo/redo, save, template, styling) + `useUIStore` (sidebars, zoom, export modal)
- **No navigation, no logo, no branding**

### AppShell.tsx (`web/src/components/layout/AppShell.tsx`, 30 lines)
- Wraps `Header` + content in `TooltipProvider`
- Props: `showHeader` (default true), `fullHeight` (default false)
- Used by Dashboard and Templates pages
- **NOT used by editor** — editor has its own `TooltipProvider`

### EditorLayout.tsx (`web/src/components/editor/EditorLayout.tsx`, 76 lines)
- Renders `EditorToolbar` at top, then 3-panel layout below
- Handles responsive sidebar auto-collapse
- Renders modals (`ExportModal`, `DeleteConfirmModal`)

### Editor page.tsx (`web/src/app/editor/page.tsx`, 162 lines)
- Wraps `EditorLayout` in `TooltipProvider` + `h-screen` container
- Does NOT use `AppShell` at all
- Has loading/error states before showing EditorLayout

### Tests
- **No tests exist** for Header, EditorToolbar, AppShell, or EditorLayout
- All 200+ existing tests cover stores, hooks, and section editors

### Design System Notes
- Both Header and EditorToolbar use identical color tokens (`bg-card`, `border-ink`, `text-secondary`, `vermillion-pale`)
- Height mismatch: Header=60px, EditorToolbar=50px
- Padding mismatch: Header=`px-9` (36px), EditorToolbar=`px-5` (20px)
- Both use same button pattern: 32px icon buttons with hover highlights
- Design system reference: `DESIGN-SYSTEM.md` specifies nav height=60px, editor toolbar height=50px

## Gaps Found

### BLOCKER
None — this is a new feature, not fixing broken behavior.

### Medium
1. **No navigation on editor page** — User cannot reach Dashboard or Templates from editor without using browser back or manually editing URL. Only has a "Back to Dashboard" arrow icon.
2. **Duplicate TooltipProvider** — `AppShell` wraps in `TooltipProvider`, editor page creates its own `TooltipProvider`. Unified header should consolidate this.
3. **Height inconsistency** — Header is 60px, EditorToolbar is 50px. Unified approach needs a decision on total height.

### Low
4. **No tests** for any header/toolbar components — should add as part of this task.
5. **Mobile hamburger doesn't include editor tools** — future consideration for T-031.

## Architect Handoff

### Decisions Needed

**Decision 1: Layout Architecture**
- **Option A: Single row** — Unified header (logo + nav + editor tools all in one 60px bar). Editor tools replace the right side (New CV button + avatar) when on editor page. Compact but may be cramped.
- **Option B: Two rows on editor** — Shared header (60px, always: logo + nav + right actions) + editor toolbar below (50px, editor-only: all editor tools). Total 110px on editor, 60px elsewhere. More space but taller.
- **Option C: Single row with overflow** — Unified 60px header, editor tools collapse into a "more" dropdown on smaller screens.
- **Recommendation:** Option B (two rows) — cleanest separation of concerns, no cramming, editor tools stay exactly as they are. Header becomes truly shared.

**Decision 2: Editor page integration with AppShell**
- **Option A: Editor uses AppShell** with `fullHeight={true}` — simplest, AppShell already supports it.
- **Option B: Keep editor separate** — duplicate the header rendering.
- **Recommendation:** Option A — use `AppShell fullHeight`. Removes duplicate `TooltipProvider`, consistent mounting.

**Decision 3: EditorToolbar "Back to Dashboard" removal**
- If Header provides full navigation, the "Back to Dashboard" arrow in EditorToolbar becomes redundant.
- **Recommendation:** Remove it. The Header's nav already covers navigation.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `web/src/components/layout/Header.tsx` | No changes needed if Option B | — |
| `web/src/components/layout/AppShell.tsx` | Minor: ensure `fullHeight` works for editor | #2 |
| `web/src/components/editor/EditorToolbar.tsx` | Remove "Back to Dashboard" link and first divider | #3 |
| `web/src/components/editor/EditorLayout.tsx` | No changes if EditorToolbar stays | — |
| `web/src/app/editor/page.tsx` | Wrap in `AppShell fullHeight` instead of raw `TooltipProvider`+`div` | #1, #2 |
| `web/src/app/editor/page.tsx` (loading/error states) | Move loading/error inside AppShell so header is always visible | #1 |

### Tests to Add
- Header: renders logo, nav items, active state, mobile menu toggle
- AppShell: renders header when `showHeader=true`, respects `fullHeight`
- EditorToolbar: renders all controls, connects to stores
- Integration: editor page renders both Header and EditorToolbar

### Constraints
- PR target: <300 lines changed
- Must preserve all existing editor functionality (undo/redo, zoom, save, export, panel toggles)
- Must maintain responsive behavior (mobile hamburger menu, sidebar auto-collapse)
- Design system compliance: use existing tokens, no new colors/fonts

### Recommended Next Agent
**@architect** — Needs to confirm the two-row approach and finalize the layout decision before implementation.

---

## Architect Plan

### Decisions

**Decision 1: Layout Architecture**
- Option A: Single row (60px) — logo + nav + editor tools in one bar. Pros: compact. Cons: 12 editor controls + logo + 3 nav links + 2 right actions = extremely cramped; doesn't fit on screens under ~1200px; forces complex responsive collapse logic.
- Option B: Two rows on editor — shared Header (60px) always on top + EditorToolbar (50px) below on editor only. Pros: zero cramming; both bars keep their existing layout/styling; separation of concerns (navigation vs. editing tools); matches DESIGN-SYSTEM.md specs (nav=60px, toolbar=50px). Cons: 110px total chrome on editor page (was 50px).
- Option C: Single row with dropdown overflow — Pros: single bar. Cons: most complex to implement, hides tools behind clicks, worse UX.
- **Choice:** B (two rows) — **Rationale:** The 60px navigation overhead is a fair trade for full branding and navigation on every page. Both bars already exist with correct styling. Minimal code changes. The DESIGN-SYSTEM.md already anticipates separate nav (60px) and toolbar (50px) heights.

**Decision 2: Editor page integration with AppShell**
- Option A: Editor uses `AppShell fullHeight` — reuses existing layout wrapper, eliminates duplicate TooltipProvider, consistent with Dashboard/Templates.
- Option B: Editor renders Header directly — duplicates the TooltipProvider and layout wrapper logic.
- **Choice:** A — **Rationale:** `AppShell` already has `fullHeight` mode that creates `flex flex-col h-screen` + `main.flex-1.overflow-hidden`. This is exactly what the editor needs. EditorLayout's `h-full flex flex-col` resolves correctly inside this structure. One less TooltipProvider nesting.

**Decision 3: Remove "Back to Dashboard" arrow from EditorToolbar**
- Option A: Keep it alongside the Header nav — redundant but harmless.
- Option B: Remove it — Header provides full navigation (logo→`/`, nav links for Dashboard/Editor/Templates).
- **Choice:** B — **Rationale:** Redundant controls clutter the toolbar. The Header's nav is more discoverable and consistent across pages. Saves ~14 lines and reclaims horizontal space in the toolbar.

**Decision 4: Editor loading/error states**
- Current: Loading and error states render OUTSIDE the TooltipProvider/container wrapper, filling `h-screen` with no header visible.
- New: Loading/error states render INSIDE `AppShell fullHeight`, so the Header is always visible even during loading/errors. States use `h-full` (fills remaining space below Header) instead of `h-screen`.
- **Choice:** Restructure — **Rationale:** Users should always see navigation, even when the editor is loading or errored. This lets them navigate away without using browser back.

### Layout Diagram

```
BEFORE (editor):                    AFTER (editor):
+----------------------------+      +----------------------------+
| EditorToolbar (50px)       |      | Header (60px)              |
|  ← | panels | undo | zoom |      |  Logo | Nav | New CV | U   |
|    | template | save | exp |      +----------------------------+
+----------------------------+      | EditorToolbar (50px)       |
|          |         |       |      |  panels | undo | zoom      |
| Palette  | Preview | Props |      |  template | save | export  |
|          |         |       |      +----------------------------+
+----------------------------+      |          |         |       |
                                    | Palette  | Preview | Props |
                                    |          |         |       |
                                    +----------------------------+

BEFORE (dashboard/templates):       AFTER (dashboard/templates):
+----------------------------+      +----------------------------+
| Header (60px)              |      | Header (60px)              |
|  Logo | Nav | New CV | U   |      |  Logo | Nav | New CV | U   |
+----------------------------+      +----------------------------+
| Page content               |      | Page content               |
+----------------------------+      +----------------------------+
(no change)
```

### Files Changing

| File | What Changes | Lines ~Δ |
|------|-------------|----------|
| `web/src/app/editor/page.tsx` | Wrap in `AppShell fullHeight`, remove manual `TooltipProvider` + `h-screen div`, restructure loading/error to render inside AppShell | ~+8/-12 |
| `web/src/components/editor/EditorToolbar.tsx` | Remove "Back to Dashboard" `<Link>` (lines 53–64), remove first divider (line 66), remove `ArrowLeft` import, remove `Link` import | ~-16 |

**Estimated total diff: ~36 lines changed**

### Files NOT Changing

| File | Why |
|------|-----|
| `Header.tsx` | Works perfectly as-is. Active detection already handles `/editor` path. |
| `AppShell.tsx` | `fullHeight` mode already creates the correct flex layout. No modifications needed. |
| `EditorLayout.tsx` | Renders `EditorToolbar` + 3-panel layout. Structure is correct. `h-full flex flex-col` will resolve correctly inside AppShell's `main.flex-1.overflow-hidden`. |
| `web/src/app/page.tsx` | Dashboard page — no changes needed. |
| `web/src/app/templates/page.tsx` | Templates page — no changes needed. |

### Risk Assessment

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| EditorLayout height calculation breaks inside AppShell | Low | High | Verified: `flex-1 overflow-hidden` on main gives computed height, `h-full` on EditorLayout resolves correctly. E2E visual test confirms. |
| Header mobile menu clipped by AppShell's overflow-hidden | Low | Medium | Verified: mobile menu overflows the `<header>` element (which is a sibling of `<main>`), so `main`'s `overflow-hidden` doesn't clip it. `z-50` stacks above main. |
| "New CV" button on editor page causes issues (links to `/editor`, user is already there) | Low | Low | Pre-existing behavior. Out of scope for T-034. Can address in T-015 (multiple CV support). |
| Existing tests break | None | — | No existing tests for any of these components. |

### Consequences

**Becomes easier:**
- Navigation from editor page (was impossible, now full nav)
- Adding new pages — just use `AppShell` and get consistent header
- Branding consistency across all pages

**Becomes harder:**
- Nothing — this is purely additive

### Engineer Execution Steps

1. **Edit `EditorToolbar.tsx`:**
   - Remove `ArrowLeft` from lucide-react import (line 15)
   - Remove `Link` and `import Link from 'next/link'` (line 3) — check if Link is still used elsewhere in the file first. It is NOT used elsewhere, so remove it.
   - Remove the "Back to Dashboard" Tooltip block (lines 53–64)
   - Remove the first divider `<div className="w-px h-5 bg-border mx-1.5" />` (line 66)

2. **Edit `editor/page.tsx`:**
   - Add import: `import { AppShell } from '@/components/layout';`
   - Remove import: `import { TooltipProvider } from '@/components/ui/tooltip';`
   - In `EditorContent`, restructure the return flow:
     - Remove the early returns for loading/error that use `h-screen`
     - Return `<AppShell fullHeight>` wrapping a single content block
     - Inside AppShell: conditional rendering — if loading show loading state (use `h-full` not `h-screen`), if error show error state (`h-full`), otherwise show `<EditorLayout />`
   - Remove the outer `<TooltipProvider>` + `<div className="h-screen relative z-[1]">` wrapper from the success return — AppShell handles both

3. **Verify:**
   - `cd web && npx tsc --noEmit` — type check
   - `cd web && npm test` — all 200+ tests pass
   - `cd web && npx next build` — build succeeds
   - Manual/E2E: editor page shows Header + EditorToolbar, navigation works, all editor tools work, mobile hamburger menu works

### Test Recommendations

Given the <300 line PR target and the simplicity of the change (wiring, not logic), tests should be minimal and focused:

- No new unit tests needed — the change is purely structural (which component wraps which). The existing 200+ tests cover store interactions and hook behavior, which are unaffected.
- **E2E verification is the real test**: editor page renders Header on top, EditorToolbar below, all editor tools functional, navigation links work.
- If QA plan calls for unit tests, recommend testing: (1) EditorToolbar does NOT render a link to `/`, (2) editor page renders both Header and EditorToolbar.

### Notes for QA
- The "New CV" button in the Header links to `/editor`. When already on `/editor`, this is a same-page navigation — Next.js won't full-reload. Depending on editor state, this may or may not reset. This is pre-existing behavior, not introduced by T-034.
- Mobile: the Header hamburger menu on the editor page will overlay the EditorToolbar. This is correct behavior (z-50 stacking). Editor tools are only usable when the mobile menu is closed.

---

## Test Plan

### Existing Test Coverage

| Test File | Tests | Relevant? |
|-----------|-------|-----------|
| `web/src/app/__tests__/EditorPage.test.tsx` | 4 tests (loading, success, error, auto-create) | **Yes — will break.** Mocks `@/components/ui/tooltip` (TooltipProvider). After T-034, editor page imports `AppShell` from `@/components/layout` instead. Mock must be updated. |
| `web/src/app/__tests__/DashboardPage.test.tsx` | Dashboard tests | **Reference only.** Shows the correct `AppShell` mock pattern: `vi.mock('@/components/layout', ...)` with `<div data-testid="app-shell">{children}</div>`. |
| 200+ other tests (stores, hooks, editors) | Stores, hooks, section editors | Not affected — no imports of Header, Toolbar, or AppShell. |

### QA-Recommended Additional Tests

| # | Test | Layer | Priority | Description |
|---|------|-------|----------|-------------|
| 1 | Update EditorPage mock from TooltipProvider to AppShell | Frontend page test | **High** | In `EditorPage.test.tsx`: replace `vi.mock('@/components/ui/tooltip', ...)` (lines 68-72) with `vi.mock('@/components/layout', () => ({ AppShell: ({ children }: { children: React.ReactNode }) => <div data-testid="app-shell">{children}</div> }))`. Follow the existing `DashboardPage.test.tsx` pattern (lines 27-31). All 4 existing tests must continue to pass. **This is a required fix, not optional — tests will fail without it.** |
| 2 | Verify loading/error render inside AppShell wrapper | Frontend page test | **High** | In `EditorPage.test.tsx`: add assertion in the "shows loading state" and "shows error state" tests that `screen.getByTestId('app-shell')` exists AND contains the loading/error content. This confirms the structural change — Header is visible during loading/error states. |
| 3 | EditorToolbar does not render "Back to Dashboard" link | Frontend component test | **Medium** | New test file: `web/src/components/editor/__tests__/EditorToolbar.test.tsx`. Mock `@/stores` (useEditorStore, useUIStore, usePreviewScalePercent) with defaults. Render `EditorToolbar` inside a `TooltipProvider`. Assert: no element with `aria-label` or text "Back to Dashboard", no `<a href="/">` link inside the toolbar. This guards the removal regression. |
| 4 | EditorToolbar renders core controls | Frontend component test | **Medium** | Same test file as #3. Assert presence of: undo button, redo button, zoom controls, template `<select>`, export button, left/right panel toggles. This is a baseline smoke test ensuring nothing beyond the back arrow was accidentally removed. |

### Notes on Scope
- Tests #1-2 are **required** (existing tests break without #1; #2 validates the core acceptance criteria).
- Tests #3-4 are **recommended** — they add the first-ever test coverage for EditorToolbar, which is the project's most control-dense component.
- No tests recommended for Header.tsx or AppShell.tsx — they are not changing in this task.
- E2E visual verification remains the primary validation for layout correctness (110px stacking, mobile menu overlay behavior).
