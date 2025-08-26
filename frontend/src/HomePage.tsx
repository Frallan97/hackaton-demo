import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from './AuthContext';
import { useAppSelector, useAppDispatch } from './store/hooks';
import { showSuccess, showError } from './store/slices/uiSlice';
import { useSetupFirstAdminMutation } from './store/api';
import config from './config';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Loader2, Shield, Building2, User, Calendar, Mail, Key, RefreshCw, LogOut } from 'lucide-react';
import { ThemeToggle } from './components/ThemeToggle';
import { ReduxDemo } from './components/ReduxDemo';
import { StripeDemo } from './components/StripeDemo';

const HomePage: React.FC = () => {
  const { user, handleLogout, refreshToken, hasRole, authenticatedFetch, setError } = useAuth();
  const dispatch = useAppDispatch();
  const [setupFirstAdmin, { isLoading: setupLoading }] = useSetupFirstAdminMutation();

  const handleSetupAdmin = async (): Promise<void> => {
    try {
      const result = await setupFirstAdmin().unwrap();
      dispatch(showSuccess('Admin setup successful! Please refresh your token or login again to see admin features.'));
    } catch (err: any) {
      if (err.status === 409) {
        dispatch(showError('Admin user already exists in the system.'));
      } else {
        dispatch(showError('Setup failed: ' + (err.data?.message || err.message || 'Unknown error')));
      }
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 p-4">
      <div className="max-w-4xl mx-auto">
        {/* Header with theme toggle */}
        <div className="flex items-center justify-between mb-8">
          <div className="text-center flex-1">
            <h1 className="text-4xl font-bold text-gray-900 dark:text-white mb-4">Hackaton Demo</h1>
            <p className="text-xl text-gray-600 dark:text-gray-300">Welcome to the React + Go micro-app with RBAC!</p>
          </div>
          <div className="flex items-center gap-4">
            <ThemeToggle />
            <Button
              onClick={handleLogout}
              variant="outline"
              className="border-red-200 text-red-600 hover:bg-red-50 dark:border-red-800 dark:text-red-400 dark:hover:bg-red-900/20"
            >
              <LogOut className="w-4 h-4 mr-2" />
              Logout
            </Button>
          </div>
        </div>
        
        <Card className="mb-8 shadow-xl dark:bg-gray-800 dark:border-gray-700">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-2xl dark:text-white">
              <User className="w-6 h-6 text-blue-600 dark:text-blue-400" />
              Welcome, {user?.name}!
            </CardTitle>
            <CardDescription className="dark:text-gray-400">Your account information and permissions</CardDescription>
          </CardHeader>
          
          <CardContent className="space-y-6">
            <div className="flex items-start gap-6">
              <Avatar className="w-20 h-20 border-4 border-blue-500 dark:border-blue-400">
                <AvatarImage src={user?.picture} alt="Profile" />
                <AvatarFallback className="text-2xl bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-400">
                  {user?.name?.charAt(0)?.toUpperCase() || 'U'}
                </AvatarFallback>
              </Avatar>
              
              <div className="flex-1 space-y-3">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="flex items-center gap-2">
                    <Mail className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="font-medium dark:text-gray-300">Email:</span>
                    <span className="text-gray-700 dark:text-gray-300">{user?.email}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Key className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="font-medium dark:text-gray-300">User ID:</span>
                    <span className="text-gray-700 dark:text-gray-300">{user?.id}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="font-medium dark:text-gray-300">Last Login:</span>
                    <span className="text-gray-700 dark:text-gray-300">
                      {user?.last_login_at ? new Date(user.last_login_at).toLocaleString() : 'N/A'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                    <span className="font-medium dark:text-gray-300">Created:</span>
                    <span className="text-gray-700 dark:text-gray-300">
                      {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <Separator className="dark:bg-gray-700" />

            {/* Display user roles */}
            {user?.roles && user.roles.length > 0 && (
              <div>
                <h3 className="flex items-center gap-2 text-lg font-semibold mb-3 dark:text-white">
                  <Shield className="w-5 h-5 text-blue-600 dark:text-blue-400" />
                  Your Roles
                </h3>
                <div className="flex gap-2 flex-wrap">
                  {user.roles.map(role => (
                    <Badge key={role.id} variant="default" className="bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600">
                      {role.name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {/* Display user organizations */}
            {user?.organizations && user.organizations.length > 0 && (
              <div>
                <h3 className="flex items-center gap-2 text-lg font-semibold mb-3 dark:text-white">
                  <Building2 className="w-5 h-5 text-green-600 dark:text-green-400" />
                  Your Organizations
                </h3>
                <div className="flex gap-2 flex-wrap">
                  {user.organizations.map(org => (
                    <Badge key={org.id} variant="secondary" className="bg-green-600 hover:bg-green-700 text-white dark:bg-green-500 dark:hover:bg-green-600">
                      {org.name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Redux Demo Component */}
        <div className="mb-8">
          <ReduxDemo />
        </div>

        {/* Stripe Integration Demo */}
        <div className="mb-8">
          <StripeDemo />
        </div>

        {/* Navigation and action buttons */}
        <div className="flex flex-wrap gap-4 justify-center mb-8">
          {hasRole('admin') && (
            <Link to="/admin">
              <Button variant="destructive" size="lg" className="shadow-lg">
                <Shield className="w-4 h-4 mr-2" />
                Admin Dashboard
              </Button>
            </Link>
          )}
          
          {!hasRole('admin') && (
            <Button
              onClick={handleSetupAdmin}
              disabled={setupLoading}
              variant="outline"
              size="lg"
              className="shadow-lg border-yellow-500 text-yellow-700 hover:bg-yellow-50"
            >
              {setupLoading ? (
                <>
                  <Loader2 className="w-4 h-4 mr-2 animate-spin" />
                  Setting up...
                </>
              ) : (
                <>
                  <Shield className="w-4 h-4 mr-2" />
                  Make Me Admin
                </>
              )}
            </Button>
          )}
        </div>

        <div className="flex flex-wrap gap-4 justify-center">
          <Button 
            onClick={refreshToken}
            variant="outline"
            className="shadow-lg"
          >
            <RefreshCw className="w-4 h-4 mr-2" />
            Refresh Token
          </Button>
          <Button 
            onClick={handleLogout}
            variant="secondary"
            className="shadow-lg"
          >
            <LogOut className="w-4 h-4 mr-2" />
            Sign Out
          </Button>
        </div>
      </div>
    </div>
  );
};

export default HomePage; 