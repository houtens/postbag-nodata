-- name: CreateTitle :one
INSERT INTO titles (
    name
) VALUES (
    $1
) RETURNING *;

-- name: GetTitle :one
SELECT * FROM titles
WHERE id = $1 LIMIT 1;

-- name: ListTitles :many
SELECT * FROM titles
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteTitle :exec
DELETE FROM titles where id = $1;

-- name: TruncateTitles :exec
truncate titles cascade;
