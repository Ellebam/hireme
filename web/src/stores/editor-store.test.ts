/**
 * Editor Store Tests
 */

import { describe, it, expect, beforeEach, vi } from 'vitest';
import { useEditorStore } from './editor-store';
import { mockCV, minimalCV, createMockCV } from '@/test/mocks/cv';
import type { CVContent } from '@/types/cv';

describe('EditorStore', () => {
  // Reset store before each test
  beforeEach(() => {
    useEditorStore.setState({
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
    });
  });

  // ========================================================================
  // Initialization
  // ========================================================================

  describe('setCV', () => {
    it('should set CV and parse content', () => {
      useEditorStore.getState().setCV(mockCV);

      const state = useEditorStore.getState();
      expect(state.cv).toEqual(mockCV);
      expect(state.cvContent).toEqual(mockCV.content);
      expect(state.isDirty).toBe(false);
      expect(state.saveStatus).toBe('idle');
    });

    it('should initialize history with current content', () => {
      useEditorStore.getState().setCV(mockCV);

      const state = useEditorStore.getState();
      expect(state.history).toHaveLength(1);
      expect(state.historyIndex).toBe(0);
    });

    it('should clear previous selection', () => {
      // Set some selection first
      useEditorStore.setState({ selectedSectionId: 'old-section' });

      useEditorStore.getState().setCV(mockCV);

      const state = useEditorStore.getState();
      expect(state.selectedSectionId).toBeNull();
      expect(state.selectedEntryId).toBeNull();
    });
  });

  describe('reset', () => {
    it('should reset all state to initial values', () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().selectSection('sec-personal');

      useEditorStore.getState().reset();

      const state = useEditorStore.getState();
      expect(state.cv).toBeNull();
      expect(state.cvContent).toBeNull();
      expect(state.selectedSectionId).toBeNull();
      expect(state.isDirty).toBe(false);
    });
  });

  // ========================================================================
  // Section Operations
  // ========================================================================

  describe('updateSection', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should update section and mark dirty', () => {
      useEditorStore.getState().updateSection('sec-personal', { visible: false });

      const state = useEditorStore.getState();
      const section = state.cvContent?.sections.find((s) => s.id === 'sec-personal');
      expect(section?.visible).toBe(false);
      expect(state.isDirty).toBe(true);
    });

    it('should push to history before updating', () => {
      const initialHistoryLength = useEditorStore.getState().history.length;

      useEditorStore.getState().updateSection('sec-personal', { title: 'New Title' });

      expect(useEditorStore.getState().history.length).toBe(initialHistoryLength + 1);
    });

    it('should not modify other sections', () => {
      const originalSummary = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.type === 'summary'
      );

      useEditorStore.getState().updateSection('sec-personal', { title: 'Changed' });

      const newSummary = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.type === 'summary'
      );
      expect(newSummary).toEqual(originalSummary);
    });
  });

  describe('addSection', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(minimalCV);
    });

    it('should add new section at the end', () => {
      const newId = useEditorStore.getState().addSection('experience');

      const state = useEditorStore.getState();
      const sections = state.cvContent?.sections || [];
      expect(sections).toHaveLength(2);
      expect(sections[1].id).toBe(newId);
      expect(sections[1].type).toBe('experience');
    });

    it('should add section after specified section', () => {
      // Add second section
      useEditorStore.getState().addSection('summary');

      // Add third section after first
      const newId = useEditorStore.getState().addSection('skills', 'sec-1');

      const sections = useEditorStore.getState().cvContent?.sections || [];
      expect(sections[1].id).toBe(newId);
      expect(sections[1].type).toBe('skills');
    });

    it('should select the new section', () => {
      const newId = useEditorStore.getState().addSection('education');

      expect(useEditorStore.getState().selectedSectionId).toBe(newId);
    });

    it('should return the new section ID', () => {
      const newId = useEditorStore.getState().addSection('languages');

      expect(newId).toBeTruthy();
      expect(typeof newId).toBe('string');
    });

    it('should recalculate order for all sections', () => {
      useEditorStore.getState().addSection('summary');
      useEditorStore.getState().addSection('skills');

      const sections = useEditorStore.getState().cvContent?.sections || [];
      sections.forEach((section, index) => {
        expect(section.order).toBe(index);
      });
    });
  });

  describe('deleteSection', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should remove section', () => {
      const initialLength = useEditorStore.getState().cvContent?.sections.length || 0;

      useEditorStore.getState().deleteSection('sec-languages');

      const sections = useEditorStore.getState().cvContent?.sections || [];
      expect(sections).toHaveLength(initialLength - 1);
      expect(sections.find((s) => s.id === 'sec-languages')).toBeUndefined();
    });

    it('should clear selection if deleted section was selected', () => {
      useEditorStore.getState().selectSection('sec-languages');
      expect(useEditorStore.getState().selectedSectionId).toBe('sec-languages');

      useEditorStore.getState().deleteSection('sec-languages');

      expect(useEditorStore.getState().selectedSectionId).toBeNull();
    });

    it('should recalculate order after deletion', () => {
      useEditorStore.getState().deleteSection('sec-summary');

      const sections = useEditorStore.getState().cvContent?.sections || [];
      sections.forEach((section, index) => {
        expect(section.order).toBe(index);
      });
    });
  });

  describe('reorderSections', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should move section from one position to another', () => {
      // Move languages (last) to position of summary (second)
      useEditorStore.getState().reorderSections('sec-languages', 'sec-summary');

      const sections = useEditorStore.getState().cvContent?.sections || [];
      expect(sections[1].id).toBe('sec-languages');
    });

    it('should not change anything if same ID', () => {
      const originalSections = [...(useEditorStore.getState().cvContent?.sections || [])];

      useEditorStore.getState().reorderSections('sec-personal', 'sec-personal');

      const newSections = useEditorStore.getState().cvContent?.sections || [];
      expect(newSections.map((s) => s.id)).toEqual(originalSections.map((s) => s.id));
    });

    it('should mark as dirty', () => {
      useEditorStore.getState().reorderSections('sec-skills', 'sec-personal');

      expect(useEditorStore.getState().isDirty).toBe(true);
    });
  });

  describe('toggleSectionVisibility', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should toggle visibility from true to false', () => {
      const section = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.id === 'sec-personal'
      );
      expect(section?.visible).toBe(true);

      useEditorStore.getState().toggleSectionVisibility('sec-personal');

      const updated = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.id === 'sec-personal'
      );
      expect(updated?.visible).toBe(false);
    });

    it('should toggle visibility from false to true', () => {
      useEditorStore.getState().updateSection('sec-personal', { visible: false });

      useEditorStore.getState().toggleSectionVisibility('sec-personal');

      const section = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.id === 'sec-personal'
      );
      expect(section?.visible).toBe(true);
    });
  });

  // ========================================================================
  // Entry Operations
  // ========================================================================

  describe('reorderEntries', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should reorder entries within a section', () => {
      const section = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.id === 'sec-experience'
      );
      const content = section?.content as { entries: { id: string }[] };
      const originalOrder = content.entries.map((e) => e.id);

      useEditorStore.getState().reorderEntries('sec-experience', 'exp-2', 'exp-1');

      const updated = useEditorStore.getState().cvContent?.sections.find(
        (s) => s.id === 'sec-experience'
      );
      const updatedContent = updated?.content as { entries: { id: string }[] };
      expect(updatedContent.entries[0].id).toBe('exp-2');
      expect(updatedContent.entries[1].id).toBe('exp-1');
    });
  });

  // ========================================================================
  // Selection
  // ========================================================================

  describe('selectSection', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should set selectedSectionId', () => {
      useEditorStore.getState().selectSection('sec-experience');

      expect(useEditorStore.getState().selectedSectionId).toBe('sec-experience');
    });

    it('should clear entry selection when selecting section', () => {
      useEditorStore.setState({ selectedEntryId: 'some-entry' });

      useEditorStore.getState().selectSection('sec-education');

      expect(useEditorStore.getState().selectedEntryId).toBeNull();
    });

    it('should allow null to clear selection', () => {
      useEditorStore.getState().selectSection('sec-skills');
      useEditorStore.getState().selectSection(null);

      expect(useEditorStore.getState().selectedSectionId).toBeNull();
    });
  });

  // ========================================================================
  // History (Undo/Redo)
  // ========================================================================

  describe('undo/redo', () => {
    beforeEach(() => {
      useEditorStore.getState().setCV(mockCV);
    });

    it('should undo last change', () => {
      const originalTitle = useEditorStore.getState().cvContent?.sections[0].title;

      useEditorStore.getState().updateSection('sec-personal', { title: 'Changed' });
      expect(useEditorStore.getState().cvContent?.sections[0].title).toBe('Changed');

      useEditorStore.getState().undo();

      expect(useEditorStore.getState().cvContent?.sections[0].title).toBe(originalTitle);
    });

    it('should redo undone change', () => {
      // First check visibility is true
      expect(useEditorStore.getState().cvContent?.sections[0].visible).toBe(true);

      // Change visibility to false
      useEditorStore.getState().updateSection('sec-personal', { visible: false });
      expect(useEditorStore.getState().cvContent?.sections[0].visible).toBe(false);

      // Undo - should be true again
      useEditorStore.getState().undo();
      expect(useEditorStore.getState().cvContent?.sections[0].visible).toBe(true);

      // Redo - should be false again
      useEditorStore.getState().redo();
      expect(useEditorStore.getState().cvContent?.sections[0].visible).toBe(false);
    });

    it('canUndo should return false at beginning', () => {
      expect(useEditorStore.getState().canUndo()).toBe(false);
    });

    it('canUndo should return true after change', () => {
      useEditorStore.getState().updateSection('sec-personal', { title: 'Changed' });

      expect(useEditorStore.getState().canUndo()).toBe(true);
    });

    it('canRedo should return false when at latest', () => {
      useEditorStore.getState().updateSection('sec-personal', { title: 'Changed' });

      expect(useEditorStore.getState().canRedo()).toBe(false);
    });

    it('canRedo should return true after undo', () => {
      useEditorStore.getState().updateSection('sec-personal', { title: 'Changed' });
      useEditorStore.getState().undo();

      expect(useEditorStore.getState().canRedo()).toBe(true);
    });
  });

  // ========================================================================
  // Save State
  // ========================================================================

  describe('save state management', () => {
    it('should track save status transitions', () => {
      useEditorStore.getState().setCV(mockCV);

      expect(useEditorStore.getState().saveStatus).toBe('idle');

      useEditorStore.getState().markSaving();
      expect(useEditorStore.getState().saveStatus).toBe('saving');

      useEditorStore.getState().markSaved();
      expect(useEditorStore.getState().saveStatus).toBe('saved');
      expect(useEditorStore.getState().isDirty).toBe(false);
      expect(useEditorStore.getState().lastSavedAt).toBeInstanceOf(Date);
    });

    it('should handle save errors', () => {
      useEditorStore.getState().markError('Network error');

      expect(useEditorStore.getState().saveStatus).toBe('error');
      expect(useEditorStore.getState().saveError).toBe('Network error');
    });

    it('markSaved should clear error', () => {
      useEditorStore.getState().markError('Some error');
      useEditorStore.getState().markSaved();

      expect(useEditorStore.getState().saveError).toBeNull();
    });
  });

  // ========================================================================
  // setSaveNow
  // ========================================================================

  describe('setSaveNow', () => {
    it('should store a callback function', () => {
      const mockSave = vi.fn();
      useEditorStore.getState().setSaveNow(mockSave);

      expect(useEditorStore.getState().saveNow).toBe(mockSave);
    });

    it('should allow clearing with null', () => {
      useEditorStore.getState().setSaveNow(vi.fn());
      useEditorStore.getState().setSaveNow(null);

      expect(useEditorStore.getState().saveNow).toBeNull();
    });

    it('should replace existing callback', () => {
      const first = vi.fn();
      const second = vi.fn();

      useEditorStore.getState().setSaveNow(first);
      useEditorStore.getState().setSaveNow(second);

      expect(useEditorStore.getState().saveNow).toBe(second);
    });
  });

  // ========================================================================
  // Styling & Template
  // ========================================================================

  describe('updateStyling', () => {
    it('should update styling properties', () => {
      useEditorStore.getState().setCV(mockCV);

      useEditorStore.getState().updateStyling({ primaryColor: '#ff0000' });

      const state = useEditorStore.getState();
      expect(state.cvContent?.styling?.primaryColor).toBe('#ff0000');
      expect(state.isDirty).toBe(true);
    });

    it('should merge with existing styling', () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().updateStyling({ primaryColor: '#ff0000' });
      useEditorStore.getState().updateStyling({ secondaryColor: '#00ff00' });

      const styling = useEditorStore.getState().cvContent?.styling;
      expect(styling?.primaryColor).toBe('#ff0000');
      expect(styling?.secondaryColor).toBe('#00ff00');
    });

    it('should do nothing when no CV is loaded', () => {
      useEditorStore.getState().updateStyling({ primaryColor: '#ff0000' });
      expect(useEditorStore.getState().cvContent).toBeNull();
    });

    it('should push to history', () => {
      useEditorStore.getState().setCV(mockCV);
      const historyLenBefore = useEditorStore.getState().history.length;

      useEditorStore.getState().updateStyling({ primaryColor: '#ff0000' });

      expect(useEditorStore.getState().history.length).toBe(historyLenBefore + 1);
    });
  });

  describe('updateTemplateId', () => {
    it('should update templateId', () => {
      useEditorStore.getState().setCV(mockCV);

      useEditorStore.getState().updateTemplateId('visionary');

      expect(useEditorStore.getState().cvContent?.templateId).toBe('visionary');
      expect(useEditorStore.getState().isDirty).toBe(true);
    });

    it('should do nothing when no CV is loaded', () => {
      useEditorStore.getState().updateTemplateId('classic');
      expect(useEditorStore.getState().cvContent).toBeNull();
    });

    it('should push to history', () => {
      useEditorStore.getState().setCV(mockCV);
      const historyLenBefore = useEditorStore.getState().history.length;

      useEditorStore.getState().updateTemplateId('modern');

      expect(useEditorStore.getState().history.length).toBe(historyLenBefore + 1);
    });
  });
});
