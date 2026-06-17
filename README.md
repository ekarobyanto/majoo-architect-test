# Go Backend Boilerplate

A production-ready Go backend boilerplate using a modular (feature-based) architecture, designed for performance, maintainability, and ease of use.

## Tech Stack

- **Framework**: [Fiber](https://gofiber.io/) - High-performance Express-inspired web framework for Go.
- **Database**: PostgreSQL - Reliable and robust relational database.
- **Query Builder**: [sqlx](https://github.com/jmoiron/sqlx) - Lightweight extensions to `database/sql` for simpler data mapping.
- **Testing**: [Ginkgo](https://onsi.github.io/ginkgo/) & [Gomega](https://onsi.github.io/gomega/) - BDD-style testing framework for clear and expressive tests.
- **Config**: [Viper](https://github.com/spf13/viper) - Comprehensive configuration solution.

## Architecture

The project follows a **Modular (Feature-based) Architecture**. Instead of grouping by technical layers (handlers, services, repositories) at the top level, code is organized by business features (modules).

### Directory Structure

- `cmd/api/`: Application entry point.
- `config/`: Configuration loading and structures.
- `database/`: Database migrations and seeds.
  - `migrations/`: SQL migration files.
  - `seeds/`: Database seed scripts.
- `internal/modules/`: Business logic organized by feature (e.g., `health`).
  - `domain/`: Interfaces, DTOs, and domain models.
  - `handler/`: HTTP handlers and routing.
  - `service/`: Business logic implementation.
  - `repository/`: Data access layer.
- `internal/platform/`: Infrastructure and shared components (database, server).
- `docs/`: Project documentation and AI guides.
- `tests/`: Integration tests.

## Getting Started

### Prerequisites

- Go 1.25 or higher
- PostgreSQL

### Setup

1. **Download dependencies**:
   ```bash
   go mod download
   ```

2. **Environment Configuration**:
   Copy `.env.example` to `.env` and update the database credentials.
   ```bash
   cp .env.example .env
   ```

3. **Running the Application**:
   ```bash
   go run cmd/api/main.go
   ```

## Testing

The project uses Ginkgo for BDD-style testing.

### Run all tests
```bash
ginkgo ./...
```

### Run tests with Go tool
```bash
go test ./...
```

### Test Coverage
To view test coverage, run:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Test Report

Current test status as of 2026-06-17:

| Component | Status | Test Type |
|-----------|--------|-----------|
| Health Module | ✅ Passing | Unit / Handler |
| Database Connection | ✅ Passing | Integration |
| Server Initialization | ✅ Passing | Integration |
| Configuration Loader | ✅ Passing | Unit |

Total tests: 13 passed, 0 failed.

## Justification of Technology Choices

- **Fiber**: Chosen for its extreme performance and low memory footprint. Its Express-like API makes it intuitive for developers coming from other ecosystems.
- **sqlx**: Provides a good balance between raw SQL control and the convenience of mapping rows to structs without the overhead of a full ORM.
- **PostgreSQL**: The industry standard for relational data, offering strong consistency and advanced features.
- **Ginkgo**: Enables Behavioral-Driven Development (BDD), making tests more readable and better documented as specifications.

## Limitations and Future Improvements

- **Authentication**: JWT or OAuth2 integration is planned.
- **Migrations**: Database migration tool (e.g., `golang-migrate`) integration.
- **Docker Support**: Containerization for easier deployment.
- **CI/CD**: GitHub Actions or GitLab CI pipelines.
