package controllers

import (
	"net/http"

	"github.com/frallan97/hackaton-demo-backend/database"
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

// HealthHandler responds with {"status":"ok"}
// @Summary     Health check
// @Description Returns 200 if DB is reachable
// @Tags        health
// @Produce     json
// @Success     200  {object}  map[string]string
// @Router      /health [get]
func (hc *HealthController) HealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !hc.dbManager.IsConnected() {
			http.Error(w, `{"status":"db unavailable"}`, http.StatusServiceUnavailable)
			return
		}
		if err := hc.dbManager.DB.Ping(); err != nil {
			http.Error(w, `{"status":"error"}`, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	}
}
