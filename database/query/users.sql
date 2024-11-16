-- name: CreateUser :one
INSERT INTO users (
    first_name,
    last_name,
    alt_name,
    email,
    password_hash,
    absp_num,
    club_id,
    title_id,
    role_id,
    is_deceased,
    x_life,
    x_post,
    x_id,
    avatar
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
) RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY last_name, first_name
LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE users
SET
first_name = $2,
last_name = $3,
alt_name = $4,
email = $5,
password_hash = $6,
absp_num = $7,
club_id = $8,
title_id = $9,
role_id = $10,
is_deceased = $11,
x_life = $12,
x_post = $13,
x_id = $14,
avatar = $15
WHERE id = $1
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
where id = $1;

-- NOTE: custom queries

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- NOTE: seed
-- name: GetUserByXID :one
SELECT * FROM users
WHERE x_id = $1 LIMIT 1;

-- NOTE: seed
-- name: ListUsersByXLife :many
SELECT * FROM users WHERE x_life = 't' ORDER BY id;

-- name: TruncateUsers :exec
truncate users cascade;

-- name: UpdateUserSetToken :one
UPDATE users SET pw_token = $2, pw_token_expiry = NOW() + '15 minutes' WHERE id = $1 RETURNING *;

-- name: GetUserByValidToken :one
SELECT * FROM users WHERE pw_token = $1 AND pw_token_expiry > NOW() LIMIT 1;

-- name: UpdateUserExpireToken :one
UPDATE users SET pw_token_expiry = NOW() WHERE pw_token = $1 RETURNING *;

-- name: UpdateUserPasswordHash :one
UPDATE users SET password_hash = $3 WHERE id = $1 and pw_token = $2 RETURNING *;

-- NOTE: custom query for users used in player dropdowns

-- name: GetActiveUsers :many
SELECT id, first_name, last_name, alt_name FROM users WHERE is_deceased IS false order by last_name;

