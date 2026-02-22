import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@/test/utils';
import userEvent from '@testing-library/user-event';

const mockAddSection = vi.fn();
const mockReorderSections = vi.fn();
const mockSelectSection = vi.fn();
const mockOpenDeleteConfirmModal = vi.fn();

vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    addSection: mockAddSection,
    reorderSections: mockReorderSections,
    selectSection: mockSelectSection,
    selectedSectionId: null,
  }),
  useSortedSections: () => [
    {
      id: 'sec-personal',
      type: 'personal',
      order: 0,
      visible: true,
      content: { firstName: 'John', lastName: 'Doe' },
    },
    {
      id: 'sec-experience',
      type: 'experience',
      order: 1,
      visible: true,
      title: 'Work Experience',
      content: { entries: [] },
    },
    {
      id: 'sec-skills',
      type: 'skills',
      order: 2,
      visible: false,
      content: { categories: [] },
    },
  ],
  useUIStore: () => ({
    openDeleteConfirmModal: mockOpenDeleteConfirmModal,
  }),
}));

vi.mock('@/lib/dnd', () => ({
  useDndSensors: () => [],
  collisionDetection: undefined,
  sortStrategy: undefined,
  useSortableItem: () => ({
    ref: vi.fn(),
    style: {},
    isDragging: false,
    isOver: false,
    dragHandleProps: {},
  }),
}));

vi.mock('@dnd-kit/core', () => ({
  DndContext: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  DragOverlay: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="drag-overlay">{children}</div>
  ),
}));

vi.mock('@dnd-kit/sortable', () => ({
  SortableContext: ({ children }: { children: React.ReactNode }) => <>{children}</>,
}));

vi.mock('@/components/ui/tooltip', () => ({
  Tooltip: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  TooltipTrigger: ({ children }: { children: React.ReactNode }) => <>{children}</>,
  TooltipContent: () => null,
}));

import { SectionPalette } from '../SectionPalette';

describe('SectionPalette — delete button', () => {
  beforeEach(() => {
    mockAddSection.mockReset();
    mockReorderSections.mockReset();
    mockSelectSection.mockReset();
    mockOpenDeleteConfirmModal.mockReset();
  });

  it('renders a delete button for each section in the structure list', () => {
    render(<SectionPalette />);

    expect(screen.getByLabelText('Delete Personal Info')).toBeInTheDocument();
    expect(screen.getByLabelText('Delete Work Experience')).toBeInTheDocument();
    expect(screen.getByLabelText('Delete Skills')).toBeInTheDocument();
  });

  it('uses custom title in aria-label when set, falls back to SECTION_LABELS otherwise', () => {
    render(<SectionPalette />);

    // sec-experience has title: 'Work Experience' (custom)
    expect(screen.getByLabelText('Delete Work Experience')).toBeInTheDocument();

    // sec-personal has no custom title — falls back to SECTION_LABELS['personal'] = 'Personal Info'
    expect(screen.getByLabelText('Delete Personal Info')).toBeInTheDocument();
  });

  it('calls openDeleteConfirmModal with correct section ID on click', async () => {
    const user = userEvent.setup();
    render(<SectionPalette />);

    await user.click(screen.getByLabelText('Delete Work Experience'));

    expect(mockOpenDeleteConfirmModal).toHaveBeenCalledTimes(1);
    expect(mockOpenDeleteConfirmModal).toHaveBeenCalledWith('sec-experience');
  });

  it('does not trigger section selection when clicking delete', async () => {
    const user = userEvent.setup();
    render(<SectionPalette />);

    await user.click(screen.getByLabelText('Delete Personal Info'));

    expect(mockOpenDeleteConfirmModal).toHaveBeenCalledTimes(1);
    expect(mockSelectSection).not.toHaveBeenCalled();
  });

  it('renders delete button for hidden sections', () => {
    render(<SectionPalette />);

    // sec-skills has visible: false
    expect(screen.getByLabelText('Delete Skills')).toBeInTheDocument();
  });
});
