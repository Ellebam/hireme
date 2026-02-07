'use client';

import { useEffect } from 'react';
import { useUIStore } from '@/stores';
import { EditorToolbar } from './EditorToolbar';
import { SectionPalette } from './SectionPalette';
import { CVPreview } from './CVPreview';
import { PropertiesPanel } from './PropertiesPanel';
import { ExportModal, DeleteConfirmModal } from './modals';
import { cn } from '@/lib/utils';

const BREAKPOINT_MD = 768;
const BREAKPOINT_LG = 1024;

export function EditorLayout() {
  const { leftSidebarOpen, rightSidebarOpen, setLeftSidebarOpen, setRightSidebarOpen } = useUIStore();

  // Auto-collapse sidebars on small screens
  useEffect(() => {
    function handleResize() {
      const width = window.innerWidth;
      if (width < BREAKPOINT_MD) {
        setLeftSidebarOpen(false);
        setRightSidebarOpen(false);
      } else if (width < BREAKPOINT_LG) {
        setLeftSidebarOpen(false);
      }
    }

    // Check on mount
    handleResize();

    window.addEventListener('resize', handleResize);
    return () => window.removeEventListener('resize', handleResize);
    // Only run on mount
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  return (
    <div className="h-full flex flex-col">
      {/* Toolbar */}
      <EditorToolbar />

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left Sidebar - Section Palette */}
        <div
          className={cn(
            'border-r flex-shrink-0 transition-all duration-200 overflow-hidden',
            leftSidebarOpen ? 'w-60' : 'w-0 border-r-0'
          )}
        >
          <SectionPalette />
        </div>

        {/* Center - CV Preview */}
        <CVPreview />

        {/* Right Sidebar - Properties Panel */}
        <div
          className={cn(
            'border-l flex-shrink-0 transition-all duration-200 overflow-hidden',
            rightSidebarOpen ? 'w-80' : 'w-0 border-l-0'
          )}
        >
          <PropertiesPanel />
        </div>
      </div>

      {/* Modals */}
      <ExportModal />
      <DeleteConfirmModal />
    </div>
  );
}
