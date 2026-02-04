'use client';

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
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  Tooltip,
  TooltipContent,
  TooltipTrigger,
} from '@/components/ui/tooltip';
import { useEditorStore, useSortedSections } from '@/stores';
import { SECTION_LABELS, type SectionType } from '@/types/cv';
import { cn } from '@/lib/utils';

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
];

export function SectionPalette() {
  const { addSection, cvContent } = useEditorStore();
  const sections = useSortedSections();

  // Check which sections already exist
  const existingSectionTypes = new Set(sections.map((s) => s.type));

  const handleAddSection = (type: SectionType) => {
    addSection(type);
  };

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

        {/* Structure View */}
        <div>
          <h3 className="text-xs font-medium text-muted-foreground mb-2 px-1">
            CV Structure
          </h3>
          <div className="space-y-1">
            {sections.map((section) => {
              const Icon = sectionIcons[section.type] || Plus;
              return (
                <div
                  key={section.id}
                  className={cn(
                    'flex items-center gap-2 px-2 py-1.5 rounded text-sm',
                    section.visible !== false
                      ? 'text-foreground'
                      : 'text-muted-foreground'
                  )}
                >
                  <Icon className="h-3.5 w-3.5 shrink-0" />
                  <span className="truncate">
                    {section.title || SECTION_LABELS[section.type]}
                  </span>
                  {section.visible === false && (
                    <span className="text-xs text-muted-foreground ml-auto">
                      hidden
                    </span>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
