import React from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '../ui/card';
import { Badge } from '../ui/badge';
import { CheckCircle, AlertCircle, XCircle, CreditCard, Clock } from 'lucide-react';

export interface PaymentStatusProps {
  status: 'active' | 'inactive' | 'pending' | 'failed' | 'cancelled' | 'none';
  title?: string;
  description?: string;
  amount?: number;
  currency?: string;
  nextBilling?: string;
  planName?: string;
  showCard?: boolean;
  className?: string;
  formatPrice?: (amount: number, currency?: string) => string;
}

export const PaymentStatus: React.FC<PaymentStatusProps> = ({
  status,
  title,
  description,
  amount,
  currency = 'usd',
  nextBilling,
  planName,
  showCard = true,
  className = '',
  formatPrice,
}) => {
  const getStatusConfig = () => {
    switch (status) {
      case 'active':
        return {
          icon: <CheckCircle className="h-5 w-5 text-green-600" />,
          badgeClass: 'bg-green-100 text-green-800',
          badgeText: 'Active',
          title: title || 'Payment Active',
          description: description || 'Your payment is active and working properly.',
        };
      case 'pending':
        return {
          icon: <Clock className="h-5 w-5 text-yellow-600" />,
          badgeClass: 'bg-yellow-100 text-yellow-800',
          badgeText: 'Pending',
          title: title || 'Payment Pending',
          description: description || 'Your payment is being processed.',
        };
      case 'failed':
        return {
          icon: <XCircle className="h-5 w-5 text-red-600" />,
          badgeClass: 'bg-red-100 text-red-800',
          badgeText: 'Failed',
          title: title || 'Payment Failed',
          description: description || 'There was an issue with your payment.',
        };
      case 'cancelled':
        return {
          icon: <AlertCircle className="h-5 w-5 text-gray-600" />,
          badgeClass: 'bg-gray-100 text-gray-800',
          badgeText: 'Cancelled',
          title: title || 'Payment Cancelled',
          description: description || 'Your payment has been cancelled.',
        };
      case 'inactive':
        return {
          icon: <CreditCard className="h-5 w-5 text-gray-400" />,
          badgeClass: 'bg-gray-100 text-gray-800',
          badgeText: 'Inactive',
          title: title || 'No Active Payment',
          description: description || 'You don\'t have any active payments.',
        };
      case 'none':
      default:
        return {
          icon: <CreditCard className="h-5 w-5 text-gray-400" />,
          badgeClass: 'bg-gray-100 text-gray-800',
          badgeText: 'No Payment',
          title: title || 'No Payment Method',
          description: description || 'No payment method configured.',
        };
    }
  };

  const statusConfig = getStatusConfig();

  const content = (
    <div className="flex items-center justify-between">
      <div className="flex items-center space-x-3">
        <div className="p-2 bg-blue-100 rounded-lg">
          {statusConfig.icon}
        </div>
        <div>
          <h3 className="text-lg font-semibold">{statusConfig.title}</h3>
          <p className="text-gray-600 text-sm">{statusConfig.description}</p>
          {planName && (
            <p className="text-blue-600 text-sm font-medium mt-1">{planName}</p>
          )}
        </div>
      </div>
      
      <div className="text-right">
        <Badge className={statusConfig.badgeClass}>
          {statusConfig.badgeText}
        </Badge>
        {amount && formatPrice && (
          <div className="text-lg font-semibold text-gray-900 mt-1">
            {formatPrice(amount, currency)}
          </div>
        )}
        {nextBilling && (
          <div className="text-xs text-gray-500 mt-1">
            Next: {new Date(nextBilling).toLocaleDateString()}
          </div>
        )}
      </div>
    </div>
  );

  if (!showCard) {
    return <div className={className}>{content}</div>;
  }

  return (
    <Card className={className}>
      <CardContent className="pt-6">
        {content}
      </CardContent>
    </Card>
  );
}; 