# T-033: Editor Pane Design System Alignment

## Summary
- **Goal:** Align editor pane (toolbar, sidebars, properties panel, section editors) with the Editorial Craft design system so the app looks visually consistent across all pages
- **Acceptance:** Editor visually matches dashboard/header — same typography conventions, border weights, hover effects, color usage, and component styles
- **Branch:** `feat/t-033-editor-design-alignment`

---

## Investigation Checklist
- [x] Document design system tokens and patterns (CSS vars, tailwind config)
- [x] Screenshot dashboard (reference) and editor (current)
- [x] Read all editor components: EditorToolbar, SectionPalette, PropertiesPanel, CVPreview
- [x] Read all section editors: Personal, Summary, Experience, Education, Skills, Languages, Certifications, Projects
- [x] Read shared UI components: Label, Input, Button, Card
- [x] Compare design system patterns against editor usage
- [x] Identify styling gaps with severity

## Findings

### What's Already Aligned (Good)
The T-025/T-030 design overhaul already updated many editor elements:

| Component | Pattern | Status |
|-----------|---------|--------|
| SectionPalette header | `font-serif text-sm font-semibold`, `border-b-2 border-ink` | Aligned |
| PropertiesPanel header | Same serif header, `border-b-2 border-ink` | Aligned |
| CVPreview container | `bg-secondary p-8`, paper with `shadow-offset-lg animate-paper-drop` | Aligned |
| Sidebar section items | `hover:bg-[hsl(var(--vermillion-pale))]`, accent border on select | Aligned |
| Toolbar icon buttons | vermillion-pale hover, text-secondary color | Aligned |
| Template selector | `border-2 border-input bg-secondary` + focus states | Aligned |
| Save status indicator | `font-mono text-[0.6875rem] uppercase tracking-[0.03em]` | Aligned |
| Export button | Uses `Button` component with proper variant | Aligned |
| Input/Textarea | `border-2 border-input bg-secondary` + focus with `shadow-offset-sm` | Aligned |

### What's NOT Aligned (Gaps)

#### GAP 1: Toolbar border & height mismatch — MEDIUM
- **Current:** `h-[50px] border-b border-dashed border-border` (1px dashed, 50px)
- **Design system:** Header uses `h-[60px] border-b-2 border-ink` (2px solid, 60px)
- **Issue:** Toolbar feels lighter/thinner than the header, creates visual weight imbalance
- **File:** `EditorToolbar.tsx:52`
- **Fix:** Change to `border-b-2 border-ink` (keep height at 50px — toolbar should be compact)

#### GAP 2: Native `<select>` elements use 1px borders — MEDIUM
Three native selects in editors don't match the design system's `border-2` pattern:

| Location | Current | Should Be |
|----------|---------|-----------|
| `PersonalEditor.tsx:159` | `border border-input bg-background` | `border-2 border-input bg-secondary` + focus states |
| `LanguagesEditor.tsx:120` | `border border-input bg-background` | Same |
| `SkillsEditor.tsx:248` | `border rounded bg-background` | Same (smaller variant) |

**Fix:** Align all native selects to `border-2 border-input bg-secondary px-3 text-sm focus:outline-none focus:border-primary focus:bg-card focus:shadow-offset-sm transition-all duration-150`

#### GAP 3: Entry cards use 1px border, no hover effect — MEDIUM
Entry cards (Experience, Languages, Certifications, Projects) use subtle 1px borders and `hover:bg-accent/50` (which renders as a semi-transparent vermillion wash — too strong and inconsistent).

| Location | Current |
|----------|---------|
| `ExperienceEditor.tsx:129` | `rounded-lg border bg-card hover:bg-accent/50` |
| `LanguagesEditor.tsx:103` | Same |
| `CertificationsEditor.tsx` | Same |
| `ProjectsEditor.tsx` | Same |
| `SkillsEditor.tsx:163` (category card) | `rounded-lg border bg-card p-4` |

**Design system pattern (from dashboard document row):**
- `border-b border-dashed border-border` with `hover:bg-[hsl(var(--vermillion-pale))] hover:shadow-offset-md`
- Or for cards: `border-2 border-border` with subtle hover lift

**Fix:** Change to `border-2` borders and use `hover:bg-[hsl(var(--vermillion-pale))]` (subtle tint) instead of `hover:bg-accent/50` (too strong). Optionally add `hover:shadow-offset-sm` for paper lift effect.

