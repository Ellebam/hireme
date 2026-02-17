import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, setupUser } from '@/test/utils';

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

  describe('interactions', () => {
    it('should add new category when Add Category is clicked', async () => {
      const user = setupUser();
      render(
        <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
      );

      await user.click(screen.getByRole('button', { name: /add category/i }));

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        categories: [
          ...mockSkillsContent.categories,
          expect.objectContaining({ name: 'New Category', skills: [] }),
        ],
      });
    });

    it('should enter edit mode when category name is clicked', async () => {
      const user = setupUser();
      render(
        <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
      );

      await user.click(screen.getByText('Languages'));

      expect(screen.getByDisplayValue('Languages')).toBeInTheDocument();
    });

    it('should delete category', async () => {
      const user = setupUser();
      render(
        <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
      );

      // Trash2 icon buttons on categories — one per category
      const deleteButtons = screen.getAllByRole('button').filter(
        (btn) => !btn.textContent?.trim() && btn.className.includes('destructive')
      );
      await user.click(deleteButtons[0]);

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        categories: [mockSkillsContent.categories[1]],
      });
    });

    it('should add skill to category', async () => {
      const user = setupUser();
      render(
        <SkillsEditor sectionId="sec-1" content={mockSkillsContent} />
      );

      // There are two "Add Skill" buttons (one per category)
      const addSkillButtons = screen.getAllByRole('button', { name: /add skill/i });
      await user.click(addSkillButtons[0]);

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        categories: expect.arrayContaining([
          expect.objectContaining({
            id: 'cat-1',
            skills: [
              ...mockSkillsContent.categories[0].skills,
              { name: '', level: undefined },
            ],
          }),
        ]),
      });
    });

    it('should show edit input when new empty skill is added', () => {
      const contentWithEmptySkill = {
        categories: [
          {
            id: 'cat-1',
            name: 'Test',
            skills: [{ name: '', level: undefined }],
          },
        ],
      };
      render(
        <SkillsEditor sectionId="sec-1" content={contentWithEmptySkill} />
      );

      expect(screen.getByPlaceholderText('Skill name')).toBeInTheDocument();
    });

    it('should delete skill when X is clicked', async () => {
      const user = setupUser();
      const singleCategoryContent = {
        categories: [
          {
            id: 'cat-1',
            name: 'Languages',
            skills: [
              { name: 'TypeScript', level: 'expert' as const },
              { name: 'Go', level: 'advanced' as const },
            ],
          },
        ],
      };
      render(
        <SkillsEditor sectionId="sec-1" content={singleCategoryContent} />
      );

      // X buttons on skill tags — hidden until hover, but still in DOM
      // Each skill tag has a delete button; find the first one
      const deleteButtons = screen.getAllByRole('button').filter(
        (btn) => !btn.textContent?.trim() && !btn.className.includes('destructive')
      );
      // First non-destructive icon button after category delete is a skill X button
      // More reliable: find buttons inside skill tags (the small X buttons)
      const skillDeleteBtn = screen.getAllByRole('button').find(
        (btn) => btn.closest('[class*="rounded-md bg-muted"]') && !btn.textContent?.trim()
      );
      if (skillDeleteBtn) await user.click(skillDeleteBtn);

      expect(mockUpdateSectionContent).toHaveBeenCalledWith('sec-1', {
        categories: [
          expect.objectContaining({
            id: 'cat-1',
            skills: [{ name: 'Go', level: 'advanced' }],
          }),
        ],
      });
    });
  });
});
