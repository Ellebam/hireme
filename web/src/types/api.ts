/**
 * API Types
 * Matches Go API responses and requests
 */

import type { CVContent } from './cv';

// ============================================================================
// Common Types
// ============================================================================

/** API error response */
export interface ApiError {
  message: string;
  code?: string;
  details?: Record<string, string>;
}

/** Validation error response */
export interface ValidationError {
  field: string;
  message: string;
}

// ============================================================================
// User Types
// ============================================================================

export interface User {
  id: string;
  email: string;
  name: string;
  avatarUrl?: string;
  storageUsedBytes: number;
  storageLimitBytes: number;
  createdAt: string;
  updatedAt: string;
}

// ============================================================================
// CV Types
// ============================================================================

/** CV as returned by API */
export interface CV {
  id: string;
  title: string;
  schemaVersion: string;
  content: CVContent;
  createdAt: string;
  updatedAt: string;
}

/** Create CV request */
export interface CreateCVRequest {
  title: string;
  content: CVContent;
}

/** Update CV request */
export interface UpdateCVRequest {
  title?: string;
  content?: CVContent;
}

// ============================================================================
// Asset Types
// ============================================================================

export interface Asset {
  id: string;
  userId: string;
  filename: string;
  contentType: string;
  sizeBytes: number;
  checksum: string;
  width?: number;
  height?: number;
  createdAt: string;
}

export interface UploadAssetResponse {
  asset: Asset;
}

// ============================================================================
// Export Types
// ============================================================================

export type ExportFormat = 'pdf' | 'docx' | 'json';

export interface ExportRequest {
  format: ExportFormat;
  cvId: string;
}

// ============================================================================
// API Response Wrappers
// ============================================================================

/** Paginated list response */
export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}

/** Single item response */
export interface SingleResponse<T> {
  data: T;
}
