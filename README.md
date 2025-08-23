# React Go App - Hackathon Template

A full-stack micro-app with React frontend, Go backend, PostgreSQL database, Redis caching, and NATS event messaging. Perfect for hackathons and rapid prototyping.

## 🏗️ Architecture

```
hackaton-demo/
├── frontend/          # React app (Vite + TypeScript)
├── backend/           # Go API server (modular architecture)
│   ├── controllers/   # HTTP request handlers
│   ├── models/        # Data structures
│   ├── handlers/      # Route management
│   ├── middleware/    # HTTP middleware (logging, CORS, RBAC)
│   ├── config/        # Configuration management
│   ├── database/      # Database connection management
│   ├── events/        # NATS event bus system
│   └── services/      # Business logic services
├── charts/            # Helm chart for Kubernetes deployment
└── docker-compose.yml # Local development environment
```

## 🚀 Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.23+
- Bun (for frontend)

### One-Command Setup
```bash
# Clone and setup
git clone <your-repo>
cd hackaton-demo

# Copy environment file
cp env.example .env

# Edit .env with your Google OAuth credentials
nano .env

# Start everything
docker-compose up --build

# Access your app
# Frontend: http://localhost:3000
# Backend: http://localhost:8080
# API docs: http://localhost:8080/docs
```

## 🔧 Environment Setup

### Required Environment Variables
```bash
# Google OAuth (required for authentication)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret

# Security
JWT_SECRET_KEY=your-secret-key-change-in-production

# Optional (have defaults)
GOOGLE_REDIRECT_URL=http://localhost:3000/auth/callback
SERVER_PORT=8080
```

### Google OAuth Setup
1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing
3. Enable Google+ API
4. Configure OAuth consent screen (External user type)
5. Create OAuth 2.0 credentials (Web application)
6. Add authorized origins: `http://localhost:3000`
7. Add redirect URIs: `http://localhost:3000/auth/callback`
8. Copy Client ID and Secret to your `.env` file

## 🏃‍♂️ Development Mode

### Start Only Database Services
```bash
# Start database, Redis, and NATS
docker-compose up db redis nats -d
```

### Backend Development
```bash
cd backend
go mod tidy
go run main.go
```

### Frontend Development
```bash
cd frontend
bun install
bun run dev
```

## 🔐 Authentication & RBAC

### Google OAuth Flow
1. User clicks login → redirected to Google
2. Google returns authorization code
3. Backend exchanges code for user info
4. User created/updated in database
5. JWT tokens generated and returned
6. NATS events published for user actions

### Role-Based Access Control
- **admin**: Full system access, manage users/roles/organizations
- **manager**: View and manage roles/organizations
- **editor**: Edit content (placeholder)
- **reader**: Read-only access (placeholder)

### First Admin Setup
```bash
# After first user login, make them admin
curl -X POST http://localhost:8080/api/setup/first-admin
```

## 📡 NATS Event System

### What is NATS?
- High-performance messaging system (10M+ messages/second)
- JetStream persistence for reliable delivery
- Auto-scaling and clustering support
- Cloud-native design

### Event Types
```go
// User events
EventTypeUserCreated = "user.created"
EventTypeUserLogin   = "user.login"
EventTypeUserLogout  = "user.logout"

// Authentication events
EventTypeAuthSuccess = "auth.success"
EventTypeAuthFailure = "auth.failure"

// Admin events
EventTypeAdminAction = "admin.action"
```

### Event Flow
1. **Publishing**: Events sent to NATS subjects (e.g., `events.users.user.login`)
2. **Persistence**: Stored in JetStream for reliability
3. **Processing**: Handlers run asynchronously in goroutines
4. **Subscription**: Real-time event delivery to subscribers

### Example Usage
```go
// Publish user login event
err := eventService.PublishUserLogin(userID, email, name)

// Subscribe to user events
userEvents, err := eventService.SubscribeToUserEvents(userID)
```

## 🗄️ Database Schema

### Core Tables
- **users**: User accounts with Google OAuth integration
- **roles**: Role definitions (admin, manager, editor, reader)
- **organizations**: Group management
- **user_roles**: User-role assignments
- **user_organizations**: User-organization memberships
- **messages**: Simple message storage

