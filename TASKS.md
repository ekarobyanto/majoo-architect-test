# Blog System Implementation Tasks

* [x] Task 0: Foundation

  * [x] Define API response contract
  * [x] Define AppError structure
  * [x] Define error codes
  * [x] Implement global error handler
  * [x] Define validation error response format
  * [x] Define common response helpers
  * [x] Generate documentation

* [x] Task 1: Authentication Module - Register

  * [x] Define interfaces and DTOs in `internal/modules/auth/domain`
  * [x] Implement Repository
  * [x] Implement Service
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Write Ginkgo tests
  * [x] Generate documentation

* [x] Task 2: Authentication Module - Login

  * [x] Define Login DTOs
  * [x] Implement Repository
  * [x] Implement Service
  * [x] Implement Handler
  * [x] Add request validation
  * [x] Write Ginkgo tests
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

* [x] Task 7: Final Integration

  * [x] Register all routes
  * [x] Verify endpoint specifications
  * [x] Verify authentication flow
  * [x] Verify authorization flow
  * [x] Verify validation responses
  * [x] Verify transaction handling
  * [x] Verify pagination responses
  * [x] Generate OpenAPI/Swagger documentation
  * [x] Run end-to-end tests
  * [x] Final review

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
