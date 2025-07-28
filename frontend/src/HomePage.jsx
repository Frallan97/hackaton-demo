import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from './AuthContext';
import config from './config.js';

const HomePage = () => {
  const { user, handleLogout, refreshToken, hasRole, authenticatedFetch, setError } = useAuth();
  const [setupLoading, setSetupLoading] = useState(false);

  const handleSetupAdmin = async () => {
    setSetupLoading(true);
    setError('');
    
    try {
      const response = await fetch(`${config.apiBaseUrl}/api/setup/first-admin`, {
        method: 'POST'
      });

      if (response.ok) {
        const data = await response.json();
        alert('Success: ' + data.message + '. Please refresh your token or login again to see admin features.');
      } else {
        const errorText = await response.text();
        if (response.status === 409) {
          alert('Admin user already exists in the system.');
        } else {
          setError('Setup failed: ' + errorText);
        }
      }
    } catch (err) {
      setError('Setup failed: ' + err.message);
    } finally {
      setSetupLoading(false);
    }
  };

  return (
    <div style={{ 
      fontFamily: 'Arial, sans-serif', 
      maxWidth: '800px', 
      margin: '50px auto', 
      padding: '20px',
      textAlign: 'center'
    }}>
      <h1>React Go App</h1>
      <p>Welcome to the React + Go micro-app with RBAC!</p>
      
      <div style={{ 
        backgroundColor: '#f8f9fa', 
        padding: '20px', 
        borderRadius: '8px', 
        marginBottom: '20px',
        textAlign: 'left'
      }}>
        <h2 style={{ marginTop: 0, color: '#333' }}>Welcome, {user?.name}!</h2>
        
        <div style={{ display: 'flex', alignItems: 'center', gap: '20px', marginBottom: '20px' }}>
          {user?.picture && (
            <img 
              src={user.picture} 
              alt="Profile" 
              style={{ 
                width: '80px', 
                height: '80px', 
                borderRadius: '50%',
                border: '3px solid #4285f4'
              }} 
            />
          )}
          <div>
            <p><strong>Email:</strong> {user?.email}</p>
            <p><strong>Name:</strong> {user?.name}</p>
            <p><strong>User ID:</strong> {user?.id}</p>
            <p><strong>Last Login:</strong> {user?.last_login_at ? new Date(user.last_login_at).toLocaleString() : 'N/A'}</p>
            <p><strong>Account Created:</strong> {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}</p>
          </div>
        </div>

        {/* Display user roles */}
        {user?.roles && user.roles.length > 0 && (
          <div style={{ marginBottom: '20px' }}>
            <h3>Your Roles:</h3>
            <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
              {user.roles.map(role => (
                <span 
                  key={role.id}
                  style={{
                    backgroundColor: '#007bff',
                    color: 'white',
                    padding: '4px 8px',
                    borderRadius: '4px',
                    fontSize: '12px'
                  }}
                >
                  {role.name}
                </span>
              ))}
            </div>
          </div>
        )}

        {/* Display user organizations */}
        {user?.organizations && user.organizations.length > 0 && (
          <div style={{ marginBottom: '20px' }}>
            <h3>Your Organizations:</h3>
            <div style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}>
              {user.organizations.map(org => (
                <span 
                  key={org.id}
                  style={{
                    backgroundColor: '#28a745',
                    color: 'white',
                    padding: '4px 8px',
                    borderRadius: '4px',
                    fontSize: '12px'
                  }}
                >
                  {org.name}
                </span>
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Navigation and action buttons */}
      <div style={{ marginBottom: '20px' }}>
        {hasRole('admin') && (
          <Link 
            to="/admin"
            style={{
              padding: '12px 24px',
              fontSize: '16px',
              backgroundColor: '#dc3545',
              color: 'white',
              textDecoration: 'none',
              borderRadius: '4px',
              marginRight: '10px',
              display: 'inline-block'
            }}
          >
            Admin Dashboard
          </Link>
        )}
        
        {!hasRole('admin') && (
          <button
            onClick={handleSetupAdmin}
            disabled={setupLoading}
            style={{
              padding: '12px 24px',
              fontSize: '16px',
              backgroundColor: '#ffc107',
              color: '#212529',
              border: 'none',
              borderRadius: '4px',
              cursor: setupLoading ? 'not-allowed' : 'pointer',
              opacity: setupLoading ? 0.7 : 1,
              marginRight: '10px'
            }}
          >
            {setupLoading ? 'Setting up...' : 'Make Me Admin'}
          </button>
        )}
      </div>

      <div style={{ display: 'flex', gap: '10px', justifyContent: 'center' }}>
        <button 
          onClick={refreshToken}
          style={{
            padding: '8px 16px',
            fontSize: '14px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Refresh Token
        </button>
        <button 
          onClick={handleLogout}
          style={{
            padding: '8px 16px',
            fontSize: '14px',
            backgroundColor: '#6c757d',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Sign Out
        </button>
      </div>
    </div>
  );
};

export default HomePage;