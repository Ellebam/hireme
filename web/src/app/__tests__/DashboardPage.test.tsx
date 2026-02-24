import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@/test/utils';

// Mock API
const mockListCVs = vi.fn();
const mockDeleteCV = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    cv: {
      list: (...args: unknown[]) => mockListCVs(...args),
      delete: (...args: unknown[]) => mockDeleteCV(...args),
    },
  },
  ApiError: class ApiError extends Error {
    status: number;
    constructor(status: number, message: string) {
      super(message);
      this.status = status;
    }
    get isNotFound() {
      return this.status === 404;
    }
  },
}));

// Mock AppShell to avoid rendering Header with its dependencies
vi.mock('@/components/layout', () => ({
  AppShell: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="app-shell">{children}</div>
  ),
}));

// Mock next/link to avoid IntersectionObserver issues in tests
vi.mock('next/link', () => ({
  default: ({
    children,
    href,
    ...props
  }: {
    children: React.ReactNode;
    href: string;
    [key: string]: unknown;
  }) => (
    <a href={href} {...props}>
      {children}
    </a>
  ),
}));

import DashboardPage from '../page';
import { mockCV, createMockCV } from '@/test/mocks/cv';

describe('DashboardPage', () => {
  beforeEach(() => {
    mockListCVs.mockReset();
    mockDeleteCV.mockReset();
  });

  it('shows loading skeleton initially', () => {
    mockListCVs.mockReturnValue(new Promise(() => {})); // never resolves
    render(<DashboardPage />);

    expect(screen.getByText('Your Workspace')).toBeInTheDocument();
  });

  it('renders CV card after data loads', async () => {
    mockListCVs.mockResolvedValue([mockCV]);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Software Engineer CV')).toBeInTheDocument();
    });

    expect(screen.getByText('modern')).toBeInTheDocument();
    expect(screen.getByText('8 sections')).toBeInTheDocument();
  });

  it('renders multiple CV cards', async () => {
    const cv2 = createMockCV({ id: 'cv-456', title: 'Design CV' });
    mockListCVs.mockResolvedValue([mockCV, cv2]);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Software Engineer CV')).toBeInTheDocument();
      expect(screen.getByText('Design CV')).toBeInTheDocument();
    });

    // Document count badge is dynamic
    expect(screen.getByText('2', { selector: 'span' })).toBeInTheDocument();
  });

  it('shows empty state when list returns empty array', async () => {
    mockListCVs.mockResolvedValue([]);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('No CVs yet')).toBeInTheDocument();
    });
  });

  it('shows error state on API failure', async () => {
    mockListCVs.mockRejectedValue(new Error('Network error'));
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Failed to load CVs')).toBeInTheDocument();
    });

    expect(screen.getByText('Try Again')).toBeInTheDocument();
  });

  it('edit links include CV ID', async () => {
    mockListCVs.mockResolvedValue([mockCV]);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Software Engineer CV')).toBeInTheDocument();
    });

    const editLinks = screen.getAllByRole('link').filter(
      (link) => link.getAttribute('href') === `/editor/${mockCV.id}`
    );
    expect(editLinks.length).toBeGreaterThan(0);
  });

  it('create links point to /templates', async () => {
    mockListCVs.mockResolvedValue([mockCV]);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Software Engineer CV')).toBeInTheDocument();
    });

    const createLink = screen.getByText('Create New Document').closest('a');
    expect(createLink).toHaveAttribute('href', '/templates');
  });
});
