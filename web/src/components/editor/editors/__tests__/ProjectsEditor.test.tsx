import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { ProjectsEditor } from '../ProjectsEditor';
import { mockProjectsContent } from '@/test/mocks/cv';

describe('ProjectsEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders entry cards with name and role', () => {
    render(
      <ProjectsEditor sectionId="sec-1" content={mockProjectsContent} />
    );

    expect(
      screen.getByText('Open Source CLI Tool')
    ).toBeInTheDocument();
    expect(screen.getByText(/Lead Developer · Go, Docker, GitHub Actions/)).toBeInTheDocument();
    expect(
      screen.getByText('E-Commerce Platform')
    ).toBeInTheDocument();
  });

  it('shows entry count', () => {
    render(
      <ProjectsEditor sectionId="sec-1" content={mockProjectsContent} />
    );

    expect(screen.getByText('2 entries')).toBeInTheDocument();
  });

  it('shows empty state when no entries', () => {
    render(
      <ProjectsEditor sectionId="sec-1" content={{ entries: [] }} />
    );

    expect(
      screen.getByText('No projects added')
    ).toBeInTheDocument();
  });
});
