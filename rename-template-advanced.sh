#!/bin/bash

# Advanced Template Renamer Script
# Usage: ./rename-template-advanced.sh <new-app-name> [github-username] [domain]
# Example: ./rename-template-advanced.sh inventory-app myusername mydomain.com

set -e

if [ $# -lt 1 ]; then
    echo "Usage: $0 <new-app-name> [github-username] [domain]"
    echo "Example: $0 inventory-app myusername mydomain.com"
    echo ""
    echo "If github-username is not provided, will use 'frallan97'"
    echo "If domain is not provided, will use 'franssjostrom.com'"
    exit 1
fi

NEW_APP_NAME=$1
GITHUB_USERNAME=${2:-"frallan97"}
DOMAIN=${3:-"franssjostrom.com"}

OLD_APP_NAME="hackaton-demo"
OLD_APP_NAME_FRONTEND="hackaton-demo-frontend"
OLD_APP_NAME_BACKEND="hackaton-demo-backend"
NEW_APP_NAME_FRONTEND="${NEW_APP_NAME}-frontend"
NEW_APP_NAME_BACKEND="${NEW_APP_NAME}-backend"

OLD_GITHUB_USERNAME="frallan97"
OLD_DOMAIN="franssjostrom.com"

echo "Advanced Template Renamer"
echo "========================"
echo "New App Name: $NEW_APP_NAME"
echo "GitHub Username: $GITHUB_USERNAME"
echo "Domain: $DOMAIN"
echo ""
echo "Frontend: $OLD_APP_NAME_FRONTEND -> $NEW_APP_NAME_FRONTEND"
echo "Backend: $OLD_APP_NAME_BACKEND -> $NEW_APP_NAME_BACKEND"
echo "GitHub: $OLD_GITHUB_USERNAME -> $GITHUB_USERNAME"
echo "Domain: $OLD_DOMAIN -> $DOMAIN"
echo ""

# Function to replace text in files
replace_in_files() {
    local pattern=$1
    local replacement=$2
    local description=$3
    
    echo "Replacing '$pattern' with '$replacement' in $description..."
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

# Replace app name references
replace_in_files "$OLD_APP_NAME_FRONTEND" "$NEW_APP_NAME_FRONTEND" "frontend package names"
replace_in_files "$OLD_APP_NAME_BACKEND" "$NEW_APP_NAME_BACKEND" "backend module names"
replace_in_files "$OLD_APP_NAME" "$NEW_APP_NAME" "general app references"

# Replace GitHub username references
replace_in_files "github.com/$OLD_GITHUB_USERNAME/$OLD_APP_NAME_BACKEND" "github.com/$GITHUB_USERNAME/$NEW_APP_NAME_BACKEND" "Go module paths"
replace_in_files "ghcr.io/$OLD_GITHUB_USERNAME/$OLD_APP_NAME_FRONTEND" "ghcr.io/$GITHUB_USERNAME/$NEW_APP_NAME_FRONTEND" "frontend Docker images"
replace_in_files "ghcr.io/$OLD_GITHUB_USERNAME/$OLD_APP_NAME_BACKEND" "ghcr.io/$GITHUB_USERNAME/$NEW_APP_NAME_BACKEND" "backend Docker images"

# Replace domain references
replace_in_files "$OLD_APP_NAME.web.$OLD_DOMAIN" "$NEW_APP_NAME.web.$DOMAIN" "domain references"
# Note: We don't replace general domain references as it could break other URLs

# Replace JWT issuer
replace_in_files "Issuer:    \"$OLD_APP_NAME\"" "Issuer:    \"$NEW_APP_NAME\"" "JWT issuer"

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

# Update template configuration file
if [ -f "template-config.yaml" ]; then
    echo "Updating template configuration..."
    sed -i "s/base_name: \"$OLD_APP_NAME\"/base_name: \"$NEW_APP_NAME\"/g" template-config.yaml
    sed -i "s/frontend_name: \"$OLD_APP_NAME_FRONTEND\"/frontend_name: \"$NEW_APP_NAME_FRONTEND\"/g" template-config.yaml
    sed -i "s/backend_name: \"$OLD_APP_NAME_BACKEND\"/backend_name: \"$NEW_APP_NAME_BACKEND\"/g" template-config.yaml
    sed -i "s/username: \"$OLD_GITHUB_USERNAME\"/username: \"$GITHUB_USERNAME\"/g" template-config.yaml
    sed -i "s/base: \"$OLD_DOMAIN\"/base: \"$DOMAIN\"/g" template-config.yaml
    sed -i "s/subdomain: \"$OLD_APP_NAME\"/subdomain: \"$NEW_APP_NAME\"/g" template-config.yaml
    sed -i "s/full_domain: \"$OLD_APP_NAME.web.$OLD_DOMAIN\"/full_domain: \"$NEW_APP_NAME.web.$DOMAIN\"/g" template-config.yaml
    sed -i "s/issuer: \"$OLD_APP_NAME\"/issuer: \"$NEW_APP_NAME\"/g" template-config.yaml
    sed -i "s/frontend_image: \"ghcr.io\/$OLD_GITHUB_USERNAME\/$OLD_APP_NAME_FRONTEND\"/frontend_image: \"ghcr.io\/$GITHUB_USERNAME\/$NEW_APP_NAME_FRONTEND\"/g" template-config.yaml
    sed -i "s/backend_image: \"ghcr.io\/$OLD_GITHUB_USERNAME\/$OLD_APP_NAME_BACKEND\"/backend_image: \"ghcr.io\/$GITHUB_USERNAME\/$NEW_APP_NAME_BACKEND\"/g" template-config.yaml
    sed -i "s/chart_name: \"$OLD_APP_NAME\"/chart_name: \"$NEW_APP_NAME\"/g" template-config.yaml
    sed -i "s/chart_path: \"charts\/$OLD_APP_NAME\"/chart_path: \"charts\/$NEW_APP_NAME\"/g" template-config.yaml
fi

echo ""
echo "✅ Template renamed successfully!"
echo ""
echo "Next steps:"
echo "1. Update your GitHub repository name to '$NEW_APP_NAME'"
echo "2. Update your GitHub Container Registry repository names:"
echo "   - $NEW_APP_NAME_FRONTEND"
echo "   - $NEW_APP_NAME_BACKEND"
echo "3. Update your domain DNS settings: $NEW_APP_NAME.web.$DOMAIN"
echo "4. Update your Google OAuth configuration with new redirect URIs"
echo "5. Review and update any environment-specific configurations"
echo ""
echo "Files that were updated:"
echo "- Go module paths in backend/"
echo "- Package names in frontend/"
echo "- Helm chart configuration"
echo "- CI/CD pipeline configuration"
echo "- Documentation files"
echo "- Docker configurations"
echo "- Template configuration"
echo ""
echo "GitHub Container Registry images:"
echo "- ghcr.io/$GITHUB_USERNAME/$NEW_APP_NAME_FRONTEND"
echo "- ghcr.io/$GITHUB_USERNAME/$NEW_APP_NAME_BACKEND" 