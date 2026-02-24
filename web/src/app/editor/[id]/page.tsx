'use client';

import { Suspense, useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { Loader2 } from 'lucide-react';
import { AppShell } from '@/components/layout';
import { EditorLayout } from '@/components/editor';
import { useEditorStore } from '@/stores';
import { useAutoSave, useKeyboardShortcuts } from '@/hooks';
import { api, ApiError } from '@/lib/api';
import { logger } from '@/lib/logger';

function EditorContent() {
  const router = useRouter();
  const params = useParams<{ id: string }>();
  const { setCV, reset } = useEditorStore();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Enable auto-save and keyboard shortcuts
  useAutoSave();
  useKeyboardShortcuts();

  useEffect(() => {
    let isActive = true;

    async function loadCV() {
      logger.info('Editor', 'Loading CV...', { id: params.id });

      try {
        const cv = await api.cv.get(params.id);
        if (!isActive) return;
        logger.info('Editor', 'CV loaded from API', {
          id: cv.id,
          title: cv.title,
          sectionCount: cv.content?.sections?.length,
        });
        setCV(cv);
      } catch (err) {
        if (!isActive) return;
        if (err instanceof ApiError && (err.isNotFound || err.isForbidden)) {
          logger.warn('Editor', 'CV not found or access denied', { id: params.id });
          router.replace('/');
          return;
        }
        logger.error('Editor', 'Failed to load CV', err);
        setError('Failed to load CV. Please try again.');
      } finally {
        if (isActive) {
          setLoading(false);
        }
      }
    }

    loadCV();

    return () => {
      isActive = false;
      reset();
    };
  }, [params.id, setCV, reset, router]);

  return (
    <AppShell fullHeight>
      {loading ? (
        <div className="h-full flex items-center justify-center">
          <div className="flex flex-col items-center gap-4">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-sm text-muted-foreground">Loading editor...</p>
          </div>
        </div>
      ) : error ? (
        <div className="h-full flex items-center justify-center">
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
      ) : (
        <EditorLayout />
      )}
    </AppShell>
  );
}

export default function EditorIdPage() {
  return (
    <Suspense
      fallback={
        <div className="h-screen flex items-center justify-center">
          <div className="flex flex-col items-center gap-4">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
            <p className="text-sm text-muted-foreground">Loading editor...</p>
          </div>
        </div>
      }
    >
      <EditorContent />
    </Suspense>
  );
}
