package services

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/frallan97/react-go-app-backend/models"
)

// RoleService handles role-related business logic
type RoleService struct {
	db *sql.DB
}

// NewRoleService creates a new role service
func NewRoleService(db *sql.DB) *RoleService {
	return &RoleService{db: db}
}

// GetAllRoles retrieves all roles from the database
func (rs *RoleService) GetAllRoles() ([]models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name`
	
	rows, err := rs.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %w", err)
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

// GetRoleByID retrieves a role by its ID
func (rs *RoleService) GetRoleByID(id int) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`
	
	var role models.Role
	err := rs.db.QueryRow(query, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to query role: %w", err)
	}

	return &role, nil
}

// GetRoleByName retrieves a role by its name
func (rs *RoleService) GetRoleByName(name string) (*models.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`
	
	var role models.Role
	err := rs.db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to query role: %w", err)
	}

	return &role, nil
}

// CreateRole creates a new role
func (rs *RoleService) CreateRole(roleCreate models.RoleCreate) (*models.Role, error) {
	query := `INSERT INTO roles (name, description) VALUES ($1, $2) RETURNING id, name, description, created_at, updated_at`
	
	var role models.Role
	err := rs.db.QueryRow(query, roleCreate.Name, roleCreate.Description).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create role: %w", err)
	}

	return &role, nil
}

// UpdateRole updates an existing role
func (rs *RoleService) UpdateRole(id int, roleUpdate models.RoleUpdate) (*models.Role, error) {
	query := `UPDATE roles SET name = $1, description = $2 WHERE id = $3 RETURNING id, name, description, created_at, updated_at`
	
	var role models.Role
	err := rs.db.QueryRow(query, roleUpdate.Name, roleUpdate.Description, id).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("role not found")
		}
		return nil, fmt.Errorf("failed to update role: %w", err)
	}

	return &role, nil
}

// DeleteRole deletes a role by its ID
func (rs *RoleService) DeleteRole(id int) error {
	query := `DELETE FROM roles WHERE id = $1`
	
	result, err := rs.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete role: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("role not found")
	}

	return nil
}

// OrganizationService handles organization-related business logic
type OrganizationService struct {
	db *sql.DB
}

// NewOrganizationService creates a new organization service
func NewOrganizationService(db *sql.DB) *OrganizationService {
	return &OrganizationService{db: db}
}

// GetAllOrganizations retrieves all organizations from the database
func (os *OrganizationService) GetAllOrganizations() ([]models.Organization, error) {
	query := `SELECT id, name, description, metadata, created_at, updated_at FROM organizations ORDER BY name`
	
	rows, err := os.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query organizations: %w", err)
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

// GetOrganizationByID retrieves an organization by its ID
func (os *OrganizationService) GetOrganizationByID(id int) (*models.Organization, error) {
	query := `SELECT id, name, description, metadata, created_at, updated_at FROM organizations WHERE id = $1`
	
	var org models.Organization
	var metadataJSON []byte
	err := os.db.QueryRow(query, id).Scan(&org.ID, &org.Name, &org.Description, &metadataJSON, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to query organization: %w", err)
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

	return &org, nil
}

// CreateOrganization creates a new organization
func (os *OrganizationService) CreateOrganization(orgCreate models.OrganizationCreate) (*models.Organization, error) {
	metadataJSON, err := json.Marshal(orgCreate.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `INSERT INTO organizations (name, description, metadata) VALUES ($1, $2, $3) RETURNING id, name, description, metadata, created_at, updated_at`
	
	var org models.Organization
	var returnedMetadataJSON []byte
	err = os.db.QueryRow(query, orgCreate.Name, orgCreate.Description, metadataJSON).Scan(&org.ID, &org.Name, &org.Description, &returnedMetadataJSON, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create organization: %w", err)
	}

	// Parse JSON metadata
	if len(returnedMetadataJSON) > 0 {
		err = json.Unmarshal(returnedMetadataJSON, &org.Metadata)
		if err != nil {
			org.Metadata = make(map[string]interface{})
		}
	} else {
		org.Metadata = make(map[string]interface{})
	}

	return &org, nil
}

// UpdateOrganization updates an existing organization
func (os *OrganizationService) UpdateOrganization(id int, orgUpdate models.OrganizationUpdate) (*models.Organization, error) {
	metadataJSON, err := json.Marshal(orgUpdate.Metadata)
	if err != nil {
		metadataJSON = []byte("{}")
	}

	query := `UPDATE organizations SET name = $1, description = $2, metadata = $3 WHERE id = $4 RETURNING id, name, description, metadata, created_at, updated_at`
	
	var org models.Organization
	var returnedMetadataJSON []byte
	err = os.db.QueryRow(query, orgUpdate.Name, orgUpdate.Description, metadataJSON, id).Scan(&org.ID, &org.Name, &org.Description, &returnedMetadataJSON, &org.CreatedAt, &org.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("organization not found")
		}
		return nil, fmt.Errorf("failed to update organization: %w", err)
	}

	// Parse JSON metadata
	if len(returnedMetadataJSON) > 0 {
		err = json.Unmarshal(returnedMetadataJSON, &org.Metadata)
		if err != nil {
			org.Metadata = make(map[string]interface{})
		}
	} else {
		org.Metadata = make(map[string]interface{})
	}

	return &org, nil
}

// DeleteOrganization deletes an organization by its ID
func (os *OrganizationService) DeleteOrganization(id int) error {
	query := `DELETE FROM organizations WHERE id = $1`
	
	result, err := os.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete organization: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("organization not found")
	}

	return nil
}