/**
 * CV Schema Types
 * Generated from /schemas/cv-schema.json (v1.0.0)
 */

// ============================================================================
// Enums & Literals
// ============================================================================

export const SCHEMA_VERSION = '1.0.0' as const;

export const TEMPLATE_IDS = ['classic', 'modern', 'minimal'] as const;
export type TemplateId = (typeof TEMPLATE_IDS)[number];

export const LOCALES = ['en', 'de'] as const;
export type Locale = (typeof LOCALES)[number];

export const SECTION_TYPES = [
  'personal',
  'summary',
  'experience',
  'education',
  'skills',
  'languages',
  'certifications',
  'projects',
  'awards',
  'publications',
  'references',
  'custom',
] as const;
export type SectionType = (typeof SECTION_TYPES)[number];

export const LINK_TYPES = [
  'linkedin',
  'github',
  'twitter',
  'website',
  'portfolio',
  'other',
] as const;
export type LinkType = (typeof LINK_TYPES)[number];

export const SKILL_LEVELS = [
  'beginner',
  'intermediate',
  'advanced',
  'expert',
] as const;
export type SkillLevel = (typeof SKILL_LEVELS)[number];

export const LANGUAGE_PROFICIENCIES = [
  'native',
  'fluent',
  'advanced',
  'intermediate',
  'basic',
] as const;
export type LanguageProficiency = (typeof LANGUAGE_PROFICIENCIES)[number];

export const FONT_FAMILIES = [
  'inter',
  'roboto',
  'opensans',
  'lato',
  'merriweather',
] as const;
export type FontFamily = (typeof FONT_FAMILIES)[number];

export const FONT_SIZES = ['small', 'medium', 'large'] as const;
export type FontSize = (typeof FONT_SIZES)[number];

export const LINE_HEIGHTS = ['compact', 'normal', 'relaxed'] as const;
export type LineHeight = (typeof LINE_HEIGHTS)[number];

// ============================================================================
// Content Types
// ============================================================================

/** Social/professional link */
export interface ProfileLink {
  type: LinkType;
  url: string;
  label?: string;
}

/** Personal information section content */
export interface PersonalContent {
  firstName?: string;
  lastName?: string;
  jobTitle?: string;
  email?: string;
  phone?: string;
  location?: string;
  portraitAssetId?: string | null;
  links?: ProfileLink[];
}

/** Professional summary/objective */
export interface SummaryContent {
  text: string;
}

/** Work experience entry */
export interface ExperienceEntry {
  id: string;
  company: string;
  position: string;
  location?: string;
  startDate: string; // YYYY-MM
  endDate?: string | null; // YYYY-MM or null if current
  current?: boolean;
  description?: string;
  highlights?: string[];
}

export interface ExperienceContent {
  entries: ExperienceEntry[];
}

/** Education entry */
export interface EducationEntry {
  id: string;
  institution: string;
  degree: string;
  field?: string;
  location?: string;
  startDate?: string;
  endDate?: string | null;
  current?: boolean;
  grade?: string;
  description?: string;
}

export interface EducationContent {
  entries: EducationEntry[];
}

/** Individual skill with optional level */
export interface Skill {
  name: string;
  level?: SkillLevel;
}

/** Skill category grouping */
export interface SkillCategory {
  id: string;
  name: string;
  skills: Skill[];
}

export interface SkillsContent {
  categories: SkillCategory[];
}

/** Language entry */
export interface LanguageEntry {
  language: string;
  proficiency: LanguageProficiency;
  certification?: string;
}

export interface LanguagesContent {
  entries: LanguageEntry[];
}

/** Certification entry */
export interface CertificationEntry {
  id: string;
  name: string;
  issuer: string;
  date?: string;
  expiryDate?: string | null;
  credentialId?: string;
  url?: string;
}

export interface CertificationsContent {
  entries: CertificationEntry[];
}

/** Project entry */
export interface ProjectEntry {
  id: string;
  name: string;
  role?: string;
  description: string;
  url?: string;
  technologies?: string[];
  startDate?: string;
  endDate?: string | null;
}

export interface ProjectsContent {
  entries: ProjectEntry[];
}

/** Award entry */
export interface AwardEntry {
  id: string;
  title: string;
  issuer?: string;
  date?: string;
  description?: string;
}

export interface AwardsContent {
  entries: AwardEntry[];
}

