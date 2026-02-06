/**
 * Test Utilities
 * Helper functions for testing React components
 */

import React, { ReactElement } from 'react';
import { render, RenderOptions } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

// ============================================================================
// Providers Wrapper
// ============================================================================

interface ProvidersProps {
  children: React.ReactNode;
}

/**
 * Wraps components with all necessary providers for testing
 */
function AllProviders({ children }: ProvidersProps) {
  return <>{children}</>;
}

// ============================================================================
// Custom Render
// ============================================================================

interface CustomRenderOptions extends Omit<RenderOptions, 'wrapper'> {
  // Add any custom options here
}

/**
 * Custom render function that wraps components with providers
 */
function customRender(ui: ReactElement, options?: CustomRenderOptions) {
  return render(ui, {
    wrapper: AllProviders,
    ...options,
  });
}

// ============================================================================
// User Event Setup
// ============================================================================

/**
 * Setup user event for testing user interactions
 */
function setupUser() {
  return userEvent.setup();
}

// ============================================================================
// Async Helpers
// ============================================================================

/**
 * Wait for a condition to be true
 */
async function waitForCondition(
  condition: () => boolean,
  timeout = 5000,
  interval = 50
): Promise<void> {
  const startTime = Date.now();
  while (!condition()) {
    if (Date.now() - startTime > timeout) {
      throw new Error('Condition not met within timeout');
    }
    await new Promise((resolve) => setTimeout(resolve, interval));
  }
}

/**
 * Wait for next tick
 */
function waitForNextTick(): Promise<void> {
  return new Promise((resolve) => setTimeout(resolve, 0));
}

// ============================================================================
// Store Helpers
// ============================================================================

/**
 * Reset a Zustand store to initial state
 */
function resetStore<T extends object>(
  useStore: { setState: (state: T) => void; getInitialState?: () => T },
  initialState: T
) {
  useStore.setState(initialState);
}

// ============================================================================
// Exports
// ============================================================================

export * from '@testing-library/react';
export { customRender as render, setupUser, waitForCondition, waitForNextTick, resetStore };
