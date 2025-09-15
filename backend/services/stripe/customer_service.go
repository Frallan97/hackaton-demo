package stripe

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/frallan97/hackaton-demo-backend/models"
	"github.com/stripe/stripe-go/v76"
)

// CustomerService handles customer-related operations
type CustomerService struct {
	db           *sql.DB
	stripeClient *StripeClient
}

// NewCustomerService creates a new customer service
func NewCustomerService(db *sql.DB, stripeClient *StripeClient) *CustomerService {
	return &CustomerService{
		db:           db,
		stripeClient: stripeClient,
	}
}

// CreateCustomer creates a new customer in both Stripe and the database
func (s *CustomerService) CreateCustomer(userID int, email, name string) (*models.StripeCustomer, error) {
	// Create customer in Stripe
	stripeParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	stripeCustomer, err := s.stripeClient.CreateCustomer(stripeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create Stripe customer: %w", err)
	}

	// Store customer in database
	dbCustomer, err := s.storeCustomerInDB(userID, stripeCustomer.ID, email)
	if err != nil {
		return nil, fmt.Errorf("failed to store customer in database: %w", err)
	}

	return dbCustomer, nil
}

// GetCustomerByUserID retrieves a customer by user ID
func (s *CustomerService) GetCustomerByUserID(userID int) (*models.StripeCustomer, error) {
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
			return nil, nil // Customer not found
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

// GetCustomerByStripeID retrieves a customer by Stripe ID
func (s *CustomerService) GetCustomerByStripeID(stripeID string) (*models.StripeCustomer, error) {
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
			return nil, nil // Customer not found
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}

	return &customer, nil
}

// GetOrCreateCustomer gets an existing customer or creates a new one
func (s *CustomerService) GetOrCreateCustomer(userID int, email, name string) (*models.StripeCustomer, error) {
	// Try to get existing customer
	customer, err := s.GetCustomerByUserID(userID)
	if err != nil {
		return nil, err
	}

	// If customer exists, return it
	if customer != nil {
		return customer, nil
	}

	// Create new customer
	return s.CreateCustomer(userID, email, name)
}

// UpdateCustomer updates customer information in both Stripe and database
func (s *CustomerService) UpdateCustomer(userID int, email, name string) (*models.StripeCustomer, error) {
	// Get existing customer
	customer, err := s.GetCustomerByUserID(userID)
	if err != nil {
		return nil, err
	}

	if customer == nil {
		return nil, fmt.Errorf("customer not found for user ID: %d", userID)
	}

	// Update in Stripe
	stripeParams := &stripe.CustomerParams{
		Email: stripe.String(email),
		Name:  stripe.String(name),
	}

	_, err = s.stripeClient.UpdateCustomer(customer.StripeID, stripeParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update Stripe customer: %w", err)
	}

	// Update in database
	query := `
		UPDATE stripe_customers 
		SET email = $1, updated_at = $2
		WHERE user_id = $3
		RETURNING id, user_id, stripe_id, email, default_source, created_at, updated_at
	`

	var updatedCustomer models.StripeCustomer
	err = s.db.QueryRow(query, email, time.Now(), userID).Scan(
		&updatedCustomer.ID,
		&updatedCustomer.UserID,
		&updatedCustomer.StripeID,
		&updatedCustomer.Email,
		&updatedCustomer.DefaultSource,
		&updatedCustomer.CreatedAt,
		&updatedCustomer.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update customer in database: %w", err)
	}

	return &updatedCustomer, nil
}

// DeleteCustomer removes a customer from both Stripe and database
func (s *CustomerService) DeleteCustomer(userID int) error {
	// Get customer first
	customer, err := s.GetCustomerByUserID(userID)
	if err != nil {
		return err
	}

	if customer == nil {
		return fmt.Errorf("customer not found for user ID: %d", userID)
	}

	// Note: Stripe doesn't allow deleting customers, only updating them
	// We'll just remove from our database
	query := `DELETE FROM stripe_customers WHERE user_id = $1`
	_, err = s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete customer from database: %w", err)
	}

	return nil
}

// ListCustomers lists all customers with pagination
func (s *CustomerService) ListCustomers(offset, limit int) ([]*models.StripeCustomer, error) {
	query := `
		SELECT id, user_id, stripe_id, email, default_source, created_at, updated_at
		FROM stripe_customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer rows.Close()

	var customers []*models.StripeCustomer
	for rows.Next() {
		var customer models.StripeCustomer
		err := rows.Scan(
			&customer.ID,
			&customer.UserID,
			&customer.StripeID,
			&customer.Email,
			&customer.DefaultSource,
			&customer.CreatedAt,
			&customer.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan customer: %w", err)
		}
		customers = append(customers, &customer)
	}

	return customers, nil
}

// GetCustomerCount returns the total number of customers
func (s *CustomerService) GetCustomerCount() (int, error) {
	query := `SELECT COUNT(*) FROM stripe_customers`
	var count int
	err := s.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get customer count: %w", err)
	}
	return count, nil
}

// Private helper methods

func (s *CustomerService) storeCustomerInDB(userID int, stripeID, email string) (*models.StripeCustomer, error) {
	query := `
		INSERT INTO stripe_customers (user_id, stripe_id, email, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $4)
		RETURNING id, user_id, stripe_id, email, default_source, created_at, updated_at
	`

	var customer models.StripeCustomer
	err := s.db.QueryRow(
		query,
		userID,
		stripeID,
		email,
		time.Now(),
	).Scan(
		&customer.ID,
		&customer.UserID,
		&customer.StripeID,
		&customer.Email,
		&customer.DefaultSource,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &customer, nil
}
