# HireMe Design System — Editorial Craft

## 1. Design Philosophy

The Editorial Craft design system draws from newspaper and magazine design traditions — where typography leads, whitespace breathes, and every element earns its place on the page.

**Core Principles:**

1. **Sharp, not soft.** No border-radius by default. Elements have crisp edges. Rounded corners are reserved for specific exceptions (template cards, pill tags).
2. **Ink on paper.** The palette is built around warm, natural tones — cream backgrounds, ink-dark text, vermillion accents. Every surface feels like quality stationery.
3. **Typography is hierarchy.** Three font families serve distinct roles: serif for headings and titles (authority), sans-serif for body text and UI (clarity), monospace for metadata and labels (precision).
4. **Offset, not float.** Shadows are solid and directional (3px 3px 0), not blurred. Elements feel stamped onto the page, not hovering above it.
5. **Dashed lines tell stories.** Dashed borders are a deliberate design element — they separate content like column rules in a newspaper, creating rhythm without heaviness.

---

## 2. Color Tokens

### 2.1 Semantic Tokens (shadcn-compatible)

These tokens replace the existing CSS variables in `globals.css`. They use HSL format for shadcn/ui compatibility.

| Token | Light Mode | Dark Mode | Description |
|-------|-----------|-----------|-------------|
| `--background` | `42.0 41.7% 95.3%` | `30.0 13.0% 9.0%` | Page background (cream / ink) |
| `--foreground` | `36.9 11.9% 21.4%` | `41.2 34.8% 91.0%` | Body text (text / cream-deep) |
| `--card` | `45.0 100.0% 98.4%` | `25.7 8.4% 16.3%` | Card/panel surfaces (paper / charcoal) |
| `--card-foreground` | `30.0 13.0% 9.0%` | `41.2 34.8% 91.0%` | Card text (ink / cream-deep) |
| `--popover` | `45.0 100.0% 98.4%` | `25.7 8.4% 16.3%` | Popover/dropdown bg (paper / charcoal) |
| `--popover-foreground` | `30.0 13.0% 9.0%` | `41.2 34.8% 91.0%` | Popover text (ink / cream-deep) |
| `--primary` | `30.0 13.0% 9.0%` | `42.0 41.7% 95.3%` | Primary buttons, heavy borders (ink / cream) |
| `--primary-foreground` | `45.0 100.0% 98.4%` | `30.0 13.0% 9.0%` | Text on primary (paper / ink) |
| `--secondary` | `41.2 34.8% 91.0%` | `30.0 10.0% 13.0%` | Secondary surfaces (cream-deep / dark warm) |
| `--secondary-foreground` | `36.9 11.9% 21.4%` | `39.1 28.4% 84.1%` | Text on secondary (text / cream-dark) |
| `--muted` | `41.2 34.8% 91.0%` | `30.0 10.0% 13.0%` | Muted backgrounds (cream-deep / dark warm) |
| `--muted-foreground` | `33.0 10.3% 62.0%` | `33.0 10.3% 62.0%` | Muted text — same both modes (text-muted) |
| `--accent` | `5.6 63.4% 46.1%` | `5.6 63.4% 46.1%` | Accent color — same both modes (vermillion) |
| `--accent-foreground` | `45.0 100.0% 98.4%` | `42.0 41.7% 95.3%` | Text on accent (paper / cream) |
| `--destructive` | `5.6 63.4% 46.1%` | `5.6 70.0% 52.0%` | Destructive actions (vermillion / brighter vermillion) |
| `--destructive-foreground` | `45.0 100.0% 98.4%` | `42.0 41.7% 95.3%` | Text on destructive (paper / cream) |
| `--border` | `38.2 22.0% 80.4%` | `30.0 8.0% 22.0%` | Default border color |
| `--input` | `38.2 22.0% 80.4%` | `30.0 8.0% 22.0%` | Input border color |
| `--ring` | `30.0 13.0% 9.0%` | `42.0 41.7% 95.3%` | Focus ring (ink / cream) |
| `--radius` | `0rem` | `0rem` | Default border radius — sharp corners |

### 2.2 Custom Tokens

These extend beyond shadcn's standard set to capture the prototype's richer vocabulary.

| Token | Light Mode | Dark Mode | Description |
|-------|-----------|-----------|-------------|
| `--cream-dark` | `39.1 28.4% 84.1%` | `25.0 10.0% 6.0%` | Offset shadow color |
| `--ink` | `30.0 13.0% 9.0%` | `42.0 41.7% 95.3%` | Heavy borders, emphasis (inverts with mode) |
| `--sienna` | `26.8 38.1% 43.7%` | `26.8 38.1% 53.0%` | Warm accent — italic text, job titles, avatar |
| `--vermillion-pale` | `5.6 63.4% 46.1% / 0.06` | `5.6 63.4% 46.1% / 0.12` | Subtle hover background (doubled alpha in dark) |
| `--text-secondary` | `32.5 11.5% 40.8%` | `33.0 12.0% 65.0%` | Secondary text color |

### 2.3 Copy-Pasteable CSS

