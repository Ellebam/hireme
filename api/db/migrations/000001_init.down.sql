-- HireMe Initial Schema Rollback
-- Migration: 000001_init
-- Description: Drops all base tables

-- Drop triggers first
DROP TRIGGER IF EXISTS update_cvs_updated_at ON cvs;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS audit_log;
DROP TABLE IF EXISTS export_jobs;
DROP TABLE IF EXISTS assets;
DROP TABLE IF EXISTS cvs;
DROP TABLE IF EXISTS users;

-- Extensions are left intact (they might be used by other schemas)