### Migrations
Database migrations run automatically on startup. Tables are created in order:
1. Messages table
2. Users table
3. Roles table
4. Organizations table
5. Junction tables for relationships

## 📚 API Endpoints

### Health & Documentation
- `GET /health` - Health check with database status
- `GET /docs` - Swagger UI documentation

### Authentication
- `GET /api/auth/google/url` - Get Google OAuth URL
- `POST /api/auth/google/login` - Complete OAuth login
- `GET /api/auth/me` - Get current user info
- `POST /api/auth/refresh` - Refresh access token
- `POST /api/auth/logout` - Logout

### Messages
- `GET /api/messages` - List messages
- `POST /api/messages` - Create message

### RBAC Management (Admin/Manager)
- `GET /api/roles` - List roles
- `POST /api/roles` - Create role
- `GET /api/organizations` - List organizations
- `POST /api/organizations` - Create organization

### Admin Operations (Admin only)
- `GET /api/admin/users` - List all users with roles
- `POST /api/admin/assign-role` - Assign role to user
- `POST /api/admin/assign-organization` - Add user to organization

### Setup
- `POST /api/setup/first-admin` - Make first user admin

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./...
```

### Frontend Tests
```bash
cd frontend
bun test
```

### Integration Tests
```bash
# Start services
docker-compose up db redis nats -d

# Run tests
go test -tags=integration ./...

# Cleanup
docker-compose down
```

## 🚀 Deployment

### Docker Compose (Local)
```bash
docker-compose up --build
```

### Kubernetes (Production)
```bash
# Deploy with Helm
helm install hackaton-demo charts/hackaton-demo/

# Or via Argo CD (configured for auto-sync)
```

### CI/CD Pipeline
- GitHub Actions for building and testing
- Docker image building and pushing
- Helm chart version bumping
- Argo CD for Kubernetes deployment

## 📊 Monitoring

### Health Checks
- Database connectivity monitoring
- NATS connection health
- Redis connection status
- Automatic reconnection handling

### Event Statistics
```bash
# Get event bus stats
curl http://localhost:8080/api/events/stats
```

### NATS Dashboard
- Visit: http://localhost:8222
- Monitor message rates and performance
- View JetStream streams and storage

## 🔒 Security Features

### Authentication
- Google OAuth 2.0 integration
- JWT token-based sessions
- Automatic token refresh
- Secure token storage

### Authorization
- Role-based access control (RBAC)
- Protected API endpoints
- Frontend route guards
- Admin-only operations

### Security Headers
- CORS configuration
- Rate limiting on auth endpoints
- Input validation
- SQL injection protection

## 🛠️ Development Workflow

### Adding New Features

#### 1. Create Model
```go
// models/example.go
type Example struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}
```

#### 2. Create Controller
```go
// controllers/example_controller.go
type ExampleController struct {
    dbManager *database.DBManager
}

func (ec *ExampleController) ExampleHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Handle request
    }
}
```

#### 3. Add to Router
```go
// handlers/router.go
mux.HandleFunc("/api/example", r.exampleController.ExampleHandler())
```

#### 4. Add Events (Optional)
```go
// Publish event for tracking
eventService.PublishExampleEvent(example.ID, example.Name)
```

### Database Changes
1. Create migration files in `backend/migrations/`
2. Update models in `backend/models/`
3. Test with `go run main.go` (migrations run automatically)

## 🐛 Troubleshooting

### Common Issues

#### Backend Can't Connect to Database
```bash
# Check if database is running
docker-compose ps db

# Check database logs
docker-compose logs db

# Verify connection string in .env
DB_URL=postgres://postgres:postgres@db:5432/postgres?sslmode=disable
```

#### OAuth Errors
- Verify `GOOGLE_CLIENT_ID` and `GOOGLE_CLIENT_SECRET` in `.env`
- Check redirect URI matches exactly in Google Console
- Ensure OAuth consent screen is configured

#### NATS Connection Issues
```bash
# Check NATS status
docker-compose ps nats

# Test connectivity
curl http://localhost:8222/healthz

