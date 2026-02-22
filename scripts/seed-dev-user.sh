#!/bin/bash
# Seed development user for auth bypass mode
# This creates a user record that matches AUTH_BYPASS_USER_ID

set -e

# Load environment variables
if [ -f .env.local ]; then
    export $(grep -v '^#' .env.local | xargs)
fi

# Default values if not set
DATABASE_URL="${DATABASE_URL:-postgres://hireme:hireme_dev_password@localhost:5433/hireme?sslmode=disable}"
DEV_USER_ID="${AUTH_BYPASS_USER_ID:-dev-user-001}"

echo "🌱 Seeding development user..."

# Using docker to run psql against the containerized postgres
docker exec -i hireme-postgres psql -U hireme -d hireme << EOF
-- Insert dev user if not exists
INSERT INTO users (
    id,
    external_id,
    provider,
    email,
    email_verified,
    display_name,
    tier,
    cv_limit,
    storage_limit_bytes,
    locale,
    created_at,
    updated_at
) VALUES (
    '${DEV_USER_ID}',
    'dev-external-id',
    'development',
    'dev@hireme.local',
    true,
    'Development User',
    'power',
    10,
    52428800,  -- 50MB for dev
    'en',
    NOW(),
    NOW()
) ON CONFLICT (id) DO UPDATE SET
    updated_at = NOW();
EOF

echo "✅ Development user seeded with ID: ${DEV_USER_ID}"
echo "   Email: dev@hireme.local"
echo "   Tier: power (10 CVs, 50MB storage)"
