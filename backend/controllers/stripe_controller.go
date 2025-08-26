package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/frallan97/hackaton-demo-backend/services"
	"github.com/frallan97/hackaton-demo-backend/utils"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/webhook"
)

// StripeController handles Stripe-related HTTP requests
type StripeController struct {
	stripeService       *services.StripeService
	subscriptionService *services.SubscriptionService
	config              *config.Config
}

// NewStripeController creates a new Stripe controller
func NewStripeController(stripeService *services.StripeService, subscriptionService *services.SubscriptionService, config *config.Config) *StripeController {
	return &StripeController{
		stripeService:       stripeService,
		subscriptionService: subscriptionService,
		config:              config,
	}
}

// CreateCheckoutSessionHandler handles creating a new checkout session
func (c *StripeController) CreateCheckoutSessionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Get user ID from context (set by auth middleware)
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		var req models.CreateCheckoutSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.WriteBadRequest(w, "Invalid request body", err)
			return
		}

		// Validate request
		if req.PlanID == "" || req.SuccessURL == "" || req.CancelURL == "" {
			utils.WriteBadRequest(w, "Missing required fields", nil)
			return
		}

		// Create checkout session
		session, err := c.stripeService.CreateCheckoutSession(userID, req.PlanID, req.SuccessURL, req.CancelURL)
		if err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to create checkout session: %v", err), err)
			return
		}

		utils.WriteOK(w, session, "Checkout session created successfully")
	}
}

// WebhookHandler handles Stripe webhook events
func (c *StripeController) WebhookHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			utils.WriteBadRequest(w, "Failed to read request body", err)
			return
		}

		// Verify webhook signature
		event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), c.config.StripeWebhookSecret)
		if err != nil {
			utils.WriteBadRequest(w, "Invalid webhook signature", err)
			return
		}

		// Handle the event
		if err := c.handleWebhookEvent(event); err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to handle webhook: %v", err), err)
			return
		}

		utils.WriteOK(w, nil, "Webhook processed successfully")
	}
}

// GetAvailablePlansHandler returns available subscription plans
func (c *StripeController) GetAvailablePlansHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		plans := c.stripeService.GetAvailablePlans()
		utils.WriteOK(w, plans, "Plans retrieved successfully")
	}
}

// GetUserSubscriptionHandler returns the current user's subscription
func (c *StripeController) GetUserSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		subscription, err := c.subscriptionService.GetUserSubscriptionStatus(userID)
		if err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to get subscription: %v", err), err)
			return
		}

		if subscription == nil {
			utils.WriteOK(w, map[string]interface{}{
				"subscription": nil,
				"status":       "none",
			}, "No active subscription found")
			return
		}

		utils.WriteOK(w, subscription, "Subscription retrieved successfully")
	}
}

// GetUserSubscriptionHistoryHandler returns the user's subscription history
func (c *StripeController) GetUserSubscriptionHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		history, err := c.subscriptionService.GetUserSubscriptionHistory(userID)
		if err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to get subscription history: %v", err), err)
			return
		}

		utils.WriteOK(w, history, "Subscription history retrieved successfully")
	}
}

// GetUserPaymentHistoryHandler returns the user's payment history
func (c *StripeController) GetUserPaymentHistoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		history, err := c.subscriptionService.GetUserPaymentHistory(userID)
		if err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to get payment history: %v", err), err)
			return
		}

		utils.WriteOK(w, history, "Payment history retrieved successfully")
	}
}

// CancelSubscriptionHandler cancels the user's subscription
func (c *StripeController) CancelSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		if err := c.subscriptionService.CancelSubscription(userID); err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to cancel subscription: %v", err), err)
			return
		}

		utils.WriteOK(w, nil, "Subscription cancelled successfully")
	}
}

// ReactivateSubscriptionHandler reactivates a cancelled subscription
func (c *StripeController) ReactivateSubscriptionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteMethodNotAllowed(w, "POST")
			return
		}

		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(int)
		if !ok {
			utils.WriteUnauthorized(w, "User not authenticated")
			return
		}

		if err := c.subscriptionService.ReactivateSubscription(userID); err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to reactivate subscription: %v", err), err)
			return
		}

		utils.WriteOK(w, nil, "Subscription reactivated successfully")
	}
}

// GetSubscriptionMetricsHandler returns subscription metrics (admin only)
func (c *StripeController) GetSubscriptionMetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			utils.WriteMethodNotAllowed(w, "GET")
			return
		}

		// Check if user is admin (this should be handled by RBAC middleware)
		// For now, we'll assume this endpoint is protected by middleware

		metrics, err := c.subscriptionService.GetSubscriptionMetrics()
		if err != nil {
			utils.WriteInternalServerError(w, fmt.Sprintf("Failed to get metrics: %v", err), err)
			return
		}

		utils.WriteOK(w, metrics, "Metrics retrieved successfully")
	}
}

// handleWebhookEvent processes Stripe webhook events
func (c *StripeController) handleWebhookEvent(event stripe.Event) error {
	switch event.Type {
	case "checkout.session.completed":
		return c.handleCheckoutSessionCompleted(event)
	case "customer.subscription.created":
		return c.handleSubscriptionCreated(event)
	case "customer.subscription.updated":
		return c.handleSubscriptionUpdated(event)
	case "customer.subscription.deleted":
		return c.handleSubscriptionDeleted(event)
	case "invoice.payment_succeeded":
		return c.handlePaymentSucceeded(event)
	case "invoice.payment_failed":
		return c.handlePaymentFailed(event)
	default:
		// Log unhandled events
		fmt.Printf("Unhandled event type: %s\n", event.Type)
		return nil
	}
}

