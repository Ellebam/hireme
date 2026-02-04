# Frontend MVP Architecture Plan

## Executive Summary

This document outlines the frontend MVP implementation plan for HireMe CV Editor, designed around:
- **Stack:** Next.js 14, React 18, Tailwind CSS, shadcn/ui, Zustand, dnd-kit, next-intl
- **Testing:** Vitest with React Testing Library
- **Backend:** Go API on localhost:8080 (functional)

---

## 1. Page Structure

### Routes

```
web/src/app/
├── layout.tsx                    # Root layout
├── page.tsx                      # Landing page
├── login/page.tsx                # Login (simplified for MVP)
└── [locale]/
    ├── layout.tsx                # Locale provider
    ├── (app)/                    # Authenticated routes
    │   ├── layout.tsx            # App shell + auth guard
    │   ├── dashboard/page.tsx    # CV list, create new
    │   └── editor/[cvId]/page.tsx # Main CV editor
    └── (marketing)/              # Public pages (future)
```

### MVP Pages

| Page | Route | Purpose | Rendering |
|------|-------|---------|-----------|
| Landing | `/` | Marketing, CTA | SSR |
| Login | `/login` | Auth entry | CSR |
| Dashboard | `/[locale]/dashboard` | CV list | SSR + CSR |
| Editor | `/[locale]/editor/[cvId]` | CV editing | CSR |

---

## 2. Editor Layout Architecture

### Three-Column Layout

```
+------------------------------------------------------------------+
|                         Top Toolbar                               |
|  [Undo] [Redo] | [Template: Modern v] | [Export v] | [Auto-saved] |
+------------------------------------------------------------------+
|           |                              |                        |
|  Section  |      CV Preview              |    Properties          |
|  Palette  |      (Live Render)           |    Panel               |
|  (240px)  |      (Flex)                  |    (320px)             |
|           |                              |                        |
+------------------------------------------------------------------+
```

### Zone Breakdown

| Zone | Width | Content | Mobile |
|------|-------|---------|--------|
| Left Sidebar | 240px | Section palette, structure tree | Bottom sheet |
| Center | Flex | CV preview (A4 aspect) | Full width |
| Right Sidebar | 320px | Properties panel | Bottom sheet |
| Top Toolbar | 100% | Actions, save status | Sticky, compact |

### Responsive Breakpoints

- **Desktop (≥1280px):** Full three-column
- **Tablet (768-1279px):** Two-column (toggle right)
- **Mobile (<768px):** Single column + FAB

---

## 3. Component Hierarchy

### Core Components

```
EditorLayout
├── EditorToolbar
│   ├── UndoRedoButtons
│   ├── TemplateSelector
│   ├── ExportDropdown
│   └── SaveStatus
├── SectionPalette (left)
│   ├── SectionButton
│   └── CVStructureTree
│       └── SectionTreeItem (draggable)
├── CVPreview (center)
│   └── SectionRenderer
│       ├── PersonalSection
│       ├── SummarySection
│       ├── ExperienceSection
│       ├── EducationSection
│       ├── SkillsSection
│       └── LanguagesSection
└── PropertiesPanel (right)
    ├── SectionProperties
    └── ContentEditors
        ├── PersonalEditor
        ├── ExperienceEntryEditor
        └── ...
```

### shadcn/ui Components Needed

- `button`, `card`, `dialog`, `dropdown-menu`, `input`, `label`
- `select`, `tabs`, `textarea`, `toast`, `tooltip` (new)
- `popover` (new), `switch` (new), `slider` (new)

---

## 4. State Management

### Zustand Stores

#### EditorStore (`stores/editor-store.ts`)

```typescript
interface EditorState {
  cv: CV | null;
  cvContent: CVContent | null;
  selectedSectionId: string | null;
  selectedEntryId: string | null;
  isDirty: boolean;
  saveStatus: 'saved' | 'saving' | 'error' | 'unsaved';

  // History
  history: CVContent[];
  historyIndex: number;

  // Actions
  setCV: (cv: CV) => void;
  updateSection: (id: string, updates: Partial<CVSection>) => void;
  addSection: (type: SectionType) => void;
  deleteSection: (id: string) => void;
  reorderSections: (activeId: string, overId: string) => void;
  undo: () => void;
  redo: () => void;
}
```

#### UIStore (`stores/ui-store.ts`)

```typescript
interface UIState {
  leftSidebarOpen: boolean;
  rightSidebarOpen: boolean;
  exportModalOpen: boolean;
  previewScale: number;

  toggleLeftSidebar: () => void;
  toggleRightSidebar: () => void;
}
```

### Auto-Save Strategy

- Debounced save (2 second delay)
- Optimistic updates with rollback on error
- Visual save status indicator

---

## 5. Drag & Drop

### Where DnD is Needed

| Context | Items | Purpose |
|---------|-------|---------|
| Section reordering | CVSection[] | Change CV layout |
| Experience entries | ExperienceEntry[] | Reorder work history |
| Education entries | EducationEntry[] | Reorder education |
| Skill categories | SkillCategory[] | Reorder skill groups |

