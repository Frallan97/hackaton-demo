package services

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/frallan97/react-go-app-backend/models"
)

// UserService handles user-related database operations
type UserService struct {
	db *sql.DB
}

// NewUserService creates a new user service
func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a new user in the database
func (u *UserService) CreateUser(userData *models.UserCreate) (*models.User, error) {
	query := `
		INSERT INTO users (email, name, picture, google_id, is_active, last_login_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, email, name, picture, google_id, is_active, last_login_at, created_at, updated_at
	`

	now := time.Now()
	user := &models.User{}

	err := u.db.QueryRow(
		query,
		userData.Email,
		userData.Name,
		userData.Picture,
		userData.GoogleID,
		true, // is_active
		now,  // last_login_at
		now,  // created_at
		now,  // updated_at
	).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

// GetUserByGoogleID retrieves a user by their Google ID
func (u *UserService) GetUserByGoogleID(googleID string) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE google_id = $1 AND is_active = true
	`

	user := &models.User{}
	err := u.db.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by Google ID: %w", err)
	}

	return user, nil
}

// GetUserByID retrieves a user by their ID
func (u *UserService) GetUserByID(userID int) (*models.User, error) {
	query := `
		SELECT id, email, name, picture, google_id, is_active, last_login_at, created_at, updated_at
		FROM users
		WHERE id = $1 AND is_active = true
	`

	user := &models.User{}
	err := u.db.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return user, nil
}

// UpdateUserLastLogin updates the last login timestamp for a user
func (u *UserService) UpdateUserLastLogin(userID int) error {
	query := `
		UPDATE users
		SET last_login_at = $1, updated_at = $2
		WHERE id = $3
	`

	now := time.Now()
	_, err := u.db.Exec(query, now, now, userID)
	if err != nil {
		return fmt.Errorf("failed to update user last login: %w", err)
	}

	return nil
}

// UpdateUserProfile updates a user's profile information
func (u *UserService) UpdateUserProfile(userID int, name, picture string) (*models.User, error) {
	query := `
		UPDATE users
		SET name = $1, picture = $2, updated_at = $3
		WHERE id = $4
		RETURNING id, email, name, picture, google_id, is_active, last_login_at, created_at, updated_at
	`

	now := time.Now()
	user := &models.User{}

	err := u.db.QueryRow(query, name, picture, now, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.IsActive,
		&user.LastLoginAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	return user, nil
}
