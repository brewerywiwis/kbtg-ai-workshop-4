# GitHub Copilot Instructions

This file provides custom instructions for GitHub Copilot to understand and follow the specific patterns, architecture, and coding standards used in this LBK Points Transfer API project.

## Code Style

### Go Code Standards

- Follow Go's official formatting using `gofmt`
- Use meaningful variable and function names in English
- Follow Go naming conventions:
  - Use camelCase for unexported functions and variables
  - Use PascalCase for exported functions, types, and variables
  - Use ALL_CAPS for constants
- Keep functions small and focused on a single responsibility
- Use early returns to reduce nesting
- Always handle errors explicitly - never ignore them
- Use context.Context for cancellation and timeouts in service methods
- Add comprehensive error messages with context

### JSON and API Standards

- Use camelCase for JSON field names
- Always include proper HTTP status codes
- Use consistent error response format:
  ```json
  {
    "error": "ERROR_TYPE",
    "message": "Human readable message"
  }
  ```
- Include proper validation with meaningful error messages

### Database Standards

- Use snake_case for database column names
- Always use transactions for operations that modify multiple tables
- Include proper timestamps (created_at, updated_at)
- Use proper foreign key constraints

## Things to Avoid

### Anti-patterns

- **DO NOT** put business logic in handlers - keep them thin
- **DO NOT** directly access database from handlers
- **DO NOT** ignore errors or use `_` for error handling
- **DO NOT** use global variables for application state
- **DO NOT** hardcode database credentials or configuration
- **DO NOT** mix domain models with database models
- **DO NOT** create circular dependencies between packages
- **DO NOT** use `interface{}` unless absolutely necessary - prefer typed interfaces

### Code Smells

- Avoid long parameter lists - use struct parameters instead
- Don't create god objects or oversized structs
- Avoid deep nesting - prefer early returns
- Don't duplicate validation logic across layers
- Avoid tight coupling between packages

### Security Issues

- Never log sensitive data (passwords, tokens, personal information)
- Always validate and sanitize user input
- Use proper SQL parameterization to prevent injection
- Don't expose internal error details to clients

## Code Architecture

### Project Structure (Standard Go Layout)

```
├── cmd/server/           # Application entry point
├── internal/             # Private application code
│   ├── domain/          # Core business entities and rules
│   ├── port/            # Interfaces (repository contracts)
│   ├── adapter/         # External integrations (database, HTTP)
│   ├── service/         # Business logic layer
│   ├── handler/         # HTTP request handlers
│   └── app/             # Application initialization
├── configs/             # Configuration files
├── docs/                # Documentation
└── api/                 # API specifications
```

### Hexagonal Architecture Principles

This project follows hexagonal (ports and adapters) architecture:

1. **Domain Layer** (`internal/domain/`):

   - Contains core business entities (User, Transfer, PointLedger)
   - Should have NO external dependencies
   - Pure business logic and rules

2. **Port Layer** (`internal/port/`):

   - Defines interfaces/contracts for external dependencies
   - Repository interfaces live here
   - No implementation details

3. **Adapter Layer** (`internal/adapter/`):

   - Implements port interfaces
   - Database repositories (SQLite implementations)
   - External service integrations

4. **Service Layer** (`internal/service/`):

   - Orchestrates business operations
   - Uses domain entities and port interfaces
   - Contains transaction logic and validation

5. **Handler Layer** (`internal/handler/`):
   - HTTP request/response handling
   - Input validation and transformation
   - Delegates to service layer

### Dependency Flow

```
Handler → Service → Port Interface ← Adapter
    ↓         ↓            ↓
  Domain ← Domain    Domain Entity
```

### Code Generation Guidelines

When generating new code:

1. **New API Endpoints**:

   - Add handler function in appropriate handler file
   - Add service method for business logic
   - Update repository interface if needed
   - Implement repository method in adapter
   - Add proper error handling at each layer

2. **New Domain Entities**:

   - Create in `internal/domain/`
   - Add repository interface in `internal/port/`
   - Implement repository in `internal/adapter/`
   - Add service methods in `internal/service/`
   - Create handlers in `internal/handler/`

3. **Database Operations**:
   - Always use transactions for multi-table operations
   - Include proper error handling and rollback
   - Use prepared statements for SQL queries
   - Follow the existing SQLite patterns

### Testing Approach

- Unit tests for domain logic
- Integration tests for repository implementations
- API tests for handler functions
- Mock interfaces for testing isolation

### Example Code Patterns

#### Handler Pattern

```go
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
    var user domain.User
    if err := c.BodyParser(&user); err != nil {
        return c.Status(400).JSON(fiber.Map{
            "error": "INVALID_REQUEST",
            "message": "Invalid request body",
        })
    }

    if err := h.userService.Create(&user); err != nil {
        // Handle service errors appropriately
    }

    return c.Status(201).JSON(user)
}
```

#### Service Pattern

```go
func (s *UserService) Create(user *domain.User) error {
    // Validate business rules
    if err := s.validateUser(user); err != nil {
        return err
    }

    // Set timestamps
    now := time.Now()
    user.CreatedAt = now
    user.UpdatedAt = now

    // Delegate to repository
    return s.repo.Create(user)
}
```

#### Repository Pattern

```go
func (r *SqliteUserRepository) Create(user *domain.User) error {
    query := `INSERT INTO users (name, phone, email, member_since, membership_level, member_id, points, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

    result, err := r.db.Exec(query, user.Name, user.Phone, user.Email, user.MemberSince, user.MembershipLevel, user.MemberID, user.Points, user.CreatedAt, user.UpdatedAt)
    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    id, _ := result.LastInsertId()
    user.ID = int(id)
    return nil
}
```

Follow these patterns when generating or suggesting code for this project.
