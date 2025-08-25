import { api } from '../../../shared/store/api';

export const dashboardApi = api.injectEndpoints({
  endpoints: (builder) => ({
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

    // Health check
    getHealthCheck: builder.query<any, void>({
      query: () => '/health',
    }),
  }),
});

export const {
  useGetMessagesQuery,
  useCreateMessageMutation,
  useGetHealthCheckQuery,
} = dashboardApi;