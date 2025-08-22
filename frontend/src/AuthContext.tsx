import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import config from './config';

// Types
interface User {
  id: number;
  email: string;
  name: string;
  picture?: string;
  google_id: string;
  is_active: boolean;
  last_login_at: string;
  created_at: string;
  updated_at: string;
  roles?: Role[];
  organizations?: Organization[];
}

interface Role {
  id: number;
  name: string;
  description?: string;
}

interface Organization {
  id: number;
  name: string;
  description?: string;
}

interface AuthContextType {
  isLoggedIn: boolean;
  user: User | null;
  loading: boolean;
  error: string;
  handleGoogleLogin: () => Promise<void>;
  handleLogout: () => Promise<void>;
  refreshToken: () => Promise<boolean>;
  authenticatedFetch: (url: string, options?: RequestInit) => Promise<Response>;
  hasRole: (roleName: string) => boolean;
  setError: (error: string) => void;
}

// Create Auth Context
const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const useAuth = (): AuthContextType => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState<boolean>(false);
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');

  // Check if user is already logged in on component mount
  useEffect(() => {
    checkAuthStatus();
  }, []);

  // Check if there's an authorization code in the URL (OAuth callback)
  useEffect(() => {
    const urlParams = new URLSearchParams(window.location.search);
    const code = urlParams.get('code');
    
    if (code) {
      handleOAuthCallback(code);
    }
  }, []);

  const checkAuthStatus = async (): Promise<void> => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      setLoading(false);
      return;
    }

    try {
      const response = await fetch(`${config.apiBaseUrl}/api/auth/me`, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const userData: User = await response.json();
        
        // Try to get user roles and organizations from admin endpoint
        try {
          const [rolesRes, orgsRes] = await Promise.all([
            fetch(`${config.apiBaseUrl}/api/admin/user-roles?id=${userData.id}`, {
              headers: { 'Authorization': `Bearer ${token}` }
            }),
            fetch(`${config.apiBaseUrl}/api/admin/user-organizations?id=${userData.id}`, {
              headers: { 'Authorization': `Bearer ${token}` }
            })
          ]);

          if (rolesRes.ok) {
            const roles: Role[] = await rolesRes.json();
            userData.roles = roles;
          } else {
            userData.roles = [];
          }

          if (orgsRes.ok) {
            const organizations: Organization[] = await orgsRes.json();
            userData.organizations = organizations;
          } else {
            userData.organizations = [];
          }
        } catch (roleOrgError) {
          // If we can't get roles/orgs, that's OK - user might not have admin access
          console.log('Could not fetch user roles/organizations:', roleOrgError);
          userData.roles = [];
          userData.organizations = [];
        }

        setUser(userData);
        setIsLoggedIn(true);
      } else {
        // Token is invalid, clear it
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      localStorage.removeItem('access_token');
      localStorage.removeItem('refresh_token');
    } finally {
      setLoading(false);
    }
  };

  const handleGoogleLogin = async (): Promise<void> => {
    setLoading(true);
    setError('');

    try {
      // Get the OAuth URL from our backend
      const response = await fetch(`${config.apiBaseUrl}/api/auth/google/url`);
      const data = await response.json();
      
      // Redirect to Google OAuth
      window.location.href = data.auth_url;
    } catch (error) {
      setError('Failed to start login process. Please try again.');
      setLoading(false);
    }
  };

  const handleOAuthCallback = async (code: string): Promise<void> => {
    setLoading(true);
    setError('');

    try {
      const response = await fetch(`${config.apiBaseUrl}/api/auth/google/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ code })
      });

      if (response.ok) {
        const authData = await response.json();
        
        // Store tokens
        localStorage.setItem('access_token', authData.access_token);
        localStorage.setItem('refresh_token', authData.refresh_token);
        
        // Set user data
        setUser(authData.user);
        setIsLoggedIn(true);
        
        // Clear the URL parameters
        window.history.replaceState({}, document.title, window.location.pathname);
      } else {
        const errorData = await response.text();
        setError(`Login failed: ${errorData}`);
      }
    } catch (error) {
      setError('Login failed. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = async (): Promise<void> => {
    try {
      // Call logout endpoint
      await fetch(`${config.apiBaseUrl}/api/auth/logout`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('access_token')}`
        }
      });
    } catch (error) {
      console.error('Logout request failed:', error);
    }

    // Clear local storage and state
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setIsLoggedIn(false);
    setUser(null);
    setError('');
  };

  const refreshToken = async (): Promise<boolean> => {
    const refreshToken = localStorage.getItem('refresh_token');
    if (!refreshToken) return false;

    try {
      const response = await fetch(`${config.apiBaseUrl}/api/auth/refresh`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ refresh_token: refreshToken })
      });

      if (response.ok) {
        const data = await response.json();
        localStorage.setItem('access_token', data.access_token);
        return true;
      }
    } catch (error) {
      console.error('Token refresh failed:', error);
    }

    return false;
  };

  // Helper function to make authenticated API calls
  const authenticatedFetch = async (url: string, options: RequestInit = {}): Promise<Response> => {
    const token = localStorage.getItem('access_token');
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };

    if (token) {
      headers.Authorization = `Bearer ${token}`;
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      });

      // If token expired, try to refresh
      if (response.status === 401) {
        const refreshed = await refreshToken();
        if (refreshed) {
          // Retry the request with new token
          const newToken = localStorage.getItem('access_token');
          headers.Authorization = `Bearer ${newToken}`;
          return await fetch(url, {
            ...options,
            headers,
          });
        } else {
          // Refresh failed, logout user
          handleLogout();
          throw new Error('Authentication failed');
        }
      }

      return response;
    } catch (error) {
      throw error;
    }
  };

  // Check if user has a specific role
  const hasRole = (roleName: string): boolean => {
    return user?.roles?.some(role => role.name === roleName) || false;
  };

  const value: AuthContextType = {
    isLoggedIn,
    user,
    loading,
    error,
    handleGoogleLogin,
    handleLogout,
    refreshToken,
    authenticatedFetch,
    hasRole,
    setError
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
}; 