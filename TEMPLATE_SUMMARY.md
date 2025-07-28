# Template System Summary

## Overview

This repository is now set up as a **reusable template** for creating new React + Go applications. The template includes:

- ✅ **Automated renaming scripts** for easy customization
- ✅ **Complete CI/CD pipeline** with GitHub Actions
- ✅ **Kubernetes deployment** with Helm charts
- ✅ **OAuth authentication** with Google
- ✅ **Role-based access control (RBAC)**
- ✅ **Docker containerization**
- ✅ **Database migrations**

## Quick Usage

### 1. Basic Rename (Same GitHub username/domain)
```bash
./rename-template.sh inventory-app
```

### 2. Advanced Rename (Custom GitHub username/domain)
```bash
./rename-template-advanced.sh inventory-app myusername mydomain.com
```

## What Gets Updated

The scripts automatically update **all** references across the entire codebase:

### Backend (Go)
- Module paths: `github.com/frallan97/hackaton-demo-backend` → `github.com/myusername/inventory-app-backend`
- Import statements in all `.go` files
- JWT issuer configuration
- API documentation titles

### Frontend (React)
- Package name: `hackaton-demo-frontend` → `inventory-app-frontend`
- Docker service references
- Build configurations

### Kubernetes/Helm
- Chart name: `hackaton-demo` → `inventory-app`
- Image references: `ghcr.io/frallan97/hackaton-demo-frontend` → `ghcr.io/myusername/inventory-app-frontend`
- Domain configuration: `hackaton-demo.web.franssjostrom.com` → `inventory-app.web.mydomain.com`

### CI/CD Pipeline
- GitHub Container Registry image names
- Helm chart paths
- Build and deployment configurations

### Documentation
- README files
- Configuration examples
- Domain references

## Template Files Created

1. **`rename-template.sh`** - Basic renaming script
2. **`rename-template-advanced.sh`** - Advanced script with custom GitHub/domain
3. **`template-config.yaml`** - Configuration file defining all template variables
4. **`TEMPLATE_USAGE.md`** - Comprehensive usage guide
5. **`TEMPLATE_SUMMARY.md`** - This summary document

## Example: Creating an Inventory Management App

```bash
# 1. Clone this template
git clone https://github.com/frallan97/hackaton-demo.git inventory-app
cd inventory-app

# 2. Rename everything
./rename-template-advanced.sh inventory-app myusername mydomain.com

# 3. Update Git remote
git remote set-url origin https://github.com/myusername/inventory-app.git

# 4. Create GitHub Container Registry repositories
# - inventory-app-frontend
# - inventory-app-backend

# 5. Update environment variables
cp env.example .env
# Edit .env with your specific values

# 6. Test locally
docker-compose up -d

# 7. Deploy to Kubernetes
helm install inventory-app charts/inventory-app/
```

## Benefits

### For You (Template Creator)
- ✅ Easy to maintain one template for multiple projects
- ✅ Consistent architecture across all applications
- ✅ Automated setup reduces manual errors

### For Template Users
- ✅ One-command setup for new applications
- ✅ Consistent, battle-tested architecture
- ✅ Full CI/CD pipeline ready to go
- ✅ Production-ready Kubernetes deployment

## Next Steps

1. **Test the template** by running the rename script on a copy
2. **Create GitHub repositories** for your new applications
3. **Set up GitHub Container Registry** repositories
4. **Configure OAuth** for your new domains
5. **Deploy and test** the complete pipeline

## Support

If you encounter issues:
1. Check the `TEMPLATE_USAGE.md` for troubleshooting
2. Verify all GitHub repositories and permissions
3. Review the CI/CD workflow configuration
4. Check domain DNS settings

---

**This template system makes it trivial to spin up new production-ready applications with consistent architecture and full automation!** 🚀 