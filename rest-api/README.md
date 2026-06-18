# Simple Blog API

Simple Blog API is a Go backend service for a small blog platform. It supports user registration and login, JWT authentication, role and ownership authorization, post CRUD, comment creation/update/delete, health checks, PostgreSQL migrations, Swagger documentation, Docker execution, and automated tests.

## Table of Contents

- [Simple Blog API](#simple-blog-api)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Setup and Running](#setup-and-running)
    - [Prerequisites](#prerequisites)
    - [Environment](#environment)
    - [Run Locally](#run-locally)
    - [Run with Docker Compose](#run-with-docker-compose)
    - [Manual and Makefile Alternatives](#manual-and-makefile-alternatives)
    - [Useful Make Commands](#useful-make-commands)
  - [API Documentation](#api-documentation)
  - [Architecture](#architecture)
    - [Directory Structure](#directory-structure)
    - [Request Flow](#request-flow)
    - [Module Pattern](#module-pattern)
  - [Technology Choices](#technology-choices)
  - [Database Schema](#database-schema)
  - [Testing and Coverage](#testing-and-coverage)
  - [Known Limitations](#known-limitations)
  - [Future Improvements](#future-improvements)

## Features

- Authentication with register and login endpoints.
- Password hashing with bcrypt.
- JWT access tokens with configurable expiration.
- Role-aware user context stored in JWT claims.
- Ownership/admin authorization for protected post and comment mutations.
- Public post listing with pagination.
- PostgreSQL schema migrations and seed data.
- Swagger/OpenAPI documentation served by the app.
- Docker and Docker Compose support.
- Unit, handler, repository, service, middleware, and integration tests.

## Setup and Running

### Prerequisites

- Go 1.25 or newer.
- PostgreSQL 16 or compatible.
- Docker and Docker Compose, if running with containers.
- `migrate` CLI, if applying migrations through the Makefile.
- `psql`, if running seed files through the Makefile.
- `swag`, if regenerating Swagger docs.

### Environment

Copy the example environment file and update values if needed.

```bash
cp .env.example .env
```

Default local values:

```env
PORT=8080
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRATION_HOURS=24

DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=root
DB_NAME=blog
DB_SSLMODE=disable
```

The app accepts either `DATABASE_URL` or the split `DB_*` variables.

```env
DATABASE_URL=postgres://admin:root@localhost:5432/blog?sslmode=disable
```

### Run Locally

Install dependencies:

```bash
make deps
```

Manual alternative:

```bash
go mod download
go mod tidy
```

Apply database migrations:

```bash
make migrate-up
```

Manual alternative:

```bash
migrate -path ./database/migrations \
  -database "postgres://admin:root@localhost:5432/blog?sslmode=disable" \
  up
```

Seed initial data:

```bash
make seed
```

Manual alternative:

```bash
PGPASSWORD=root psql -h localhost -p 5432 -U admin -d blog -f database/seeds/001_initial_data.sql
```

Run the API:

```bash
make run
```

Manual alternative:

```bash
go run ./cmd/api/main.go
```

The server starts on:

```text
http://localhost:8080
```

### Run with Docker Compose

Start PostgreSQL and the API together:

```bash
docker compose up --build
```

The compose file starts:

- PostgreSQL on `localhost:5432`
- API on `localhost:8080`

### Manual and Makefile Alternatives

| Workflow | Makefile command | Manual command |
| --- | --- | --- |
| Install dependencies | `make deps` | `go mod download` then `go mod tidy` |
| Build binary | `make build` | `mkdir -p bin` then `go build -o bin/api ./cmd/api` |
| Run API | `make run` | `go run ./cmd/api/main.go` |
| Fresh local run | `make fresh-run` | Run migration, seed, build, then `./bin/api` |
| Run tests | `make test` | `go test ./...` |
| Run Ginkgo tests | `make ginkgo` | `ginkgo ./...` |
| Coverage summary | `make test-coverage` | `go test -coverprofile=coverage.out ./...` then `go tool cover -func=coverage.out` |
| HTML coverage | `make test-coverage-html` | `go test -coverprofile=coverage.out ./...` then `go tool cover -html=coverage.out -o coverage.html` |
| Apply migrations | `make migrate-up` | `migrate -path ./database/migrations -database "postgres://admin:root@localhost:5432/blog?sslmode=disable" up` |
| Roll back latest migration | `make migrate-down` | `migrate -path ./database/migrations -database "postgres://admin:root@localhost:5432/blog?sslmode=disable" down 1` |
| Migration status | `make migrate-status` | `migrate -path ./database/migrations -database "postgres://admin:root@localhost:5432/blog?sslmode=disable" version` |
| Create migration | `make migrate-create name=add_new_table` | `migrate create -ext sql -dir ./database/migrations -seq add_new_table` |
| Run seed | `make seed` | `PGPASSWORD=root psql -h localhost -p 5432 -U admin -d blog -f database/seeds/001_initial_data.sql` |
| Generate Swagger docs | `make swagger` | `swag init -g cmd/api/main.go -o ./docs/swagger --parseDependency --parseInternal` |

### Useful Make Commands

```bash
make deps                 # Download and tidy dependencies
make build                # Build binary into bin/api
make run                  # Run cmd/api/main.go
make fresh-run            # Clean, migrate, seed, build, and run
make test                 # Run all tests with go test
make ginkgo               # Run all tests with Ginkgo
make test-coverage        # Generate coverage summary
make test-coverage-html   # Generate coverage.html
make migrate-up           # Apply pending migrations
make migrate-down         # Roll back the latest migration
make migrate-status       # Show migration version
make swagger              # Regenerate Swagger docs
```

## API Documentation

Swagger files are generated under [docs/swagger](docs/swagger). When the app is running, open:

```text
http://localhost:8080/swagger/index.html
```

Main endpoints:

| Method | Path | Access |
| --- | --- | --- |
| `GET` | `/health` | Public |
| `POST` | `/auth/register` | Public |
| `POST` | `/auth/login` | Public |
| `GET` | `/posts` | Public |
| `GET` | `/posts/:id` | Public |
| `POST` | `/posts` | Authenticated |
| `PUT` | `/posts/:id` | Owner or admin |
| `DELETE` | `/posts/:id` | Owner or admin |
| `POST` | `/posts/:id/comments` | Authenticated |
| `PUT` | `/comments/:id` | Owner or admin |
| `DELETE` | `/comments/:id` | Owner or admin |

## Architecture

The project uses a modular, feature-based architecture. Each business module owns its domain contracts, DTOs, handler, service, repository, provider, and router. Shared infrastructure lives under `internal/platform`.

### Directory Structure

| Path | Purpose |
| --- | --- |
| [cmd/api](cmd/api) | Application entry point and graceful shutdown. |
| [config](config) | App, auth, and database configuration loading. |
| [database/migrations](database/migrations) | SQL migration files. |
| [database/seeds](database/seeds) | Initial seed data. |
| [docs/swagger](docs/swagger) | Generated OpenAPI/Swagger artifacts. |
| [internal/modules/auth](internal/modules/auth) | Register, login, JWT middleware, role checks, ownership checks. |
| [internal/modules/posts](internal/modules/posts) | Post CRUD and pagination. |
| [internal/modules/comments](internal/modules/comments) | Comment create, update, and delete flows. |
| [internal/modules/health](internal/modules/health) | Health check endpoint. |
| [internal/platform/database](internal/platform/database) | PostgreSQL connection and transaction helper. |
| [internal/platform/errors](internal/platform/errors) | Application errors and global Fiber error handler. |
| [internal/platform/response](internal/platform/response) | Standard API response shape. |
| [internal/platform/server](internal/platform/server) | Fiber setup, middleware, Swagger, and route registration. |
| [tests/integration](tests/integration) | API-level integration tests. |

### Request Flow

1. `cmd/api/main.go` loads configuration from `.env` or environment variables.
2. `internal/platform/database` creates the PostgreSQL connection using `sqlx`.
3. `internal/platform/di` wires repositories, services, handlers, routers, and shared infrastructure.
4. `internal/platform/server` creates the Fiber app, global middleware, Swagger route, and module routes.
5. Handlers parse and validate HTTP input.
6. Services enforce business rules such as uniqueness, ownership, and authorization.
7. Repositories execute SQL through `sqlx`.
8. Responses are returned through the shared response and error contracts.

### Module Pattern

Most feature modules follow this shape:

```text
domain      # Interfaces, request DTOs, response DTOs
handler     # HTTP request parsing and response formatting
service     # Business logic and authorization decisions
repository  # SQL queries and persistence
router      # Route registration
provider    # Wire provider set
```

This keeps business capability code together and avoids scattering one feature across many technical folders.

## Technology Choices

- Fiber: lightweight HTTP framework with an expressive routing and middleware API.
- PostgreSQL: reliable relational database for transactional blog data and ownership relationships.
- sqlx: keeps SQL explicit while reducing boilerplate around scanning and query execution.
- Viper: centralizes `.env` and environment variable configuration.
- JWT: stateless access-token authentication for API clients.
- bcrypt: standard password hashing algorithm with adaptive cost.
- Google Wire: compile-time dependency injection without runtime reflection.
- Ginkgo and Gomega: readable behavior-oriented tests for handlers, services, repositories, and integration flows.
- Testify and sqlmock: focused assertions and repository-level SQL mocking.
- Swagger/OpenAPI: interactive API documentation and machine-readable endpoint contract.
- Docker: repeatable local and containerized execution.

## Database Schema

Migrations live in [database/migrations](database/migrations). Seed scripts live in [database/seeds](database/seeds).

The current schema includes:

- `users`
- `roles`
- `user_roles`
- `posts`
- `comments`
- indexes for post and comment lookups

```mermaid
erDiagram
  USERS ||--o{ POSTS : writes
  USERS ||--o{ COMMENTS : writes
  POSTS ||--o{ COMMENTS : has
  USERS ||--o{ USER_ROLES : assigned
  ROLES ||--o{ USER_ROLES : granted

  USERS {
    UUID id PK
    VARCHAR username
    VARCHAR email
    TEXT password_hash
    TIMESTAMPTZ created_at
    TIMESTAMPTZ updated_at
  }

  ROLES {
    UUID id PK
    VARCHAR name
    TIMESTAMPTZ created_at
  }

  USER_ROLES {
    UUID user_id PK
    UUID role_id PK
    TIMESTAMPTZ assigned_at
  }

  POSTS {
    UUID id PK
    UUID author_id FK
    VARCHAR title
    TEXT content
    TIMESTAMPTZ created_at
    TIMESTAMPTZ updated_at
  }

  COMMENTS {
    UUID id PK
    UUID post_id FK
    UUID author_id FK
    TEXT content
    TIMESTAMPTZ created_at
    TIMESTAMPTZ updated_at
  }
```

## Testing and Coverage

Run all tests:

```bash
go test ./...
```

Makefile alternative:

```bash
make test
```

Run Ginkgo tests:

```bash
ginkgo ./...
```

Makefile alternative:

```bash
make ginkgo
```

Generate coverage summary:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out
```

Makefile alternative:

```bash
make test-coverage
```

Generate HTML coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

Makefile alternative:

```bash
make test-coverage-html
```

Current coverage from [coverage.out](coverage.out):

```text
total: (statements) 95.9%
```

Coverage highlights:

| Area | Coverage |
| --- | --- |
| Config env binding | 100.0% |
| Auth middleware and authorization helpers | 92.0% - 100.0% |
| Auth service | 95.0% - 96.7% |
| Auth repository | 87.5% - 100.0% |
| Post service | 88.9% - 100.0% |
| Post repository | 100.0% |
| Comment service | 100.0% |
| Comment repository | 100.0% |
| Health module | 100.0% |
| Database connection and transaction helper | 91.7% - 92.9% |

The repository also includes [coverage.html](coverage.html) for browser-based inspection of covered and uncovered statements.

## Known Limitations

- No refresh token mechanism yet; login returns only a JWT access token.
- No token revocation or logout blacklist, so issued access tokens remain valid until expiration.
- No rate limiting on authentication or write endpoints.
- No soft delete; deleting posts and comments permanently removes records.
- No admin management endpoints for assigning or revoking roles after registration.
- Pagination is implemented for the public posts list, but there is no comment listing endpoint with pagination.
- No full-text search or filtering for posts.
- No file upload support for post images or attachments.
- No CI/CD pipeline configuration is included.
- Docker setup is suitable for local development, but production deployment hardening is not included.

## Future Improvements

- Add refresh tokens with rotation and reuse detection.
- Add logout/token revocation support backed by Redis or a database table.
- Add rate limiting and request throttling for login, registration, and write-heavy endpoints.
- Add soft delete columns such as `deleted_at` and update queries to exclude deleted records.
- Add admin APIs for user management, role assignment, and moderation.
- Add paginated comment listing per post.
- Add post search, filtering, and sorting by author, keyword, and date.
- Add observability with structured logs, request IDs, metrics, and tracing.
- Add CI/CD workflow for linting, tests, coverage, Docker build, and deployment.
- Add production configuration examples for secrets, TLS, database pooling, and migrations.
- Add Redis-backed session management to store active user sessions, refresh token metadata, token blacklists, and expiration policies.
- Add Redis-based caching with cache invalidation strategies for posts and comment counts to improve read performance and reduce database contention on high-traffic endpoints.
