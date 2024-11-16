package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

type Rating struct {
	ID     uuid.UUID
	Rank   int
	Name   string
	Club   string
	Rating int
}

type RatingsList []Rating

// RatingsList renders the latest ratings
func (s Service) GetRatingsList() (RatingsList, error) {
	q := models.New(s.DB)

	// recent tournaments are the past 20 or over X years
	arg := models.ListRatingsRankParams{
		Limit:  10000,
		Offset: 0,
	}

	ratings, err := q.ListRatingsRank(context.Background(), arg)
	if err != nil {
		fmt.Println(err)
	}

	ratingsList := RatingsList{}

	for _, r := range ratings {
		ratingsList = append(
			ratingsList,
			Rating{
				ID:     r.UserID,
				Rank:   int(r.RowNumber),
				Name:   fullName(r.FirstName, r.LastName, r.Title),
				Club:   r.Club.String,
				Rating: int(r.EndRating.Int32),
			},
		)
	}

	return ratingsList, nil
}

// FilterRatings returns a slice of ratings filtered by the query string
func (s Service) FilterRatings(query string) (RatingsList, error) {
	ratings, err := s.GetRatingsList()
	if err != nil {
		return RatingsList{}, err
	}

	filteredRatings := RatingsList{}

	for _, r := range ratings {
		found := false
		if strings.Contains(strings.ToLower(r.Name), query) {
			found = true
		}
		if strings.Contains(strings.ToLower(r.Club), query) {
			found = true
		}

		if found {
			filteredRatings = append(filteredRatings, r)
		}
	}

	return filteredRatings, nil
}
