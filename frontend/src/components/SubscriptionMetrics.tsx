import React, { useState, useEffect } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from './ui/card';
import { Badge } from './ui/badge';
import { Button } from './ui/button';
import { TrendingUp, Users, DollarSign, BarChart3, RefreshCw } from 'lucide-react';
import { stripeService, SubscriptionMetrics as Metrics } from '../services/stripeService';

export const SubscriptionMetrics: React.FC = () => {
  const [metrics, setMetrics] = useState<Metrics | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchMetrics();
  }, []);

  const fetchMetrics = async () => {
    try {
      setLoading(true);
      const subscriptionMetrics = await stripeService.getSubscriptionMetrics();
      setMetrics(subscriptionMetrics);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch subscription metrics');
    } finally {
      setLoading(false);
    }
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount / 100); // Convert cents to dollars
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
      <Card>
        <CardContent className="pt-6">
          <div className="text-center">
            <p className="text-red-600 mb-4">{error}</p>
            <Button onClick={fetchMetrics} variant="outline">
              Try Again
            </Button>
          </div>
        </CardContent>
      </Card>
    );
  }

  if (!metrics) {
    return (
      <Card>
        <CardContent className="pt-6">
          <div className="text-center">
            <p className="text-gray-600">No metrics available</p>
          </div>
        </CardContent>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold text-gray-900">Subscription Analytics</h2>
          <p className="text-gray-600">Overview of subscription performance and revenue</p>
        </div>
        <Button onClick={fetchMetrics} variant="outline" size="sm">
          <RefreshCw className="h-4 w-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* Key Metrics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Active Subscriptions</CardTitle>
            <Users className="h-4 w-4 text-blue-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-blue-600">
              {metrics.active_subscriptions}
            </div>
            <p className="text-xs text-gray-600">
              Currently active subscriptions
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Cancelled Subscriptions</CardTitle>
            <TrendingUp className="h-4 w-4 text-red-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-red-600">
              {metrics.cancelled_subscriptions}
            </div>
            <p className="text-xs text-gray-600">
              Total cancelled subscriptions
            </p>
          </CardContent>
        </Card>

        <Card>
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-sm font-medium">Total Revenue</CardTitle>
            <DollarSign className="h-4 w-4 text-green-600" />
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold text-green-600">
              {formatCurrency(metrics.total_revenue_cents)}
            </div>
            <p className="text-xs text-gray-600">
              Lifetime revenue from subscriptions
            </p>
          </CardContent>
        </Card>
      </div>

      {/* Plan Distribution */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5 text-purple-600" />
            Plan Distribution
          </CardTitle>
          <CardDescription>
            Breakdown of active subscriptions by plan type
          </CardDescription>
        </CardHeader>
        <CardContent>
          {Object.keys(metrics.plan_distribution).length > 0 ? (
            <div className="space-y-4">
              {Object.entries(metrics.plan_distribution).map(([planName, count]) => (
                <div key={planName} className="flex items-center justify-between">
                  <div className="flex items-center gap-3">
                    <div className="w-4 h-4 rounded-full bg-purple-500"></div>
                    <span className="font-medium text-gray-900">{planName}</span>
                  </div>
                  <div className="flex items-center gap-3">
                    <div className="w-32 bg-gray-200 rounded-full h-2">
                      <div
                        className="bg-purple-500 h-2 rounded-full"
                        style={{
                          width: `${(count / metrics.active_subscriptions) * 100}%`
                        }}
                      ></div>
                    </div>
                    <Badge variant="secondary">
                      {count} ({((count / metrics.active_subscriptions) * 100).toFixed(1)}%)
                    </Badge>
                  </div>
                </div>
              ))}
            </div>
          ) : (
            <div className="text-center py-8 text-gray-500">
              No plan distribution data available
            </div>
          )}
        </CardContent>
      </Card>

      {/* Additional Insights */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Subscription Health</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Active Rate</span>
              <Badge variant="outline" className="font-medium">
                {metrics.active_subscriptions > 0 
                  ? ((metrics.active_subscriptions / (metrics.active_subscriptions + metrics.cancelled_subscriptions)) * 100).toFixed(1)
                  : 0}%
              </Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Cancellation Rate</span>
              <Badge variant="outline" className="font-medium">
                {metrics.cancelled_subscriptions > 0
                  ? ((metrics.cancelled_subscriptions / (metrics.active_subscriptions + metrics.cancelled_subscriptions)) * 100).toFixed(1)
                  : 0}%
              </Badge>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Revenue Insights</CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Avg Revenue per Sub</span>
              <Badge variant="outline" className="font-medium">
                {metrics.active_subscriptions > 0
                  ? formatCurrency(metrics.total_revenue_cents / metrics.active_subscriptions)
                  : '$0.00'}
              </Badge>
            </div>
            <div className="flex items-center justify-between">
              <span className="text-sm text-gray-600">Total Subscriptions</span>
              <Badge variant="outline" className="font-medium">
                {metrics.active_subscriptions + metrics.cancelled_subscriptions}
              </Badge>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}; 