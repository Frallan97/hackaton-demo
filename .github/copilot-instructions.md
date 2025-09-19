# GitHub Copilot Instructions

This file provides context and guidelines for GitHub Copilot when working with the Hackaton Web4 project.

## Project Context

This is a full-stack web application with:
- **Backend**: Go REST API with PostgreSQL database
- **Frontend**: React with TypeScript, Vite, and Tailwind CSS
- **Development**: Hot reloading with Air (backend) and Vite (frontend)
- **Package Management**: Go modules (backend), Bun (frontend)

## Architecture Overview

```
hackaton-web4/
├── backend/                 # Go REST API
│   ├── controllers/         # HTTP request handlers
│   ├── models/             # Data structures
│   ├── services/           # Business logic
│   ├── middleware/         # HTTP middleware
│   ├── database/           # Database connection and migrations
│   ├── config/             # Configuration management
│   └── .air.toml           # Hot reloading configuration
├── frontend/               # React TypeScript app
│   ├── src/components/     # Reusable UI components
│   ├── src/admin/          # Admin dashboard components
│   └── src/lib/           # Utility functions
└── docker-compose.yml     # Local development setup
```

## Code Generation Guidelines

### Backend (Go)

When generating Go code:

1. **File Structure**:
   - Controllers: `*_controller.go` in `backend/controllers/`
   - Models: `*.go` (singular names) in `backend/models/`
   - Services: `*_service.go` in `backend/services/`
   - Middleware: descriptive names like `logging.go` in `backend/middleware/`

2. **Code Style**:
   - Follow standard Go conventions and use `gofmt`
   - Use snake_case for Go identifiers
   - Include proper error handling
   - Use structured logging with context

3. **Database Operations**:
   - Use parameterized queries for security
   - Implement proper connection pooling
   - Handle database errors gracefully
   - Use migrations for schema changes

4. **API Design**:
   - Follow RESTful conventions
   - Use consistent error responses
   - Include proper HTTP status codes
   - Add Swagger documentation comments

5. **Authentication**:
   - Implement Google OAuth integration
   - Use JWT tokens for session management
   - Include RBAC (Role-Based Access Control)
   - Protect routes with middleware

### Frontend (React/TypeScript)

When generating React/TypeScript code:

1. **File Structure**:
   - Components: `PascalCase.tsx` in appropriate directories
   - Hooks: `use*.ts` in `src/lib/` or component directories
   - Utilities: `camelCase.ts` in `src/lib/`
   - Pages: `PascalCase.tsx` in `src/`

2. **Code Style**:
   - Use functional components with hooks
   - Use camelCase for TypeScript/JavaScript identifiers
   - Include proper TypeScript types
   - Use descriptive component and prop names

3. **UI/Styling**:
   - Use Tailwind CSS classes exclusively
   - Avoid custom CSS when possible
   - Use Radix UI components for complex interactions
   - Follow modern UI/UX best practices

4. **State Management**:
   - Use React hooks for local state
   - Implement proper error boundaries
   - Use context for global state when needed
   - Handle loading and error states

## Development Patterns

### Environment Configuration

- Use `.env` files for local development
- Never include secrets in code
- Use environment variables for configuration
- Provide sensible defaults for development

### Error Handling

- Always handle errors gracefully
- Log errors with appropriate context
- Return meaningful error messages
- Use proper HTTP status codes

### Performance Considerations

- Optimize database queries
- Use appropriate indexes
- Implement caching where beneficial
- Minimize bundle sizes
- Use debug mode conditionally

### Security Best Practices

- Validate all inputs
- Use parameterized database queries
- Implement proper CORS policies
- Secure JWT token handling
- Follow OAuth 2.0 best practices

## Development Commands

### Backend Development
```bash
# Install dependencies
go mod tidy

# Install Air for hot reloading
go install github.com/air-verse/air@latest

# Development with hot reloading
DEBUG=true ENVIRONMENT=development air

# Normal run
go run main.go
```

### Frontend Development
```bash
# Install dependencies
bun install

# Development server
bun dev

# Build for production
bun run build
```

## Common Code Patterns

### Go Controller Pattern
```go
type ExampleController struct {
    dbManager *database.DBManager
}

func NewExampleController(dbManager *database.DBManager) *ExampleController {
    return &ExampleController{dbManager: dbManager}
}

func (ec *ExampleController) HandleExample() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Implementation
    }
}
```

### React Component Pattern
```tsx
interface ExampleProps {
    title: string;
    onAction?: () => void;
}

export function Example({ title, onAction }: ExampleProps) {
    return (
        <div className="p-4 bg-white rounded-lg shadow">
            <h2 className="text-xl font-semibold">{title}</h2>
            {onAction && (
                <button 
                    onClick={onAction}
                    className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
                >
                    Action
                </button>
            )}
        </div>
    );
}
```

### Database Migration Pattern
```sql
-- migrations/YYYYMMDDHHMMSS_description.up.sql
CREATE TABLE example (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- migrations/YYYYMMDDHHMMSS_description.down.sql
DROP TABLE IF EXISTS example;
```

## Testing Guidelines

- Write unit tests for business logic
- Test API endpoints thoroughly
- Test React components with proper mocking
- Use meaningful test descriptions
- Test both success and error scenarios

## Documentation Standards

- Use clear, descriptive comments
- Document API endpoints with Swagger
- Keep README files up to date
- Document environment variables
- Explain complex business logic

## Deployment Considerations

- Application is containerized with Docker
- Uses multi-stage builds for production
- Environment-specific configurations
- Health checks for all services
- Kubernetes-ready with Helm charts

When generating code, prioritize:
1. Security and proper error handling
2. Performance and scalability
3. Maintainability and readability
4. Following established patterns
5. Comprehensive testing coverage 