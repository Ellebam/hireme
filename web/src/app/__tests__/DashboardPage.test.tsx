import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@/test/utils';

// Mock API
const mockGetCV = vi.fn();
const mockDeleteCV = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    cv: {
      get: (...args: unknown[]) => mockGetCV(...args),
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
import { mockCV } from '@/test/mocks/cv';

// Import the mocked ApiError for use in tests
const { ApiError } = await import('@/lib/api');

describe('DashboardPage', () => {
  beforeEach(() => {
    mockGetCV.mockReset();
    mockDeleteCV.mockReset();
  });

  it('shows loading skeleton initially', () => {
    mockGetCV.mockReturnValue(new Promise(() => {})); // never resolves
    render(<DashboardPage />);

    // During loading, the heading is still shown but content is skeleton
    expect(screen.getByText('Dashboard')).toBeInTheDocument();
  });

  it('renders CV card after data loads', async () => {
    mockGetCV.mockResolvedValue(mockCV);
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Software Engineer CV')).toBeInTheDocument();
    });

    expect(screen.getByText('modern')).toBeInTheDocument();
    expect(screen.getByText('8 sections')).toBeInTheDocument();
  });

  it('shows empty state when API returns 404', async () => {
    mockGetCV.mockRejectedValue(new ApiError(404, 'not found'));
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('No CV yet')).toBeInTheDocument();
    });
  });

  it('shows error state on API failure', async () => {
    mockGetCV.mockRejectedValue(new Error('Network error'));
    render(<DashboardPage />);

    await waitFor(() => {
      expect(screen.getByText('Failed to load CV')).toBeInTheDocument();
    });

    expect(screen.getByText('Try Again')).toBeInTheDocument();
  });
});
