/**
 * Frontend Logger
 *
 * Configurable logging with levels for easier debugging.
 * Set LOG_LEVEL in environment or localStorage to control output.
 *
 * Levels: debug < info < warn < error < none
 *
 * Usage:
 *   import { logger } from '@/lib/logger';
 *   logger.debug('API', 'Fetching CV');
 *   logger.info('Store', 'CV loaded', { id: cv.id });
 *   logger.error('API', 'Request failed', error);
 *
 * Enable debug logging:
 *   localStorage.setItem('LOG_LEVEL', 'debug');
 *   location.reload();
 */

export type LogLevel = 'debug' | 'info' | 'warn' | 'error' | 'none';

const LOG_LEVELS: Record<LogLevel, number> = {
  debug: 0,
  info: 1,
  warn: 2,
  error: 3,
  none: 4,
};

function getLogLevel(): LogLevel {
  // Check localStorage first (for browser debugging)
  if (typeof window !== 'undefined') {
    const stored = localStorage.getItem('LOG_LEVEL');
    if (stored && stored in LOG_LEVELS) {
      return stored as LogLevel;
    }
  }

  // Check environment variable
  const envLevel = process.env.NEXT_PUBLIC_LOG_LEVEL;
  if (envLevel && envLevel in LOG_LEVELS) {
    return envLevel as LogLevel;
  }

  // Default: info in development, warn in production
  return process.env.NODE_ENV === 'development' ? 'info' : 'warn';
}

function shouldLog(level: LogLevel): boolean {
  const currentLevel = getLogLevel();
  return LOG_LEVELS[level] >= LOG_LEVELS[currentLevel];
}

function formatMessage(category: string, message: string): string {
  const timestamp = new Date().toISOString().substring(11, 23);
  return `[${timestamp}] [${category}] ${message}`;
}

export const logger = {
  /**
   * Debug level - detailed information for troubleshooting
   */
  debug(category: string, message: string, data?: unknown): void {
    if (!shouldLog('debug')) return;
    const formatted = formatMessage(category, message);
    if (data !== undefined) {
      console.debug(formatted, data);
    } else {
      console.debug(formatted);
    }
  },

  /**
   * Info level - general operational information
   */
  info(category: string, message: string, data?: unknown): void {
    if (!shouldLog('info')) return;
    const formatted = formatMessage(category, message);
    if (data !== undefined) {
      console.info(formatted, data);
    } else {
      console.info(formatted);
    }
  },

  /**
   * Warn level - potential issues that don't prevent operation
   */
  warn(category: string, message: string, data?: unknown): void {
    if (!shouldLog('warn')) return;
    const formatted = formatMessage(category, message);
    if (data !== undefined) {
      console.warn(formatted, data);
    } else {
      console.warn(formatted);
    }
  },

  /**
   * Error level - errors that affect operation
   */
  error(category: string, message: string, data?: unknown): void {
    if (!shouldLog('error')) return;
    const formatted = formatMessage(category, message);
    if (data !== undefined) {
      console.error(formatted, data);
    } else {
      console.error(formatted);
    }
  },

  /**
   * Get current log level
   */
  getLevel(): LogLevel {
    return getLogLevel();
  },

  /**
   * Set log level (persists in localStorage)
   */
  setLevel(level: LogLevel): void {
    if (typeof window !== 'undefined') {
      localStorage.setItem('LOG_LEVEL', level);
    }
  },
};

export default logger;
