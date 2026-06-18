# Project Planning & Architecture

## Architecture Overview

This project uses a **Modular (Feature-based) Architecture**. The goal is to keep code related to a single business feature together, making it easier to understand, test, and maintain.

### Top-Level Structure

- **`cmd/`**: Entry points for the application. Each subdirectory here is a separate binary.
- **`config/`**: Handles application configuration via environment variables.
- **`internal/`**: Private application code. Code here cannot be imported by other projects.
    - **`modules/`**: Contains business features. Each module is self-contained.
    - **`platform/`**: Shared infrastructure code (e.g., database connections, server setup).
- **`pkg/`**: Public utilities that could be shared with other projects (currently empty).
- **`tests/`**: Cross-module integration tests.

## Module Structure

Each module in `internal/modules/` follows a standardized structure to ensure consistency:

### 1. Domain (`domain/`)
- Defines the "what".
- Contains interfaces for Repository and Service.
- Defines Domain Models (structs used across the module).
- Defines DTOs (Data Transfer Objects) for requests and responses.

### 2. Repository (`repository/`)
- Handles "how" data is stored and retrieved.
- Implements the Repository interface defined in Domain.
- Uses `sqlx` for database operations.
- Isolated from business logic and HTTP details.

### 3. Service (`service/`)
- Implements the "business logic".
- Orchestrates calls to one or more Repositories.
- Implements the Service interface defined in Domain.
- Isolated from HTTP details and database implementation details (uses Repository interfaces).

### 4. Handler (`handler/`)
- Handles the "delivery mechanism" (HTTP).
- Defines Fiber routes and maps incoming requests to Service calls.
- Validates input and formats output (JSON).
- Orchestrates the module's dependencies.

## Layer Responsibilities

| Layer | Responsible For | Should NOT know about |
|-------|-----------------|-----------------------|
| **Handler** | Routing, Request Parsing, Validation, Response Formatting | SQL, Transaction management, Business rules complexity |
| **Service** | Business Logic, Validation (Business), Orchestration | HTTP (fiber.Ctx), SQL Queries, Database Drivers |
| **Repository** | SQL Queries, Data Mapping (sqlx), Persistence logic | HTTP, Business rules (should be dumb) |
| **Domain** | Definitions, Interfaces, Models | Implementation details |

## Communication Flow

`Request -> Handler -> Service -> Repository -> Database`

Modules communicate with each other through interfaces to avoid tight coupling. If Module A needs something from Module B, it should ideally go through Module B's Service interface.
