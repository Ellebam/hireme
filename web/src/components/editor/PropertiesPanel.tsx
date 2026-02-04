'use client';

import { Trash2, Eye, EyeOff } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Separator } from '@/components/ui/separator';
import { useEditorStore, useSelectedSection } from '@/stores';
import { useUIStore } from '@/stores';
import { SECTION_LABELS } from '@/types/cv';
import { cn } from '@/lib/utils';

export function PropertiesPanel() {
  const section = useSelectedSection();
  const { toggleSectionVisibility, deleteSection, selectSection } = useEditorStore();
  const { openDeleteConfirmModal } = useUIStore();

  if (!section) {
    return (
      <div className="h-full flex flex-col bg-muted/30">
        <div className="p-4 border-b">
          <h2 className="font-semibold text-sm">Properties</h2>
        </div>
        <div className="flex-1 flex items-center justify-center p-4">
          <p className="text-sm text-muted-foreground text-center">
            Select a section from the preview to edit its properties
          </p>
        </div>
      </div>
    );
  }

  const handleDelete = () => {
    openDeleteConfirmModal(section.id);
  };

  const handleToggleVisibility = () => {
    toggleSectionVisibility(section.id);
  };

  return (
    <div className="h-full flex flex-col bg-muted/30">
      {/* Header */}
      <div className="p-4 border-b">
        <div className="flex items-center justify-between">
          <h2 className="font-semibold text-sm">
            {section.title || SECTION_LABELS[section.type]}
          </h2>
          <div className="flex items-center gap-1">
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7"
              onClick={handleToggleVisibility}
            >
              {section.visible !== false ? (
                <Eye className="h-4 w-4" />
              ) : (
                <EyeOff className="h-4 w-4 text-muted-foreground" />
              )}
            </Button>
            <Button
              variant="ghost"
              size="icon"
              className="h-7 w-7 text-destructive hover:text-destructive"
              onClick={handleDelete}
            >
              <Trash2 className="h-4 w-4" />
            </Button>
          </div>
        </div>
        <p className="text-xs text-muted-foreground mt-1">
          {section.visible !== false
            ? 'This section will appear in exports'
            : 'This section is hidden from exports'}
        </p>
      </div>

      {/* Content Editor */}
      <div className="flex-1 overflow-y-auto p-4">
        <SectionEditor section={section} />
      </div>
    </div>
  );
}

import type { CVSection } from '@/types/cv';

function SectionEditor({ section }: { section: CVSection }) {
  const { updateSectionContent } = useEditorStore();

  // For MVP, show a simple form based on section type
  // This will be expanded with proper form fields later

  return (
    <div className="space-y-4">
      <div className="p-4 rounded-lg bg-muted/50 border border-dashed">
        <p className="text-sm text-muted-foreground text-center">
          Edit {SECTION_LABELS[section.type]} content
        </p>
        <p className="text-xs text-muted-foreground text-center mt-2">
          Full editing forms will be implemented in Phase 4
        </p>
      </div>

      {/* Quick preview of current content */}
      <div>
        <h3 className="text-xs font-medium text-muted-foreground mb-2">
          Current Content
        </h3>
        <pre className="text-xs bg-muted p-3 rounded overflow-auto max-h-[300px]">
          {JSON.stringify(section.content, null, 2)}
        </pre>
      </div>
    </div>
  );
}