```css
@layer base {
  :root {
    --background: 42.0 41.7% 95.3%;
    --foreground: 36.9 11.9% 21.4%;
    --card: 45.0 100.0% 98.4%;
    --card-foreground: 30.0 13.0% 9.0%;
    --popover: 45.0 100.0% 98.4%;
    --popover-foreground: 30.0 13.0% 9.0%;
    --primary: 30.0 13.0% 9.0%;
    --primary-foreground: 45.0 100.0% 98.4%;
    --secondary: 41.2 34.8% 91.0%;
    --secondary-foreground: 36.9 11.9% 21.4%;
    --muted: 41.2 34.8% 91.0%;
    --muted-foreground: 33.0 10.3% 62.0%;
    --accent: 5.6 63.4% 46.1%;
    --accent-foreground: 45.0 100.0% 98.4%;
    --destructive: 5.6 63.4% 46.1%;
    --destructive-foreground: 45.0 100.0% 98.4%;
    --border: 38.2 22.0% 80.4%;
    --input: 38.2 22.0% 80.4%;
    --ring: 30.0 13.0% 9.0%;
    --radius: 0rem;

    /* Custom tokens */
    --cream-dark: 39.1 28.4% 84.1%;
    --ink: 30.0 13.0% 9.0%;
    --sienna: 26.8 38.1% 43.7%;
    --vermillion-pale: 5.6 63.4% 46.1% / 0.06;
    --text-secondary: 32.5 11.5% 40.8%;
  }

  .dark {
    --background: 30.0 13.0% 9.0%;
    --foreground: 41.2 34.8% 91.0%;
    --card: 25.7 8.4% 16.3%;
    --card-foreground: 41.2 34.8% 91.0%;
    --popover: 25.7 8.4% 16.3%;
    --popover-foreground: 41.2 34.8% 91.0%;
    --primary: 42.0 41.7% 95.3%;
    --primary-foreground: 30.0 13.0% 9.0%;
    --secondary: 30.0 10.0% 13.0%;
    --secondary-foreground: 39.1 28.4% 84.1%;
    --muted: 30.0 10.0% 13.0%;
    --muted-foreground: 33.0 10.3% 62.0%;
    --accent: 5.6 63.4% 46.1%;
    --accent-foreground: 42.0 41.7% 95.3%;
    --destructive: 5.6 70.0% 52.0%;
    --destructive-foreground: 42.0 41.7% 95.3%;
    --border: 30.0 8.0% 22.0%;
    --input: 30.0 8.0% 22.0%;
    --ring: 42.0 41.7% 95.3%;
    --radius: 0rem;

    /* Custom tokens */
    --cream-dark: 25.0 10.0% 6.0%;
    --ink: 42.0 41.7% 95.3%;
    --sienna: 26.8 38.1% 53.0%;
    --vermillion-pale: 5.6 63.4% 46.1% / 0.12;
    --text-secondary: 33.0 12.0% 65.0%;
  }
}
```

### 2.4 Named Color Reference

Quick-reference mapping from prototype names to hex values (light mode only — for visual comparison with the HTML prototype):

| Name | Hex | Swatch Description |
|------|-----|--------------------|
| Cream | `#f8f5ee` | Warm off-white — page background |
| Cream Deep | `#f0ebe0` | Slightly deeper cream — editor center, badges |
| Cream Dark | `#e2dacb` | Warm tan — offset shadow color |
| Paper | `#fffdf7` | Near-white with warm tint — card surfaces |
| Ink | `#1a1714` | Near-black with warm undertone — headings, borders |
| Charcoal | `#2d2926` | Dark warm gray — dark mode card surface |
| Vermillion | `#c0392b` | Rich red — accent, active states, section titles |
| Sienna | `#9a6b45` | Warm brown — italic text, avatars, job titles |
| Text | `#3d3830` | Dark warm gray — body text |
| Text Secondary | `#74695c` | Medium warm gray — descriptions, nav links |
| Text Muted | `#a89f94` | Light warm gray — metadata, dates, labels |
| Border | `#d8d0c2` | Light warm tan — standard borders |
| Border Heavy | `#c4b9a8` | Medium warm tan — heavy borders (reserved) |

---

## 3. Typography

### 3.1 Font Stack

| Role | Font Family | CSS Variable | Tailwind Class | `next/font/google` Import |
|------|------------|-------------|----------------|--------------------------|
| Serif | Newsreader | `--font-serif` | `font-serif` | `Newsreader({ subsets: ['latin'], variable: '--font-serif', display: 'swap' })` |
| Sans | Source Sans 3 | `--font-sans` | `font-sans` | `Source_Sans_3({ subsets: ['latin'], variable: '--font-sans', display: 'swap' })` |
| Mono | JetBrains Mono | `--font-mono` | `font-mono` | `JetBrains_Mono({ subsets: ['latin'], variable: '--font-mono', display: 'swap' })` |

All three are **variable fonts** — `next/font/google` loads the full weight range automatically. No `weight` array is needed in the import config.

**Tailwind config extension:**
```ts
fontFamily: {
  serif: ['var(--font-serif)', 'Georgia', 'serif'],
  sans: ['var(--font-sans)', 'system-ui', 'sans-serif'],
  mono: ['var(--font-mono)', 'Menlo', 'monospace'],
},
```

**Layout.tsx body classes:**
```tsx
<body className={`${newsreader.variable} ${sourceSans3.variable} ${jetbrainsMono.variable} font-sans antialiased`}>
```

### 3.2 Type Scale

