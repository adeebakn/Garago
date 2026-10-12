# Garago

Garago is a cloud-based business management system designed for automobile accessories, customization, and service businesses.

## Features

- Customer management
- Vehicle management
- Service management
- Quotation management
- Work order management
- Invoice management
- Payment tracking
- Dashboard and reporting
- PostgreSQL database
- REST APIs

## Technology Stack

- Go
- Gin
- PostgreSQL
- SQLC
- Docker
- Git & GitHub

## Backend Architecture

Client / React Frontend
        ↓
     Gin API
        ↓
      Go
        ↓
      SQLC
        ↓
   PostgreSQL

## Business Workflow

Customer
   ↓
Quotation
   ↓
Customer Approval
   ↓
Work Order
   ↓
Invoice
   ↓
Payment

## API Modules

- Customers
- Vehicles
- Services
- Quotations
- Work Orders
- Invoices
- Payments
- Dashboard

## Running the Project

### 1. Start PostgreSQL

docker compose up -d

### 2. Set the Database URL

Create a .env file:

DATABASE_URL=postgres://autoerp:autoerp@localhost:5432/autoerp

### 3. Run the Application

go run .

The API runs on:

http://localhost:8080

## Project Structure

AutoERP/
├── db/sqlc/          # SQLC generated code
├── migrations/       # Database migrations
├── queries/          # SQL queries
├── *_handler.go      # HTTP handlers
├── *.go              # Application logic
├── docker-compose.yml
├── sqlc.yaml
├── go.mod
└── go.sum

## Status

The admin backend and React/Vite frontend are in active development. The
customer-facing portal and authentication are deferred.

## Database migrations

Apply the SQL files in `migrations/` in numeric order using a PostgreSQL
migration tool or `psql`. When a file in `queries/` changes, run
`sqlc generate` from the repository root to refresh `db/sqlc`.

## Frontend

```powershell
cd frontend
Copy-Item .env.example .env
npm install
npm run dev
```

Set `VITE_API_BASE_URL` in `frontend/.env` when the API is not running at
`http://localhost:8080/api`.

## Manual acceptance checklist

1. Create a customer, vehicle, and service.
2. Create a quotation with one or more service lines, then approve it.
3. Create a work order from that approved quotation.
4. Create one invoice for the work order; a second attempt must return a
   conflict.
5. Record an advance and a final payment; verify the invoice changes from
   `Unpaid` to `Partial` and then `Paid`.
