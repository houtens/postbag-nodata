// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: users.sql

package models

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
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
) RETURNING id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry
`

type CreateUserParams struct {
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	AltName      sql.NullString `db:"alt_name"`
	Email        sql.NullString `db:"email"`
	PasswordHash sql.NullString `db:"password_hash"`
	AbspNum      sql.NullInt32  `db:"absp_num"`
	ClubID       uuid.NullUUID  `db:"club_id"`
	TitleID      uuid.NullUUID  `db:"title_id"`
	RoleID       uuid.UUID      `db:"role_id"`
	IsDeceased   bool           `db:"is_deceased"`
	XLife        bool           `db:"x_life"`
	XPost        bool           `db:"x_post"`
	XID          string         `db:"x_id"`
	Avatar       sql.NullString `db:"avatar"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.FirstName,
		arg.LastName,
		arg.AltName,
		arg.Email,
		arg.PasswordHash,
		arg.AbspNum,
		arg.ClubID,
		arg.TitleID,
		arg.RoleID,
		arg.IsDeceased,
		arg.XLife,
		arg.XPost,
		arg.XID,
		arg.Avatar,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE FROM users
where id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUser, id)
	return err
}

const getActiveUsers = `-- name: GetActiveUsers :many

SELECT id, first_name, last_name, alt_name FROM users WHERE is_deceased IS false order by last_name
`

type GetActiveUsersRow struct {
	ID        uuid.UUID      `db:"id"`
	FirstName string         `db:"first_name"`
	LastName  string         `db:"last_name"`
	AltName   sql.NullString `db:"alt_name"`
}

// NOTE: custom query for users used in player dropdowns
func (q *Queries) GetActiveUsers(ctx context.Context) ([]GetActiveUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getActiveUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetActiveUsersRow{}
	for rows.Next() {
		var i GetActiveUsersRow
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AltName,
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

const getUser = `-- name: GetUser :one
SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one

SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users WHERE email = $1 LIMIT 1
`

// NOTE: custom queries
func (q *Queries) GetUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const getUserByValidToken = `-- name: GetUserByValidToken :one
SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users WHERE pw_token = $1 AND pw_token_expiry > NOW() LIMIT 1
`

func (q *Queries) GetUserByValidToken(ctx context.Context, pwToken sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByValidToken, pwToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const getUserByXID = `-- name: GetUserByXID :one
SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users
WHERE x_id = $1 LIMIT 1
`

// NOTE: seed
func (q *Queries) GetUserByXID(ctx context.Context, xID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByXID, xID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const listUsers = `-- name: ListUsers :many
SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users
ORDER BY last_name, first_name
LIMIT $1
OFFSET $2
`

type ListUsersParams struct {
	Limit  int32 `db:"limit"`
	Offset int32 `db:"offset"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AltName,
			&i.Email,
			&i.PasswordHash,
			&i.AbspNum,
			&i.ClubID,
			&i.TitleID,
			&i.RoleID,
			&i.XLife,
			&i.XPost,
			&i.XID,
			&i.IsDeceased,
			&i.Avatar,
			&i.PwToken,
			&i.PwTokenExpiry,
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

const listUsersByXLife = `-- name: ListUsersByXLife :many
SELECT id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry FROM users WHERE x_life = 't' ORDER BY id
`

// NOTE: seed
func (q *Queries) ListUsersByXLife(ctx context.Context) ([]User, error) {
	rows, err := q.db.QueryContext(ctx, listUsersByXLife)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.AltName,
			&i.Email,
			&i.PasswordHash,
			&i.AbspNum,
			&i.ClubID,
			&i.TitleID,
			&i.RoleID,
			&i.XLife,
			&i.XPost,
			&i.XID,
			&i.IsDeceased,
			&i.Avatar,
			&i.PwToken,
			&i.PwTokenExpiry,
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

const truncateUsers = `-- name: TruncateUsers :exec
truncate users cascade
`

func (q *Queries) TruncateUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, truncateUsers)
	return err
}

const updateUser = `-- name: UpdateUser :one
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
RETURNING id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry
`

type UpdateUserParams struct {
	ID           uuid.UUID      `db:"id"`
	FirstName    string         `db:"first_name"`
	LastName     string         `db:"last_name"`
	AltName      sql.NullString `db:"alt_name"`
	Email        sql.NullString `db:"email"`
	PasswordHash sql.NullString `db:"password_hash"`
	AbspNum      sql.NullInt32  `db:"absp_num"`
	ClubID       uuid.NullUUID  `db:"club_id"`
	TitleID      uuid.NullUUID  `db:"title_id"`
	RoleID       uuid.UUID      `db:"role_id"`
	IsDeceased   bool           `db:"is_deceased"`
	XLife        bool           `db:"x_life"`
	XPost        bool           `db:"x_post"`
	XID          string         `db:"x_id"`
	Avatar       sql.NullString `db:"avatar"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.FirstName,
		arg.LastName,
		arg.AltName,
		arg.Email,
		arg.PasswordHash,
		arg.AbspNum,
		arg.ClubID,
		arg.TitleID,
		arg.RoleID,
		arg.IsDeceased,
		arg.XLife,
		arg.XPost,
		arg.XID,
		arg.Avatar,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const updateUserExpireToken = `-- name: UpdateUserExpireToken :one
UPDATE users SET pw_token_expiry = NOW() WHERE pw_token = $1 RETURNING id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry
`

func (q *Queries) UpdateUserExpireToken(ctx context.Context, pwToken sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserExpireToken, pwToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const updateUserPasswordHash = `-- name: UpdateUserPasswordHash :one
UPDATE users SET password_hash = $3 WHERE id = $1 and pw_token = $2 RETURNING id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry
`

type UpdateUserPasswordHashParams struct {
	ID           uuid.UUID      `db:"id"`
	PwToken      sql.NullString `db:"pw_token"`
	PasswordHash sql.NullString `db:"password_hash"`
}

func (q *Queries) UpdateUserPasswordHash(ctx context.Context, arg UpdateUserPasswordHashParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserPasswordHash, arg.ID, arg.PwToken, arg.PasswordHash)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}

const updateUserSetToken = `-- name: UpdateUserSetToken :one
UPDATE users SET pw_token = $2, pw_token_expiry = NOW() + '15 minutes' WHERE id = $1 RETURNING id, first_name, last_name, alt_name, email, password_hash, absp_num, club_id, title_id, role_id, x_life, x_post, x_id, is_deceased, avatar, pw_token, pw_token_expiry
`

type UpdateUserSetTokenParams struct {
	ID      uuid.UUID      `db:"id"`
	PwToken sql.NullString `db:"pw_token"`
}

func (q *Queries) UpdateUserSetToken(ctx context.Context, arg UpdateUserSetTokenParams) (User, error) {
	row := q.db.QueryRowContext(ctx, updateUserSetToken, arg.ID, arg.PwToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.AltName,
		&i.Email,
		&i.PasswordHash,
		&i.AbspNum,
		&i.ClubID,
		&i.TitleID,
		&i.RoleID,
		&i.XLife,
		&i.XPost,
		&i.XID,
		&i.IsDeceased,
		&i.Avatar,
		&i.PwToken,
		&i.PwTokenExpiry,
	)
	return i, err
}