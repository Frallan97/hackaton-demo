package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frallan97/hackaton-demo-backend/controllers"
	"github.com/frallan97/hackaton-demo-backend/database"
	"github.com/frallan97/hackaton-demo-backend/events"
	"github.com/frallan97/hackaton-demo-backend/middleware"
	"github.com/frallan97/hackaton-demo-backend/services"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router handles all routing for the application
type Router struct {
	loginRateLimiter       *middleware.RateLimiter
	healthController       *controllers.HealthController
	messageController      *controllers.MessageController
	authController         *controllers.AuthController
	roleController         *controllers.RoleController
	organizationController *controllers.OrganizationController
	adminController        *controllers.AdminController
	setupController        *controllers.SetupController
	rbacMiddleware         *middleware.RBACMiddleware
	eventService           *events.EventService
}

// NewRouter creates a new router with all controllers
func NewRouter(dbManager *database.DBManager, userService *services.UserService, jwtService *services.JWTService, googleOAuthService *services.GoogleOAuthService, eventService *events.EventService) *Router {
	// Create rate limiter for login endpoint: 5 requests per minute
	loginRateLimiter := middleware.NewRateLimiter(5, time.Minute)
	adminService := services.NewAdminService(dbManager.DB)
	rbacMiddleware := middleware.NewRBACMiddleware(jwtService, adminService)

	return &Router{
		loginRateLimiter:       loginRateLimiter,
		healthController:       controllers.NewHealthController(dbManager),
		messageController:      controllers.NewMessageController(dbManager),
		authController:         controllers.NewAuthController(dbManager, userService, jwtService, googleOAuthService, eventService),
		roleController:         controllers.NewRoleController(dbManager),
		organizationController: controllers.NewOrganizationController(dbManager),
		adminController:        controllers.NewAdminController(dbManager),
		setupController:        controllers.NewSetupController(dbManager),
		rbacMiddleware:         rbacMiddleware,
		eventService:           eventService,
	}
}

// SetupRoutes configures all routes for the application
func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", r.healthController.HealthHandler())

	// Event monitoring endpoint (admin only)
	mux.Handle("/api/events/stats", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.getEventStats)))

	// API endpoints
	mux.HandleFunc("/api/messages", r.messageController.MessagesHandler())

	// Authentication endpoints with rate limiting on login
	loginHandler := middleware.RateLimitMiddleware(r.loginRateLimiter)(http.HandlerFunc(r.authController.GoogleLoginHandler()))
	mux.Handle("/api/auth/google/login", loginHandler)
	mux.HandleFunc("/api/auth/google/url", r.authController.GetAuthURLHandler())
	mux.HandleFunc("/api/auth/refresh", r.authController.RefreshTokenHandler())
	mux.HandleFunc("/api/auth/me", r.authController.GetMeHandler())
	mux.HandleFunc("/api/auth/logout", r.authController.LogoutHandler())

	// Setup endpoints - for initial admin setup
	mux.HandleFunc("/api/setup/first-admin", r.setupController.MakeFirstUserAdminHandler())

	// RBAC endpoints - require authentication
	mux.Handle("/api/roles", r.rbacMiddleware.RequireAnyRole([]string{"admin", "manager"})(http.HandlerFunc(r.roleController.RolesHandler())))
	mux.Handle("/api/organizations", r.rbacMiddleware.RequireAnyRole([]string{"admin", "manager"})(http.HandlerFunc(r.organizationController.OrganizationsHandler())))

	// Admin endpoints - require admin role
	mux.Handle("/api/admin/users", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.GetAllUsersHandler())))
	mux.Handle("/api/admin/assign-role", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.AssignRoleHandler())))
	mux.Handle("/api/admin/remove-role", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.RemoveRoleHandler())))
	mux.Handle("/api/admin/assign-organization", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.AssignOrganizationHandler())))
	mux.Handle("/api/admin/remove-organization", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.RemoveOrganizationHandler())))
	mux.Handle("/api/admin/user-roles", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.GetUserRolesHandler())))
	mux.Handle("/api/admin/user-organizations", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.adminController.GetUserOrganizationsHandler())))

	// Swagger documentation
	mux.Handle("/docs/", httpSwagger.WrapHandler)

	// Apply middleware - CORS must be first to handle preflight requests
	handler := middleware.CORSMiddleware(mux)
	handler = middleware.LoggingMiddleware(handler)

	return handler
}

// getEventStats returns event bus statistics
func (r *Router) getEventStats(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		http.Error(w, "method not allowed", 405)
		return
	}

	stats := r.eventService.GetEventStats()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
