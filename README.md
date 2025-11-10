# LBK Points Transfer API

A RESTful API service built with Go and Fiber framework for managing user accounts and point transfers, designed following hexagonal architecture principles.

## Features

- **User Management**: CRUD operations for user accounts with Thai banking profile support
- **Point Transfers**: Secure point transfer system with transaction ledger
- **Hexagonal Architecture**: Clean separation of concerns with ports and adapters pattern
- **SQLite Database**: Lightweight database with automatic migrations
- **REST API**: JSON-based API endpoints with proper HTTP status codes

## Project Structure

This project follows the [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

```
├── api/                    # API documentation and specifications
├── bin/                    # Compiled binaries
├── cmd/
│   └── server/            # Application entry point
├── configs/               # Configuration files and templates
├── docs/                  # Database schema and documentation
├── internal/              # Private application code
│   ├── adapter/           # External adapters (database, HTTP)
│   ├── domain/            # Core business entities
│   ├── handler/           # HTTP request handlers
│   ├── port/              # Application interfaces/ports
│   └── service/           # Business logic layer
└── Makefile              # Build automation
```

## Getting Started

### Prerequisites

- Go 1.19 or higher
- SQLite3

### Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd workshop4-backend
```

2. Install dependencies:

```bash
make deps
```

3. Build the application:

```bash
make build
```

4. Run the server:

```bash
make run
```

The server will start on `http://localhost:3000`

### Development

For development with hot reload:

```bash
make dev
```

## API Endpoints

### Users

- `GET /users` - List all users
- `GET /users/:id` - Get user by ID
- `POST /users` - Create new user
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

### Transfers

- `POST /transfers` - Create point transfer
- `GET /transfers/:id` - Get transfer details
- `GET /users/:id/transfers` - Get user's transfer history

## Database Schema

The application uses SQLite with the following tables:

- **users**: User account information
- **transfers**: Point transfer transactions
- **point_ledger**: Point balance ledger

See [docs/database.md](docs/database.md) for detailed schema information and ER diagram.

## Architecture

This project implements hexagonal architecture (ports and adapters) pattern:

- **Domain Layer**: Core business entities (`internal/domain/`)
- **Port Layer**: Interfaces defining contracts (`internal/port/`)
- **Adapter Layer**: External integrations (`internal/adapter/`)
- **Service Layer**: Business logic (`internal/service/`)
- **Handler Layer**: HTTP handlers (`internal/handler/`)

## Configuration

Configuration can be managed through `configs/app.yaml`. See `configs/README.md` for available options.

## Available Commands

```bash
make build      # Build the application binary
make test       # Run tests
make clean      # Clean build artifacts
make run        # Build and run the application
make dev        # Run in development mode
make deps       # Download dependencies
make tidy       # Tidy go modules
make lint       # Run linter
make format     # Format code
make db-reset   # Reset database
make test-unit  # Run unit tests only
make test-coverage  # Run tests with coverage
make test-coverage-html  # Generate HTML coverage report
make help       # Show available commands
```

## Testing

Run the test suite:

```bash
make test
```

Run unit tests only:

```bash
make test-unit
```

Run tests with coverage:

```bash
make test-coverage
```

Generate HTML coverage report:

```bash
make test-coverage-html
```

### Testing Strategy

- **Domain Tests**: Unit tests for domain model validation logic (100% coverage)
- **Service Tests**: Unit tests for business logic with mocked repositories
- **Integration Tests**: Would test the full flow with real database (separate from unit tests)

The unit tests focus on:

- Input validation
- Business rule enforcement
- Error handling
- Mock-based testing for external dependencies

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and linting
5. Submit a pull request

## License

This project is licensed under the MIT License.
