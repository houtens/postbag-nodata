package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTournament(t *testing.T) Tournament {
	tournamentState := createRandomTournamentState(t)

	arg := CreateTournamentParams{
		Name: randomString(24),
		ShortName: sql.NullString{
			String: randomName(),
			Valid:  true,
		},
		State: tournamentState.ID,
		StartDate: sql.NullTime{
			Time:  randomDate(),
			Valid: true,
		},
		EndDate: sql.NullTime{
			Time:  randomDate(),
			Valid: true,
		},
		NumDivisions: sql.NullInt32{
			Int32: int32(randomInt(1, 12)),
			Valid: true,
		},
		NumRounds: sql.NullInt32{
			Int32: int32(randomInt(1, 21)),
			Valid: true,
		},
		NumEntries: sql.NullInt32{
			Int32: int32(randomInt(14, 88)),
			Valid: true,
		},
		IsLocked: false,
	}

	tournament, err := db.CreateTournament(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tournament)

	return tournament
}

func TestCreateTournament(t *testing.T) {
	createRandomTournament(t)
}

func TestGetTournament(t *testing.T) {
	tournament1 := createRandomTournament(t)
	tournament2, err := db.GetTournament(context.Background(), tournament1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tournament2)
}

func TestListTournament(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTournament(t)
	}

	arg := ListTournamentsParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := db.ListTournaments(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateTournament(t *testing.T) {
	tournament1 := createRandomTournament(t)
	tournamentState := createRandomTournamentState(t)

	arg := UpdateTournamentParams{
		ID:           tournament1.ID,
		Name:         randomString(24),
		ShortName:    sql.NullString{String: randomString(12), Valid: true},
		StartDate:    sql.NullTime{Time: randomDate(), Valid: true},
		EndDate:      sql.NullTime{Time: randomDate(), Valid: true},
		State:        tournamentState.ID,
		NumDivisions: sql.NullInt32{Int32: int32(randomInt(1, 22))},
		NumRounds:    tournament1.NumRounds,
		NumEntries:   tournament1.NumEntries,
		IsLocked:     tournament1.IsLocked,
		CreatorID:    tournament1.CreatorID,
		OrganiserID:  tournament1.OrganiserID,
		DirectorID:   tournament1.DirectorID,
		CoperatorID:  tournament1.CoperatorID,
		XID:          tournament1.XID,
	}

	tournament2, err := db.UpdateTournament(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tournament2)

	require.Equal(t, tournament1.ID, tournament2.ID)
	require.NotEqual(t, tournament1.Name, tournament2.Name)
	require.NotEqual(t, tournament1.ShortName, tournament2.ShortName)
	require.NotEqual(t, tournament1.StartDate, tournament2.StartDate)
	require.NotEqual(t, tournament1.State, tournament2.State)
	require.NotEqual(t, tournament1.EndDate, tournament2.EndDate)
	require.NotEqual(t, tournament1.NumDivisions, tournament2.NumDivisions)
	require.Equal(t, tournament1.CreatorID, tournament2.CreatorID)
}

func TestDeleteTournament(t *testing.T) {
	tournament1 := createRandomTournament(t)
	db.DeleteTournament(context.Background(), tournament1.ID)

	tournament2, err := db.GetUser(context.Background(), tournament1.ID)
	require.Error(t, err)
	require.Empty(t, tournament2)

}
