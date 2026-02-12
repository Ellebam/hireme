import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';

const mockUpdateSectionContent = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    updateSectionContent: mockUpdateSectionContent,
  }),
}));

import { PersonalEditor } from '../PersonalEditor';
import { mockPersonalContent } from '@/test/mocks/cv';

describe('PersonalEditor', () => {
  beforeEach(() => {
    mockUpdateSectionContent.mockReset();
  });

  it('renders labels and input values from content', () => {
    render(
      <PersonalEditor sectionId="sec-1" content={mockPersonalContent} />
    );

    expect(screen.getByLabelText(/first name/i)).toHaveValue('John');
    expect(screen.getByLabelText(/last name/i)).toHaveValue('Doe');
    expect(screen.getByLabelText(/job title/i)).toHaveValue(
      'Senior Software Engineer'
    );
    expect(screen.getByLabelText(/email/i)).toHaveValue(
      'john.doe@example.com'
    );
    expect(screen.getByLabelText(/phone/i)).toHaveValue('+1 555-123-4567');
    expect(screen.getByLabelText(/location/i)).toHaveValue(
      'San Francisco, CA'
    );
  });

  it('renders link URL inputs', () => {
    render(
      <PersonalEditor sectionId="sec-1" content={mockPersonalContent} />
    );

    expect(
      screen.getByDisplayValue('https://linkedin.com/in/johndoe')
    ).toBeInTheDocument();
    expect(
      screen.getByDisplayValue('https://github.com/johndoe')
    ).toBeInTheDocument();
  });

  it('shows empty links state when no links', () => {
    render(
      <PersonalEditor
        sectionId="sec-1"
        content={{ ...mockPersonalContent, links: [] }}
      />
    );

    expect(screen.getByText('No links added yet')).toBeInTheDocument();
  });
});
