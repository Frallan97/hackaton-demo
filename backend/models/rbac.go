package models

import (
	"time"
)

// Role represents a role in the system
type Role struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// RoleCreate represents the data needed to create a new role
type RoleCreate struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}

// RoleUpdate represents the data needed to update a role
type RoleUpdate struct {
	Name        string `json:"name" validate:"required,max=50"`
	Description string `json:"description"`
}

// Organization represents an organization in the system
type Organization struct {
	ID          int                    `json:"id" db:"id"`
	Name        string                 `json:"name" db:"name"`
	Description string                 `json:"description" db:"description"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}

// OrganizationCreate represents the data needed to create a new organization
type OrganizationCreate struct {
	Name        string                 `json:"name" validate:"required,max=255"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// OrganizationUpdate represents the data needed to update an organization
type OrganizationUpdate struct {
	Name        string                 `json:"name" validate:"required,max=255"`
	Description string                 `json:"description"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID     int       `json:"user_id" db:"user_id"`
	RoleID     int       `json:"role_id" db:"role_id"`
	AssignedAt time.Time `json:"assigned_at" db:"assigned_at"`
	AssignedBy *int      `json:"assigned_by" db:"assigned_by"`
}

// UserOrganization represents the many-to-many relationship between users and organizations
type UserOrganization struct {
	UserID         int       `json:"user_id" db:"user_id"`
	OrganizationID int       `json:"organization_id" db:"organization_id"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
	Role           string    `json:"role" db:"role"`
}

// UserWithRoles represents a user with their assigned roles
type UserWithRoles struct {
	User
	Roles []Role `json:"roles"`
}

// UserWithOrganizations represents a user with their organization memberships
type UserWithOrganizations struct {
	User
	Organizations []Organization `json:"organizations"`
}

// UserWithRolesAndOrganizations represents a user with both roles and organizations
type UserWithRolesAndOrganizations struct {
	User
	Roles         []Role         `json:"roles"`
	Organizations []Organization `json:"organizations"`
}

// RoleAssignmentRequest represents a request to assign a role to a user
type RoleAssignmentRequest struct {
	UserID int `json:"user_id" validate:"required"`
	RoleID int `json:"role_id" validate:"required"`
}

// OrganizationMembershipRequest represents a request to add a user to an organization
type OrganizationMembershipRequest struct {
	UserID         int    `json:"user_id" validate:"required"`
	OrganizationID int    `json:"organization_id" validate:"required"`
	Role           string `json:"role" validate:"required"`
}