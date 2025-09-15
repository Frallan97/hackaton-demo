package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
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
	stripeController       *controllers.StripeController
	rbacMiddleware         *middleware.RBACMiddleware
	eventService           *events.EventService
}

// NewRouter creates a new router with all controllers
func NewRouter(dbManager *database.DBManager, userService *services.UserService, jwtService *services.JWTService, googleOAuthService *services.GoogleOAuthService, eventService *events.EventService, config *config.Config) *Router {
	// Create rate limiter for login endpoint: 5 requests per minute
	loginRateLimiter := middleware.NewRateLimiter(5, time.Minute)
	adminService := services.NewAdminService(dbManager.DB)
	roleService := services.NewRoleService(dbManager.DB)
	rbacMiddleware := middleware.NewRBACMiddleware(jwtService, adminService)

	// Initialize Stripe services
	stripeService := services.NewStripeService(dbManager.DB, config)
	subscriptionService := services.NewSubscriptionService(dbManager.DB, stripeService)

	return &Router{
		loginRateLimiter:       loginRateLimiter,
		healthController:       controllers.NewHealthController(dbManager),
		messageController:      controllers.NewMessageController(dbManager),
		authController:         controllers.NewAuthController(dbManager, userService, jwtService, googleOAuthService, eventService, roleService, adminService),
		roleController:         controllers.NewRoleController(dbManager),
		organizationController: controllers.NewOrganizationController(dbManager),
		adminController:        controllers.NewAdminController(dbManager),
		setupController:        controllers.NewSetupController(dbManager, jwtService, config),
		stripeController:       controllers.NewStripeController(stripeService, subscriptionService, config),
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
	mux.HandleFunc("/api/setup/dev-token", r.setupController.GenerateDevTokenHandler())

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

	// Stripe endpoints - public endpoints
	mux.HandleFunc("/api/stripe/webhook", r.stripeController.WebhookHandler())
	mux.HandleFunc("/api/stripe/plans", r.stripeController.GetAvailablePlansHandler())

	// Stripe endpoints - require authentication
	mux.Handle("/api/stripe/checkout", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.CreateCheckoutSessionHandler())))
	mux.Handle("/api/stripe/subscription", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.GetUserSubscriptionHandler())))
	mux.Handle("/api/stripe/subscription/history", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.GetUserSubscriptionHistoryHandler())))
	mux.Handle("/api/stripe/payments", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.GetUserPaymentHistoryHandler())))
	mux.Handle("/api/stripe/subscription/cancel", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.CancelSubscriptionHandler())))
	mux.Handle("/api/stripe/subscription/reactivate", r.rbacMiddleware.RequireAnyRole([]string{"user", "admin", "manager"})(http.HandlerFunc(r.stripeController.ReactivateSubscriptionHandler())))

	// Stripe admin endpoints - require admin role
	mux.Handle("/api/stripe/admin/metrics", r.rbacMiddleware.RequireRole("admin")(http.HandlerFunc(r.stripeController.GetSubscriptionMetricsHandler())))

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
