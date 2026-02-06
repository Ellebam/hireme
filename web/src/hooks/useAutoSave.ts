'use client';

import { useEffect, useRef, useCallback } from 'react';
import { useEditorStore } from '@/stores';
import { api } from '@/lib/api';
import { debounce } from '@/lib/utils';

const AUTOSAVE_DELAY = 2000; // 2 seconds

export function useAutoSave() {
  const { cv, cvContent, isDirty, markSaving, markSaved, markError } = useEditorStore();
  const saveInProgressRef = useRef(false);

  // Create debounced save function
  const debouncedSave = useCallback(
    debounce(async () => {
      const currentCV = useEditorStore.getState().cv;
      const currentContent = useEditorStore.getState().cvContent;
      const currentIsDirty = useEditorStore.getState().isDirty;

      if (!currentCV || !currentContent || !currentIsDirty || saveInProgressRef.current) {
        return;
      }

      saveInProgressRef.current = true;
      markSaving();

      try {
        await api.cv.update(currentCV.id, { content: currentContent });
        markSaved();
      } catch (error) {
        console.error('Auto-save failed:', error);
        markError(error instanceof Error ? error.message : 'Failed to save');
      } finally {
        saveInProgressRef.current = false;
      }
    }, AUTOSAVE_DELAY),
    [markSaving, markSaved, markError]
  );

  // Trigger auto-save when content changes
  useEffect(() => {
    if (isDirty && cv && cvContent) {
      debouncedSave();
    }
  }, [isDirty, cv, cvContent, debouncedSave]);

  // Cancel pending saves on unmount
  useEffect(() => {
    return () => {
      debouncedSave.cancel();
    };
  }, [debouncedSave]);
}
