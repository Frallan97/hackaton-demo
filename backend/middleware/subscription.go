package middleware

import (
	"context"
	"net/http"

	"github.com/frallan97/hackaton-demo-backend/services"
)

// SubscriptionMiddleware provides middleware for subscription-based access control
type SubscriptionMiddleware struct {
	subscriptionService *services.SubscriptionService
}

// NewSubscriptionMiddleware creates a new subscription middleware
func NewSubscriptionMiddleware(subscriptionService *services.SubscriptionService) *SubscriptionMiddleware {
	return &SubscriptionMiddleware{
		subscriptionService: subscriptionService,
	}
}

// RequireSubscription creates middleware that requires an active subscription
func (m *SubscriptionMiddleware) RequireSubscription() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("user_id").(int)
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if user has active subscription
			isSubscribed, err := m.subscriptionService.IsUserSubscribed(userID)
			if err != nil {
				http.Error(w, "Failed to check subscription status", http.StatusInternalServerError)
				return
			}

			if !isSubscribed {
				http.Error(w, "Active subscription required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequirePlan creates middleware that requires a specific subscription plan
func (m *SubscriptionMiddleware) RequirePlan(requiredPlan string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("user_id").(int)
			if !ok {
				http.Error(w, "User not authenticated", http.StatusUnauthorized)
				return
			}

			// Check if user has access to required plan
			hasAccess, err := m.subscriptionService.HasUserAccess(userID, requiredPlan)
			if err != nil {
				http.Error(w, "Failed to check subscription access", http.StatusInternalServerError)
				return
			}

			if !hasAccess {
				http.Error(w, "Higher subscription plan required", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// AddSubscriptionContext adds subscription information to the request context
func (m *SubscriptionMiddleware) AddSubscriptionContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get user ID from context (set by auth middleware)
			userID, ok := r.Context().Value("user_id").(int)
			if !ok {
				next.ServeHTTP(w, r)
				return
			}

			// Get subscription status
			subscription, err := m.subscriptionService.GetUserSubscriptionStatus(userID)
			if err != nil {
				// Log error but don't fail the request
				next.ServeHTTP(w, r)
				return
			}

			// Add subscription info to context
			ctx := context.WithValue(r.Context(), "subscription", subscription)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