# Check logs
docker-compose logs nats
```

#### Port Conflicts
```bash
# Check what's using ports
lsof -i :3000  # Frontend
lsof -i :8080  # Backend
lsof -i :5432  # Database
lsof -i :6379  # Redis
lsof -i :4222  # NATS
```

### Debug Mode
```bash
# Enable debug logging
export LOG_LEVEL=debug

# Start backend with verbose logging
go run main.go
```

## 📈 Performance

### Current Capabilities
- **HTTP Requests**: Concurrent handling (Go goroutines)
- **Database**: Connection pooling and health monitoring
- **Events**: Asynchronous processing via NATS
- **Caching**: Redis integration ready

### Optimization Tips
- Use database indexes for frequently queried fields
- Implement Redis caching for expensive operations
- Use NATS for background processing
- Monitor database connection pool usage

## 🔮 Future Enhancements

### Planned Features
- [ ] Real-time notifications via WebSockets
- [ ] Advanced caching strategies
- [ ] API rate limiting and throttling
- [ ] Comprehensive logging and monitoring
- [ ] Multi-tenant support
- [ ] Advanced RBAC with custom permissions

### Contributing
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📚 Resources

### Documentation
- [Go Documentation](https://golang.org/doc/)
- [React Documentation](https://reactjs.org/docs/)
- [NATS Documentation](https://docs.nats.io/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

### Learning Resources
- [Go by Example](https://gobyexample.com/)
- [React Tutorial](https://reactjs.org/tutorial/tutorial.html)
- [NATS Best Practices](https://docs.nats.io/running-a-nats-service/configuration)

---

**🎉 Happy Hacking! This template gives you a solid foundation for building amazing applications quickly.**

For questions or issues, check the troubleshooting section above or create an issue in the repository.

## 🎨 Frontend Features

The React frontend includes:
- **Google OAuth Integration**: Complete OAuth 2.0 flow with Google
- **User Profile Display**: Shows user information from Google
- **Token Management**: Automatic token refresh and logout
- **Responsive Design**: Clean, modern UI with Google branding
- **Error Handling**: Comprehensive error messages and loading states
- **Dark Mode Support**: Automatic theme switching with system preference detection
- **Theme Persistence**: Remembers user's theme choice across sessions
- **Redux Toolkit State Management**: Centralized state management across all pages
- **RTK Query API Handling**: Automatic caching, loading states, and error handling

## 🏗️ Frontend Architecture

### **State Management with Redux Toolkit**
- **Centralized Store**: Single source of truth for all application state
- **Slice-based Architecture**: Organized state management with auth, theme, and UI slices
- **TypeScript Support**: Full type safety throughout the Redux store
- **DevTools Integration**: Redux DevTools for debugging and state inspection

### **API Management with RTK Query**
- **Automatic Caching**: Intelligent caching with automatic background updates
- **Loading States**: Built-in loading and error state management
- **Optimistic Updates**: Immediate UI updates with automatic rollback on errors
- **Request Deduplication**: Prevents duplicate API calls
- **Automatic Re-fetching**: Smart re-fetching on focus and reconnection

### **Store Structure**
```
store/
├── index.ts          # Main store configuration
├── api.ts            # RTK Query API endpoints
├── hooks.ts          # Typed Redux hooks
└── slices/
    ├── authSlice.ts  # Authentication state
    ├── themeSlice.ts # Dark/light mode state
    └── uiSlice.ts    # UI notifications and loading states
```

### **Usage Examples**

#### **State Selection**
```tsx
// Select state from Redux store
const theme = useAppSelector((state) => state.theme.theme);
const user = useAppSelector((state) => state.auth.user);
```

#### **API Calls**
```tsx
// RTK Query hooks for API operations
const { data: users, isLoading, error } = useGetUsersQuery();
const [createUser, { isLoading: isCreating }] = useCreateUserMutation();

// Mutations with automatic cache invalidation
await createUser(userData).unwrap();
```

#### **Dispatching Actions**
```tsx
// Dispatch actions to update state
const dispatch = useAppDispatch();
dispatch(showSuccess('Operation completed!'));
dispatch(toggleTheme());
```

