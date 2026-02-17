# T-023: README Overhaul — Content + Design

## Summary
- **Goal:** Update README.md to accurately reflect current project state and improve visual design/structure
- **Acceptance:** README has correct info, no dead links, professional design, good first impression for visitors
- **Branch:** `docs/t-023-readme-overhaul`

---

## Investigation Checklist
- [x] Audit current README for factual errors
- [x] Check all links (internal and external)
- [x] Identify missing/outdated sections
- [x] Review inspiration (OpenClaw README)
- [x] Catalog actual project assets (screenshots, badges, CI)
- [x] Verify tech stack versions
- [x] Check template names and feature claims

## Findings

### Factual Errors in Current README

| Issue | Current | Correct | Severity |
|-------|---------|---------|----------|
| Tech stack: Next.js | 14 | 15.5.12 | **BLOCKER** |
| Tech stack: React | 18 | 19.0.0 | **BLOCKER** |
| Template names | "Classic, Modern, Minimal" | Classic, Modern, Visionary | **BLOCKER** |
| Export formats | "PDF, DOCX, JSON, YAML" | PDF, DOCX, JSON (no YAML) | **BLOCKER** |
| Export engine | "Gotenberg (HTML → PDF/DOCX)" | Gotenberg (HTML→PDF), godocx (CV data→DOCX) | Medium |
| Clone URL | `yourusername/hireme` | `Ellebam/hireme` | **BLOCKER** |
| Feature: Multi-language | Listed as feature | Not implemented | **BLOCKER** |
| Feature: Responsive | Listed as feature | Known broken (T-017) | Medium |
| Section count | Not mentioned | 8 section editors | Low |

### Dead Links

| Link | Target | Status |
|------|--------|--------|
| `docs/ARCHITECTURE.md` | Architecture docs | **MISSING** — file does not exist |
| `docs/API.md` | API reference | **MISSING** — file does not exist |
| `docs/DEPLOYMENT.md` | Deployment guide | **MISSING** — file does not exist |
| `docs/DEVELOPMENT.md` | Development guide | EXISTS |

### Missing Content (vs. what the project actually has)

1. **No app screenshots** — `data/st.png` is an unrelated image (red chair). No editor/dashboard screenshots exist.
   - **UPDATE:** User has added `docs/hireme1.png` (penguin mascot logo) and `docs/hireme2.png` (penguin hero image). Use hireme1 as centered logo, hireme2 as hero image.
2. **No CI badge** — GitHub Actions CI exists (`.github/workflows/ci.yml`) but no badge in README
3. **Test count not mentioned** — 170 passing Vitest tests, Go tests across all layers
4. **Key features not highlighted:**
   - Drag & drop section reordering (dnd-kit)
   - Live A4 preview with zoom (50–200%)
   - Undo/redo (50-step history)
   - Auto-save with debounce
   - Keyboard shortcuts
   - 8 section types (Personal, Summary, Experience, Education, Skills, Languages, Certifications, Projects)
5. **No mention of Zustand state management or schema-driven architecture**
6. **GitHub repo has no description or topics set**

### Design Issues (vs. Inspiration)

The OpenClaw README uses:
- Centered header with logo + tagline
- Badge row (CI, release, license, community)
- Progressive disclosure: quick-start → features → architecture → contributing
- Highlights section linking to deeper docs
- Clear multi-audience targeting (users, developers, operators)

Current HireMe README:
- Plain text header, no visual branding
- Only one badge (MIT license)
- Feature list is short and partially wrong
- References non-existent doc pages
- No screenshots or visual demonstration

## Gaps Found

### BLOCKER — Must fix
1. **6 factual errors** — wrong versions, wrong template names, wrong export formats, claims unimplemented features (multi-language, YAML export)
2. **3 dead doc links** — `docs/ARCHITECTURE.md`, `docs/API.md`, `docs/DEPLOYMENT.md` don't exist
3. **Wrong clone URL** — `yourusername` should be `Ellebam`

### Medium — Should fix
1. **No screenshots** — Need to capture editor and dashboard screenshots (or defer as separate task)
2. **No CI badge** — CI exists, badge should be added
3. **Missing feature highlights** — Many shipped features not mentioned
4. **Export engine description wrong** — DOCX uses godocx, not Gotenberg

### Low — Nice to fix
1. **No GitHub repo description/topics** — `gh repo edit` can set these
2. **No CONTRIBUTING.md** — README mentions contributing but has no guide
3. **Could add star history or project status indicators**

## Architect Handoff

### Decisions Needed

**QD-1: Screenshot strategy**
- **Option A:** Capture screenshots now and include in README (adds image files to repo)
- **Option B:** Skip screenshots for now, add placeholder text; create a separate task for screenshots later
- **Recommendation:** Option B — T-023 is XS, screenshots are a rabbit hole (need running app, good sample data, right viewport). Add a note "Screenshots coming soon" or just omit.

**QD-2: Dead doc links**
- **Option A:** Remove links entirely — honest about current state
- **Option B:** Keep links as "TODO" with placeholders
- **Option C:** Create stub files with "Coming soon"
- **Recommendation:** Option A — remove dead links. Link to `CONTEXT.md` for architecture info since it's the actual source of truth.

**QD-3: How much design polish**
- **Option A:** Just fix facts + structure (truly XS)
- **Option B:** Full visual redesign inspired by OpenClaw (centered header, badge row, progressive disclosure, ASCII diagrams)
- **Recommendation:** Option B — this is the whole point of the task. The README is the project's front door.

### Files in Scope

| File | Change | Gap |
|------|--------|-----|
| `README.md` | Full rewrite — fix facts, restructure, improve design | All blockers + medium |
| (optional) GitHub repo metadata | Set description + topics via `gh repo edit` | Low |

### Tests to Add/Update
- None — documentation only task

### Constraints
- XS task — keep scope tight, no new files beyond README
- No screenshots needed (defer to future task)
- Must not break existing LICENSE badge link

### Recommended Next Agent
**@engineer** — This is a straightforward rewrite. All facts are documented above; no architectural decisions beyond the three QDs listed. Engineer can execute directly after decisions are made.

## Test Plan
- Visual review: read the rendered README on GitHub after push
- Link check: verify all internal links resolve
- Factual check: cross-reference tech stack, features, commands against CONTEXT.md
