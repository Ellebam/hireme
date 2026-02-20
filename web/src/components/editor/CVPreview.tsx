'use client';

import { useCallback } from 'react';
import { useVisibleSections, useEditorStore } from '@/stores';
import { useUIStore } from '@/stores';
import { ClassicTemplate, ModernTemplate, VisionaryTemplate } from '@/components/templates';

// A4 dimensions at 96 DPI
const A4_WIDTH = 794;
const A4_HEIGHT = 1123;

export function CVPreview() {
  const sections = useVisibleSections();
  const { cvContent, selectedSectionId, selectSection } = useEditorStore();
  const { previewScale, setRightSidebarOpen } = useUIStore();

  const handleSectionClick = useCallback(
    (id: string) => {
      selectSection(id);
    },
    [selectSection]
  );

  const handleSectionDoubleClick = useCallback(
    (id: string) => {
      selectSection(id);
      setRightSidebarOpen(true);
    },
    [selectSection, setRightSidebarOpen]
  );

  if (!cvContent) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted-foreground">
        Loading...
      </div>
    );
  }

  const templateId = cvContent.templateId || 'classic';

  const templateProps = {
    sections,
    styling: cvContent.styling,
    selectedSectionId,
    onSectionClick: handleSectionClick,
    onSectionDoubleClick: handleSectionDoubleClick,
  };

  return (
    <div className="flex-1 overflow-auto bg-secondary p-8">
      <div className="flex justify-center">
        <div
          style={{
            width: A4_WIDTH * previewScale,
            minHeight: A4_HEIGHT * previewScale,
          }}
        >
          <div
            className="bg-card border-2 border-border shadow-offset-lg animate-paper-drop overflow-hidden"
            style={{
              width: A4_WIDTH,
              minHeight: A4_HEIGHT,
              transform: `scale(${previewScale})`,
              transformOrigin: 'top center',
            }}
          >
            {sections.length === 0 ? (
              <EmptyPreview />
            ) : templateId === 'modern' ? (
              <ModernTemplate {...templateProps} />
            ) : templateId === 'visionary' ? (
              <VisionaryTemplate {...templateProps} />
            ) : (
              <ClassicTemplate {...templateProps} />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function EmptyPreview() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] text-center p-12">
      <h3 className="font-serif text-lg font-semibold text-muted-foreground">
        Your CV is empty
      </h3>
      <p className="text-sm text-[hsl(var(--text-secondary))] mt-2 max-w-xs">
        Add sections from the palette on the left to start building your CV
      </p>
    </div>
  );
}
