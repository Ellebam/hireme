import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { EducationEditor } from '../EducationEditor';
import { mockEducationContent } from '@/test/mocks/cv';

describe('EducationEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders entry cards with degree and institution', () => {
    render(
      <EducationEditor sectionId="sec-1" content={mockEducationContent} />
    );

    expect(
      screen.getByText(/Master of Science/)
    ).toBeInTheDocument();
    expect(screen.getByText(/Stanford University/)).toBeInTheDocument();
    expect(
      screen.getByText(/Bachelor of Science/)
    ).toBeInTheDocument();
    expect(screen.getByText(/UC Berkeley/)).toBeInTheDocument();
  });

  it('shows entry count', () => {
    render(
      <EducationEditor sectionId="sec-1" content={mockEducationContent} />
    );

    expect(screen.getByText('2 entries')).toBeInTheDocument();
  });

  it('shows empty state when no entries', () => {
    render(
      <EducationEditor sectionId="sec-1" content={{ entries: [] }} />
    );

    expect(screen.getByText('No education added')).toBeInTheDocument();
  });
});
