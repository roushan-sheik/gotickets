# Project Structure

This project follows a standard Go project layout, utilizing a layered architecture (Handler, Service, Repository) to maintain a clean separation of concerns.

## Directory Layout

```text
.
├── cmd/                        # Main applications for this project
│   └── main.go                 # Application entry point, router setup, database connection, and dependency injection
├── internal/                   # Private application and library code
│   ├── config/                 # Configuration loading and environment variables handling
│   ├── httpresponse/           # Standardized HTTP response wrappers and error handling
│   │   └── error.go            # Custom HTTP error structure implementing the `error` interface
│   └── user/                   # User domain module
│       ├── dto/                # Data Transfer Objects (Request/Response schemas)
│       │   ├── request.go      # Input payload structures for User APIs
│       │   └── response.go     # Output payload structures for User APIs
│       ├── entity.go           # Database models/entities for the User domain (GORM)
│       ├── handler.go          # HTTP transport layer (Echo handlers mapping requests to services)
│       ├── repository.go       # Data access layer (Database interactions for Users)
│       └── service.go          # Business logic layer (Core application rules for Users)
├── go.mod                      # Go module dependencies file
└── go.sum                      # Go module checksums file
```

## Architecture Layers (User Domain)

1. **Handler (`handler.go`)**: The presentation layer. It is responsible for parsing HTTP requests, input validation, and formatting HTTP responses. It delegates the actual work to the Service layer.
2. **Service (`service.go`)**: The business logic layer. It contains the core logic of the application. It receives validated data from handlers, applies business rules, and communicates with the Repository layer.
3. **Repository (`repository.go`)**: The data access layer. It is responsible for direct interactions with the database (e.g., executing SQL queries via GORM). It abstracts the database operations away from the Service layer.
4. **Entity (`entity.go`)**: The database schema representation. Used by GORM for migrations and querying.
5. **DTO (`dto/`)**: Structures used specifically for API requests and responses, keeping the API contract decoupled from the internal database entities.
