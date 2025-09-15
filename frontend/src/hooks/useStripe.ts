import { useState, useCallback } from 'react';
import { stripeService, PaymentPlan, Payment, CreateCheckoutSessionRequest } from '../services/stripeService';

export interface UseStripeReturn {
  // State
  plans: PaymentPlan[];
  payments: Payment[];
  loading: boolean;
  error: string | null;
  
  // Actions
  loadPlans: () => Promise<void>;
  loadPayments: () => Promise<void>;
  createCheckout: (request: CreateCheckoutSessionRequest) => Promise<string>;
  clearError: () => void;
  
  // Utilities
  formatPrice: (amount: number, currency?: string) => string;
}

export const useStripe = (): UseStripeReturn => {
  const [plans, setPlans] = useState<PaymentPlan[]>([]);
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadPlans = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const availablePlans = await stripeService.getAvailablePlans();
      setPlans(availablePlans);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load payment plans');
    } finally {
      setLoading(false);
    }
  }, []);

  const loadPayments = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      const paymentHistory = await stripeService.getPaymentHistory();
      setPayments(paymentHistory);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load payment history');
    } finally {
      setLoading(false);
    }
  }, []);

  const createCheckout = useCallback(async (request: CreateCheckoutSessionRequest): Promise<string> => {
    try {
      setLoading(true);
      setError(null);
      const session = await stripeService.createCheckoutSession(request);
      return session.url;
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Failed to create checkout session';
      setError(errorMessage);
      throw new Error(errorMessage);
    } finally {
      setLoading(false);
    }
  }, []);

  const clearError = useCallback(() => {
    setError(null);
  }, []);

  const formatPrice = useCallback((amount: number, currency = 'usd') => {
    return stripeService.formatPrice(amount, currency);
  }, []);

  return {
    plans,
    payments,
    loading,
    error,
    loadPlans,
    loadPayments,
    createCheckout,
    clearError,
    formatPrice,
  };
}; 