| Element | Font | Size | Weight | Line-Height | Letter-Spacing | Tailwind Recipe |
|---------|------|------|--------|-------------|----------------|-----------------|
| Hero title | Serif | 56px (3.5rem) | 700 | 1.08 | -0.03em | `font-serif text-[3.5rem] font-bold leading-[1.08] tracking-[-0.03em]` |
| Hero title emphasis | Serif | 56px | 400 italic | 1.08 | -0.03em | `font-serif text-[3.5rem] font-normal italic leading-[1.08] tracking-[-0.03em] text-sienna` |
| Page section title | Serif | 26px (1.625rem) | 600 | 1.2 | -0.02em | `font-serif text-[1.625rem] font-semibold tracking-[-0.02em]` |
| Template page title | Serif | 44px (2.75rem) | 700 | 1.1 | -0.03em | `font-serif text-[2.75rem] font-bold leading-[1.1] tracking-[-0.03em]` |
| Modal title | Serif | 26px | 700 | 1.2 | -0.02em | `font-serif text-[1.625rem] font-bold tracking-[-0.02em]` |
| Nav logo text | Serif | 24px (1.5rem) | 700 | 1 | -0.03em | `font-serif text-2xl font-bold tracking-[-0.03em]` |
| Card/row title | Serif | 19px (1.1875rem) | 600 | 1.3 | -0.01em | `font-serif text-[1.1875rem] font-semibold tracking-[-0.01em]` |
| Panel title | Serif | 16px (1rem) | 600 | 1.3 | 0 | `font-serif text-base font-semibold` |
| CV name | Serif | 34px (2.125rem) | 700 | 1 | -0.03em | `font-serif text-[2.125rem] font-bold leading-none tracking-[-0.03em]` |
| CV job title | Serif | 16px | 400 italic | 1.4 | 0 | `font-serif text-base font-normal italic text-sienna` |
| CV entry title | Serif | 15px (0.9375rem) | 600 | 1.4 | 0 | `font-serif text-[0.9375rem] font-semibold` |
| Hero body | Sans | 17px (1.0625rem) | 400 | 1.7 | 0 | `text-[1.0625rem] leading-[1.7] text-text-secondary` |
| Body text | Sans | 14px (0.875rem) | 400 | 1.75 | 0 | `text-sm leading-[1.75]` |
| Body text small | Sans | 13px (0.8125rem) | 400 | 1.5 | 0 | `text-[0.8125rem] leading-normal` |
| Button text | Sans | 13px | 600 | 1 | 1px (0.0625em) | `text-[0.8125rem] font-semibold uppercase tracking-[0.0625em]` |
| CV section title | Sans | 11px (0.6875rem) | 700 | 1 | 3px (0.1875em) | `text-[0.6875rem] font-bold uppercase tracking-[0.1875em] text-accent` |
| Property label | Mono | 11px | 500 | 1 | 0.5px | `font-mono text-[0.6875rem] font-medium uppercase tracking-[0.03em] text-muted-foreground` |
| Rule label | Mono | 11px | 500 | 1 | 3px | `font-mono text-[0.6875rem] font-medium uppercase tracking-[0.1875em] text-accent` |
| Row number | Mono | 13px | 400 | 1 | 0 | `font-mono text-[0.8125rem] text-muted-foreground` |
| Badge/count | Mono | 12px (0.75rem) | 400 | 1 | 0 | `font-mono text-xs text-muted-foreground bg-secondary px-2 py-0.5` |
| Zoom display | Mono | 12px | 400 | 1 | 0 | `font-mono text-xs text-text-secondary` |
| Save indicator | Mono | 11px | 400 | 1 | 0.5px | `font-mono text-[0.6875rem] uppercase tracking-[0.03em] text-muted-foreground` |
| Template tag | Mono | 10px (0.625rem) | 500 | 1 | 1px | `font-mono text-[0.625rem] font-medium uppercase tracking-[0.0625em]` |

### 3.3 Usage Rules

| Context | Font | Why |
|---------|------|-----|
| Page titles, section headings, modal titles | Serif (Newsreader) | Authority, editorial weight |
| CV names, entry titles, card titles | Serif (Newsreader) | Content emphasis, document feel |
| Logo wordmark | Serif (Newsreader) | Brand identity |
| Body text, descriptions, paragraphs | Sans (Source Sans 3) | Readability at small sizes |
| Buttons, form labels (non-mono) | Sans (Source Sans 3) | UI clarity |
| CV section titles (uppercase) | Sans (Source Sans 3) | Clear hierarchy without competing with serif |
| Metadata: dates, counts, row numbers | Mono (JetBrains Mono) | Tabular alignment, technical precision |
| Labels: property labels, rules, tags | Mono (JetBrains Mono) | Systematic, structured feel |
| Save/status indicators, zoom display | Mono (JetBrains Mono) | Dashboard instrument aesthetic |
| Code or technical content | Mono (JetBrains Mono) | Expected convention |

**Rule of thumb:** If it's a *heading* or *name*, use serif. If it's *readable content* or *UI control*, use sans. If it's *metadata*, *label*, or *number*, use mono.

---

## 4. Spacing

### 4.1 Component Spacing Reference

| Context | Value | Tailwind | Notes |
|---------|-------|----------|-------|
| Page horizontal padding | 40px | `px-10` | Consistent across dashboard sections |
| Hero top padding | 72px | `pt-[72px]` | Below nav |
| Nav height | 60px | `h-[60px]` | Fixed top bar |
| Nav horizontal padding | 36px | `px-9` | Slightly less than page padding |
| Editor toolbar height | 50px | `h-[50px]` | Below nav |
| Editor toolbar padding | 0 20px | `px-5` | Horizontal only |
| Left sidebar width | 248px | `w-[248px]` | Section palette |
| Right sidebar width | 320px | `w-[320px]` | Properties panel |
| CV paper padding | 52px 48px | `px-12 py-[52px]` | Inner content area |
| Panel head padding | 18px 20px | `px-5 py-[18px]` | Section/properties panel headers |
| Section item padding | 10px 12px | `px-3 py-2.5` | Items in section list |
| Section item gap | 14px | `gap-3.5` | Between icon/number and label |
| Property section padding | 16px 20px | `px-5 py-4` | Properties panel content area |
| Input padding | 9px 12px | `px-3 py-[9px]` | Text inputs and textareas |
| Button padding (default) | 10px 20px | `px-5 py-2.5` | Standard button size |
| Button padding (small) | 6px 14px | `px-3.5 py-1.5` | Toolbar export button |
| Button gap (icon) | 8px | `gap-2` | Between icon and text |
| Card info padding | 20px | `p-5` | Template card info section |
| Modal padding | 36px | `p-9` | Export modal |
| Hero actions gap | 12px | `gap-3` | Between action buttons |
| Nav links gap | 0 | `gap-0` | Links touch, padding creates spacing |
| Nav link padding | 18px 20px | `px-5 py-[18px]` | Click area for nav items |
| Template grid gap | 16px | `gap-4` | Between template cards |
| CV block margin | 22px bottom | `mb-[22px]` | Between CV sections |
| Format option padding | 16px 0 | `py-4` | Export format rows |

### 4.2 Base Unit

The prototype doesn't use a strict 4px/8px grid — spacing values are custom per-context. Use the exact values above rather than trying to snap to a grid.

---

## 5. Shadows

### 5.1 Offset Shadow System

All shadows use the **cream-dark** color for a solid offset effect — no blur, no spread. This is the defining visual characteristic of the design.

