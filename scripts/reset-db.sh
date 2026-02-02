#!/bin/bash
# Reset the database by dropping and recreating it
# WARNING: This deletes all data!

set -e

echo "⚠️  Resetting database..."

# Drop and recreate database
docker exec -i hireme-postgres psql -U hireme -d postgres << EOF
-- Terminate existing connections
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'hireme'
  AND pid <> pg_backend_pid();

-- Drop and recreate
DROP DATABASE IF EXISTS hireme;
CREATE DATABASE hireme OWNER hireme;
EOF

echo "✅ Database reset complete"
echo "   Run 'task db:migrate' to apply migrations"
echo "   Run 'task db:seed' to seed dev user"
