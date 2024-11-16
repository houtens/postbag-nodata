-- name: CreateCountry :one
INSERT INTO countries (
    name,
    flag,
    code,
    priority,
    x_id
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetCountry :one
SELECT * FROM countries
WHERE id = $1 LIMIT 1;

-- name: ListCountries :many
SELECT * FROM countries
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateCountry :one
UPDATE countries
SET flag = $2, code = $3, priority = $4, x_id = $5
WHERE id = $1
RETURNING *;

-- name: DeleteCountry :exec
DELETE FROM countries
WHERE id = $1;

-- name: GetCountryByName :one
SELECT * FROM countries
WHERE name = $1 LIMIT 1;

-- name: GetCountryByXID :one
SELECT * FROM countries 
WHERE x_id = $1 LIMIT 1;

-- name: TruncateCountries :exec
truncate countries cascade;
