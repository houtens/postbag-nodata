// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: contacts.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createContact = `-- name: CreateContact :one
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
) RETURNING id, address1, address2, address3, address4, postcode, country_id, phone, mobile, user_id, notes, created_at, updated_at
`

type CreateContactParams struct {
	UserID    uuid.UUID      `db:"user_id"`
	Address1  sql.NullString `db:"address1"`
	Address2  sql.NullString `db:"address2"`
	Address3  sql.NullString `db:"address3"`
	Address4  sql.NullString `db:"address4"`
	Postcode  sql.NullString `db:"postcode"`
	CountryID uuid.NullUUID  `db:"country_id"`
	Phone     sql.NullString `db:"phone"`
	Mobile    sql.NullString `db:"mobile"`
	Notes     sql.NullString `db:"notes"`
}

func (q *Queries) CreateContact(ctx context.Context, arg CreateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, createContact,
		arg.UserID,
		arg.Address1,
		arg.Address2,
		arg.Address3,
		arg.Address4,
		arg.Postcode,
		arg.CountryID,
		arg.Phone,
		arg.Mobile,
		arg.Notes,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Address1,
		&i.Address2,
		&i.Address3,
		&i.Address4,
		&i.Postcode,
		&i.CountryID,
		&i.Phone,
		&i.Mobile,
		&i.UserID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteContact = `-- name: DeleteContact :exec
DELETE FROM contacts
WHERE id = $1
`

func (q *Queries) DeleteContact(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteContact, id)
	return err
}

const getContact = `-- name: GetContact :one
SELECT id, address1, address2, address3, address4, postcode, country_id, phone, mobile, user_id, notes, created_at, updated_at FROM contacts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetContact(ctx context.Context, id uuid.UUID) (Contact, error) {
	row := q.db.QueryRowContext(ctx, getContact, id)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Address1,
		&i.Address2,
		&i.Address3,
		&i.Address4,
		&i.Postcode,
		&i.CountryID,
		&i.Phone,
		&i.Mobile,
		&i.UserID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getContactByUserID = `-- name: GetContactByUserID :one
SELECT id, address1, address2, address3, address4, postcode, country_id, phone, mobile, user_id, notes, created_at, updated_at FROM contacts WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetContactByUserID(ctx context.Context, userID uuid.UUID) (Contact, error) {
	row := q.db.QueryRowContext(ctx, getContactByUserID, userID)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Address1,
		&i.Address2,
		&i.Address3,
		&i.Address4,
		&i.Postcode,
		&i.CountryID,
		&i.Phone,
		&i.Mobile,
		&i.UserID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listContacts = `-- name: ListContacts :many
SELECT id, address1, address2, address3, address4, postcode, country_id, phone, mobile, user_id, notes, created_at, updated_at FROM contacts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListContactsParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) ListContacts(ctx context.Context, arg ListContactsParams) ([]Contact, error) {
	rows, err := q.db.QueryContext(ctx, listContacts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Contact{}
	for rows.Next() {
		var i Contact
		if err := rows.Scan(
			&i.ID,
			&i.Address1,
			&i.Address2,
			&i.Address3,
			&i.Address4,
			&i.Postcode,
			&i.CountryID,
			&i.Phone,
			&i.Mobile,
			&i.UserID,
			&i.Notes,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const truncateContacts = `-- name: TruncateContacts :exec
truncate contacts cascade
`

func (q *Queries) TruncateContacts(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, truncateContacts)
	return err
}

const updateContact = `-- name: UpdateContact :one
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
RETURNING id, address1, address2, address3, address4, postcode, country_id, phone, mobile, user_id, notes, created_at, updated_at
`

type UpdateContactParams struct {
	ID        uuid.UUID      `db:"id"`
	UserID    uuid.UUID      `db:"user_id"`
	Address1  sql.NullString `db:"address1"`
	Address2  sql.NullString `db:"address2"`
	Address3  sql.NullString `db:"address3"`
	Address4  sql.NullString `db:"address4"`
	Postcode  sql.NullString `db:"postcode"`
	CountryID uuid.NullUUID  `db:"country_id"`
	Phone     sql.NullString `db:"phone"`
	Mobile    sql.NullString `db:"mobile"`
	Notes     sql.NullString `db:"notes"`
}

func (q *Queries) UpdateContact(ctx context.Context, arg UpdateContactParams) (Contact, error) {
	row := q.db.QueryRowContext(ctx, updateContact,
		arg.ID,
		arg.UserID,
		arg.Address1,
		arg.Address2,
		arg.Address3,
		arg.Address4,
		arg.Postcode,
		arg.CountryID,
		arg.Phone,
		arg.Mobile,
		arg.Notes,
	)
	var i Contact
	err := row.Scan(
		&i.ID,
		&i.Address1,
		&i.Address2,
		&i.Address3,
		&i.Address4,
		&i.Postcode,
		&i.CountryID,
		&i.Phone,
		&i.Mobile,
		&i.UserID,
		&i.Notes,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}