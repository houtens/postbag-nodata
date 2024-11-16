-- name: CreateContact :one
INSERT INTO contacts (
    user_id,
    address1,
    address2,
    address3,
    address4,
    postcode,
    country_id,
    phone,
    mobile,
    notes
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetContact :one
SELECT * FROM contacts
WHERE id = $1 LIMIT 1;

-- name: GetContactByUserID :one
SELECT * FROM contacts WHERE user_id = $1 LIMIT 1;

-- name: ListContacts :many
SELECT * FROM contacts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateContact :one
UPDATE contacts
SET 
user_id = $2,
address1 = $3,
address2 = $4,
address3 = $5,
address4 = $6,
postcode = $7,
country_id = $8,
phone = $9,
mobile = $10,
notes = $11
WHERE id = $1
RETURNING *;

-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1;

-- name: TruncateContacts :exec
truncate contacts cascade;

