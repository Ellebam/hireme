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

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

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
    return this.status === 422;
  }
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

    const response = await fetch(url, {
      ...options,
      headers: {
        'Content-Type': 'application/json',
        ...options?.headers,
      },
      credentials: 'include', // For cookies/auth
    });

    // Handle no-content responses
    if (response.status === 204) {
      return undefined as T;
    }

    // Parse response body
    const body = await response.json().catch(() => ({}));

    if (!response.ok) {
      throw new ApiError(
        response.status,
        body.message || 'Request failed',
        body.code,
        body.details
      );
    }

    return body as T;
  }

  async get<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'GET' });
  }

  async post<T>(endpoint: string, data?: unknown): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    });
  }

  async put<T>(endpoint: string, data: unknown): Promise<T> {
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
    const response = await fetch(url, {
      method: 'POST',
      body: formData,
      credentials: 'include',
    });

    const body = await response.json().catch(() => ({}));

    if (!response.ok) {
      throw new ApiError(
        response.status,
        body.message || 'Upload failed',
        body.code
      );
    }

    return body as T;
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
  getMe: () => client.get<User>('/api/v1/user/me'),
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
    const response = await fetch(url, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ cvId }),
      credentials: 'include',
    });

    if (!response.ok) {
      const body = await response.json().catch(() => ({}));
      throw new ApiError(
        response.status,
        body.message || 'Export failed',
        body.code
      );
    }

    return response.blob();
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
