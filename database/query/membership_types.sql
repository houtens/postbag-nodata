-- name: CreateMembershipType :one
INSERT INTO membership_types (
    name,
    code,
    num_years,
    is_junior,
    is_post,
    is_life
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: GetMembershipType :one
SELECT * FROM membership_types WHERE id = $1 LIMIT 1;

-- name: GetMembershipTypeByName :one
SELECT * FROM membership_types WHERE name = $1 LIMIT 1;

-- name: ListMembershipTypes :many
SELECT * FROM membership_types ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateMembershipType :one
UPDATE membership_types
SET
    name = $2,
    code = $3,
    num_years = $4,
    is_junior = $5,
    is_post = $6,
    is_life = $7
WHERE
    id = $1
RETURNING *;

-- name: DeleteMembershipType :exec
DELETE FROM membership_types WHERE id = $1;


-- name: TruncateMembershipTypes :exec
truncate membership_types cascade;
