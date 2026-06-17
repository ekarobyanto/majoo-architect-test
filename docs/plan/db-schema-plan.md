# Blog System Database Migration Plan

## Overview

This document outlines the database migration strategy for the Blog System API.

### Core Requirements

* User Authentication & Authorization
* CRUD Operations for Posts
* CRUD Operations for Comments
* Role-Based Access Control (RBAC)
* UUID-based Primary Keys
* PostgreSQL Database
* Transaction Support
* Scalable and Maintainable Schema Design

---

# Migration Order

```text
000001_enable_extensions

000002_create_users
000003_create_roles
000004_create_user_roles
000005_seed_roles

000006_create_posts
000007_create_comments

000008_create_posts_indexes
000009_create_comments_indexes
```

---

# Phase 0 — Database Extensions

## Migration

```text
000001_enable_extensions
```

### Purpose

Enable PostgreSQL extensions required by the application.

### Extension

```sql
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
```

### Reason

Provides:

* `gen_random_uuid()`

Used for UUID primary key generation.

---

# Phase 1 — Identity & Access Management

## Migration

```text
000002_create_users
```

### Table: users

| Column        | Type         | Constraints                   |
| ------------- | ------------ | ----------------------------- |
| id            | UUID         | PK, DEFAULT gen_random_uuid() |
| username      | VARCHAR(50)  | UNIQUE, NOT NULL              |
| email         | VARCHAR(255) | UNIQUE, NOT NULL              |
| password_hash | TEXT         | NOT NULL                      |
| created_at    | TIMESTAMP    | NOT NULL                      |
| updated_at    | TIMESTAMP    | NOT NULL                      |

### Notes

* Username must be unique.
* Email must be unique.
* Passwords are stored as hashes only.

---

## Migration

```text
000003_create_roles
```

### Table: roles

| Column     | Type        | Constraints                   |
| ---------- | ----------- | ----------------------------- |
| id         | UUID        | PK, DEFAULT gen_random_uuid() |
| name       | VARCHAR(50) | UNIQUE, NOT NULL              |
| created_at | TIMESTAMP   | NOT NULL                      |

### Notes

Stores application roles.

Examples:

* admin
* user

---

## Migration

```text
000004_create_user_roles
```

### Table: user_roles

| Column      | Type      | Constraints  |
| ----------- | --------- | ------------ |
| user_id     | UUID      | FK users(id) |
| role_id     | UUID      | FK roles(id) |
| assigned_at | TIMESTAMP | NOT NULL     |

### Constraints

```text
PRIMARY KEY (user_id, role_id)
```

### Relationships

```text
users
  |
  +----< user_roles >----+
                         |
                      roles
```

### Notes

Supports many-to-many role assignments.

---

## Migration

```text
000005_seed_roles
```

### Seed Data

```sql
INSERT INTO roles(name)
VALUES
('admin'),
('user');
```

---

# Phase 2 — Blog Domain

## Migration

```text
000006_create_posts
```

### Table: posts

| Column     | Type         | Constraints                   |
| ---------- | ------------ | ----------------------------- |
| id         | UUID         | PK, DEFAULT gen_random_uuid() |
| author_id  | UUID         | FK users(id)                  |
| title      | VARCHAR(255) | NOT NULL                      |
| content    | TEXT         | NOT NULL                      |
| created_at | TIMESTAMP    | NOT NULL                      |
| updated_at | TIMESTAMP    | NOT NULL                      |

### Relationships

```text
users
  |
  +----< posts
```

### Notes

Represents blog articles created by users.

---

## Migration

```text
000007_create_comments
```

### Table: comments

| Column     | Type      | Constraints                   |
| ---------- | --------- | ----------------------------- |
| id         | UUID      | PK, DEFAULT gen_random_uuid() |
| post_id    | UUID      | FK posts(id)                  |
| author_id  | UUID      | FK users(id)                  |
| content    | TEXT      | NOT NULL                      |
| created_at | TIMESTAMP | NOT NULL                      |
| updated_at | TIMESTAMP | NOT NULL                      |

### Relationships

```text
posts
  |
  +----< comments

users
  |
  +----< comments
```

### Notes

Represents comments made by users on blog posts.

---

# Foreign Key Strategy

## User Deletion

### Posts

```sql
ON DELETE CASCADE
```

### Comments

```sql
ON DELETE CASCADE
```

### Reason

When a user is deleted:

* All authored posts are removed.
* All authored comments are removed.

---

## Post Deletion

### Comments

```sql
ON DELETE CASCADE
```

### Reason

When a post is deleted:

* All related comments are automatically removed.

---

# Phase 3 — Performance Optimization

## Migration

```text
000008_create_posts_indexes
```

### Indexes

```sql
CREATE INDEX idx_posts_author_id
ON posts(author_id);

CREATE INDEX idx_posts_created_at
ON posts(created_at DESC);
```

### Purpose

Optimize:

* Post listing
* Author filtering
* Recent post retrieval

---

## Migration

```text
000009_create_comments_indexes
```

### Indexes

```sql
CREATE INDEX idx_comments_post_id
ON comments(post_id);

CREATE INDEX idx_comments_author_id
ON comments(author_id);

CREATE INDEX idx_comments_created_at
ON comments(created_at DESC);
```

### Purpose

Optimize:

* Comment retrieval by post
* Author filtering
* Recent comment retrieval

---

# Entity Relationship Diagram

```text
users
  │
  ├──< user_roles >── roles
  │
  ├──< posts
  │        │
  │        └──< comments
  │
  └──< comments
```

---

# Authorization Model

## Roles

### admin

Permissions:

* Manage all posts
* Manage all comments
* Manage users
* Assign roles

### user

Permissions:

* Create posts
* Update own posts
* Delete own posts
* Create comments
* Update own comments
* Delete own comments

---

# Ownership Rules

## Posts

```text
post.author_id == current_user.id
```

Required for:

* Update Post
* Delete Post

Unless user has:

```text
admin
```

role.

---

## Comments

```text
comment.author_id == current_user.id
```

Required for:

* Update Comment
* Delete Comment

Unless user has:

```text
admin
```

role.
