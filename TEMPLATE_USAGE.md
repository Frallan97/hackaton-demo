# Template Usage Guide

This repository serves as a template for creating new React + Go applications with full CI/CD pipeline, Kubernetes deployment, and OAuth authentication.

## Quick Start

### 1. Use as GitHub Template

1. Click "Use this template" on GitHub
2. Choose "Create a new repository"
3. Name your new repository (e.g., `inventory-app`)
4. Clone the new repository locally

### 2. Rename the Template

Use the provided script to rename all references:

```bash
./rename-template.sh inventory-app
```

This will automatically update:
- Go module paths
- Package names
- Docker image names
- Helm chart names
- CI/CD pipeline references
- Documentation

## Manual Renaming Process

If you prefer to rename manually, here are the key files to update:

### 1. Backend (Go)
- `backend/go.mod` - Update module name
- `backend/main.go` - Update imports and API title
- All Go files in `backend/` - Update import paths
- `backend/services/jwt_service.go` - Update JWT issuer

### 2. Frontend (React)
- `frontend/package.json` - Update package name
- `frontend/Dockerfile` - Update service references

### 3. Kubernetes/Helm
- `charts/hackaton-demo/Chart.yaml` - Update chart name
- `charts/hackaton-demo/values.yaml` - Update image names and domain
- Rename `charts/hackaton-demo/` directory

### 4. CI/CD Pipeline
- `.github/workflows/ci.yaml` - Update image names and paths
- Update GitHub Container Registry repository names

### 5. Documentation
- `README.md` - Update all references
- Update domain references in documentation

## Template Variables

The template uses these naming conventions:

| Variable | Pattern | Example |
|----------|---------|---------|
| Base App Name | `{app-name}` | `inventory-app` |
| Frontend Name | `{app-name}-frontend` | `inventory-app-frontend` |
| Backend Name | `{app-name}-backend` | `inventory-app-backend` |
| Docker Images | `ghcr.io/{username}/{app-name}-{frontend/backend}` | `ghcr.io/frallan97/inventory-app-frontend` |
| Domain | `{app-name}.web.{domain}` | `inventory-app.web.franssjostrom.com` |

## Post-Rename Checklist

After running the rename script, verify these items:

### 1. GitHub Setup
- [ ] Repository name matches your new app name
- [ ] GitHub Container Registry repositories exist:
  - `{app-name}-frontend`
  - `{app-name}-backend`

### 2. Domain Configuration
- [ ] Update DNS settings for your new domain
- [ ] Update SSL certificates if using custom domains

### 3. Environment Variables
- [ ] Update `.env` file with new app-specific values
- [ ] Update Kubernetes secrets if needed

### 4. OAuth Configuration
- [ ] Update Google OAuth client configuration
- [ ] Update redirect URIs for new domain

### 5. Testing
- [ ] Run `go mod tidy` in backend directory
- [ ] Run `npm install` in frontend directory
- [ ] Test local development setup
- [ ] Test CI/CD pipeline

## Example: Creating an Inventory App

```bash
# 1. Clone the template
git clone https://github.com/frallan97/hackaton-demo.git inventory-app
cd inventory-app

# 2. Rename the template
./rename-template.sh inventory-app

# 3. Update GitHub repository
git remote set-url origin https://github.com/frallan97/inventory-app.git

# 4. Create new GitHub repositories for container images
# - inventory-app-frontend
# - inventory-app-backend

# 5. Update environment variables
cp env.example .env
# Edit .env with your specific values

# 6. Test the setup
cd backend && go mod tidy
cd ../frontend && npm install
cd .. && docker-compose up -d
```

## Troubleshooting

### Common Issues

1. **Import errors in Go files**
   - Run `go mod tidy` in the backend directory
   - Verify all import paths are updated

2. **Package name errors in frontend**
   - Run `npm install` in the frontend directory
   - Check `package.json` for correct name

3. **Docker build failures**
   - Verify Dockerfile references are updated
   - Check image names in docker-compose.yml

4. **CI/CD pipeline failures**
   - Verify GitHub Container Registry repositories exist
   - Check workflow file paths and image names

### Getting Help

If you encounter issues:
1. Check the logs in the specific component (backend/frontend)
2. Verify all file paths and names are correctly updated
3. Ensure GitHub repositories and permissions are set up correctly
4. Review the CI/CD workflow configuration

## Contributing to the Template

To improve this template:
1. Make changes to the base template
2. Update `template-config.yaml` if adding new variables
3. Test the rename script with your changes
4. Update this documentation
5. Create a pull request

## License

This template is provided as-is for educational and development purposes. 