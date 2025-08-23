#!/bin/bash

# Hackathon Demo Setup Script
# This script sets up the development environment quickly

set -e

echo "🚀 Setting up Hackathon Demo..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp env.example .env
    echo "⚠️  Please edit .env with your Google OAuth credentials!"
    echo "   Required: GOOGLE_CLIENT_ID and GOOGLE_CLIENT_SECRET"
    echo "   You can get these from: https://console.cloud.google.com/"
    echo ""
    read -p "Press Enter after you've updated .env..."
else
    echo "✅ .env file already exists"
fi

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

echo "🐳 Starting database services..."
docker-compose up db redis nats -d

echo "⏳ Waiting for services to be ready..."
sleep 10

echo "🔍 Checking service status..."
docker-compose ps

echo ""
echo "🎉 Setup complete! Next steps:"
echo ""
echo "1. Backend development:"
echo "   cd backend && go run main.go"
echo ""
echo "2. Frontend development:"
echo "   cd frontend && bun install && bun run dev"
echo ""
echo "3. Access your app:"
echo "   Frontend: http://localhost:3000"
echo "   Backend:  http://localhost:8080"
echo "   API docs: http://localhost:8080/docs"
echo ""
echo "4. Stop services:"
echo "   docker-compose down"
echo ""
echo "Happy hacking! 🚀" 