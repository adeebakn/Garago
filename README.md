# Garago

Garago is a business management application for automobile accessory, customization, and service shops. It combines a Go REST API with a React dashboard for managing the full customer-to-payment workflow.

## Features

- Cookie-based authentication with Redis-backed sessions
- Password reset by email
- Customer and vehicle management
- Service catalog
- Quotations with line items and approval status
- Work orders with lifecycle tracking
- Invoices, payments, balances, and payment receipts
- Dashboard reporting
- Organization profile and document branding
- PDF-ready quotation and invoice views
- Soft deletion for business records

## Stack

- **Backend:** Go, Gin, pgx, SQLC
- **Frontend:** React 19, Vite, React Router, Axios
- **Data:** PostgreSQL 17 and Redis 7
- **Infrastructure:** Docker Compose

## Prerequisites

- Go 1.26.4 or newer
- Node.js and npm compatible with Vite 8
- Docker with Docker Compose
- An SMTP account for password-reset email

## Local setup

### 1. Start PostgreSQL and Redis

```powershell
docker compose up -d
```

### 2. Configure the backend

Copy the example environment file and replace the placeholder secrets and SMTP credentials:

```powershell
Copy-Item .env.example .env
```

The API requires these settings:

| Variable | Purpose | Local default/example |
| --- | --- | --- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://garago:garago@localhost:5432/garago` |
| `JWT_SECRET` | Session signing key (at least 32 characters) | Replace the example value |
| `JWT_EXPIRY_HOURS` | Session lifetime | `24` |
| `REDIS_URL` | Redis connection string | `redis://localhost:6379/0` |
| `COOKIE_SECURE` | Send the session cookie over HTTPS only | `false` for local development |
| `FRONTEND_URL` | Allowed CORS origin and password-reset URL base | `http://localhost:5173` |
| `SERVER_ADDRESS` | API listen address | `:8080` |
| `SMTP_HOST`, `SMTP_PORT` | SMTP server | `smtp.gmail.com`, `587` |
| `SMTP_USERNAME`, `SMTP_PASSWORD` | SMTP credentials | Use an app password where required |
| `SMTP_FROM`, `SMTP_FROM_NAME` | Sender address and display name | Your sender address, `Garago` |

The backend loads `.env` automatically. All SMTP values are currently required at startup.

### 3. Apply database migrations

Run every SQL file in `migrations/` in numeric order. With the Compose database and PowerShell:

```powershell
Get-ChildItem .\migrations\*.sql |
  Sort-Object Name |
  ForEach-Object {
    Get-Content $_.FullName -Raw |
      docker compose exec -T postgres psql -U garago -d garago -v ON_ERROR_STOP=1
  }
```

Migrations are not run automatically by the application.

### 4. Run the API

```powershell
go run .
```

The explicit production entry point is also available with `go run ./cmd/api`.

The API is available at `http://localhost:8080/api`. Verify it with `GET /api/health`.

### 5. Run the frontend

In a second terminal:

```powershell
Set-Location frontend
Copy-Item .env.example .env
npm install
npm run dev
```

Open `http://localhost:5173`. Set `VITE_API_BASE_URL` in `frontend/.env` if the API is hosted somewhere other than `http://localhost:8080/api`.

Create the first administrator through the sign-up screen, then sign in to access the dashboard.

## Development commands

### Backend

```powershell
go test ./...
go build ./...
sqlc generate
```

Run `sqlc generate` after changing files in `queries/` or the database schema used by SQLC.

### Frontend

```powershell
Set-Location frontend
npm run lint
npm run build
npm run preview
```

## Application workflow

```text
Customer + Vehicle
        |
        v
    Quotation
        |
        v
 Customer approval
        |
        v
    Work order
        |
        v
      Invoice
        |
        v
      Payment
```

### History and deletion rules

Garago preserves the identity attached to operational and financial records:

- A customer can be deleted only when no vehicle or business-history record has ever been linked to that customer.
- A vehicle can be deleted only when no quotation or work order has ever been linked to that vehicle.
- Archived quotations, work orders, invoices, and payments still count as history.
- Delete requests that would break this audit trail return HTTP `409 Conflict` and leave the record unchanged.

These rules are enforced by both the backend and PostgreSQL triggers; the frontend confirmation text is informational and is not the security boundary. This prevents older application processes, maintenance scripts, or direct SQL from bypassing audit-history protection.

## API overview

Public endpoints include health, organization details, sign-up, sign-in, and password reset. Authenticated endpoints cover:

- Customers and vehicles
- Services
- Quotations and quotation status
- Work orders and work-order status
- Invoices and payments
- Dashboard data
- Organization settings

Authentication uses the HTTP-only `garago_session` cookie, so browser requests must include credentials. The frontend's Axios client is configured accordingly.

## Project structure

```text
Garago/
├── cmd/api/              # Production API executable
├── frontend/             # React/Vite application
│   └── src/
│       ├── auth/         # Authentication screens and state
│       ├── components/   # Lazy-loaded dashboard and business modules
│       ├── hooks/        # Reusable lifecycle hooks
│       ├── organization/ # Organization state and provider
│       ├── ui/           # Shared interaction state such as confirmation dialogs
│       └── utils/        # Shared API-error and pagination utilities
├── internal/
│   ├── config/           # Environment loading and startup validation
│   ├── apis/             # HTTP server, routing, middleware, and feature handlers
│   │   ├── server.go     # Dependency lifecycle and route registration
│   │   ├── auth.go       # Authentication and session endpoints
│   │   ├── *_handler.go  # One handler module per business capability
│   │   ├── *_types.go    # Transport request/response types by capability
│   │   └── http_helpers.go
│   ├── database/         # PostgreSQL connection lifecycle
│   │   └── postgres/     # Generated PostgreSQL query code
│   └── platform/
│       └── mailer/       # SMTP password-reset delivery
├── migrations/           # Ordered PostgreSQL migrations
├── queries/              # SQLC query definitions
├── docker-compose.yml    # PostgreSQL and Redis services
├── sqlc.yaml
└── go.mod
```

## Acceptance smoke test

1. Sign up and sign in.
2. Configure the organization profile.
3. Create a customer, vehicle, and service.
4. Create a quotation with one or more service lines and approve it.
5. Create a work order from the approved quotation and progress its status.
6. Create an invoice for the work order. A second active invoice for the same work order should be rejected.
7. Record partial and final payments and confirm that the invoice moves from `Unpaid` to `Partial` to `Paid`.
