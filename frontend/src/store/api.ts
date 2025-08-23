import { createApi, fetchBaseQuery } from '@reduxjs/toolkit/query/react';
import config from '../config';

// Create base query with authentication and response transformation
const baseQuery = fetchBaseQuery({
  baseUrl: config.apiBaseUrl,
  prepareHeaders: (headers, { getState }) => {
    // Get the token from the auth state
    const token = (getState() as any).auth.accessToken;
    
    if (token) {
      headers.set('authorization', `Bearer ${token}`);
    }
    
    headers.set('Content-Type', 'application/json');
    return headers;
  },
  // Transform all responses to handle the standardized format
  responseHandler: async (response) => {
    const data = await response.json();
    
    // If the response has our standardized format, extract the data
    if (data && typeof data === 'object' && 'success' in data) {
      if (data.success && data.data !== undefined) {
        return data.data;
      } else {
        // Handle error responses
        throw new Error(data.error || 'Request failed');
      }
    }
    
    // Return original response if not in standardized format
    return data;
  },
});

// Create the API
export const api = createApi({
  reducerPath: 'api',
  baseQuery,
  tagTypes: ['User', 'Role', 'Organization', 'Message', 'Auth'],
  endpoints: (builder) => ({
    // Auth endpoints
    getGoogleAuthUrl: builder.query<{ auth_url: string; state: string }, void>({
      query: () => '/api/auth/google/url',
      providesTags: ['Auth'],
    }),

    googleLogin: builder.mutation<{
      user: any;
      access_token: string;
      refresh_token: string;
      token_type: string;
      expires_in: number;
    }, { code: string }>({
      query: (credentials) => ({
        url: '/api/auth/google/login',
        method: 'POST',
        body: credentials,
      }),
      invalidatesTags: ['Auth', 'User'],
    }),

    refreshToken: builder.mutation<{
      access_token: string;
      token_type: string;
      expires_in: string;
    }, { refresh_token: string }>({
      query: (credentials) => ({
        url: '/api/auth/refresh',
        method: 'POST',
        body: credentials,
      }),
      invalidatesTags: ['Auth'],
    }),

    getCurrentUser: builder.query<any, void>({
      query: () => '/api/auth/me',
      providesTags: ['User'],
    }),

    logout: builder.mutation<void, void>({
      query: () => ({
        url: '/api/auth/logout',
        method: 'POST',
      }),
      invalidatesTags: ['Auth', 'User'],
    }),

    // User management endpoints
    getUsers: builder.query<any[], void>({
      query: () => '/api/admin/users',
      providesTags: ['User'],
    }),

    assignRole: builder.mutation<void, { user_id: number; role_id: number }>({
      query: (data) => ({
        url: '/api/admin/assign-role',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['User'],
    }),

    removeRole: builder.mutation<void, { user_id: number; role_id: number }>({
      query: (data) => ({
        url: '/api/admin/remove-role',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['User'],
    }),

    assignOrganization: builder.mutation<void, { user_id: number; organization_id: number }>({
      query: (data) => ({
        url: '/api/admin/assign-organization',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['User'],
    }),

    removeOrganization: builder.mutation<void, { user_id: number; organization_id: number }>({
      query: (data) => ({
        url: '/api/admin/remove-organization',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['User'],
    }),

    // Role management endpoints
    getRoles: builder.query<any[], void>({
      query: () => '/api/roles',
      providesTags: ['Role'],
    }),

    createRole: builder.mutation<any, { name: string; description?: string }>({
      query: (data) => ({
        url: '/api/roles',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Role'],
    }),

    updateRole: builder.mutation<any, { id: number; name: string; description?: string }>({
      query: ({ id, ...data }) => ({
        url: `/api/roles?id=${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Role'],
    }),

    deleteRole: builder.mutation<void, { id: number }>({
      query: ({ id }) => ({
        url: `/api/roles?id=${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Role'],
    }),

    // Organization management endpoints
    getOrganizations: builder.query<any[], void>({
      query: () => '/api/organizations',
      providesTags: ['Organization'],
    }),

    createOrganization: builder.mutation<any, { name: string; description?: string }>({
      query: (data) => ({
        url: '/api/organizations',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Organization'],
    }),

    updateOrganization: builder.mutation<any, { id: number; name: string; description?: string }>({
      query: ({ id, ...data }) => ({
        url: `/api/organizations?id=${id}`,
        method: 'PUT',
        body: data,
      }),
      invalidatesTags: ['Organization'],
    }),

    deleteOrganization: builder.mutation<void, { id: number }>({
      query: ({ id }) => ({
        url: `/api/organizations?id=${id}`,
        method: 'DELETE',
      }),
      invalidatesTags: ['Organization'],
    }),

    // Message endpoints
    getMessages: builder.query<any[], void>({
      query: () => '/api/messages',
      providesTags: ['Message'],
    }),

    createMessage: builder.mutation<{ id: number }, { content: string }>({
      query: (data) => ({
        url: '/api/messages',
        method: 'POST',
        body: data,
      }),
      invalidatesTags: ['Message'],
    }),

    // Setup endpoints
    setupFirstAdmin: builder.mutation<any, void>({
      query: () => ({
        url: '/api/setup/first-admin',
        method: 'POST',
      }),
      invalidatesTags: ['User', 'Role'],
    }),

    // Health check
    getHealth: builder.query<any, void>({
      query: () => '/health',
      providesTags: ['Auth'],
    }),
  }),
});

// Export hooks for use in components
export const {
  // Auth hooks
  useGetGoogleAuthUrlQuery,
  useGoogleLoginMutation,
  useRefreshTokenMutation,
  useGetCurrentUserQuery,
  useLogoutMutation,
  
  // User management hooks
  useGetUsersQuery,
  useAssignRoleMutation,
  useRemoveRoleMutation,
  useAssignOrganizationMutation,
  useRemoveOrganizationMutation,
  
  // Role management hooks
  useGetRolesQuery,
  useCreateRoleMutation,
  useUpdateRoleMutation,
  useDeleteRoleMutation,
  
  // Organization management hooks
  useGetOrganizationsQuery,
  useCreateOrganizationMutation,
  useUpdateOrganizationMutation,
  useDeleteOrganizationMutation,
  
  // Message hooks
  useGetMessagesQuery,
  useCreateMessageMutation,
  
  // Setup hooks
  useSetupFirstAdminMutation,
  
  // Health hooks
  useGetHealthQuery,
} = api; 