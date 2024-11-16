-- name: CreateClub :one
INSERT INTO clubs (
    name,
    county,
    website,
    is_active,
    phone,
    email,
    contact_name,
    country_id,
    x_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING *;

-- name: GetClub :one
SELECT * FROM clubs WHERE id = $1 LIMIT 1;

-- name: ListClubs :many
SELECT * FROM clubs ORDER BY name LIMIT $1 OFFSET $2;

-- name: UpdateClub :one
UPDATE clubs
SET
  name = $2,
  county = $3,
  website = $4,
  is_active = $5,
  phone = $6,
  email = $7,
  contact_name = $8,
  country_id = $9,
  x_id = $10,
  updated_at = now()
WHERE
  id = $1
RETURNING *;

-- name: DeleteClub :exec
DELETE FROM clubs WHERE id = $1;


-- name: ListClubsWithCountry :many
SELECT c.id, c.name, c.county, r.name as country
FROM clubs c JOIN countries r ON c.country_id = r.id 
ORDER BY c.name 
OFFSET $1
LIMIT $2;

-- name: GetClubWithCountry :one
SELECT * FROM clubs join countries on clubs.country_id = countries.id
WHERE clubs.id = $1 LIMIT 1;

-- name: TruncateClubs :exec
truncate clubs cascade;
