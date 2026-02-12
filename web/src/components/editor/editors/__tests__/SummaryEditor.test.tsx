import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { SummaryEditor } from '../SummaryEditor';
import { mockSummaryContent } from '@/test/mocks/cv';

describe('SummaryEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders label and textarea with content', () => {
    render(
      <SummaryEditor sectionId="sec-1" content={mockSummaryContent} />
    );

    expect(screen.getByLabelText(/professional summary/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/professional summary/i)).toHaveValue(
      mockSummaryContent.text
    );
  });

  it('shows character count', () => {
    render(
      <SummaryEditor sectionId="sec-1" content={mockSummaryContent} />
    );

    const expectedCount = mockSummaryContent.text.length;
    expect(
      screen.getByText(`${expectedCount} characters`)
    ).toBeInTheDocument();
  });

  it('renders tips section', () => {
    render(
      <SummaryEditor sectionId="sec-1" content={mockSummaryContent} />
    );

    expect(
      screen.getByText('Tips for a great summary:')
    ).toBeInTheDocument();
  });
});
