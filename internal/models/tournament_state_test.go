package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTournamentState(t *testing.T) TournamentState {
	arg := CreateTournamentStateParams{
		Name: randomString(15),
		Code: randomString(12),
	}
	tournamentState, err := db.CreateTournamentState(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tournamentState)

	return tournamentState
}

func TestCreateTournamentState(t *testing.T) {
	createRandomTournamentState(t)
}

func TestGetTournamentState(t *testing.T) {
	tournamentState1 := createRandomTournamentState(t)
	tournamentState2, err := db.GetTournamentState(context.Background(), tournamentState1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, tournamentState2)
}

func TestListTournamentStates(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTournamentState(t)
	}

	arg := ListTournamentStatesParams{
		Limit:  5,
		Offset: 5,
	}

	tournamentStates, err := db.ListTournamentStates(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tournamentStates)
	require.Len(t, tournamentStates, 5)
}

func TestUpdateTournamentState(t *testing.T) {
	tournamentState1 := createRandomTournamentState(t)

	arg := UpdateTournamentStateParams{
		ID:   tournamentState1.ID,
		Name: randomString(18),
		Code: randomString(15),
	}

	tournamentState2, err := db.UpdateTournamentState(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, tournamentState2)
	require.NotEqual(t, tournamentState1.Name, tournamentState2.Name)
	require.NotEqual(t, tournamentState1.Code, tournamentState2.Code)
}

func TestDeleteTournamentState(t *testing.T) {
	tournamentState1 := createRandomTournamentState(t)
	db.DeleteTournamentState(context.Background(), tournamentState1.ID)

	tournamentState2, err := db.GetTournamentState(context.Background(), tournamentState1.ID)
	require.Error(t, err)
	require.Empty(t, tournamentState2)
}
