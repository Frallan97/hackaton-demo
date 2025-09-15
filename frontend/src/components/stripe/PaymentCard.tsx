import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Button } from '../ui/button';
import { Badge } from '../ui/badge';
import { Check, Zap } from 'lucide-react';
import { PaymentPlan } from '../../services/stripeService';

export interface PaymentCardProps {
  plan: PaymentPlan;
  onSelect?: (plan: PaymentPlan) => void;
  isSelected?: boolean;
  isPurchased?: boolean;
  loading?: boolean;
  className?: string;
  variant?: 'default' | 'compact' | 'featured';
  showFeatures?: boolean;
  formatPrice: (amount: number, currency?: string) => string;
}

export const PaymentCard: React.FC<PaymentCardProps> = ({
  plan,
  onSelect,
  isSelected = false,
  isPurchased = false,
  loading = false,
  className = '',
  variant = 'default',
  showFeatures = true,
  formatPrice,
}) => {
  const handleSelect = () => {
    if (!isPurchased && !loading && onSelect) {
      onSelect(plan);
    }
  };

  const getCardStyle = () => {
    if (variant === 'featured') {
      return 'bg-gradient-to-br from-blue-50 to-indigo-100 border-blue-200 ring-2 ring-blue-300';
    }
    if (variant === 'compact') {
      return 'bg-white border-gray-200';
    }
    return 'bg-gradient-to-br from-blue-50 to-indigo-100 border-blue-200';
  };

  const getButtonText = () => {
    if (isPurchased) return 'Already Purchased';
    if (loading) return 'Processing...';
    if (variant === 'compact') return 'Select';
    return '🚀 Purchase Now';
  };

  return (
    <Card 
      className={`relative transition-all duration-200 hover:shadow-lg ${getCardStyle()} ${
        isSelected ? 'ring-2 ring-blue-500 shadow-lg' : ''
      } ${className}`}
    >
      {isPurchased && (
        <Badge className="absolute -top-2 left-1/2 transform -translate-x-1/2 bg-green-600 text-white">
          Purchased
        </Badge>
      )}

      {variant === 'featured' && !isPurchased && (
        <Badge className="absolute -top-2 left-1/2 transform -translate-x-1/2 bg-blue-600 text-white">
          Recommended
        </Badge>
      )}

      <CardHeader className={`text-center ${variant === 'compact' ? 'pb-2' : 'pb-4'}`}>
        {variant !== 'compact' && (
          <div className="flex justify-center mb-4">
            <div className="p-3 bg-blue-600 rounded-full">
              <Zap className="h-8 w-8 text-white" />
            </div>
          </div>
        )}
        
        <CardTitle className={variant === 'compact' ? 'text-lg' : 'text-2xl font-bold text-gray-900'}>
          {plan.name}
        </CardTitle>
        
        <CardDescription className={`text-gray-600 ${variant === 'compact' ? 'text-sm' : 'text-lg'}`}>
          {plan.description}
        </CardDescription>
      </CardHeader>

      <CardContent className="text-center">
        <div className={variant === 'compact' ? 'mb-4' : 'mb-8'}>
          <div className={`font-bold text-gray-900 ${variant === 'compact' ? 'text-2xl mb-1' : 'text-4xl mb-2'}`}>
            {formatPrice(plan.price, plan.currency)}
          </div>
          <div className={`text-gray-600 ${variant === 'compact' ? 'text-sm' : 'text-lg'}`}>
            one-time payment
          </div>
          {variant !== 'compact' && (
            <div className="text-sm text-blue-600 mt-2">
              💳 Card payments (📱 Swish ready when configured)
            </div>
          )}
        </div>

        {showFeatures && plan.features && plan.features.length > 0 && (
          <div className={`text-left bg-white rounded-lg p-4 shadow-sm ${variant === 'compact' ? 'mb-4 space-y-2' : 'mb-8 space-y-4'}`}>
            {variant !== 'compact' && (
              <h4 className="font-semibold text-gray-900 text-center mb-3">What's included:</h4>
            )}
            {plan.features.map((feature, index) => (
              <div key={index} className="flex items-center">
                <Check className={`text-green-600 mr-3 flex-shrink-0 ${variant === 'compact' ? 'h-4 w-4' : 'h-5 w-5'}`} />
                <span className={`text-gray-700 ${variant === 'compact' ? 'text-sm' : ''}`}>{feature}</span>
              </div>
            ))}
          </div>
        )}

        {onSelect && (
          <Button
            onClick={handleSelect}
            className={`w-full ${variant === 'compact' ? 'py-2' : 'text-lg py-3'} ${
              isPurchased
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-blue-600 hover:bg-blue-700 shadow-lg hover:shadow-xl transform hover:scale-105 transition-all duration-200'
            }`}
            disabled={isPurchased || loading}
          >
            {getButtonText()}
          </Button>
        )}
      </CardContent>
    </Card>
  );
}; 