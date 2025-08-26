import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Separator } from './ui/separator';
import { Calendar, CreditCard, AlertCircle, CheckCircle, XCircle } from 'lucide-react';
import { stripeService, Subscription } from '../services/stripeService';

interface SubscriptionManagementProps {
  onSubscriptionUpdate?: () => void;
}

export const SubscriptionManagement: React.FC<SubscriptionManagementProps> = ({
  onSubscriptionUpdate,
}) => {
  const [subscription, setSubscription] = useState<Subscription | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [actionLoading, setActionLoading] = useState(false);

  useEffect(() => {
    fetchSubscription();
  }, []);

  const fetchSubscription = async () => {
    try {
      setLoading(true);
      const currentSubscription = await stripeService.getCurrentSubscription();
      setSubscription(currentSubscription);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch subscription');
    } finally {
      setLoading(false);
    }
  };

  const handleCancelSubscription = async () => {
    if (!subscription) return;

    try {
      setActionLoading(true);
      await stripeService.cancelSubscription();
      await fetchSubscription();
      onSubscriptionUpdate?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to cancel subscription');
    } finally {
      setActionLoading(false);
    }
  };

  const handleReactivateSubscription = async () => {
    if (!subscription) return;

    try {
      setActionLoading(true);
      await stripeService.reactivateSubscription();
      await fetchSubscription();
      onSubscriptionUpdate?.();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to reactivate subscription');
    } finally {
      setActionLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  const getStatusIcon = (status: string, cancelAtPeriodEnd: boolean) => {
    if (status === 'active' && !cancelAtPeriodEnd) {
      return <CheckCircle className="h-5 w-5 text-green-600" />;
    } else if (status === 'active' && cancelAtPeriodEnd) {
      return <AlertCircle className="h-5 w-5 text-yellow-600" />;
    } else {
      return <XCircle className="h-5 w-5 text-red-600" />;
    }
  };

  const getStatusText = (status: string, cancelAtPeriodEnd: boolean) => {
    if (status === 'active' && !cancelAtPeriodEnd) {
      return 'Active';
    } else if (status === 'active' && cancelAtPeriodEnd) {
      return 'Active (Cancelling)';
    } else if (status === 'canceled') {
      return 'Cancelled';
    } else if (status === 'expired') {
      return 'Expired';
    } else {
      return status;
    }
  };

  const getStatusColor = (status: string, cancelAtPeriodEnd: boolean) => {
    if (status === 'active' && !cancelAtPeriodEnd) {
      return 'bg-green-100 text-green-800';
    } else if (status === 'active' && cancelAtPeriodEnd) {
      return 'bg-yellow-100 text-yellow-800';
    } else {
      return 'bg-red-100 text-red-800';
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900"></div>
      </div>
    );
  }

  if (error) {
    return (
      <Card className="max-w-2xl mx-auto">
        <CardContent className="pt-6">
          <div className="text-center">
            <XCircle className="h-12 w-12 text-red-600 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">Error</h3>
            <p className="text-red-600 mb-4">{error}</p>
            <Button onClick={fetchSubscription} variant="outline">
              Try Again
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!subscription) {
    return (
      <Card className="max-w-2xl mx-auto">
        <CardContent className="pt-6">
          <div className="text-center">
            <CreditCard className="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <h3 className="text-lg font-semibold text-gray-900 mb-2">No Active Subscription</h3>
            <p className="text-gray-600 mb-4">
              You don't have an active subscription. Choose a plan to get started.
            </p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="max-w-2xl mx-auto">
      <CardHeader>
        <div className="flex items-center justify-between">
          <div>
            <CardTitle className="text-xl">Subscription Details</CardTitle>
            <CardDescription>
              Manage your {subscription.plan_name} subscription
            </CardDescription>
          </div>
          <Badge className={getStatusColor(subscription.status, subscription.cancel_at_period_end)}>
            <div className="flex items-center gap-2">
              {getStatusIcon(subscription.status, subscription.cancel_at_period_end)}
              {getStatusText(subscription.status, subscription.cancel_at_period_end)}
            </div>
          </Badge>
        </div>
      </CardHeader>

      <CardContent className="space-y-6">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-500">Plan</label>
            <p className="text-gray-900 font-medium">{subscription.plan_name}</p>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-500">Status</label>
            <p className="text-gray-900 font-medium">
              {getStatusText(subscription.status, subscription.cancel_at_period_end)}
            </p>
          </div>
        </div>

        <Separator />

        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-500">Current Period Start</label>
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-gray-400" />
              <p className="text-gray-900">{formatDate(subscription.current_period_start)}</p>
            </div>
          </div>
          <div className="space-y-2">
            <label className="text-sm font-medium text-gray-500">Current Period End</label>
            <div className="flex items-center gap-2">
              <Calendar className="h-4 w-4 text-gray-400" />
              <p className="text-gray-900">{formatDate(subscription.current_period_end)}</p>
            </div>
          </div>
        </div>

        <Separator />

        <div className="space-y-4">
          <h4 className="font-medium text-gray-900">Actions</h4>
          
          {subscription.status === 'active' && !subscription.cancel_at_period_end && (
            <Button
              onClick={handleCancelSubscription}
              variant="outline"
              className="w-full"
              disabled={actionLoading}
            >
              {actionLoading ? 'Cancelling...' : 'Cancel Subscription'}
            </Button>
          )}

          {subscription.status === 'active' && subscription.cancel_at_period_end && (
            <div className="space-y-3">
              <div className="p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                <p className="text-sm text-yellow-800">
                  Your subscription will be cancelled at the end of the current billing period.
                </p>
              </div>
              <Button
                onClick={handleReactivateSubscription}
                className="w-full"
                disabled={actionLoading}
              >
                {actionLoading ? 'Reactivating...' : 'Reactivate Subscription'}
              </Button>
            </div>
          )}

          {subscription.status === 'canceled' && (
            <div className="p-3 bg-gray-50 border border-gray-200 rounded-lg">
              <p className="text-sm text-gray-600">
                Your subscription has been cancelled. You can start a new subscription at any time.
              </p>
            </div>
          )}
        </div>

        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg">
            <p className="text-sm text-red-800">{error}</p>
          </div>
        )}
      </CardContent>
    </Card>
  );
}; 