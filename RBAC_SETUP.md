# RBAC Setup Guide

This guide explains how to set up and use the Role-Based Access Control (RBAC) system in the React-Go app.

## Overview

The RBAC system provides:
- **Roles**: admin, manager, editor, reader (with default roles created automatically)
- **Organizations**: Groups for users with optional membership
- **Admin Dashboard**: Complete UI for managing users, roles, and organizations
- **Protected Routes**: Role-based access control for frontend routes

## Initial Setup

### 1. Database Migration

The RBAC tables are created automatically when the backend starts. The following tables are created:
- `roles` - Contains the role definitions
- `organizations` - Contains organization information  
- `user_roles` - Junction table for user-role assignments
- `user_organizations` - Junction table for user-organization memberships

### 2. Setting Up the First Admin User

To set up the first admin user:

1. **Start the backend and database**:
   ```bash
   # Start database and Redis
   docker-compose up db redis -d
   
   # Start backend
   cd backend
   go run main.go
   ```

2. **Start the frontend**:
   ```bash
   # Start frontend
   cd frontend  
   npm run dev
   ```

3. **Login with Google OAuth**:
   - Go to http://localhost:3000
   - Login with your Google account
   - This creates your user account in the system

4. **Assign admin role to first user**:
   ```bash
   # Make API call to assign admin role to first user
   curl -X POST http://localhost:8080/api/setup/first-admin
   ```
   
   This endpoint:
   - Checks if any admin users already exist
   - If none exist, assigns admin role to the first user in the system
   - Returns success message with user details

5. **Refresh your frontend session**:
   - Click "Refresh Token" or logout and login again
   - You should now see "Admin Dashboard" link in the home page

## Using the Admin Dashboard

### Accessing the Admin Dashboard

1. Login with a user that has the `admin` role
2. Navigate to `/admin` or click "Admin Dashboard" link on home page
3. The dashboard has three tabs: Users, Roles, Organizations

### Users Tab

- **View all users**: See all users with their roles and organizations
- **Assign roles**: Click the "+" button next to roles, select a role to assign
- **Remove roles**: Click the "×" button next to any assigned role
- **Add to organization**: Click the "+" button next to organizations, select an organization
- **Remove from organization**: Click the "×" button next to any organization membership

### Roles Tab

- **View all roles**: See all available roles in the system
- **Create new role**: Click "Create Role" button, fill in name and description
- **Default roles**: admin, manager, editor, reader are created automatically

### Organizations Tab

- **View all organizations**: See all organizations in the system
- **Create new organization**: Click "Create Organization" button, fill in details
- **Metadata support**: Organizations support JSONB metadata for additional information

## API Endpoints

### Role Management (admin/manager access required)
- `GET /api/roles` - List all roles
- `POST /api/roles` - Create new role
- `PUT /api/roles?id={id}` - Update role
- `DELETE /api/roles?id={id}` - Delete role

### Organization Management (admin/manager access required)  
- `GET /api/organizations` - List all organizations
- `POST /api/organizations` - Create new organization
- `PUT /api/organizations?id={id}` - Update organization
- `DELETE /api/organizations?id={id}` - Delete organization

### Admin Management (admin access required)
- `GET /api/admin/users` - List all users with roles and organizations
- `POST /api/admin/assign-role` - Assign role to user
- `POST /api/admin/remove-role` - Remove role from user
- `POST /api/admin/assign-organization` - Add user to organization
- `POST /api/admin/remove-organization` - Remove user from organization
- `GET /api/admin/user-roles?id={user_id}` - Get user's roles
- `GET /api/admin/user-organizations?id={user_id}` - Get user's organizations

### Setup (no auth required)
- `POST /api/setup/first-admin` - Assign admin role to first user (if no admin exists)

## Role Permissions

### Default Roles
- **admin**: Full system access, can manage all users, roles, and organizations
- **manager**: Can view and manage roles and organizations (not implemented in current UI)
- **editor**: Can edit content (not implemented in current UI)  
- **reader**: Read-only access (not implemented in current UI)

### Route Protection
- `/` - Requires authentication (any logged-in user)
- `/admin` - Requires admin role
- `/login` - Public route

## Troubleshooting

### Common Issues

1. **Can't access admin dashboard**:
   - Ensure you have admin role assigned
   - Check that you've refreshed your token after role assignment
   - Verify the user has the admin role in the database

2. **Setup endpoint fails**:
   - Make sure database is running and migrations have been applied
   - Ensure at least one user exists in the system (login first)
   - Check that no admin user already exists

3. **Token refresh issues**:
   - Clear localStorage and login again
   - Check that JWT_SECRET_KEY is consistent between sessions
   - Verify backend is running and accessible

### Database Queries for Debugging

```sql
-- Check all users and their roles
SELECT u.name, u.email, r.name as role_name 
FROM users u 
LEFT JOIN user_roles ur ON u.id = ur.user_id 
LEFT JOIN roles r ON ur.role_id = r.id;

-- Check all roles
SELECT * FROM roles;

-- Check all organizations  
SELECT * FROM organizations;

-- Manually assign admin role to user ID 1
INSERT INTO user_roles (user_id, role_id, assigned_by) 
VALUES (1, (SELECT id FROM roles WHERE name = 'admin'), 1);
```

## Security Considerations

- Admin role provides full system access - assign carefully
- JWT tokens contain user ID for role validation
- All admin endpoints require valid JWT token with admin role
- Database foreign keys ensure referential integrity
- Role assignments are tracked with assigned_by field for audit purposes

## Extending the System

### Adding New Roles
1. Use the admin dashboard to create new roles
2. Add endpoint protection in backend routes using RBAC middleware
3. Add frontend route guards in ProtectedRoute component

### Adding Organization Features
1. Organizations support JSONB metadata for custom fields
2. User-organization relationships include a role field for organization-level permissions
3. Extend the UI to show organization-specific role management

### Custom Permissions
1. Extend the Role model with permissions field
2. Implement permission-based middleware alongside role-based middleware
3. Add fine-grained permission management in admin dashboard