| Name | CSS Value | Tailwind Class | Usage |
|------|-----------|---------------|-------|
| `shadow-offset-sm` | `3px 3px 0 hsl(var(--cream-dark))` | `shadow-offset-sm` | Buttons (default), focused inputs |
| `shadow-offset-sm-hover` | `4px 4px 0 hsl(var(--cream-dark))` | `shadow-offset-sm-hover` | Button hover state (with `-translate-x-px -translate-y-px`) |
| `shadow-offset-md` | `4px 4px 0 hsl(var(--cream-dark))` | `shadow-offset-md` | CV row hover, general emphasis |
| `shadow-offset-lg` | `6px 6px 0 hsl(var(--cream-dark))` | `shadow-offset-lg` | CV paper preview |
| `shadow-offset-xl` | `8px 8px 0 hsl(var(--cream-dark))` | `shadow-offset-xl` | Modals |
| `shadow-card-hover` | `0 12px 32px hsla(var(--accent), 0.08)` | `shadow-card-hover` | Template cards only (exception: soft shadow) |

**Tailwind config extension:**
```ts
boxShadow: {
  'offset-sm': '3px 3px 0 hsl(var(--cream-dark))',
  'offset-sm-hover': '4px 4px 0 hsl(var(--cream-dark))',
  'offset-md': '4px 4px 0 hsl(var(--cream-dark))',
  'offset-lg': '6px 6px 0 hsl(var(--cream-dark))',
  'offset-xl': '8px 8px 0 hsl(var(--cream-dark))',
  'card-hover': '0 12px 32px hsl(var(--accent) / 0.08)',
},
```

### 5.2 Shadow Interaction Pattern

Buttons use a **lift** effect: on hover, they translate up-left by 1px while the shadow grows.

```
Default:  shadow-offset-sm
Hover:    -translate-x-px -translate-y-px shadow-offset-sm-hover
Active:   translate-x-0 translate-y-0 shadow-none  (pressed down)
```

---

## 6. Borders

### 6.1 Border Vocabulary

| Style | CSS | Tailwind Recipe | Usage |
|-------|-----|----------------|-------|
| Heavy | `2px solid hsl(var(--ink))` | `border-2 border-ink` | Nav bottom, panel heads, CV contact divider, page dividers |
| Standard | `2px solid hsl(var(--border))` | `border-2 border-border` | Editor panel borders, input borders, CV paper |
| Light dashed | `1px dashed hsl(var(--border))` | `border border-dashed border-border` | Row separators, editor toolbar bottom, CV section titles |
| Card | `1px solid hsl(var(--border))` | `border border-border` | Template cards |
| Accent | `2px solid hsl(var(--accent))` | `border-2 border-accent` | Logo stamp |
| Active indicator | `3px solid hsl(var(--accent))` | `border-l-[3px] border-accent` | Active section item (left border) |
| Active placeholder | `3px solid transparent` | `border-l-[3px] border-transparent` | Inactive section items (prevents layout shift) |

### 6.2 Border Radius

Default is **0** (sharp corners). Exceptions are explicit:

| Value | Tailwind | Usage |
|-------|----------|-------|
| `0` | `rounded-none` (or omit — default) | Buttons, inputs, textareas, modals, nav, CV paper, logo stamp |
| `8px` | `rounded-lg` | Template cards |
| `4px` | `rounded` | Template mini-preview thumbnails |
| `20px` | `rounded-full` (or `rounded-[20px]`) | Tags, badges (pill shape) |
| `50%` | `rounded-full` | Save indicator dot |

**Tailwind config:**
```ts
borderRadius: {
  lg: '8px',
  md: '4px',
  sm: '2px',
  // Default --radius is 0rem, so `rounded-lg` = 8px, not calc(var(--radius))
},
```

Note: Setting `--radius: 0rem` means the existing `calc(var(--radius) - 2px)` formulas all resolve to 0 or negative (clamped to 0). This effectively makes all shadcn components sharp-cornered by default. Template cards that need rounding should use explicit `rounded-lg`.

---

## 7. Animations & Transitions

### 7.1 Keyframe Animations

| Name | Keyframes | Duration / Easing | Usage |
|------|-----------|-------------------|-------|
| `page-in` | `opacity: 0 → 1` | `0.4s ease` | Page transitions |
| `slide-in` | `opacity: 0, translateX(-12px) → opacity: 1, translateX(0)` | `0.5s ease` | Hero elements (stagger with `animation-delay`) |
| `paper-drop` | `opacity: 0, translateY(12px) rotate(-0.3deg) → opacity: 1, translateY(0) rotate(0)` | `0.5s cubic-bezier(0.16, 1, 0.3, 1)` | CV paper appearing in editor |
| `fade-in` | `opacity: 0 → 1` | `0.2s ease` | Modal backdrop |
| `modal-drop` | `opacity: 0, translateY(8px) → opacity: 1, translateY(0)` | `0.3s ease` | Modal content |

**Tailwind config extension:**
```ts
keyframes: {
  'page-in': {
    from: { opacity: '0' },
    to: { opacity: '1' },
  },
  'slide-in': {
    from: { opacity: '0', transform: 'translateX(-12px)' },
    to: { opacity: '1', transform: 'translateX(0)' },
  },
  'paper-drop': {
    from: { opacity: '0', transform: 'translateY(12px) rotate(-0.3deg)' },
    to: { opacity: '1', transform: 'translateY(0) rotate(0)' },
  },
  'fade-in': {
    from: { opacity: '0' },
    to: { opacity: '1' },
  },
  'modal-drop': {
    from: { opacity: '0', transform: 'translateY(8px)' },
    to: { opacity: '1', transform: 'translateY(0)' },
  },
  // Keep existing accordion animations
  'accordion-down': {
    from: { height: '0' },
    to: { height: 'var(--radix-accordion-content-height)' },
  },
  'accordion-up': {
    from: { height: 'var(--radix-accordion-content-height)' },
    to: { height: '0' },
  },
},
animation: {
  'page-in': 'page-in 0.4s ease',
  'slide-in': 'slide-in 0.5s ease forwards',
  'paper-drop': 'paper-drop 0.5s cubic-bezier(0.16, 1, 0.3, 1)',
  'fade-in': 'fade-in 0.2s ease',
  'modal-drop': 'modal-drop 0.3s ease',
  'accordion-down': 'accordion-down 0.2s ease-out',
  'accordion-up': 'accordion-up 0.2s ease-out',
},
```

