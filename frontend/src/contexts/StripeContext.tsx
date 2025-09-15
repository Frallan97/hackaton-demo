import React, { createContext, useContext, ReactNode } from 'react';
import { useStripe, UseStripeReturn } from '../hooks/useStripe';

const StripeContext = createContext<UseStripeReturn | undefined>(undefined);

export const useStripeContext = (): UseStripeReturn => {
  const context = useContext(StripeContext);
  if (!context) {
    throw new Error('useStripeContext must be used within a StripeProvider');
  }
  return context;
};

interface StripeProviderProps {
  children: ReactNode;
}

export const StripeProvider: React.FC<StripeProviderProps> = ({ children }) => {
  const stripeHook = useStripe();

  return (
    <StripeContext.Provider value={stripeHook}>
      {children}
    </StripeContext.Provider>
  );
}; 