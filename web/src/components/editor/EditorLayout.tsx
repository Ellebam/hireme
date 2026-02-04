'use client';

import { useUIStore } from '@/stores';
import { EditorToolbar } from './EditorToolbar';
import { SectionPalette } from './SectionPalette';
import { CVPreview } from './CVPreview';
import { PropertiesPanel } from './PropertiesPanel';
import { ExportModal, DeleteConfirmModal } from './modals';
import { cn } from '@/lib/utils';

export function EditorLayout() {
  const { leftSidebarOpen, rightSidebarOpen } = useUIStore();

  return (
    <div className="h-full flex flex-col">
      {/* Toolbar */}
      <EditorToolbar />

      {/* Main Content */}
      <div className="flex-1 flex overflow-hidden">
        {/* Left Sidebar - Section Palette */}
        <div
          className={cn(
            'w-60 border-r flex-shrink-0 transition-all duration-200 overflow-hidden',
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
            'w-80 border-l flex-shrink-0 transition-all duration-200 overflow-hidden',
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
