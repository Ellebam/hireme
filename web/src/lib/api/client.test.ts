/**
 * API Client Tests
 */

import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { api, ApiError } from './client';

// Mock fetch
const mockFetch = vi.fn();
global.fetch = mockFetch;

describe('ApiClient', () => {
  beforeEach(() => {
    mockFetch.mockReset();
  });

  afterEach(() => {
    vi.clearAllMocks();
  });

  describe('ApiError', () => {
    it('should have correct properties', () => {
      // ApiError constructor: (status, message, code?, details?)
      const error = new ApiError(404, 'Test error', undefined, { field: 'validation failed' });

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
  });

  describe('cvApi', () => {
    it('should get CV', async () => {
      const mockCV = { id: 'cv-123', title: 'My CV' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockCV),
      });

      const result = await api.cv.get();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result).toEqual(mockCV);
    });

    it('should create CV', async () => {
      const mockCV = { id: 'cv-123', title: 'New CV' };
      const createData = { title: 'New CV', content: {} };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockCV),
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

    it('should update CV', async () => {
      const mockCV = { id: 'cv-123', title: 'Updated CV' };
      const updateData = { title: 'Updated CV' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockCV),
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
        json: () => Promise.resolve({}),
      });

      await api.cv.delete('cv-123');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/cv/cv-123'),
        expect.objectContaining({ method: 'DELETE' })
      );
    });
  });

  describe('userApi', () => {
    it('should get current user', async () => {
      const mockUser = { id: 'user-123', email: 'test@example.com' };
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockUser),
      });

      const result = await api.user.getMe();

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/user/me'),
        expect.objectContaining({ method: 'GET' })
      );
      expect(result).toEqual(mockUser);
    });
  });
});
