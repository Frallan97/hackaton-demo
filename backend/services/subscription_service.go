package services

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/frallan97/hackaton-demo-backend/models"
)

// SubscriptionService handles subscription business logic
type SubscriptionService struct {
	db            *sql.DB
	stripeService *StripeService
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(db *sql.DB, stripeService *StripeService) *SubscriptionService {
	return &SubscriptionService{
		db:            db,
		stripeService: stripeService,
	}
}

// GetUserSubscriptionStatus returns the current subscription status for a user
func (s *SubscriptionService) GetUserSubscriptionStatus(userID int) (*models.Subscription, error) {
	// Get the most recent active subscription
	query := `
		SELECT id, user_id, stripe_customer_id, stripe_sub_id, status, plan_id, plan_name,
		       current_period_start, current_period_end, cancel_at_period_end, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1 AND status = 'active'
		ORDER BY created_at DESC
		LIMIT 1
	`

	var sub models.Subscription
	err := s.db.QueryRow(query, userID).Scan(
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
		return nil, fmt.Errorf("failed to get subscription status: %w", err)
	}

	return &sub, nil
}

// IsUserSubscribed checks if a user has an active subscription
func (s *SubscriptionService) IsUserSubscribed(userID int) (bool, error) {
	sub, err := s.GetUserSubscriptionStatus(userID)
	if err != nil {
		return false, err
	}
	return sub != nil && sub.Status == "active", nil
}

// HasUserAccess checks if a user has access to a specific feature based on their subscription
func (s *SubscriptionService) HasUserAccess(userID int, requiredPlan string) (bool, error) {
	sub, err := s.GetUserSubscriptionStatus(userID)
	if err != nil {
		return false, err
	}

	if sub == nil || sub.Status != "active" {
		return false, nil
	}

	// Check if user's plan meets the requirement
	switch requiredPlan {
	case "basic":
		return true, nil // All active subscriptions have basic access
	case "pro":
		return sub.PlanName == "Pro Plan" || sub.PlanName == "Enterprise Plan", nil
	case "enterprise":
		return sub.PlanName == "Enterprise Plan", nil
	default:
		return false, fmt.Errorf("unknown required plan: %s", requiredPlan)
	}
}

// GetUserSubscriptionHistory returns all subscription history for a user
func (s *SubscriptionService) GetUserSubscriptionHistory(userID int) ([]*models.Subscription, error) {
	return s.stripeService.GetUserSubscriptions(userID)
}

// GetUserPaymentHistory returns all payment history for a user
func (s *SubscriptionService) GetUserPaymentHistory(userID int) ([]*models.Payment, error) {
	return s.stripeService.GetUserPayments(userID)
}

// CancelSubscription cancels a user's subscription at the end of the current period
func (s *SubscriptionService) CancelSubscription(userID int) error {
	sub, err := s.GetUserSubscriptionStatus(userID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		return fmt.Errorf("no active subscription found")
	}

	// Update subscription to cancel at period end
	query := `
		UPDATE subscriptions 
		SET cancel_at_period_end = true, updated_at = $1
		WHERE id = $2
	`

	_, err = s.db.Exec(query, time.Now(), sub.ID)
	if err != nil {
		return fmt.Errorf("failed to cancel subscription: %w", err)
	}

	log.Printf("Subscription %s marked for cancellation at period end", sub.StripeSubID)
	return nil
}

// ReactivateSubscription reactivates a cancelled subscription
func (s *SubscriptionService) ReactivateSubscription(userID int) error {
	sub, err := s.GetUserSubscriptionStatus(userID)
	if err != nil {
		return fmt.Errorf("failed to get subscription: %w", err)
	}

	if sub == nil {
		return fmt.Errorf("no active subscription found")
	}

	if !sub.CancelAtPeriodEnd {
		return fmt.Errorf("subscription is not cancelled")
	}

	// Update subscription to not cancel at period end
	query := `
		UPDATE subscriptions 
		SET cancel_at_period_end = false, updated_at = $1
		WHERE id = $2
	`

	_, err = s.db.Exec(query, time.Now(), sub.ID)
	if err != nil {
		return fmt.Errorf("failed to reactivate subscription: %w", err)
	}

	log.Printf("Subscription %s reactivated", sub.StripeSubID)
	return nil
}

// GetSubscriptionMetrics returns subscription metrics for admin purposes
func (s *SubscriptionService) GetSubscriptionMetrics() (map[string]interface{}, error) {
	metrics := make(map[string]interface{})

	// Total active subscriptions
	var activeCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE status = 'active'").Scan(&activeCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscription count: %w", err)
	}
	metrics["active_subscriptions"] = activeCount

	// Total cancelled subscriptions
	var cancelledCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM subscriptions WHERE status = 'canceled'").Scan(&cancelledCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get cancelled subscription count: %w", err)
	}
	metrics["cancelled_subscriptions"] = cancelledCount

	// Total revenue (sum of all successful payments)
	var totalRevenue int64
	err = s.db.QueryRow("SELECT COALESCE(SUM(amount), 0) FROM payments WHERE status = 'succeeded'").Scan(&totalRevenue)
	if err != nil {
		return nil, fmt.Errorf("failed to get total revenue: %w", err)
	}
	metrics["total_revenue_cents"] = totalRevenue

	// Plan distribution
	planQuery := `
		SELECT plan_name, COUNT(*) 
		FROM subscriptions 
		WHERE status = 'active' 
		GROUP BY plan_name
	`
	rows, err := s.db.Query(planQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan distribution: %w", err)
	}
	defer rows.Close()

	planDistribution := make(map[string]int)
	for rows.Next() {
		var planName string
		var count int
		err := rows.Scan(&planName, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan plan distribution: %w", err)
		}
		planDistribution[planName] = count
	}
	metrics["plan_distribution"] = planDistribution

	return metrics, nil
}

// CleanupExpiredSubscriptions removes expired subscriptions and updates user status
func (s *SubscriptionService) CleanupExpiredSubscriptions() error {
	// Find expired subscriptions
	query := `
		SELECT id, user_id, stripe_sub_id
		FROM subscriptions
		WHERE status = 'active' AND current_period_end < $1
	`

	rows, err := s.db.Query(query, time.Now())
	if err != nil {
		return fmt.Errorf("failed to query expired subscriptions: %w", err)
	}
	defer rows.Close()

	var expiredSubs []struct {
		ID          int
		UserID      int
		StripeSubID string
	}

	for rows.Next() {
		var sub struct {
			ID          int
			UserID      int
			StripeSubID string
		}
		err := rows.Scan(&sub.ID, &sub.UserID, &sub.StripeSubID)
		if err != nil {
			return fmt.Errorf("failed to scan expired subscription: %w", err)
		}
		expiredSubs = append(expiredSubs, sub)
	}

	// Update expired subscriptions
	for _, sub := range expiredSubs {
		// Update subscription status
		_, err := s.db.Exec(`
			UPDATE subscriptions 
			SET status = 'expired', updated_at = $1
			WHERE id = $2
		`, time.Now(), sub.ID)
		if err != nil {
			log.Printf("Failed to update expired subscription %d: %v", sub.ID, err)
			continue
		}

		// Update user subscription status
		_, err = s.db.Exec(`
			UPDATE users 
			SET subscription_status = 'expired', subscription_expires_at = $1
			WHERE id = $2
		`, time.Now(), sub.UserID)
		if err != nil {
			log.Printf("Failed to update user %d subscription status: %v", sub.UserID, err)
		}

		log.Printf("Marked subscription %s as expired", sub.StripeSubID)
	}

	return nil
}
