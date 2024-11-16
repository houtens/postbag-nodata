package models

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createRandomResult(t *testing.T) Result {

	// create players and tournament
	player1 := createRandomUser(t)
	player2 := createRandomUser(t)
	tournament := createRandomTournament(t)

	arg := CreateResultParams{
		Player1ID:    uuid.NullUUID{UUID: player1.ID, Valid: true},
		Player2ID:    uuid.NullUUID{UUID: player2.ID, Valid: true},
		Score1:       int32(randomInt(300, 500)),
		Score2:       int32(randomInt(300, 500)),
		TournamentID: tournament.ID,
		Type:         0,
		RoundNum:     int32(randomInt(1, 16)),
		IsLocked:     false,
	}

	result, err := db.CreateResult(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	return result
}

func TestCreateResult(t *testing.T) {
	createRandomResult(t)
}

func TestGetResult(t *testing.T) {
	result1 := createRandomResult(t)
	result2, err := db.GetResult(context.Background(), result1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, result2)
}

func TestListResults(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomResult(t)
	}

	arg := ListResultsParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := db.ListResults(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateResult(t *testing.T) {
	result1 := createRandomResult(t)
	player1 := createRandomUser(t)
	player2 := createRandomUser(t)
	tournament1 := createRandomTournament(t)

	arg := UpdateResultParams{
		ID:           result1.ID,
		Player1ID:    uuid.NullUUID{UUID: player1.ID, Valid: true},
		Player2ID:    uuid.NullUUID{UUID: player2.ID, Valid: true},
		Score1:       int32(randomInt(300, 500)),
		Score2:       int32(randomInt(300, 500)),
		TournamentID: tournament1.ID,
		Type:         123,
		RoundNum:     int32(randomInt(1, 22)),
		IsLocked:     true,
	}

	result2, err := db.UpdateResult(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, result2)

	require.Equal(t, result1.ID, result2.ID)
	require.NotEqual(t, result1.Player1ID, result2.Player1ID)
	require.NotEqual(t, result1.Player2ID, result2.Player2ID)
}

func TestDeleteResult(t *testing.T) {
	result1 := createRandomResult(t)
	db.DeleteResult(context.Background(), result1.ID)

	result2, err := db.GetUser(context.Background(), result1.ID)
	require.Error(t, err)
	require.Empty(t, result2)
}
