-- name: CreatePaymentType :one
INSERT INTO payment_types (name) VALUES ($1) RETURNING *;

-- name: GetPaymentType :one
SELECT * FROM payment_types WHERE id = $1 LIMIT 1;

-- name: ListPaymentTypes :many
SELECT * FROM payment_types
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePaymentType :one
UPDATE payment_types
SET name = $2 WHERE id = $1 RETURNING *;

-- name: DeletePaymentType :exec
DELETE FROM payment_types WHERE id = $1;

-- name: GetPaymentTypeByName :one
SELECT * FROM payment_types WHERE name = $1 LIMIT 1;

-- name: TruncatePaymentTypes :exec
truncate payment_types cascade;

