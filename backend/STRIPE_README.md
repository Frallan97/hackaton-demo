# Stripe Integration

This module provides a complete Stripe integration for subscription management, payments, and webhook handling.

## Features

- **Customer Management**: Create and manage Stripe customers linked to users
- **Subscription Management**: Handle subscription creation, updates, and cancellation
- **Payment Tracking**: Record and track all payment attempts and results
- **Webhook Handling**: Process Stripe webhook events automatically
- **Access Control**: Middleware for subscription-based feature gating
- **Admin Metrics**: Subscription and revenue analytics for administrators

## Configuration

Add the following environment variables to your `.env` file:

```bash
# Stripe Configuration
STRIPE_SECRET_KEY=sk_test_your_stripe_secret_key_here
STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key_here
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret_here
STRIPE_ENDPOINT_SECRET=whsec_your_endpoint_secret_here
```

## API Endpoints

### Public Endpoints
- `GET /api/stripe/plans` - Get available subscription plans
- `POST /api/stripe/webhook` - Stripe webhook endpoint

### Protected Endpoints (Require Authentication)
- `POST /api/stripe/checkout` - Create checkout session
- `GET /api/stripe/subscription` - Get current subscription
- `GET /api/stripe/subscription/history` - Get subscription history
- `GET /api/stripe/payments` - Get payment history
- `POST /api/stripe/subscription/cancel` - Cancel subscription
- `POST /api/stripe/subscription/reactivate` - Reactivate subscription

### Admin Endpoints (Require Admin Role)
- `GET /api/stripe/admin/metrics` - Get subscription metrics

## Usage Examples

### Creating a Checkout Session

```typescript
const response = await fetch('/api/stripe/checkout', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    plan_id: 'price_basic_monthly',
    success_url: 'https://yourapp.com/success',
    cancel_url: 'https://yourapp.com/cancel'
  })
});

const { url } = await response.json();
window.location.href = url; // Redirect to Stripe Checkout
```

### Checking Subscription Status

```typescript
const response = await fetch('/api/stripe/subscription', {
  headers: {
    'Authorization': `Bearer ${token}`
  }
});

const subscription = await response.json();
if (subscription.data && subscription.data.status === 'active') {
  // User has active subscription
}
```

## Middleware Usage

### Require Active Subscription

```go
// Apply to routes that require any active subscription
mux.Handle("/api/premium-feature", 
    subscriptionMiddleware.RequireSubscription()(
        http.HandlerFunc(handler)
    ))
```

### Require Specific Plan

```go
// Apply to routes that require a specific plan level
mux.Handle("/api/enterprise-feature", 
    subscriptionMiddleware.RequirePlan("enterprise")(
        http.HandlerFunc(handler)
    ))
```

### Add Subscription Context

```go
// Add subscription info to request context
mux.Handle("/api/feature", 
    subscriptionMiddleware.AddSubscriptionContext()(
        http.HandlerFunc(handler)
    ))
```

## Webhook Events Handled

- `checkout.session.completed` - Checkout completed
- `customer.subscription.created` - New subscription
- `customer.subscription.updated` - Subscription updated
- `customer.subscription.deleted` - Subscription cancelled
- `invoice.payment_succeeded` - Payment successful
- `invoice.payment_failed` - Payment failed

## Database Schema

The integration creates the following tables:

- `stripe_customers` - Links users to Stripe customers
- `subscriptions` - Tracks user subscriptions
- `payments` - Records payment history
- Updates `users` table with subscription fields

## Testing

1. Set up Stripe test keys in your environment
2. Use Stripe's test card numbers for testing payments
3. Use Stripe CLI to test webhooks locally:
   ```bash
   stripe listen --forward-to localhost:8080/api/stripe/webhook
   ```

## Security Notes

- Webhook signatures are verified using `Stripe-Signature` header
- All endpoints (except webhook) require authentication
- Admin endpoints require admin role
- Customer data is linked to authenticated users only

## Error Handling

The integration includes comprehensive error handling:
- Database connection errors
- Stripe API errors
- Webhook signature verification failures
- Invalid subscription states
- Payment processing failures

## Monitoring

- Logs all webhook events
- Tracks subscription lifecycle changes
- Records payment success/failure
- Provides admin metrics dashboard 