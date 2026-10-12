-- name: CreateService :one
INSERT INTO services (name, description, price)
VALUES ($1, $2, $3)
RETURNING id, name, description, price, created_at;


-- name: ListServices :many
SELECT id, name, description, price, created_at
FROM services
ORDER BY id;


-- name: GetService :one
SELECT id, name, description, price, created_at
FROM services
WHERE id = $1;


-- name: UpdateService :one
UPDATE services
SET
    name = $1,
    description = $2,
    price = $3
WHERE id = $4
RETURNING id, name, description, price, created_at;


-- name: DeleteService :exec
DELETE FROM services
WHERE id = $1;