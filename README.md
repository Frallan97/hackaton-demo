# React Go App

A full-stack micro-app with React frontend, Go backend, and PostgreSQL database.

## Structure

```
hackaton-demo/
├── frontend/          # React app (Vite)
├── backend/           # Go API server (modular architecture)
│   ├── controllers/   # HTTP request handlers
│   ├── models/        # Data structures
│   ├── handlers/      # Route management
│   ├── middleware/    # HTTP middleware
│   ├── config/        # Configuration management
│   └── database/      # Database connection management
├── charts/            # Helm chart for Kubernetes
└── .github/           # CI/CD workflows
```

## Local Development

### Prerequisites
- Docker & Docker Compose (for database and Redis)
- Go 1.23+
- Bun (for frontend)

### Quick Start with Docker (All Services)
```bash
# Set up environment variables first
cp env.example .env
# Edit .env with your Google OAuth credentials

# Start everything locally including frontend and backend
docker-compose up --build

# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# API docs: http://localhost:8080/docs
```

### Development Mode (Recommended)
This approach runs frontend and backend locally for faster development, while keeping database and Redis in Docker.

#### 1. Set Up Environment Variables
```bash
# Copy the example environment file
cp env.example .env

# Edit .env with your actual values
# You'll need to get Google OAuth credentials from Google Cloud Console
nano .env  # or use your preferred editor
```

#### 2. Start Database and Redis
```bash
# Start only database and Redis services
docker-compose up db redis -d
```

#### 3. Backend Development
```bash
cd backend

# Install dependencies
go mod tidy

# Run the server (will use environment variables from .env)
go run main.go

# Or build and run
go build -o server .
./server
```

The backend will be available at: http://localhost:8080

#### 4. Frontend Development
```bash
cd frontend

# Install dependencies
bun install

# Start development server
bun run dev
```

The frontend will be available at: http://localhost:3000

### Environment Variables Setup

The application uses environment variables for configuration. Copy `env.example` to `.env` and update the values:

```bash
cp env.example .env
```

**Required for Google OAuth:**
- `GOOGLE_CLIENT_ID` - Your Google OAuth client ID
- `GOOGLE_CLIENT_SECRET` - Your Google OAuth client secret

**Security:**
- `JWT_SECRET_KEY` - Secret key for JWT token signing (change in production)

**Optional (have defaults):**
- `GOOGLE_REDIRECT_URL` - OAuth redirect URL
- `SERVER_PORT` - Backend server port
- Database configuration variables

### Backend Architecture
The backend uses a modular controller-based architecture:
- **Controllers**: Handle HTTP requests and business logic
- **Models**: Define data structures
- **Middleware**: Cross-cutting concerns (logging, CORS)
- **Config**: Environment-based configuration
- **Database**: Connection management and monitoring

See [backend/README.md](backend/README.md) for detailed backend documentation.

### Frontend Features
The React frontend includes:
- **Google OAuth Integration**: Complete OAuth 2.0 flow with Google
- **User Profile Display**: Shows user information from Google
- **Token Management**: Automatic token refresh and logout
- **Responsive Design**: Clean, modern UI with Google branding
- **Error Handling**: Comprehensive error messages and loading states

See [DEMO_OAUTH_FLOW.md](DEMO_OAUTH_FLOW.md) for testing the complete OAuth flow.

### Stopping Services
```bash
# Stop database and Redis
docker-compose down

# Stop all services (if running full docker-compose)
docker-compose down
```

### Troubleshooting
- **Backend can't connect to database**: Ensure `docker-compose up db redis -d` is running
- **Port conflicts**: Check if ports 3000, 8080, 5432, or 6379 are already in use
- **Database connection issues**: Verify PostgreSQL is running with `docker-compose ps`
- **Frontend build issues**: Try `bun install` in the frontend directory
- **OAuth errors**: Check that `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` are set in `.env`
- **JWT errors**: Verify `JWT_SECRET_KEY` is set in `.env`
- **Kubernetes nginx proxy errors**: The frontend uses environment variable `BACKEND_URL` to connect to backend service
  - Docker Compose: `http://backend:8080`
  - Kubernetes: `http://hackaton-demo-backend:8080` (or your release name)

## Deployment

### Kubernetes (via Argo CD)
```bash
# Deploy to cluster
helm install hackaton-demo charts/hackaton-demo/

# Or via Argo CD Application
# (configured to auto-sync from Git)
```

### CI/CD
- GitHub Actions builds and pushes Docker images
- Automatically bumps Helm chart versions
- Argo CD deploys to Kubernetes cluster
- **Swagger Documentation**: Automatically regenerated during CI/CD pipeline
  - Uses `swag` tool to generate API documentation from code annotations
  - Available at `/docs` endpoint when backend is running
  - Updated automatically on every deployment

## API Endpoints

### Health & Documentation
- `GET /health` - Health check
- `GET /docs` - Swagger UI

### Messages
- `GET /api/messages` - List messages
- `POST /api/messages` - Create message

### Authentication (Google OAuth)
- `GET /api/auth/google/url` - Get Google OAuth authorization URL
- `POST /api/auth/google/login` - Complete Google OAuth login with authorization code (rate limited: 5 requests per minute per IP)
- `GET /api/auth/me` - Get current user information (requires Bearer token)
- `POST /api/auth/refresh` - Refresh access token using refresh token
- `POST /api/auth/logout` - Logout (client-side token removal)

## Environment Variables

### Backend
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: postgres)
- `DB_URL` - Full database URL (overrides individual vars)
- `SERVER_PORT` - Server port (default: 8080)

### Authentication (Google OAuth)
- `JWT_SECRET_KEY` - Secret key for JWT token signing (default: your-secret-key-change-in-production)
- `GOOGLE_CLIENT_ID` - Google OAuth client ID (required for authentication)
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret (required for authentication)
- `GOOGLE_REDIRECT_URL` - OAuth redirect URL (default: http://localhost:3000/auth/callback)

### Development Tips
- When running backend locally, it will automatically connect to the PostgreSQL instance running in Docker
- The backend includes automatic database connection monitoring and health checks
- All API endpoints are documented via Swagger at http://localhost:8080/docs
- **Google OAuth Setup**: You'll need to create a Google OAuth application in the Google Cloud Console and set the client ID and secret
- **Production OAuth**: See [PRODUCTION_OAUTH_SETUP.md](PRODUCTION_OAUTH_SETUP.md) for production deployment configuration

