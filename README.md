# flick-api

A production-ready REST API for managing movies, built with Go. Features user authentication, role-based access control, rate limiting, email notifications, and comprehensive monitoring.

## Overview

**flick-api** is a RESTful API service that provides endpoints for managing a movies database with user accounts and permission-based access control. The API is production-ready with features like automatic rate limiting, CORS support, structured logging, metrics tracking, and graceful shutdown handling.

### Key Features

- 🎬 **Movies Management** - Create, read, update, and delete movie records
- 👥 **User Authentication** - JWT-based token authentication with email verification
- 🔐 **Permission-Based Access Control** - Fine-grained permissions for read/write operations
- 📧 **Email Notifications** - SMTP integration for user registration confirmations
- ⚡ **Rate Limiting** - Token bucket-based request throttling
- 🌐 **CORS Support** - Cross-origin resource sharing enabled
- 📊 **Metrics & Monitoring** - Built-in metrics endpoint for performance tracking
- 🛡️ **Error Handling** - Comprehensive error handling with proper HTTP responses
- 🔄 **Graceful Shutdown** - Clean server shutdown with background task completion

## Technology Stack

- **Language**: Go 1.26.1
- **Database**: PostgreSQL
- **Router**: [julienschmidt/httprouter](https://github.com/julienschmidt/httprouter)
- **Authentication**: JWT tokens with bcrypt password hashing
- **Rate Limiting**: Token bucket algorithm (golang.org/x/time/rate)
- **Email**: SMTP via go-mail

## Project Structure

```
flick-api/
├── bin/                          # Compiled binaries (development & production)
├── cmd/api/                      # Application entry point
│   ├── main.go                   # Server initialization and configuration
│   ├── routes.go                 # API route definitions
│   ├── server.go                 # HTTP server setup and shutdown handling
│   ├── middleware.go             # Middleware implementations
│   ├── handlers.go               # HTTP request handlers
│   ├── errors.go                 # Error response formatting
│   ├── context.go                # Request context utilities
│   └── helpers.go                # Helper functions
├── internal/
│   ├── data/                     # Database models and access layers
│   │   ├── models.go             # Data models
│   │   ├── movies.go             # Movie model operations
│   │   ├── users.go              # User model operations
│   │   ├── tokens.go             # Token model operations
│   │   ├── permissions.go        # Permission model operations
│   │   ├── filters.go            # Database filtering utilities
│   │   └── runtime.go            # Runtime information
│   ├── mailer/                   # Email service
│   │   ├── mailer.go             # SMTP email implementation
│   │   └── templates/            # Email templates
│   ├── validator/                # Input validation
│   │   └── validator.go          # Validation helpers
│   └── vcs/                      # Version control system info
│       └── vcs.go                # Version information
├── migrations/                   # Database migrations (SQL)
├── remote/                       # Production deployment configuration
├── go.mod                        # Go module definition
├── Makefile                      # Build and development tasks
└── README.md                     # This file
```

## API Endpoints

### Health Check
- `GET /v1/healthcheck` - Server health status

### Movies
- `GET /v1/movies` - List all movies (requires `movies:read`)
- `POST /v1/movies` - Create a new movie (requires `movies:write`)
- `GET /v1/movies/:id` - Get a movie by ID (requires `movies:read`)
- `PATCH /v1/movies/:id` - Update a movie (requires `movies:write`)
- `DELETE /v1/movies/:id` - Delete a movie (requires `movies:write`)

### Users
- `POST /v1/users` - Register a new user
- `PUT /v1/users/activated` - Activate a user account

### Authentication
- `POST /v1/tokens/authentication` - Create an authentication token

### Monitoring
- `GET /debug/vars` - Runtime metrics and statistics

## Getting Started

### Prerequisites

- Go 1.26.1 or later
- PostgreSQL 12+
- `migrate` CLI tool (for database migrations)
- SMTP server credentials (for email notifications)

### Installation

1. **Clone the repository**
   ```sh
   git clone <repository-url>
   cd flick-api
   ```

2. **Set environment variables**
   Create a `.envrc` file in the project root:
   ```sh
   export FLICK_DB_DSN="user=postgres password=yourpassword dbname=flick sslmode=disable"
   export FLICK_SMTP_HOST="smtp.gmail.com"
   export FLICK_SMTP_PORT="587"
   export FLICK_SMTP_USERNAME="your-email@gmail.com"
   export FLICK_SMTP_PASSWORD="your-app-password"
   ```

3. **Build the application**
   ```sh
   make build/api
   ```

4. **Run database migrations**
   ```sh
   make db/migrations/up
   ```

5. **Start the server**
   ```sh
   make run/api
   ```

   The API will be available at `http://localhost:4000` by default.

## Development

### Available Commands

Use `make help` to see all available commands. Key commands include:

```sh
# Build the application
make build/api

# Run the application with database DSN
make run/api

# Database migrations
make db/migrations/create name="your_migration_name"
make db/migrations/up
make db/migrations/down
make db/migrations/version
```

### Configuration

Configuration is managed via command-line flags in `cmd/api/main.go`:

- `-port` - Server port (default: 4000)
- `-env` - Environment (default: development)
- `-db-dsn` - Database connection string (required)
- `-db-max-open-conns` - Max open database connections
- `-db-max-idle-conns` - Max idle database connections
- `-db-max-idle-time` - Max idle time for connections
- `-limiter-rps` - Rate limiter requests per second
- `-limiter-burst` - Rate limiter burst size
- `-limiter-enabled` - Enable rate limiting
- `-smtp-host` - SMTP host for email
- `-smtp-port` - SMTP port
- `-smtp-username` - SMTP username
- `-smtp-password` - SMTP password
- `-smtp-sender` - Email sender address
- `-cors-trusted-origins` - Comma-separated list of trusted CORS origins

### Example Usage

**Create a new user:**
```sh
curl -X POST http://localhost:4000/v1/users \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "name": "John Doe", "password": "securepassword"}'
```

**Authenticate:**
```sh
curl -X POST http://localhost:4000/v1/tokens/authentication \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "securepassword"}'
```

**List movies (requires authentication):**
```sh
curl -X GET http://localhost:4000/v1/movies \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Database Migrations

Migrations are managed using the `migrate` tool and located in the `migrations/` directory. Each migration is numbered sequentially with `.up.sql` and `.down.sql` files.

### Create a new migration:
```sh
make db/migrations/create name="add_user_preferences_table"
```

### Apply migrations:
```sh
make db/migrations/up
```

### Rollback one migration:
```sh
make db/migrations/down
```

### Check current version:
```sh
make db/migrations/version
```

## Deployment

### Building for Production

```sh
make build/api
```

This generates:
- `bin/api` - Local binary (macOS/Linux)
- `bin/linux_amd64/api` - Linux x86_64 binary

### Environment Configuration

For production, set appropriate environment variables:

```sh
export FLICK_ENV="production"
export FLICK_DB_DSN="postgresql://user:password@host:5432/flick"
export FLICK_LIMITER_ENABLED="true"
export FLICK_LIMITER_RPS="2"
export FLICK_CORS_ORIGINS="https://yourdomain.com"
```

### Running in Production

```sh
./bin/linux_amd64/api \
  -port=4000 \
  -env=production \
  -db-dsn=$FLICK_DB_DSN \
  -limiter-enabled=true
```

## Architecture

### Middleware Stack

The API applies middleware in order:
1. **Metrics** - Track request statistics
2. **Panic Recovery** - Recover from panics and return proper error responses
3. **CORS** - Enable cross-origin requests
4. **Rate Limiting** - Throttle requests based on configured limits
5. **Authentication** - Validate and extract JWT tokens

### Authentication Flow

1. User registers with email and password
2. Password is hashed using bcrypt
3. User receives activation email
4. User activates account by clicking link
5. User authenticates with credentials
6. API returns JWT token
7. Token is included in `Authorization: Bearer` header for subsequent requests

### Permission Model

Permissions are fine-grained and composable:
- `movies:read` - Read movie data
- `movies:write` - Create, update, or delete movies
- Additional permissions can be defined and assigned per user

## Monitoring and Debugging

### Health Check
```sh
curl http://localhost:4000/v1/healthcheck
```

### Metrics
```sh
curl http://localhost:4000/debug/vars
```

Returns runtime statistics including goroutine count, allocated memory, and custom application metrics.

### Structured Logging

The application uses `log/slog` for structured logging. Logs include contextual information and are designed for easy parsing and analysis.

## Contributing

When contributing to this project:

1. Follow Go code conventions and best practices
2. Keep the project structure organized
3. Write clear commit messages
4. Test your changes before submitting
5. Update documentation as needed

## License

[Insert your license information here]

## Support

For issues, questions, or contributions, please refer to the repository's issue tracker.