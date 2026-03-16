# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the application
go run ./cmd/main.go

# Build
go build -o animap-go-api ./cmd/main.go

# Run all tests
go test ./...

# Run a single test
go test ./internal/... -run TestFunctionName -v

# Start local PostgreSQL (required before running)
docker-compose up -d

# Tidy dependencies
go mod tidy
```

## Architecture

This is a **Clean Architecture** Go REST API for an anime database platform. Dependency flow: `HTTP Handlers → Services → Repositories → PostgreSQL`.

All dependencies are wired manually via constructor injection in `cmd/main.go` — there is no DI framework.

### Layer Responsibilities

| Layer | Location | Responsibility |
|---|---|---|
| HTTP Handlers | `internal/adapters/https/` | Bind requests, call services, return JSON |
| Services | `internal/core/services/` | Business logic, coordinate repositories |
| Repositories | `internal/adapters/repositories/` | GORM queries against PostgreSQL |
| Entities | `internal/core/entities/` | GORM-mapped database models |
| DTOs | `internal/core/dtos/` | Request/response structs for the HTTP layer |

### Key Patterns

- **Read/Write DB Separation**: All repositories receive both a primary (write) and replica (read) `*gorm.DB`. Use `r.db` for writes and `r.dbReplica` for reads.
- **Custom Error Type**: Always return `*errs.AppError` from services. Use `errs.NewNotFoundError()`, `errs.NewBadRequestError()`, `errs.NewUnauthorizedError()`, or `errs.NewUnexpectedError()`. The `handleError()` in `internal/adapters/https/handler.go` maps these to HTTP status codes.
- **JWT Auth via Cookie**: `AuthRequired` middleware (`internal/middleware/auth.go`) extracts a `"jwt"` cookie, validates it using `JWT_SECRET` env var, and sets `userId` in the Gin context.
- **Auto-migration on startup**: `internal/database/database.go` runs GORM `AutoMigrate` for all entities on every startup.

### Configuration

- `config.yaml` — primary config (DB connections, AWS, MAL client ID, app port). Loaded by Viper. **Git-ignored.**
- `.env` — secrets (`JWT_SECRET`, `GIN_MODE`, `AWS_S3_BUCKET`). **Git-ignored.**
- `docker-compose.yml` — runs PostgreSQL locally on port 5432 with credentials matching `config.yaml`.

### External Integrations

- **AWS S3** (`internal/adapters/aws/`): Generates presigned upload URLs for user avatars.
- **MyAnimeList API** (`internal/adapters/external_api/`): Used for bulk anime data import/migration via `HttpAnimeMigrateHandler`.
- **WebSocket** (`internal/adapters/websocket/`): Real-time support, wired in routes.

### Adding a New Feature

Follow this pattern for any new domain (e.g., "review"):
1. Add entity in `internal/core/entities/` and register it in `database.go` AutoMigrate.
2. Add request/response DTOs in `internal/core/dtos/`.
3. Define repository interface + GORM implementation in `internal/adapters/repositories/`.
4. Define service interface + implementation in `internal/core/services/`.
5. Add HTTP handler in `internal/adapters/https/`.
6. Wire everything in `cmd/main.go` and add routes in `internal/route/route.go`.