### dnd-kit Configuration

```typescript
useSensors(
  useSensor(PointerSensor, { activationConstraint: { distance: 8 } }),
  useSensor(KeyboardSensor, { coordinateGetter: sortableKeyboardCoordinates })
);
```

---

## 6. Interaction Patterns

### Hover & Tooltips

| Element | Hover | Tooltip |
|---------|-------|---------|
| Section in preview | Blue outline | "Click to edit" |
| Drag handle | Cursor: grab | "Drag to reorder" |
| Visibility toggle | Highlight | "Hide/Show in export" |

### Edit Methods

| Content | Method | Reason |
|---------|--------|--------|
| Section title | Inline | Quick edit |
| Experience entry | Modal | Many fields |
| Summary text | Panel inline | Single field |
| Dates | Popover picker | Better UX |

### Feedback

- **Save success:** Toast + green checkmark
- **Section added:** Toast + auto-select
- **Delete:** Toast with undo option
- **Drag:** Visual opacity + shadow

---

## 7. Testing Strategy

### Priorities

| Priority | Category | Coverage | Tests |
|----------|----------|----------|-------|
| P0 | Zustand stores | 90% | 25-30 |
| P0 | Critical flows | 80% | 15-20 |
| P1 | UI components | 70% | 40-50 |
| P1 | API client | 80% | 10-15 |
| P2 | Drag & drop | 60% | 10-15 |

### Test Setup

```typescript
// src/test/setup.ts
import '@testing-library/jest-dom';
import { vi } from 'vitest';

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
  useParams: () => ({}),
}));
```

---

## 8. Implementation Order

### Phase 1: Foundation (Days 1-2)
- [ ] TypeScript types for CV schema
- [ ] API client with typed methods
- [ ] Zustand editor store (basic)
- [ ] Zustand UI store
- [ ] Vitest setup + test utilities
- [ ] Store tests

### Phase 2: Layout & Navigation (Days 2-3)
- [ ] App layout shell
- [ ] Dashboard page (CV list)
- [ ] Editor page skeleton
- [ ] Three-column layout
- [ ] Responsive layout

### Phase 3: CV Preview (Days 3-4)
- [ ] Section renderer framework
- [ ] PersonalSection renderer
- [ ] SummarySection renderer
- [ ] ExperienceSection renderer
- [ ] EducationSection renderer
- [ ] SkillsSection renderer
- [ ] LanguagesSection renderer

### Phase 4: Section Editing (Days 4-6)
- [ ] Properties panel framework
- [ ] PersonalEditor
- [ ] SummaryEditor
- [ ] ExperienceEntryEditor (modal)
- [ ] EducationEntryEditor (modal)
- [ ] SkillsEditor
- [ ] Add/delete section

### Phase 5: Drag & Drop (Day 6-7)
- [ ] dnd-kit setup
- [ ] Section reordering
- [ ] Entry reordering

### Phase 6: Polish & UX (Days 7-8)
- [ ] Toolbar (undo/redo)
- [ ] Auto-save
- [ ] Tooltips
- [ ] Toast notifications
- [ ] Export dropdown

### Phase 7: Testing (Days 8-9)
- [ ] Component tests
- [ ] Integration tests
- [ ] Mobile testing

---

## 9. MVP Scope

### In Scope
- Single CV editing
- 6 section types: Personal, Summary, Experience, Education, Skills, Languages
- Drag & drop reordering
- Live preview
- Auto-save + undo/redo
- 3 templates
- Export (PDF/DOCX/JSON)
- Responsive design
- i18n (en/de)

### Out of Scope (Future)
- Multiple CVs
- AI assistance
- Job matching
- Collaboration
- OAuth (using bypass)
- Premium features

---

## 10. File Structure

```
web/src/
├── app/
│   ├── globals.css
│   ├── layout.tsx
│   ├── page.tsx
│   ├── login/page.tsx
│   └── [locale]/
│       ├── layout.tsx
│       └── (app)/
│           ├── layout.tsx
│           ├── dashboard/page.tsx
│           └── editor/[cvId]/page.tsx
├── components/
│   ├── ui/                      # shadcn
│   ├── layout/                  # AppShell, Header, Sidebar
│   ├── editor/                  # Editor components
│   │   ├── sections/            # Section renderers
│   │   ├── editors/             # Content editors
│   │   └── modals/              # Export, Template modals
│   └── shared/                  # Loading, Empty, Error
├── stores/
│   ├── editor-store.ts
│   └── ui-store.ts
├── lib/
│   ├── api/client.ts
│   └── dnd/config.ts
├── hooks/
├── types/
│   ├── cv.ts
│   └── api.ts
├── i18n/
└── test/
```

---

## Critical References

1. **CV Schema:** `/home/user/hireme/schemas/cv-schema.json`
2. **Web Conventions:** `/home/user/hireme/web/CLAUDE.md`
3. **API Handlers:** `/home/user/hireme/api/internal/handler/cv.go`
4. **Sample Data:** `/home/user/hireme/api/db/seed/dev_seed.sql`
