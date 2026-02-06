/**
 * API Client
 * Typed HTTP client for the HireMe Go backend
 */

import type {
  CV,
  CreateCVRequest,
  UpdateCVRequest,
  User,
  Asset,
  ExportFormat,
} from '@/types/api';
import { logger } from '@/lib/logger';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
const REQUEST_TIMEOUT_MS = 30000; // 30 seconds

// ============================================================================
// Error Handling
// ============================================================================

export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
    public code?: string,
    public details?: Record<string, string>
  ) {
    super(message);
    this.name = 'ApiError';
  }

  get isNotFound(): boolean {
    return this.status === 404;
  }

  get isUnauthorized(): boolean {
    return this.status === 401;
  }

  get isForbidden(): boolean {
    return this.status === 403;
  }

  get isValidationError(): boolean {
    return this.status === 422 || this.status === 400;
  }
}

// ============================================================================
// API Response Types
// ============================================================================

/** Standard API response wrapper from backend */
interface ApiResponse<T> {
  data?: T;
  error?: {
    code: string;
    message: string;
    field?: string;
  };
}

// ============================================================================
// Base Client
// ============================================================================

class ApiClient {
  private baseUrl: string;

  constructor(baseUrl: string) {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options?: RequestInit
  ): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    const method = options?.method || 'GET';

    logger.debug('API', `${method} ${endpoint}`);

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS);

    try {
      const response = await fetch(url, {
        ...options,
        signal: controller.signal,
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        credentials: 'include', // For cookies/auth
      });

      // Handle no-content responses
      if (response.status === 204) {
        logger.debug('API', `${method} ${endpoint} -> 204 No Content`);
        return undefined as T;
      }

      // Parse response body
      const body: ApiResponse<T> = await response.json().catch(() => ({}));

      logger.debug('API', `${method} ${endpoint} -> ${response.status}`, {
        hasData: !!body.data,
        hasError: !!body.error,
      });

      // Check for error response
      if (!response.ok || body.error) {
        const error = body.error;
        const errorMessage = error?.message || 'Request failed';
        const errorCode = error?.code;

        logger.error('API', `${method} ${endpoint} failed: ${errorMessage}`, {
          status: response.status,
          code: errorCode,
        });

        throw new ApiError(response.status, errorMessage, errorCode);
      }

      // Unwrap the data from the response wrapper
      if (body.data !== undefined) {
        return body.data;
      }

      // Fallback: return body directly if no wrapper (shouldn't happen with our API)
      logger.warn('API', `Response missing data wrapper for ${endpoint}`);
      return body as unknown as T;
    } catch (err) {
      // Re-throw ApiErrors as-is
      if (err instanceof ApiError) {
        throw err;
      }

      // Handle abort/timeout errors
      if (err instanceof DOMException && err.name === 'AbortError') {
        logger.error('API', `Request timeout for ${method} ${endpoint}`);
        throw new ApiError(0, 'Request timed out - please try again');
      }

      // Handle network errors
      logger.error('API', `Network error for ${method} ${endpoint}`);
      throw new ApiError(0, 'Network error - please check your connection');
    } finally {
      clearTimeout(timeoutId);
    }
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    logger.debug('API', `POST ${endpoint}`, { hasPayload: data !== undefined });
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data: unknown): Promise<T> {
    logger.debug('API', `PUT ${endpoint}`, { hasPayload: true });
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: JSON.stringify(data),
    });
  }

  async delete<T = void>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' });
  }

  async upload<T>(endpoint: string, file: File): Promise<T> {
    const formData = new FormData();
    formData.append('file', file);

    const url = `${this.baseUrl}${endpoint}`;
    logger.debug('API', `UPLOAD ${endpoint}`, { filename: file.name, size: file.size });

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS * 2); // Double timeout for uploads

    try {
      const response = await fetch(url, {
        method: 'POST',
        body: formData,
        credentials: 'include',
        signal: controller.signal,
      });

      const body: ApiResponse<T> = await response.json().catch(() => ({}));

      if (!response.ok || body.error) {
        const error = body.error;
        logger.error('API', `Upload failed: ${error?.message || 'Unknown error'}`);
        throw new ApiError(
          response.status,
          error?.message || 'Upload failed',
          error?.code
        );
      }

      // Unwrap the data
      if (body.data !== undefined) {
        return body.data;
      }

      return body as unknown as T;
    } catch (err) {
      if (err instanceof ApiError) throw err;
      if (err instanceof DOMException && err.name === 'AbortError') {
        throw new ApiError(0, 'Upload timed out - please try again');
      }
      logger.error('API', 'Upload network error');
      throw new ApiError(0, 'Upload failed - network error');
    } finally {
      clearTimeout(timeoutId);
    }
  }
}

// ============================================================================
// API Instance
// ============================================================================

const client = new ApiClient(API_BASE);

// ============================================================================
// User API
// ============================================================================

export const userApi = {
  /** Get current authenticated user */
  getMe: () => client.get<User>('/api/v1/users/me'),
};

// ============================================================================
// CV API
// ============================================================================

export const cvApi = {
  /** Get current user's CV */
  get: () => client.get<CV>('/api/v1/cv'),

  /** Create a new CV */
  create: (data: CreateCVRequest) => client.post<CV>('/api/v1/cv', data),

  /** Update an existing CV */
  update: (id: string, data: UpdateCVRequest) =>
    client.put<CV>(`/api/v1/cv/${id}`, data),

  /** Delete a CV */
  delete: (id: string) => client.delete(`/api/v1/cv/${id}`),
};

// ============================================================================
// Asset API
// ============================================================================

export const assetApi = {
  /** Upload an asset (image) */
  upload: (file: File) =>
    client.upload<{ asset: Asset }>('/api/v1/assets', file),

  /** Get asset metadata */
  get: (id: string) => client.get<Asset>(`/api/v1/assets/${id}`),

  /** Get asset file content (returns URL for direct use) */
  getFileUrl: (id: string) => `${API_BASE}/api/v1/assets/${id}/file`,

  /** Delete an asset */
  delete: (id: string) => client.delete(`/api/v1/assets/${id}`),
};

// ============================================================================
// Export API
// ============================================================================

export const exportApi = {
  /** Export CV to specified format */
  export: async (cvId: string, format: ExportFormat): Promise<Blob> => {
    const url = `${API_BASE}/api/v1/export/${format}`;
    logger.debug('API', `Export ${format}`, { cvId });

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), REQUEST_TIMEOUT_MS * 2);

    try {
      const response = await fetch(url, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ cvId }),
        credentials: 'include',
        signal: controller.signal,
      });

      if (!response.ok) {
        const body = await response.json().catch(() => ({}));
        const error = body.error || body;
        logger.error('API', `Export failed: ${error.message || 'Unknown error'}`);
        throw new ApiError(
          response.status,
          error.message || 'Export failed',
          error.code
        );
      }

      return response.blob();
    } finally {
      clearTimeout(timeoutId);
    }
  },

  /** Get download URL for export */
  getDownloadUrl: (format: ExportFormat) =>
    `${API_BASE}/api/v1/export/${format}`,
};

// ============================================================================
// Convenience Export
// ============================================================================

export const api = {
  user: userApi,
  cv: cvApi,
  asset: assetApi,
  export: exportApi,
};

export default api;
