package stripe

import (
	"encoding/json"
	"net/http"

	"github.com/frallan97/hackaton-demo-backend/models"
	stripeService "github.com/frallan97/hackaton-demo-backend/services/stripe"
	"github.com/frallan97/hackaton-demo-backend/utils"
)

// StripeController handles all Stripe-related HTTP requests
type StripeController struct {
	stripeManager *stripeService.StripeManager
}

// NewStripeController creates a new Stripe controller
func NewStripeController(stripeManager *stripeService.StripeManager) *StripeController {
	return &StripeController{
		stripeManager: stripeManager,
	}
}

// GetPlansHandler returns available payment plans
func (sc *StripeController) GetPlansHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		plans := sc.stripeManager.Plan.GetAvailablePlans()
		utils.WriteOK(w, plans, "Plans retrieved successfully")
	}
}

// CreateCheckoutSessionHandler creates a new checkout session
func (sc *StripeController) CreateCheckoutSessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		var request models.CreateCheckoutSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.WriteBadRequest(w, "Invalid request data", err)
			return
		}

		// Validate plan exists
		if err := sc.stripeManager.Plan.ValidatePlan(request.PlanID); err != nil {
			utils.WriteBadRequest(w, "Invalid plan", err)
			return
		}

		// Create checkout session
		session, err := sc.stripeManager.Payment.CreateCheckoutSession(
			userID,
			request.PlanID,
			request.SuccessURL,
			request.CancelURL,
		)
		if err != nil {
			utils.WriteInternalServerError(w, "Failed to create checkout session", err)
			return
		}

		utils.WriteOK(w, session, "Checkout session created successfully")
	}
}

// GetPaymentHistoryHandler returns user's payment history
func (sc *StripeController) GetPaymentHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("userID").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		payments, err := sc.stripeManager.Payment.GetPaymentHistory(userID)
		if err != nil {
			utils.WriteInternalServerError(w, "Failed to get payment history", err)
			return
		}

		response := map[string]interface{}{
			"data": payments,
		}
		utils.WriteOK(w, response, "Payment history retrieved successfully")
	}
}

// GetPaymentMetricsHandler returns payment metrics (admin only)
func (sc *StripeController) GetPaymentMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		metrics, err := sc.stripeManager.Payment.GetPaymentMetrics()
		if err != nil {
			utils.WriteInternalServerError(w, "Failed to get payment metrics", err)
			return
		}

		utils.WriteOK(w, metrics, "Payment metrics retrieved successfully")
	}
}

// HealthCheckHandler checks Stripe service health
func (sc *StripeController) HealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		err := sc.stripeManager.HealthCheck()
		if err != nil {
			utils.WriteError(w, http.StatusServiceUnavailable, "Stripe service unavailable", err)
			return
		}

		response := map[string]interface{}{
			"status":  "healthy",
			"service": "stripe",
		}
		utils.WriteOK(w, response, "Stripe service is healthy")
	}
}
