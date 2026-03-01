# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

```bash
# Run the application (loads .env automatically)
make start

# Run directly
go run ./cmd/main.go

# Build binary
go build -o main ./cmd/main.go

# Start infrastructure (PostgreSQL + app) via Docker
docker compose -f manifests/docker-compose.yaml --env-file manifests/.env.docker up --build

# Tidy dependencies
go mod tidy
```

There are no tests in this codebase currently.

## Architecture

This is an **Employee Management REST API** built with Go Fiber, using Uber Dig for dependency injection.

### Dependency Injection Flow

All components are wired through Uber Dig in `internal/containers/container.go`:

```
Config → Database → Fiber App → Router (with middleware) → Service → StartServer(:3000)
```

Entry point is `cmd/main.go`, which builds the DI container and invokes `StartServer`.

### Key Layers

- **`internal/config/`** — Loads env vars (DB, JWKS_URL, Vault) with defaults via godotenv
- **`internal/database/`** — GORM/PostgreSQL connection, auto-migrates the `Employees` model
- **`internal/routes/`** — Defines endpoints and attaches RBAC middleware per route
- **`internal/middlewares/`** — Keycloak JWT validation (JWKS) and role-based access control
- **`internal/service/`** — Business logic (CRUD operations, Keycloak login/refresh, Vault encrypt/decrypt)
- **`internal/models/`** — `Employees` GORM model plus generic `ApiRequest`/`ApiResponse` wrappers

### Authentication & Authorization

Keycloak handles auth. The middleware chain is: **JWT validation → role extraction from `realm_access.roles` claim → RBAC check**.

Three roles: `user` (read-only), `manager` (read + update), `admin` (full CRUD). Role requirements are defined per-route in `internal/routes/setup_router.go`.

### Vault Encryption

Email and Mobile fields are encrypted at rest via HashiCorp Vault's transit engine (`internal/service/vault.go`). Each employee gets a Vault entity ID (`{employeeID}pii`). Fields are encrypted on create/update and decrypted on read.

### API Structure

All requests/responses use generic wrappers (`ApiRequest`/`ApiResponse` in `internal/models/models.go`) with metadata (ID, version, timestamp, params).

**Public:** `POST /api/login`, `POST /api/refresh`
**Protected:** All under `/api/v1/` — `listEmployees`, `addEmployee`, `getEmployee/{id}`, `updateEmployee/{id}`, `deleteEmployee/{id}`

### Supporting Files

- `scripts/openapi.yaml` — OpenAPI 3.0.3 spec
- `scripts/api.postman_collection.json` — Postman collection for manual testing
- `scripts/realm-export.json` — Keycloak realm config (import into Keycloak for roles/clients)
- `.env.example` — Reference for required environment variables
