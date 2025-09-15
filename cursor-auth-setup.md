# Cursor AI Authentication Setup

This guide helps you authenticate Cursor with your app for better AI development flow.

## 🔑 Step 1: Generate Development Token

Run this command to generate a 30-day development token:

```bash
curl -X POST http://localhost:8080/api/setup/dev-token
```

**Example Response:**
```json
{
  "expires_days": 30,
  "message": "Development token generated successfully", 
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "usage": "Add this as 'Authorization: Bearer <token>' header in your requests",
  "user_email": "your-email@gmail.com",
  "user_id": 1
}
```

## 🔧 Step 2: Configure Cursor

### Option A: Use Cursor's HTTP Client
1. **Open Cursor's HTTP Client** (if available)
2. **Add this header to all requests:**
   ```
   Authorization: Bearer YOUR_TOKEN_HERE
   ```

### Option B: Create a `.cursorrules` File
Create a `.cursorrules` file in your project root:

```
# Authentication for API requests
When making HTTP requests to the backend API, always include:
Authorization: Bearer YOUR_TOKEN_HERE

# API Base URL
Backend API: http://localhost:8080
Frontend: http://localhost:3000

# Common endpoints:
- Health: GET http://localhost:8080/health
- Auth me: GET http://localhost:8080/api/auth/me
- Stripe plans: GET http://localhost:8080/api/stripe/plans
- Messages: GET http://localhost:8080/api/messages
```

### Option C: Environment Variable
Add to your shell profile (`.bashrc`, `.zshrc`, etc.):

```bash
export CURSOR_DEV_TOKEN="YOUR_TOKEN_HERE"
export BACKEND_URL="http://localhost:8080"
```

## 🧪 Step 3: Test Authentication

Test that Cursor can authenticate with your API:

```bash
curl -H "Authorization: Bearer YOUR_TOKEN_HERE" http://localhost:8080/api/auth/me
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "your-email@gmail.com",
    "name": "Your Name",
    "roles": ["user"]
  }
}
```

## 🔄 Token Renewal

The development token lasts 30 days. To generate a new one:

```bash
curl -X POST http://localhost:8080/api/setup/dev-token
```

## 🛡️ Security Notes

- **Development Only**: This endpoint is automatically disabled in production
- **Environment Protected**: Only works when `ENVIRONMENT=development` in .env
- **Don't Commit**: Never commit tokens to version control
- **Regenerate**: Generate new tokens regularly
- **Production**: Use proper OAuth flow in production

### Environment Protection
The `/api/setup/dev-token` endpoint is protected by environment checking:
- **Development**: `ENVIRONMENT=development` - Endpoint available
- **Production**: `ENVIRONMENT=production` - Endpoint returns 403 Forbidden

## 📝 Usage Examples for Cursor

Now Cursor can help you with authenticated requests:

**Example prompts:**
- "Test the Stripe checkout endpoint"
- "Get my current user info from the API"
- "Create a new message via the API"
- "Check my subscription status"

Cursor will automatically include the authentication header when making requests to your backend. 