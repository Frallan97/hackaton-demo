# Frontend Stripe Integration

This directory contains the frontend components and services for the Stripe integration demo.

## Components Overview

### 🎯 **StripeDemo** - Main Demo Component
The primary component displayed on the home page for logged-in users. Features:
- **Subscription Plans Tab**: Display available plans with pricing and features
- **Manage Subscription Tab**: View and manage current subscription
- **How It Works Tab**: Technical details and feature explanations
- **Real-time Status**: Shows current subscription status and plan

### 📋 **SubscriptionPlans** - Plan Selection
Displays available subscription plans with:
- Visual plan cards with icons and colors
- Feature lists for each plan
- Pricing display (automatically converts cents to dollars)
- Current plan highlighting
- Plan selection buttons

### ⚙️ **SubscriptionManagement** - Subscription Control
Allows users to:
- View current subscription details
- See billing period information
- Cancel subscriptions (at period end)
- Reactivate cancelled subscriptions
- View subscription status and history

### 💳 **PaymentHistory** - Payment Tracking
Shows payment history with:
- Payment amounts and status
- Transaction dates and descriptions
- Status badges (success, failed, pending)
- Admin view with totals and summaries

### 📊 **SubscriptionMetrics** - Admin Analytics
Provides administrators with:
- Active vs. cancelled subscription counts
- Total revenue tracking
- Plan distribution charts
- Subscription health metrics
- Revenue insights and averages

## Services

### **stripeService** - API Integration
Handles all communication with the backend Stripe endpoints:
- Plan fetching
- Checkout session creation
- Subscription management
- Payment history retrieval
- Admin metrics access

## Configuration

### Environment Variables
Create a `.env` file in the frontend directory:

```bash
# Stripe Configuration
VITE_STRIPE_PUBLISHABLE_KEY=pk_test_your_stripe_publishable_key_here

# Backend API URL
VITE_BACKEND_URL=http://localhost:8080
```

### Stripe Setup
1. Get your Stripe publishable key from the Stripe dashboard
2. Add it to the environment variables
3. Ensure the backend is running with Stripe configuration

## Usage Examples

### Basic Integration
```tsx
import { StripeDemo } from './components/StripeDemo';

function App() {
  return (
    <div>
      <h1>My App</h1>
      <StripeDemo />
    </div>
  );
}
```

### Standalone Components
```tsx
import { SubscriptionPlans } from './components/SubscriptionPlans';
import { SubscriptionManagement } from './components/SubscriptionManagement';

function SubscriptionPage() {
  const handlePlanSelect = (plan) => {
    console.log('Selected plan:', plan);
  };

  return (
    <div>
      <SubscriptionPlans onPlanSelect={handlePlanSelect} />
      <SubscriptionManagement />
    </div>
  );
}
```

### Admin Dashboard
```tsx
import { SubscriptionMetrics } from './components/SubscriptionMetrics';
import { PaymentHistory } from './components/PaymentHistory';

function AdminDashboard() {
  return (
    <div>
      <SubscriptionMetrics />
      <PaymentHistory isAdmin={true} />
    </div>
  );
}
```

## Features

### ✅ **Complete Subscription Lifecycle**
- Plan selection and checkout
- Subscription management
- Cancellation and reactivation
- Payment tracking

### 🔒 **Security & Authentication**
- JWT token integration
- Protected API endpoints
- User-specific data access
- Admin role verification

### 📱 **Responsive Design**
- Mobile-first approach
- Tailwind CSS styling
- Consistent UI components
- Accessible design patterns

### 🚀 **Real-time Updates**
- Automatic status refresh
- Live subscription updates
- Payment confirmation
- Error handling and recovery

## Component Props

### SubscriptionPlans
```tsx
interface SubscriptionPlansProps {
  onPlanSelect?: (plan: SubscriptionPlan) => void;
  showCurrentPlan?: boolean;
  currentPlanId?: string;
}
```

### SubscriptionManagement
```tsx
interface SubscriptionManagementProps {
  onSubscriptionUpdate?: () => void;
}
```

### PaymentHistory
```tsx
interface PaymentHistoryProps {
  userId?: number;
  isAdmin?: boolean;
}
```

## Styling

All components use Tailwind CSS classes and follow the existing design system:
- Consistent color schemes
- Responsive breakpoints
- Hover and focus states
- Loading and error states
- Icon integration with Lucide React

## Error Handling

Components include comprehensive error handling:
- API failure recovery
- Network error messages
- User-friendly error displays
- Retry mechanisms
- Fallback states

## Testing

To test the integration:

1. **Set up Stripe test keys**
2. **Create test products/prices** in Stripe dashboard
3. **Use test card numbers** for payments
4. **Monitor webhook events** in Stripe dashboard
5. **Check backend logs** for webhook processing

## Dependencies

- `@stripe/stripe-js`: Stripe JavaScript SDK
- `@stripe/react-stripe-js`: React components for Stripe
- `lucide-react`: Icon library
- `tailwindcss`: CSS framework
- Custom UI components from `./ui/`

## Browser Support

- Modern browsers (Chrome, Firefox, Safari, Edge)
- ES6+ features
- CSS Grid and Flexbox
- Local Storage for authentication

## Performance

- Lazy loading of components
- Efficient state management
- Minimal re-renders
- Optimized API calls
- Responsive image handling

## Accessibility

- ARIA labels and descriptions
- Keyboard navigation support
- Screen reader compatibility
- High contrast support
- Focus management

## Future Enhancements

- **Real-time updates** with WebSocket
- **Advanced analytics** and reporting
- **Multi-currency support**
- **Subscription upgrades/downgrades**
- **Bulk operations** for admins
- **Export functionality** for data
- **Advanced filtering** and search
- **Mobile app integration** 