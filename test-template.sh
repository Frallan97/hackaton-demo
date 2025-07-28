#!/bin/bash

# Test Template Renaming Script
# This script tests the template renaming functionality

set -e

echo "🧪 Testing Template Renaming System"
echo "==================================="

# Test 1: Basic rename functionality
echo ""
echo "Test 1: Basic rename (hackaton-demo -> inventory-app)"
echo "-----------------------------------------------------"

# Create a test directory structure
mkdir -p test-template/charts/hackaton-demo
mkdir -p test-template/backend
mkdir -p test-template/frontend

# Create test files with references
cat > test-template/backend/go.mod << EOF
module github.com/frallan97/hackaton-demo-backend
EOF

cat > test-template/frontend/package.json << EOF
{
  "name": "hackaton-demo-frontend"
}
EOF

cat > test-template/charts/hackaton-demo/values.yaml << EOF
frontend:
  repository: ghcr.io/frallan97/hackaton-demo-frontend
backend:
  repository: ghcr.io/frallan97/hackaton-demo-backend
ingress:
  hosts:
    - host: hackaton-demo.web.hackaton-demo.web.franssjostrom.com
EOF

# Copy the rename script
cp rename-template.sh test-template/
cd test-template

# Run the rename script
echo "Running rename script..."
./rename-template.sh inventory-app

# Verify the changes
echo ""
echo "Verifying changes..."
echo "-------------------"

echo "✅ Go module:"
grep "inventory-app-backend" backend/go.mod || echo "❌ Go module not updated"

echo "✅ Frontend package:"
grep "inventory-app-frontend" frontend/package.json || echo "❌ Frontend package not updated"

echo "✅ Docker images:"
grep "inventory-app-frontend" charts/inventory-app/values.yaml || echo "❌ Frontend image not updated"
grep "inventory-app-backend" charts/inventory-app/values.yaml || echo "❌ Backend image not updated"

echo "✅ Domain:"
grep "inventory-app.web.franssjostrom.com" charts/inventory-app/values.yaml || echo "❌ Domain not updated"

echo "✅ Chart directory renamed:"
if [ -d "charts/inventory-app" ]; then
    echo "✅ charts/inventory-app exists"
else
    echo "❌ charts/inventory-app does not exist"
fi

# Cleanup
cd ..
rm -rf test-template

echo ""
echo "🎉 Template test completed!"
echo ""
echo "If all checks passed, your template system is working correctly."
echo "You can now use it to create new applications:"
echo ""
echo "  ./rename-template.sh your-app-name" 