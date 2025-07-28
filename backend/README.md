# Backend API

A Go-based REST API with a clean, modular architecture.

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
└── main.go            # Application entry point
```

## Features

- **Modular Architecture**: Clean separation of concerns
- **Middleware Support**: Logging and CORS middleware included
- **Configuration Management**: Environment-based configuration
- **Database Monitoring**: Automatic connection health monitoring
- **Swagger Documentation**: Auto-generated API docs at `/docs`
- **Health Checks**: Database connectivity monitoring

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

### Running Locally
```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go

# Or build and run
go build -o server .
./server
```

### Environment Variables
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
    "github.com/frallan97/react-go-app-backend/database"
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