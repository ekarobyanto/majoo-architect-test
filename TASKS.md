# Blog System Implementation Tasks

* [ ] Task 0: Foundation

  * [ ] Define API response contract
  * [ ] Define AppError structure
  * [ ] Define error codes
  * [ ] Implement global error handler
  * [ ] Define validation error response format
* [ ] Define common response helpers
  * [ ] Generate documentation

* [x] Task 1: Authentication Module - Register

  * [x] Define interfaces and DTOs in `internal/modules/auth/domain`
  * [x] Implement Repository

    * [x] Create user
    * [x] Get role by name
    * [x] Assign role to user
    * [x] Check username existence
    * [x] Check email existence
  * [x] Implement Service

    * [x] Validate uniqueness
    * [x] Hash password
    * [x] Create user transaction
    * [x] Assign default role
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Write Ginkgo tests

    * [x] Unit tests
    * [x] Integration tests
  * [x] Generate documentation

* [x] Task 2: Authentication Module - Login

  * [x] Define Login DTOs
  * [x] Implement Repository

    * [x] Get user by email
    * [x] Get roles by user id
  * [x] Implement Service

    * [x] Verify password
    * [x] Load user roles
    * [x] Generate JWT
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Write Ginkgo tests

    * [x] Unit tests
    * [x] Integration tests
  * [x] Generate documentation

* [x] Task 3: Authentication Middleware

  * [x] Implement JWT middleware
  * [x] Extract bearer token
  * [x] Validate JWT signature
  * [x] Validate JWT expiration
  * [x] Extract user claims
  * [x] Inject user id into request context
  * [x] Inject user roles into request context
  * [x] Write tests
  * [x] Generate documentation

* [x] Task 4: Authorization

  * [x] Implement role-based authorization
  * [x] Implement ownership-based authorization
  * [x] Implement post ownership verification
  * [x] Implement comment ownership verification
  * [x] Implement admin access verification
  * [x] Write tests
  * [x] Generate documentation

* [x] Task 5: Posts Module

  * [x] Define DTOs
  * [x] Implement Repository
  * [x] Implement Service
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Integrate authorization
  * [x] Implement pagination
  * [x] Implement delete transaction
  * [x] Write Ginkgo tests
  * [x] Generate documentation

* [x] Task 6: Comments Module

  * [x] Define DTOs
  * [x] Implement Repository
  * [x] Implement Service
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Integrate authorization
  * [x] Write Ginkgo tests
  * [x] Generate documentation

* [ ] Task 7: Final Integration

  * [ ] Register all routes
  * [ ] Verify endpoint specifications
  * [ ] Verify authentication flow
  * [ ] Verify authorization flow
  * [ ] Verify validation responses
  * [ ] Verify transaction handling
  * [ ] Verify pagination responses
  * [ ] Generate OpenAPI/Swagger documentation
  * [ ] Run end-to-end tests
  * [ ] Final review

## API Endpoints

### Authentication

* [x] POST `/auth/register`
* [x] POST `/auth/login`

### Posts

* [x] POST `/posts`
* [x] GET `/posts`
* [x] GET `/posts/{id}`
* [x] PUT `/posts/{id}`
* [x] DELETE `/posts/{id}`

### Comments

* [x] POST `/posts/{id}/comments`
* [x] PUT `/comments/{id}`
* [x] DELETE `/comments/{id}`
