-- name: SeedMembership :one
INSERT INTO memberships (
    user_id,
    cost,
    membership_type_id,
    payment_type_id,
    expires_at,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: CreateMembership :one
INSERT INTO memberships (
    user_id,
    cost,
    membership_type_id,
    payment_type_id,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetMembership :one
SELECT * FROM memberships WHERE id = $1 LIMIT 1;

-- name: ListMemberships :many
SELECT * FROM memberships ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateMembership :one
UPDATE memberships
SET
    user_id = $2,
    cost = $3,
    membership_type_id = $4,
    payment_type_id = $5,
    expires_at = $6
WHERE
    id = $1
RETURNING *;

-- name: DeleteMembership :exec
DELETE FROM memberships WHERE id = $1;

-- name: GetValidMembership :one
SELECT CASE WHEN expires_at > NOW() THEN 1 ELSE 0 END FROM memberships WHERE user_id = $1 ORDER BY expires_at DESC LIMIT 1;

-- name: TruncateMemberships :exec
truncate memberships cascade;

