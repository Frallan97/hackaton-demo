import React from 'react';
import { Navigate } from 'react-router-dom';
import { useAuth } from './AuthContext';

const ProtectedRoute = ({ children, requiredRole = null }) => {
  const { isLoggedIn, user, loading } = useAuth();

  if (loading) {
    return (
      <div style={{ 
        fontFamily: 'Arial, sans-serif', 
        maxWidth: '600px', 
        margin: '50px auto', 
        padding: '20px',
        textAlign: 'center'
      }}>
        <h1>React Go App</h1>
        <div style={{ marginTop: '20px' }}>
          <p>Loading...</p>
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
        <div style={{ 
          fontFamily: 'Arial, sans-serif', 
          maxWidth: '600px', 
          margin: '50px auto', 
          padding: '20px',
          textAlign: 'center'
        }}>
          <h1>Access Denied</h1>
          <p>You don't have permission to access this page.</p>
          <p>Required role: {requiredRole}</p>
        </div>
      );
    }
  }

  return children;
};

export default ProtectedRoute;