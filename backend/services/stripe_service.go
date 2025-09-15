package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
)

// StripeService handles all Stripe-related operations
type StripeService struct {
	db     *sql.DB
	config *config.Config
}

// NewStripeService creates a new Stripe service
func NewStripeService(db *sql.DB, config *config.Config) *StripeService {
	// Set Stripe API key
	stripe.Key = config.StripeSecretKey

	return &StripeService{
		db:     db,
		config: config,
	}
}

// CreateCustomer creates a new Stripe customer
func (s *StripeService) CreateCustomer(userID int, email, name string) (*models.StripeCustomer, error) {
	// Create customer in Stripe
	customerParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	stripeCustomer, err := customer.New(customerParams)
	if err != nil {
		log.Printf("Failed to create Stripe customer: %v", err)
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Store customer in database
	query := `
		INSERT INTO stripe_customers (user_id, stripe_id, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, user_id, stripe_id, email, default_source, created_at, updated_at
	`

	var dbCustomer models.StripeCustomer
	err = s.db.QueryRow(
		query,
		userID,
		stripeCustomer.ID,
		email,
		time.Now(),
	).Scan(
		&dbCustomer.ID,
		&dbCustomer.UserID,
		&dbCustomer.StripeID,
		&dbCustomer.Email,
		&dbCustomer.DefaultSource,
		&dbCustomer.CreatedAt,
		&dbCustomer.UpdatedAt,
	)

	if err != nil {
		log.Printf("Failed to store Stripe customer in database: %v", err)
		return nil, fmt.Errorf("failed to store customer: %w", err)
	}

	return &dbCustomer, nil
}

// GetCustomerByUserID retrieves a Stripe customer by user ID
func (s *StripeService) GetCustomerByUserID(userID int) (*models.StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_id, email, default_source, created_at, updated_at
		FROM stripe_customers
		WHERE user_id = $1
	`

	var customer models.StripeCustomer
	err := s.db.QueryRow(query, userID).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.StripeID,
		&customer.Email,
		&customer.DefaultSource,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

// GetCustomerByStripeID retrieves a Stripe customer by Stripe ID
func (s *StripeService) GetCustomerByStripeID(stripeID string) (*models.StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_id, email, default_source, created_at, updated_at
		FROM stripe_customers
		WHERE stripe_id = $1
	`

	var customer models.StripeCustomer
	err := s.db.QueryRow(query, stripeID).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.StripeID,
		&customer.Email,
		&customer.DefaultSource,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

// CreateCheckoutSession creates a new Stripe checkout session
func (s *StripeService) CreateCheckoutSession(userID int, planID, successURL, cancelURL string) (*models.CreateCheckoutSessionResponse, error) {
	// Get or create customer
	customer, err := s.GetCustomerByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		// Get user info to create customer
		var email, name string
		err = s.db.QueryRow("SELECT email, name FROM users WHERE id = $1", userID).Scan(&email, &name)
		if err != nil {
			return nil, fmt.Errorf("failed to get user info: %w", err)
		}

		customer, err = s.CreateCustomer(userID, email, name)
		if err != nil {
			return nil, fmt.Errorf("failed to create customer: %w", err)
		}
	}

	// Create checkout session
	sessionParams := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customer.StripeID),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
			// Note: Swish requires special setup in Stripe dashboard and is region-specific
			// stripe.String("swish"), // Uncomment when Swish is enabled in your Stripe account
		},
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(planID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(successURL),
		CancelURL:  stripe.String(cancelURL),
	}

	session, err := session.New(sessionParams)
	if err != nil {
		log.Printf("Failed to create checkout session: %v", err)
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &models.CreateCheckoutSessionResponse{
		SessionID: session.ID,
		URL:       session.URL,
	}, nil
}

