#!/bin/bash
# Initial development environment setup
# Run this once after cloning the repository

set -e

echo "🚀 Setting up HireMe development environment..."
echo ""

# Check dependencies
echo "📋 Checking dependencies..."

check_command() {
    if command -v $1 &> /dev/null; then
        echo "  ✅ $1 found: $($1 --version 2>/dev/null | head -n 1)"
        return 0
    else
        echo "  ❌ $1 not found"
        return 1
    fi
}

MISSING_DEPS=0
check_command "go" || MISSING_DEPS=1
check_command "node" || MISSING_DEPS=1
check_command "docker" || MISSING_DEPS=1
check_command "task" || { echo "  ⚠️  task not found - install from https://taskfile.dev"; MISSING_DEPS=1; }

if [ $MISSING_DEPS -eq 1 ]; then
    echo ""
    echo "❌ Missing required dependencies. Please install them and try again."
    exit 1
fi

echo ""

# Create .env.local
if [ ! -f .env.local ]; then
    echo "📝 Creating .env.local from .env.example..."
    cp .env.example .env.local
    echo "  ✅ Created .env.local"
else
    echo "📝 .env.local already exists, skipping..."
fi

echo ""

# Create data directories
echo "📁 Creating data directories..."
mkdir -p data/uploads
mkdir -p data/postgres
echo "  ✅ Created data/uploads"
echo "  ✅ Created data/postgres"

echo ""

# Make scripts executable
echo "🔧 Making scripts executable..."
chmod +x scripts/*.sh
echo "  ✅ Scripts are now executable"

echo ""

# Start infrastructure
echo "🐳 Starting infrastructure (PostgreSQL, Gotenberg)..."
task infra:up

echo ""

# Install Go dependencies
echo "📦 Installing Go dependencies..."
cd api && go mod download && go mod tidy && cd ..
echo "  ✅ Go dependencies installed"

echo ""

# Install Node dependencies
echo "📦 Installing Node.js dependencies..."
cd web && npm install && cd ..
echo "  ✅ Node dependencies installed"

echo ""

# Run migrations
echo "🗃️  Running database migrations..."
task db:migrate

echo ""

# Seed dev user
echo "🌱 Seeding development user..."
task db:seed

echo ""
echo "════════════════════════════════════════════════════════════════════"
echo "✅ Setup complete!"
echo "════════════════════════════════════════════════════════════════════"
echo ""
echo "Next steps:"
echo "  1. Review .env.local and adjust settings if needed"
echo "  2. Start the API server:     task api:dev"
echo "  3. Start the web server:     task web:dev"
echo ""
echo "Development servers will run at:"
echo "  • API:  http://localhost:8080"
echo "  • Web:  http://localhost:3000"
echo ""
echo "Auth bypass is enabled by default. See .env.local to configure."
echo ""
