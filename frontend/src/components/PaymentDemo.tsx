import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Separator } from './ui/separator';
import { CreditCard, Crown, Zap, CheckCircle, AlertCircle } from 'lucide-react';
import { 
  PaymentList, 
  PaymentHistory, 
  PaymentStatus, 
  StripeError, 
  StripeLoading,
  useStripe 
} from './stripe';
import { PaymentPlan } from '../services/stripeService';
import { stripeConfig } from '../config/stripe';

export const PaymentDemo: React.FC = () => {
  const {
    plans,
    payments,
    loading,
    error,
    loadPlans,
    loadPayments,
    createCheckout,
    clearError,
    formatPrice,
  } = useStripe();

  const [checkoutLoading, setCheckoutLoading] = useState(false);

  useEffect(() => {
    loadPlans();
  }, [loadPlans]);

  const handlePlanSelect = async (plan: PaymentPlan) => {
    try {
      setCheckoutLoading(true);
      
      // Create checkout session
      const checkoutUrl = await createCheckout({
        plan_id: plan.id,
        success_url: `${window.location.origin}/?payment=success`,
        cancel_url: `${window.location.origin}/?payment=cancelled`,
      });

      // Redirect to Stripe Checkout
      window.location.href = checkoutUrl;
    } catch (err) {
      // Error is already handled by the useStripe hook
      setCheckoutLoading(false);
    }
  };

  if (!stripeConfig.isEnabled) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center">
            <CreditCard className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold mb-2">Stripe Not Configured</h3>
            <p className="text-gray-600">
              Please configure your Stripe keys to enable payment functionality.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Status Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="p-2 bg-blue-100 rounded-lg">
                <CreditCard className="h-6 w-6 text-blue-600" />
              </div>
              <div>
                <CardTitle className="text-xl">Payment Integration</CardTitle>
                <CardDescription>
                  Secure one-time payments powered by Stripe
                </CardDescription>
              </div>
            </div>
            <div className="flex items-center space-x-2">
              <CheckCircle className="h-5 w-5 text-green-600" />
              <span className="text-sm font-medium text-green-600">Active</span>
            </div>
          </div>
        </CardHeader>
      </Card>

      {/* Error Display */}
      {error && (
        <StripeError 
          error={error}
          onDismiss={clearError}
          title="Payment Error"
          variant="banner"
        />
      )}

      {/* Main Content Tabs */}
      <Tabs defaultValue="plans" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="plans" className="flex items-center gap-2">
            <Crown className="h-4 w-4" />
            Payment Plans
          </TabsTrigger>
          <TabsTrigger value="history" className="flex items-center gap-2">
            <CreditCard className="h-4 w-4" />
            Payment History
          </TabsTrigger>
          <TabsTrigger value="info" className="flex items-center gap-2">
            <Zap className="h-4 w-4" />
            How It Works
          </TabsTrigger>
        </TabsList>

        <TabsContent value="plans" className="space-y-6">
          <div className="text-center">
            <h3 className="text-2xl font-semibold mb-2">Test Payment Integration</h3>
            <p className="text-gray-600 text-lg">
              Try our payment system with cards. Swish support ready when configured.
            </p>
          </div>
          
          {loading ? (
            <StripeLoading message="Loading payment plans..." />
          ) : (
            <PaymentList 
              plans={plans}
              onPlanSelect={handlePlanSelect}
              loading={checkoutLoading}
              layout="single"
              formatPrice={formatPrice}
            />
          )}

          {checkoutLoading && (
            <div className="text-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto mb-2"></div>
              <p className="text-gray-600">Preparing checkout...</p>
            </div>
          )}
        </TabsContent>

        <TabsContent value="history" className="space-y-6">
          <div className="text-center">
            <h3 className="text-xl font-semibold mb-2">Payment History</h3>
            <p className="text-gray-600">
              View all your past payments and transactions.
            </p>
          </div>
          
          <PaymentHistory />
        </TabsContent>

        <TabsContent value="info" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <CheckCircle className="h-5 w-5 text-green-600" />
                  <span>Secure Payments</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600">
                  All payments are processed securely through Stripe, ensuring your 
                  payment information is protected with industry-standard encryption.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Zap className="h-5 w-5 text-blue-600" />
                  <span>One-Time Payment</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600">
                  Pay once and get lifetime access to your selected plan. 
                  No recurring charges or hidden fees.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <CreditCard className="h-5 w-5 text-purple-600" />
                  <span>Card & Swish Payments</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600">
                  Support for credit cards and debit cards. Swish payments can be 
                  enabled with additional Stripe configuration for Swedish customers.
                </p>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center space-x-2">
                  <Crown className="h-5 w-5 text-yellow-600" />
                  <span>Instant Access</span>
                </CardTitle>
              </CardHeader>
              <CardContent>
                <p className="text-gray-600">
                  Get immediate access to your purchased features as soon as 
                  your payment is confirmed.
                </p>
              </CardContent>
            </Card>
          </div>

          <Separator />

          <div className="text-center">
            <h4 className="text-lg font-semibold mb-4">Technical Implementation</h4>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
              <div className="p-4 bg-gray-50 rounded-lg">
                <h5 className="font-medium mb-2">Frontend</h5>
                <p className="text-gray-600">React with TypeScript, Tailwind CSS, and Redux Toolkit</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg">
                <h5 className="font-medium mb-2">Backend</h5>
                <p className="text-gray-600">Go with PostgreSQL database and JWT authentication</p>
              </div>
              <div className="p-4 bg-gray-50 rounded-lg">
                <h5 className="font-medium mb-2">Payments</h5>
                <p className="text-gray-600">Stripe Checkout with webhook integration</p>
              </div>
            </div>
          </div>
        </TabsContent>
      </Tabs>
    </div>
  );
}; 