/** Publication entry */
export interface PublicationEntry {
  id: string;
  title: string;
  publisher?: string;
  date?: string;
  url?: string;
  description?: string;
}

export interface PublicationsContent {
  entries: PublicationEntry[];
}

/** Reference entry */
export interface ReferenceEntry {
  id: string;
  name: string;
  position?: string;
  company?: string;
  email?: string;
  phone?: string;
  relationship?: string;
}

export interface ReferencesContent {
  available?: boolean;
  entries?: ReferenceEntry[];
}

/** Custom section with markdown */
export interface CustomContent {
  markdown: string;
}

/** Union of all section content types */
export type SectionContent =
  | PersonalContent
  | SummaryContent
  | ExperienceContent
  | EducationContent
  | SkillsContent
  | LanguagesContent
  | CertificationsContent
  | ProjectsContent
  | AwardsContent
  | PublicationsContent
  | ReferencesContent
  | CustomContent;

// ============================================================================
// Section Type
// ============================================================================

/** CV section with typed content */
export interface CVSection<T extends SectionType = SectionType> {
  id: string;
  type: T;
  order: number;
  visible?: boolean;
  title?: string;
  content: T extends 'personal'
    ? PersonalContent
    : T extends 'summary'
      ? SummaryContent
      : T extends 'experience'
        ? ExperienceContent
        : T extends 'education'
          ? EducationContent
          : T extends 'skills'
            ? SkillsContent
            : T extends 'languages'
              ? LanguagesContent
              : T extends 'certifications'
                ? CertificationsContent
                : T extends 'projects'
                  ? ProjectsContent
                  : T extends 'awards'
                    ? AwardsContent
                    : T extends 'publications'
                      ? PublicationsContent
                      : T extends 'references'
                        ? ReferencesContent
                        : T extends 'custom'
                          ? CustomContent
                          : SectionContent;
}

// ============================================================================
// Styling
// ============================================================================

export interface CVStyling {
  primaryColor?: string; // Hex color, default #2563eb
  secondaryColor?: string; // Hex color, default #64748b
  fontFamily?: FontFamily;
  fontSize?: FontSize;
  lineHeight?: LineHeight;
  showIcons?: boolean;
}

// ============================================================================
// CV Content (Root)
// ============================================================================

/** Full CV content structure */
export interface CVContent {
  schemaVersion: typeof SCHEMA_VERSION;
  templateId: TemplateId;
  locale?: Locale;
  title?: string;
  sections: CVSection[];
  styling?: CVStyling;
}

// ============================================================================
// Helpers
// ============================================================================

/** Type guard for section types */
export function isSectionType(type: string): type is SectionType {
  return SECTION_TYPES.includes(type as SectionType);
}

/** Get typed section content */
export function getSectionContent<T extends SectionType>(
  section: CVSection,
  expectedType: T
): CVSection<T>['content'] | null {
  if (section.type !== expectedType) return null;
  return section.content as CVSection<T>['content'];
}

/** Default content for each section type */
export function getDefaultContent(type: SectionType): SectionContent {
  switch (type) {
    case 'personal':
      return { links: [] } satisfies PersonalContent;
    case 'summary':
      return { text: '' } satisfies SummaryContent;
    case 'experience':
      return { entries: [] } satisfies ExperienceContent;
    case 'education':
      return { entries: [] } satisfies EducationContent;
    case 'skills':
      return { categories: [] } satisfies SkillsContent;
    case 'languages':
      return { entries: [] } satisfies LanguagesContent;
    case 'certifications':
      return { entries: [] } satisfies CertificationsContent;
    case 'projects':
      return { entries: [] } satisfies ProjectsContent;
    case 'awards':
      return { entries: [] } satisfies AwardsContent;
    case 'publications':
      return { entries: [] } satisfies PublicationsContent;
    case 'references':
      return { available: true, entries: [] } satisfies ReferencesContent;
    case 'custom':
      return { markdown: '' } satisfies CustomContent;
  }
}

/** Section type labels for display */
export const SECTION_LABELS: Record<SectionType, string> = {
  personal: 'Personal Info',
  summary: 'Summary',
  experience: 'Experience',
  education: 'Education',
  skills: 'Skills',
  languages: 'Languages',
  certifications: 'Certifications',
  projects: 'Projects',
  awards: 'Awards',
  publications: 'Publications',
  references: 'References',
  custom: 'Custom Section',
};
