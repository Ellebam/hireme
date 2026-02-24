import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@/test/utils';

// Mock next/navigation
const mockReplace = vi.fn();
const mockPush = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: mockReplace, push: mockPush }),
  useParams: () => ({ id: 'cv-123' }),
}));

// Mock API
const mockGetCV = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    cv: {
      get: (...args: unknown[]) => mockGetCV(...args),
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
    get isForbidden() {
      return this.status === 403;
    }
  },
}));

// Mock logger
vi.mock('@/lib/logger', () => ({
  logger: {
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
}));

// Mock hooks
vi.mock('@/hooks', () => ({
  useAutoSave: vi.fn(),
  useKeyboardShortcuts: vi.fn(),
}));

// Mock store
const mockSetCV = vi.fn();
const mockReset = vi.fn();
vi.mock('@/stores', () => ({
  useEditorStore: () => ({
    setCV: mockSetCV,
    reset: mockReset,
  }),
}));

// Mock EditorLayout
vi.mock('@/components/editor', () => ({
  EditorLayout: () => <div data-testid="editor-layout">EditorLayout</div>,
}));

// Mock AppShell
vi.mock('@/components/layout', () => ({
  AppShell: ({ children }: { children: React.ReactNode }) => (
    <div data-testid="app-shell">{children}</div>
  ),
}));

import EditorIdPage from '../editor/[id]/page';
import { mockCV } from '@/test/mocks/cv';

const { ApiError } = await import('@/lib/api');

describe('EditorIdPage', () => {
  beforeEach(() => {
    mockGetCV.mockReset();
    mockSetCV.mockReset();
    mockReset.mockReset();
    mockReplace.mockReset();
    mockPush.mockReset();
  });

  it('shows loading state initially', () => {
    mockGetCV.mockReturnValue(new Promise(() => {}));
    render(<EditorIdPage />);

    expect(screen.getByText('Loading editor...')).toBeInTheDocument();
  });

  it('renders EditorLayout on success', async () => {
    mockGetCV.mockResolvedValue(mockCV);
    render(<EditorIdPage />);

    await waitFor(() => {
      expect(screen.getByTestId('editor-layout')).toBeInTheDocument();
    });

    expect(mockGetCV).toHaveBeenCalledWith('cv-123');
    expect(mockSetCV).toHaveBeenCalledWith(mockCV);
  });

  it('redirects to dashboard on 404', async () => {
    mockGetCV.mockRejectedValue(new ApiError(404, 'not found'));
    render(<EditorIdPage />);

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith('/');
    });
  });

  it('redirects to dashboard on 403', async () => {
    mockGetCV.mockRejectedValue(new ApiError(403, 'forbidden'));
    render(<EditorIdPage />);

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith('/');
    });
  });

  it('shows error state on other API failure', async () => {
    mockGetCV.mockRejectedValue(new Error('Network error'));
    render(<EditorIdPage />);

    await waitFor(() => {
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });
  });

  it('does not auto-create CV on missing', async () => {
    const mockCreateCV = vi.fn();
    mockGetCV.mockRejectedValue(new ApiError(404, 'not found'));
    render(<EditorIdPage />);

    await waitFor(() => {
      expect(mockReplace).toHaveBeenCalledWith('/');
    });

    // Verify no create call was made (unlike old editor page which auto-created)
    expect(mockCreateCV).not.toHaveBeenCalled();
  });
});
