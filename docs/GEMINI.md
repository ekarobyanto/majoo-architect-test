# AI Agent Guide (GEMINI)

This guide provides context for Gemini (and other AI agents) to understand the architecture and conventions of this Go backend boilerplate.

## Core Mandates

- **Modular Architecture**: Always follow the feature-based modular structure in `internal/modules/`. Do not create top-level `handlers/`, `services/`, etc.
- **TDD Preference**: Prefer Test-Driven Development. Create/update Ginkgo tests when adding or modifying features.
- **Interface First**: Define interfaces in the `domain/` package of a module before implementing them in `service/` or `repository/`.
- **Fiber & sqlx**: Use Fiber for HTTP and sqlx for database operations. Avoid adding GORM or other heavy ORMs unless explicitly requested.

## Project Structure Reference

```
internal/modules/<feature>/
├── domain/       # Interfaces and DTOs
├── handler/      # Fiber handlers & routing
├── service/      # Business logic implementation
└── repository/   # Data access implementation
```

## Common Tasks

### Adding a New Module
1. Create the directory structure under `internal/modules/`.
2. Define interfaces and models in `domain/`.
3. Implement the Repository and write Ginkgo tests.
4. Implement the Service and write Ginkgo tests.
5. Implement the Handler and write Ginkgo tests.
6. Register the new module's routes in `internal/platform/server/server.go`.

### Database Operations
- Use `sqlx` tags in structs: `` `db:"column_name"` ``.
- Use named queries or standard SQL as appropriate.
- Ensure repositories are tested with `go-sqlmock` or a real test database.

### Testing Conventions
- **Unit Tests**: Place in the same package as the code (e.g., `health_service_test.go` in `service/`).
- **Integration Tests**: Place in `tests/integration/`.
- **BDD Style**: Use `Describe`, `Context`, `It` from Ginkgo.
