/**
 * Keyboard Shortcuts Hook Tests
 */

import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook } from '@testing-library/react';
import { fireEvent } from '@testing-library/react';

const mockUndo = vi.fn();
const mockRedo = vi.fn();
const mockCanUndo = vi.fn(() => true);
const mockCanRedo = vi.fn(() => true);
const mockZoomIn = vi.fn();
const mockZoomOut = vi.fn();
const mockResetZoom = vi.fn();

vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    undo: mockUndo,
    redo: mockRedo,
    canUndo: mockCanUndo,
    canRedo: mockCanRedo,
  }),
  useUIStore: () => ({
    zoomIn: mockZoomIn,
    zoomOut: mockZoomOut,
    resetZoom: mockResetZoom,
  }),
}));

import { useKeyboardShortcuts } from '../useKeyboardShortcuts';

describe('useKeyboardShortcuts', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockCanUndo.mockReturnValue(true);
    mockCanRedo.mockReturnValue(true);
  });

  it('should call undo on Ctrl+Z', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'z', ctrlKey: true });

    expect(mockUndo).toHaveBeenCalled();
  });

  it('should call redo on Ctrl+Y', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'y', ctrlKey: true });

    expect(mockRedo).toHaveBeenCalled();
  });

  it('should call redo on Ctrl+Shift+Z', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'z', ctrlKey: true, shiftKey: true });

    expect(mockRedo).toHaveBeenCalled();
  });

  it('should not undo when canUndo is false', () => {
    mockCanUndo.mockReturnValue(false);
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'z', ctrlKey: true });

    expect(mockUndo).not.toHaveBeenCalled();
  });

  it('should not redo when canRedo is false', () => {
    mockCanRedo.mockReturnValue(false);
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'y', ctrlKey: true });

    expect(mockRedo).not.toHaveBeenCalled();
  });

  it('should call zoomIn on Ctrl+=', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: '=', ctrlKey: true });

    expect(mockZoomIn).toHaveBeenCalled();
  });

  it('should call zoomOut on Ctrl+-', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: '-', ctrlKey: true });

    expect(mockZoomOut).toHaveBeenCalled();
  });

  it('should call resetZoom on Ctrl+0', () => {
    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: '0', ctrlKey: true });

    expect(mockResetZoom).toHaveBeenCalled();
  });

  it('should ignore shortcuts when target is input', () => {
    renderHook(() => useKeyboardShortcuts());

    const input = document.createElement('input');
    document.body.appendChild(input);
    fireEvent.keyDown(input, { key: 'z', ctrlKey: true });

    expect(mockUndo).not.toHaveBeenCalled();
    document.body.removeChild(input);
  });

  it('should ignore shortcuts when target is textarea', () => {
    renderHook(() => useKeyboardShortcuts());

    const textarea = document.createElement('textarea');
    document.body.appendChild(textarea);
    fireEvent.keyDown(textarea, { key: 'z', ctrlKey: true });

    expect(mockUndo).not.toHaveBeenCalled();
    document.body.removeChild(textarea);
  });

  it('should ignore shortcuts when target is contentEditable', () => {
    renderHook(() => useKeyboardShortcuts());

    const div = document.createElement('div');
    div.contentEditable = 'true';
    // jsdom doesn't compute isContentEditable from the attribute
    Object.defineProperty(div, 'isContentEditable', { value: true });
    document.body.appendChild(div);
    fireEvent.keyDown(div, { key: 'z', ctrlKey: true });

    expect(mockUndo).not.toHaveBeenCalled();
    document.body.removeChild(div);
  });

  it('should call preventDefault for matched shortcuts', () => {
    renderHook(() => useKeyboardShortcuts());

    const event = new KeyboardEvent('keydown', {
      key: 'z',
      ctrlKey: true,
      bubbles: true,
      cancelable: true,
    });
    const preventDefaultSpy = vi.spyOn(event, 'preventDefault');
    window.dispatchEvent(event);

    expect(preventDefaultSpy).toHaveBeenCalled();
  });

  it('should use metaKey on Mac platform', () => {
    const originalPlatform = navigator.platform;
    Object.defineProperty(navigator, 'platform', {
      value: 'MacIntel',
      configurable: true,
    });

    renderHook(() => useKeyboardShortcuts());
    fireEvent.keyDown(window, { key: 'z', metaKey: true });

    expect(mockUndo).toHaveBeenCalled();

    Object.defineProperty(navigator, 'platform', {
      value: originalPlatform,
      configurable: true,
    });
  });
});