**Stagger pattern for hero elements:** Use `animation-delay` via arbitrary values:
```html
<div class="animate-slide-in opacity-0 [animation-delay:0.1s]">Rule</div>
<h1 class="animate-slide-in opacity-0 [animation-delay:0.2s]">Title</h1>
<p class="animate-slide-in opacity-0 [animation-delay:0.35s]">Body</p>
<div class="animate-slide-in opacity-0 [animation-delay:0.45s]">Actions</div>
```

Note: `opacity-0` is needed as initial state since `animation-fill-mode: forwards` (in the `forwards` keyword) keeps the final frame.

### 7.2 Transitions

| Element | Property | Duration | Tailwind Recipe |
|---------|----------|----------|----------------|
| Nav links | `color` | `200ms ease` | `transition-colors duration-200` |
| Buttons (all) | `all` | `200ms ease` | `transition-all duration-200` |
| CV rows | `all` | `200ms ease` | `transition-all duration-200` |
| Section items | `all` | `150ms ease` | `transition-all duration-150` |
| Editor bar buttons | `all` | `150ms ease` | `transition-all duration-150` |
| Inputs (focus) | `all` | `150ms ease` | `transition-all duration-150` |
| Template cards | `all` | `250ms ease` | `transition-all duration-250` |
| Drag handles | `opacity` | `150ms` | `transition-opacity duration-150` |

---

## 8. Component Recipes

### 8.1 Navigation

Fixed top bar with logo, nav links, and user actions.

```
Container:
  fixed top-0 left-0 right-0 z-50
  h-[60px] bg-card border-b-2 border-ink
  flex items-center px-9

Logo stamp:
  w-[34px] h-[34px] border-2 border-accent
  flex items-center justify-center
  font-serif font-bold text-[17px] text-accent
  -rotate-[2deg]

Logo text:
  font-serif text-2xl font-bold text-primary tracking-[-0.03em]

Nav link:
  px-5 py-[18px]
  text-xs font-semibold uppercase tracking-[0.125em]
  text-text-secondary
  transition-colors duration-200
  hover:text-primary
  relative

Nav link (active):
  text-primary
  after:absolute after:bottom-[-2px] after:left-5 after:right-5
  after:h-0.5 after:bg-accent

User avatar:
  w-8 h-8 border-2 border-border
  flex items-center justify-center
  font-serif font-semibold text-[13px] text-sienna
  cursor-pointer
```

### 8.2 Hero Section

Editorial-style page header with rule pattern.

```
Container:
  pt-[72px] px-10 max-w-[840px]

Rule pattern:
  flex items-center gap-3 mb-5
  animate-slide-in opacity-0 [animation-delay:0.1s]

Rule line:
  w-9 h-0.5 bg-accent

Rule text:
  font-mono text-[0.6875rem] font-medium uppercase tracking-[0.1875em] text-accent

Title:
  font-serif text-[3.5rem] font-bold leading-[1.08] tracking-[-0.03em] text-primary
  mb-5 animate-slide-in opacity-0 [animation-delay:0.2s]

Title emphasis (italic words):
  font-normal italic text-sienna

Body:
  text-[1.0625rem] text-text-secondary leading-[1.7]
  max-w-[500px] mb-9
  animate-slide-in opacity-0 [animation-delay:0.35s]

Actions:
  flex gap-3 mb-14
  animate-slide-in opacity-0 [animation-delay:0.45s]
```

### 8.3 Document List

Grid-based row layout for CV documents.

```
Section header:
  flex items-baseline justify-between py-6 pb-4

Section title:
  font-serif text-[1.625rem] font-semibold tracking-[-0.02em]

Section count badge:
  font-mono text-xs text-muted-foreground bg-secondary px-2 py-0.5

Row:
  grid grid-cols-[48px_1fr_120px_100px_40px] items-center gap-5
  py-5 border-b border-dashed border-border
  cursor-pointer transition-all duration-200
  hover:bg-[hsl(var(--vermillion-pale))] hover:pl-3 hover:-ml-3 hover:-mr-3 hover:pr-3
  hover:shadow-offset-md

Row number:
  font-mono text-[0.8125rem] text-muted-foreground

Row title:
  font-serif text-[1.1875rem] font-semibold text-primary tracking-[-0.01em]

Row description:
  text-[0.8125rem] text-text-secondary

Row template:
  text-[0.8125rem] font-medium text-text-secondary

Row date:
  font-mono text-xs text-muted-foreground

Row action (arrow):
  text-lg text-muted-foreground flex justify-end
  hover:text-accent

Create row:
  flex items-center gap-4 py-6
  cursor-pointer text-muted-foreground transition-all duration-200
  hover:text-accent

Create plus:
  font-serif text-2xl font-light

Create label:
  text-[0.8125rem] font-semibold uppercase tracking-[0.0625em]
```

### 8.4 Buttons

Four variants, all with sharp corners and offset shadows.

```
Base:
  inline-flex items-center gap-2
  px-5 py-2.5
  text-[0.8125rem] font-semibold uppercase tracking-[0.0625em]
  cursor-pointer transition-all duration-200

Primary (ink background):
  bg-primary text-primary-foreground
  shadow-offset-sm
  hover:-translate-x-px hover:-translate-y-px hover:shadow-offset-sm-hover
  active:translate-x-0 active:translate-y-0 active:shadow-none

Red (vermillion background):
  bg-accent text-accent-foreground
  shadow-offset-sm
  hover:-translate-x-px hover:-translate-y-px hover:shadow-offset-sm-hover
  active:translate-x-0 active:translate-y-0 active:shadow-none

Outline:
  bg-transparent text-primary border-2 border-primary
  hover:bg-primary hover:text-primary-foreground

Ghost:
  bg-transparent text-text-secondary px-3 py-2.5
  hover:text-primary
```

### 8.5 Editor Toolbar

Horizontal bar below nav with tool buttons and status.

```
Container:
  h-[50px] bg-card border-b border-dashed border-border
  flex items-center px-5 gap-1

Tool button:
  w-8 h-8 flex items-center justify-center
  border border-transparent bg-transparent
  text-text-secondary text-[15px]
  cursor-pointer transition-all duration-150
  hover:text-primary hover:bg-[hsl(var(--vermillion-pale))]
  hover:border-accent/10

Tool button (active):
  text-accent bg-[hsl(var(--vermillion-pale))]

Separator:
  w-px h-5 bg-border mx-1.5

Zoom display:
  font-mono text-xs text-text-secondary px-1.5

Save indicator:
  ml-auto flex items-center gap-1.5
  font-mono text-[0.6875rem] text-muted-foreground uppercase tracking-[0.03em]

Save indicator dot:
  w-1.5 h-1.5 rounded-full bg-[#3d7a3d]

Export button (small):
  px-3.5 py-1.5 text-[0.6875rem]
  (uses primary button variant at reduced size)
```

