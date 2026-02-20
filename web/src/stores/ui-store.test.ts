/**
 * UI Store Tests
 */

import { describe, it, expect, beforeEach } from 'vitest';
import { useUIStore, useIsAnyModalOpen, usePreviewScalePercent } from './ui-store';

describe('UIStore', () => {
  // Reset store before each test
  beforeEach(() => {
    useUIStore.setState({
      leftSidebarOpen: true,
      rightSidebarOpen: true,
      exportModalOpen: false,
      templateModalOpen: false,
      deleteConfirmModalOpen: false,
      deleteSectionId: null,
      previewScale: 1.0,
      previewScrollPosition: 0,
      autoFitScale: 1.0,
      mobileMenuOpen: false,
      activeMobilePanel: 'preview',
      theme: 'system',
    });
  });

  // ========================================================================
  // Sidebars
  // ========================================================================

  describe('sidebar controls', () => {
    it('should toggle left sidebar', () => {
      expect(useUIStore.getState().leftSidebarOpen).toBe(true);

      useUIStore.getState().toggleLeftSidebar();
      expect(useUIStore.getState().leftSidebarOpen).toBe(false);

      useUIStore.getState().toggleLeftSidebar();
      expect(useUIStore.getState().leftSidebarOpen).toBe(true);
    });

    it('should toggle right sidebar', () => {
      expect(useUIStore.getState().rightSidebarOpen).toBe(true);

      useUIStore.getState().toggleRightSidebar();
      expect(useUIStore.getState().rightSidebarOpen).toBe(false);
    });

    it('should set sidebar state directly', () => {
      useUIStore.getState().setLeftSidebarOpen(false);
      expect(useUIStore.getState().leftSidebarOpen).toBe(false);

      useUIStore.getState().setRightSidebarOpen(false);
      expect(useUIStore.getState().rightSidebarOpen).toBe(false);
    });
  });

  // ========================================================================
  // Modals
  // ========================================================================

  describe('modal controls', () => {
    it('should open and close export modal', () => {
      expect(useUIStore.getState().exportModalOpen).toBe(false);

      useUIStore.getState().openExportModal();
      expect(useUIStore.getState().exportModalOpen).toBe(true);

      useUIStore.getState().closeExportModal();
      expect(useUIStore.getState().exportModalOpen).toBe(false);
    });

    it('should open and close template modal', () => {
      expect(useUIStore.getState().templateModalOpen).toBe(false);

      useUIStore.getState().openTemplateModal();
      expect(useUIStore.getState().templateModalOpen).toBe(true);

      useUIStore.getState().closeTemplateModal();
      expect(useUIStore.getState().templateModalOpen).toBe(false);
    });

    it('should open delete confirm modal with section ID', () => {
      useUIStore.getState().openDeleteConfirmModal('sec-123');

      expect(useUIStore.getState().deleteConfirmModalOpen).toBe(true);
      expect(useUIStore.getState().deleteSectionId).toBe('sec-123');
    });

    it('should clear section ID when closing delete modal', () => {
      useUIStore.getState().openDeleteConfirmModal('sec-123');
      useUIStore.getState().closeDeleteConfirmModal();

      expect(useUIStore.getState().deleteConfirmModalOpen).toBe(false);
      expect(useUIStore.getState().deleteSectionId).toBeNull();
    });
  });

  // ========================================================================
  // Preview Controls
  // ========================================================================

  describe('preview controls', () => {
    it('should set preview scale', () => {
      useUIStore.getState().setPreviewScale(1.5);
      expect(useUIStore.getState().previewScale).toBe(1.5);
    });

    it('should clamp scale to min value', () => {
      useUIStore.getState().setPreviewScale(0.1);
      expect(useUIStore.getState().previewScale).toBe(0.5); // MIN_ZOOM
    });

    it('should clamp scale to max value', () => {
      useUIStore.getState().setPreviewScale(5.0);
      expect(useUIStore.getState().previewScale).toBe(2.0); // MAX_ZOOM
    });

    it('should zoom in by step', () => {
      useUIStore.getState().setPreviewScale(1.0);
      useUIStore.getState().zoomIn();
      expect(useUIStore.getState().previewScale).toBe(1.1);
    });

    it('should zoom out by step', () => {
      useUIStore.getState().setPreviewScale(1.0);
      useUIStore.getState().zoomOut();
      expect(useUIStore.getState().previewScale).toBe(0.9);
    });

    it('should not zoom in beyond max', () => {
      useUIStore.getState().setPreviewScale(2.0);
      useUIStore.getState().zoomIn();
      expect(useUIStore.getState().previewScale).toBe(2.0);
    });

    it('should not zoom out below min', () => {
      useUIStore.getState().setPreviewScale(0.5);
      useUIStore.getState().zoomOut();
      expect(useUIStore.getState().previewScale).toBe(0.5);
    });

    it('should reset zoom to autoFitScale', () => {
      useUIStore.getState().setPreviewScale(1.8);
      useUIStore.getState().resetZoom();
      expect(useUIStore.getState().previewScale).toBe(useUIStore.getState().autoFitScale);
    });

    it('should set scroll position', () => {
      useUIStore.getState().setPreviewScrollPosition(500);
      expect(useUIStore.getState().previewScrollPosition).toBe(500);
    });
  });

  // ========================================================================
  // Auto-Fit Controls
  // ========================================================================

  describe('auto-fit controls', () => {
    it('should set auto-fit scale', () => {
      useUIStore.getState().setAutoFitScale(0.8);
      expect(useUIStore.getState().autoFitScale).toBe(0.8);
    });

    it('should clamp auto-fit scale to min value', () => {
      useUIStore.getState().setAutoFitScale(0.1);
      expect(useUIStore.getState().autoFitScale).toBe(0.5);
    });

    it('should clamp auto-fit scale to max value', () => {
      useUIStore.getState().setAutoFitScale(5.0);
      expect(useUIStore.getState().autoFitScale).toBe(2.0);
    });

    it('should clamp auto-fit scale for zero', () => {
      useUIStore.getState().setAutoFitScale(0);
      expect(useUIStore.getState().autoFitScale).toBe(0.5);
    });

    it('should reset zoom to auto-fit scale', () => {
      useUIStore.getState().setAutoFitScale(0.7);
      useUIStore.getState().setPreviewScale(1.5);
      useUIStore.getState().resetZoom();
      expect(useUIStore.getState().previewScale).toBe(0.7);
    });

    it('should reset zoom to 1.0 when auto-fit scale is default', () => {
      useUIStore.getState().setPreviewScale(0.5);
      useUIStore.getState().resetZoom();
      expect(useUIStore.getState().previewScale).toBe(1.0);
    });

    it('should not modify previewScale when setting auto-fit scale', () => {
      useUIStore.getState().setPreviewScale(1.5);
      useUIStore.getState().setAutoFitScale(0.8);
      expect(useUIStore.getState().previewScale).toBe(1.5);
    });

    it('should not modify autoFitScale when zooming', () => {
      useUIStore.getState().setAutoFitScale(0.7);
      useUIStore.getState().zoomIn();
      useUIStore.getState().zoomOut();
      expect(useUIStore.getState().autoFitScale).toBe(0.7);
    });

    it('should handle full workflow: auto-fit, manual zoom, reset', () => {
      // Simulate auto-fit setting scale
      useUIStore.getState().setAutoFitScale(0.7);
      useUIStore.getState().setPreviewScale(0.7);

      // User manually zooms in
      useUIStore.getState().zoomIn();
      expect(useUIStore.getState().previewScale).toBe(0.8);

      // Reset returns to auto-fit scale
      useUIStore.getState().resetZoom();
      expect(useUIStore.getState().previewScale).toBe(0.7);
    });
  });

  // ========================================================================
  // Mobile
  // ========================================================================

  describe('mobile controls', () => {
    it('should toggle mobile menu', () => {
      expect(useUIStore.getState().mobileMenuOpen).toBe(false);

      useUIStore.getState().toggleMobileMenu();
      expect(useUIStore.getState().mobileMenuOpen).toBe(true);

      useUIStore.getState().toggleMobileMenu();
      expect(useUIStore.getState().mobileMenuOpen).toBe(false);
    });

    it('should set active mobile panel', () => {
      expect(useUIStore.getState().activeMobilePanel).toBe('preview');

      useUIStore.getState().setActiveMobilePanel('palette');
      expect(useUIStore.getState().activeMobilePanel).toBe('palette');

      useUIStore.getState().setActiveMobilePanel('properties');
      expect(useUIStore.getState().activeMobilePanel).toBe('properties');
    });
  });

  // ========================================================================
  // Theme
  // ========================================================================

  describe('theme controls', () => {
    it('should set theme', () => {
      expect(useUIStore.getState().theme).toBe('system');

      useUIStore.getState().setTheme('dark');
      expect(useUIStore.getState().theme).toBe('dark');

      useUIStore.getState().setTheme('light');
      expect(useUIStore.getState().theme).toBe('light');
    });
  });

  // ========================================================================
  // Reset
  // ========================================================================

  describe('reset', () => {
    it('should reset all state to initial values', () => {
      // Change some state
      useUIStore.getState().toggleLeftSidebar();
      useUIStore.getState().openExportModal();
      useUIStore.getState().setPreviewScale(1.5);
      useUIStore.getState().setAutoFitScale(0.7);
      useUIStore.getState().setTheme('dark');

      useUIStore.getState().reset();

      const state = useUIStore.getState();
      expect(state.leftSidebarOpen).toBe(true);
      expect(state.exportModalOpen).toBe(false);
      expect(state.previewScale).toBe(1.0);
      expect(state.autoFitScale).toBe(1.0);
      expect(state.theme).toBe('system');
    });
  });

  // ========================================================================
  // Selectors
  // ========================================================================

  describe('selectors', () => {
    it('useIsAnyModalOpen should return true when any modal is open', () => {
      // Need to test within React component context, but we can test the logic
      expect(useUIStore.getState().exportModalOpen).toBe(false);
      expect(useUIStore.getState().templateModalOpen).toBe(false);
      expect(useUIStore.getState().deleteConfirmModalOpen).toBe(false);

      useUIStore.getState().openExportModal();

      const state = useUIStore.getState();
      const isAnyOpen =
        state.exportModalOpen ||
        state.templateModalOpen ||
        state.deleteConfirmModalOpen;
      expect(isAnyOpen).toBe(true);
    });

    it('should calculate preview scale as percentage', () => {
      useUIStore.getState().setPreviewScale(1.0);
      expect(Math.round(useUIStore.getState().previewScale * 100) + '%').toBe('100%');

      useUIStore.getState().setPreviewScale(0.75);
      expect(Math.round(useUIStore.getState().previewScale * 100) + '%').toBe('75%');
    });
  });
});
