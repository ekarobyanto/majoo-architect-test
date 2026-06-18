# Simple Blog

A production-ready Go backend boilerplate built with a modular, feature-based architecture. The repository includes the full backend stack expected in a typical service deliverable: database schema and migrations, Swagger/OpenAPI documentation, Docker configuration, and a complete setup guide.

## Deliverables

- Database schema and migration files are stored under [database/migrations](database/migrations) and [database/seeds](database/seeds).
- API documentation is generated from Swagger annotations and published under [docs/swagger](docs/swagger), with the UI available at `/swagger/index.html` when the app is running.
- Docker configuration is provided through [Dockerfile](Dockerfile) and [docker-compose.yml](docker-compose.yml).
- This README covers setup, architecture, technology choices, and known limitations.

## Overview

The service is organized around business capabilities rather than technical layers. Each module owns its handler, service, repository, domain objects, and routing, which keeps related logic close together and easier to maintain.

### Main Modules

- `auth`: registration, login, JWT authentication, and ownership/authorization checks.
- `posts`: post CRUD operations and pagination.
- `comments`: comment CRUD operations tied to posts.
- `health`: health-check endpoint for readiness and monitoring.

## Database Schema and Migrations

The database schema is expressed as SQL migration files in [database/migrations](database/migrations). The repository includes migrations for:

- PostgreSQL extensions
- `users`
- `roles`
- `user_roles`
- `posts`
- `comments`
- supporting indexes for posts and comments

Seed data lives in [database/seeds](database/seeds) and is intended to bootstrap the initial role records and other starter data.

### ERD

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

## API Documentation

The API is documented with Swagger/OpenAPI annotations in the Go source. The generated artifacts live in [docs/swagger](docs/swagger), and the Swagger UI is served by the application itself.

When the server is running, open:

- `http://localhost:8080/swagger/index.html`

## Architecture

The project follows a modular feature-based layout. Instead of placing all handlers together, all services together, and all repositories together, each feature owns its own slice of the codebase.

### Directory Structure

- [cmd/api](cmd/api): application entry point.
- [config](config): environment and config loading.
- [database/migrations](database/migrations): SQL migrations.
- [database/seeds](database/seeds): seed scripts.
- [docs/swagger](docs/swagger): generated Swagger/OpenAPI artifacts.
- [internal/modules](internal/modules): feature modules.
- [internal/platform](internal/platform): shared infrastructure such as database and server setup.
- [tests/integration](tests/integration): integration tests.

### Request Flow

1. `cmd/api/main.go` loads configuration.
2. `internal/platform/database` opens the PostgreSQL connection.
3. `internal/platform/di` wires the module providers.
4. Module routers register HTTP routes in Fiber.
5. Handlers validate requests, call services, and return structured responses.

## Setup and Running Instructions

### Prerequisites

- Go 1.25 or higher
- PostgreSQL 16 or compatible
- Docker and Docker Compose if you want containerized runs

### Local Setup

1. Install dependencies.

```bash
go mod download
```

2. Create your environment file.

```bash
cp .env.example .env
```

3. Configure either `DATABASE_URL` or the split `DB_*` variables.

```bash
DATABASE_URL=postgres://admin:root@localhost:5432/blog?sslmode=disable
```

4. Run the application.

```bash
go run cmd/api/main.go
```

### Using the Makefile

The Makefile exposes the common workflows for this project:

```bash
make deps        # Download and tidy dependencies
make build       # Build the binary into bin/api
make run         # Run the application directly with go run
make fresh-run   # Clean, migrate, seed, rebuild, and run the app from scratch
make test        # Run the full Go test suite
make help        # Show all available targets
```

If you want the shortest path to a local run, use `make run` after configuring `.env`.

## Database Migrations

Migration files live in [database/migrations](database/migrations), and the Makefile provides helper targets for applying them.

```bash
make migrate-up       # Apply all pending migrations
make migrate-down     # Roll back the last migration
make migrate-status   # Show the current migration version
make migrate-create name=add_new_table  # Create a new SQL migration pair
```

These targets require the `migrate` CLI to be installed and available on your `PATH`. They use the database settings from `.env` to build the connection string.

### Docker

Build the image:

```bash
docker build -t majoo-architect-app .
```

Run against an existing database:

```bash
docker run --rm -p 8080:8080 \
  -e PORT=8080 \
  -e DATABASE_URL='postgres://admin:root@host.docker.internal:5432/blog?sslmode=disable' \
  majoo-architect-app
```

Run the app and database together:

```bash
docker compose up --build
```

## Testing

The project uses Ginkgo and Gomega for behavior-focused tests, plus Go’s standard testing package for unit and integration-style checks.

### Run all tests

```bash
ginkgo ./...
```

### Run tests with Go

```bash
go test ./...
```

### Coverage Report

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Technology Choices

- [Fiber](https://gofiber.io/): fast HTTP framework with low overhead and an expressive API.
- [PostgreSQL](https://www.postgresql.org/): reliable relational database with strong consistency and mature tooling.
- [sqlx](https://github.com/jmoiron/sqlx): keeps SQL explicit while making row mapping simpler than raw `database/sql`.
- [Viper](https://github.com/spf13/viper): flexible environment and config loader.
- [Ginkgo](https://onsi.github.io/ginkgo/) and [Gomega](https://onsi.github.io/gomega/): readable test style for behavior-driven scenarios.
- [Swagger/OpenAPI](https://swagger.io/specification/): machine-readable API contract and interactive documentation.

## Known Limitations and Future Improvements

- Authentication currently uses JWT-based plumbing; OAuth2 or refresh-token flows are not implemented.
- Migrations are SQL-based, but a dedicated migration runner is not bundled yet.
- The Docker setup is ready for local development, but production hardening is still needed.
- CI/CD pipelines are not included yet.
- Observability enhancements such as structured logging, metrics, and tracing are still future work.
