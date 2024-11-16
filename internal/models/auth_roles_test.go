package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func CreateRandomAuthRole(t *testing.T) AuthRole {
	arg := CreateAuthRoleParams{
		Name:               "Random Privilege",
		CanLogin:           randomBool(),
		IsGuest:            randomBool(),
		IsMembersAdmin:     randomBool(),
		IsClubsAdmin:       randomBool(),
		IsRatingsAdmin:     randomBool(),
		IsTournamentsAdmin: randomBool(),
		IsSuperAdmin:       randomBool(),
	}

	authrole, err := db.CreateAuthRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, authrole)

	require.Equal(t, arg.Name, authrole.Name)
	require.Equal(t, arg.CanLogin, authrole.CanLogin)
	require.Equal(t, arg.IsGuest, authrole.IsGuest)
	require.Equal(t, arg.IsMembersAdmin, authrole.IsMembersAdmin)
	require.Equal(t, arg.IsClubsAdmin, authrole.IsClubsAdmin)
	require.Equal(t, arg.IsRatingsAdmin, authrole.IsRatingsAdmin)
	require.Equal(t, arg.IsTournamentsAdmin, authrole.IsTournamentsAdmin)
	require.Equal(t, arg.IsSuperAdmin, authrole.IsSuperAdmin)

	require.NotZero(t, authrole.ID)

	return authrole
}

func TestCreateAuthRole(t *testing.T) {
	CreateRandomAuthRole(t)
}

func TestGetAuthRole(t *testing.T) {
	authrole1 := CreateRandomAuthRole(t)

	authrole2, err := db.GetAuthRole(context.Background(), authrole1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, authrole2)

	require.Equal(t, authrole1.ID, authrole2.ID)
	require.Equal(t, authrole1.CanLogin, authrole2.CanLogin)
	require.Equal(t, authrole1.IsGuest, authrole2.IsGuest)
	require.Equal(t, authrole1.IsMembersAdmin, authrole2.IsMembersAdmin)
	require.Equal(t, authrole1.IsClubsAdmin, authrole2.IsClubsAdmin)
	require.Equal(t, authrole1.IsRatingsAdmin, authrole2.IsRatingsAdmin)
	require.Equal(t, authrole1.IsTournamentsAdmin, authrole2.IsTournamentsAdmin)
	require.Equal(t, authrole1.IsSuperAdmin, authrole2.IsSuperAdmin)

	require.NotZero(t, authrole2.ID)
}

func TestListAuthRole(t *testing.T) {
	for i := 0; i < 10; i++ {
		CreateRandomAuthRole(t)
	}

	arg := ListAuthRolesParams{
		Limit:  5,
		Offset: 5,
	}

	authroles, err := db.ListAuthRoles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, authroles, 5)

	for _, authrole := range authroles {
		require.NotEmpty(t, authrole)
	}
}

func TestUpdateAuthRole(t *testing.T) {
	authrole1 := CreateRandomAuthRole(t)

	arg := UpdateAuthRoleParams{
		ID:                 authrole1.ID,
		Name:               authrole1.Name,
		CanLogin:           !randomBool(),
		IsGuest:            !randomBool(),
		IsMembersAdmin:     !randomBool(),
		IsClubsAdmin:       !randomBool(),
		IsRatingsAdmin:     !randomBool(),
		IsTournamentsAdmin: !randomBool(),
		IsSuperAdmin:       !randomBool(),
	}

	authrole2, err := db.UpdateAuthRole(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, authrole2)

	require.Equal(t, authrole1.ID, authrole2.ID)
	require.Equal(t, authrole1.Name, authrole2.Name)
	require.Equal(t, arg.CanLogin, authrole2.CanLogin)
	require.Equal(t, arg.IsGuest, authrole2.IsGuest)
	require.Equal(t, arg.IsMembersAdmin, authrole2.IsMembersAdmin)
	require.Equal(t, arg.IsClubsAdmin, authrole2.IsClubsAdmin)
	require.Equal(t, arg.IsRatingsAdmin, authrole2.IsRatingsAdmin)
	require.Equal(t, arg.IsTournamentsAdmin, authrole2.IsTournamentsAdmin)
	require.Equal(t, arg.IsSuperAdmin, authrole2.IsSuperAdmin)
}
