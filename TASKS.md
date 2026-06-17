# Blog System Implementation Tasks

* [ ] Task 0: Foundation

  * [ ] Define API response contract
  * [ ] Define AppError structure
  * [ ] Define error codes
  * [ ] Implement global error handler
  * [ ] Define validation error response format
* [ ] Define common response helpers
  * [ ] Generate documentation

* [ ] Task 1: Authentication Module - Register

  * [ ] Define interfaces and DTOs in `internal/modules/auth/domain`
  * [ ] Implement Repository

    * [ ] Create user
    * [ ] Get role by name
    * [ ] Assign role to user
    * [ ] Check username existence
    * [ ] Check email existence
  * [ ] Implement Service

    * [ ] Validate uniqueness
    * [ ] Hash password
    * [ ] Create user transaction
    * [ ] Assign default role
  * [ ] Implement Handler
  * [ ] Add request validation
  * [ ] Write Ginkgo tests

    * [ ] Unit tests
    * [ ] Integration tests
  * [ ] Generate documentation

* [ ] Task 2: Authentication Module - Login

  * [ ] Define Login DTOs
  * [ ] Implement Repository

    * [ ] Get user by email
    * [ ] Get roles by user id
  * [ ] Implement Service

    * [ ] Verify password
    * [ ] Load user roles
    * [ ] Generate JWT
  * [ ] Implement Handler
  * [ ] Add request validation
  * [ ] Write Ginkgo tests

    * [ ] Unit tests
    * [ ] Integration tests
  * [ ] Generate documentation

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

* [ ] POST `/auth/register`
* [ ] POST `/auth/login`

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
