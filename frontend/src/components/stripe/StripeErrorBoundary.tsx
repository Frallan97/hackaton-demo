import React from 'react';
import { Card, CardContent } from '../ui/card';
import { Button } from '../ui/button';
import { AlertCircle, RefreshCw } from 'lucide-react';

export interface StripeErrorProps {
  error: string;
  onRetry?: () => void;
  onDismiss?: () => void;
  title?: string;
  showRetry?: boolean;
  variant?: 'inline' | 'card' | 'banner';
  className?: string;
}

export const StripeError: React.FC<StripeErrorProps> = ({
  error,
  onRetry,
  onDismiss,
  title = 'Payment Error',
  showRetry = true,
  variant = 'card',
  className = '',
}) => {
  const content = (
    <div className="flex items-start space-x-3">
      <AlertCircle className="h-5 w-5 text-red-600 flex-shrink-0 mt-0.5" />
      <div className="flex-1">
        <h3 className="text-red-800 font-medium">{title}</h3>
        <p className="text-red-600 text-sm mt-1">{error}</p>
      </div>
      <div className="flex space-x-2">
        {showRetry && onRetry && (
          <Button
            onClick={onRetry}
            variant="outline"
            size="sm"
            className="border-red-200 text-red-700 hover:bg-red-50"
          >
            <RefreshCw className="h-4 w-4 mr-1" />
            Retry
          </Button>
        )}
        {onDismiss && (
          <Button
            onClick={onDismiss}
            variant="outline"
            size="sm"
            className="border-red-200 text-red-700 hover:bg-red-50"
          >
            Dismiss
          </Button>
        )}
      </div>
    </div>
  );

  if (variant === 'inline') {
    return <div className={`text-red-600 ${className}`}>{content}</div>;
  }

  if (variant === 'banner') {
    return (
      <div className={`bg-red-50 border border-red-200 rounded-lg p-4 ${className}`}>
        {content}
      </div>
    );
  }

  return (
    <Card className={`border-red-200 bg-red-50 ${className}`}>
      <CardContent className="pt-6">
        {content}
      </CardContent>
    </Card>
  );
};

export interface StripeLoadingProps {
  message?: string;
  variant?: 'inline' | 'card' | 'overlay';
  className?: string;
}

export const StripeLoading: React.FC<StripeLoadingProps> = ({
  message = 'Loading...',
  variant = 'inline',
  className = '',
}) => {
  const spinner = (
    <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
  );

  const content = (
    <div className="flex items-center justify-center space-x-3 py-8">
      {spinner}
      <span className="text-gray-600">{message}</span>
    </div>
  );

  if (variant === 'inline') {
    return <div className={className}>{content}</div>;
  }

  if (variant === 'overlay') {
    return (
      <div className={`absolute inset-0 bg-white bg-opacity-75 flex items-center justify-center ${className}`}>
        {content}
      </div>
    );
  }

  return (
    <Card className={className}>
      <CardContent>
        {content}
      </CardContent>
    </Card>
  );
}; 