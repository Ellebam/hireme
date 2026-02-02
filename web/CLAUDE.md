# HireMe Web — Frontend AI Context

## Overview

Next.js 14 application using App Router with TypeScript, Tailwind CSS, and shadcn/ui.

**Port:** 3000 (default Next.js dev server)

## Directory Structure

```
web/src/
├── app/
│   ├── [locale]/              # i18n routing (en, de)
│   │   ├── layout.tsx         # Root layout with providers
│   │   ├── page.tsx           # Landing page (SSR)
│   │   ├── (marketing)/       # Marketing pages (SSR)
│   │   │   ├── pricing/
│   │   │   └── about/
│   │   └── (app)/             # Authenticated app (mostly CSR)
│   │       ├── layout.tsx     # Auth guard
│   │       ├── dashboard/
│   │       └── editor/
│   ├── api/                   # API routes (auth callbacks)
│   ├── globals.css            # Tailwind base styles
│   └── providers.tsx          # Client providers wrapper
├── components/
│   ├── ui/                    # shadcn/ui components
│   ├── layout/                # Header, Footer, Sidebar
│   ├── editor/                # CV Editor components
│   │   ├── CVEditor.tsx       # Main editor container
│   │   ├── SectionList.tsx    # DnD list of sections
│   │   ├── LivePreview.tsx    # Real-time CV preview
│   │   └── sections/          # Section-specific editors
│   └── shared/                # Reusable components
├── lib/
│   ├── api/                   # API client and typed fetchers
│   ├── utils.ts               # General utilities
│   └── cn.ts                  # Tailwind class merge utility
├── stores/                    # Zustand state stores
├── hooks/                     # Custom React hooks
├── types/                     # TypeScript type definitions
└── i18n/                      # Internationalization
```

## Code Conventions

### Component Pattern
```tsx
// components/editor/SectionItem.tsx
"use client"; // Only if using hooks/interactivity

import { useState } from "react";
import { useSortable } from "@dnd-kit/sortable";
import { CSS } from "@dnd-kit/utilities";
import { GripVertical, Trash2 } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/cn";
import type { CVSection } from "@/types/cv";

interface SectionItemProps {
  section: CVSection;
  onDelete: (id: string) => void;
  onEdit: (section: CVSection) => void;
}

export function SectionItem({ section, onDelete, onEdit }: SectionItemProps) {
  const { attributes, listeners, setNodeRef, transform, transition, isDragging } = 
    useSortable({ id: section.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        "flex items-center gap-2 p-4 bg-white rounded-lg border",
        isDragging && "opacity-50"
      )}
    >
      <button {...attributes} {...listeners} className="cursor-grab">
        <GripVertical className="h-5 w-5 text-gray-400" />
      </button>
      
      <div className="flex-1">
        <h3 className="font-medium">{section.type}</h3>
      </div>
      
      <Button variant="ghost" size="sm" onClick={() => onEdit(section)}>
        Edit
      </Button>
      <Button variant="ghost" size="sm" onClick={() => onDelete(section.id)}>
        <Trash2 className="h-4 w-4" />
      </Button>
    </div>
  );
}
```

### Zustand Store Pattern
```tsx
// stores/editor-store.ts
import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";
import type { CV, CVSection } from "@/types/cv";

interface EditorState {
  cv: CV | null;
  selectedSectionId: string | null;
  isDirty: boolean;
  
  // Actions
  setCV: (cv: CV) => void;
  updateSection: (sectionId: string, content: Partial<CVSection>) => void;
  reorderSections: (activeId: string, overId: string) => void;
  addSection: (type: CVSection["type"]) => void;
  deleteSection: (id: string) => void;
  selectSection: (id: string | null) => void;
  markClean: () => void;
}

export const useEditorStore = create<EditorState>()(
  devtools(
    persist(
      (set, get) => ({
        cv: null,
        selectedSectionId: null,
        isDirty: false,

        setCV: (cv) => set({ cv, isDirty: false }),
        
        updateSection: (sectionId, content) => set((state) => {
          if (!state.cv) return state;
          return {
            cv: {
              ...state.cv,
              sections: state.cv.sections.map((s) =>
                s.id === sectionId ? { ...s, ...content } : s
              ),
            },
            isDirty: true,
          };
        }),
        
        // ... more actions
      }),
      { name: "hireme-editor" }
    )
  )
);
```

