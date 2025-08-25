import { api } from '../../../shared/store/api';

export const adminApi = api.injectEndpoints({
  endpoints: (builder) => ({
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

    // Setup endpoints
    setupFirstAdmin: builder.mutation<any, void>({
      query: () => ({
        url: '/api/setup/first-admin',
        method: 'POST',
      }),
      invalidatesTags: ['User', 'Role'],
    }),
  }),
});

export const {
  useGetUsersQuery,
  useAssignRoleMutation,
  useRemoveRoleMutation,
  useAssignOrganizationMutation,
  useRemoveOrganizationMutation,
  useGetRolesQuery,
  useCreateRoleMutation,
  useUpdateRoleMutation,
  useDeleteRoleMutation,
  useGetOrganizationsQuery,
  useCreateOrganizationMutation,
  useUpdateOrganizationMutation,
  useDeleteOrganizationMutation,
  useSetupFirstAdminMutation,
} = adminApi;