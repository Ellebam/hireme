/**
 * Vitest Setup
 * Global test configuration and mocks
 */

import '@testing-library/jest-dom';
import { vi, beforeEach, afterEach } from 'vitest';

// ============================================================================
// Next.js Mocks
// ============================================================================

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    refresh: vi.fn(),
    prefetch: vi.fn(),
  }),
  useParams: () => ({}),
  useSearchParams: () => new URLSearchParams(),
  usePathname: () => '/',
}));

// ============================================================================
// next-intl Mocks
// ============================================================================

vi.mock('next-intl', () => ({
  useTranslations: () => (key: string) => key,
  useLocale: () => 'en',
  useMessages: () => ({}),
  useTimeZone: () => 'UTC',
  useNow: () => new Date(),
  NextIntlClientProvider: ({ children }: { children: React.ReactNode }) =>
    children,
}));

// ============================================================================
// LocalStorage Mock (for Zustand persist)
// ============================================================================

const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    getItem: (key: string) => store[key] ?? null,
    setItem: (key: string, value: string) => {
      store[key] = value;
    },
    removeItem: (key: string) => {
      delete store[key];
    },
    clear: () => {
      store = {};
    },
    get length() {
      return Object.keys(store).length;
    },
    key: (index: number) => Object.keys(store)[index] ?? null,
  };
})();

Object.defineProperty(window, 'localStorage', {
  value: localStorageMock,
  writable: true,
});

// ============================================================================
// Window/Browser Mocks
// ============================================================================

// Mock matchMedia for responsive tests
Object.defineProperty(window, 'matchMedia', {
  writable: true,
  value: vi.fn().mockImplementation((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock ResizeObserver
global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}));

// Mock IntersectionObserver
global.IntersectionObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
  root: null,
  rootMargin: '',
  thresholds: [],
}));

// ============================================================================
// Crypto Mock (for generateId)
// ============================================================================

let uuidCounter = 0;
Object.defineProperty(global, 'crypto', {
  value: {
    randomUUID: () => {
      uuidCounter++;
      return `test-uuid-${uuidCounter.toString().padStart(4, '0')}`;
    },
    getRandomValues: <T extends ArrayBufferView | null>(array: T): T => {
      if (array) {
        const uint8Array = array as unknown as Uint8Array;
        for (let i = 0; i < uint8Array.length; i++) {
          uint8Array[i] = (i * 17 + 42) % 256;
        }
      }
      return array;
    },
  },
});

// ============================================================================
// Zustand Store Reset
// ============================================================================

// Reset stores between tests to prevent state leakage
beforeEach(() => {
  // Clear localStorage to reset persisted store state
  window.localStorage.clear();
  // Reset UUID counter for deterministic IDs
  uuidCounter = 0;
});

afterEach(() => {
  vi.clearAllMocks();
});

// ============================================================================
// Console Suppression (optional)
// ============================================================================

// Suppress specific console warnings during tests
const originalError = console.error;
console.error = (...args: unknown[]) => {
  // Filter out React act() warnings that are often noise in tests
  if (
    typeof args[0] === 'string' &&
    args[0].includes('Warning: ReactDOM.render is no longer supported')
  ) {
    return;
  }
  originalError.call(console, ...args);
};
