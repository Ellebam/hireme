#!/bin/bash
# Generate TypeScript types from JSON Schema
# This ensures frontend types match the backend CV schema

set -e

SCHEMA_FILE="schemas/cv-schema.json"
OUTPUT_FILE="web/src/types/cv-generated.ts"

echo "📝 Generating TypeScript types from JSON Schema..."

# Check if json-schema-to-typescript is installed
if ! command -v npx &> /dev/null; then
    echo "❌ npx not found. Please install Node.js"
    exit 1
fi

# Check if schema file exists
if [ ! -f "$SCHEMA_FILE" ]; then
    echo "❌ Schema file not found: $SCHEMA_FILE"
    exit 1
fi

# Create output directory if needed
mkdir -p "$(dirname "$OUTPUT_FILE")"

# Generate types using json-schema-to-typescript
cd web && npx json-schema-to-typescript "../$SCHEMA_FILE" -o "../$OUTPUT_FILE" --bannerComment "/* eslint-disable */
/**
 * AUTO-GENERATED FILE - DO NOT EDIT
 * Generated from: schemas/cv-schema.json
 * Run 'task generate:types' to regenerate
 */"

echo "✅ TypeScript types generated: $OUTPUT_FILE"
