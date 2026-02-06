'use client';

import { useVisibleSections, useEditorStore } from '@/stores';
import { useUIStore } from '@/stores';
import { SECTION_LABELS, type CVSection } from '@/types/cv';
import { cn } from '@/lib/utils';

// A4 dimensions at 96 DPI
const A4_WIDTH = 794;
const A4_HEIGHT = 1123;

export function CVPreview() {
  const sections = useVisibleSections();
  const { cvContent, selectedSectionId, selectSection } = useEditorStore();
  const { previewScale } = useUIStore();

  if (!cvContent) {
    return (
      <div className="flex-1 flex items-center justify-center text-muted-foreground">
        Loading...
      </div>
    );
  }

  return (
    <div className="flex-1 overflow-auto bg-muted/50 p-8">
      <div className="flex justify-center">
        <div
          style={{
            width: A4_WIDTH * previewScale,
            minHeight: A4_HEIGHT * previewScale,
          }}
        >
          <div
            className="bg-white shadow-lg rounded"
            style={{
              width: A4_WIDTH,
              minHeight: A4_HEIGHT,
              transform: `scale(${previewScale})`,
              transformOrigin: 'top center',
            }}
          >
            {/* CV Content */}
            <div className="p-12">
              {sections.length === 0 ? (
                <EmptyPreview />
              ) : (
                <div className="space-y-6">
                  {sections.map((section) => (
                    <PreviewSection
                      key={section.id}
                      section={section}
                      isSelected={section.id === selectedSectionId}
                      onClick={() => selectSection(section.id)}
                    />
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function EmptyPreview() {
  return (
    <div className="flex flex-col items-center justify-center min-h-[400px] text-center">
      <div className="text-6xl mb-4 opacity-20">📄</div>
      <h3 className="text-lg font-medium text-muted-foreground">
        Your CV is empty
      </h3>
      <p className="text-sm text-muted-foreground mt-2 max-w-xs">
        Add sections from the palette on the left to start building your CV
      </p>
    </div>
  );
}

interface PreviewSectionProps {
  section: CVSection;
  isSelected: boolean;
  onClick: () => void;
}

function PreviewSection({ section, isSelected, onClick }: PreviewSectionProps) {
  return (
    <div
      className={cn(
        'relative p-4 -mx-4 rounded-lg cursor-pointer transition-all',
        'hover:bg-primary/5',
        isSelected && 'bg-primary/10 ring-2 ring-primary/30'
      )}
      onClick={onClick}
    >
      {/* Section Title */}
      {section.type !== 'personal' && (
        <h2 className="text-lg font-semibold text-primary mb-3 uppercase tracking-wide">
          {section.title || SECTION_LABELS[section.type]}
        </h2>
      )}

      {/* Section Content Renderer */}
      <SectionContentRenderer section={section} />

      {/* Selection Indicator */}
      {isSelected && (
        <div className="absolute top-2 right-2 px-2 py-1 bg-primary text-primary-foreground text-xs rounded">
          Click to edit
        </div>
      )}
    </div>
  );
}

function SectionContentRenderer({ section }: { section: CVSection }) {
  switch (section.type) {
    case 'personal':
      return <PersonalPreview content={section.content} />;
    case 'summary':
      return <SummaryPreview content={section.content} />;
    case 'experience':
      return <ExperiencePreview content={section.content} />;
    case 'education':
      return <EducationPreview content={section.content} />;
    case 'skills':
      return <SkillsPreview content={section.content} />;
    case 'languages':
      return <LanguagesPreview content={section.content} />;
    default:
      return (
        <div className="text-muted-foreground text-sm italic">
          {SECTION_LABELS[section.type]} content
        </div>
      );
  }
}

function PersonalPreview({ content }: { content: unknown }) {
  const data = content as {
    firstName?: string;
    lastName?: string;
    jobTitle?: string;
    email?: string;
    phone?: string;
    location?: string;
  };

  const fullName = [data.firstName, data.lastName].filter(Boolean).join(' ');

  return (
    <div className="text-center pb-6 border-b">
      <h1 className="text-3xl font-bold text-foreground">
        {fullName || 'Your Name'}
      </h1>
      {data.jobTitle && (
        <p className="text-lg text-primary mt-1">{data.jobTitle}</p>
      )}
      <div className="flex items-center justify-center gap-4 mt-3 text-sm text-muted-foreground flex-wrap">
        {data.email && <span>{data.email}</span>}
        {data.phone && <span>{data.phone}</span>}
        {data.location && <span>{data.location}</span>}
      </div>
    </div>
  );
}

function SummaryPreview({ content }: { content: unknown }) {
  const data = content as { text?: string };
  return (
    <p className="text-sm leading-relaxed text-muted-foreground">
      {data.text || 'Add your professional summary here...'}
    </p>
  );
}

function ExperiencePreview({ content }: { content: unknown }) {
  const data = content as {
    entries?: Array<{
      id: string;
      company?: string;
      position?: string;
      location?: string;
      startDate?: string;
      endDate?: string | null;
      current?: boolean;
      description?: string;
      highlights?: string[];
    }>;
  };

  if (!data.entries?.length) {
    return (
      <p className="text-sm text-muted-foreground italic">
        Add your work experience...
      </p>
    );
  }

  return (
    <div className="space-y-4">
      {data.entries.map((entry) => (
        <div key={entry.id}>
          <div className="flex items-start justify-between">
            <div>
              <h3 className="font-semibold">{entry.position || 'Position'}</h3>
              <p className="text-sm text-muted-foreground">
                {entry.company}
                {entry.location && ` · ${entry.location}`}
              </p>
            </div>
            <div className="text-sm text-muted-foreground text-right">
              {entry.startDate}
              {' - '}
              {entry.current ? 'Present' : entry.endDate}
            </div>
          </div>
          {entry.description && (
            <p className="text-sm mt-2">{entry.description}</p>
          )}
          {entry.highlights && entry.highlights.length > 0 && (
            <ul className="mt-2 space-y-1">
              {entry.highlights.map((h, i) => (
                <li key={i} className="text-sm text-muted-foreground flex">
                  <span className="mr-2">•</span>
                  {h}
                </li>
              ))}
            </ul>
          )}
        </div>
      ))}
    </div>
  );
}

function EducationPreview({ content }: { content: unknown }) {
  const data = content as {
    entries?: Array<{
      id: string;
      institution?: string;
      degree?: string;
      field?: string;
      startDate?: string;
      endDate?: string | null;
      grade?: string;
    }>;
  };

  if (!data.entries?.length) {
    return (
      <p className="text-sm text-muted-foreground italic">
        Add your education...
      </p>
    );
  }

  return (
    <div className="space-y-3">
      {data.entries.map((entry) => (
        <div key={entry.id} className="flex items-start justify-between">
          <div>
            <h3 className="font-semibold">
              {entry.degree}
              {entry.field && ` in ${entry.field}`}
            </h3>
            <p className="text-sm text-muted-foreground">{entry.institution}</p>
          </div>
          <div className="text-sm text-muted-foreground text-right">
            {entry.startDate && `${entry.startDate} - `}
            {entry.endDate || 'Present'}
            {entry.grade && <div className="text-xs">{entry.grade}</div>}
          </div>
        </div>
      ))}
    </div>
  );
}

function SkillsPreview({ content }: { content: unknown }) {
  const data = content as {
    categories?: Array<{
      id: string;
      name?: string;
      skills?: Array<{ name: string; level?: string }>;
    }>;
  };

  if (!data.categories?.length) {
    return (
      <p className="text-sm text-muted-foreground italic">Add your skills...</p>
    );
  }

  return (
    <div className="space-y-3">
      {data.categories.map((cat) => (
        <div key={cat.id}>
          <h3 className="font-medium text-sm mb-1">{cat.name}</h3>
          <div className="flex flex-wrap gap-2">
            {cat.skills?.map((skill, i) => (
              <span
                key={i}
                className="px-2 py-1 bg-muted rounded text-xs"
              >
                {skill.name}
              </span>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}

function LanguagesPreview({ content }: { content: unknown }) {
  const data = content as {
    entries?: Array<{
      language?: string;
      proficiency?: string;
    }>;
  };

  if (!data.entries?.length) {
    return (
      <p className="text-sm text-muted-foreground italic">
        Add your languages...
      </p>
    );
  }

  return (
    <div className="flex flex-wrap gap-4">
      {data.entries.map((entry, i) => (
        <div key={i} className="text-sm">
          <span className="font-medium">{entry.language}</span>
          <span className="text-muted-foreground ml-2 capitalize">
            ({entry.proficiency})
          </span>
        </div>
      ))}
    </div>
  );
}
