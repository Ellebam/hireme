import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@/test/utils';

// Mock API
const mockGetCV = vi.fn();
const mockCreateCV = vi.fn();
vi.mock('@/lib/api', () => ({
  api: {
    cv: {
      get: (...args: unknown[]) => mockGetCV(...args),
      create: (...args: unknown[]) => mockCreateCV(...args),
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

// Mock templates
vi.mock('@/lib/templates', () => ({
  starterTemplate: () => ({
    schemaVersion: '1.0',
    templateId: 'classic',
    title: 'My CV',
    sections: [],
  }),
}));

// Mock EditorLayout to avoid cascading child dependencies
vi.mock('@/components/editor', () => ({
  EditorLayout: () => <div data-testid="editor-layout">EditorLayout</div>,
}));

// Mock TooltipProvider
vi.mock('@/components/ui/tooltip', () => ({
  TooltipProvider: ({ children }: { children: React.ReactNode }) => (
    <>{children}</>
  ),
}));

import EditorPage from '../../app/editor/page';
import { mockCV } from '@/test/mocks/cv';

const { ApiError } = await import('@/lib/api');

describe('EditorPage', () => {
  beforeEach(() => {
    mockGetCV.mockReset();
    mockCreateCV.mockReset();
    mockSetCV.mockReset();
    mockReset.mockReset();
  });

  it('shows loading state initially', () => {
    mockGetCV.mockReturnValue(new Promise(() => {}));
    render(<EditorPage />);

    expect(screen.getByText('Loading editor...')).toBeInTheDocument();
  });

  it('renders EditorLayout on success', async () => {
    mockGetCV.mockResolvedValue(mockCV);
    render(<EditorPage />);

    await waitFor(() => {
      expect(screen.getByTestId('editor-layout')).toBeInTheDocument();
    });
  });

  it('shows error state on API failure', async () => {
    mockGetCV.mockRejectedValue(new Error('Network error'));
    render(<EditorPage />);

    await waitFor(() => {
      expect(screen.getByText('Something went wrong')).toBeInTheDocument();
    });
  });

  it('auto-creates CV when none exists (404)', async () => {
    mockGetCV.mockRejectedValue(new ApiError(404, 'not found'));
    mockCreateCV.mockResolvedValue(mockCV);
    render(<EditorPage />);

    await waitFor(() => {
      expect(screen.getByTestId('editor-layout')).toBeInTheDocument();
    });

    expect(mockCreateCV).toHaveBeenCalled();
  });
});
