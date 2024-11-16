package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

// return true when the deviation is large
func checkDeviation(c, p int32) bool {
	// if previous is zero, whatever current is will do
	if p == 0 {
		return true
	}
	// if current is zero, we do not want to recommend this rating
	if c == 0 {
		return false
	}

	// small difference between c and p should all current to be recommended
	if math.Abs(float64(c-p)) < 20 {
		return true
	}

	// failure of Abs condition
	return false
}

func getPrevStartRating(pID uuid.UUID) int32 {
	p, err := db.GetRating(context.Background(), pID)
	if err != nil {
		log.Fatal(err)
	}

	return p.StartRating.Int32
}

func getPrevEndRating(pID uuid.UUID) int32 {
	p, err := db.GetRating(context.Background(), pID)
	if err != nil {
		log.Fatal(err)
	}

	return p.EndRating.Int32
}

func updateRating(pID uuid.UUID, startRating, endRating int32) {

	var err error

	// Only trust the ID of the previous rating
	p, err := db.GetRating(context.Background(), pID)
	if err != nil {
		log.Fatal(err)
	}
	if endRating > 0 && endRating < 50 {
		endRating = 50
	}

	// Update the previous end rating with the current start rating
	arg1 := models.UpdateRatingParams{
		ID:            p.ID,
		UserID:        p.UserID,
		TournamentID:  p.TournamentID,
		Division:      p.Division,
		NumGames:      p.NumGames,
		StartRating:   sql.NullInt32{Int32: startRating, Valid: true},
		EndRating:     sql.NullInt32{Int32: endRating, Valid: true},
		RatingPoints:  p.RatingPoints,
		OppRatingsSum: p.OppRatingsSum,
		NumWins:       p.NumWins,
		IsLocked:      p.IsLocked,
		XID:           p.XID,
	}

	_, err = db.UpdateRating(context.Background(), arg1)
	if err != nil {
		log.Fatal(err)
	}
}

func calcAverageRating(points, games int32) int32 {
	var result int32
	if games > 0 && points > 0 {
		result = int32(float64(points) / float64(games))
	}
	return result
}

// startRating sr, ratingPoints rp, numGames ng
func calcRating(start, points, games int32) int32 {
	// can this be caught outside the function call
	if start == 0 || points == 0 || games == 0 {
		return 0
	}

	// averagePoints is just TR
	averagePoints := points / games
	// ratings difference
	ratingDifference := averagePoints - start
	if start == 0 {
		ratingDifference = averagePoints
	}
	// new predicted rating
	predicted := float64(ratingDifference*games) / 120
	// round the ratings difference
	return int32(math.Round(predicted*1.0 + float64(start)))
}

func withinBounds(a, b int32) bool {
	// If either is zero then deny it
	if a == 0 || b == 0 {
		return false
	}

	// var ratio float64
	var ratio = math.Abs(1.0 - float64(b)/float64(a))
	return ratio < 0.2
}

func endRatings() {
	fmt.Printf("Fixing end ratings... ")
	now := time.Now()

	arg := models.ListUsersParams{
		Limit:  10000,
		Offset: 0,
	}

	users, err := db.ListUsers(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	var previous models.ListRatingsByUserRow
	var played, points int32
	var startRating int32
	var newRating int32
	var isProvisional = true

	for _, u := range users {

		ratings, err := db.ListRatingsByUser(context.Background(), u.ID)
		if err != nil {
			log.Fatal(err)
		}

		// Initialise loop values for each user
		played, points = 0, 0
		startRating = 0
		isProvisional = true
		newRating = 0

		previous = models.ListRatingsByUserRow{}

		for i, r := range ratings {
			// fmt.Println("--")
			// fmt.Println("current:", r.ID, r.StartRating.Int32)
			// fmt.Println("previous:", previous.ID, previous.StartRating.Int32)
			// fmt.Println("--")

			// Established players start with a non-zero rating and so the 50 game rules does not apply
			if i == 0 && r.StartRating.Int32 > 0 {
				isProvisional = false
			}
			// Players become established after 50 games
			if isProvisional && played >= 50 {
				isProvisional = false
			}

			startRating = r.StartRating.Int32

			// Determine which rating to use as the tournament start rating
			// On first iteration we simply use the startRating
			if i > 0 {
				// 1. If startRating is zero and previous startRating is not
				if startRating == 0 && previous.EndRating.Int32 > 0 {
					startRating = previous.EndRating.Int32
				}

				// 2. startRating is more than 20% away from previous startRating
				if !checkDeviation(startRating, previous.StartRating.Int32) && !isProvisional {
					// start and previous are wildly different, do not trust start rating
					// fmt.Println("20% rule - using the previous rating")
					startRating = previous.EndRating.Int32
				}

				// 3. Conditionally update the previous end rating; set within the loop and update database
				if startRating == r.StartRating.Int32 && startRating != previous.EndRating.Int32 {
					// fmt.Println("updating previous end rating with", startRating)
					previous.EndRating.Int32 = startRating
					// fmt.Println("=", previous.EndRating)
					// fmt.Println(previous.EndRating.Int32, r.EndRating.Int32)
					updateRating(previous.ID, previous.StartRating.Int32, previous.EndRating.Int32)
				}
			}

			// Update number of games played
			played = played + r.NumGames.Int32
			points = points + r.RatingPoints.Int32

			// Estimate the current rating -- handle ratings, points and games = 0
			if isProvisional {
				newRating = calcAverageRating(points, played)
			} else {
				newRating = calcRating(startRating, r.RatingPoints.Int32, r.NumGames.Int32)
			}

			// Update the current rating
			r.EndRating.Int32 = newRating
			updateRating(r.ID, startRating, newRating)
			// Show the player/rating details
			// fmt.Println("*", u.FirstName, u.LastName, u.ID, r.StartRating.Int32, r.EndRating.Int32, r.NumGames.Int32, r.RatingPoints.Int32, "[", startRating, "]", newRating)

			// Save current rating as previous for next loop
			previous = r
			// fmt.Println("update:", previous.ID)
			previous.StartRating.Int32 = startRating
		}
	}
	fmt.Println(time.Since(now))

}
