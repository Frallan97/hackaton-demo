import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from './AuthContext';
import config from './config';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Separator } from '@/components/ui/separator';
import { Loader2, Shield, Building2, User, Calendar, Mail, Key, RefreshCw, LogOut } from 'lucide-react';

const HomePage: React.FC = () => {
  const { user, handleLogout, refreshToken, hasRole, authenticatedFetch, setError } = useAuth();
  const [setupLoading, setSetupLoading] = useState<boolean>(false);

  const handleSetupAdmin = async (): Promise<void> => {
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
      setError('Setup failed: ' + (err as Error).message);
    } finally {
      setSetupLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-8">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">Hackaton Demo</h1>
          <p className="text-xl text-gray-600">Welcome to the React + Go micro-app with RBAC!</p>
        </div>
        
        <Card className="mb-8 shadow-xl">
          <CardHeader>
            <CardTitle className="flex items-center gap-3 text-2xl">
              <User className="w-6 h-6 text-blue-600" />
              Welcome, {user?.name}!
            </CardTitle>
            <CardDescription>Your account information and permissions</CardDescription>
          </CardHeader>
          
          <CardContent className="space-y-6">
            <div className="flex items-start gap-6">
              <Avatar className="w-20 h-20 border-4 border-blue-500">
                <AvatarImage src={user?.picture} alt="Profile" />
                <AvatarFallback className="text-2xl bg-blue-100 text-blue-600">
                  {user?.name?.charAt(0)?.toUpperCase() || 'U'}
                </AvatarFallback>
              </Avatar>
              
              <div className="flex-1 space-y-3">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="flex items-center gap-2">
                    <Mail className="w-4 h-4 text-gray-500" />
                    <span className="font-medium">Email:</span>
                    <span className="text-gray-700">{user?.email}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Key className="w-4 h-4 text-gray-500" />
                    <span className="font-medium">User ID:</span>
                    <span className="text-gray-700">{user?.id}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4 text-gray-500" />
                    <span className="font-medium">Last Login:</span>
                    <span className="text-gray-700">
                      {user?.last_login_at ? new Date(user.last_login_at).toLocaleString() : 'N/A'}
                    </span>
                  </div>
                  <div className="flex items-center gap-2">
                    <Calendar className="w-4 h-4 text-gray-500" />
                    <span className="font-medium">Created:</span>
                    <span className="text-gray-700">
                      {user?.created_at ? new Date(user.created_at).toLocaleDateString() : 'N/A'}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <Separator />

            {/* Display user roles */}
            {user?.roles && user.roles.length > 0 && (
              <div>
                <h3 className="flex items-center gap-2 text-lg font-semibold mb-3">
                  <Shield className="w-5 h-5 text-blue-600" />
                  Your Roles
                </h3>
                <div className="flex gap-2 flex-wrap">
                  {user.roles.map(role => (
                    <Badge key={role.id} variant="default" className="bg-blue-600 hover:bg-blue-700">
                      {role.name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}

            {/* Display user organizations */}
            {user?.organizations && user.organizations.length > 0 && (
              <div>
                <h3 className="flex items-center gap-2 text-lg font-semibold mb-3">
                  <Building2 className="w-5 h-5 text-green-600" />
                  Your Organizations
                </h3>
                <div className="flex gap-2 flex-wrap">
                  {user.organizations.map(org => (
                    <Badge key={org.id} variant="secondary" className="bg-green-600 hover:bg-green-700 text-white">
                      {org.name}
                    </Badge>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

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