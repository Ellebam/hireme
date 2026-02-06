'use client';

import { useEffect } from 'react';
import { useEditorStore } from '@/stores';
import { useUIStore } from '@/stores';

export function useKeyboardShortcuts() {
  const { undo, redo, canUndo, canRedo } = useEditorStore();
  const { zoomIn, zoomOut, resetZoom } = useUIStore();

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Ignore if typing in an input
      const target = e.target as HTMLElement;
      if (
        target.tagName === 'INPUT' ||
        target.tagName === 'TEXTAREA' ||
        target.isContentEditable
      ) {
        return;
      }

      const isMac = navigator.platform.toUpperCase().indexOf('MAC') >= 0;
      const modifier = isMac ? e.metaKey : e.ctrlKey;
      const key = e.key.toLowerCase();

      // Undo: Ctrl/Cmd + Z
      if (modifier && key === 'z' && !e.shiftKey) {
        e.preventDefault();
        if (canUndo()) {
          undo();
        }
        return;
      }

      // Redo: Ctrl/Cmd + Y or Ctrl/Cmd + Shift + Z
      if (
        (modifier && key === 'y') ||
        (modifier && key === 'z' && e.shiftKey)
      ) {
        e.preventDefault();
        if (canRedo()) {
          redo();
        }
        return;
      }

      // Zoom in: Ctrl/Cmd + =
      if (modifier && (key === '=' || key === '+')) {
        e.preventDefault();
        zoomIn();
        return;
      }

      // Zoom out: Ctrl/Cmd + -
      if (modifier && key === '-') {
        e.preventDefault();
        zoomOut();
        return;
      }

      // Reset zoom: Ctrl/Cmd + 0
      if (modifier && key === '0') {
        e.preventDefault();
        resetZoom();
        return;
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [undo, redo, canUndo, canRedo, zoomIn, zoomOut, resetZoom]);
}
