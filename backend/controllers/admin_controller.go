package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/frallan97/react-go-app-backend/database"
	"github.com/frallan97/react-go-app-backend/middleware"
	"github.com/frallan97/react-go-app-backend/models"
	"github.com/frallan97/react-go-app-backend/services"
)

// AdminController handles admin-related HTTP requests
type AdminController struct {
	adminService *services.AdminService
	roleService  *services.RoleService
	orgService   *services.OrganizationService
}

// NewAdminController creates a new admin controller
func NewAdminController(dbManager *database.DBManager) *AdminController {
	return &AdminController{
		adminService: services.NewAdminService(dbManager.DB),
		roleService:  services.NewRoleService(dbManager.DB),
		orgService:   services.NewOrganizationService(dbManager.DB),
	}
}

// GetAllUsersHandler returns all users with their roles and organizations
// @Summary Get all users with roles and organizations
// @Description Get all users with their assigned roles and organization memberships (Admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.UserWithRolesAndOrganizations
// @Router /api/admin/users [get]
func (ac *AdminController) GetAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		users, err := ac.adminService.GetAllUsersWithRolesAndOrganizations()
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

// AssignRoleHandler assigns a role to a user
// @Summary Assign role to user
// @Description Assign a role to a user (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RoleAssignmentRequest true "Role assignment request"
// @Success 200 {string} string "Role assigned successfully"
// @Router /api/admin/assign-role [post]
func (ac *AdminController) AssignRoleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RoleAssignmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 || req.RoleID == 0 {
			http.Error(w, "User ID and Role ID are required", http.StatusBadRequest)
			return
		}

		// Get the admin user ID from context
		adminUserID, ok := middleware.GetUserIDFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		err := ac.adminService.AssignRoleToUser(req.UserID, req.RoleID, adminUserID)
		if err != nil {
			if err.Error() == "user already has this role" {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Role assigned successfully"})
	}
}

// RemoveRoleHandler removes a role from a user
// @Summary Remove role from user
// @Description Remove a role from a user (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.RoleAssignmentRequest true "Role assignment request"
// @Success 200 {string} string "Role removed successfully"
// @Router /api/admin/remove-role [post]
func (ac *AdminController) RemoveRoleHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.RoleAssignmentRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 || req.RoleID == 0 {
			http.Error(w, "User ID and Role ID are required", http.StatusBadRequest)
			return
		}

		err := ac.adminService.RemoveRoleFromUser(req.UserID, req.RoleID)
		if err != nil {
			if err.Error() == "user does not have this role" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Role removed successfully"})
	}
}

// AssignOrganizationHandler adds a user to an organization
// @Summary Add user to organization
// @Description Add a user to an organization with a specific role (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.OrganizationMembershipRequest true "Organization membership request"
// @Success 200 {string} string "User added to organization successfully"
// @Router /api/admin/assign-organization [post]
func (ac *AdminController) AssignOrganizationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.OrganizationMembershipRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 || req.OrganizationID == 0 || req.Role == "" {
			http.Error(w, "User ID, Organization ID, and Role are required", http.StatusBadRequest)
			return
		}

		err := ac.adminService.AddUserToOrganization(req.UserID, req.OrganizationID, req.Role)
		if err != nil {
			if err.Error() == "user is already a member of this organization" {
				http.Error(w, err.Error(), http.StatusConflict)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User added to organization successfully"})
	}
}

// RemoveOrganizationHandler removes a user from an organization
// @Summary Remove user from organization
// @Description Remove a user from an organization (Admin only)
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.OrganizationMembershipRequest true "Organization membership request"
// @Success 200 {string} string "User removed from organization successfully"
// @Router /api/admin/remove-organization [post]
func (ac *AdminController) RemoveOrganizationHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req models.OrganizationMembershipRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if req.UserID == 0 || req.OrganizationID == 0 {
			http.Error(w, "User ID and Organization ID are required", http.StatusBadRequest)
			return
		}

		err := ac.adminService.RemoveUserFromOrganization(req.UserID, req.OrganizationID)
		if err != nil {
			if err.Error() == "user is not a member of this organization" {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User removed from organization successfully"})
	}
}

// GetUserRolesHandler gets roles for a specific user
// @Summary Get user roles
// @Description Get all roles assigned to a specific user (Admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id query int true "User ID"
// @Success 200 {array} models.Role
// @Router /api/admin/user-roles [get]
func (ac *AdminController) GetUserRolesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.URL.Query().Get("id")
		if userIDStr == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		roles, err := ac.adminService.GetUserRoles(userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(roles)
	}
}

// GetUserOrganizationsHandler gets organizations for a specific user
// @Summary Get user organizations
// @Description Get all organizations a specific user belongs to (Admin only)
// @Tags admin
// @Produce json
// @Security BearerAuth
// @Param id query int true "User ID"
// @Success 200 {array} models.Organization
// @Router /api/admin/user-organizations [get]
func (ac *AdminController) GetUserOrganizationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		userIDStr := r.URL.Query().Get("id")
		if userIDStr == "" {
			http.Error(w, "User ID is required", http.StatusBadRequest)
			return
		}

		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		orgs, err := ac.adminService.GetUserOrganizations(userID)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(orgs)
	}
}