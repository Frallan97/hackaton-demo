import { api } from '../store/api';

export interface SubscriptionPlan {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
  interval: string;
  features: string[];
}

export interface CreateCheckoutSessionRequest {
  plan_id: string;
  success_url: string;
  cancel_url: string;
}

export interface CreateCheckoutSessionResponse {
  session_id: string;
  url: string;
}

export interface Subscription {
  id: number;
  user_id: number;
  stripe_customer_id: number;
  stripe_sub_id: string;
  status: string;
  plan_id: string;
  plan_name: string;
  current_period_start: string;
  current_period_end: string;
  cancel_at_period_end: boolean;
  created_at: string;
  updated_at: string;
}

export interface Payment {
  id: number;
  user_id: number;
  stripe_customer_id: number;
  stripe_payment_id: string;
  amount: number;
  currency: string;
  status: string;
  description: string;
  created_at: string;
}

export interface SubscriptionMetrics {
  active_subscriptions: number;
  cancelled_subscriptions: number;
  total_revenue_cents: number;
  plan_distribution: Record<string, number>;
}

class StripeService {
  private baseUrl = '/api/stripe';

  // Get available subscription plans
  async getAvailablePlans(): Promise<SubscriptionPlan[]> {
    const response = await fetch(`${this.baseUrl}/plans`);
    if (!response.ok) {
      throw new Error('Failed to fetch plans');
    }
    const data = await response.json();
    return data.data;
  }

  // Create checkout session
  async createCheckoutSession(request: CreateCheckoutSessionRequest): Promise<CreateCheckoutSessionResponse> {
    const response = await fetch(`${this.baseUrl}/checkout`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
      body: JSON.stringify(request),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to create checkout session');
    }

    const data = await response.json();
    return data.data;
  }

  // Get current user subscription
  async getCurrentSubscription(): Promise<Subscription | null> {
    const response = await fetch(`${this.baseUrl}/subscription`, {
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch subscription');
    }

    const data = await response.json();
    return data.data.subscription;
  }

  // Get subscription history
  async getSubscriptionHistory(): Promise<Subscription[]> {
    const response = await fetch(`${this.baseUrl}/subscription/history`, {
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch subscription history');
    }

    const data = await response.json();
    return data.data;
  }

  // Get payment history
  async getPaymentHistory(): Promise<Payment[]> {
    const response = await fetch(`${this.baseUrl}/payments`, {
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch payment history');
    }

    const data = await response.json();
    return data.data;
  }

  // Cancel subscription
  async cancelSubscription(): Promise<void> {
    const response = await fetch(`${this.baseUrl}/subscription/cancel`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to cancel subscription');
    }
  }

  // Reactivate subscription
  async reactivateSubscription(): Promise<void> {
    const response = await fetch(`${this.baseUrl}/subscription/reactivate`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to reactivate subscription');
    }
  }

  // Get subscription metrics (admin only)
  async getSubscriptionMetrics(): Promise<SubscriptionMetrics> {
    const response = await fetch(`${this.baseUrl}/admin/metrics`, {
      headers: {
        'Authorization': `Bearer ${this.getAuthToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Failed to fetch subscription metrics');
    }

    const data = await response.json();
    return data.data;
  }

  // Helper method to get auth token
  private getAuthToken(): string {
    const token = localStorage.getItem('accessToken');
    if (!token) {
      throw new Error('No authentication token found');
    }
    return token;
  }

  // Format price for display
  formatPrice(amount: number, currency: string = 'usd'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency.toUpperCase(),
    }).format(amount / 100); // Convert cents to dollars
  }

  // Check if subscription is active
  isSubscriptionActive(subscription: Subscription | null): boolean {
    return subscription?.status === 'active';
  }

  // Check if subscription is cancelled
  isSubscriptionCancelled(subscription: Subscription | null): boolean {
    return subscription?.cancel_at_period_end === true;
  }

  // Get subscription status text
  getSubscriptionStatusText(subscription: Subscription | null): string {
    if (!subscription) return 'No subscription';
    if (subscription.status === 'active') {
      return subscription.cancel_at_period_end ? 'Active (Cancelling)' : 'Active';
    }
    return subscription.status;
  }
}

export const stripeService = new StripeService(); 