import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { SkillsEditor } from '../SkillsEditor';
import { mockSkillsContent } from '@/test/mocks/cv';

describe('SkillsEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders category cards with names', () => {
    render(
      <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
    );

    expect(screen.getByText('Languages')).toBeInTheDocument();
    expect(screen.getByText('Frameworks')).toBeInTheDocument();
    expect(screen.getByText('2 categories')).toBeInTheDocument();
  });

  it('renders skill tags with levels', () => {
    render(
      <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
    );

    expect(screen.getByText('TypeScript')).toBeInTheDocument();
    expect(screen.getByText('React')).toBeInTheDocument();
    expect(screen.getByText('Go')).toBeInTheDocument();
    expect(screen.getAllByText('(expert)')).toHaveLength(2);
    expect(screen.getAllByText('(advanced)')).toHaveLength(4);
  });

  it('shows empty state when no categories', () => {
    render(
      <SkillsEditor sectionId="sec-1" content={{ categories: [] }} />
    );

    expect(
      screen.getByText('No skill categories added')
    ).toBeInTheDocument();
  });
});
