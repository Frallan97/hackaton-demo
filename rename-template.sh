#!/bin/bash

# Template Renamer Script
# Usage: ./rename-template.sh <new-app-name>
# Example: ./rename-template.sh inventory-app

set -e

if [ $# -eq 0 ]; then
    echo "Usage: $0 <new-app-name>"
    echo "Example: $0 inventory-app"
    exit 1
fi

NEW_APP_NAME=$1
OLD_APP_NAME="hackaton-demo"
OLD_APP_NAME_FRONTEND="hackaton-demo-frontend"
OLD_APP_NAME_BACKEND="hackaton-demo-backend"
NEW_APP_NAME_FRONTEND="${NEW_APP_NAME}-frontend"
NEW_APP_NAME_BACKEND="${NEW_APP_NAME}-backend"

echo "Renaming template from '$OLD_APP_NAME' to '$NEW_APP_NAME'"
echo "Frontend: $OLD_APP_NAME_FRONTEND -> $NEW_APP_NAME_FRONTEND"
echo "Backend: $OLD_APP_NAME_BACKEND -> $NEW_APP_NAME_BACKEND"
echo ""

# Function to replace text in files
replace_in_files() {
    local pattern=$1
    local replacement=$2
    local files=$3
    
    echo "Replacing '$pattern' with '$replacement' in $files files..."
    find . -type f \( -name "*.go" -o -name "*.yaml" -o -name "*.yml" -o -name "*.json" -o -name "*.md" -o -name "*.sh" -o -name "Dockerfile" -o -name "package.json" -o -name "go.mod" \) \
        -not -path "./.git/*" \
        -not -path "./node_modules/*" \
        -exec sed -i "s/$pattern/$replacement/g" {} \;
}

# Function to rename directories
rename_directories() {
    local old_name=$1
    local new_name=$2
    
    if [ -d "$old_name" ]; then
        echo "Renaming directory: $old_name -> $new_name"
        mv "$old_name" "$new_name"
    fi
}

# Replace all text references
replace_in_files "$OLD_APP_NAME_FRONTEND" "$NEW_APP_NAME_FRONTEND" "package.json, Dockerfile, CI files"
replace_in_files "$OLD_APP_NAME_BACKEND" "$NEW_APP_NAME_BACKEND" "Go files, CI files"
replace_in_files "$OLD_APP_NAME" "$NEW_APP_NAME" "Chart.yaml, values.yaml, README files"

# Replace GitHub repository references
replace_in_files "github.com/frallan97/$OLD_APP_NAME_BACKEND" "github.com/frallan97/$NEW_APP_NAME_BACKEND" "Go files"
replace_in_files "ghcr.io/frallan97/$OLD_APP_NAME_FRONTEND" "ghcr.io/frallan97/$NEW_APP_NAME_FRONTEND" "values.yaml, CI files"
replace_in_files "ghcr.io/frallan97/$OLD_APP_NAME_BACKEND" "ghcr.io/frallan97/$NEW_APP_NAME_BACKEND" "values.yaml, CI files"

# Replace domain references
replace_in_files "$OLD_APP_NAME.web.franssjostrom.com" "$NEW_APP_NAME.web.franssjostrom.com" "values.yaml"

# Replace JWT issuer
replace_in_files "Issuer:    \"$OLD_APP_NAME\"" "Issuer:    \"$NEW_APP_NAME\"" "JWT service"

# Rename chart directory (only if it's different)
if [ "$OLD_APP_NAME" != "$NEW_APP_NAME" ]; then
    rename_directories "charts/$OLD_APP_NAME" "charts/$NEW_APP_NAME"
fi

# Update CI workflow paths
if [ -f ".github/workflows/ci.yaml" ]; then
    echo "Updating CI workflow paths..."
    sed -i "s|charts/$OLD_APP_NAME/|charts/$NEW_APP_NAME/|g" .github/workflows/ci.yaml
fi

# Update docker-compose service names if they exist
if [ -f "docker-compose.yml" ]; then
    echo "Updating docker-compose service names..."
    sed -i "s/$OLD_APP_NAME_FRONTEND/$NEW_APP_NAME_FRONTEND/g" docker-compose.yml
    sed -i "s/$OLD_APP_NAME_BACKEND/$NEW_APP_NAME_BACKEND/g" docker-compose.yml
fi

echo ""
echo "✅ Template renamed successfully!"
echo ""
echo "Next steps:"
echo "1. Update your GitHub repository name to '$NEW_APP_NAME'"
echo "2. Update your GitHub Container Registry repository names:"
echo "   - $NEW_APP_NAME_FRONTEND"
echo "   - $NEW_APP_NAME_BACKEND"
echo "3. Update your domain DNS settings if using custom domains"
echo "4. Review and update any environment-specific configurations"
echo ""
echo "Files that were updated:"
echo "- Go module paths in backend/"
echo "- Package names in frontend/"
echo "- Helm chart configuration"
echo "- CI/CD pipeline configuration"
echo "- Documentation files"
echo "- Docker configurations" 