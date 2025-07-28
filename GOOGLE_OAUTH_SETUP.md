# Google OAuth Setup Guide

This guide will help you set up Google OAuth for the backend authentication system.

## Prerequisites

1. A Google account
2. Access to Google Cloud Console

## Step 1: Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select an existing one
3. Enable the Google+ API (if not already enabled)

## Step 2: Configure OAuth Consent Screen

1. In the Google Cloud Console, go to **APIs & Services** > **OAuth consent screen**
2. Choose **External** user type (unless you have a Google Workspace organization)
3. Fill in the required information:
   - **App name**: Your application name
   - **User support email**: Your email address
   - **Developer contact information**: Your email address
4. Add the following scopes:
   - `https://www.googleapis.com/auth/userinfo.email`
   - `https://www.googleapis.com/auth/userinfo.profile`
5. Add test users (your email address) if you're in testing mode
6. Save and continue

## Step 3: Create OAuth 2.0 Credentials

1. Go to **APIs & Services** > **Credentials**
2. Click **Create Credentials** > **OAuth 2.0 Client IDs**
3. Choose **Web application** as the application type
4. Fill in the details:
   - **Name**: Your application name
   - **Authorized JavaScript origins**:
     - `http://localhost:3000` (for development)
     - `https://yourdomain.com` (for production)
   - **Authorized redirect URIs**:
     - `http://localhost:3000/auth/callback` (for development)
     - `https://yourdomain.com/auth/callback` (for production)
5. Click **Create**

## Step 4: Get Your Credentials

After creating the OAuth 2.0 client, you'll get:
- **Client ID**: A long string ending with `.apps.googleusercontent.com`
- **Client Secret**: A secret string (keep this secure!)

## Step 5: Configure Environment Variables

Set these environment variables in your backend:

```bash
# Required for authentication
export GOOGLE_CLIENT_ID="your-client-id.apps.googleusercontent.com"
export GOOGLE_CLIENT_SECRET="your-client-secret"

# Optional (defaults shown)
export GOOGLE_REDIRECT_URL="http://localhost:3000/auth/callback"
export JWT_SECRET_KEY="your-secret-key-change-in-production"
```

## Step 6: Frontend Integration

Your frontend needs to:

1. **Get the OAuth URL**:
   ```javascript
   const response = await fetch('/api/auth/google/url');
   const { auth_url } = await response.json();
   window.location.href = auth_url;
   ```

2. **Handle the callback**:
   ```javascript
   // Extract the authorization code from URL parameters
   const urlParams = new URLSearchParams(window.location.search);
   const code = urlParams.get('code');
   
   if (code) {
     const response = await fetch('/api/auth/google/login', {
       method: 'POST',
       headers: { 'Content-Type': 'application/json' },
       body: JSON.stringify({ code })
     });
     
     const authData = await response.json();
     // Store tokens securely (e.g., in httpOnly cookies or secure storage)
     localStorage.setItem('access_token', authData.access_token);
     localStorage.setItem('refresh_token', authData.refresh_token);
   }
   ```

3. **Use the access token**:
   ```javascript
   const response = await fetch('/api/auth/me', {
     headers: {
       'Authorization': `Bearer ${localStorage.getItem('access_token')}`
     }
   });
   const user = await response.json();
   ```

## Testing the Setup

1. **Start the backend**:
   ```bash
   cd backend
   go run main.go
   ```

2. **Test the OAuth URL endpoint**:
   ```bash
   curl http://localhost:8080/api/auth/google/url
   ```

3. **Test the health endpoint**:
   ```bash
   curl http://localhost:8080/health
   ```

## Security Considerations

1. **Never expose your client secret** in frontend code
2. **Use HTTPS in production**
3. **Store JWT tokens securely** (httpOnly cookies recommended)
4. **Implement proper token refresh logic**
5. **Add rate limiting** for authentication endpoints
6. **Use environment variables** for all sensitive configuration

## Troubleshooting

### Common Issues

1. **"redirect_uri_mismatch" error**:
   - Ensure the redirect URI in Google Console matches exactly
   - Check for trailing slashes or protocol differences

2. **"invalid_client" error**:
   - Verify your client ID and secret are correct
   - Ensure the OAuth consent screen is properly configured

3. **"access_denied" error**:
   - Check if the user is in the test users list (if in testing mode)
   - Verify the required scopes are added

4. **CORS errors**:
   - Ensure the backend CORS middleware is configured correctly
   - Check that the frontend origin is allowed

### Debug Mode

To enable debug logging, set the log level in your backend:

```go
log.SetLevel(log.DebugLevel)
```

## Production Deployment

For production deployment:

1. **Update redirect URIs** in Google Console to use your production domain
2. **Use a strong JWT secret key**
3. **Enable HTTPS** for all OAuth endpoints
4. **Implement proper session management**
5. **Add monitoring and logging** for authentication events
6. **Consider implementing token blacklisting** for enhanced security

## Support

If you encounter issues:

1. Check the Google Cloud Console for any error messages
2. Review the backend logs for detailed error information
3. Verify all environment variables are set correctly
4. Test with the provided curl commands
5. Check the Swagger documentation at `http://localhost:8080/docs` 