### 8.6 Section Panel (Left Sidebar)

Ordered list of CV sections with active state.

```
Container:
  w-[248px] bg-card border-r-2 border-border
  flex flex-col

Panel head:
  px-5 py-[18px] border-b-2 border-ink

Panel title:
  font-serif text-base font-semibold text-primary

Panel subtitle:
  font-mono text-[0.6875rem] text-muted-foreground mt-0.5

Section list:
  p-2 flex-1 overflow-y-auto

Section item:
  flex items-center gap-3.5
  px-3 py-2.5 cursor-pointer
  transition-all duration-150
  border-l-[3px] border-transparent
  mb-px
  hover:bg-black/[0.02]

Section item (active):
  border-l-accent bg-[hsl(var(--vermillion-pale))]

Section number:
  font-mono text-[0.6875rem] text-muted-foreground w-5

Section label:
  text-sm font-medium text-text-secondary

Section label (active):
  text-primary font-semibold

Drag handle:
  ml-auto opacity-0 text-[0.6875rem] text-muted-foreground font-mono
  transition-opacity duration-150
  group-hover:opacity-50
```

### 8.7 CV Preview

Paper metaphor with offset shadow.

```
Container (editor center):
  flex-1 bg-secondary
  flex items-start justify-center
  p-10 overflow-auto

CV paper:
  w-[580px] min-h-[820px]
  bg-card border-2 border-border
  shadow-offset-lg
  px-12 py-[52px]
  text-primary
  animate-paper-drop

CV name:
  font-serif text-[2.125rem] font-bold tracking-[-0.03em] leading-none mb-1

CV job title:
  font-serif text-base font-normal italic text-sienna mb-3.5

CV contact bar:
  flex gap-4
  font-mono text-[0.6875rem] text-text-secondary
  pb-[18px] mb-[22px]
  border-b-2 border-ink

CV section title:
  text-[0.6875rem] font-bold uppercase tracking-[0.1875em]
  text-accent
  mb-2.5 pb-1.5
  border-b border-dashed border-border

CV block:
  mb-[22px]

CV entry title:
  font-serif text-[0.9375rem] font-semibold text-primary

CV entry subtitle:
  text-[0.8125rem] text-text-secondary italic mb-1

CV entry body:
  text-[0.8125rem] text-foreground leading-[1.75]

CV skills (inline list):
  flex flex-wrap
  font-mono text-[0.6875rem] text-foreground
  (items separated by " · " using CSS ::after)
```

### 8.8 Properties Panel (Right Sidebar)

Form panel for editing section content.

```
Container:
  w-[320px] bg-card border-l-2 border-border
  flex flex-col

(Uses same panel head as section panel — 8.6)

Properties content:
  px-5 py-4 flex-1 overflow-y-auto

Property label:
  font-mono text-[0.6875rem] font-medium uppercase tracking-[0.03em]
  text-muted-foreground mb-1.5

Text input:
  w-full px-3 py-[9px]
  bg-secondary border-2 border-border
  text-foreground text-sm
  transition-all duration-150
  outline-none mb-3
  focus:border-primary focus:bg-card focus:shadow-offset-sm

Textarea:
  w-full px-3 py-[9px]
  bg-secondary border-2 border-border
  text-foreground text-sm
  resize-y min-h-[80px]
  outline-none
  transition-all duration-150
  focus:border-primary focus:bg-card focus:shadow-offset-sm
```

### 8.9 Template Cards

Grid of selectable CV templates with hover lift.

```
Page container:
  max-w-[960px] mx-auto px-10 py-[60px] pb-20

(Uses hero rule pattern + serif title from 8.2)

Grid:
  grid grid-cols-3 gap-4

Card:
  bg-card border border-border rounded-lg
  cursor-pointer transition-all duration-250 overflow-hidden
  hover:border-accent/20 hover:shadow-card-hover hover:-translate-y-1

Preview area:
  h-[200px] flex items-center justify-center bg-secondary

Mini paper thumbnail:
  w-[110px] h-[155px] bg-card border border-border rounded
  shadow-[0_6px_16px_rgba(0,0,0,0.06)]
  p-2.5 transition-transform duration-250
  group-hover:-translate-y-1.5 group-hover:-rotate-1

Card info:
  p-5

Card number:
  font-mono text-[0.6875rem] text-muted-foreground mb-1.5

Card name:
  font-serif text-[1.1875rem] font-semibold text-primary tracking-[-0.01em] mb-1

Card description:
  text-[0.8125rem] text-text-secondary leading-normal

Card tag:
  inline-block mt-2.5
  font-mono text-[0.625rem] font-medium uppercase tracking-[0.0625em]
  text-accent bg-[hsl(var(--vermillion-pale))]
  px-3 py-1 rounded-full
```

### 8.10 Export Modal

Centered dialog with format selection list.

```
Backdrop:
  fixed inset-0 z-[200]
  bg-primary/55
  flex items-center justify-center
  animate-fade-in

Modal:
  bg-card border-2 border-primary
  p-9 max-w-[440px] w-[90%]
  shadow-offset-xl
  animate-modal-drop

Modal title:
  font-serif text-[1.625rem] font-bold text-primary tracking-[-0.02em] mb-1.5

Modal description:
  text-sm text-text-secondary mb-7 leading-normal

Format options container:
  flex flex-col mb-7 border-t border-dashed border-border

Format option:
  flex items-center gap-4
  py-4 border-b border-dashed border-border
  cursor-pointer transition-all duration-150
  hover:bg-[hsl(var(--vermillion-pale))] hover:pl-2

Format option (selected):
  bg-[hsl(var(--vermillion-pale))] pl-2
  border-l-[3px] border-accent

Format icon:
  font-mono text-xs font-semibold w-11 text-center

Format icon colors:
  PDF: text-accent
  DOCX: text-[#2563eb]
  JSON: text-primary

Format name:
  font-semibold text-sm text-primary

Format subtitle:
  text-xs text-muted-foreground

Modal footer:
  flex gap-3 justify-end
```