### API Client Pattern
```tsx
// lib/api/client.ts
const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const response = await fetch(url, {
      ...options,
      headers: {
        "Content-Type": "application/json",
        ...options?.headers,
      },
      credentials: "include", // For cookies/auth
    });

    if (!response.ok) {
      const error = await response.json().catch(() => ({}));
      throw new ApiError(response.status, error.message || "Request failed");
    }

    return response.json();
  }

  get<T>(endpoint: string) {
    return this.request<T>(endpoint, { method: "GET" });
  }

  post<T>(endpoint: string, data: unknown) {
    return this.request<T>(endpoint, {
      method: "POST",
      body: JSON.stringify(data),
    });
  }
  
  // ... put, delete, etc.
}

export const api = new ApiClient(API_BASE);
```

## Routing

### Locale-Based Routing
All pages are under `[locale]/` for i18n support:
- `/en/` → English
- `/de/` → German

### Route Groups
- `(marketing)/` — Public pages (SSR for SEO)
- `(app)/` — Authenticated pages (auth guard in layout)

### Auth Guard Pattern
```tsx
// app/[locale]/(app)/layout.tsx
import { redirect } from "next/navigation";
import { getSession } from "@/lib/auth";

export default async function AppLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const session = await getSession();
  
  if (!session) {
    redirect("/login");
  }

  return (
    <div className="min-h-screen flex">
      <Sidebar />
      <main className="flex-1">{children}</main>
    </div>
  );
}
```

## State Management

### When to Use What
| State Type | Solution |
|------------|----------|
| CV editor state | Zustand `useEditorStore` |
| UI state (modals, sidebars) | Zustand `useUIStore` |
| Server data (fetching) | React Query or SWR |
| Form state | React Hook Form |
| URL state | Next.js `useSearchParams` |

## Styling

### Tailwind + shadcn/ui
```tsx
// Use cn() for conditional classes
import { cn } from "@/lib/cn";

<div className={cn(
  "p-4 rounded-lg",
  isActive && "bg-blue-50 border-blue-200",
  isDisabled && "opacity-50 cursor-not-allowed"
)} />
```

### Component Installation
```bash
# Add shadcn/ui components as needed
npx shadcn-ui@latest add button
npx shadcn-ui@latest add card
npx shadcn-ui@latest add dialog
```

## i18n

Using `next-intl`:

```tsx
// Access translations in components
import { useTranslations } from "next-intl";

export function Header() {
  const t = useTranslations("Header");
  return <h1>{t("title")}</h1>;
}
```

Translation files in `src/i18n/`:
- `en.json` — English
- `de.json` — German

## Testing

```bash
task web:test        # Run Vitest
task web:test:watch  # Watch mode
```

### Test Pattern
```tsx
// components/editor/SectionItem.test.tsx
import { render, screen, fireEvent } from "@testing-library/react";
import { SectionItem } from "./SectionItem";

describe("SectionItem", () => {
  const mockSection = {
    id: "1",
    type: "experience" as const,
    content: {},
    order: 0,
  };

  it("renders section type", () => {
    render(
      <SectionItem 
        section={mockSection} 
        onDelete={vi.fn()} 
        onEdit={vi.fn()} 
      />
    );
    expect(screen.getByText("experience")).toBeInTheDocument();
  });

  it("calls onDelete when delete button clicked", () => {
    const onDelete = vi.fn();
    render(
      <SectionItem 
        section={mockSection} 
        onDelete={onDelete} 
        onEdit={vi.fn()} 
      />
    );
    fireEvent.click(screen.getByRole("button", { name: /delete/i }));
    expect(onDelete).toHaveBeenCalledWith("1");
  });
});
```

## Common Tasks

### Adding a New Page
1. Create `app/[locale]/(group)/pagename/page.tsx`
2. Add translations to `i18n/*.json`
3. Add to navigation if needed

### Adding a New Section Type
1. Create `components/editor/sections/NewSection.tsx`
2. Add type to `types/cv.ts`
3. Update `SectionList.tsx` to handle new type
4. Add default content in editor store

### Adding a UI Component
1. Check if shadcn/ui has it: `npx shadcn-ui@latest add <name>`
2. If custom, create in `components/ui/` or `components/shared/`

## Environment Variables

```env
# .env.local
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_APP_URL=http://localhost:3000

# Auth (when enabled)
NEXTAUTH_URL=http://localhost:3000
NEXTAUTH_SECRET=dev-secret-change-in-prod
GOOGLE_CLIENT_ID=...
GOOGLE_CLIENT_SECRET=...
```

## Performance Notes

- Use `"use client"` only when necessary (hooks, events, browser APIs)
- Prefer Server Components for static content
- Use `loading.tsx` for suspense boundaries
- Lazy load heavy editor components
