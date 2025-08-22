import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../AuthContext';
import config from '../config';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Badge } from '@/components/ui/badge';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Loader2, ArrowLeft, Users, Shield, Building2, Plus, Minus } from 'lucide-react';

interface User {
  id: number;
  name: string;
  email: string;
  picture?: string;
  is_active: boolean;
  last_login_at: string;
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

const AdminDashboard: React.FC = () => {
  const { authenticatedFetch, handleLogout } = useAuth();
  const [users, setUsers] = useState<User[]>([]);
  const [roles, setRoles] = useState<Role[]>([]);
  const [organizations, setOrganizations] = useState<Organization[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string>('');
  const [activeTab, setActiveTab] = useState<string>('users');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async (): Promise<void> => {
    setLoading(true);
    setError('');
    
    try {
      // Load all data in parallel
      const [usersRes, rolesRes, orgsRes] = await Promise.all([
        authenticatedFetch(`${config.apiBaseUrl}/api/admin/users`),
        authenticatedFetch(`${config.apiBaseUrl}/api/roles`),
        authenticatedFetch(`${config.apiBaseUrl}/api/organizations`)
      ]);

      if (usersRes.ok && rolesRes.ok && orgsRes.ok) {
        const [usersData, rolesData, orgsData] = await Promise.all([
          usersRes.json(),
          rolesRes.json(),
          orgsRes.json()
        ]);

        setUsers(usersData);
        setRoles(rolesData);
        setOrganizations(orgsData);
      } else {
        throw new Error('Failed to load data');
      }
    } catch (err) {
      setError('Failed to load admin data: ' + (err as Error).message);
      console.error('Admin data loading error:', err);
    } finally {
      setLoading(false);
    }
  };

  const assignRole = async (userId: number, roleId: number): Promise<void> => {
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/admin/assign-role`, {
        method: 'POST',
        body: JSON.stringify({ user_id: userId, role_id: roleId })
      });

      if (response.ok) {
        loadData(); // Reload data
      } else {
        const errorText = await response.text();
        setError('Failed to assign role: ' + errorText);
      }
    } catch (err) {
      setError('Failed to assign role: ' + (err as Error).message);
    }
  };

  const removeRole = async (userId: number, roleId: number): Promise<void> => {
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/admin/remove-role`, {
        method: 'POST',
        body: JSON.stringify({ user_id: userId, role_id: roleId })
      });

      if (response.ok) {
        loadData(); // Reload data
      } else {
        const errorText = await response.text();
        setError('Failed to remove role: ' + errorText);
      }
    } catch (err) {
      setError('Failed to remove role: ' + (err as Error).message);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
        <div className="max-w-6xl mx-auto">
          <div className="text-center">
            <Loader2 className="w-8 h-8 animate-spin mx-auto mb-4" />
            <p>Loading admin dashboard...</p>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 p-4">
      <div className="max-w-6xl mx-auto">
        {/* Header */}
        <div className="flex items-center justify-between mb-8">
          <div className="flex items-center gap-4">
            <Link to="/">
              <Button variant="outline" size="sm">
                <ArrowLeft className="w-4 h-4 mr-2" />
                Back to Home
              </Button>
            </Link>
            <div>
              <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
              <p className="text-gray-600">Manage users, roles, and organizations</p>
            </div>
          </div>
          <Button onClick={handleLogout} variant="outline">
            Sign Out
          </Button>
        </div>

        {/* Error Display */}
        {error && (
          <Alert variant="destructive" className="mb-6">
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={setActiveTab} className="space-y-6">
          <TabsList className="grid w-full grid-cols-3">
            <TabsTrigger value="users" className="flex items-center gap-2">
              <Users className="w-4 h-4" />
              Users ({users.length})
            </TabsTrigger>
            <TabsTrigger value="roles" className="flex items-center gap-2">
              <Shield className="w-4 h-4" />
              Roles ({roles.length})
            </TabsTrigger>
            <TabsTrigger value="organizations" className="flex items-center gap-2">
              <Building2 className="w-4 h-4" />
              Organizations ({organizations.length})
            </TabsTrigger>
          </TabsList>

          {/* Users Tab */}
          <TabsContent value="users" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>User Management</CardTitle>
                <CardDescription>View and manage user accounts</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4">
                  {users.map((user) => (
                    <Card key={user.id} className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-4">
                          <div className="w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                            <span className="text-blue-600 font-semibold">
                              {user.name.charAt(0).toUpperCase()}
                            </span>
                          </div>
                          <div>
                            <h3 className="font-semibold">{user.name}</h3>
                            <p className="text-sm text-gray-600">{user.email}</p>
                            <div className="flex gap-2 mt-2">
                              {user.roles?.map((role) => (
                                <Badge key={role.id} variant="secondary">
                                  {role.name}
                                </Badge>
                              ))}
                            </div>
                          </div>
                        </div>
                        <div className="flex gap-2">
                          {roles.map((role) => {
                            const hasRole = user.roles?.some(r => r.id === role.id);
                            return (
                              <Button
                                key={role.id}
                                size="sm"
                                variant={hasRole ? "destructive" : "outline"}
                                onClick={() => hasRole 
                                  ? removeRole(user.id, role.id)
                                  : assignRole(user.id, role.id)
                                }
                              >
                                {hasRole ? <Minus className="w-3 h-3" /> : <Plus className="w-3 h-3" />}
                                {role.name}
                              </Button>
                            );
                          })}
                        </div>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Roles Tab */}
          <TabsContent value="roles" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Role Management</CardTitle>
                <CardDescription>View available roles in the system</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4">
                  {roles.map((role) => (
                    <Card key={role.id} className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold">{role.name}</h3>
                          {role.description && (
                            <p className="text-sm text-gray-600">{role.description}</p>
                          )}
                        </div>
                        <Badge variant="outline">{role.name}</Badge>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>

          {/* Organizations Tab */}
          <TabsContent value="organizations" className="space-y-4">
            <Card>
              <CardHeader>
                <CardTitle>Organization Management</CardTitle>
                <CardDescription>View available organizations</CardDescription>
              </CardHeader>
              <CardContent>
                <div className="grid gap-4">
                  {organizations.map((org) => (
                    <Card key={org.id} className="p-4">
                      <div className="flex items-center justify-between">
                        <div>
                          <h3 className="font-semibold">{org.name}</h3>
                          {org.description && (
                            <p className="text-sm text-gray-600">{org.description}</p>
                          )}
                        </div>
                        <Badge variant="outline">{org.name}</Badge>
                      </div>
                    </Card>
                  ))}
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </div>
    </div>
  );
};

export default AdminDashboard; 