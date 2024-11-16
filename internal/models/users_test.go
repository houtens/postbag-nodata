package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomUser(t *testing.T) User {
	auth_role := CreateRandomAuthRole(t)
	club := createRandomClub(t)
	title := createRandomTitle(t)

	arg := CreateUserParams{
		FirstName:    randomName(),
		LastName:     randomName(),
		AltName:      sql.NullString{String: randomName(), Valid: true},
		Email:        sql.NullString{String: randomString(32), Valid: true},
		PasswordHash: sql.NullString{String: randomString(32), Valid: true},
		AbspNum:      sql.NullInt32{Int32: int32(randomInt(0, 12000)), Valid: true},
		ClubID:       uuid.NullUUID{UUID: club.ID, Valid: true},
		TitleID:      uuid.NullUUID{UUID: title.ID, Valid: true},
		RoleID:       auth_role.ID,
		IsDeceased:   randomBool(),
		XLife:        randomBool(),
	}

	user, err := db.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := db.GetUser(context.Background(), user1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, user2)
}

func TestListUsers(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomUser(t)
	}

	arg := ListUsersParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := db.ListUsers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateUser(t *testing.T) {
	user1 := createRandomUser(t)

	arg := UpdateUserParams{
		ID:           user1.ID,
		FirstName:    randomName(),
		LastName:     randomName(),
		AltName:      user1.AltName,
		ClubID:       user1.ClubID,
		TitleID:      user1.TitleID,
		Email:        sql.NullString{String: randomString(32), Valid: true},
		PasswordHash: sql.NullString{String: randomString(32), Valid: true},
		AbspNum:      sql.NullInt32{Int32: int32(randomInt(15000, 20000)), Valid: true},
		RoleID:       user1.RoleID,
	}

	user2, err := db.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.ID, user2.ID)
	require.NotEqual(t, user1.FirstName, user2.FirstName)
	require.NotEqual(t, user1.LastName, user2.LastName)
	require.NotEqual(t, user1.Email, user2.Email)
	require.NotEqual(t, user1.PasswordHash, user2.PasswordHash)
	require.NotEqual(t, user1.AbspNum, user2.AbspNum)
}

func TestDeleteUser(t *testing.T) {
	user1 := createRandomUser(t)
	err := db.DeleteUser(context.Background(), user1.ID)
	require.NoError(t, err)

	user2, err := db.GetUser(context.Background(), user1.ID)
	require.Error(t, err)
	require.Empty(t, user2)

}
