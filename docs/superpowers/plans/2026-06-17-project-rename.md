# Project Rename Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rename project from `github.com/user/simple-blog` to `github.com/user/simple-blog`.

**Architecture:** Global string replacement of module path in code, configuration, and documentation.

**Tech Stack:** Go, Wire, Swag.

---

### Task 1: Update go.mod

**Files:**
- Modify: `go.mod`

- [ ] **Step 1: Update module path in go.mod**

```go
module github.com/user/simple-blog
```

- [ ] **Step 2: Run go mod tidy**

Run: `go mod tidy`
Expected: `go.sum` updated, no errors.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: rename module to github.com/user/simple-blog"
```

### Task 2: Global Search and Replace Imports

**Files:**
- Modify: All `.go` files in `cmd/`, `internal/`, `config/`, `models/`, `tests/`.

- [ ] **Step 1: Replace imports using sed**

Run:
```bash
find . -type f -name "*.go" -not -path "./vendor/*" -exec sed -i 's|github.com/user/simple-blog|github.com/user/simple-blog|g' {} +
```

- [ ] **Step 2: Verify replacements**

Run: `grep -r "simple-blog" .`
Expected: No matches in Go files.

- [ ] **Step 3: Commit**

```bash
git add .
git commit -m "refactor: update all imports to new module path"
```

### Task 3: Update Documentation and Plans

**Files:**
- Modify: `README.md`, `docs/plan/`, `docs/swagger/`, `Makefile`.

- [ ] **Step 1: Replace project name in documentation**

Run:
```bash
find . -type f \( -name "*.md" -o -name "*.json" -o -name "*.yaml" -o -name "Makefile" \) -not -path "./.git/*" -exec sed -i 's|simple-blog|simple-blog|g' {} +
find . -type f \( -name "*.md" -o -name "*.json" -o -name "*.yaml" -o -name "Makefile" \) -not -path "./.git/*" -exec sed -i 's|Simple Blog|Simple Blog|g' {} +
```

- [ ] **Step 2: Commit**

```bash
git add .
git commit -m "docs: update project name in documentation"
```

### Task 4: Regenerate Code and Docs

**Files:**
- Modify: `internal/platform/di/wire_gen.go`, `docs/swagger/*`

- [ ] **Step 1: Regenerate Wire DI**

Run: `wire ./internal/platform/di`
Expected: `internal/platform/di/wire_gen.go` updated.

- [ ] **Step 2: Regenerate Swagger**

Run: `make swagger`
Expected: `docs/swagger/` updated.

- [ ] **Step 3: Commit**

```bash
git add .
git commit -m "chore: regenerate wire and swagger docs"
```

### Task 5: Final Validation

- [ ] **Step 1: Run tests**

Run: `go test ./...`
Expected: All tests PASS.

- [ ] **Step 2: Build app**

Run: `make build`
Expected: Binary `bin/api` created successfully.
