-- name: CreateAuthRole :one
 INSERT INTO auth_roles (
    name,
    can_login,
    is_guest,
    is_members_admin,
    is_clubs_admin,
    is_ratings_admin,
    is_tournaments_admin,
    is_super_admin
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetAuthRole :one
SELECT * FROM auth_roles
WHERE id = $1 LIMIT 1;

-- name: ListAuthRoles :many
SELECT * FROM auth_roles
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAuthRole :one
UPDATE auth_roles
SET
name = $2,
can_login = $3,
is_guest = $4,
is_members_admin = $5,
is_clubs_admin = $6,
is_ratings_admin = $7,
is_tournaments_admin = $8,
is_super_admin = $9
WHERE id = $1 RETURNING *;

-- name: DeleteAuthRole :exec
DELETE FROM auth_roles
WHERE id = $1;

-- name: TruncateAuthRole :exec
truncate auth_roles cascade;
