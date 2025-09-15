package stripe

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/stripe/stripe-go/v76"
)

// PaymentService handles payment-related operations
type PaymentService struct {
	db              *sql.DB
	stripeClient    *StripeClient
	customerService *CustomerService
}

// NewPaymentService creates a new payment service
func NewPaymentService(db *sql.DB, stripeClient *StripeClient, customerService *CustomerService) *PaymentService {
	return &PaymentService{
		db:              db,
		stripeClient:    stripeClient,
		customerService: customerService,
	}
}

// CreateCheckoutSession creates a new Stripe checkout session
func (s *PaymentService) CreateCheckoutSession(userID int, planID, successURL, cancelURL string) (*models.CreateCheckoutSessionResponse, error) {
	// Get user info to create/get customer
	var email, name string
	err := s.db.QueryRow("SELECT email, name FROM users WHERE id = $1", userID).Scan(&email, &name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Get or create customer
	customer, err := s.customerService.GetOrCreateCustomer(userID, email, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create customer: %w", err)
	}

	// Create checkout session parameters
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

	// Create session in Stripe
	session, err := s.stripeClient.CreateCheckoutSession(sessionParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create checkout session: %w", err)
	}

	return &models.CreateCheckoutSessionResponse{
		SessionID: session.ID,
		URL:       session.URL,
	}, nil
}

// CreatePaymentIntent creates a new payment intent for direct payments
func (s *PaymentService) CreatePaymentIntent(userID int, amount int64, currency string) (*stripe.PaymentIntent, error) {
	// Get user info
	var email, name string
	err := s.db.QueryRow("SELECT email, name FROM users WHERE id = $1", userID).Scan(&email, &name)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Get or create customer
	customer, err := s.customerService.GetOrCreateCustomer(userID, email, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create customer: %w", err)
	}

	// Create payment intent parameters
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(amount),
		Currency: stripe.String(currency),
		Customer: stripe.String(customer.StripeID),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	// Create payment intent in Stripe
	intent, err := s.stripeClient.CreatePaymentIntent(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment intent: %w", err)
	}

	return intent, nil
}

// RecordPayment records a successful payment in the database
func (s *PaymentService) RecordPayment(userID int, stripePaymentID string, amount int64, currency, status, description string) (*models.Payment, error) {
	// Get customer
	customer, err := s.customerService.GetCustomerByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	if customer == nil {
		return nil, fmt.Errorf("customer not found for user ID: %d", userID)
	}

	// Insert payment record
	query := `
		INSERT INTO payments (user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at
	`

	var payment models.Payment
	err = s.db.QueryRow(
		query,
		userID,
		customer.ID,
		stripePaymentID,
		amount,
		currency,
		status,
		description,
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
		return nil, fmt.Errorf("failed to record payment: %w", err)
	}

	return &payment, nil
}

// GetPaymentHistory retrieves payment history for a user
func (s *PaymentService) GetPaymentHistory(userID int) ([]*models.Payment, error) {
	query := `
		SELECT p.id, p.user_id, p.stripe_customer_id, p.stripe_payment_id, 
		       p.amount, p.currency, p.status, p.description, p.created_at
		FROM payments p
		INNER JOIN stripe_customers sc ON p.stripe_customer_id = sc.id
		WHERE p.user_id = $1
		ORDER BY p.created_at DESC
	`

	rows, err := s.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment history: %w", err)
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

// GetPaymentByStripeID retrieves a payment by Stripe payment ID
func (s *PaymentService) GetPaymentByStripeID(stripePaymentID string) (*models.Payment, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at
		FROM payments
		WHERE stripe_payment_id = $1
	`

	var payment models.Payment
	err := s.db.QueryRow(query, stripePaymentID).Scan(
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
		if err == sql.ErrNoRows {
			return nil, nil // Payment not found
		}
		return nil, fmt.Errorf("failed to get payment: %w", err)
	}

	return &payment, nil
}

// UpdatePaymentStatus updates the status of a payment
func (s *PaymentService) UpdatePaymentStatus(stripePaymentID, status string) error {
	query := `UPDATE payments SET status = $1 WHERE stripe_payment_id = $2`

	result, err := s.db.Exec(query, status, stripePaymentID)
	if err != nil {
		return fmt.Errorf("failed to update payment status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment not found with Stripe ID: %s", stripePaymentID)
	}

	return nil
}

// GetPaymentMetrics returns payment metrics for analytics
func (s *PaymentService) GetPaymentMetrics() (*models.PaymentMetrics, error) {
	// Get total payments count
	var totalPayments int
	err := s.db.QueryRow("SELECT COUNT(*) FROM payments WHERE status = 'succeeded'").Scan(&totalPayments)
	if err != nil {
		return nil, fmt.Errorf("failed to get total payments: %w", err)
	}

	// Get total revenue
	var totalRevenue sql.NullInt64
	err = s.db.QueryRow("SELECT SUM(amount) FROM payments WHERE status = 'succeeded'").Scan(&totalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}

	// Get plan distribution (this would need to be enhanced based on your plan tracking)
	planDistribution := make(map[string]int)
	// This is a simplified version - you might want to track plan information in payments table

	metrics := &models.PaymentMetrics{
		TotalPayments:     totalPayments,
		TotalRevenueCents: int(totalRevenue.Int64),
		PlanDistribution:  planDistribution,
	}

	return metrics, nil
}

// ListAllPayments lists all payments with pagination (admin function)
func (s *PaymentService) ListAllPayments(offset, limit int) ([]*models.Payment, error) {
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_payment_id, amount, currency, status, description, created_at
		FROM payments
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list payments: %w", err)
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
