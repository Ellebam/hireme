/// <reference types="vitest" />
import { defineConfig } from 'vitest/config';
import path from 'path';

export default defineConfig({
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
    include: ['src/**/*.test.{ts,tsx}'],
    exclude: ['node_modules', '.next'],
    coverage: {
      provider: 'v8',
      reporter: ['text', 'html', 'lcov'],
      include: ['src/**'],
      exclude: [
        'src/test/',
        '**/*.d.ts',
        '**/*.config.{ts,js,mjs,cjs}',
        '**/types/',
      ],
      thresholds: {
        // Match current coverage — raise incrementally as more tests are added
        statements: 18,
        branches: 40,
        functions: 30,
        lines: 18,
      },
    },
    // Clear mocks between tests
    clearMocks: true,
    restoreMocks: true,
  },
  resolve: {
    alias: {
      '@': path.resolve(__dirname, './src'),
    },
  },
});
