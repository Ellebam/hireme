import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@/test/utils';

// Mock next/navigation
const mockReplace = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace: mockReplace }),
}));

import EditorRedirectPage from '../editor/page';

describe('EditorRedirectPage', () => {
  beforeEach(() => {
    mockReplace.mockReset();
  });

  it('redirects to dashboard', () => {
    render(<EditorRedirectPage />);
    expect(mockReplace).toHaveBeenCalledWith('/');
  });
});
