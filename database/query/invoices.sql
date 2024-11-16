-- name: SeedInvoice :one
INSERT INTO invoices (
    tournament_id,
    num_players,
    num_non_members,
    num_games,
    is_multiday,
    is_overseas,
    levy_cost,
    extras_cost,
    total_cost,
    is_paid,
    description,
    extras_comment,
    comment,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING *;

-- name: CreateInvoice :one
INSERT INTO invoices (
    tournament_id,
    num_players,
    num_non_members,
    num_games,
    is_multiday,
    is_overseas,
    levy_cost,
    extras_cost,
    total_cost,
    is_paid,
    description,
    extras_comment,
    comment
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
) RETURNING *;

-- name: GetInvoice :one
SELECT * FROM invoices WHERE id = $1 LIMIT 1;

-- name: ListInvoices :many
SELECT * FROM invoices ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateInvoice :one
UPDATE invoices
SET
    tournament_id = $2,
    num_players = $3,
    num_non_members = $4,
    num_games = $5,
    is_multiday = $6,
    is_overseas = $7,
    levy_cost = $8,
    extras_cost = $9,
    total_cost = $10,
    is_paid = $11,
    description = $12,
    extras_comment = $13,
    comment = $14
WHERE
  id = $1
RETURNING *;

-- name: DeleteInvoice :exec
DELETE FROM invoices WHERE id = $1;

-- name: TruncateInvoices :exec
truncate invoices cascade;