// GetSubscription retrieves a subscription by Stripe subscription ID
func (s *StripeService) GetSubscription(stripeSubID string) (*models.Subscription, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_sub_id, status, plan_id, plan_name,
		       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		FROM subscriptions
		WHERE stripe_sub_id = $1
	`

	var sub models.Subscription
	err := s.db.QueryRow(query, stripeSubID).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.StripeCustomerID,
		&sub.StripeSubID,
		&sub.Status,
		&sub.PlanID,
		&sub.PlanName,
		&sub.CurrentPeriodStart,
		&sub.CurrentPeriodEnd,
		&sub.CancelAtPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

// CreateSubscription creates a new subscription record
func (s *StripeService) CreateSubscription(subData *models.SubscriptionCreate, periodStart, periodEnd time.Time) (*models.Subscription, error) {
	query := `
		INSERT INTO subscriptions (user_id, stripe_customer_id, stripe_sub_id, status, plan_id, plan_name,
		                         current_period_start, current_period_end, created_at, updated_at)
		VALUES ($1, $2, $3, 'active', $4, $5, $6, $7, $8, $8)
		RETURNING id, user_id, stripe_customer_id, stripe_sub_id, status, plan_id, plan_name,
		          current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
	`

	var sub models.Subscription
	err := s.db.QueryRow(
		query,
		subData.UserID,
		subData.StripeCustomerID,
		subData.StripeSubID,
		subData.PlanID,
		subData.PlanName,
		periodStart,
		periodEnd,
		time.Now(),
	).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.StripeCustomerID,
		&sub.StripeSubID,
		&sub.Status,
		&sub.PlanID,
		&sub.PlanName,
		&sub.CurrentPeriodStart,
		&sub.CurrentPeriodEnd,
		&sub.CancelAtPeriodEnd,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		log.Printf("Failed to create subscription: %v", err)
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Update user subscription status
	_, err = s.db.Exec(`
		UPDATE users 
		SET subscription_status = 'active', subscription_plan = $1, subscription_expires_at = $2
		WHERE id = $3
	`, subData.PlanName, periodEnd, subData.UserID)

	if err != nil {
		log.Printf("Warning: Failed to update user subscription status: %v", err)
	}

	return &sub, nil
}

// UpdateSubscription updates an existing subscription
func (s *StripeService) UpdateSubscription(stripeSubID, status string, periodStart, periodEnd time.Time, cancelAtPeriodEnd bool) error {
	query := `
		UPDATE subscriptions 
		SET status = $1, current_period_start = $2, current_period_end = $3, 
		    cancel_at_period_end = $4, updated_at = $5
		WHERE stripe_sub_id = $6
	`

	result, err := s.db.Exec(query, status, periodStart, periodEnd, cancelAtPeriodEnd, time.Now(), stripeSubID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found: %s", stripeSubID)
	}

	// Update user subscription status if subscription is cancelled
	if status == "canceled" || status == "unpaid" {
		_, err = s.db.Exec(`
			UPDATE users 
			SET subscription_status = 'inactive', subscription_expires_at = $1
			WHERE id = (SELECT user_id FROM subscriptions WHERE stripe_sub_id = $2)
		`, time.Now(), stripeSubID)

		if err != nil {
			log.Printf("Warning: Failed to update user subscription status: %v", err)
		}
	}

	return nil
}

// CreatePayment creates a new payment record
func (s *StripeService) CreatePayment(paymentData *models.PaymentCreate) (*models.Payment, error) {
	query := `
		INSERT INTO payments (user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at
	`

	var payment models.Payment
	err := s.db.QueryRow(
		query,
		paymentData.UserID,
		paymentData.StripeCustomerID,
		paymentData.StripePaymentID,
		paymentData.Amount,
		paymentData.Currency,
		paymentData.Status,
		paymentData.Description,
		time.Now(),
	).Scan(
		&payment.ID,
		&payment.UserID,
		&payment.StripeCustomerID,
		&payment.StripePaymentID,
		&payment.Amount,
		&payment.Currency,
		&payment.Status,
		&payment.Description,
		&payment.CreatedAt,
	)

	if err != nil {
		log.Printf("Failed to create payment: %v", err)
		return nil, fmt.Errorf("failed to create payment: %w", err)
	}

	return &payment, nil
}

// GetUserSubscriptions retrieves all subscriptions for a user
func (s *StripeService) GetUserSubscriptions(userID int) ([]*models.Subscription, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_sub_id, status, plan_id, plan_name,
		       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*models.Subscription
	for rows.Next() {
		var sub models.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.UserID,
			&sub.StripeCustomerID,
			&sub.StripeSubID,
			&sub.Status,
			&sub.PlanID,
			&sub.PlanName,
			&sub.CurrentPeriodStart,
			&sub.CurrentPeriodEnd,
			&sub.CancelAtPeriodEnd,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, &sub)
	}

	return subscriptions, nil
}

// GetUserPayments retrieves all payments for a user
func (s *StripeService) GetUserPayments(userID int) ([]*models.Payment, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at
		FROM payments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query payments: %w", err)
	}
	defer rows.Close()

	var payments []*models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(
			&payment.ID,
			&payment.UserID,
			&payment.StripeCustomerID,
			&payment.StripePaymentID,
			&payment.Amount,
			&payment.Currency,
			&payment.Status,
			&payment.Description,
			&payment.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan payment: %w", err)
		}
		payments = append(payments, &payment)
	}

	return payments, nil
}

// GetAvailablePlans returns available payment plans
func (s *StripeService) GetAvailablePlans() []*models.PaymentPlan {
	// This could be fetched from Stripe API or stored in database
	// For now, returning a single test plan
	return []*models.PaymentPlan{
		{
			ID:          "price_1S7hcfAeXvIjnXEPpXj1morV",
			Name:        "Test Payment",
			Description: "Test payment with card and Swish support",
			Price:       999, // $9.99 in cents
			Currency:    "usd",
			Features:    []string{"Test payment functionality", "Card payments", "Swish payments", "Payment history"},
		},
	}
}
