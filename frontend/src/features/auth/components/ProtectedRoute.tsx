import React, { ReactNode } from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from './AuthContext';

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRole?: string | null;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children, requiredRole = null }) => {
  const { isLoggedIn, user, loading } = useAuth();

  if (loading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
        <div className="text-center max-w-md mx-auto p-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-6">Hackaton Demo</h1>
          <div className="bg-white rounded-lg shadow-lg p-6">
            <p className="text-lg text-gray-700">Loading...</p>
          </div>
        </div>
      </div>
    );
  }

  if (!isLoggedIn) {
    return <Navigate to="/login" replace />;
  }

  // Check role requirement
  if (requiredRole && user) {
    const hasRequiredRole = user.roles?.some(role => role.name === requiredRole) || false;
    if (!hasRequiredRole) {
      return (
        <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-red-50 to-pink-100">
          <div className="text-center max-w-md mx-auto p-8">
            <h1 className="text-3xl font-bold text-red-900 mb-6">Access Denied</h1>
            <div className="bg-white rounded-lg shadow-lg p-6">
              <p className="text-lg text-red-700 mb-2">You don't have permission to access this page.</p>
              <p className="text-sm text-red-600">Required role: {requiredRole}</p>
            </div>
          </div>
        </div>
      );
    }
  }

  return <>{children}</>;
};

export default ProtectedRoute; 