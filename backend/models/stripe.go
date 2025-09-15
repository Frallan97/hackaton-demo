package models

import (
	"database/sql"
	"time"
)

// StripeCustomer represents a Stripe customer linked to a user
type StripeCustomer struct {
	ID            int            `json:"id" db:"id"`
	UserID        int            `json:"user_id" db:"user_id"`
	StripeID      string         `json:"stripe_id" db:"stripe_id"`
	Email         string         `json:"email" db:"email"`
	DefaultSource sql.NullString `json:"default_source" db:"default_source"`
	CreatedAt     time.Time      `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at" db:"updated_at"`
}

// StripeCustomerCreate represents the data needed to create a new Stripe customer
type StripeCustomerCreate struct {
	UserID   int    `json:"user_id"`
	StripeID string `json:"stripe_id"`
	Email    string `json:"email"`
}

// Subscription represents a user subscription
type Subscription struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	StripeCustomerID   int       `json:"stripe_customer_id" db:"stripe_customer_id"`
	StripeSubID        string    `json:"stripe_sub_id" db:"stripe_sub_id"`
	Status             string    `json:"status" db:"status"`
	PlanID             string    `json:"plan_id" db:"plan_id"`
	PlanName           string    `json:"plan_name" db:"plan_name"`
	CurrentPeriodStart time.Time `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end" db:"current_period_end"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end" db:"cancel_at_period_end"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// SubscriptionCreate represents the data needed to create a new subscription
type SubscriptionCreate struct {
	UserID           int    `json:"user_id"`
	StripeCustomerID int    `json:"stripe_customer_id"`
	StripeSubID      string `json:"stripe_sub_id"`
	PlanID           string `json:"plan_id"`
	PlanName         string `json:"plan_name"`
}

// Payment represents a payment record
type Payment struct {
	ID               int       `json:"id" db:"id"`
	UserID           int       `json:"user_id" db:"user_id"`
	StripeCustomerID int       `json:"stripe_customer_id" db:"stripe_customer_id"`
	StripePaymentID  string    `json:"stripe_payment_id" db:"stripe_payment_id"`
	Amount           int64     `json:"amount" db:"amount"`
	Currency         string    `json:"currency" db:"currency"`
	Status           string    `json:"status" db:"status"`
	Description      string    `json:"description" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// PaymentCreate represents the data needed to create a new payment record
type PaymentCreate struct {
	UserID           int    `json:"user_id"`
	StripeCustomerID int    `json:"stripe_customer_id"`
	StripePaymentID  string `json:"stripe_payment_id"`
	Amount           int64  `json:"amount"`
	Currency         string `json:"currency"`
	Status           string `json:"status"`
	Description      string `json:"description"`
}

// PaymentPlan represents a one-time payment plan
type PaymentPlan struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       int64    `json:"price"`
	Currency    string   `json:"currency"`
	Features    []string `json:"features"`
}

// CreateCheckoutSessionRequest represents a request to create a checkout session
type CreateCheckoutSessionRequest struct {
	PlanID     string `json:"plan_id" validate:"required"`
	SuccessURL string `json:"success_url" validate:"required"`
	CancelURL  string `json:"cancel_url" validate:"required"`
}

// CreateCheckoutSessionResponse represents the response from creating a checkout session
type CreateCheckoutSessionResponse struct {
	SessionID string `json:"session_id"`
	URL       string `json:"url"`
}

// PaymentMetrics represents payment analytics data
type PaymentMetrics struct {
	TotalPayments     int            `json:"total_payments"`
	TotalRevenueCents int            `json:"total_revenue_cents"`
	PlanDistribution  map[string]int `json:"plan_distribution"`
}

// WebhookEvent represents a Stripe webhook event
type WebhookEvent struct {
	Type string                 `json:"type"`
	Data map[string]interface{} `json:"data"`
}
