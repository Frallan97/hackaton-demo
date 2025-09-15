import React from 'react';
import { PaymentCard } from './PaymentCard';
import { PaymentPlan } from '../../services/stripeService';

export interface PaymentListProps {
  plans: PaymentPlan[];
  onPlanSelect?: (plan: PaymentPlan) => void;
  selectedPlanId?: string;
  purchasedPlanIds?: string[];
  loading?: boolean;
  layout?: 'grid' | 'single' | 'horizontal';
  variant?: 'default' | 'compact' | 'featured';
  showFeatures?: boolean;
  className?: string;
  formatPrice: (amount: number, currency?: string) => string;
}

export const PaymentList: React.FC<PaymentListProps> = ({
  plans,
  onPlanSelect,
  selectedPlanId,
  purchasedPlanIds = [],
  loading = false,
  layout = 'single',
  variant = 'default',
  showFeatures = true,
  className = '',
  formatPrice,
}) => {
  if (!plans || plans.length === 0) {
    return (
      <div className="text-center py-8">
        <p className="text-gray-600">No payment plans available.</p>
      </div>
    );
  }

  const getLayoutClasses = () => {
    switch (layout) {
      case 'grid':
        return 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6';
      case 'horizontal':
        return 'flex flex-wrap gap-4 justify-center';
      case 'single':
      default:
        return 'flex justify-center';
    }
  };

  const getCardClasses = () => {
    switch (layout) {
      case 'single':
        return 'w-full max-w-md';
      case 'horizontal':
        return 'flex-1 min-w-64 max-w-sm';
      case 'grid':
      default:
        return '';
    }
  };

  return (
    <div className={`${getLayoutClasses()} ${className}`}>
      {plans.map((plan, index) => (
        <PaymentCard
          key={plan.id}
          plan={plan}
          onSelect={onPlanSelect}
          isSelected={selectedPlanId === plan.id}
          isPurchased={purchasedPlanIds.includes(plan.id)}
          loading={loading}
          variant={plans.length === 1 ? 'featured' : variant}
          showFeatures={showFeatures}
          formatPrice={formatPrice}
          className={getCardClasses()}
        />
      ))}
    </div>
  );
}; 