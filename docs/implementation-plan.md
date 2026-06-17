# Blog System API - Implementation Plan

## Phase 1 - Authentication Module

### Register

#### Validation

* Username required
* Username length 3-50 characters
* Email required
* Valid email format
* Password required
* Password minimum 8 characters

#### Business Flow

* Check username uniqueness
* Check email uniqueness
* Hash password
* Create user
* Assign default role
* Commit transaction

#### Repository Requirements

* Create user
* Get role by name
* Assign role to user

#### Transaction

* Create user
* Assign user role

---

### Login

#### Validation

* Email required
* Password required

#### Business Flow

* Find user by email
* Verify password
* Load user roles
* Generate JWT token

#### Repository Requirements

* Get user by email
* Get roles by user id

---

## Phase 2 - Authentication Middleware

### JWT Middleware

* Extract bearer token
* Validate JWT signature
* Validate JWT expiration
* Extract user claims
* Inject user context into request

### User Context

* User ID
* User Roles

---

## Phase 3 - Authorization

### Post Authorization

* Owner can update post
* Owner can delete post
* Admin can update any post
* Admin can delete any post

### Comment Authorization

* Owner can update comment
* Owner can delete comment
* Admin can update any comment
* Admin can delete any comment

---

## Phase 4 - Posts Module

### Create Post

#### Validation

* Title required
* Title maximum 255 characters
* Content required

#### Business Flow

* Validate request
* Create post

---

### Get Posts

#### Features

* Pagination

#### Query Parameters

* page
* limit

---

### Get Post Detail

#### Features

* Post information
* Author information
* Comments list

---

### Update Post

#### Validation

* Title required
* Content required

#### Business Flow

* Find post
* Verify ownership
* Update post

---

### Delete Post

#### Business Flow

* Find post
* Verify ownership
* Delete comments
* Delete post

#### Transaction

* Delete comments
* Delete post

---

## Phase 5 - Comments Module

### Create Comment

#### Validation

* Content required
* Content maximum 1000 characters

#### Business Flow

* Verify post exists
* Create comment

---

### Update Comment

#### Validation

* Content required
* Content maximum 1000 characters

#### Business Flow

* Find comment
* Verify ownership
* Update comment

---

### Delete Comment

#### Business Flow

* Find comment
* Verify ownership
* Delete comment

---

## Phase 6 - Validation

### User Registration

* Username required
* Username length 3-50 characters
* Email required
* Valid email format
* Password required
* Password minimum 8 characters

### User Login

* Email required
* Password required

### Post

* Title required
* Title maximum 255 characters
* Content required

### Comment

* Content required
* Content maximum 1000 characters

---

## Phase 7 - Error Handling

### Validation Error

* Invalid request payload
* Missing required field
* Invalid field format

### Business Error

* User not found
* Post not found
* Comment not found
* Duplicate username
* Duplicate email

### Authorization Error

* Unauthorized
* Forbidden

### Server Error

* Internal server error

---

## Phase 8 - API Endpoints

### Authentication

* POST /auth/register
* POST /auth/login

### Posts

* POST /posts
* GET /posts
* GET /posts/{id}
* PUT /posts/{id}
* DELETE /posts/{id}

### Comments

* POST /posts/{id}/comments
* PUT /comments/{id}
* DELETE /comments/{id}

---

## Phase 9 - Testing

### Authentication

* Register success
* Register duplicate username
* Register duplicate email
* Login success
* Login invalid password

### Posts

* Create post
* Get posts
* Get post detail
* Update own post
* Update other user's post
* Delete own post
* Delete other user's post

### Comments

* Create comment
* Update own comment
* Update other user's comment
* Delete own comment
* Delete other user's comment

### Authorization

* Access protected endpoint without token
* Access protected endpoint with invalid token
* Access admin functionality as user
* Access admin functionality as admin

### Transactions

* Register user transaction rollback
* Delete post transaction rollback
