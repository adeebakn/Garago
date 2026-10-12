-- name: CreateCustomer :one
INSERT INTO customers (name,phone,email)
VALUES ($1,$2,$3)
RETURNING id, name, phone,email,created_at;

-- name: ListCustomers :many
SELECT id, name, phone, email, created_at
FROM customers
WHERE deleted_at IS NULL
ORDER BY id;

-- name: GetCustomer :one
SELECT id,name,phone,email,created_at
FROM customers
WHERE id=$1 AND deleted_at IS NULL;

-- name: UpdateCustomer :one
UPDATE customers
SET
    name=$1,
    phone=$2,
    email=$3
WHERE id=$4
  AND deleted_at IS NULL
RETURNING id,name,phone,email,created_at;
