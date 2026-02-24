'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { Loader2 } from 'lucide-react';
import { AppShell } from '@/components/layout';
import { api, ApiError } from '@/lib/api';
import { starterTemplate } from '@/lib/templates';
import type { TemplateId } from '@/types/cv';

const templates = [
  {
    id: 'classic' as TemplateId,
    name: 'Classic',
    number: '01',
    description:
      'Traditional single-column layout. Clean hierarchy, familiar structure.',
    tag: 'Popular',
  },
  {
    id: 'modern' as TemplateId,
    name: 'Modern',
    number: '02',
    description:
      'Contemporary design with timeline accent. Stand out with style.',
    tag: 'Recommended',
  },
  {
    id: 'visionary' as TemplateId,
    name: 'Visionary',
    number: '03',
    description:
      'Bold two-column layout with sidebar. Maximum visual impact.',
    tag: 'New',
  },
];

function ClassicPreview() {
  return (
    <div className="w-[110px] h-[155px] bg-card border border-border rounded shadow-[0_6px_16px_rgba(0,0,0,0.06)] p-2.5 transition-transform duration-300 group-hover:-translate-y-1.5 group-hover:-rotate-1">
      <div className="h-[5px] w-[70%] bg-primary rounded-sm mb-1.5" />
      <div className="h-[3px] w-[50%] bg-accent/50 rounded-sm mb-3" />
      <div className="space-y-1.5">
        <div className="h-[2px] w-full bg-border rounded-sm" />
        <div className="h-[2px] w-[90%] bg-border rounded-sm" />
        <div className="h-[2px] w-[80%] bg-border rounded-sm" />
        <div className="h-[2px] w-full bg-border rounded-sm" />
        <div className="h-[2px] w-[85%] bg-border rounded-sm" />
      </div>
    </div>
  );
}

function ModernPreview() {
  return (
    <div className="w-[110px] h-[155px] bg-card border border-border rounded shadow-[0_6px_16px_rgba(0,0,0,0.06)] p-2.5 transition-transform duration-300 group-hover:-translate-y-1.5 group-hover:-rotate-1">
      <div className="h-[5px] w-[70%] bg-primary rounded-sm mb-1.5" />
      <div className="h-[3px] w-[50%] bg-accent/50 rounded-sm mb-3" />
      <div className="flex gap-2">
        <div className="w-[3px] bg-accent/60 self-stretch rounded-sm" />
        <div className="flex-1 space-y-1.5">
          <div className="h-[2px] w-full bg-border rounded-sm" />
          <div className="h-[2px] w-[90%] bg-border rounded-sm" />
          <div className="h-[2px] w-[80%] bg-border rounded-sm" />
          <div className="h-[2px] w-full bg-border rounded-sm" />
          <div className="h-[2px] w-[85%] bg-border rounded-sm" />
        </div>
      </div>
    </div>
  );
}

function VisionaryPreview() {
  return (
    <div className="w-[110px] h-[155px] bg-card border border-border rounded shadow-[0_6px_16px_rgba(0,0,0,0.06)] flex overflow-hidden transition-transform duration-300 group-hover:-translate-y-1.5 group-hover:-rotate-1">
      <div className="w-[35%] bg-primary" />
      <div className="flex-1 p-2 space-y-1.5">
        <div className="h-[3px] w-[80%] bg-border rounded-sm" />
        <div className="h-[2px] w-full bg-border rounded-sm" />
        <div className="h-[2px] w-[90%] bg-border rounded-sm" />
        <div className="h-[2px] w-[70%] bg-border rounded-sm" />
        <div className="h-[2px] w-full bg-border rounded-sm" />
        <div className="h-[2px] w-[85%] bg-border rounded-sm" />
      </div>
    </div>
  );
}

const previewComponents: Record<string, () => React.JSX.Element> = {
  classic: ClassicPreview,
  modern: ModernPreview,
  visionary: VisionaryPreview,
};

export default function TemplatesPage() {
  const router = useRouter();
  const [creating, setCreating] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  async function handleTemplateClick(templateId: TemplateId) {
    if (creating) return;
    setCreating(templateId);
    setError(null);

    try {
      const content = starterTemplate();
      content.templateId = templateId;
      const newCV = await api.cv.create({
        title: 'Untitled CV',
        content,
      });
      router.push(`/editor/${newCV.id}`);
    } catch (err) {
      if (err instanceof ApiError && err.isForbidden) {
        setError('CV limit reached. Delete an existing CV to create a new one.');
      } else {
        setError('Failed to create CV. Please try again.');
      }
      setCreating(null);
    }
  }

  return (
    <AppShell>
      <div className="max-w-[960px] px-10 pt-[60px] pb-20">
        {/* Hero */}
        <div className="flex items-center gap-3 mb-5 animate-slide-in opacity-0 [animation-delay:0.1s]">
          <div className="w-9 h-0.5 bg-accent" />
          <span className="font-mono text-[0.6875rem] font-medium uppercase tracking-[0.1875em] text-accent">
            Templates
          </span>
        </div>

        <h1 className="font-serif text-[2.75rem] font-bold leading-[1.08] tracking-[-0.03em] text-primary mb-5 animate-slide-in opacity-0 [animation-delay:0.2s]">
          Choose a{' '}
          <em className="font-normal italic text-sienna">format.</em>
        </h1>

        <p className="text-base text-[hsl(var(--text-secondary))] leading-[1.7] max-w-[460px] mb-11 animate-slide-in opacity-0 [animation-delay:0.35s]">
          Each template shapes perception. The right one amplifies your story.
        </p>

        {error && (
          <div className="mb-6 p-4 border border-destructive/20 bg-destructive/5 rounded-lg text-sm text-destructive">
            {error}
          </div>
        )}

        {/* Template Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4 animate-slide-in opacity-0 [animation-delay:0.45s]">
          {templates.map((template) => {
            const Preview = previewComponents[template.id];
            const isCreating = creating === template.id;
            return (
              <button
                key={template.id}
                onClick={() => handleTemplateClick(template.id)}
                disabled={creating !== null}
                className="group bg-card border border-border rounded-lg cursor-pointer transition-all duration-[250ms] ease-out hover:border-accent/20 hover:shadow-[0_12px_32px_rgba(192,57,43,0.08)] hover:-translate-y-1 text-left disabled:opacity-60 disabled:cursor-not-allowed"
              >
                {/* Preview Area */}
                <div className="h-[200px] bg-secondary flex items-center justify-center rounded-t-lg">
                  {isCreating ? (
                    <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
                  ) : (
                    <Preview />
                  )}
                </div>

                {/* Info */}
                <div className="p-5">
                  <span className="font-mono text-[0.6875rem] text-muted-foreground">
                    {template.number}
                  </span>
                  <h3 className="font-serif text-[1.1875rem] font-semibold text-primary tracking-[-0.01em] mt-0.5">
                    {template.name}
                  </h3>
                  <p className="text-[0.8125rem] text-[hsl(var(--text-secondary))] leading-[1.5] mt-1.5 mb-3">
                    {template.description}
                  </p>
                  <span className="inline-block font-mono text-[0.625rem] font-medium uppercase tracking-[0.0625em] text-accent bg-[hsl(var(--vermillion-pale))] px-3 py-1 rounded-full">
                    {template.tag}
                  </span>
                </div>
              </button>
            );
          })}
        </div>
      </div>
    </AppShell>
  );
}
