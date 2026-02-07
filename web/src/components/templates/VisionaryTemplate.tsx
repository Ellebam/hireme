'use client';

import { cn } from '@/lib/utils';
import type {
  CVSection,
  CVStyling,
  PersonalContent,
  SummaryContent,
  ExperienceContent,
  EducationContent,
  SkillsContent,
  LanguagesContent,
} from '@/types/cv';
import { SECTION_LABELS } from '@/types/cv';

interface TemplateProps {
  sections: CVSection[];
  styling?: CVStyling;
  selectedSectionId: string | null;
  onSectionClick: (id: string) => void;
  onSectionDoubleClick: (id: string) => void;
}

const SIDEBAR_TYPES = new Set(['personal', 'skills', 'languages']);

function formatDate(dateStr: string): string {
  if (!dateStr) return '';
  const [year, month] = dateStr.split('-');
  if (!year) return '';
  if (!month) return year;
  const date = new Date(Number(year), Number(month) - 1);
  return date.toLocaleDateString('en', { year: 'numeric', month: 'short' });
}

export function VisionaryTemplate({
  sections,
  styling,
  selectedSectionId,
  onSectionClick,
  onSectionDoubleClick,
}: TemplateProps) {
  const primaryColor = styling?.primaryColor || '#2563eb';
  const secondaryColor = styling?.secondaryColor || '#64748b';

  const visibleSections = sections.filter((s) => s.visible !== false);

  const sidebarSections = visibleSections
    .filter((s) => SIDEBAR_TYPES.has(s.type))
    .sort((a, b) => a.order - b.order);

  const mainSections = visibleSections
    .filter((s) => !SIDEBAR_TYPES.has(s.type))
    .sort((a, b) => a.order - b.order);

  // Extract personal info for the header
  const personalSection = visibleSections.find((s) => s.type === 'personal');
  const personal = personalSection?.content as PersonalContent | undefined;

  return (
    <div className="flex min-h-full">
      {/* Sidebar */}
      <div
        className="w-[250px] flex-shrink-0"
        style={{ backgroundColor: primaryColor }}
      >
        <div className="p-8 text-white">
          {sidebarSections.map((section) => {
            const isSelected = selectedSectionId === section.id;

            return (
              <div
                key={section.id}
                className={cn(
                  'relative p-3 -mx-3 rounded-lg cursor-pointer transition-all',
                  'hover:opacity-90',
                  isSelected && 'ring-2 ring-white/50'
                )}
                onClick={() => onSectionClick(section.id)}
                onDoubleClick={() => onSectionDoubleClick(section.id)}
              >
                {section.type === 'personal' && (
                  <SidebarPersonal
                    content={section.content as PersonalContent}
                  />
                )}
                {section.type === 'skills' && (
                  <SidebarSkills
                    content={section.content as SkillsContent}
                    title={section.title || SECTION_LABELS[section.type]}
                  />
                )}
                {section.type === 'languages' && (
                  <SidebarLanguages
                    content={section.content as LanguagesContent}
                    title={section.title || SECTION_LABELS[section.type]}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Main content */}
      <div className="flex-1 p-8">
        {/* Header strip with name and title */}
        {personal && (personal.firstName || personal.lastName) && (
          <div className="mb-8 pb-4 border-b-2" style={{ borderColor: primaryColor }}>
            <h1
              className="text-3xl font-bold tracking-tight"
              style={{ color: primaryColor }}
            >
              {[personal.firstName, personal.lastName].filter(Boolean).join(' ')}
            </h1>
            {personal.jobTitle && (
              <p
                className="text-lg mt-1"
                style={{ color: secondaryColor }}
              >
                {personal.jobTitle}
              </p>
            )}
          </div>
        )}

        {mainSections.map((section) => {
          const isSelected = selectedSectionId === section.id;

          return (
            <div
              key={section.id}
              className={cn(
                'relative p-3 -mx-3 rounded-lg cursor-pointer transition-all',
                'hover:opacity-90',
                isSelected && 'ring-2 ring-primary/30'
              )}
              onClick={() => onSectionClick(section.id)}
              onDoubleClick={() => onSectionDoubleClick(section.id)}
            >
              {isSelected && (
                <div className="absolute top-2 right-2 px-2 py-1 bg-primary text-primary-foreground text-xs rounded">
                  Click to edit
                </div>
              )}

              {section.type === 'summary' && (
                <MainSummary
                  content={section.content as SummaryContent}
                  title={section.title || SECTION_LABELS[section.type]}
                  primaryColor={primaryColor}
                  secondaryColor={secondaryColor}
                />
              )}
              {section.type === 'experience' && (
                <MainExperience
                  content={section.content as ExperienceContent}
                  title={section.title || SECTION_LABELS[section.type]}
                  primaryColor={primaryColor}
                  secondaryColor={secondaryColor}
                />
              )}
              {section.type === 'education' && (
                <MainEducation
                  content={section.content as EducationContent}
                  title={section.title || SECTION_LABELS[section.type]}
                  primaryColor={primaryColor}
                  secondaryColor={secondaryColor}
                />
              )}
              {/* Fallback for any other section type in main */}
              {!['summary', 'experience', 'education'].includes(section.type) && (
                <MainGeneric
                  title={section.title || SECTION_LABELS[section.type] || 'Section'}
                  primaryColor={primaryColor}
                />
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

// =============================================================================
// Sidebar Components
// =============================================================================

function SidebarPersonal({ content }: { content: PersonalContent }) {
  const fullName = [content.firstName, content.lastName]
    .filter(Boolean)
    .join(' ');

  return (
    <div className="mb-6">
      {fullName && (
        <h2 className="text-2xl font-bold mb-1">{fullName}</h2>
      )}
      {content.jobTitle && (
        <p className="text-white/80 text-sm mb-4">{content.jobTitle}</p>
      )}

      <div className="space-y-2 text-sm">
        {content.email && (
          <div className="flex items-start gap-2">
            <span className="text-white/60 flex-shrink-0">Email</span>
            <span className="text-white/90 break-all">{content.email}</span>
          </div>
        )}
        {content.phone && (
          <div className="flex items-start gap-2">
            <span className="text-white/60 flex-shrink-0">Phone</span>
            <span className="text-white/90">{content.phone}</span>
          </div>
        )}
        {content.location && (
          <div className="flex items-start gap-2">
            <span className="text-white/60 flex-shrink-0">Location</span>
            <span className="text-white/90">{content.location}</span>
          </div>
        )}
        {content.links && content.links.length > 0 && (
          <div className="mt-3 pt-3 border-t border-white/20 space-y-1">
            {content.links.map((link, i) => (
              <div key={i} className="text-white/90 text-xs break-all">
                {link.label || link.url}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function SidebarSkills({
  content,
  title,
}: {
  content: SkillsContent;
  title: string;
}) {
  return (
    <div className="mb-6">
      <h3 className="text-sm font-semibold uppercase tracking-wider text-white/80 mb-3">
        {title}
      </h3>
      {content.categories.map((category) => (
        <div key={category.id} className="mb-3">
          {category.name && (
            <p className="text-xs font-medium text-white/70 mb-1">
              {category.name}
            </p>
          )}
          <ul className="space-y-1">
            {category.skills.map((skill, i) => (
              <li key={i} className="text-sm text-white/90 flex items-center gap-2">
                <span className="w-1 h-1 rounded-full bg-white/50 flex-shrink-0" />
                <span>{skill.name}</span>
                {skill.level && (
                  <span className="text-white/50 text-xs ml-auto">
                    {skill.level}
                  </span>
                )}
              </li>
            ))}
          </ul>
        </div>
      ))}
    </div>
  );
}

function SidebarLanguages({
  content,
  title,
}: {
  content: LanguagesContent;
  title: string;
}) {
  return (
    <div className="mb-6">
      <h3 className="text-sm font-semibold uppercase tracking-wider text-white/80 mb-3">
        {title}
      </h3>
      <ul className="space-y-2">
        {content.entries.map((entry) => (
          <li key={entry.id} className="text-sm">
            <span className="text-white/90">{entry.language}</span>
            <span className="text-white/50 ml-2 text-xs">
              {entry.proficiency}
            </span>
            {entry.certification && (
              <span className="text-white/40 ml-1 text-xs">
                ({entry.certification})
              </span>
            )}
          </li>
        ))}
      </ul>
    </div>
  );
}

// =============================================================================
// Main Content Components
// =============================================================================

function MainSectionTitle({
  title,
  primaryColor,
}: {
  title: string;
  primaryColor: string;
}) {
  return (
    <h2
      className="text-sm font-semibold uppercase tracking-wider mb-3 pb-2 border-b-2"
      style={{ color: primaryColor, borderColor: primaryColor }}
    >
      {title}
    </h2>
  );
}

function MainSummary({
  content,
  title,
  primaryColor,
  secondaryColor,
}: {
  content: SummaryContent;
  title: string;
  primaryColor: string;
  secondaryColor: string;
}) {
  return (
    <div className="mb-6">
      <MainSectionTitle title={title} primaryColor={primaryColor} />
      {content.text && (
        <p className="text-sm leading-relaxed" style={{ color: secondaryColor }}>
          {content.text}
        </p>
      )}
    </div>
  );
}

function MainExperience({
  content,
  title,
  primaryColor,
  secondaryColor,
}: {
  content: ExperienceContent;
  title: string;
  primaryColor: string;
  secondaryColor: string;
}) {
  return (
    <div className="mb-6">
      <MainSectionTitle title={title} primaryColor={primaryColor} />
      <div className="space-y-5">
        {content.entries.map((entry) => {
          const startDate = entry.startDate ? formatDate(entry.startDate) : '';
          const endDate = entry.current
            ? 'Present'
            : entry.endDate
              ? formatDate(entry.endDate)
              : '';
          const dateRange = [startDate, endDate].filter(Boolean).join(' - ');

          return (
            <div key={entry.id}>
              <div className="flex justify-between items-start mb-1">
                <div>
                  <h3 className="text-sm font-semibold text-gray-900">
                    {entry.position}
                  </h3>
                  <p className="text-sm" style={{ color: primaryColor }}>
                    {entry.company}
                    {entry.location && (
                      <span style={{ color: secondaryColor }}>
                        {' '}
                        &middot; {entry.location}
                      </span>
                    )}
                  </p>
                </div>
                {dateRange && (
                  <span
                    className="text-xs flex-shrink-0 ml-4"
                    style={{ color: secondaryColor }}
                  >
                    {dateRange}
                  </span>
                )}
              </div>
              {entry.description && (
                <p
                  className="text-xs mt-1 leading-relaxed"
                  style={{ color: secondaryColor }}
                >
                  {entry.description}
                </p>
              )}
              {entry.highlights && entry.highlights.length > 0 && (
                <ul className="mt-2 space-y-1">
                  {entry.highlights.map((highlight, i) => (
                    <li
                      key={i}
                      className="text-xs flex items-start gap-2"
                      style={{ color: secondaryColor }}
                    >
                      <span
                        className="w-1.5 h-1.5 rounded-full mt-1 flex-shrink-0"
                        style={{ backgroundColor: primaryColor }}
                      />
                      {highlight}
                    </li>
                  ))}
                </ul>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

function MainEducation({
  content,
  title,
  primaryColor,
  secondaryColor,
}: {
  content: EducationContent;
  title: string;
  primaryColor: string;
  secondaryColor: string;
}) {
  return (
    <div className="mb-6">
      <MainSectionTitle title={title} primaryColor={primaryColor} />
      <div className="space-y-4">
        {content.entries.map((entry) => {
          const startDate = entry.startDate ? formatDate(entry.startDate) : '';
          const endDate = entry.current
            ? 'Present'
            : entry.endDate
              ? formatDate(entry.endDate)
              : '';
          const dateRange = [startDate, endDate].filter(Boolean).join(' - ');

          return (
            <div key={entry.id}>
              <div className="flex justify-between items-start mb-1">
                <div>
                  <h3 className="text-sm font-semibold text-gray-900">
                    {entry.degree}
                    {entry.field && (
                      <span className="font-normal"> in {entry.field}</span>
                    )}
                  </h3>
                  <p className="text-sm" style={{ color: primaryColor }}>
                    {entry.institution}
                    {entry.location && (
                      <span style={{ color: secondaryColor }}>
                        {' '}
                        &middot; {entry.location}
                      </span>
                    )}
                  </p>
                </div>
                {dateRange && (
                  <span
                    className="text-xs flex-shrink-0 ml-4"
                    style={{ color: secondaryColor }}
                  >
                    {dateRange}
                  </span>
                )}
              </div>
              {entry.grade && (
                <p className="text-xs mt-1" style={{ color: secondaryColor }}>
                  Grade: {entry.grade}
                </p>
              )}
              {entry.description && (
                <p
                  className="text-xs mt-1 leading-relaxed"
                  style={{ color: secondaryColor }}
                >
                  {entry.description}
                </p>
              )}
            </div>
          );
        })}
      </div>
    </div>
  );
}

function MainGeneric({
  title,
  primaryColor,
}: {
  title: string;
  primaryColor: string;
}) {
  return (
    <div className="mb-6">
      <MainSectionTitle title={title} primaryColor={primaryColor} />
      <p className="text-xs text-gray-400 italic">
        Section content will appear here.
      </p>
    </div>
  );
}
