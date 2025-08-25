import { api } from '../../../shared/store/api';

export const authApi = api.injectEndpoints({
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
  }),
});

export const {
  useGetGoogleAuthUrlQuery,
  useGoogleLoginMutation,
  useRefreshTokenMutation,
  useGetCurrentUserQuery,
  useLogoutMutation,
} = authApi;