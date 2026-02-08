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
  LanguageProficiency,
} from '@/types/cv';
import { SECTION_LABELS } from '@/types/cv';

interface TemplateProps {
  sections: CVSection[];
  styling?: CVStyling;
  selectedSectionId: string | null;
  onSectionClick: (id: string) => void;
  onSectionDoubleClick: (id: string) => void;
}

const proficiencyLevels: Record<LanguageProficiency, number> = {
  native: 5,
  fluent: 4,
  advanced: 3,
  intermediate: 2,
  basic: 1,
};

function formatDate(dateStr: string): string {
  if (!dateStr) return '';
  const [year, month] = dateStr.split('-');
  if (!month) return year;
  const date = new Date(Number(year), Number(month) - 1);
  return date.toLocaleDateString('en', { year: 'numeric', month: 'short' });
}

export function ModernTemplate({
  sections,
  styling,
  selectedSectionId,
  onSectionClick,
  onSectionDoubleClick,
}: TemplateProps) {
  const primaryColor = styling?.primaryColor || '#2563eb';
  const secondaryColor = styling?.secondaryColor || '#64748b';

  const visibleSections = sections
    .filter((s) => s.visible !== false)
    .sort((a, b) => a.order - b.order);

  function renderSectionWrapper(section: CVSection, children: React.ReactNode) {
    const isSelected = selectedSectionId === section.id;
    return (
      <div
        key={section.id}
        className={cn(
          'relative p-4 -mx-4 rounded-lg cursor-pointer transition-all',
          'hover:bg-primary/5',
          isSelected && 'bg-primary/10 ring-2 ring-primary/30'
        )}
        onClick={() => onSectionClick(section.id)}
        onDoubleClick={() => onSectionDoubleClick(section.id)}
      >
        {isSelected && (
          <div className="absolute top-2 right-2 px-2 py-1 bg-primary text-primary-foreground text-xs rounded">
            Click to edit
          </div>
        )}
        {children}
      </div>
    );
  }

  function renderSectionTitle(section: CVSection) {
    const label = section.title || SECTION_LABELS[section.type] || section.type;
    return (
      <h2
        className="text-lg font-bold border-l-[3px] pl-3 mb-4"
        style={{ borderColor: primaryColor }}
      >
        {label}
      </h2>
    );
  }

  function renderPersonal(section: CVSection) {
    const content = section.content as PersonalContent;
    const fullName = [content.firstName, content.lastName].filter(Boolean).join(' ');
    const contactItems: string[] = [];
    if (content.email) contactItems.push(content.email);
    if (content.phone) contactItems.push(content.phone);
    if (content.location) contactItems.push(content.location);

    return renderSectionWrapper(
      section,
      <div>
        {fullName && (
          <h1 className="text-3xl font-bold tracking-tight">{fullName}</h1>
        )}
        {content.jobTitle && (
          <p className="text-lg mt-1" style={{ color: primaryColor }}>
            {content.jobTitle}
          </p>
        )}
        {contactItems.length > 0 && (
          <div
            className="flex flex-wrap gap-3 mt-3 text-sm"
            style={{ color: secondaryColor }}
          >
            {contactItems.map((item, i) => (
              <span key={i} className="flex items-center gap-1">
                {i > 0 && <span className="mr-1">|</span>}
                {item}
              </span>
            ))}
          </div>
        )}
        {content.links && content.links.length > 0 && (
          <div
            className="flex flex-wrap gap-3 mt-2 text-sm"
            style={{ color: secondaryColor }}
          >
            {content.links.map((link, i) => (
              <span key={i} className="underline">
                {link.label || link.url}
              </span>
            ))}
          </div>
        )}
      </div>
    );
  }

  function renderSummary(section: CVSection) {
    const content = section.content as SummaryContent;
    if (!content.text) return renderSectionWrapper(section, <>{renderSectionTitle(section)}</>);
    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <p className="text-sm leading-relaxed" style={{ color: secondaryColor }}>
          {content.text}
        </p>
      </div>
    );
  }

  function renderExperience(section: CVSection) {
    const content = section.content as ExperienceContent;
    const entries = content.entries || [];

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <div className="relative pl-6">
          {/* Timeline line */}
          <div
            className="absolute left-[7px] top-2 bottom-0 w-0.5"
            style={{ backgroundColor: primaryColor + '30' }}
          />
          {entries.map((entry) => {
            const startFormatted = formatDate(entry.startDate);
            const endFormatted = entry.current
              ? 'Present'
              : entry.endDate
                ? formatDate(entry.endDate)
                : '';
            const dateRange = [startFormatted, endFormatted]
              .filter(Boolean)
              .join(' - ');

            return (
              <div key={entry.id} className="relative mb-4">
                {/* Dot */}
                <div
                  className="absolute -left-6 top-1.5 w-3 h-3 rounded-full border-2 bg-white"
                  style={{ borderColor: primaryColor }}
                />
                {/* Content */}
                <div>
                  <div className="flex flex-wrap items-baseline justify-between gap-2">
                    <h3 className="font-semibold text-sm">{entry.position}</h3>
                    {dateRange && (
                      <span className="text-xs" style={{ color: secondaryColor }}>
                        {dateRange}
                      </span>
                    )}
                  </div>
                  <div className="flex flex-wrap items-center gap-2 mt-0.5">
                    <span className="text-sm" style={{ color: primaryColor }}>
                      {entry.company}
                    </span>
                    {entry.location && (
                      <span className="text-xs" style={{ color: secondaryColor }}>
                        | {entry.location}
                      </span>
                    )}
                  </div>
                  {entry.description && (
                    <p
                      className="text-sm mt-1 leading-relaxed"
                      style={{ color: secondaryColor }}
                    >
                      {entry.description}
                    </p>
                  )}
                  {entry.highlights && entry.highlights.length > 0 && (
                    <ul className="mt-1 space-y-0.5">
                      {entry.highlights.map((highlight, i) => (
                        <li
                          key={i}
                          className="text-sm flex items-start gap-2"
                          style={{ color: secondaryColor }}
                        >
                          <span
                            className="mt-1.5 w-1.5 h-1.5 rounded-full flex-shrink-0"
                            style={{ backgroundColor: primaryColor }}
                          />
                          {highlight}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  }

  function renderEducation(section: CVSection) {
    const content = section.content as EducationContent;
    const entries = content.entries || [];

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <div className="space-y-3">
          {entries.map((entry) => {
            const startFormatted = entry.startDate ? formatDate(entry.startDate) : '';
            const endFormatted = entry.current
              ? 'Present'
              : entry.endDate
                ? formatDate(entry.endDate)
                : '';
            const dateRange = [startFormatted, endFormatted]
              .filter(Boolean)
              .join(' - ');
            const degreeField = [entry.degree, entry.field]
              .filter(Boolean)
              .join(' in ');

            return (
              <div key={entry.id}>
                <div className="flex flex-wrap items-baseline justify-between gap-2">
                  <h3 className="font-semibold text-sm">{degreeField || entry.degree}</h3>
                  {dateRange && (
                    <span className="text-xs" style={{ color: secondaryColor }}>
                      {dateRange}
                    </span>
                  )}
                </div>
                <div className="flex flex-wrap items-center gap-2 mt-0.5">
                  <span className="text-sm" style={{ color: primaryColor }}>
                    {entry.institution}
                  </span>
                  {entry.location && (
                    <span className="text-xs" style={{ color: secondaryColor }}>
                      | {entry.location}
                    </span>
                  )}
                </div>
                {entry.grade && (
                  <p className="text-xs mt-0.5" style={{ color: secondaryColor }}>
                    Grade: {entry.grade}
                  </p>
                )}
                {entry.description && (
                  <p
                    className="text-sm mt-1 leading-relaxed"
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

  function renderSkills(section: CVSection) {
    const content = section.content as SkillsContent;
    const categories = content.categories || [];

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <div className="space-y-3">
          {categories.map((category) => (
            <div key={category.id}>
              {category.name && (
                <h3 className="text-sm font-medium mb-1.5">{category.name}</h3>
              )}
              <div className="flex flex-wrap gap-1.5">
                {category.skills.map((skill, i) => (
                  <span
                    key={i}
                    className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    style={{
                      backgroundColor: primaryColor + '15',
                      color: primaryColor,
                    }}
                  >
                    {skill.name}
                  </span>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    );
  }

  function renderLanguages(section: CVSection) {
    const content = section.content as LanguagesContent;
    const entries = content.entries || [];

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <div className="space-y-2.5">
          {entries.map((entry) => {
            const level = proficiencyLevels[entry.proficiency] || 1;
            return (
              <div key={entry.id} className="flex items-center gap-3">
                <span className="text-sm font-medium w-28 flex-shrink-0">
                  {entry.language}
                </span>
                <div className="flex gap-1">
                  {Array.from({ length: 5 }).map((_, i) => (
                    <div
                      key={i}
                      className="w-6 h-2 rounded-sm"
                      style={{
                        backgroundColor:
                          i < level ? primaryColor : primaryColor + '20',
                      }}
                    />
                  ))}
                </div>
                <span className="text-xs capitalize" style={{ color: secondaryColor }}>
                  {entry.proficiency}
                </span>
              </div>
            );
          })}
        </div>
      </div>
    );
  }

  function renderFallback(section: CVSection) {
    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <p className="text-sm" style={{ color: secondaryColor }}>
          {section.type} section
        </p>
      </div>
    );
  }

  function renderSection(section: CVSection) {
    switch (section.type) {
      case 'personal':
        return renderPersonal(section);
      case 'summary':
        return renderSummary(section);
      case 'experience':
        return renderExperience(section);
      case 'education':
        return renderEducation(section);
      case 'skills':
        return renderSkills(section);
      case 'languages':
        return renderLanguages(section);
      default:
        return renderFallback(section);
    }
  }

  return (
    <div
      style={
        {
          '--cv-primary': primaryColor,
          '--cv-secondary': secondaryColor,
        } as React.CSSProperties
      }
    >
      <div className="p-12">
        <div className="space-y-6">
          {visibleSections.map((section) => renderSection(section))}
        </div>
      </div>
    </div>
  );
}
