'use client';

import type React from 'react';
import {
  User,
  FileText,
  Briefcase,
  GraduationCap,
  Wrench,
  Languages,
  Award,
  FolderKanban,
  Trophy,
  BookOpen,
  Users,
  Plus,
  GripVertical,
} from 'lucide-react';
import {
  DndContext,
  DragEndEvent,
  DragOverlay,
  DragStartEvent,
} from '@dnd-kit/core';
import { SortableContext } from '@dnd-kit/sortable';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useEditorStore, useSortedSections } from '@/stores';
import { SECTION_LABELS, type SectionType, type CVSection } from '@/types/cv';
import { cn } from '@/lib/utils';
import { useDndSensors, collisionDetection, sortStrategy, useSortableItem } from '@/lib/dnd';
import { useState } from 'react';

const sectionIcons: Record<SectionType, React.ElementType> = {
  personal: User,
  summary: FileText,
  experience: Briefcase,
  education: GraduationCap,
  skills: Wrench,
  languages: Languages,
  certifications: Award,
  projects: FolderKanban,
  awards: Trophy,
  publications: BookOpen,
  references: Users,
  custom: Plus,
};

// MVP section types (in preferred order)
const mvpSectionTypes: SectionType[] = [
  'personal',
  'summary',
  'experience',
  'education',
  'skills',
  'languages',
  'certifications',
  'projects',
];

export function SectionPalette() {
  const { addSection, reorderSections, selectSection, selectedSectionId } = useEditorStore();
  const sections = useSortedSections();
  const sensors = useDndSensors();
  const [activeId, setActiveId] = useState<string | null>(null);

  // Check which sections already exist
  const existingSectionTypes = new Set(sections.map((s) => s.type));

  const handleAddSection = (type: SectionType) => {
    addSection(type);
  };

  const handleDragStart = (event: DragStartEvent) => {
    setActiveId(event.active.id as string);
  };

  const handleDragEnd = (event: DragEndEvent) => {
    const { active, over } = event;
    setActiveId(null);

    if (over && active.id !== over.id) {
      reorderSections(active.id as string, over.id as string);
    }
  };

  const activeSection = activeId
    ? sections.find((s) => s.id === activeId)
    : null;

  return (
    <div className="h-full flex flex-col bg-muted/30">
      {/* Header */}
      <div className="p-4 border-b">
        <h2 className="font-semibold text-sm">Add Section</h2>
        <p className="text-xs text-muted-foreground mt-1">
          Click to add a new section to your CV
        </p>
      </div>

      {/* Section Buttons */}
      <div className="flex-1 overflow-y-auto p-3">
        <div className="grid gap-2">
          {mvpSectionTypes.map((type) => {
            const Icon = sectionIcons[type];
            const label = SECTION_LABELS[type];
            const exists = existingSectionTypes.has(type);

            // Personal section can only exist once
            const isDisabled = type === 'personal' && exists;

            return (
              <Tooltip key={type}>
                <TooltipTrigger asChild>
                  <Button
                    variant="outline"
                    size="sm"
                    className={cn(
                      'justify-start h-10 px-3',
                      exists && type !== 'personal' && 'border-primary/20 bg-primary/5',
                      isDisabled && 'opacity-50 cursor-not-allowed'
                    )}
                    onClick={() => !isDisabled && handleAddSection(type)}
                    disabled={isDisabled}
                  >
                    <Icon className="h-4 w-4 mr-2 shrink-0" />
                    <span className="truncate">{label}</span>
                    {exists && type !== 'personal' && (
                      <span className="ml-auto text-xs text-muted-foreground">
                        +
                      </span>
                    )}
                  </Button>
                </TooltipTrigger>
                <TooltipContent side="right">
                  {isDisabled
                    ? 'Personal info already exists'
                    : exists
                      ? `Add another ${label} section`
                      : `Add ${label} section`}
                </TooltipContent>
              </Tooltip>
            );
          })}
        </div>

        {/* Divider */}
        <div className="my-4 border-t" />

        {/* Structure View with Drag & Drop */}
        <div>
          <h3 className="text-xs font-medium text-muted-foreground mb-2 px-1">
            CV Structure
            <span className="ml-2 text-muted-foreground/70">(drag to reorder)</span>
          </h3>
          <DndContext
            sensors={sensors}
            collisionDetection={collisionDetection}
            onDragStart={handleDragStart}
            onDragEnd={handleDragEnd}
          >
            <SortableContext items={sections.map((s) => s.id)} strategy={sortStrategy}>
              <div className="space-y-1">
                {sections.map((section) => (
                  <SortableSectionItem
                    key={section.id}
                    section={section}
                    isSelected={section.id === selectedSectionId}
                    onSelect={() => selectSection(section.id)}
                  />
                ))}
              </div>
            </SortableContext>

            {/* Drag Overlay */}
            <DragOverlay>
              {activeSection && (
                <SectionItemContent
                  section={activeSection}
                  isSelected={false}
                  isDragging
                />
              )}
            </DragOverlay>
          </DndContext>
        </div>
      </div>
    </div>
  );
}

interface SortableSectionItemProps {
  section: CVSection;
  isSelected: boolean;
  onSelect: () => void;
}

function SortableSectionItem({ section, isSelected, onSelect }: SortableSectionItemProps) {
  const { ref, style, isDragging, dragHandleProps } = useSortableItem(section.id);

  return (
    <div ref={ref} style={style}>
      <SectionItemContent
        section={section}
        isSelected={isSelected}
        isDragging={isDragging}
        dragHandleProps={dragHandleProps}
        onClick={onSelect}
      />
    </div>
  );
}

interface SectionItemContentProps {
  section: CVSection;
  isSelected: boolean;
  isDragging?: boolean;
  dragHandleProps?: Record<string, unknown>;
  onClick?: () => void;
}

function SectionItemContent({
  section,
  isSelected,
  isDragging,
  dragHandleProps,
  onClick,
}: SectionItemContentProps) {
  const Icon = sectionIcons[section.type] || Plus;

  return (
    <div
      className={cn(
        'flex items-center gap-2 px-2 py-1.5 rounded text-sm cursor-pointer',
        'hover:bg-accent transition-colors',
        isSelected && 'bg-primary/10 text-primary',
        section.visible === false && 'opacity-60',
        isDragging && 'bg-accent shadow-lg ring-2 ring-primary/20'
      )}
      onClick={onClick}
    >
      <div {...dragHandleProps} className="cursor-grab touch-none">
        <GripVertical className="h-3.5 w-3.5 text-muted-foreground" />
      </div>
      <Icon className="h-3.5 w-3.5 shrink-0" />
      <span className="truncate flex-1">
        {section.title || SECTION_LABELS[section.type]}
      </span>
      {section.visible === false && (
        <span className="text-xs text-muted-foreground">hidden</span>
      )}
    </div>
  );
}
