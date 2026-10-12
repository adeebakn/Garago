# Garago backend architecture

Garago uses a modular-monolith structure. Packages are organized by responsibility, and implementation packages live under `internal/` so they cannot be imported accidentally by external modules.

## Dependency direction

```text
cmd/api
  └── internal/apis
        ├── internal/config
        ├── internal/database
        │     └── internal/database/postgres
        ├── internal/platform/mailer
```

Dependencies point inward from the executable toward internal implementation packages. Infrastructure packages do not import the HTTP layer.

## Package ownership

- `cmd/api`: process entry point, operating-system signals, and final error reporting.
- `internal/config`: environment loading, defaults, and startup validation.
- `internal/apis`: HTTP server lifecycle, route registration, authentication middleware, transport models, and one handler module per business capability.
- `internal/database`: PostgreSQL connection lifecycle and readiness checks.
- `internal/database/postgres`: generated PostgreSQL query code. Do not edit these files manually.
- `internal/platform/mailer`: SMTP implementation for password-reset delivery.
- `queries`: reviewed SQL used as SQLC input.
- `migrations`: ordered database schema and integrity changes.

## File conventions

- `*_handler.go`: HTTP behavior for one business capability.
- `*_types.go`: request and response transport types for that capability.
- `*_test.go`: tests next to the package behavior they protect.
- Generated files remain isolated under `internal/database/postgres`.

## Commands

```powershell
go run ./cmd/api
go test ./...
go vet ./...
go build ./...
sqlc generate
```

When introducing a new capability, add its route and handler to `internal/apis`, put reusable infrastructure behind a focused package under `internal/platform`, and add database operations to `queries` rather than embedding routine SQL in handlers.
