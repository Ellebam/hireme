'use client';

/**
 * UI Store
 * State management for UI elements (panels, modals, etc.)
 */

import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

// ============================================================================
// Types
// ============================================================================

interface UIState {
  // Sidebars
  leftSidebarOpen: boolean;
  rightSidebarOpen: boolean;

  // Modals
  exportModalOpen: boolean;
  templateModalOpen: boolean;
  deleteConfirmModalOpen: boolean;
  deleteSectionId: string | null;

  // Preview
  previewScale: number;
  previewScrollPosition: number;

  // Mobile
  mobileMenuOpen: boolean;
  activeMobilePanel: 'palette' | 'preview' | 'properties';

  // Theme (future)
  theme: 'light' | 'dark' | 'system';
}

interface UIActions {
  // Sidebars
  toggleLeftSidebar: () => void;
  toggleRightSidebar: () => void;
  setLeftSidebarOpen: (open: boolean) => void;
  setRightSidebarOpen: (open: boolean) => void;

  // Modals
  openExportModal: () => void;
  closeExportModal: () => void;
  openTemplateModal: () => void;
  closeTemplateModal: () => void;
  openDeleteConfirmModal: (sectionId: string) => void;
  closeDeleteConfirmModal: () => void;

  // Preview
  setPreviewScale: (scale: number) => void;
  zoomIn: () => void;
  zoomOut: () => void;
  resetZoom: () => void;
  setPreviewScrollPosition: (position: number) => void;

  // Mobile
  toggleMobileMenu: () => void;
  setActiveMobilePanel: (panel: 'palette' | 'preview' | 'properties') => void;

  // Theme
  setTheme: (theme: 'light' | 'dark' | 'system') => void;

  // Reset
  reset: () => void;
}

type UIStore = UIState & UIActions;

// ============================================================================
// Constants
// ============================================================================

const ZOOM_STEP = 0.1;
const MIN_ZOOM = 0.5;
const MAX_ZOOM = 2.0;
const DEFAULT_ZOOM = 1.0;

const initialState: UIState = {
  leftSidebarOpen: true,
  rightSidebarOpen: true,
  exportModalOpen: false,
  templateModalOpen: false,
  deleteConfirmModalOpen: false,
  deleteSectionId: null,
  previewScale: DEFAULT_ZOOM,
  previewScrollPosition: 0,
  mobileMenuOpen: false,
  activeMobilePanel: 'preview',
  theme: 'system',
};

// ============================================================================
// Store
// ============================================================================

export const useUIStore = create<UIStore>()(
  devtools(
    persist(
      (set, get) => ({
        ...initialState,

        // ======================================================================
        // Sidebars
        // ======================================================================

        toggleLeftSidebar: () => {
          set((state) => ({ leftSidebarOpen: !state.leftSidebarOpen }));
        },

        toggleRightSidebar: () => {
          set((state) => ({ rightSidebarOpen: !state.rightSidebarOpen }));
        },

        setLeftSidebarOpen: (open) => {
          set({ leftSidebarOpen: open });
        },

        setRightSidebarOpen: (open) => {
          set({ rightSidebarOpen: open });
        },

        // ======================================================================
        // Modals
        // ======================================================================

        openExportModal: () => {
          set({ exportModalOpen: true });
        },

        closeExportModal: () => {
          set({ exportModalOpen: false });
        },

        openTemplateModal: () => {
          set({ templateModalOpen: true });
        },

        closeTemplateModal: () => {
          set({ templateModalOpen: false });
        },

        openDeleteConfirmModal: (sectionId) => {
          set({ deleteConfirmModalOpen: true, deleteSectionId: sectionId });
        },

        closeDeleteConfirmModal: () => {
          set({ deleteConfirmModalOpen: false, deleteSectionId: null });
        },

        // ======================================================================
        // Preview
        // ======================================================================

        setPreviewScale: (scale) => {
          const clampedScale = Math.min(MAX_ZOOM, Math.max(MIN_ZOOM, scale));
          set({ previewScale: clampedScale });
        },

        zoomIn: () => {
          const { previewScale } = get();
          const newScale = Math.min(MAX_ZOOM, previewScale + ZOOM_STEP);
          set({ previewScale: Math.round(newScale * 10) / 10 });
        },

        zoomOut: () => {
          const { previewScale } = get();
          const newScale = Math.max(MIN_ZOOM, previewScale - ZOOM_STEP);
          set({ previewScale: Math.round(newScale * 10) / 10 });
        },

        resetZoom: () => {
          set({ previewScale: DEFAULT_ZOOM });
        },

        setPreviewScrollPosition: (position) => {
          set({ previewScrollPosition: position });
        },

        // ======================================================================
        // Mobile
        // ======================================================================

        toggleMobileMenu: () => {
          set((state) => ({ mobileMenuOpen: !state.mobileMenuOpen }));
        },

        setActiveMobilePanel: (panel) => {
          set({ activeMobilePanel: panel });
        },

        // ======================================================================
        // Theme
        // ======================================================================

        setTheme: (theme) => {
          set({ theme });
        },

        // ======================================================================
        // Reset
        // ======================================================================

        reset: () => {
          set(initialState);
        },
      }),
      {
        name: 'hireme-ui',
        partialize: (state) => ({
          // Only persist these UI preferences
          leftSidebarOpen: state.leftSidebarOpen,
          rightSidebarOpen: state.rightSidebarOpen,
          previewScale: state.previewScale,
          theme: state.theme,
        }),
      }
    ),
    { name: 'hireme-ui' }
  )
);

// ============================================================================
// Selectors
// ============================================================================

/** Check if any modal is open */
export const useIsAnyModalOpen = () => {
  return useUIStore(
    (state) =>
      state.exportModalOpen ||
      state.templateModalOpen ||
      state.deleteConfirmModalOpen
  );
};

/** Get preview scale as percentage string */
export const usePreviewScalePercent = () => {
  return useUIStore((state) => `${Math.round(state.previewScale * 100)}%`);
};
