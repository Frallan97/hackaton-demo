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

// OrganizationController handles organization-related HTTP requests
type OrganizationController struct {
	orgService *services.OrganizationService
}

// NewOrganizationController creates a new organization controller
func NewOrganizationController(dbManager *database.DBManager) *OrganizationController {
	return &OrganizationController{
		orgService: services.NewOrganizationService(dbManager.DB),
	}
}

// OrganizationsHandler handles organization CRUD operations
// @Summary Organization operations
// @Description Handle organization CRUD operations
// @Tags organizations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Router /api/organizations [get,post,put,delete]
func (oc *OrganizationController) OrganizationsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			oc.handleGetOrganizations(w, r)
		case http.MethodPost:
			oc.handleCreateOrganization(w, r)
		case http.MethodPut:
			oc.handleUpdateOrganization(w, r)
		case http.MethodDelete:
			oc.handleDeleteOrganization(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (oc *OrganizationController) handleGetOrganizations(w http.ResponseWriter, r *http.Request) {
	// Check if specific organization ID is requested
	orgIDStr := r.URL.Query().Get("id")
	if orgIDStr != "" {
		orgID, err := strconv.Atoi(orgIDStr)
		if err != nil {
			http.Error(w, "Invalid organization ID", http.StatusBadRequest)
			return
		}

		org, err := oc.orgService.GetOrganizationByID(orgID)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				http.Error(w, "Organization not found", http.StatusNotFound)
			} else {
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(org)
		return
	}

	// Get all organizations
	orgs, err := oc.orgService.GetAllOrganizations()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orgs)
}

func (oc *OrganizationController) handleCreateOrganization(w http.ResponseWriter, r *http.Request) {
	var orgCreate models.OrganizationCreate
	if err := json.NewDecoder(r.Body).Decode(&orgCreate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if orgCreate.Name == "" {
		http.Error(w, "Organization name is required", http.StatusBadRequest)
		return
	}

	// Initialize metadata if nil
	if orgCreate.Metadata == nil {
		orgCreate.Metadata = make(map[string]interface{})
	}

	org, err := oc.orgService.CreateOrganization(orgCreate)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Organization name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(org)
}

func (oc *OrganizationController) handleUpdateOrganization(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("id")
	if orgIDStr == "" {
		http.Error(w, "Organization ID is required", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	var orgUpdate models.OrganizationUpdate
	if err := json.NewDecoder(r.Body).Decode(&orgUpdate); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if orgUpdate.Name == "" {
		http.Error(w, "Organization name is required", http.StatusBadRequest)
		return
	}

	// Initialize metadata if nil
	if orgUpdate.Metadata == nil {
		orgUpdate.Metadata = make(map[string]interface{})
	}

	org, err := oc.orgService.UpdateOrganization(orgID, orgUpdate)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Organization not found", http.StatusNotFound)
		} else if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "unique") {
			http.Error(w, "Organization name already exists", http.StatusConflict)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

func (oc *OrganizationController) handleDeleteOrganization(w http.ResponseWriter, r *http.Request) {
	orgIDStr := r.URL.Query().Get("id")
	if orgIDStr == "" {
		http.Error(w, "Organization ID is required", http.StatusBadRequest)
		return
	}

	orgID, err := strconv.Atoi(orgIDStr)
	if err != nil {
		http.Error(w, "Invalid organization ID", http.StatusBadRequest)
		return
	}

	err = oc.orgService.DeleteOrganization(orgID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Organization not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}