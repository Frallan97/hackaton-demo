package controllers

import (
	"net/http"
	"time"

	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/utils"
)

// HealthController handles health-related endpoints
type HealthController struct {
	dbManager *database.DBManager
}

// NewHealthController creates a new health controller
func NewHealthController(dbManager *database.DBManager) *HealthController {
	return &HealthController{
		dbManager: dbManager,
	}
}

// HealthResponse represents the health check response data
type HealthResponse struct {
	Status    string    `json:"status"`
	Database  string    `json:"database"`
	Timestamp time.Time `json:"timestamp"`
	Uptime    string    `json:"uptime,omitempty"`
}

// HealthHandler responds with standardized health status
// @Summary     Health check
// @Description Returns 200 if DB is reachable
// @Tags        health
// @Produce     json
// @Success     200  {object}  utils.APIResponse{data=HealthResponse}
// @Failure     503  {object}  utils.APIResponse
// @Failure     500  {object}  utils.APIResponse
// @Router      /health [get]
func (hc *HealthController) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check database connection
		if !hc.dbManager.IsConnected() {
			utils.WriteError(w, http.StatusServiceUnavailable, "Database connection unavailable", nil)
			return
		}

		// Test database ping
		if err := hc.dbManager.DB.Ping(); err != nil {
			utils.WriteInternalServerError(w, "Database ping failed", err)
			return
		}

		// Create health response
		healthData := &HealthResponse{
			Status:    "healthy",
			Database:  "connected",
			Timestamp: time.Now(),
		}

		utils.WriteOK(w, healthData, "Service is healthy")
	}
}