### 8.11 Dividers

Two types used throughout the app.

```
Heavy divider:
  border-t-2 border-primary mx-10

Light dashed divider:
  border-t border-dashed border-border mx-10
```

### 8.12 Background Texture

Subtle SVG crosshatch pattern overlaid on the page.

```
Applied via CSS pseudo-element on body::before:
  fixed inset-0 pointer-events-none z-0
  background-image: url("data:image/svg+xml,...")
  (crosshatch SVG at 1.5% opacity, 60x60px tile)

Content wrapper needs:
  relative z-[1]
```

The SVG data URL from the prototype:
```
data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23000000' fill-opacity='0.015'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E
```

Note: In dark mode, consider inverting the fill color to white at the same low opacity, or using `filter: invert(1)` on the pseudo-element.

### 8.13 Scrollbar

Custom scrollbar to match the editorial aesthetic.

```css
::-webkit-scrollbar { width: 4px; }
::-webkit-scrollbar-track { background: hsl(var(--background)); }
::-webkit-scrollbar-thumb { background: hsl(var(--cream-dark)); }
```

---

## 9. Implementation Guide (for T-025)

### 9.1 globals.css

Replace the full content of `web/src/app/globals.css` with:

```css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --background: 42.0 41.7% 95.3%;
    --foreground: 36.9 11.9% 21.4%;
    --card: 45.0 100.0% 98.4%;
    --card-foreground: 30.0 13.0% 9.0%;
    --popover: 45.0 100.0% 98.4%;
    --popover-foreground: 30.0 13.0% 9.0%;
    --primary: 30.0 13.0% 9.0%;
    --primary-foreground: 45.0 100.0% 98.4%;
    --secondary: 41.2 34.8% 91.0%;
    --secondary-foreground: 36.9 11.9% 21.4%;
    --muted: 41.2 34.8% 91.0%;
    --muted-foreground: 33.0 10.3% 62.0%;
    --accent: 5.6 63.4% 46.1%;
    --accent-foreground: 45.0 100.0% 98.4%;
    --destructive: 5.6 63.4% 46.1%;
    --destructive-foreground: 45.0 100.0% 98.4%;
    --border: 38.2 22.0% 80.4%;
    --input: 38.2 22.0% 80.4%;
    --ring: 30.0 13.0% 9.0%;
    --radius: 0rem;

    --cream-dark: 39.1 28.4% 84.1%;
    --ink: 30.0 13.0% 9.0%;
    --sienna: 26.8 38.1% 43.7%;
    --vermillion-pale: 5.6 63.4% 46.1% / 0.06;
    --text-secondary: 32.5 11.5% 40.8%;
  }

  .dark {
    --background: 30.0 13.0% 9.0%;
    --foreground: 41.2 34.8% 91.0%;
    --card: 25.7 8.4% 16.3%;
    --card-foreground: 41.2 34.8% 91.0%;
    --popover: 25.7 8.4% 16.3%;
    --popover-foreground: 41.2 34.8% 91.0%;
    --primary: 42.0 41.7% 95.3%;
    --primary-foreground: 30.0 13.0% 9.0%;
    --secondary: 30.0 10.0% 13.0%;
    --secondary-foreground: 39.1 28.4% 84.1%;
    --muted: 30.0 10.0% 13.0%;
    --muted-foreground: 33.0 10.3% 62.0%;
    --accent: 5.6 63.4% 46.1%;
    --accent-foreground: 42.0 41.7% 95.3%;
    --destructive: 5.6 70.0% 52.0%;
    --destructive-foreground: 42.0 41.7% 95.3%;
    --border: 30.0 8.0% 22.0%;
    --input: 30.0 8.0% 22.0%;
    --ring: 42.0 41.7% 95.3%;
    --radius: 0rem;

    --cream-dark: 25.0 10.0% 6.0%;
    --ink: 42.0 41.7% 95.3%;
    --sienna: 26.8 38.1% 53.0%;
    --vermillion-pale: 5.6 63.4% 46.1% / 0.12;
    --text-secondary: 33.0 12.0% 65.0%;
  }
}

@layer base {
  * {
    @apply border-border;
  }
  body {
    @apply bg-background text-foreground;
  }
}

/* Scrollbar */
::-webkit-scrollbar { width: 4px; }
::-webkit-scrollbar-track { background: hsl(var(--background)); }
::-webkit-scrollbar-thumb { background: hsl(var(--cream-dark)); }

/* Background texture */
body::before {
  content: '';
  position: fixed;
  inset: 0;
  background-image: url("data:image/svg+xml,%3Csvg width='60' height='60' viewBox='0 0 60 60' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-rule='evenodd'%3E%3Cg fill='%23000000' fill-opacity='0.015'%3E%3Cpath d='M36 34v-4h-2v4h-4v2h4v4h2v-4h4v-2h-4zm0-30V0h-2v4h-4v2h4v4h2V6h4V4h-4zM6 34v-4H4v4H0v2h4v4h2v-4h4v-2H6zM6 4V0H4v4H0v2h4v4h2V6h4V4H6z'/%3E%3C/g%3E%3C/g%3E%3C/svg%3E");
  pointer-events: none;
  z-index: 0;
}

.dark body::before {
  filter: invert(1);
}
```

### 9.2 tailwind.config.ts

Extend the existing Tailwind config with the following additions:

