package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/services"
)

// SetupController handles initial setup operations
type SetupController struct {
	adminService *services.AdminService
	userService  *services.UserService
	roleService  *services.RoleService
	jwtService   *services.JWTService
	config       *config.Config
}

// NewSetupController creates a new setup controller
func NewSetupController(dbManager *database.DBManager, jwtService *services.JWTService, config *config.Config) *SetupController {
	return &SetupController{
		adminService: services.NewAdminService(dbManager.DB),
		userService:  services.NewUserService(dbManager.DB),
		roleService:  services.NewRoleService(dbManager.DB),
		jwtService:   jwtService,
		config:       config,
	}
}

// MakeFirstUserAdminHandler assigns admin role to the first user in the system
// This is a convenience endpoint for initial setup
// @Summary Make first user admin
// @Description Assigns admin role to the first user if no admin exists yet
// @Tags setup
// @Produce json
// @Success 200 {string} string "Admin role assigned successfully"
// @Router /api/setup/first-admin [post]
func (sc *SetupController) MakeFirstUserAdminHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if there are any admins already
		users, err := sc.adminService.GetAllUsersWithRolesAndOrganizations()
		if err != nil {
			http.Error(w, "Failed to check existing users", http.StatusInternalServerError)
			return
		}

		// Check if any user already has admin role
		hasAdmin := false
		for _, user := range users {
			for _, role := range user.Roles {
				if role.Name == "admin" {
					hasAdmin = true
					break
				}
			}
			if hasAdmin {
				break
			}
		}

		if hasAdmin {
			http.Error(w, "Admin user already exists", http.StatusConflict)
			return
		}

		if len(users) == 0 {
			http.Error(w, "No users found in system", http.StatusNotFound)
			return
		}

		// Get the first user
		firstUser := users[0]

		// Get admin role
		adminRole, err := sc.roleService.GetRoleByName("admin")
		if err != nil {
			http.Error(w, "Admin role not found", http.StatusInternalServerError)
			return
		}

		// Assign admin role to first user (self-assigned)
		err = sc.adminService.AssignRoleToUser(firstUser.ID, adminRole.ID, firstUser.ID)
		if err != nil {
			http.Error(w, "Failed to assign admin role", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":   "Admin role assigned successfully to first user",
			"user_id":   firstUser.ID,
			"user_name": firstUser.Name,
		})
	}
}

// GenerateDevTokenHandler creates a long-lived development token for Cursor AI
// This should only be used in development environments
func (sc *SetupController) GenerateDevTokenHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Check if we're in development mode
		if sc.config.Environment != "development" {
			http.Error(w, "Development token endpoint is only available in development mode", http.StatusForbidden)
			return
		}

		// Get the first user (or create a dev user)
		users, err := sc.adminService.GetAllUsersWithRolesAndOrganizations()
		if err != nil {
			http.Error(w, "Failed to get users", http.StatusInternalServerError)
			return
		}

		if len(users) == 0 {
			http.Error(w, "No users found. Please create a user first by logging in.", http.StatusNotFound)
			return
		}

		// Use the first user
		user := users[0]

		// Generate the token using JWT service (which has the correct secret key)
		tokenString, _, err := sc.jwtService.GenerateTokens(&user.User)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":      "Development token generated successfully",
			"token":        tokenString,
			"user_id":      user.ID,
			"user_email":   user.Email,
			"expires_days": 30,
			"usage":        "Add this as 'Authorization: Bearer <token>' header in your requests",
		})
	}
}
