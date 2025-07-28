package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/frallan97/hackaton-demo-backend/services"
)

// RBACMiddleware provides role-based access control
type RBACMiddleware struct {
	jwtService   *services.JWTService
	adminService *services.AdminService
}

// NewRBACMiddleware creates a new RBAC middleware
func NewRBACMiddleware(jwtService *services.JWTService, adminService *services.AdminService) *RBACMiddleware {
	return &RBACMiddleware{
		jwtService:   jwtService,
		adminService: adminService,
	}
}

// RequireRole returns a middleware that requires a specific role
func (rbac *RBACMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := rbac.getUserIDFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hasRole, err := rbac.adminService.UserHasRole(userID, role)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !hasRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID to context for use in handlers
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAnyRole returns a middleware that requires any of the specified roles
func (rbac *RBACMiddleware) RequireAnyRole(roles []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := rbac.getUserIDFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			hasAnyRole := false
			for _, role := range roles {
				hasRole, err := rbac.adminService.UserHasRole(userID, role)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				if hasRole {
					hasAnyRole = true
					break
				}
			}

			if !hasAnyRole {
				http.Error(w, "Forbidden: insufficient permissions", http.StatusForbidden)
				return
			}

			// Add user ID to context for use in handlers
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireAuth returns a middleware that requires authentication but no specific role
func (rbac *RBACMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, err := rbac.getUserIDFromRequest(r)
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Add user ID to context for use in handlers
			ctx := context.WithValue(r.Context(), "userID", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// getUserIDFromRequest extracts and validates the user ID from the JWT token in the request
func (rbac *RBACMiddleware) getUserIDFromRequest(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, &AuthError{Message: "authorization header required"}
	}

	// Extract token from "Bearer <token>"
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return 0, &AuthError{Message: "invalid authorization header format"}
	}

	token := parts[1]
	claims, err := rbac.jwtService.ValidateToken(token)
	if err != nil {
		return 0, &AuthError{Message: "invalid token"}
	}

	return claims.UserID, nil
}

// AuthError represents an authentication error
type AuthError struct {
	Message string
}

func (e *AuthError) Error() string {
	return e.Message
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value("userID").(int)
	return userID, ok
}