// handleCheckoutSessionCompleted processes completed checkout sessions
func (c *StripeController) handleCheckoutSessionCompleted(event stripe.Event) error {
	var session stripe.CheckoutSession
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		return fmt.Errorf("failed to unmarshal checkout session: %w", err)
	}

	// The subscription will be handled by the subscription.created event
	fmt.Printf("Checkout session completed: %s\n", session.ID)
	return nil
}

// handleSubscriptionCreated processes new subscription creation
func (c *StripeController) handleSubscriptionCreated(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Get customer info
	customer, err := c.stripeService.GetCustomerByStripeID(sub.Customer.ID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return fmt.Errorf("customer not found for subscription: %s", sub.ID)
	}

	// Get plan details
	planName := "Unknown Plan"
	if sub.Items != nil && len(sub.Items.Data) > 0 {
		// You might want to fetch plan details from Stripe API here
		planName = fmt.Sprintf("Plan %s", sub.Items.Data[0].Price.ID)
	}

	// Create subscription record
	subData := &models.SubscriptionCreate{
		UserID:           customer.UserID,
		StripeCustomerID: customer.ID,
		StripeSubID:      sub.ID,
		PlanID:           sub.Items.Data[0].Price.ID,
		PlanName:         planName,
	}

	_, err = c.stripeService.CreateSubscription(
		subData,
		time.Unix(sub.CurrentPeriodStart, 0),
		time.Unix(sub.CurrentPeriodEnd, 0),
	)

	if err != nil {
		return fmt.Errorf("failed to create subscription record: %w", err)
	}

	fmt.Printf("Subscription created: %s for user %d\n", sub.ID, customer.UserID)
	return nil
}

// handleSubscriptionUpdated processes subscription updates
func (c *StripeController) handleSubscriptionUpdated(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Update subscription record
	err := c.stripeService.UpdateSubscription(
		sub.ID,
		string(sub.Status),
		time.Unix(sub.CurrentPeriodStart, 0),
		time.Unix(sub.CurrentPeriodEnd, 0),
		sub.CancelAtPeriodEnd,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	fmt.Printf("Subscription updated: %s\n", sub.ID)
	return nil
}

// handleSubscriptionDeleted processes subscription deletion
func (c *StripeController) handleSubscriptionDeleted(event stripe.Event) error {
	var sub stripe.Subscription
	if err := json.Unmarshal(event.Data.Raw, &sub); err != nil {
		return fmt.Errorf("failed to unmarshal subscription: %w", err)
	}

	// Update subscription record
	err := c.stripeService.UpdateSubscription(
		sub.ID,
		"canceled",
		time.Unix(sub.CurrentPeriodStart, 0),
		time.Unix(sub.CurrentPeriodEnd, 0),
		true,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	fmt.Printf("Subscription deleted: %s\n", sub.ID)
	return nil
}

// handlePaymentSucceeded processes successful payments
func (c *StripeController) handlePaymentSucceeded(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	// Get customer info
	customer, err := c.stripeService.GetCustomerByStripeID(invoice.Customer.ID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return fmt.Errorf("customer not found for invoice: %s", invoice.ID)
	}

	// Create payment record
	paymentData := &models.PaymentCreate{
		UserID:           customer.UserID,
		StripeCustomerID: customer.ID,
		StripePaymentID:  invoice.PaymentIntent.ID,
		Amount:           invoice.AmountPaid,
		Currency:         string(invoice.Currency),
		Status:           "succeeded",
		Description:      fmt.Sprintf("Payment for invoice %s", invoice.ID),
	}

	_, err = c.stripeService.CreatePayment(paymentData)
	if err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}

	fmt.Printf("Payment succeeded: %s for user %d\n", invoice.PaymentIntent.ID, customer.UserID)
	return nil
}

// handlePaymentFailed processes failed payments
func (c *StripeController) handlePaymentFailed(event stripe.Event) error {
	var invoice stripe.Invoice
	if err := json.Unmarshal(event.Data.Raw, &invoice); err != nil {
		return fmt.Errorf("failed to unmarshal invoice: %w", err)
	}

	// Get customer info
	customer, err := c.stripeService.GetCustomerByStripeID(invoice.Customer.ID)
	if err != nil {
		return fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return fmt.Errorf("customer not found for invoice: %s", invoice.ID)
	}

	// Create payment record for failed payment
	paymentData := &models.PaymentCreate{
		UserID:           customer.UserID,
		StripeCustomerID: customer.ID,
		StripePaymentID:  invoice.PaymentIntent.ID,
		Amount:           invoice.AmountDue,
		Currency:         string(invoice.Currency),
		Status:           "failed",
		Description:      fmt.Sprintf("Failed payment for invoice %s", invoice.ID),
	}

	_, err = c.stripeService.CreatePayment(paymentData)
	if err != nil {
		return fmt.Errorf("failed to create payment record: %w", err)
	}

	fmt.Printf("Payment failed: %s for user %d\n", invoice.PaymentIntent.ID, customer.UserID)
	return nil
}
