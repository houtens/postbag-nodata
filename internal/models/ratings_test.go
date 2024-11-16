package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomRating(t *testing.T) Rating {
	user := createRandomUser(t)
	tournament := createRandomTournament(t)

	arg := CreateRatingParams{
		UserID:       user.ID,
		TournamentID: tournament.ID,
	}

	rating, err := db.CreateRating(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rating)

	return rating
}

func TestCreateRating(t *testing.T) {
	createRandomRating(t)
}

func TestGetRating(t *testing.T) {
	rating1 := createRandomRating(t)
	rating2, err := db.GetRating(context.Background(), rating1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, rating2)
}

func TestListRatings(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomRating(t)
	}

	arg := ListRatingsParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := db.ListRatings(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateRating(t *testing.T) {
	rating1 := createRandomRating(t)
	tournament := createRandomTournament(t)

	arg := UpdateRatingParams{
		ID:            rating1.ID,
		UserID:        rating1.UserID,
		TournamentID:  tournament.ID,
		Division:      rating1.Division,
		NumGames:      rating1.NumGames,
		StartRating:   rating1.StartRating,
		RatingPoints:  rating1.RatingPoints,
		OppRatingsSum: rating1.OppRatingsSum,
		NumWins:       rating1.NumWins,
		IsLocked:      rating1.IsLocked,
		XID:           rating1.XID,
	}

	rating2, err := db.UpdateRating(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, rating2)
	require.NotEqual(t, rating1.TournamentID, rating2.TournamentID)
}

func TestDeleteRating(t *testing.T) {
	rating1 := createRandomRating(t)
	db.DeleteRating(context.Background(), rating1.ID)

	rating2, err := db.GetUser(context.Background(), rating1.ID)
	require.Error(t, err)
	require.Empty(t, rating2)
}