#### GAP 4: Skill tags hover to solid accent — LOW
- **Current:** `bg-muted hover:bg-accent` — hovers to solid vermillion red, very aggressive
- **Design system:** Hover states use `hover:bg-[hsl(var(--vermillion-pale))]` (subtle 6% opacity tint)
- **File:** `SkillsEditor.tsx:270`
- **Fix:** Change to `hover:bg-[hsl(var(--vermillion-pale))]`

#### GAP 5: Form field labels lack editorial typography — MEDIUM
- **Current (Label component):** `text-sm font-medium` — plain, no editorial character
- **Design system pattern:** Uses `font-mono text-[0.6875rem] uppercase tracking-[0.05em]` for section labels
- **Issue:** The Label component is shared across the entire app, so changing it directly would affect non-editor uses
- **File:** `label.tsx:8-10`
- **Recommendation:** Don't change the global Label component. Instead, apply editorial label styling directly in properties panel editors using a utility class or wrapper. The properties panel headers already use `font-serif` — field labels should use `font-mono uppercase tracking-[0.05em]` to match the dashboard's metadata treatment.

#### GAP 6: Summary editor uses raw Tailwind colors — LOW
- **Current:** `text-amber-500` and `text-green-500` for character count
- **Design system:** Uses `text-accent` (vermillion), `text-sienna`, `text-[hsl(var(--text-secondary))]`
- **File:** `SummaryEditor.tsx:47-52`
- **Fix:** Use `text-sienna` for warning range, `text-accent` or keep `text-green-500` for good range (acceptable divergence since green feedback is universally understood)

