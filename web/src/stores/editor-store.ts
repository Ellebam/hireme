'use client';

/**
 * Editor Store
 * Main state management for the CV editor
 */

import { create } from 'zustand';
import { devtools } from 'zustand/middleware';
import { immer } from 'zustand/middleware/immer';
import type { CV } from '@/types/api';
import type {
  CVContent,
  CVSection,
  SectionType,
  SectionContent,
} from '@/types/cv';
import { getDefaultContent, SCHEMA_VERSION } from '@/types/cv';
import { generateId } from '@/lib/utils';
import { logger } from '@/lib/logger';

// ============================================================================
// Types
// ============================================================================

export type SaveStatus = 'idle' | 'saving' | 'saved' | 'error';

interface EditorState {
  // Data
  cv: CV | null;
  cvContent: CVContent | null;

  // Selection
  selectedSectionId: string | null;
  selectedEntryId: string | null;

  // Edit state
  isDirty: boolean;
  saveStatus: SaveStatus;
  lastSavedAt: Date | null;
  saveError: string | null;

  // Manual save callback (registered by useAutoSave)
  saveNow: (() => void) | null;

  // History (undo/redo)
  history: CVContent[];
  historyIndex: number;
}

interface EditorActions {
  // Initialization
  setCV: (cv: CV) => void;
  reset: () => void;

  // Section operations
  updateSection: (sectionId: string, updates: Partial<CVSection>) => void;
  updateSectionContent: (sectionId: string, content: SectionContent) => void;
  addSection: (type: SectionType, afterSectionId?: string) => string;
  deleteSection: (sectionId: string) => void;
  reorderSections: (activeId: string, overId: string) => void;
  toggleSectionVisibility: (sectionId: string) => void;

  // Entry operations (for sections with entries array)
  reorderEntries: (
    sectionId: string,
    activeId: string,
    overId: string
  ) => void;

  // Selection
  selectSection: (id: string | null) => void;
  selectEntry: (id: string | null) => void;

  // History
  undo: () => void;
  redo: () => void;
  canUndo: () => boolean;
  canRedo: () => boolean;
  pushHistory: () => void;

  // Save state
  markDirty: () => void;
  markSaving: () => void;
  markSaved: () => void;
  markError: (error: string) => void;
  setSaveNow: (fn: (() => void) | null) => void;

  // Direct content update (for batch operations)
  setContent: (content: CVContent) => void;
}

type EditorStore = EditorState & EditorActions;

// ============================================================================
// Constants
// ============================================================================

const MAX_HISTORY = 50;

const initialState: EditorState = {
  cv: null,
  cvContent: null,
  selectedSectionId: null,
  selectedEntryId: null,
  isDirty: false,
  saveStatus: 'idle',
  lastSavedAt: null,
  saveError: null,
  saveNow: null,
  history: [],
  historyIndex: -1,
};

// ============================================================================
// Store
// ============================================================================

