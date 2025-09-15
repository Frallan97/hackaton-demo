package stripe

import (
	"log"

	"github.com/frallan97/hackaton-demo-backend/config"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/paymentintent"
	"github.com/stripe/stripe-go/v76/price"
	"github.com/stripe/stripe-go/v76/product"
	"github.com/stripe/stripe-go/v76/subscription"
)

// StripeClient handles all direct Stripe API interactions
type StripeClient struct {
	config *config.Config
}

// NewStripeClient creates a new Stripe client
func NewStripeClient(config *config.Config) *StripeClient {
	// Set Stripe API key
	stripe.Key = config.StripeSecretKey

	return &StripeClient{
		config: config,
	}
}

// Customer operations
func (c *StripeClient) CreateCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	customer, err := customer.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreateCustomer: %v", err)
		return nil, err
	}
	return customer, nil
}

func (c *StripeClient) GetCustomer(customerID string) (*stripe.Customer, error) {
	customer, err := customer.Get(customerID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetCustomer: %v", err)
		return nil, err
	}
	return customer, nil
}

func (c *StripeClient) UpdateCustomer(customerID string, params *stripe.CustomerParams) (*stripe.Customer, error) {
	customer, err := customer.Update(customerID, params)
	if err != nil {
		log.Printf("Stripe API error - UpdateCustomer: %v", err)
		return nil, err
	}
	return customer, nil
}

// Product operations
func (c *StripeClient) CreateProduct(params *stripe.ProductParams) (*stripe.Product, error) {
	product, err := product.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreateProduct: %v", err)
		return nil, err
	}
	return product, nil
}

func (c *StripeClient) GetProduct(productID string) (*stripe.Product, error) {
	product, err := product.Get(productID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetProduct: %v", err)
		return nil, err
	}
	return product, nil
}

func (c *StripeClient) ListProducts(params *stripe.ProductListParams) *product.Iter {
	return product.List(params)
}

// Price operations
func (c *StripeClient) CreatePrice(params *stripe.PriceParams) (*stripe.Price, error) {
	price, err := price.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreatePrice: %v", err)
		return nil, err
	}
	return price, nil
}

func (c *StripeClient) GetPrice(priceID string) (*stripe.Price, error) {
	price, err := price.Get(priceID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetPrice: %v", err)
		return nil, err
	}
	return price, nil
}

func (c *StripeClient) ListPrices(params *stripe.PriceListParams) *price.Iter {
	return price.List(params)
}

// Checkout Session operations
func (c *StripeClient) CreateCheckoutSession(params *stripe.CheckoutSessionParams) (*stripe.CheckoutSession, error) {
	session, err := session.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreateCheckoutSession: %v", err)
		return nil, err
	}
	return session, nil
}

func (c *StripeClient) GetCheckoutSession(sessionID string) (*stripe.CheckoutSession, error) {
	session, err := session.Get(sessionID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetCheckoutSession: %v", err)
		return nil, err
	}
	return session, nil
}

// Payment Intent operations
func (c *StripeClient) CreatePaymentIntent(params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreatePaymentIntent: %v", err)
		return nil, err
	}
	return intent, nil
}

func (c *StripeClient) GetPaymentIntent(intentID string) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.Get(intentID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetPaymentIntent: %v", err)
		return nil, err
	}
	return intent, nil
}

func (c *StripeClient) UpdatePaymentIntent(intentID string, params *stripe.PaymentIntentParams) (*stripe.PaymentIntent, error) {
	intent, err := paymentintent.Update(intentID, params)
	if err != nil {
		log.Printf("Stripe API error - UpdatePaymentIntent: %v", err)
		return nil, err
	}
	return intent, nil
}

// Subscription operations (for future use)
func (c *StripeClient) CreateSubscription(params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	subscription, err := subscription.New(params)
	if err != nil {
		log.Printf("Stripe API error - CreateSubscription: %v", err)
		return nil, err
	}
	return subscription, nil
}

func (c *StripeClient) GetSubscription(subscriptionID string) (*stripe.Subscription, error) {
	subscription, err := subscription.Get(subscriptionID, nil)
	if err != nil {
		log.Printf("Stripe API error - GetSubscription: %v", err)
		return nil, err
	}
	return subscription, nil
}

func (c *StripeClient) UpdateSubscription(subscriptionID string, params *stripe.SubscriptionParams) (*stripe.Subscription, error) {
	subscription, err := subscription.Update(subscriptionID, params)
	if err != nil {
		log.Printf("Stripe API error - UpdateSubscription: %v", err)
		return nil, err
	}
	return subscription, nil
}

func (c *StripeClient) CancelSubscription(subscriptionID string, params *stripe.SubscriptionCancelParams) (*stripe.Subscription, error) {
	subscription, err := subscription.Cancel(subscriptionID, params)
	if err != nil {
		log.Printf("Stripe API error - CancelSubscription: %v", err)
		return nil, err
	}
	return subscription, nil
}

func (c *StripeClient) ListSubscriptions(params *stripe.SubscriptionListParams) *subscription.Iter {
	return subscription.List(params)
}