```ts
import type { Config } from 'tailwindcss';

const config: Config = {
  darkMode: ['class'],
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    container: {
      center: true,
      padding: '2rem',
      screens: {
        '2xl': '1400px',
      },
    },
    extend: {
      fontFamily: {
        serif: ['var(--font-serif)', 'Georgia', 'serif'],
        sans: ['var(--font-sans)', 'system-ui', 'sans-serif'],
        mono: ['var(--font-mono)', 'Menlo', 'monospace'],
      },
      colors: {
        border: 'hsl(var(--border))',
        input: 'hsl(var(--input))',
        ring: 'hsl(var(--ring))',
        background: 'hsl(var(--background))',
        foreground: 'hsl(var(--foreground))',
        primary: {
          DEFAULT: 'hsl(var(--primary))',
          foreground: 'hsl(var(--primary-foreground))',
        },
        secondary: {
          DEFAULT: 'hsl(var(--secondary))',
          foreground: 'hsl(var(--secondary-foreground))',
        },
        destructive: {
          DEFAULT: 'hsl(var(--destructive))',
          foreground: 'hsl(var(--destructive-foreground))',
        },
        muted: {
          DEFAULT: 'hsl(var(--muted))',
          foreground: 'hsl(var(--muted-foreground))',
        },
        accent: {
          DEFAULT: 'hsl(var(--accent))',
          foreground: 'hsl(var(--accent-foreground))',
        },
        popover: {
          DEFAULT: 'hsl(var(--popover))',
          foreground: 'hsl(var(--popover-foreground))',
        },
        card: {
          DEFAULT: 'hsl(var(--card))',
          foreground: 'hsl(var(--card-foreground))',
        },
        // Custom editorial tokens
        'cream-dark': 'hsl(var(--cream-dark))',
        ink: 'hsl(var(--ink))',
        sienna: 'hsl(var(--sienna))',
        'text-secondary': 'hsl(var(--text-secondary))',
      },
      borderRadius: {
        lg: '8px',
        md: '4px',
        sm: '2px',
      },
      boxShadow: {
        'offset-sm': '3px 3px 0 hsl(var(--cream-dark))',
        'offset-sm-hover': '4px 4px 0 hsl(var(--cream-dark))',
        'offset-md': '4px 4px 0 hsl(var(--cream-dark))',
        'offset-lg': '6px 6px 0 hsl(var(--cream-dark))',
        'offset-xl': '8px 8px 0 hsl(var(--cream-dark))',
        'card-hover': '0 12px 32px hsl(var(--accent) / 0.08)',
      },
      keyframes: {
        'page-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'slide-in': {
          from: { opacity: '0', transform: 'translateX(-12px)' },
          to: { opacity: '1', transform: 'translateX(0)' },
        },
        'paper-drop': {
          from: { opacity: '0', transform: 'translateY(12px) rotate(-0.3deg)' },
          to: { opacity: '1', transform: 'translateY(0) rotate(0)' },
        },
        'fade-in': {
          from: { opacity: '0' },
          to: { opacity: '1' },
        },
        'modal-drop': {
          from: { opacity: '0', transform: 'translateY(8px)' },
          to: { opacity: '1', transform: 'translateY(0)' },
        },
        'accordion-down': {
          from: { height: '0' },
          to: { height: 'var(--radix-accordion-content-height)' },
        },
        'accordion-up': {
          from: { height: 'var(--radix-accordion-content-height)' },
          to: { height: '0' },
        },
      },
      animation: {
        'page-in': 'page-in 0.4s ease',
        'slide-in': 'slide-in 0.5s ease forwards',
        'paper-drop': 'paper-drop 0.5s cubic-bezier(0.16, 1, 0.3, 1)',
        'fade-in': 'fade-in 0.2s ease',
        'modal-drop': 'modal-drop 0.3s ease',
        'accordion-down': 'accordion-down 0.2s ease-out',
        'accordion-up': 'accordion-up 0.2s ease-out',
      },
    },
  },
  plugins: [require('tailwindcss-animate')],
};

export default config;
```

### 9.3 layout.tsx

Replace font imports in `web/src/app/layout.tsx`:

```tsx
import type { Metadata } from 'next';
import { Newsreader, Source_Sans_3, JetBrains_Mono } from 'next/font/google';
import './globals.css';

const newsreader = Newsreader({
  subsets: ['latin'],
  variable: '--font-serif',
  display: 'swap',
});

const sourceSans3 = Source_Sans_3({
  subsets: ['latin'],
  variable: '--font-sans',
  display: 'swap',
});

const jetbrainsMono = JetBrains_Mono({
  subsets: ['latin'],
  variable: '--font-mono',
  display: 'swap',
});

export const metadata: Metadata = {
  title: 'HireMe - Professional CV Builder',
  description: 'Create stunning, professional CVs with our easy-to-use drag-and-drop editor.',
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body className={`${newsreader.variable} ${sourceSans3.variable} ${jetbrainsMono.variable} font-sans antialiased`}>
        {children}
      </body>
    </html>
  );
}
```

### 9.4 File Change Summary (for T-025)

| File | Change | Scope |
|------|--------|-------|
| `web/src/app/globals.css` | Replace all CSS variables, add scrollbar + texture | Full rewrite |
| `web/tailwind.config.ts` | Add fonts, shadows, animations, custom colors, update radii | Full rewrite |
| `web/src/app/layout.tsx` | Replace Inter with 3 Google Fonts via `next/font/google` | Font imports |
| `web/src/components/ui/button.tsx` | Remap CVA variants: offset shadows, translate hover, sharp corners | Variant definitions |
| `web/src/components/ui/input.tsx` | Border focus (not ring), shadow on focus, sharp corners | Focus styles |
| `web/src/components/ui/dialog.tsx` | Ink border, offset-xl shadow, sharp corners | Container styles |
| `web/src/components/ui/card.tsx` | Sharp corners by default (explicit `rounded-lg` where needed) | Border radius |
| `web/src/components/layout/Header.tsx` | New nav pattern: fixed 60px, ink bottom border, logo stamp, uppercase nav | Full restyle |
| `web/src/app/page.tsx` | Hero section with rule pattern, document list grid rows | Full restyle |
| `web/src/components/editor/EditorLayout.tsx` | Panel borders (2px border), cream-deep center bg | Border + background |
| `web/src/components/editor/CVPreview.tsx` | Paper shadow (offset-lg), paper-drop animation | Shadow + animation |
| `web/src/components/templates/*.tsx` | CV typography: serif names, mono contacts, sans body, dashed section titles | Typography overhaul |

### 9.5 Prototype Reference

The source HTML prototype is at `design-prototypes/01-editorial-craft.html`. When in doubt during T-025 implementation, open this file in a browser as the visual reference.
