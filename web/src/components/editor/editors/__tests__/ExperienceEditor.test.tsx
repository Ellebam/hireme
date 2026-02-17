import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, setupUser } from '@/test/utils';

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

  describe('interactions', () => {
    // Button order with mockExperienceContent (2 entries):
    // [0] "Add Experience", [1] Edit entry1, [2] Delete entry1, [3] Edit entry2, [4] Delete entry2
    function getIconButtons() {
      return screen.getAllByRole('button').filter((btn) => !btn.textContent?.trim());
    }

    it('should open modal when Add Experience is clicked', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
      );

      await user.click(screen.getByRole('button', { name: /add experience/i }));

      expect(screen.getByText('Add Work Experience')).toBeInTheDocument();
    });

    it('should fill form and save new entry', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={{ entries: [] }} />
      );

      await user.click(screen.getByRole('button', { name: /add experience/i }));
      await user.type(screen.getByLabelText('Company'), 'Acme Inc.');
      await user.type(screen.getByLabelText('Position'), 'Engineer');
      await user.click(screen.getByRole('button', { name: /save/i }));

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        entries: [
          expect.objectContaining({ company: 'Acme Inc.', position: 'Engineer' }),
        ],
      });
    });

    it('should save highlights as parsed array', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={{ entries: [] }} />
      );

      await user.click(screen.getByRole('button', { name: /add experience/i }));
      await user.type(screen.getByLabelText(/key achievements/i), 'Led migration\nReduced costs');
      await user.click(screen.getByRole('button', { name: /save/i }));

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        entries: [
          expect.objectContaining({
            highlights: ['Led migration', 'Reduced costs'],
          }),
        ],
      });
    });

    it('should open modal with entry data when Edit is clicked', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
      );

      // First icon button is edit for entry 1
      const iconButtons = getIconButtons();
      await user.click(iconButtons[0]);

      expect(screen.getByDisplayValue('Tech Corp')).toBeInTheDocument();
      expect(screen.getByDisplayValue('Senior Software Engineer')).toBeInTheDocument();
    });

    it('should save edited entry', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
      );

      const iconButtons = getIconButtons();
      await user.click(iconButtons[0]); // Edit entry 1

      const positionInput = screen.getByLabelText('Position');
      await user.clear(positionInput);
      await user.type(positionInput, 'Staff Engineer');
      await user.click(screen.getByRole('button', { name: /save/i }));

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        entries: expect.arrayContaining([
          expect.objectContaining({ id: 'exp-1', position: 'Staff Engineer' }),
        ]),
      });
    });

    it('should delete entry when Delete is clicked', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
      );

      // Second icon button is delete for entry 1
      const iconButtons = getIconButtons();
      await user.click(iconButtons[1]);

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        entries: [expect.objectContaining({ id: 'exp-2' })],
      });
    });

    it('should close modal when Cancel is clicked', async () => {
      const user = setupUser();
      render(
        <ExperienceEditor sectionId="sec-1" content={mockExperienceContent} />
      );

      await user.click(screen.getByRole('button', { name: /add experience/i }));
      expect(screen.getByText('Add Work Experience')).toBeInTheDocument();

      await user.click(screen.getByRole('button', { name: /cancel/i }));

      expect(screen.queryByText('Add Work Experience')).not.toBeInTheDocument();
    });
  });
});
