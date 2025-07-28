package services

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/frallan97/hackaton-demo-backend/models"
)

// AdminService handles admin-related business logic for managing users, roles, and organizations
type AdminService struct {
	db *sql.DB
}

// NewAdminService creates a new admin service
func NewAdminService(db *sql.DB) *AdminService {
	return &AdminService{db: db}
}

// GetAllUsersWithRolesAndOrganizations retrieves all users with their roles and organizations
func (as *AdminService) GetAllUsersWithRolesAndOrganizations() ([]models.UserWithRolesAndOrganizations, error) {
	// First get all users
	users, err := as.getAllUsers()
	if err != nil {
		return nil, err
	}

	var result []models.UserWithRolesAndOrganizations
	for _, user := range users {
		userWithData := models.UserWithRolesAndOrganizations{
			User: user,
		}

		// Get roles for this user
		roles, err := as.getUserRoles(user.ID)
		if err != nil {
			return nil, err
		}
		userWithData.Roles = roles

		// Get organizations for this user
		orgs, err := as.getUserOrganizations(user.ID)
		if err != nil {
			return nil, err
		}
		userWithData.Organizations = orgs

		result = append(result, userWithData)
	}

	return result, nil
}

// AssignRoleToUser assigns a role to a user
func (as *AdminService) AssignRoleToUser(userID, roleID, assignedBy int) error {
	// Check if assignment already exists
	exists, err := as.userHasRole(userID, roleID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user already has this role")
	}

	query := `INSERT INTO user_roles (user_id, role_id, assigned_by) VALUES ($1, $2, $3)`
	_, err = as.db.Exec(query, userID, roleID, assignedBy)
	if err != nil {
		return fmt.Errorf("failed to assign role to user: %w", err)
	}

	return nil
}

// RemoveRoleFromUser removes a role from a user
func (as *AdminService) RemoveRoleFromUser(userID, roleID int) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`
	
	result, err := as.db.Exec(query, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to remove role from user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user does not have this role")
	}

	return nil
}

// AddUserToOrganization adds a user to an organization with a specific role
func (as *AdminService) AddUserToOrganization(userID, organizationID int, role string) error {
	// Check if membership already exists
	exists, err := as.userInOrganization(userID, organizationID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("user is already a member of this organization")
	}

	query := `INSERT INTO user_organizations (user_id, organization_id, role) VALUES ($1, $2, $3)`
	_, err = as.db.Exec(query, userID, organizationID, role)
	if err != nil {
		return fmt.Errorf("failed to add user to organization: %w", err)
	}

	return nil
}

// RemoveUserFromOrganization removes a user from an organization
func (as *AdminService) RemoveUserFromOrganization(userID, organizationID int) error {
	query := `DELETE FROM user_organizations WHERE user_id = $1 AND organization_id = $2`
	
	result, err := as.db.Exec(query, userID, organizationID)
	if err != nil {
		return fmt.Errorf("failed to remove user from organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user is not a member of this organization")
	}

	return nil
}

// UserHasRole checks if a user has a specific role
func (as *AdminService) UserHasRole(userID int, roleName string) (bool, error) {
	query := `
		SELECT COUNT(*) 
		FROM user_roles ur 
		JOIN roles r ON ur.role_id = r.id 
		WHERE ur.user_id = $1 AND r.name = $2
	`
	
	var count int
	err := as.db.QueryRow(query, userID, roleName).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return count > 0, nil
}

// GetUserRoles returns all roles for a specific user
func (as *AdminService) GetUserRoles(userID int) ([]models.Role, error) {
	return as.getUserRoles(userID)
}

// GetUserOrganizations returns all organizations for a specific user
func (as *AdminService) GetUserOrganizations(userID int) ([]models.Organization, error) {
	return as.getUserOrganizations(userID)
}

// Helper methods

func (as *AdminService) getAllUsers() ([]models.User, error) {
	query := `SELECT id, email, name, picture, google_id, is_active, last_login_at, created_at, updated_at FROM users ORDER BY name`
	
	rows, err := as.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Email, &user.Name, &user.Picture, &user.GoogleID, &user.IsActive, &user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (as *AdminService) getUserRoles(userID int) ([]models.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at 
		FROM roles r 
		JOIN user_roles ur ON r.id = ur.role_id 
		WHERE ur.user_id = $1 
		ORDER BY r.name
	`
	
	rows, err := as.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []models.Role
	for rows.Next() {
		var role models.Role
		err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (as *AdminService) getUserOrganizations(userID int) ([]models.Organization, error) {
	query := `
		SELECT o.id, o.name, o.description, o.metadata, o.created_at, o.updated_at 
		FROM organizations o 
		JOIN user_organizations uo ON o.id = uo.organization_id 
		WHERE uo.user_id = $1 
		ORDER BY o.name
	`
	
	rows, err := as.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user organizations: %w", err)
	}
	defer rows.Close()

	var organizations []models.Organization
	for rows.Next() {
		var org models.Organization
		var metadataJSON []byte
		err := rows.Scan(&org.ID, &org.Name, &org.Description, &metadataJSON, &org.CreatedAt, &org.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan organization: %w", err)
		}

		// Parse JSON metadata
		if len(metadataJSON) > 0 {
			err = json.Unmarshal(metadataJSON, &org.Metadata)
			if err != nil {
				org.Metadata = make(map[string]interface{})
			}
		} else {
			org.Metadata = make(map[string]interface{})
		}

		organizations = append(organizations, org)
	}

	return organizations, nil
}

func (as *AdminService) userHasRole(userID, roleID int) (bool, error) {
	query := `SELECT COUNT(*) FROM user_roles WHERE user_id = $1 AND role_id = $2`
	
	var count int
	err := as.db.QueryRow(query, userID, roleID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user role: %w", err)
	}

	return count > 0, nil
}

func (as *AdminService) userInOrganization(userID, organizationID int) (bool, error) {
	query := `SELECT COUNT(*) FROM user_organizations WHERE user_id = $1 AND organization_id = $2`
	
	var count int
	err := as.db.QueryRow(query, userID, organizationID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check user organization: %w", err)
	}

	return count > 0, nil
}