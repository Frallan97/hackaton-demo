import React, { useEffect } from 'react';
import { 
  PaymentList, 
  PaymentStatus, 
  StripeError, 
  StripeLoading,
  useStripe 
} from '../stripe';
import { PaymentPlan } from '../../services/stripeService';

/**
 * Example: Simple Payment Page
 * 
 * This demonstrates how to use the modular Stripe components
 * on any page with minimal setup.
 */
export const SimplePaymentPage: React.FC = () => {
  const {
    plans,
    loading,
    error,
    loadPlans,
    createCheckout,
    clearError,
    formatPrice,
  } = useStripe();

  useEffect(() => {
    loadPlans();
  }, [loadPlans]);

  const handlePlanSelect = async (plan: PaymentPlan) => {
    try {
      const checkoutUrl = await createCheckout({
        plan_id: plan.id,
        success_url: `${window.location.origin}/success`,
        cancel_url: `${window.location.origin}/cancel`,
      });
      window.location.href = checkoutUrl;
    } catch (err) {
      console.error('Checkout failed:', err);
    }
  };

  if (loading) {
    return <StripeLoading message="Loading payment options..." />;
  }

  if (error) {
    return (
      <StripeError 
        error={error} 
        onRetry={loadPlans}
        onDismiss={clearError}
      />
    );
  }

  return (
    <div className="max-w-2xl mx-auto p-6 space-y-6">
      <div className="text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4">
          Choose Your Plan
        </h1>
        <p className="text-gray-600">
          Simple, one-time payment. No subscriptions.
        </p>
      </div>

      {/* Payment Status */}
      <PaymentStatus 
        status="none"
        formatPrice={formatPrice}
      />

      {/* Payment Plans */}
      <PaymentList
        plans={plans}
        onPlanSelect={handlePlanSelect}
        layout="single"
        formatPrice={formatPrice}
      />
    </div>
  );
}; 