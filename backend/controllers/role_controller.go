package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/frallan97/react-go-app-backend/database"
	"github.com/frallan97/react-go-app-backend/models"
	"github.com/frallan97/react-go-app-backend/services"
)

// RoleController handles role-related HTTP requests
type RoleController struct {
	roleService *services.RoleService
}

// NewRoleController creates a new role controller
func NewRoleController(dbManager *database.DBManager) *RoleController {
	return &RoleController{
		roleService: services.NewRoleService(dbManager.DB),
	}
}

// RolesHandler handles role CRUD operations
// @Summary Role operations
// @Description Handle role CRUD operations
// @Tags roles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/roles [get,post,put,delete]
func (rc *RoleController) RolesHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			rc.handleGetRoles(w, r)
		case http.MethodPost:
			rc.handleCreateRole(w, r)
		case http.MethodPut:
			rc.handleUpdateRole(w, r)
		case http.MethodDelete:
			rc.handleDeleteRole(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (rc *RoleController) handleGetRoles(w http.ResponseWriter, r *http.Request) {
	// Check if specific role ID is requested
	roleIDStr := r.URL.Query().Get("id")
	if roleIDStr != "" {
		roleID, err := strconv.Atoi(roleIDStr)
		if err != nil {
			http.Error(w, "Invalid role ID", http.StatusBadRequest)
			return
		}

		role, err := rc.roleService.GetRoleByID(roleID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Role not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(role)
		return
	}

	// Get all roles
	roles, err := rc.roleService.GetAllRoles()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func (rc *RoleController) handleCreateRole(w http.ResponseWriter, r *http.Request) {
	var roleCreate models.RoleCreate
	if err := json.NewDecoder(r.Body).Decode(&roleCreate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if roleCreate.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	role, err := rc.roleService.CreateRole(roleCreate)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Role name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func (rc *RoleController) handleUpdateRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.URL.Query().Get("id")
	if roleIDStr == "" {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	var roleUpdate models.RoleUpdate
	if err := json.NewDecoder(r.Body).Decode(&roleUpdate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if roleUpdate.Name == "" {
		http.Error(w, "Role name is required", http.StatusBadRequest)
		return
	}

	role, err := rc.roleService.UpdateRole(roleID, roleUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Role not found", http.StatusNotFound)
		} else if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Role name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func (rc *RoleController) handleDeleteRole(w http.ResponseWriter, r *http.Request) {
	roleIDStr := r.URL.Query().Get("id")
	if roleIDStr == "" {
		http.Error(w, "Role ID is required", http.StatusBadRequest)
		return
	}

	roleID, err := strconv.Atoi(roleIDStr)
	if err != nil {
		http.Error(w, "Invalid role ID", http.StatusBadRequest)
		return
	}

	err = rc.roleService.DeleteRole(roleID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Role not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}