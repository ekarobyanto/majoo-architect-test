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

* [ ] Task 3: Authentication Middleware

  * [ ] Implement JWT middleware
  * [ ] Extract bearer token
  * [ ] Validate JWT signature
  * [ ] Validate JWT expiration
  * [ ] Extract user claims
  * [ ] Inject user id into request context
  * [ ] Inject user roles into request context
  * [ ] Write tests
  * [ ] Generate documentation

* [ ] Task 4: Authorization

  * [ ] Implement role-based authorization
  * [ ] Implement ownership-based authorization
  * [ ] Implement post ownership verification
  * [ ] Implement comment ownership verification
  * [ ] Implement admin access verification
  * [ ] Write tests
  * [ ] Generate documentation

* [ ] Task 5: Posts Module

  * [ ] Define DTOs
  * [ ] Implement Repository

    * [ ] Create post
    * [ ] Get post by id
    * [ ] Get paginated posts
    * [ ] Update post
    * [ ] Delete post
  * [ ] Implement Service

    * [ ] Create post
    * [ ] Get posts
    * [ ] Get post detail
    * [ ] Update post
    * [ ] Delete post
  * [ ] Implement Handler
  * [ ] Add request validation
  * [ ] Integrate authorization
  * [ ] Implement pagination
  * [ ] Implement delete transaction
  * [ ] Write Ginkgo tests

    * [ ] Unit tests
    * [ ] Integration tests
  * [ ] Generate documentation

* [ ] Task 6: Comments Module

  * [ ] Define DTOs
  * [ ] Implement Repository

    * [ ] Create comment
    * [ ] Get comment by id
    * [ ] Update comment
    * [ ] Delete comment
  * [ ] Implement Service

    * [ ] Verify post exists
    * [ ] Create comment
    * [ ] Update comment
    * [ ] Delete comment
  * [ ] Implement Handler
  * [ ] Add request validation
  * [ ] Integrate authorization
  * [ ] Write Ginkgo tests

    * [ ] Unit tests
    * [ ] Integration tests
  * [ ] Generate documentation

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

* [ ] POST `/posts`
* [ ] GET `/posts`
* [ ] GET `/posts/{id}`
* [ ] PUT `/posts/{id}`
* [ ] DELETE `/posts/{id}`

### Comments

* [ ] POST `/posts/{id}/comments`
* [ ] PUT `/comments/{id}`
* [ ] DELETE `/comments/{id}`
