import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Check, Crown, Zap, Building } from 'lucide-react';
import { stripeService, SubscriptionPlan } from '../services/stripeService';
import { stripeConfig } from '../config/stripe';

interface SubscriptionPlansProps {
  onPlanSelect?: (plan: SubscriptionPlan) => void;
  showCurrentPlan?: boolean;
  currentPlanId?: string;
}

const planIcons = {
  'Basic Plan': <Zap className="h-6 w-6" />,
  'Pro Plan': <Crown className="h-6 w-6" />,
  'Enterprise Plan': <Building className="h-6 w-6" />,
};

const planColors = {
  'Basic Plan': 'bg-blue-50 border-blue-200',
  'Pro Plan': 'bg-purple-50 border-purple-200',
  'Enterprise Plan': 'bg-green-50 border-green-200',
};

export const SubscriptionPlans: React.FC<SubscriptionPlansProps> = ({
  onPlanSelect,
  showCurrentPlan = false,
  currentPlanId,
}) => {
  const [plans, setPlans] = useState<SubscriptionPlan[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPlans = async () => {
      try {
        setLoading(true);
        const availablePlans = await stripeService.getAvailablePlans();
        setPlans(availablePlans);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Failed to fetch plans');
      } finally {
        setLoading(false);
      }
    };

    if (stripeConfig.isEnabled) {
      fetchPlans();
    } else {
      setLoading(false);
      setError('Stripe is not configured');
    }
  }, []);

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-center py-8">
        <p className="text-red-600 mb-4">{error}</p>
        <p className="text-gray-600 text-sm">
          Please check your Stripe configuration or try again later.
        </p>
      </div>
    );
  }

  if (!stripeConfig.isEnabled) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-600 mb-4">Stripe integration is not available</p>
        <p className="text-gray-500 text-sm">
          Please configure Stripe to view subscription plans.
        </p>
      </div>
    );
  }

  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-6 max-w-6xl mx-auto">
      {plans.map((plan) => {
        const isCurrentPlan = currentPlanId === plan.id;
        const icon = planIcons[plan.name as keyof typeof planIcons] || <Zap className="h-6 w-6" />;
        const colorClass = planColors[plan.name as keyof typeof planColors] || 'bg-gray-50 border-gray-200';

        return (
          <Card
            key={plan.id}
            className={`relative transition-all duration-200 hover:shadow-lg ${
              isCurrentPlan ? 'ring-2 ring-blue-500' : ''
            } ${colorClass}`}
          >
            {isCurrentPlan && showCurrentPlan && (
              <Badge className="absolute -top-3 left-1/2 transform -translate-x-1/2 bg-blue-600">
                Current Plan
              </Badge>
            )}
            
            <CardHeader className="text-center pb-4">
              <div className="flex justify-center mb-3">
                <div className="p-3 rounded-full bg-white shadow-sm">
                  {icon}
                </div>
              </div>
              <CardTitle className="text-xl font-bold">{plan.name}</CardTitle>
              <CardDescription className="text-gray-600">
                {plan.description}
              </CardDescription>
            </CardHeader>

            <CardContent className="text-center">
              <div className="mb-6">
                <div className="text-3xl font-bold text-gray-900">
                  {stripeService.formatPrice(plan.price, plan.currency)}
                </div>
                <div className="text-gray-600">per {plan.interval}</div>
              </div>

              <div className="space-y-3 mb-6 text-left">
                {plan.features.map((feature, index) => (
                  <div key={index} className="flex items-center">
                    <Check className="h-4 w-4 text-green-600 mr-3 flex-shrink-0" />
                    <span className="text-gray-700">{feature}</span>
                  </div>
                ))}
              </div>

              {onPlanSelect && (
                <Button
                  onClick={() => onPlanSelect(plan)}
                  className={`w-full ${
                    isCurrentPlan
                      ? 'bg-gray-400 cursor-not-allowed'
                      : 'bg-blue-600 hover:bg-blue-700'
                  }`}
                  disabled={isCurrentPlan}
                >
                  {isCurrentPlan ? 'Current Plan' : 'Select Plan'}
                </Button>
              )}
            </CardContent>
          </Card>
        );
      })}
    </div>
  );
}; 