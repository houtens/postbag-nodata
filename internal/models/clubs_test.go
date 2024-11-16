package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomClub(t *testing.T) Club {
	arg := CreateClubParams{
		Name:     randomName(),
		County:   sql.NullString{String: randomName(), Valid: true},
		Website:  sql.NullString{},
		IsActive: true,
		Phone:    sql.NullString{},
		XID:      randomID(),
	}

	club, err := db.CreateClub(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, club)

	require.Equal(t, arg.Name, club.Name)
	require.Equal(t, arg.County, club.County)
	require.Equal(t, arg.Website, club.Website)
	require.Equal(t, arg.IsActive, club.IsActive)
	require.Equal(t, arg.Phone, club.Phone)

	require.NotZero(t, club.ID)

	return club
}

func TestCreateClub(t *testing.T) {
	createRandomClub(t)
}

func TestGetClub(t *testing.T) {
	club1 := createRandomClub(t)

	club2, err := db.GetClub(context.Background(), club1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, club2)

	// deepequal
	require.Equal(t, club1.ID, club2.ID)
	require.Equal(t, club1.Name, club2.Name)
	require.Equal(t, club1.County, club2.County)
	require.Equal(t, club1.Website, club2.Website)
	require.Equal(t, club1.IsActive, club2.IsActive)
	require.Equal(t, club1.Phone, club2.Phone)
	require.Equal(t, club1.Email, club2.Email)
	require.Equal(t, club1.ContactName, club2.ContactName)
	require.Equal(t, club1.CountryID, club2.CountryID)
	require.Equal(t, club1.XID, club2.XID)
	require.Equal(t, club1.CreatedAt, club2.CreatedAt)
	require.Equal(t, club1.UpdatedAt, club2.UpdatedAt)
}

func TestListClubs(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomClub(t)
	}

	arg := ListClubsParams{
		Limit:  5,
		Offset: 5,
	}

	clubs, err := db.ListClubs(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, clubs, 5)

	for _, club := range clubs {
		require.NotEmpty(t, club)
	}
}

func TestUpdateClub(t *testing.T) {
	club1 := createRandomClub(t)

	arg := UpdateClubParams{
		ID:          club1.ID,
		Name:        club1.Name,
		County:      club1.County,
		Website:     club1.Website,
		IsActive:    club1.IsActive,
		Phone:       club1.Phone,
		Email:       club1.Email,
		ContactName: club1.ContactName,
		CountryID:   club1.CountryID,
		XID:         club1.XID,
	}

	club2, err := db.UpdateClub(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, club2)

	require.NotEqual(t, club1.UpdatedAt, club2.UpdatedAt)
}

func TestDeleteClub(t *testing.T) {
	club1 := createRandomClub(t)
	db.DeleteClub(context.Background(), club1.ID)

	club2, err := db.GetClub(context.Background(), club1.ID)

	require.Error(t, err)
	require.Empty(t, club2)
}
