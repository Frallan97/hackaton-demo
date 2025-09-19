# Backend API

A Go-based REST API with a clean, modular architecture and optimized startup performance.

## Architecture

```
backend/
├── controllers/          # HTTP request handlers
│   ├── health_controller.go
│   ├── message_controller.go
│   └── auth_controller.go
├── models/              # Data structures
│   ├── message.go
│   └── auth.go
├── handlers/            # Route management
│   └── router.go
├── middleware/          # HTTP middleware
│   └── logging.go
├── config/             # Configuration management
│   └── config.go
├── database/           # Database connection
│   └── connection.go
├── migrations/         # Database migrations
├── .air.toml          # Air configuration for hot reloading
└── main.go            # Application entry point
```

## Features

- **Modular Architecture**: Clean separation of concerns
- **Optimized Startup**: Fast database connection with exponential backoff
- **Hot Reloading**: Air integration for development
- **Middleware Support**: Logging and CORS middleware included
- **Configuration Management**: Environment-based configuration with singleton pattern
- **Database Monitoring**: Automatic connection health monitoring with adaptive intervals
- **Swagger Documentation**: Auto-generated API docs at `/docs`
- **Health Checks**: Database connectivity monitoring

## Performance Optimizations

- **Faster Startup**: Reduced database connection timeout from 30s to 15s
- **Exponential Backoff**: Smart connection retry strategy
- **Reduced Logging**: Debug-mode conditional logging to reduce noise
- **Optimized Connection Pool**: Tuned for faster startup and better resource usage
- **Singleton Configuration**: Configuration loaded once and cached
- **Adaptive Monitoring**: Different monitoring intervals for dev vs production

## API Endpoints

### Health
- `GET /health` - Health check with database status

### Messages
- `GET /api/messages` - List all messages
- `POST /api/messages` - Create a new message

### Authentication (Google OAuth)
- `GET /api/auth/google/url` - Get Google OAuth authorization URL
- `POST /api/auth/google/login` - Complete Google OAuth login with authorization code
- `GET /api/auth/me` - Get current user information (requires Bearer token)
- `POST /api/auth/refresh` - Refresh access token using refresh token
- `POST /api/auth/logout` - Logout (client-side token removal)

### Documentation
- `GET /docs` - Swagger UI documentation

## Development

### Prerequisites
- Go 1.23+
- PostgreSQL (for full functionality)
- Air (for hot reloading)

### Quick Start
```bash
# Install dependencies
go mod tidy

# Install Air for hot reloading
go install github.com/air-verse/air@latest

# Development with hot reloading (recommended)
DEBUG=true ENVIRONMENT=development air

# Or run normally without hot reloading
go run main.go
```

### Air Hot Reloading

Air is configured to provide the best development experience:

```bash
# Start with debug logging and hot reloading
DEBUG=true ENVIRONMENT=development air

# Start in production mode (less verbose logging)
ENVIRONMENT=production air

# Just run with Air (uses defaults)
air
```

The `.air.toml` configuration provides:
- **Fast rebuilds**: 1-second delay for file changes
- **Smart exclusions**: Ignores test files, migrations, and build artifacts
- **Clean interface**: Colored output and clear screen on rebuild
- **Error logging**: Build errors saved to `build-errors.log`

### Environment Variables

#### Core Configuration
- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: postgres)
- `DB_NAME` - Database name (default: postgres)
- `DB_URL` - Full database URL (overrides individual vars)
- `SERVER_PORT` - Server port (default: 8080)

#### Development & Debugging
- `DEBUG` - Enable debug logging (set to "true" for verbose logs)
- `ENVIRONMENT` - Set to "production" for optimized monitoring intervals

#### Authentication (Google OAuth)
- `JWT_SECRET_KEY` - Secret key for JWT token signing (default: your-secret-key-change-in-production)
- `GOOGLE_CLIENT_ID` - Google OAuth client ID (required for authentication)
- `GOOGLE_CLIENT_SECRET` - Google OAuth client secret (required for authentication)
- `GOOGLE_REDIRECT_URL` - OAuth redirect URL (default: http://localhost:3000/auth/callback)

## Adding New Features

### 1. Create a Model
```go
// models/example.go
package models

type Example struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

### 2. Create a Controller
```go
// controllers/example_controller.go
package controllers

import (
    "net/http"
    "github.com/frallan97/hackaton-demo-backend/database"
)

type ExampleController struct {
    dbManager *database.DBManager
}

func NewExampleController(dbManager *database.DBManager) *ExampleController {
    return &ExampleController{dbManager: dbManager}
}

func (ec *ExampleController) ExampleHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handle request
    }
}
```

### 3. Add to Router
```go
// handlers/router.go
func (r *Router) SetupRoutes() http.Handler {
    mux := http.NewServeMux()
    
    // Add your new endpoint
    mux.HandleFunc("/api/example", r.exampleController.ExampleHandler())
    
    // Apply middleware
    handler := middleware.LoggingMiddleware(mux)
    handler = middleware.CORSMiddleware(handler)
    
    return handler
}
```

## Testing

The application includes comprehensive logging and health checks. The database connection is monitored continuously and will log connection status changes.

## Deployment

The application is containerized and can be deployed using Docker or Kubernetes. See the root README for deployment instructions. 