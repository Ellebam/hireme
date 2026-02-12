import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { ExperienceEditor } from '../ExperienceEditor';
import { mockExperienceContent } from '@/test/mocks/cv';

describe('ExperienceEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders entry cards with position and company', () => {
    render(
      <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
    );

    expect(
      screen.getByText('Senior Software Engineer')
    ).toBeInTheDocument();
    expect(screen.getByText(/Tech Corp/)).toBeInTheDocument();
    expect(screen.getByText('Software Engineer')).toBeInTheDocument();
    expect(screen.getByText(/Startup Inc/)).toBeInTheDocument();
  });

  it('shows entry count', () => {
    render(
      <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
    );

    expect(screen.getByText('2 entries')).toBeInTheDocument();
  });

  it('shows empty state when no entries', () => {
    render(
      <ExperienceEditor sectionId="sec-1" content={{ entries: [] }} />
    );

    expect(
      screen.getByText('No work experience added')
    ).toBeInTheDocument();
  });
});
