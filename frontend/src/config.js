// Frontend Configuration
const config = {
  // API Base URL - will be different for development vs production
  apiBaseUrl: import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080',
  
  // OAuth Configuration
  googleClientId: import.meta.env.VITE_GOOGLE_CLIENT_ID || '',
  
  // App Configuration
  appName: 'Hackaton Demo', // This will be replaced by the template script
  appVersion: '1.0.0',
  
  // Feature Flags
  features: {
    enableAdminPanel: true,
    enableRBAC: true,
    enableOrganizations: true,
  },
  
  // UI Configuration
  ui: {
    theme: 'light', // 'light' | 'dark'
    language: 'en',
    timezone: 'UTC',
  },
  
  // Development Configuration
  dev: {
    enableDebugLogs: import.meta.env.DEV || false,
    mockApi: false,
  }
};

export default config; 