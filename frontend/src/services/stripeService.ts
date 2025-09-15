import { api } from '../store/api';
import config from '../config';

export interface PaymentPlan {
  id: string;
  name: string;
  description: string;
  price: number;
  currency: string;
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

export interface PaymentMetrics {
  total_payments: number;
  total_revenue_cents: number;
  plan_distribution: Record<string, number>;
}

class StripeService {
  private baseUrl = `${config.apiBaseUrl}/api/stripe`;

  // Get available payment plans
  async getAvailablePlans(): Promise<PaymentPlan[]> {
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
    return data.data || []; // Return empty array if data.data is null
  }



  // Helper method to get auth token
  private getAuthToken(): string {
    // Try both token names to handle different storage patterns
    const token = localStorage.getItem('access_token') || localStorage.getItem('accessToken');
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


}

export const stripeService = new StripeService(); 