export const useEditorStore = create<EditorStore>()(
  devtools(
    immer((set, get) => ({
      ...initialState,

      // ========================================================================
      // Initialization
      // ========================================================================

      setCV: (cv) => {
        logger.info('Store', 'setCV called', {
          id: cv.id,
          title: cv.title,
          hasContent: !!cv.content,
          sectionCount: cv.content?.sections?.length ?? 0,
        });

        if (!cv.content) {
          logger.error('Store', 'CV content is missing!', cv);
          return;
        }

        set((state) => {
          state.cv = cv;
          state.cvContent = cv.content;
          state.isDirty = false;
          state.saveStatus = 'idle';
          state.saveError = null;
          state.selectedSectionId = null;
          state.selectedEntryId = null;
          // Initialize history with current content
          state.history = [cv.content];
          state.historyIndex = 0;
        });

        logger.debug('Store', 'CV loaded successfully', {
          sections: cv.content.sections.map((s) => s.type),
        });
      },

      reset: () => {
        set(initialState);
      },

      // ========================================================================
      // Section Operations
      // ========================================================================

      updateSection: (sectionId, updates) => {
        const { cvContent } = get();
        if (!cvContent) return;

        set((state) => {
          if (!state.cvContent) return;
          const section = state.cvContent.sections.find(
            (s) => s.id === sectionId
          );
          if (section) {
            Object.assign(section, updates);
          }
          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      updateSectionContent: (sectionId, content) => {
        const { cvContent } = get();
        if (!cvContent) return;

        set((state) => {
          if (!state.cvContent) return;
          const section = state.cvContent.sections.find(
            (s) => s.id === sectionId
          );
          if (section) {
            section.content = content;
          }
          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      addSection: (type, afterSectionId) => {
        const { cvContent } = get();
        if (!cvContent) return '';

        const newId = generateId();

        set((state) => {
          if (!state.cvContent) return;

          // Calculate order
          let insertIndex = state.cvContent.sections.length;
          if (afterSectionId) {
            const afterIndex = state.cvContent.sections.findIndex(
              (s) => s.id === afterSectionId
            );
            if (afterIndex !== -1) {
              insertIndex = afterIndex + 1;
            }
          }

          // Create new section
          const newSection: CVSection = {
            id: newId,
            type,
            order: insertIndex,
            visible: true,
            content: getDefaultContent(type),
          };

          // Insert at position
          state.cvContent.sections.splice(insertIndex, 0, newSection);

          // Re-calculate order for all sections
          state.cvContent.sections.forEach((s, i) => {
            s.order = i;
          });

          state.isDirty = true;
          state.selectedSectionId = newId;
        });

        // Push to history AFTER making changes
        get().pushHistory();

        return newId;
      },

      deleteSection: (sectionId) => {
        const { cvContent } = get();
        if (!cvContent) return;

        set((state) => {
          if (!state.cvContent) return;

          state.cvContent.sections = state.cvContent.sections.filter(
            (s) => s.id !== sectionId
          );

          // Re-calculate order
          state.cvContent.sections.forEach((s, i) => {
            s.order = i;
          });

          // Clear selection if deleted section was selected
          if (state.selectedSectionId === sectionId) {
            state.selectedSectionId = null;
            state.selectedEntryId = null;
          }

          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      reorderSections: (activeId, overId) => {
        const { cvContent } = get();
        if (!cvContent || activeId === overId) return;

        set((state) => {
          if (!state.cvContent) return;

          const oldIndex = state.cvContent.sections.findIndex(
            (s) => s.id === activeId
          );
          const newIndex = state.cvContent.sections.findIndex(
            (s) => s.id === overId
          );

          if (oldIndex === -1 || newIndex === -1) return;

          // Remove from old position and insert at new
          const [removed] = state.cvContent.sections.splice(oldIndex, 1);
          state.cvContent.sections.splice(newIndex, 0, removed);

          // Re-calculate order
          state.cvContent.sections.forEach((s, i) => {
            s.order = i;
          });

          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      toggleSectionVisibility: (sectionId) => {
        set((state) => {
          if (!state.cvContent) return;
          const section = state.cvContent.sections.find(
            (s) => s.id === sectionId
          );
          if (section) {
            section.visible = section.visible === false ? true : false;
          }
          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      // ========================================================================
      // Entry Operations
      // ========================================================================

      reorderEntries: (sectionId, activeId, overId) => {
        const { cvContent } = get();
        if (!cvContent || activeId === overId) return;

        set((state) => {
          if (!state.cvContent) return;

          const section = state.cvContent.sections.find(
            (s) => s.id === sectionId
          );
          if (!section) return;

          // Handle sections with entries array
          const content = section.content as { entries?: { id: string }[] };
          if (!content.entries) return;

          const oldIndex = content.entries.findIndex((e) => e.id === activeId);
          const newIndex = content.entries.findIndex((e) => e.id === overId);

          if (oldIndex === -1 || newIndex === -1) return;

          const [removed] = content.entries.splice(oldIndex, 1);
          content.entries.splice(newIndex, 0, removed);

          state.isDirty = true;
        });

        // Push to history AFTER making changes
        get().pushHistory();
      },

      // ========================================================================
      // Selection
      // ========================================================================

      selectSection: (id) => {
        set((state) => {
          state.selectedSectionId = id;
          state.selectedEntryId = null;
        });
      },

      selectEntry: (id) => {
        set((state) => {
          state.selectedEntryId = id;
        });
      },

      // ========================================================================
      // History (Undo/Redo)
      // ========================================================================

      pushHistory: () => {
        // This is called AFTER making changes
        // It snapshots the current (changed) state to history
        set((state) => {
          if (!state.cvContent) return;

          // Remove any future history if we're not at the end
          const newHistory = state.history.slice(0, state.historyIndex + 1);

          // Deep clone current content and add to history
          const snapshot = JSON.parse(JSON.stringify(state.cvContent));
          newHistory.push(snapshot);

          // Limit history size
          if (newHistory.length > MAX_HISTORY) {
            newHistory.shift();
            // Adjust index if we removed from the beginning
            state.historyIndex = Math.max(0, state.historyIndex);
          }

          state.history = newHistory;
          state.historyIndex = newHistory.length - 1;
        });
      },

      undo: () => {
        const { history, historyIndex, cvContent } = get();
        if (!cvContent || historyIndex <= 0) return;

        set((state) => {
          state.historyIndex -= 1;
          state.cvContent = JSON.parse(
            JSON.stringify(state.history[state.historyIndex])
          );
          state.isDirty = true;
        });
      },

      redo: () => {
        const { history, historyIndex, cvContent } = get();
        if (!cvContent || historyIndex >= history.length - 1) return;

        set((state) => {
          state.historyIndex += 1;
          state.cvContent = JSON.parse(
            JSON.stringify(state.history[state.historyIndex])
          );
          state.isDirty = true;
        });
      },

      canUndo: () => {
        const { historyIndex } = get();
        return historyIndex > 0;
      },

      canRedo: () => {
        const { history, historyIndex } = get();
        return historyIndex < history.length - 1;
      },

      // ========================================================================
      // Save State
      // ========================================================================

      markDirty: () => {
        set((state) => {
          state.isDirty = true;
        });
      },

      markSaving: () => {
        set((state) => {
          state.saveStatus = 'saving';
          state.saveError = null;
        });
      },

      markSaved: () => {
        set((state) => {
          state.isDirty = false;
          state.saveStatus = 'saved';
          state.lastSavedAt = new Date();
          state.saveError = null;
        });
      },

      markError: (error) => {
        set((state) => {
          state.saveStatus = 'error';
          state.saveError = error;
        });
      },

      setSaveNow: (fn) => {
        set((state) => {
          state.saveNow = fn;
        });
      },

      // ========================================================================
      // Direct Content Update
      // ========================================================================

      setContent: (content) => {
        set((state) => {
          state.cvContent = content;
          state.isDirty = true;
        });
      },
    })),
    { name: 'hireme-editor' }
  )
);

// ============================================================================
// Selectors
// ============================================================================

/** Get current section by ID */
export const useSelectedSection = () => {
  return useEditorStore((state) => {
    if (!state.cvContent || !state.selectedSectionId) return null;
    return state.cvContent.sections.find(
      (s) => s.id === state.selectedSectionId
    );
  });
};

/** Get sections sorted by order */
export const useSortedSections = () => {
  return useEditorStore((state) => {
    if (!state.cvContent) return [];
    return [...state.cvContent.sections].sort((a, b) => a.order - b.order);
  });
};

/** Get visible sections only */
export const useVisibleSections = () => {
  return useEditorStore((state) => {
    if (!state.cvContent) return [];
    return state.cvContent.sections
      .filter((s) => s.visible !== false)
      .sort((a, b) => a.order - b.order);
  });
};
