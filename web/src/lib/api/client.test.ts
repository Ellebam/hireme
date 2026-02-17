/**
 * API Client Tests
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { api, ApiError } from './client';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

// Mock logger to avoid console noise
vi.mock('@/lib/logger', () => ({
  logger: {
    debug: vi.fn(),
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
  },
}));

describe('ApiClient', () => {
  beforeEach(() => {
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('ApiError', () => {
    it('should have correct properties', () => {
      const error = new ApiError(404, 'Test error', undefined, {
        field: 'validation failed',
      });

      expect(error.message).toBe('Test error');
      expect(error.status).toBe(404);
      expect(error.details).toEqual({ field: 'validation failed' });
      expect(error.isNotFound).toBe(true);
      expect(error.isUnauthorized).toBe(false);
    });

    it('should identify 401 as unauthorized', () => {
      const error = new ApiError(401, 'Unauthorized');
      expect(error.isUnauthorized).toBe(true);
    });

    it('should identify 403 as forbidden', () => {
      const error = new ApiError(403, 'Forbidden');
      expect(error.isForbidden).toBe(true);
    });

    it('should identify 422 as validation error', () => {
      const error = new ApiError(422, 'Validation failed');
      expect(error.isValidationError).toBe(true);
    });

    it('should identify 400 as validation error', () => {
      const error = new ApiError(400, 'Bad request');
      expect(error.isValidationError).toBe(true);
    });
  });

  describe('cvApi', () => {
    it('should get CV and unwrap data', async () => {
      const mockCV = { id: 'cv-123', title: 'My CV' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: mockCV }),
      });

      const result = await api.cv.get();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result).toEqual(mockCV);
    });

    it('should create CV and unwrap data', async () => {
      const mockCV = { id: 'cv-123', title: 'New CV' };
      const createData = { title: 'New CV', content: {} };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: mockCV }),
      });

      const result = await api.cv.create(createData as any);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(createData),
        })
      );
      expect(result).toEqual(mockCV);
    });

    it('should update CV and unwrap data', async () => {
      const mockCV = { id: 'cv-123', title: 'Updated CV' };
      const updateData = { title: 'Updated CV' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: mockCV }),
      });

      const result = await api.cv.update('cv-123', updateData);

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv/cv-123'),
        expect.objectContaining({
          method: 'PUT',
          body: JSON.stringify(updateData),
        })
      );
      expect(result).toEqual(mockCV);
    });

    it('should delete CV', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 204,
        json: () => Promise.resolve({}),
      });

      await api.cv.delete('cv-123');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv/cv-123'),
        expect.objectContaining({ method: 'DELETE' })
      );
    });

    it('should throw ApiError on error response', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 404,
        json: () =>
          Promise.resolve({
            error: { code: 'not_found', message: 'CV not found' },
          }),
      });

      try {
        await api.cv.get();
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(404);
        expect((err as ApiError).message).toBe('CV not found');
      }
    });
  });

  describe('userApi', () => {
    it('should get current user and unwrap data', async () => {
      const mockUser = { id: 'user-123', email: 'test@example.com' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ data: mockUser }),
      });

      const result = await api.user.getMe();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/users/me'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result).toEqual(mockUser);
    });
  });

  describe('error handling', () => {
    it('should extract error message from error wrapper', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () =>
          Promise.resolve({
            error: { code: 'internal_error', message: 'Server error' },
          }),
      });

      try {
        await api.cv.get();
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).message).toBe('Server error');
        expect((err as ApiError).code).toBe('internal_error');
      }
    });

    it('should handle network errors', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network failure'));

      try {
        await api.cv.get();
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(0);
        expect((err as ApiError).message).toContain('Network error');
      }
    });
  });

  describe('exportApi', () => {
    it('should return blob on successful export', async () => {
      const mockBlob = new Blob(['pdf-data'], { type: 'application/pdf' });
      mockFetch.mockResolvedValueOnce({
        ok: true,
        blob: () => Promise.resolve(mockBlob),
      });

      const result = await api.export.export('cv-123', 'pdf');

      expect(result).toBe(mockBlob);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/export/pdf'),
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ cvId: 'cv-123' }),
        })
      );
    });

    it('should send correct request body for docx', async () => {
      const mockBlob = new Blob(['docx-data']);
      mockFetch.mockResolvedValueOnce({
        ok: true,
        blob: () => Promise.resolve(mockBlob),
      });

      await api.export.export('cv-456', 'docx');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/export/docx'),
        expect.objectContaining({
          body: JSON.stringify({ cvId: 'cv-456' }),
        })
      );
    });

    it('should throw ApiError on HTTP error', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () =>
          Promise.resolve({
            error: { code: 'export_failed', message: 'Gotenberg error' },
          }),
      });

      try {
        await api.export.export('cv-1', 'pdf');
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(500);
        expect((err as ApiError).message).toBe('Gotenberg error');
      }
    });

    it('should throw ApiError with default message on unparseable error body', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 502,
        json: () => Promise.reject(new Error('not json')),
      });

      try {
        await api.export.export('cv-1', 'pdf');
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(502);
        expect((err as ApiError).message).toBe('Export failed');
      }
    });

    it('should throw ApiError on abort (timeout path)', async () => {
      // Directly test the abort error handling path:
      // When fetch rejects with AbortError, export should throw ApiError
      mockFetch.mockRejectedValueOnce(
        new DOMException('The operation was aborted.', 'AbortError')
      );

      try {
        await api.export.export('cv-1', 'pdf');
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(0);
        expect((err as ApiError).message).toContain('timed out');
      }
    });

    it('should throw ApiError on network error', async () => {
      mockFetch.mockRejectedValueOnce(new TypeError('Failed to fetch'));

      try {
        await api.export.export('cv-1', 'pdf');
        expect.fail('Should have thrown');
      } catch (err) {
        expect(err).toBeInstanceOf(ApiError);
        expect((err as ApiError).status).toBe(0);
        expect((err as ApiError).message).toContain('network error');
      }
    });
  });
});
