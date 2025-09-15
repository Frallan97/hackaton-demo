package stripe

import (
	"database/sql"
	"fmt"

	"github.com/frallan97/hackaton-demo-backend/models"
)

// PlanService handles plan-related operations
type PlanService struct {
	db           *sql.DB
	stripeClient *StripeClient
}

// NewPlanService creates a new plan service
func NewPlanService(db *sql.DB, stripeClient *StripeClient) *PlanService {
	return &PlanService{
		db:           db,
		stripeClient: stripeClient,
	}
}

// GetAvailablePlans returns available payment plans
// This can be extended to fetch from database or Stripe API
func (s *PlanService) GetAvailablePlans() []*models.PaymentPlan {
	// For now, returning hardcoded plans
	// TODO: Implement database storage and Stripe API fetching
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

// GetPlanByID retrieves a specific plan by ID
func (s *PlanService) GetPlanByID(planID string) (*models.PaymentPlan, error) {
	plans := s.GetAvailablePlans()
	for _, plan := range plans {
		if plan.ID == planID {
			return plan, nil
		}
	}
	return nil, fmt.Errorf("plan not found: %s", planID)
}

// CreatePlanFromStripe creates a plan from Stripe price data
func (s *PlanService) CreatePlanFromStripe(priceID string) (*models.PaymentPlan, error) {
	// Get price from Stripe
	price, err := s.stripeClient.GetPrice(priceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price from Stripe: %w", err)
	}

	// Get product information
	product, err := s.stripeClient.GetProduct(price.Product.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product from Stripe: %w", err)
	}

	// Convert to our plan model
	plan := &models.PaymentPlan{
		ID:          price.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       price.UnitAmount,
		Currency:    string(price.Currency),
		Features:    []string{}, // Features would need to be stored separately
	}

	return plan, nil
}

// ValidatePlan validates that a plan exists and is available
func (s *PlanService) ValidatePlan(planID string) error {
	plan, err := s.GetPlanByID(planID)
	if err != nil {
		return err
	}

	if plan == nil {
		return fmt.Errorf("plan not found: %s", planID)
	}

	// Additional validation logic can be added here
	// e.g., check if plan is active, not expired, etc.

	return nil
}

// Future: Database-backed plan management

// StorePlanInDB stores a plan in the database (for future implementation)
func (s *PlanService) StorePlanInDB(plan *models.PaymentPlan) error {
	// TODO: Implement database storage
	// This would allow for dynamic plan management through admin interface
	return fmt.Errorf("database storage not implemented yet")
}

// GetPlansFromDB retrieves plans from database (for future implementation)
func (s *PlanService) GetPlansFromDB() ([]*models.PaymentPlan, error) {
	// TODO: Implement database retrieval
	return nil, fmt.Errorf("database retrieval not implemented yet")
}

// UpdatePlanInDB updates a plan in the database (for future implementation)
func (s *PlanService) UpdatePlanInDB(planID string, updates *models.PaymentPlan) error {
	// TODO: Implement database updates
	return fmt.Errorf("database updates not implemented yet")
}

// DeletePlanFromDB removes a plan from the database (for future implementation)
func (s *PlanService) DeletePlanFromDB(planID string) error {
	// TODO: Implement database deletion
	return fmt.Errorf("database deletion not implemented yet")
}

// SyncPlansWithStripe synchronizes local plans with Stripe (for future implementation)
func (s *PlanService) SyncPlansWithStripe() error {
	// TODO: Implement Stripe synchronization
	// This would fetch all prices from Stripe and update local database
	return fmt.Errorf("Stripe synchronization not implemented yet")
}

// Future: Advanced plan features

// GetPlansByCategory retrieves plans by category (for future implementation)
func (s *PlanService) GetPlansByCategory(category string) ([]*models.PaymentPlan, error) {
	// TODO: Implement category filtering
	return nil, fmt.Errorf("category filtering not implemented yet")
}

// GetFeaturedPlans retrieves featured plans (for future implementation)
func (s *PlanService) GetFeaturedPlans() ([]*models.PaymentPlan, error) {
	// TODO: Implement featured plans
	return nil, fmt.Errorf("featured plans not implemented yet")
}

// GetPlanRecommendations gets recommended plans for a user (for future implementation)
func (s *PlanService) GetPlanRecommendations(userID int) ([]*models.PaymentPlan, error) {
	// TODO: Implement recommendation engine
	return nil, fmt.Errorf("plan recommendations not implemented yet")
}
