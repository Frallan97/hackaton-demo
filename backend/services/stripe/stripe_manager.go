package stripe

import (
	"database/sql"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/stripe/stripe-go/v76"
)

// StripeManager orchestrates all Stripe-related services
type StripeManager struct {
	// Core services
	Client   *StripeClient
	Customer *CustomerService
	Payment  *PaymentService
	Plan     *PlanService

	// Future services (ready for extension)
	// Subscription *SubscriptionService
	// Webhook     *WebhookService
	// Analytics   *AnalyticsService
}

// NewStripeManager creates a new Stripe manager with all services
func NewStripeManager(db *sql.DB, config *config.Config) *StripeManager {
	// Initialize core client
	stripeClient := NewStripeClient(config)

	// Initialize services
	customerService := NewCustomerService(db, stripeClient)
	paymentService := NewPaymentService(db, stripeClient, customerService)
	planService := NewPlanService(db, stripeClient)

	return &StripeManager{
		Client:   stripeClient,
		Customer: customerService,
		Payment:  paymentService,
		Plan:     planService,
	}
}

// Health check for Stripe integration
func (sm *StripeManager) HealthCheck() error {
	// Test basic Stripe connectivity by listing products
	params := &stripe.ProductListParams{}
	params.Limit = stripe.Int64(1)

	iter := sm.Client.ListProducts(params)
	if iter.Err() != nil {
		return iter.Err()
	}

	return nil
}

// Future: Add service registration methods for extensibility

// RegisterSubscriptionService adds subscription service (future implementation)
func (sm *StripeManager) RegisterSubscriptionService(service interface{}) {
	// TODO: Implement subscription service registration
	// sm.Subscription = service.(*SubscriptionService)
}

// RegisterWebhookService adds webhook service (future implementation)
func (sm *StripeManager) RegisterWebhookService(service interface{}) {
	// TODO: Implement webhook service registration
	// sm.Webhook = service.(*WebhookService)
}

// RegisterAnalyticsService adds analytics service (future implementation)
func (sm *StripeManager) RegisterAnalyticsService(service interface{}) {
	// TODO: Implement analytics service registration
	// sm.Analytics = service.(*AnalyticsService)
}
