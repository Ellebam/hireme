-- HireMe Initial Schema
-- Migration: 000001_init
-- Description: Creates all base tables for the application

-- ══════════════════════════════════════════════════════════════════════════════
-- EXTENSIONS
-- ══════════════════════════════════════════════════════════════════════════════

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ══════════════════════════════════════════════════════════════════════════════
-- USERS TABLE
-- ══════════════════════════════════════════════════════════════════════════════

CREATE TABLE users (
    id TEXT PRIMARY KEY,
    external_id TEXT NOT NULL,
    provider TEXT NOT NULL,
    email TEXT NOT NULL,
    email_verified BOOLEAN DEFAULT false,
    display_name TEXT,
    tier TEXT DEFAULT 'free' CHECK (tier IN ('free', 'pro', 'power')),
    cv_limit INTEGER DEFAULT 1,
    storage_limit_bytes BIGINT DEFAULT 5242880,
    storage_used_bytes BIGINT DEFAULT 0,
    locale TEXT DEFAULT 'en' CHECK (locale IN ('en', 'de')),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    UNIQUE(provider, external_id)
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_provider ON users(provider);

-- ══════════════════════════════════════════════════════════════════════════════
-- CVS TABLE
-- ══════════════════════════════════════════════════════════════════════════════

CREATE TABLE cvs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title TEXT NOT NULL DEFAULT 'My CV',
    schema_version TEXT NOT NULL DEFAULT '1.0.0',
    content JSONB NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT valid_content CHECK (jsonb_typeof(content) = 'object')
);

CREATE INDEX idx_cvs_user_id ON cvs(user_id);
CREATE INDEX idx_cvs_updated_at ON cvs(updated_at DESC);

-- ══════════════════════════════════════════════════════════════════════════════
-- ASSETS TABLE
-- ══════════════════════════════════════════════════════════════════════════════

CREATE TABLE assets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    original_filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes INTEGER NOT NULL,
    storage_path TEXT NOT NULL,
    storage_backend TEXT NOT NULL DEFAULT 'local',
    checksum TEXT NOT NULL,
    width INTEGER,
    height INTEGER,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMPTZ DEFAULT NOW(),
    
    CONSTRAINT valid_mime CHECK (mime_type IN ('image/jpeg', 'image/png', 'image/webp')),
    CONSTRAINT valid_backend CHECK (storage_backend IN ('local', 'r2'))
);

CREATE INDEX idx_assets_user_id ON assets(user_id);
CREATE INDEX idx_assets_checksum ON assets(checksum);

-- ══════════════════════════════════════════════════════════════════════════════
-- EXPORT JOBS TABLE
-- ══════════════════════════════════════════════════════════════════════════════

CREATE TABLE export_jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    cv_id UUID NOT NULL REFERENCES cvs(id) ON DELETE CASCADE,
    format TEXT NOT NULL CHECK (format IN ('pdf', 'docx', 'json', 'yaml')),
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'completed', 'failed')),
    result_path TEXT,
    error_message TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ
);

CREATE INDEX idx_export_jobs_user_id ON export_jobs(user_id);
CREATE INDEX idx_export_jobs_status ON export_jobs(status) WHERE status IN ('pending', 'processing');
CREATE INDEX idx_export_jobs_created_at ON export_jobs(created_at DESC);

-- ══════════════════════════════════════════════════════════════════════════════
-- AUDIT LOG TABLE
-- ══════════════════════════════════════════════════════════════════════════════

CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id TEXT REFERENCES users(id) ON DELETE SET NULL,
    action TEXT NOT NULL,
    entity_type TEXT NOT NULL,
    entity_id TEXT,
    old_value JSONB,
    new_value JSONB,
    ip_address INET,
    user_agent TEXT,
    request_id TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX idx_audit_log_entity ON audit_log(entity_type, entity_id);
CREATE INDEX idx_audit_log_created_at ON audit_log(created_at DESC);

-- Partition by month for production (optional optimization)
-- This is a simple table for MVP, can be partitioned later

-- ══════════════════════════════════════════════════════════════════════════════
-- FUNCTIONS
-- ══════════════════════════════════════════════════════════════════════════════

-- Auto-update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply to users table
CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- Apply to cvs table
CREATE TRIGGER update_cvs_updated_at
    BEFORE UPDATE ON cvs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- ══════════════════════════════════════════════════════════════════════════════
-- ROW LEVEL SECURITY (Prepared for production)
-- ══════════════════════════════════════════════════════════════════════════════

-- These policies enforce tenant isolation when RLS is enabled
-- Enable with: ALTER TABLE cvs ENABLE ROW LEVEL SECURITY;

-- Note: RLS requires setting app.current_user_id via SET LOCAL in transactions
-- This is typically done in middleware before each request

CREATE POLICY cvs_user_isolation ON cvs
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true));

CREATE POLICY assets_user_isolation ON assets
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true));

CREATE POLICY export_jobs_user_isolation ON export_jobs
    FOR ALL
    USING (user_id = current_setting('app.current_user_id', true));

-- RLS is NOT enabled by default for easier development
-- Enable in production with:
-- ALTER TABLE cvs ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE assets ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE export_jobs ENABLE ROW LEVEL SECURITY;

-- ══════════════════════════════════════════════════════════════════════════════
-- COMMENTS
-- ══════════════════════════════════════════════════════════════════════════════

COMMENT ON TABLE users IS 'User accounts linked to external identity providers';
COMMENT ON TABLE cvs IS 'CV documents stored as JSON with schema versioning';
COMMENT ON TABLE assets IS 'User-uploaded files (portraits, etc.)';
COMMENT ON TABLE export_jobs IS 'Async export job queue for PDF/DOCX generation';
COMMENT ON TABLE audit_log IS 'Immutable log of all data mutations';

COMMENT ON COLUMN users.tier IS 'Subscription tier: free (1 CV), pro (5 CVs), power (unlimited)';
COMMENT ON COLUMN users.cv_limit IS 'Maximum number of CVs allowed for this user';
COMMENT ON COLUMN cvs.content IS 'CV data following schemas/cv-schema.json';
COMMENT ON COLUMN assets.checksum IS 'SHA-256 hash for deduplication and integrity';
