'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { Loader2 } from 'lucide-react';
import { TooltipProvider } from '@/components/ui/tooltip';
import { EditorLayout } from '@/components/editor';
import { useEditorStore } from '@/stores';
import { useAutoSave, useKeyboardShortcuts } from '@/hooks';
import { api, ApiError } from '@/lib/api';
import { starterTemplate } from '@/lib/templates';
import { logger } from '@/lib/logger';

export default function EditorPage() {
  const router = useRouter();
  const { setCV, reset } = useEditorStore();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Enable auto-save and keyboard shortcuts
  useAutoSave();
  useKeyboardShortcuts();

  useEffect(() => {
    let isActive = true;

    async function loadOrCreateCV() {
      logger.info('Editor', 'Loading CV...');

      try {
        // Try to load existing CV
        const cv = await api.cv.get();
        if (!isActive) return;
        logger.info('Editor', 'CV loaded from API', {
          id: cv.id,
          title: cv.title,
          sectionCount: cv.content?.sections?.length,
        });
        setCV(cv);
      } catch (err) {
        if (!isActive) return;
        if (err instanceof ApiError && err.isNotFound) {
          logger.info('Editor', 'No CV found, creating new one');
          // No CV exists, create one with starter template
          try {
            const template = starterTemplate();
            logger.debug('Editor', 'Using starter template', {
              sectionCount: template.sections.length,
            });
            const newCV = await api.cv.create({
              title: template.title || 'My CV',
              content: template,
            });
            if (!isActive) return;
            logger.info('Editor', 'New CV created', { id: newCV.id });
            setCV(newCV);
          } catch (createErr) {
            if (!isActive) return;
            logger.error('Editor', 'Failed to create CV', createErr);
            setError('Failed to create CV. Please try again.');
          }
        } else {
          logger.error('Editor', 'Failed to load CV', err);
          setError('Failed to load CV. Please try again.');
        }
      } finally {
        if (isActive) {
          setLoading(false);
        }
      }
    }

    loadOrCreateCV();

    // Cleanup on unmount
    return () => {
      isActive = false;
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
              onClick={() => router.push('/')}
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
