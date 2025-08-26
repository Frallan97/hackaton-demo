import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Tabs, TabsContent, TabsList, TabsTrigger } from './ui/tabs';
import { Badge } from './ui/badge';
import { Separator } from './ui/separator';
import { CreditCard, Crown, Zap, Building, CheckCircle, AlertCircle } from 'lucide-react';
import { SubscriptionPlans } from './SubscriptionPlans';
import { SubscriptionManagement } from './SubscriptionManagement';
import { stripeService, SubscriptionPlan, Subscription } from '../services/stripeService';
import { stripeConfig } from '../config/stripe';

export const StripeDemo: React.FC = () => {
  const [currentSubscription, setCurrentSubscription] = useState<Subscription | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [checkoutLoading, setCheckoutLoading] = useState(false);

  useEffect(() => {
    fetchCurrentSubscription();
  }, []);

  const fetchCurrentSubscription = async () => {
    try {
      setLoading(true);
      const subscription = await stripeService.getCurrentSubscription();
      setCurrentSubscription(subscription);
    } catch (err) {
      // Don't show error for subscription fetch, just log it
      console.log('No subscription found or error fetching:', err);
    } finally {
      setLoading(false);
    }
  };

  const handlePlanSelect = async (plan: SubscriptionPlan) => {
    try {
      setCheckoutLoading(true);
      
      // Create checkout session
      const session = await stripeService.createCheckoutSession({
        plan_id: plan.id,
        success_url: `${window.location.origin}/?subscription=success`,
        cancel_url: `${window.location.origin}/?subscription=cancelled`,
      });

      // Redirect to Stripe Checkout
      window.location.href = session.url;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create checkout session');
      setCheckoutLoading(false);
    }
  };

  const handleSubscriptionUpdate = () => {
    fetchCurrentSubscription();
  };

  const getSubscriptionStatusIcon = (subscription: Subscription | null) => {
    if (!subscription) return <CreditCard className="h-5 w-5 text-gray-400" />;
    
    if (subscription.status === 'active' && !subscription.cancel_at_period_end) {
      return <CheckCircle className="h-5 w-5 text-green-600" />;
    } else if (subscription.status === 'active' && subscription.cancel_at_period_end) {
      return <AlertCircle className="h-5 w-5 text-yellow-600" />;
    } else {
      return <CreditCard className="h-5 w-5 text-gray-400" />;
    }
  };

  const getSubscriptionStatusText = (subscription: Subscription | null) => {
    if (!subscription) return 'No Subscription';
    
    if (subscription.status === 'active' && !subscription.cancel_at_period_end) {
      return 'Active';
    } else if (subscription.status === 'active' && subscription.cancel_at_period_end) {
      return 'Active (Cancelling)';
    } else if (subscription.status === 'canceled') {
      return 'Cancelled';
    } else {
      return subscription.status;
    }
  };

  const getSubscriptionStatusColor = (subscription: Subscription | null) => {
    if (!subscription) return 'bg-gray-100 text-gray-800';
    
    if (subscription.status === 'active' && !subscription.cancel_at_period_end) {
      return 'bg-green-100 text-green-800';
    } else if (subscription.status === 'active' && subscription.cancel_at_period_end) {
      return 'bg-yellow-100 text-yellow-800';
    } else {
      return 'bg-gray-100 text-gray-800';
    }
  };

  if (!stripeConfig.isEnabled) {
    return (
      <Card className="max-w-4xl mx-auto">
        <CardHeader className="text-center">
          <CreditCard className="h-12 w-12 text-gray-400 mx-auto mb-4" />
          <CardTitle className="text-2xl">Stripe Integration Demo</CardTitle>
          <CardDescription>
            Stripe integration is not configured. Please set up your Stripe keys to enable this feature.
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="text-center text-gray-600">
            <p className="mb-4">
              To enable Stripe integration, add the following environment variables:
            </p>
            <div className="bg-gray-100 p-4 rounded-lg text-left text-sm font-mono">
              <p>VITE_STRIPE_PUBLISHABLE_KEY=pk_test_your_key_here</p>
            </div>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="max-w-6xl mx-auto space-y-8">
      {/* Header */}
      <Card>
        <CardHeader className="text-center">
          <div className="flex justify-center mb-4">
            <div className="p-3 rounded-full bg-blue-100">
              <CreditCard className="h-8 w-8 text-blue-600" />
            </div>
          </div>
          <CardTitle className="text-3xl">Stripe Integration Demo</CardTitle>
          <CardDescription className="text-lg">
            Experience our subscription management system with Stripe integration
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
            <div className="flex items-center gap-2">
              {getSubscriptionStatusIcon(currentSubscription)}
              <span className="text-sm text-gray-600">Status:</span>
              <Badge className={getSubscriptionStatusColor(currentSubscription)}>
                {getSubscriptionStatusText(currentSubscription)}
              </Badge>
            </div>
            
            {currentSubscription && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-gray-600">Plan:</span>
                <Badge variant="outline" className="font-medium">
                  {currentSubscription.plan_name}
                </Badge>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Main Content Tabs */}
      <Tabs defaultValue="plans" className="w-full">
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="plans" className="flex items-center gap-2">
            <Crown className="h-4 w-4" />
            Subscription Plans
          </TabsTrigger>
          <TabsTrigger value="manage" className="flex items-center gap-2">
            <CreditCard className="h-4 w-4" />
            Manage Subscription
          </TabsTrigger>
          <TabsTrigger value="info" className="flex items-center gap-2">
            <Zap className="h-4 w-4" />
            How It Works
          </TabsTrigger>
        </TabsList>

        <TabsContent value="plans" className="space-y-6">
          <div className="text-center">
            <h3 className="text-xl font-semibold mb-2">Choose Your Plan</h3>
            <p className="text-gray-600">
              Select a subscription plan that fits your needs. All plans include our core features.
            </p>
          </div>
          
          <SubscriptionPlans 
            onPlanSelect={handlePlanSelect}
            showCurrentPlan={true}
            currentPlanId={currentSubscription?.plan_id}
          />

          {checkoutLoading && (
            <div className="text-center py-4">
              <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600 mx-auto mb-2"></div>
              <p className="text-gray-600">Preparing checkout...</p>
            </div>
          )}

          {error && (
            <Card className="border-red-200 bg-red-50">
              <CardContent className="pt-6">
                <div className="text-center text-red-800">
                  <p className="font-medium">Error</p>
                  <p className="text-sm">{error}</p>
                  <Button 
                    onClick={() => setError(null)} 
                    variant="outline" 
                    size="sm" 
                    className="mt-2"
                  >
                    Dismiss
                  </Button>
                </div>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="manage" className="space-y-6">
          <div className="text-center">
            <h3 className="text-xl font-semibold mb-2">Manage Your Subscription</h3>
            <p className="text-gray-600">
              View details, cancel, or reactivate your current subscription.
            </p>
          </div>
          
          <SubscriptionManagement onSubscriptionUpdate={handleSubscriptionUpdate} />
        </TabsContent>

        <TabsContent value="info" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Zap className="h-5 w-5 text-blue-600" />
                  How It Works
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-start gap-3">
                  <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-medium flex items-center justify-center flex-shrink-0 mt-0.5">
                    1
                  </div>
                  <p className="text-sm text-gray-600">
                    Choose a subscription plan that fits your needs
                  </p>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-medium flex items-center justify-center flex-shrink-0 mt-0.5">
                    2
                  </div>
                  <p className="text-sm text-gray-600">
                    Complete payment securely through Stripe Checkout
                  </p>
                </div>
                <div className="flex items-start gap-3">
                  <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 text-sm font-medium flex items-center justify-center flex-shrink-0 mt-0.5">
                    3
                  </div>
                  <p className="text-sm text-gray-600">
                    Access your subscription features immediately
                  </p>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Crown className="h-5 w-5 text-purple-600" />
                  Features
                </CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-gray-600">Secure payment processing</span>
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-gray-600">Automatic subscription management</span>
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-gray-600">Webhook automation</span>
                </div>
                <div className="flex items-center gap-2">
                  <CheckCircle className="h-4 w-4 text-green-600" />
                  <span className="text-sm text-gray-600">Real-time status updates</span>
                </div>
              </CardContent>
            </Card>
          </div>

          <Card>
            <CardHeader>
              <CardTitle>Technical Details</CardTitle>
              <CardDescription>
                This demo showcases a complete Stripe integration built with modern web technologies
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
                <div>
                  <h4 className="font-medium mb-2">Backend</h4>
                  <ul className="space-y-1 text-gray-600">
                    <li>• Go with clean architecture</li>
                    <li>• PostgreSQL database</li>
                    <li>• NATS event system</li>
                    <li>• JWT authentication</li>
                  </ul>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Frontend</h4>
                  <ul className="space-y-1 text-gray-600">
                    <li>• React with TypeScript</li>
                    <li>• Tailwind CSS styling</li>
                    <li>• Stripe React components</li>
                    <li>• Responsive design</li>
                  </ul>
                </div>
                <div>
                  <h4 className="font-medium mb-2">Integration</h4>
                  <ul className="space-y-1 text-gray-600">
                    <li>• Stripe Checkout</li>
                    <li>• Webhook handling</li>
                    <li>• Subscription lifecycle</li>
                    <li>• Payment tracking</li>
                  </ul>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}; 