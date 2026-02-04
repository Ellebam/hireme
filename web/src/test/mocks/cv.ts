/**
 * CV Test Mocks
 * Mock data for testing CV-related functionality
 */

import type { CV } from '@/types/api';
import type {
  CVContent,
  CVSection,
  PersonalContent,
  SummaryContent,
  ExperienceContent,
  EducationContent,
  SkillsContent,
  LanguagesContent,
  SCHEMA_VERSION,
} from '@/types/cv';

// ============================================================================
// Section Mocks
// ============================================================================

export const mockPersonalContent: PersonalContent = {
  firstName: 'John',
  lastName: 'Doe',
  jobTitle: 'Senior Software Engineer',
  email: 'john.doe@example.com',
  phone: '+1 555-123-4567',
  location: 'San Francisco, CA',
  portraitAssetId: null,
  links: [
    { type: 'linkedin', url: 'https://linkedin.com/in/johndoe' },
    { type: 'github', url: 'https://github.com/johndoe' },
  ],
};

export const mockSummaryContent: SummaryContent = {
  text: 'Experienced software engineer with 8+ years of experience building scalable web applications. Passionate about clean code, testing, and mentoring junior developers.',
};

export const mockExperienceContent: ExperienceContent = {
  entries: [
    {
      id: 'exp-1',
      company: 'Tech Corp',
      position: 'Senior Software Engineer',
      location: 'San Francisco, CA',
      startDate: '2020-01',
      endDate: null,
      current: true,
      description: 'Leading development of core platform services.',
      highlights: [
        'Architected microservices handling 1M+ requests/day',
        'Reduced deployment time by 60%',
        'Mentored team of 5 junior developers',
      ],
    },
    {
      id: 'exp-2',
      company: 'Startup Inc',
      position: 'Software Engineer',
      location: 'New York, NY',
      startDate: '2017-06',
      endDate: '2019-12',
      current: false,
      description: 'Full-stack development for e-commerce platform.',
      highlights: [
        'Built checkout flow serving 10K daily transactions',
        'Implemented real-time inventory system',
      ],
    },
  ],
};

export const mockEducationContent: EducationContent = {
  entries: [
    {
      id: 'edu-1',
      institution: 'Stanford University',
      degree: 'Master of Science',
      field: 'Computer Science',
      location: 'Stanford, CA',
      startDate: '2015-09',
      endDate: '2017-06',
      current: false,
      grade: '3.9 GPA',
    },
    {
      id: 'edu-2',
      institution: 'UC Berkeley',
      degree: 'Bachelor of Science',
      field: 'Computer Science',
      location: 'Berkeley, CA',
      startDate: '2011-09',
      endDate: '2015-05',
      current: false,
      grade: '3.8 GPA',
    },
  ],
};

export const mockSkillsContent: SkillsContent = {
  categories: [
    {
      id: 'cat-1',
      name: 'Languages',
      skills: [
        { name: 'TypeScript', level: 'expert' },
        { name: 'Go', level: 'advanced' },
        { name: 'Python', level: 'advanced' },
      ],
    },
    {
      id: 'cat-2',
      name: 'Frameworks',
      skills: [
        { name: 'React', level: 'expert' },
        { name: 'Next.js', level: 'advanced' },
        { name: 'Node.js', level: 'advanced' },
      ],
    },
  ],
};

export const mockLanguagesContent: LanguagesContent = {
  entries: [
    { language: 'English', proficiency: 'native' },
    { language: 'German', proficiency: 'fluent' },
    { language: 'Spanish', proficiency: 'intermediate' },
  ],
};

// ============================================================================
// Section Mocks
// ============================================================================

export const mockSections: CVSection[] = [
  {
    id: 'sec-personal',
    type: 'personal',
    order: 0,
    visible: true,
    content: mockPersonalContent,
  },
  {
    id: 'sec-summary',
    type: 'summary',
    order: 1,
    visible: true,
    content: mockSummaryContent,
  },
  {
    id: 'sec-experience',
    type: 'experience',
    order: 2,
    visible: true,
    title: 'Work Experience',
    content: mockExperienceContent,
  },
  {
    id: 'sec-education',
    type: 'education',
    order: 3,
    visible: true,
    content: mockEducationContent,
  },
  {
    id: 'sec-skills',
    type: 'skills',
    order: 4,
    visible: true,
    content: mockSkillsContent,
  },
  {
    id: 'sec-languages',
    type: 'languages',
    order: 5,
    visible: true,
    content: mockLanguagesContent,
  },
];

// ============================================================================
// CV Content Mock
// ============================================================================

export const mockCVContent: CVContent = {
  schemaVersion: '1.0.0',
  templateId: 'modern',
  locale: 'en',
  title: 'Software Engineer CV',
  sections: mockSections,
  styling: {
    primaryColor: '#2563eb',
    secondaryColor: '#64748b',
    fontFamily: 'inter',
    fontSize: 'medium',
    lineHeight: 'normal',
    showIcons: true,
  },
};

// ============================================================================
// Full CV Mock (API Response)
// ============================================================================

export const mockCV: CV = {
  id: 'cv-123',
  title: 'Software Engineer CV',
  schemaVersion: '1.0.0',
  content: mockCVContent,
  createdAt: '2024-01-15T10:30:00Z',
  updatedAt: '2024-02-01T14:22:00Z',
};

// ============================================================================
// Factory Functions
// ============================================================================

export function createMockCV(overrides: Partial<CV> = {}): CV {
  return {
    ...mockCV,
    ...overrides,
    id: overrides.id || `cv-${Date.now()}`,
  };
}

export function createMockCVContent(
  overrides: Partial<CVContent> = {}
): CVContent {
  return {
    ...mockCVContent,
    ...overrides,
  };
}

export function createMockSection<T extends CVSection['type']>(
  type: T,
  overrides: Partial<CVSection> = {}
): CVSection {
  const baseSection = mockSections.find((s) => s.type === type);
  if (!baseSection) {
    throw new Error(`No mock section found for type: ${type}`);
  }
  return {
    ...baseSection,
    ...overrides,
    id: overrides.id || `sec-${Date.now()}`,
  };
}

// ============================================================================
// Minimal Mocks (for unit tests)
// ============================================================================

export const minimalCVContent: CVContent = {
  schemaVersion: '1.0.0',
  templateId: 'classic',
  sections: [
    {
      id: 'sec-1',
      type: 'personal',
      order: 0,
      visible: true,
      content: { firstName: 'Test', lastName: 'User' },
    },
  ],
};

export const minimalCV: CV = {
  id: 'min-cv-1',
  title: 'Minimal CV',
  schemaVersion: '1.0.0',
  content: minimalCVContent,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};
