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

export interface TemplateProps {
  sections: CVSection[];
  styling?: CVStyling;
  selectedSectionId: string | null;
  onSectionClick: (id: string) => void;
  onSectionDoubleClick: (id: string) => void;
}

function formatDate(dateStr: string): string {
  if (!dateStr) return '';
  const [year, month] = dateStr.split('-');
  if (!year) return '';
  if (!month) return year;
  const date = new Date(parseInt(year), parseInt(month) - 1);
  return date.toLocaleDateString('en-US', { month: 'short', year: 'numeric' });
}

export function ClassicTemplate({
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
    if (section.type === 'personal') return null;

    const label = section.title || SECTION_LABELS[section.type] || section.type;

    return (
      <h2
        className="text-sm font-semibold uppercase tracking-widest mb-3 pb-2"
        style={{
          borderBottom: `1px solid ${primaryColor}`,
          color: primaryColor,
        }}
      >
        {label}
      </h2>
    );
  }

  function renderPersonal(section: CVSection) {
    const content = section.content as PersonalContent;
    const fullName = [content.firstName, content.lastName].filter(Boolean).join(' ');
    const contactParts: string[] = [];
    if (content.email) contactParts.push(content.email);
    if (content.phone) contactParts.push(content.phone);
    if (content.location) contactParts.push(content.location);

    return renderSectionWrapper(
      section,
      <div className="text-center">
        {fullName && (
          <h1 className="text-2xl font-bold tracking-tight mb-1">{fullName}</h1>
        )}
        {content.jobTitle && (
          <p className="text-base mb-3" style={{ color: secondaryColor }}>
            {content.jobTitle}
          </p>
        )}
        {contactParts.length > 0 && (
          <p className="text-sm" style={{ color: secondaryColor }}>
            {contactParts.map((part, i) => (
              <span key={i}>
                {i > 0 && <span className="mx-2">&middot;</span>}
                {part}
              </span>
            ))}
          </p>
        )}
        {content.links && content.links.length > 0 && (
          <div className="flex justify-center gap-4 mt-2">
            {content.links.map((link, i) => (
              <a
                key={i}
                href={link.url}
                target="_blank"
                rel="noopener noreferrer"
                className="text-sm underline hover:no-underline"
                style={{ color: primaryColor }}
              >
                {link.label || link.type}
              </a>
            ))}
          </div>
        )}
        {fullName && (
          <div
            className="mt-4 mx-auto"
            style={{
              borderBottom: `2px solid ${primaryColor}`,
              width: '100%',
            }}
          />
        )}
      </div>
    );
  }

  function renderSummary(section: CVSection) {
    const content = section.content as SummaryContent;

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        {content.text ? (
          <p className="text-sm leading-relaxed" style={{ color: '#374151' }}>
            {content.text}
          </p>
        ) : (
          <p className="text-sm italic" style={{ color: secondaryColor }}>
            Add your professional summary...
          </p>
        )}
      </div>
    );
  }

  function renderExperience(section: CVSection) {
    const content = section.content as ExperienceContent;

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        {content.entries.length === 0 ? (
          <p className="text-sm italic" style={{ color: secondaryColor }}>
            Add your work experience...
          </p>
        ) : (
          <div className="space-y-4">
            {content.entries.map((entry) => {
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
                <div key={entry.id}>
                  <div className="flex justify-between items-baseline">
                    <div>
                      <span className="text-sm font-bold">{entry.position}</span>
                      {entry.company && (
                        <span className="text-sm">
                          {' '}
                          at {entry.company}
                        </span>
                      )}
                      {entry.location && (
                        <span className="text-sm" style={{ color: secondaryColor }}>
                          {' '}
                          &middot; {entry.location}
                        </span>
                      )}
                    </div>
                    {dateRange && (
                      <span
                        className="text-xs whitespace-nowrap ml-4"
                        style={{ color: secondaryColor }}
                      >
                        {dateRange}
                      </span>
                    )}
                  </div>
                  {entry.description && (
                    <p className="text-sm mt-1" style={{ color: '#374151' }}>
                      {entry.description}
                    </p>
                  )}
                  {entry.highlights && entry.highlights.length > 0 && (
                    <ul className="list-disc list-inside mt-1 space-y-0.5">
                      {entry.highlights.map((highlight, i) => (
                        <li
                          key={i}
                          className="text-sm"
                          style={{ color: '#374151' }}
                        >
                          {highlight}
                        </li>
                      ))}
                    </ul>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    );
  }

  function renderEducation(section: CVSection) {
    const content = section.content as EducationContent;

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        {content.entries.length === 0 ? (
          <p className="text-sm italic" style={{ color: secondaryColor }}>
            Add your education...
          </p>
        ) : (
          <div className="space-y-4">
            {content.entries.map((entry) => {
              const startFormatted = entry.startDate
                ? formatDate(entry.startDate)
                : '';
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
                  <div className="flex justify-between items-baseline">
                    <div>
                      <span className="text-sm font-bold">{degreeField}</span>
                      {entry.institution && (
                        <span className="text-sm">
                          {' '}
                          at {entry.institution}
                        </span>
                      )}
                      {entry.location && (
                        <span className="text-sm" style={{ color: secondaryColor }}>
                          {' '}
                          &middot; {entry.location}
                        </span>
                      )}
                    </div>
                    {dateRange && (
                      <span
                        className="text-xs whitespace-nowrap ml-4"
                        style={{ color: secondaryColor }}
                      >
                        {dateRange}
                      </span>
                    )}
                  </div>
                  {entry.grade && (
                    <p className="text-sm mt-1" style={{ color: secondaryColor }}>
                      Grade: {entry.grade}
                    </p>
                  )}
                  {entry.description && (
                    <p className="text-sm mt-1" style={{ color: '#374151' }}>
                      {entry.description}
                    </p>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>
    );
  }

  function renderSkills(section: CVSection) {
    const content = section.content as SkillsContent;

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        {content.categories.length === 0 ? (
          <p className="text-sm italic" style={{ color: secondaryColor }}>
            Add your skills...
          </p>
        ) : (
          <div className="space-y-2">
            {content.categories.map((category) => (
              <div key={category.id} className="text-sm">
                <span className="font-semibold">{category.name}:</span>{' '}
                <span style={{ color: '#374151' }}>
                  {category.skills.map((skill) => skill.name).join(', ')}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    );
  }

  function renderLanguages(section: CVSection) {
    const content = section.content as LanguagesContent;

    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        {content.entries.length === 0 ? (
          <p className="text-sm italic" style={{ color: secondaryColor }}>
            Add your languages...
          </p>
        ) : (
          <div className="text-sm">
            {content.entries.map((entry, i) => (
              <span key={entry.id}>
                {i > 0 && <span className="mx-1">&middot;</span>}
                <span className="font-semibold">{entry.language}</span>
                <span style={{ color: secondaryColor }}>
                  {' '}
                  ({entry.proficiency})
                </span>
                {entry.certification && (
                  <span style={{ color: secondaryColor }}>
                    {' '}
                    - {entry.certification}
                  </span>
                )}
              </span>
            ))}
          </div>
        )}
      </div>
    );
  }

  function renderGenericSection(section: CVSection) {
    return renderSectionWrapper(
      section,
      <div>
        {renderSectionTitle(section)}
        <p className="text-sm italic" style={{ color: secondaryColor }}>
          {section.title || SECTION_LABELS[section.type] || 'Section content'}
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
        return renderGenericSection(section);
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
