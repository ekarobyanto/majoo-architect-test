# Spec: Project Rename to Simple Blog

Rename project from `github.com/user/simple-blog` to `github.com/user/simple-blog`.

## Goal
Update all module references and imports to reflect the new project name.

## Approach
1. Update `go.mod` module name.
2. Global search and replace of `github.com/user/simple-blog` with `github.com/user/simple-blog`.
3. Update `README.md` title and descriptions.
4. Regenerate dependency injection code using `wire`.
5. Regenerate Swagger documentation using `swag`.
6. Verify with tests.

## Files to Modify
- `go.mod`
- `README.md`
- All Go source files with imports.
- Swagger documentation files.
- Project plans and documentation.

## Verification
- `go mod tidy`
- `make build`
- `make test`
- `make swagger`
