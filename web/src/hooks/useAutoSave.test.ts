/**
 * useAutoSave Hook Tests
 *
 * Tests the critical save lifecycle: registration, cleanup,
 * race condition prevention, and error handling.
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useEditorStore } from '@/stores/editor-store';
import { mockCV } from '@/test/mocks/cv';

// Mock the API module
const mockUpdate = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    cv: {
      update: (...args: unknown[]) => mockUpdate(...args),
    },
  },
  ApiError: class ApiError extends Error {
    constructor(
      public status: number,
      message: string
    ) {
      super(message);
    }
  },
}));

// Mock logger
vi.mock('@/lib/logger', () => ({
  logger: {
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
}));

// Import after mocks
import { useAutoSave } from './useAutoSave';

function resetStore() {
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
}

describe('useAutoSave', () => {
  beforeEach(() => {
    resetStore();
    mockUpdate.mockReset();
    vi.useFakeTimers();
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  // ==========================================================================
  // saveNow registration lifecycle
  // ==========================================================================

  describe('saveNow registration', () => {
    it('should register saveNow callback on mount', () => {
      // Load a CV so the hook has something to work with
      useEditorStore.getState().setCV(mockCV);

      expect(useEditorStore.getState().saveNow).toBeNull();

      renderHook(() => useAutoSave());

      expect(useEditorStore.getState().saveNow).toBeInstanceOf(Function);
    });

    it('should unregister saveNow callback on unmount', () => {
      useEditorStore.getState().setCV(mockCV);

      const { unmount } = renderHook(() => useAutoSave());

      expect(useEditorStore.getState().saveNow).toBeInstanceOf(Function);

      unmount();

      expect(useEditorStore.getState().saveNow).toBeNull();
    });
  });

  // ==========================================================================
  // Immediate save (saveNow)
  // ==========================================================================

  describe('saveNow (immediate save)', () => {
    it('should save immediately when saveNow is called', async () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().markDirty();
      mockUpdate.mockResolvedValueOnce(mockCV);

      renderHook(() => useAutoSave());

      const saveNow = useEditorStore.getState().saveNow;
      expect(saveNow).toBeInstanceOf(Function);

      // Call saveNow and wait for the async save
      await act(async () => {
        saveNow!();
      });

      expect(mockUpdate).toHaveBeenCalledWith(
        mockCV.id,
        { content: mockCV.content }
      );
      expect(useEditorStore.getState().saveStatus).toBe('saved');
      expect(useEditorStore.getState().isDirty).toBe(false);
    });

    it('should not save when not dirty', async () => {
      useEditorStore.getState().setCV(mockCV);
      // isDirty is false by default after setCV

      renderHook(() => useAutoSave());

      await act(async () => {
        useEditorStore.getState().saveNow!();
      });

      expect(mockUpdate).not.toHaveBeenCalled();
    });

    it('should not save when no CV is loaded', async () => {
      // Store has no CV
      renderHook(() => useAutoSave());

      const saveNow = useEditorStore.getState().saveNow;
      if (saveNow) {
        await act(async () => {
          saveNow();
        });
      }

      expect(mockUpdate).not.toHaveBeenCalled();
    });
  });

  // ==========================================================================
  // Race condition prevention
  // ==========================================================================

  describe('concurrent save prevention', () => {
    it('should block concurrent saves via saveInProgressRef', async () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().markDirty();

      // First save takes a while
      let resolveFirst!: (value: unknown) => void;
      mockUpdate.mockImplementationOnce(
        () => new Promise((resolve) => { resolveFirst = resolve; })
      );
      // Second save would resolve immediately
      mockUpdate.mockResolvedValueOnce(mockCV);

      renderHook(() => useAutoSave());

      const saveNow = useEditorStore.getState().saveNow!;

      // Start first save (will be pending)
      await act(async () => {
        saveNow();
        // Try second save while first is in progress
        saveNow();
      });

      // Resolve the first save
      await act(async () => {
        resolveFirst(mockCV);
      });

      // Only one API call should have been made
      expect(mockUpdate).toHaveBeenCalledTimes(1);
    });
  });

  // ==========================================================================
  // Error handling
  // ==========================================================================

  describe('error handling', () => {
    it('should set error state on save failure', async () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().markDirty();
      mockUpdate.mockRejectedValueOnce(new Error('Network error'));

      renderHook(() => useAutoSave());

      await act(async () => {
        useEditorStore.getState().saveNow!();
      });

      expect(useEditorStore.getState().saveStatus).toBe('error');
      expect(useEditorStore.getState().saveError).toBe('Network error');
    });

    it('should allow retry after failure', async () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().markDirty();

      // First attempt fails
      mockUpdate.mockRejectedValueOnce(new Error('Temporary error'));

      renderHook(() => useAutoSave());

      await act(async () => {
        useEditorStore.getState().saveNow!();
      });

      expect(useEditorStore.getState().saveStatus).toBe('error');

      // Mark dirty again for retry
      useEditorStore.getState().markDirty();

      // Second attempt succeeds
      mockUpdate.mockResolvedValueOnce(mockCV);

      await act(async () => {
        useEditorStore.getState().saveNow!();
      });

      expect(useEditorStore.getState().saveStatus).toBe('saved');
      expect(useEditorStore.getState().saveError).toBeNull();
    });

    it('should handle non-Error thrown values', async () => {
      useEditorStore.getState().setCV(mockCV);
      useEditorStore.getState().markDirty();
      mockUpdate.mockRejectedValueOnce('string error');

      renderHook(() => useAutoSave());

      await act(async () => {
        useEditorStore.getState().saveNow!();
      });

      expect(useEditorStore.getState().saveStatus).toBe('error');
      expect(useEditorStore.getState().saveError).toBe('Failed to save');
    });
  });

  // ==========================================================================
  // Debounced auto-save
  // ==========================================================================

  describe('debounced auto-save', () => {
    it('should trigger save after debounce delay when dirty', async () => {
      mockUpdate.mockResolvedValue(mockCV);

      renderHook(() => useAutoSave());

      // Load CV and make dirty
      act(() => {
        useEditorStore.getState().setCV(mockCV);
        useEditorStore.getState().markDirty();
      });

      // Before debounce fires
      expect(mockUpdate).not.toHaveBeenCalled();

      // Advance past debounce delay (2000ms)
      await act(async () => {
        vi.advanceTimersByTime(2500);
      });

      expect(mockUpdate).toHaveBeenCalledTimes(1);
    });
  });
});
