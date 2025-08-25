# React Go App - Hackathon Template Instructions

Always follow these instructions first and only search for additional information or use bash commands when you encounter unexpected information that does not match the instructions here.

## Working Effectively

### Prerequisites Installation
Install required tools in this exact order:
- `curl -fsSL https://bun.sh/install | bash` -- installs Bun (JavaScript runtime)
- `source ~/.bashrc` -- reload shell to use Bun
- Verify Go is available: `go version` -- should show Go 1.23+
- Verify Docker is available: `docker --version && docker compose version`

### Bootstrap, Build, and Test the Repository
Follow these steps in exact order. NEVER CANCEL any long-running commands:

1. **Initial Setup** (takes 10-15 seconds):
   ```bash
   cp env.example .env
   docker compose up db redis nats -d
   ```
   - Database services startup: PostgreSQL, Redis, NATS
   - Wait 10 seconds for services to be ready

2. **Backend Setup** (takes 20-30 seconds):
   ```bash
   cd backend
   go mod tidy
   ```
   - Downloads Go dependencies. NEVER CANCEL - wait for completion.
   - Set timeout to 120+ seconds for this command.

3. **Frontend Setup** (takes 20-30 seconds):
   ```bash
   cd frontend
   bun install
   ```
   - Alternative: `npm install` (takes 60-90 seconds)
   - Downloads frontend dependencies. NEVER CANCEL - wait for completion.
   - Set timeout to 120+ seconds for this command.

4. **Build Frontend** (takes 5-10 seconds):
   ```bash
   cd frontend
   bun run build
   ```
   - Alternative: `npm run build`
   - Frontend builds quickly. Set timeout to 60+ seconds.

5. **Test Backend** (takes 5-15 seconds):
   ```bash
   cd backend
   go test ./...
   ```
   - Most packages have no tests currently
   - Set timeout to 60+ seconds.

### Run the Application

#### Development Mode (Recommended)
Start services and applications separately for development:

1. **Start Services** (if not already running):
   ```bash
   docker compose up db redis nats -d
   ```

2. **Start Backend** (takes 5-10 seconds to start):
   ```bash
   cd backend
   go run main.go
   ```
   - Backend starts on :8080
   - NEVER CANCEL - wait for "listening on :8080" message
   - Set timeout to 120+ seconds for first run (includes dependency compilation)

3. **Start Frontend** (takes <1 second to start):
   ```bash
   cd frontend
   bun run dev
   ```
   - Alternative: `npm run dev`
   - Frontend starts on :3000
   - Starts immediately

#### Access Points
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- API Documentation: http://localhost:8080/docs
- NATS Monitoring: http://localhost:8222

## Validation

### ALWAYS Run These Validation Steps After Making Changes

1. **Health Check**:
   ```bash
   curl -s http://localhost:8080/health
   ```
   Should return `{"success":true,"data":{"status":"healthy",...}}`

2. **API Endpoints**:
   ```bash
   curl -s http://localhost:8080/api/auth/google/url
   curl -s http://localhost:8080/docs
   ```

3. **NATS Functionality**:
   ```bash
   cd backend/cmd/test-nats
   go run main.go
   ```
   Should output "NATS test completed successfully!"

4. **Database Connection**: Backend logs should show:
   - "Database connection established successfully"
   - "NATS Event Bus initialized successfully"
   - "listening on :8080"

### Manual Testing Scenarios
After making changes, ALWAYS test these scenarios:

1. **Basic API Functionality**:
   - Health endpoint responds with 200 OK
   - OAuth URL generation works
   - Documentation is accessible

2. **Services Integration**:
   - Database connection is stable
   - NATS messaging works (run test-nats)
   - Redis connection (if using caching features)

3. **Frontend-Backend Communication**:
   - Frontend can reach backend API
   - CORS is properly configured
   - No console errors in browser

## Known Issues and Limitations

### Docker Compose Full Build
- `docker compose up --build` -- FAILS due to SSL certificate issues in CI environments
- DO NOT use this command in automated environments
- Use development mode instead (services only + manual frontend/backend)

### Environment Configuration
- Default `.env` file has placeholder OAuth credentials
- For full OAuth functionality, update `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET`
- Backend works without real OAuth credentials (uses defaults)

### Testing Infrastructure
- Backend: Basic test structure exists (`go test ./...` works)
- Frontend: No tests currently configured
- Integration tests: Available via NATS test command

## Build Times and Timeouts

Critical: ALWAYS set appropriate timeouts and NEVER CANCEL builds:

- **Service startup**: 15 seconds (set timeout: 60+ seconds)
- **go mod tidy**: 20-30 seconds (set timeout: 120+ seconds)
- **bun install**: 20-30 seconds (set timeout: 120+ seconds)
- **npm install**: 60-90 seconds (set timeout: 180+ seconds)
- **Frontend build**: 5-10 seconds (set timeout: 60+ seconds)
- **Backend startup**: 5-10 seconds, first run 30+ seconds (set timeout: 120+ seconds)
- **go test**: 5-15 seconds (set timeout: 60+ seconds)

## Common Tasks

### Repository Structure
```
hackaton-demo/
├── .github/                    # GitHub workflows and this file
├── backend/                    # Go API server
│   ├── cmd/test-nats/         # NATS functionality test
│   ├── config/                # Configuration management
│   ├── controllers/           # HTTP handlers
│   ├── database/              # Database connection
│   ├── events/                # NATS event system
│   ├── handlers/              # Route management
│   ├── middleware/            # HTTP middleware
│   ├── models/                # Data structures
│   ├── services/              # Business logic
│   └── main.go                # Application entry point
├── frontend/                  # React + TypeScript + Vite
│   ├── src/                   # Source code
│   ├── public/                # Static assets
│   └── package.json           # Dependencies
├── charts/                    # Helm charts for Kubernetes
├── scripts/                   # Setup scripts
├── docker-compose.yml         # Service definitions
└── env.example                # Environment template
```

### Quick Command Reference
```bash
# Setup
cp env.example .env
docker compose up db redis nats -d

# Backend
cd backend && go mod tidy && go run main.go

# Frontend  
cd frontend && bun install && bun run dev

# Tests
cd backend && go test ./...
cd backend/cmd/test-nats && go run main.go

# Health checks
curl http://localhost:8080/health
curl http://localhost:8222/healthz
```

### Environment Variables (from .env)
```bash
# Database (using Docker services)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres

# Server
SERVER_PORT=8080

# Security (change in production)
JWT_SECRET_KEY=your-secret-key-change-in-production

# Google OAuth (optional for development)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:3000/auth/callback
```

### Database and Services
- **PostgreSQL**: Runs on localhost:5432, auto-migrates on startup
- **Redis**: Runs on localhost:6379, ready for caching
- **NATS**: Runs on localhost:4222, with monitoring on :8222
- **Migrations**: Run automatically when backend starts
- **Event System**: NATS JetStream for reliable message processing

### Architecture Overview
- **Backend**: Modular Go architecture with RBAC, OAuth, events
- **Frontend**: React + Redux Toolkit + TypeScript + Tailwind CSS
- **Database**: PostgreSQL with automatic migrations
- **Messaging**: NATS JetStream for event-driven architecture
- **Development**: Docker Compose for services, local dev for apps

Always validate your changes using the health checks and test commands above before considering your work complete.