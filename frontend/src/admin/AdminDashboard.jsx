import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../AuthContext';
import config from '../config.js';

const AdminDashboard = () => {
  const { authenticatedFetch, handleLogout } = useAuth();
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [organizations, setOrganizations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [activeTab, setActiveTab] = useState('users');

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
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
      setError('Failed to load admin data: ' + err.message);
      console.error('Admin data loading error:', err);
    } finally {
      setLoading(false);
    }
  };

  const assignRole = async (userId, roleId) => {
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
      setError('Failed to assign role: ' + err.message);
    }
  };

  const removeRole = async (userId, roleId) => {
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
      setError('Failed to remove role: ' + err.message);
    }
  };

  const addToOrganization = async (userId, organizationId, role = 'member') => {
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/admin/assign-organization`, {
        method: 'POST',
        body: JSON.stringify({ user_id: userId, organization_id: organizationId, role })
      });

      if (response.ok) {
        loadData(); // Reload data
      } else {
        const errorText = await response.text();
        setError('Failed to add to organization: ' + errorText);
      }
    } catch (err) {
      setError('Failed to add to organization: ' + err.message);
    }
  };

  const removeFromOrganization = async (userId, organizationId) => {
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/admin/remove-organization`, {
        method: 'POST',
        body: JSON.stringify({ user_id: userId, organization_id: organizationId })
      });

      if (response.ok) {
        loadData(); // Reload data
      } else {
        const errorText = await response.text();
        setError('Failed to remove from organization: ' + errorText);
      }
    } catch (err) {
      setError('Failed to remove from organization: ' + err.message);
    }
  };

  if (loading) {
    return (
      <div style={{ fontFamily: 'Arial, sans-serif', maxWidth: '1200px', margin: '20px auto', padding: '20px' }}>
        <h1>Admin Dashboard</h1>
        <p>Loading...</p>
      </div>
    );
  }

  return (
    <div style={{ fontFamily: 'Arial, sans-serif', maxWidth: '1200px', margin: '20px auto', padding: '20px' }}>
      {/* Header */}
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h1>Admin Dashboard</h1>
        <div style={{ display: 'flex', gap: '10px' }}>
          <Link 
            to="/"
            style={{
              padding: '8px 16px',
              backgroundColor: '#6c757d',
              color: 'white',
              textDecoration: 'none',
              borderRadius: '4px'
            }}
          >
            Back to Home
          </Link>
          <button 
            onClick={handleLogout}
            style={{
              padding: '8px 16px',
              backgroundColor: '#dc3545',
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

      {/* Error Display */}
      {error && (
        <div style={{ 
          margin: '20px 0', 
          padding: '10px', 
          backgroundColor: '#f8d7da', 
          color: '#721c24', 
          borderRadius: '4px',
          border: '1px solid #f5c6cb'
        }}>
          {error}
        </div>
      )}

      {/* Tabs */}
      <div style={{ display: 'flex', borderBottom: '1px solid #ddd', marginBottom: '20px' }}>
        {[
          { key: 'users', label: 'Users' },
          { key: 'roles', label: 'Roles' },
          { key: 'organizations', label: 'Organizations' }
        ].map(tab => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            style={{
              padding: '10px 20px',
              border: 'none',
              backgroundColor: activeTab === tab.key ? '#007bff' : 'transparent',
              color: activeTab === tab.key ? 'white' : '#007bff',
              cursor: 'pointer',
              borderBottom: activeTab === tab.key ? '2px solid #007bff' : '2px solid transparent'
            }}
          >
            {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      {activeTab === 'users' && (
        <UsersTab 
          users={users} 
          roles={roles} 
          organizations={organizations}
          onAssignRole={assignRole}
          onRemoveRole={removeRole}
          onAddToOrganization={addToOrganization}
          onRemoveFromOrganization={removeFromOrganization}
        />
      )}

      {activeTab === 'roles' && (
        <RolesTab roles={roles} onDataChange={loadData} />
      )}

      {activeTab === 'organizations' && (
        <OrganizationsTab organizations={organizations} onDataChange={loadData} />
      )}
    </div>
  );
};

// Users Tab Component
const UsersTab = ({ users, roles, organizations, onAssignRole, onRemoveRole, onAddToOrganization, onRemoveFromOrganization }) => {
  return (
    <div>
      <h2>Users ({users.length})</h2>
      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ backgroundColor: '#f8f9fa' }}>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Name</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Email</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Roles</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Organizations</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Actions</th>
            </tr>
          </thead>
          <tbody>
            {users.map(user => (
              <UserRow 
                key={user.id} 
                user={user} 
                roles={roles}
                organizations={organizations}
                onAssignRole={onAssignRole}
                onRemoveRole={onRemoveRole}
                onAddToOrganization={onAddToOrganization}
                onRemoveFromOrganization={onRemoveFromOrganization}
              />
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

// User Row Component
const UserRow = ({ user, roles, organizations, onAssignRole, onRemoveRole, onAddToOrganization, onRemoveFromOrganization }) => {
  const [showRoleAssign, setShowRoleAssign] = useState(false);
  const [showOrgAssign, setShowOrgAssign] = useState(false);

  return (
    <tr>
      <td style={{ padding: '10px', border: '1px solid #ddd' }}>
        <div style={{ display: 'flex', alignItems: 'center', gap: '10px' }}>
          {user.picture && (
            <img 
              src={user.picture} 
              alt={user.name}
              style={{ width: '30px', height: '30px', borderRadius: '50%' }}
            />
          )}
          {user.name}
        </div>
      </td>
      <td style={{ padding: '10px', border: '1px solid #ddd' }}>{user.email}</td>
      <td style={{ padding: '10px', border: '1px solid #ddd' }}>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '5px' }}>
          {user.roles?.map(role => (
            <span 
              key={role.id}
              style={{
                backgroundColor: '#007bff',
                color: 'white',
                padding: '2px 6px',
                borderRadius: '3px',
                fontSize: '12px',
                display: 'flex',
                alignItems: 'center',
                gap: '5px'
              }}
            >
              {role.name}
              <button
                onClick={() => onRemoveRole(user.id, role.id)}
                style={{
                  background: 'none',
                  border: 'none',
                  color: 'white',
                  cursor: 'pointer',
                  fontSize: '12px',
                  padding: '0'
                }}
              >
                ×
              </button>
            </span>
          ))}
          <button
            onClick={() => setShowRoleAssign(!showRoleAssign)}
            style={{
              padding: '2px 6px',
              backgroundColor: '#28a745',
              color: 'white',
              border: 'none',
              borderRadius: '3px',
              fontSize: '12px',
              cursor: 'pointer'
            }}
          >
            +
          </button>
        </div>
        {showRoleAssign && (
          <div style={{ marginTop: '5px' }}>
            <select
              onChange={(e) => {
                if (e.target.value) {
                  onAssignRole(user.id, parseInt(e.target.value));
                  setShowRoleAssign(false);
                  e.target.value = '';
                }
              }}
              style={{ fontSize: '12px', padding: '2px' }}
            >
              <option value="">Select role to assign...</option>
              {roles.filter(role => !user.roles?.some(userRole => userRole.id === role.id)).map(role => (
                <option key={role.id} value={role.id}>{role.name}</option>
              ))}
            </select>
          </div>
        )}
      </td>
      <td style={{ padding: '10px', border: '1px solid #ddd' }}>
        <div style={{ display: 'flex', flexWrap: 'wrap', gap: '5px' }}>
          {user.organizations?.map(org => (
            <span 
              key={org.id}
              style={{
                backgroundColor: '#28a745',
                color: 'white',
                padding: '2px 6px',
                borderRadius: '3px',
                fontSize: '12px',
                display: 'flex',
                alignItems: 'center',
                gap: '5px'
              }}
            >
              {org.name}
              <button
                onClick={() => onRemoveFromOrganization(user.id, org.id)}
                style={{
                  background: 'none',
                  border: 'none',
                  color: 'white',
                  cursor: 'pointer',
                  fontSize: '12px',
                  padding: '0'
                }}
              >
                ×
              </button>
            </span>
          ))}
          <button
            onClick={() => setShowOrgAssign(!showOrgAssign)}
            style={{
              padding: '2px 6px',
              backgroundColor: '#17a2b8',
              color: 'white',
              border: 'none',
              borderRadius: '3px',
              fontSize: '12px',
              cursor: 'pointer'
            }}
          >
            +
          </button>
        </div>
        {showOrgAssign && (
          <div style={{ marginTop: '5px' }}>
            <select
              onChange={(e) => {
                if (e.target.value) {
                  onAddToOrganization(user.id, parseInt(e.target.value));
                  setShowOrgAssign(false);
                  e.target.value = '';
                }
              }}
              style={{ fontSize: '12px', padding: '2px' }}
            >
              <option value="">Select organization to join...</option>
              {organizations.filter(org => !user.organizations?.some(userOrg => userOrg.id === org.id)).map(org => (
                <option key={org.id} value={org.id}>{org.name}</option>
              ))}
            </select>
          </div>
        )}
      </td>
      <td style={{ padding: '10px', border: '1px solid #ddd' }}>
        <span style={{ fontSize: '12px', color: '#666' }}>
          Last login: {user.last_login_at ? new Date(user.last_login_at).toLocaleDateString() : 'Never'}
        </span>
      </td>
    </tr>
  );
};

// Roles Tab Component
const RolesTab = ({ roles, onDataChange }) => {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newRole, setNewRole] = useState({ name: '', description: '' });
  const { authenticatedFetch } = useAuth();

  const createRole = async (e) => {
    e.preventDefault();
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/roles`, {
        method: 'POST',
        body: JSON.stringify(newRole)
      });

      if (response.ok) {
        setNewRole({ name: '', description: '' });
        setShowCreateForm(false);
        onDataChange();
      } else {
        const errorText = await response.text();
        alert('Failed to create role: ' + errorText);
      }
    } catch (err) {
      alert('Failed to create role: ' + err.message);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h2>Roles ({roles.length})</h2>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          style={{
            padding: '8px 16px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Create Role
        </button>
      </div>

      {showCreateForm && (
        <form onSubmit={createRole} style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '4px' }}>
          <h3>Create New Role</h3>
          <div style={{ marginBottom: '10px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Name:</label>
            <input
              type="text"
              value={newRole.name}
              onChange={(e) => setNewRole({...newRole, name: e.target.value})}
              required
              style={{ width: '100%', padding: '8px', border: '1px solid #ddd', borderRadius: '4px' }}
            />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Description:</label>
            <textarea
              value={newRole.description}
              onChange={(e) => setNewRole({...newRole, description: e.target.value})}
              style={{ width: '100%', padding: '8px', border: '1px solid #ddd', borderRadius: '4px', height: '80px' }}
            />
          </div>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button
              type="submit"
              style={{
                padding: '8px 16px',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Create
            </button>
            <button
              type="button"
              onClick={() => setShowCreateForm(false)}
              style={{
                padding: '8px 16px',
                backgroundColor: '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ backgroundColor: '#f8f9fa' }}>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Name</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Description</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Created</th>
            </tr>
          </thead>
          <tbody>
            {roles.map(role => (
              <tr key={role.id}>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>
                  <strong>{role.name}</strong>
                </td>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>{role.description}</td>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>
                  {new Date(role.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

// Organizations Tab Component
const OrganizationsTab = ({ organizations, onDataChange }) => {
  const [showCreateForm, setShowCreateForm] = useState(false);
  const [newOrg, setNewOrg] = useState({ name: '', description: '', metadata: {} });
  const { authenticatedFetch } = useAuth();

  const createOrganization = async (e) => {
    e.preventDefault();
    try {
      const response = await authenticatedFetch(`${config.apiBaseUrl}/api/organizations`, {
        method: 'POST',
        body: JSON.stringify(newOrg)
      });

      if (response.ok) {
        setNewOrg({ name: '', description: '', metadata: {} });
        setShowCreateForm(false);
        onDataChange();
      } else {
        const errorText = await response.text();
        alert('Failed to create organization: ' + errorText);
      }
    } catch (err) {
      alert('Failed to create organization: ' + err.message);
    }
  };

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '20px' }}>
        <h2>Organizations ({organizations.length})</h2>
        <button
          onClick={() => setShowCreateForm(!showCreateForm)}
          style={{
            padding: '8px 16px',
            backgroundColor: '#28a745',
            color: 'white',
            border: 'none',
            borderRadius: '4px',
            cursor: 'pointer'
          }}
        >
          Create Organization
        </button>
      </div>

      {showCreateForm && (
        <form onSubmit={createOrganization} style={{ marginBottom: '20px', padding: '15px', border: '1px solid #ddd', borderRadius: '4px' }}>
          <h3>Create New Organization</h3>
          <div style={{ marginBottom: '10px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Name:</label>
            <input
              type="text"
              value={newOrg.name}
              onChange={(e) => setNewOrg({...newOrg, name: e.target.value})}
              required
              style={{ width: '100%', padding: '8px', border: '1px solid #ddd', borderRadius: '4px' }}
            />
          </div>
          <div style={{ marginBottom: '10px' }}>
            <label style={{ display: 'block', marginBottom: '5px' }}>Description:</label>
            <textarea
              value={newOrg.description}
              onChange={(e) => setNewOrg({...newOrg, description: e.target.value})}
              style={{ width: '100%', padding: '8px', border: '1px solid #ddd', borderRadius: '4px', height: '80px' }}
            />
          </div>
          <div style={{ display: 'flex', gap: '10px' }}>
            <button
              type="submit"
              style={{
                padding: '8px 16px',
                backgroundColor: '#007bff',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Create
            </button>
            <button
              type="button"
              onClick={() => setShowCreateForm(false)}
              style={{
                padding: '8px 16px',
                backgroundColor: '#6c757d',
                color: 'white',
                border: 'none',
                borderRadius: '4px',
                cursor: 'pointer'
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <div style={{ overflowX: 'auto' }}>
        <table style={{ width: '100%', borderCollapse: 'collapse' }}>
          <thead>
            <tr style={{ backgroundColor: '#f8f9fa' }}>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Name</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Description</th>
              <th style={{ padding: '10px', border: '1px solid #ddd', textAlign: 'left' }}>Created</th>
            </tr>
          </thead>
          <tbody>
            {organizations.map(org => (
              <tr key={org.id}>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>
                  <strong>{org.name}</strong>
                </td>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>{org.description}</td>
                <td style={{ padding: '10px', border: '1px solid #ddd' }}>
                  {new Date(org.created_at).toLocaleDateString()}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default AdminDashboard;