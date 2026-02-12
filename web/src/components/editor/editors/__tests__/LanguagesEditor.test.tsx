import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { LanguagesEditor } from '../LanguagesEditor';
import { mockLanguagesContent } from '@/test/mocks/cv';

describe('LanguagesEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders language rows with values', () => {
    render(
      <LanguagesEditor sectionId="sec-1" content={mockLanguagesContent} />
    );

    expect(screen.getByDisplayValue('English')).toBeInTheDocument();
    expect(screen.getByDisplayValue('German')).toBeInTheDocument();
    expect(screen.getByDisplayValue('Spanish')).toBeInTheDocument();
  });

  it('shows language count', () => {
    render(
      <LanguagesEditor sectionId="sec-1" content={mockLanguagesContent} />
    );

    expect(screen.getByText('3 languages')).toBeInTheDocument();
  });

  it('shows empty state when no entries', () => {
    render(
      <LanguagesEditor sectionId="sec-1" content={{ entries: [] }} />
    );

    expect(screen.getByText('No languages added')).toBeInTheDocument();
  });
});
