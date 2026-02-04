'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Loader2 } from 'lucide-react';
import { TooltipProvider } from '@/components/ui/tooltip';
import { EditorLayout } from '@/components/editor';
import { useEditorStore } from '@/stores';
import { useAutoSave, useKeyboardShortcuts } from '@/hooks';
import { api, ApiError } from '@/lib/api';
import type { CV } from '@/types/api';
import type { CVContent } from '@/types/cv';

// Default CV content for new CVs
const defaultCVContent: CVContent = {
  schemaVersion: '1.0.0',
  templateId: 'modern',
  locale: 'en',
  title: 'My CV',
  sections: [
    {
      id: 'sec-personal',
      type: 'personal',
      order: 0,
      visible: true,
      content: {
        firstName: '',
        lastName: '',
        jobTitle: '',
        email: '',
        phone: '',
        location: '',
        links: [],
      },
    },
    {
      id: 'sec-summary',
      type: 'summary',
      order: 1,
      visible: true,
      content: {
        text: '',
      },
    },
    {
      id: 'sec-experience',
      type: 'experience',
      order: 2,
      visible: true,
      title: 'Experience',
      content: {
        entries: [],
      },
    },
    {
      id: 'sec-education',
      type: 'education',
      order: 3,
      visible: true,
      title: 'Education',
      content: {
        entries: [],
      },
    },
    {
      id: 'sec-skills',
      type: 'skills',
      order: 4,
      visible: true,
      title: 'Skills',
      content: {
        categories: [],
      },
    },
  ],
  styling: {
    primaryColor: '#2563eb',
    secondaryColor: '#64748b',
    fontFamily: 'inter',
    fontSize: 'medium',
    lineHeight: 'normal',
    showIcons: true,
  },
};

export default function EditorPage() {
  const router = useRouter();
  const { setCV, reset } = useEditorStore();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Enable auto-save and keyboard shortcuts
  useAutoSave();
  useKeyboardShortcuts();

  useEffect(() => {
    async function loadOrCreateCV() {
      try {
        // Try to load existing CV
        const cv = await api.cv.get();
        setCV(cv);
      } catch (err) {
        if (err instanceof ApiError && err.isNotFound) {
          // No CV exists, create one with default content
          try {
            const newCV = await api.cv.create({
              title: 'My CV',
              content: defaultCVContent,
            });
            setCV(newCV);
          } catch (createErr) {
            console.error('Failed to create CV:', createErr);
            setError('Failed to create CV. Please try again.');
          }
        } else {
          console.error('Failed to load CV:', err);
          setError('Failed to load CV. Please try again.');
        }
      } finally {
        setLoading(false);
      }
    }

    loadOrCreateCV();

    // Cleanup on unmount
    return () => {
      reset();
    };
  }, [setCV, reset]);

  if (loading) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="flex flex-col items-center gap-4">
          <Loader2 className="h-8 w-8 animate-spin text-primary" />
          <p className="text-sm text-muted-foreground">Loading editor...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="h-screen flex items-center justify-center">
        <div className="flex flex-col items-center gap-4 text-center max-w-md px-4">
          <div className="p-4 rounded-full bg-destructive/10">
            <span className="text-4xl">⚠️</span>
          </div>
          <h1 className="text-xl font-semibold">Something went wrong</h1>
          <p className="text-muted-foreground">{error}</p>
          <div className="flex gap-2">
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-primary text-primary-foreground rounded-md text-sm font-medium hover:bg-primary/90"
            >
              Try Again
            </button>
            <button
              onClick={() => router.push('/dashboard')}
              className="px-4 py-2 bg-muted text-muted-foreground rounded-md text-sm font-medium hover:bg-muted/80"
            >
              Go to Dashboard
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <TooltipProvider delayDuration={300}>
      <div className="h-screen">
        <EditorLayout />
      </div>
    </TooltipProvider>
  );
}
