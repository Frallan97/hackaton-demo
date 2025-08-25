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

// Create the base API
export const api = createApi({
  reducerPath: 'api',
  baseQuery,
  tagTypes: ['User', 'Role', 'Organization', 'Message', 'Auth'],
  endpoints: () => ({}),
});