#### GAP 7: Empty state styling inconsistency — LOW
Empty states across editors use:
- `border border-dashed rounded-lg` (1px dashed)
- No icon consistency (some have icons, some don't)
- No editorial copy style

**Fix:** Standardize to `border-2 border-dashed border-border rounded-lg` with consistent muted icon + editorial-style "no items" text.

#### GAP 8: Properties panel content area lacks section grouping — LOW
- The properties panel form fields flow as a flat list
- Dashboard uses clear section grouping with dashed dividers
- Could add `border-t border-dashed border-border` between logical field groups
- **Recommendation:** Nice-to-have, skip for this task unless it's trivial

---

## Gaps Summary

| # | Gap | Severity | Files |
|---|-----|----------|-------|
| 1 | Toolbar border weight | MEDIUM | `EditorToolbar.tsx` |
| 2 | Native selects 1px borders | MEDIUM | `PersonalEditor.tsx`, `LanguagesEditor.tsx`, `SkillsEditor.tsx` |
| 3 | Entry cards 1px border + wrong hover | MEDIUM | `ExperienceEditor.tsx`, `LanguagesEditor.tsx`, `CertificationsEditor.tsx`, `ProjectsEditor.tsx`, `SkillsEditor.tsx` |
| 4 | Skill tag hover too strong | LOW | `SkillsEditor.tsx` |
| 5 | Form labels lack editorial typography | MEDIUM | All section editors (8 files) |
| 6 | Raw Tailwind colors in Summary | LOW | `SummaryEditor.tsx` |
| 7 | Empty state border weight | LOW | Multiple editors |
| 8 | No section grouping dividers | LOW | Properties panel editors |

---

## Architect Handoff

### Decisions Needed
1. **Label treatment approach (Quick Decision)**
   - **Option A:** Create a `PropertiesLabel` component that wraps `<Label>` with editorial classes — cleanest, no global impact
   - **Option B:** Add a `variant="editorial"` to the existing `Label` CVA — more reusable but heavier
   - **Option C:** Just add classes inline in each editor — fastest but repetitive
   - **Recommendation:** Option A — small wrapper component used only in properties panel editors

2. **Entry card hover pattern (Quick Decision)**
   - **Option A:** `hover:bg-[hsl(var(--vermillion-pale))]` only (subtle tint)
   - **Option B:** Tint + `hover:shadow-offset-sm` (paper lift, matches dashboard)
   - **Recommendation:** Option A for cards (less visual noise in a compact sidebar), Option B only if it looks good

### Files in Scope

| File | Changes | Gap |
|------|---------|-----|
| `EditorToolbar.tsx` | Change toolbar border to `border-b-2 border-ink` | #1 |
| `PersonalEditor.tsx` | Update native select styling, add editorial labels | #2, #5 |
| `SummaryEditor.tsx` | Replace raw colors with design system tokens, editorial labels | #5, #6 |
| `ExperienceEditor.tsx` | Update entry card borders + hover, editorial labels | #3, #5 |
| `EducationEditor.tsx` | Same as Experience | #3, #5 |
| `SkillsEditor.tsx` | Update category card borders, skill tag hover, native select, editorial labels | #2, #3, #4, #5 |
| `LanguagesEditor.tsx` | Update entry row borders + hover, native select, editorial labels | #2, #3, #5 |
| `CertificationsEditor.tsx` | Update entry card borders + hover, editorial labels | #3, #5 |
| `ProjectsEditor.tsx` | Update entry card borders + hover, editorial labels | #3, #5 |
| (Optional) New `PropertiesLabel` or utility | Editorial label wrapper | #5 |

### Tests to Add/Update
- No new tests needed — these are purely visual/CSS changes
- Existing 200+ tests should continue passing
- E2E visual verification via Playwright screenshots during `/local-qa`

### Constraints
- No changes to global UI components (Button, Input, Card, Label) — editor-specific only
- Keep PR focused on CSS classes — no structural/logic changes
- Verify dark mode still works after changes

### Recommended Next Agent
**@architect** for quick decisions on label treatment and hover pattern, then **@engineer** for implementation.

---

## Test Plan
- Run existing test suite (`npm test` in web/) — all 200+ tests must pass
- Visual E2E: screenshot dashboard and editor side-by-side, verify consistent visual language
- Dark mode: verify all changes render correctly in `.dark` mode
- Responsive: verify changes don't break at md/lg breakpoints where sidebars collapse

---

## Architect Plan

### Decisions

**Decision 1: Label Treatment Approach**
- Option A: Create `PropertiesLabel` wrapper component — clean, no global impact, but adds a new file + requires changing every `<Label>` in 8 editors (~45 replacements)
- Option B: Add `variant="editorial"` to Label CVA — reusable, but still ~45 prop additions across editors
- Option C: Parent CSS selector on PropertiesPanel wrapper — ONE line change in ONE file, targets all `<label>` elements within properties content area via Tailwind `[&_label]` selector
- **Choice: C** — **Rationale:** Minimal diff (1 line), zero changes to individual editors or Label component, easily reversible. The `[&_label]` selector on `PropertiesPanel.tsx:77` only targets `<label>` HTML elements within the content area — section headers use `<h2>/<h3>` and buttons use `<span>`, so there's no false positive risk.

**Decision 2: Entry Card Hover Pattern**
- Option A: Tint only — `hover:bg-[hsl(var(--vermillion-pale))]` (subtle 6% opacity vermillion wash)
- Option B: Tint + shadow — add `hover:shadow-offset-sm` for paper-lift effect (matches dashboard rows)
- **Choice: A** — **Rationale:** Entry cards live in a compact 320px sidebar. Shadow-lift would add visual clutter in tight space. The dashboard rows have 800px+ of horizontal room where shadows read well. Tint alone provides consistent editorial feedback without noise.

**Decision 3: GAP 8 (Section Grouping Dividers) — Skip**
- Investigation flagged this as "nice-to-have, skip unless trivial"
- Adding `border-t border-dashed` between field groups requires understanding the logical grouping within each editor — non-trivial to get right
- **Choice: Skip** — out of scope for this task, can be a follow-up if stakeholder wants it

### Gaps Addressed

| # | Gap | Severity | Decision |
|---|-----|----------|----------|
| 1 | Toolbar border weight | MEDIUM | Fix: `border-b-2 border-ink` |
| 2 | Native selects 1px borders | MEDIUM | Fix: match template selector pattern |
| 3 | Entry cards 1px border + wrong hover | MEDIUM | Fix: `border-2` + vermillion-pale hover |
| 4 | Skill tag hover too strong | LOW | Fix: vermillion-pale hover |
| 5 | Form labels lack editorial typography | MEDIUM | Fix: parent CSS `[&_label]` on PropertiesPanel |
| 6 | Raw Tailwind colors in Summary | LOW | Fix: `text-sienna` for warnings, keep `text-green-500` for good |
| 7 | Empty state border weight | LOW | Fix: `border-2 border-dashed border-border` |
| 8 | No section grouping dividers | LOW | **SKIP** — out of scope |

### Code Adjustments by File

#### 1. `EditorToolbar.tsx` — GAP 1
- **Line 52:** `border-b border-dashed border-border` → `border-b-2 border-ink`
- Rationale: Match header's 2px solid ink border for visual weight consistency

#### 2. `PropertiesPanel.tsx` — GAP 5
- **Line 77:** Add editorial label classes to the content wrapper div:
  ```
  <div className="flex-1 overflow-y-auto p-4 [&_label]:font-mono [&_label]:text-[0.6875rem] [&_label]:uppercase [&_label]:tracking-[0.05em]">
  ```
- Rationale: All `<label>` elements within properties panel get editorial mono/uppercase treatment. No individual editor changes needed.

#### 3. `PersonalEditor.tsx` — GAP 2, GAP 7
- **Line 159 (link type select):** `border border-input bg-background` → `border-2 border-input bg-secondary text-sm focus:outline-none focus:border-primary focus:bg-card focus:shadow-offset-sm transition-all duration-150`
- **Line 186 (empty state):** `border border-dashed rounded-lg` → `border-2 border-dashed border-border rounded-lg`

#### 4. `SummaryEditor.tsx` — GAP 6
- **Line 48:** `text-amber-500` → `text-sienna` (under minimum)
- **Line 50:** `text-amber-500` → `text-sienna` (over maximum)
- **Line 51:** Keep `text-green-500` — no design system green token, and green universally signals "good"

#### 5. `ExperienceEditor.tsx` — GAP 3, GAP 7
- **Line 90 (empty state):** `border border-dashed rounded-lg` → `border-2 border-dashed border-border rounded-lg`
- **Line 129 (entry card):** `rounded-lg border bg-card hover:bg-accent/50` → `rounded-lg border-2 bg-card hover:bg-[hsl(var(--vermillion-pale))]`

#### 6. `EducationEditor.tsx` — GAP 3, GAP 7
- **Line 91 (empty state):** Same as Experience
- **Line 131 (entry card):** Same pattern as Experience

#### 7. `SkillsEditor.tsx` — GAP 2, GAP 3, GAP 4, GAP 7
- **Line 115 (empty state):** `border border-dashed rounded-lg` → `border-2 border-dashed border-border rounded-lg`
- **Line 163 (category card):** `rounded-lg border bg-card p-4` → `rounded-lg border-2 bg-card p-4`
- **Line 248 (skill level select):** `border rounded px-1 bg-background` → `border-2 border-input rounded px-1 bg-secondary text-sm focus:outline-none focus:border-primary focus:shadow-offset-sm transition-all duration-150`
- **Line 270 (skill tag):** `hover:bg-accent` → `hover:bg-[hsl(var(--vermillion-pale))]`

#### 8. `LanguagesEditor.tsx` — GAP 2, GAP 3, GAP 7
- **Line 71 (empty state):** Same pattern
- **Line 103 (entry row):** `rounded-lg border bg-card hover:bg-accent/50` → `rounded-lg border-2 bg-card hover:bg-[hsl(var(--vermillion-pale))]`
- **Line 120 (proficiency select):** `border border-input bg-background px-3 text-sm` → `border-2 border-input bg-secondary px-3 text-sm focus:outline-none focus:border-primary focus:bg-card focus:shadow-offset-sm transition-all duration-150`

#### 9. `CertificationsEditor.tsx` — GAP 3, GAP 7
- **Line 86 (empty state):** Same pattern
- **Line 126 (entry card):** Same as Experience

#### 10. `ProjectsEditor.tsx` — GAP 3, GAP 7
- **Line 89 (empty state):** Same pattern
- **Line 129 (entry card):** Same as Experience

### Files NOT Changing

| File | Reason |
|------|--------|
| `label.tsx` | Global component — editorial treatment applied via parent CSS instead |
| `button.tsx` | Already aligned (T-025) |
| `input.tsx` / `textarea.tsx` | Already use `border-2 border-input` (T-025) |
| `SectionPalette.tsx` | Already aligned — serif headers, vermillion-pale hovers, accent borders |
| `CVPreview.tsx` | Already aligned — paper shadow, secondary bg |
| `EditorLayout.tsx` | No styling gaps found |

### Risk Assessment

| Risk | Likelihood | Mitigation |
|------|-----------|------------|
| `[&_label]` selector targets unintended elements | Low | Verified: only `<label>` HTML elements exist in editor content area; headers are `<h2>/<h3>`, buttons are `<span>` |
| `border-2` on compact cards looks too heavy | Low | Consistent with inputs/buttons that already use `border-2`; if heavy, can revert individual cards |
| `text-sienna` warning color isn't visible enough | Low | Sienna (HSL 26.8, 38.1%, 43.7%) has good contrast on both light/dark; engineer verifies visually |
| Dark mode regression | Low | All changes use CSS custom properties that have dark mode variants (`--vermillion-pale` is 12% opacity in dark) |

### Consequences

**What becomes easier:**
- Future editor components automatically get editorial labels (just use `<Label>` inside PropertiesPanel)
- Consistent hover pattern (`vermillion-pale`) across dashboard + editor reduces cognitive load for users

**What becomes harder:**
- Nothing — all changes are additive CSS class modifications

### Engineer Execution Steps

1. **Update `EditorToolbar.tsx`** — Change toolbar bottom border (GAP 1)
2. **Update `PropertiesPanel.tsx`** — Add `[&_label]` editorial classes to content wrapper (GAP 5)
3. **Update `PersonalEditor.tsx`** — Fix native select + empty state (GAP 2, 7)
4. **Update `SummaryEditor.tsx`** — Replace raw colors (GAP 6)
5. **Update `ExperienceEditor.tsx`** — Fix entry card + empty state (GAP 3, 7)
6. **Update `EducationEditor.tsx`** — Same as Experience (GAP 3, 7)
7. **Update `SkillsEditor.tsx`** — Fix category card, skill tag hover, level select, empty state (GAP 2, 3, 4, 7)
8. **Update `LanguagesEditor.tsx`** — Fix entry row, proficiency select, empty state (GAP 2, 3, 7)
9. **Update `CertificationsEditor.tsx`** — Fix entry card + empty state (GAP 3, 7)
10. **Update `ProjectsEditor.tsx`** — Fix entry card + empty state (GAP 3, 7)

### Verification Gates

1. `cd web && npm test` — all 200+ tests pass (no logic changes)
2. `cd web && npx tsc --noEmit` — type check passes
3. `cd web && npx next build` — build succeeds
4. Visual E2E — screenshot editor, compare against dashboard for consistent visual language
5. Dark mode — toggle to `.dark` and verify all changes render correctly

---

## QA Test Review

### Existing Test Coverage Audit

**18 test files reviewed** across editors, stores, hooks, API client, and page components.

**Tests affected by T-033 changes: 0** — All changes are CSS class modifications only. No logic, props, or component structure changes.

**One brittle selector noted (pre-existing, out of scope):**
- `SkillsEditor.test.tsx:162` uses `btn.closest('[class*="rounded-md bg-muted"]')` — depends on CSS class names. Our changes preserve `bg-muted` on skill tags (only `hover:bg-accent` → `hover:bg-[hsl(var(--vermillion-pale))]`), so this test is safe. However, this selector pattern is fragile and could break from unrelated Tailwind changes in the future.

### QA Assessment

**No new unit tests needed.** This is correct and appropriate because:
1. All changes are CSS class swaps — no testable logic
2. Existing 200+ tests validate that component behavior/interactions are unaffected
3. Visual correctness is only verifiable through E2E screenshots, not unit assertions

### QA-Recommended Additional Checks

| # | Check | Layer | Priority | Description |
|---|-------|-------|----------|-------------|
| 1 | Existing test suite passes | Unit | High | `cd web && npm test` — all 200+ tests must pass unchanged. Zero test modifications expected. |
| 2 | Visual E2E: editor with entries | E2E | High | Navigate to editor with seeded CV data. Screenshot each section type (Personal, Experience, Education, Skills, Languages, Certifications, Projects). Verify: 2px borders on cards, vermillion-pale hover on card mouseover, editorial mono/uppercase labels. |
| 3 | Visual E2E: empty states | E2E | Medium | Add a new empty section (e.g. Experience with 0 entries). Verify empty state has `border-2 border-dashed` treatment. |
| 4 | Visual E2E: toolbar alignment | E2E | Medium | Compare toolbar bottom border (2px solid ink) against header bottom border — should match visually in weight and color. |
| 5 | Visual E2E: Summary char count colors | E2E | Medium | Type in Summary editor — verify character count shows sienna (warm brown) when outside recommended range, green when within. |
| 6 | Dark mode regression | E2E | High | Toggle dark mode and re-check all above items. All CSS custom properties (`--vermillion-pale`, `--ink`, `--sienna`) have dark variants. |
| 7 | Native select focus states | E2E | Medium | Tab through native selects in Personal (link type), Languages (proficiency), Skills (level). Verify focus ring shows `border-primary` + `shadow-offset-sm`. |
| 8 | `[&_label]` scope verification | E2E | Medium | Verify that editorial label styling (mono, uppercase) applies ONLY within properties panel, NOT in toolbar, sidebar, or other areas. Check that header text (`<h2>`, `<h3>`) in properties panel is NOT affected. |

### Concerns

- **No automated visual regression** — This task depends entirely on manual E2E visual verification during `/local-qa`. Acceptable for a one-time design alignment PR.
- **Sienna visibility** — `text-sienna` (HSL 26.8, 38.1%, 43.7%) should be verified visually against the summary editor's light background. If it blends too much, `text-accent` (vermillion) is a fallback.
