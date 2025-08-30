#!/bin/bash

# Development script for marchive
# This script builds the frontend and starts the backend for local development

set -e

echo "🚀 Starting marchive development environment..."

# Check if we're in the right directory
if [ ! -f "go.mod" ] || [ ! -d "frontend" ]; then
    echo "❌ Error: Please run this script from the root directory of the project"
    exit 1
fi

# Check if pnpm is installed
if ! command -v pnpm &> /dev/null; then
    echo "❌ Error: pnpm is not installed. Please install it first:"
    echo "   npm install -g pnpm"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed. Please install it first."
    exit 1
fi

echo "📦 Building frontend..."
cd frontend

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "   Installing frontend dependencies..."
    pnpm install
fi

# Build the frontend
echo "   Building frontend assets..."
pnpm build

cd ..

echo "✅ Frontend built successfully!"

# Check if build directory exists
if [ ! -d "frontend/build" ]; then
    echo "❌ Error: Frontend build failed. Check the output above for errors."
    exit 1
fi

echo "🔧 Starting backend server..."
echo "   The application will be available at: http://localhost:9520"
echo "   Press Ctrl+C to stop the server"
echo ""

# Start the backend server
go run -tags=debug cmd/server/main.go
