package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/houtens/postbag/internal/models"
)

func stringToInt32(s string) int32 {
	// Parse float as some values may be float not int
	value, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Fatal("could not parse value: ", err)
	}
	return int32(value)
}

func importRatings() {
	filename := "seed/data/export/ratings.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open csv file %s\n", err)
	}
	defer f.Close()

	type RatingCSV struct {
		ID          string `csv:"id"`
		Player      string `csv:"player"`
		Tournament  string `csv:"tournament"`
		Division    string `csv:"division"`
		Team        string `csv:"team"`
		Games       string `csv:"games"`
		StartRating string `csv:"start_rating"`
		Points      string `csv:"points"`
		Date        string `csv:"date"`
		Locked      string `csv:"locked"`
		Verified    string `csv:"verified"`
	}

	rows := []RatingCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	fmt.Print("Importing ratings...")
	now := time.Now()
	// cnt := 0
	for _, r := range rows {
		// if cnt > 0 && cnt%10000 == 0 {
		// 	fmt.Printf(" %d", cnt)
		// }
		// Every rating must have a user and a tournament associated
		userID := userByXID(r.Player).UUID
		tournamentID := tournamentByXID(r.Tournament).UUID

		division := stringToInt32(r.Division)
		games := stringToInt32(r.Games)
		startRating := stringToInt32(r.StartRating)
		ratingPoints := stringToInt32(r.Points)
		isLocked := setLocked(r.Locked)

		arg := models.CreateRatingParams{
			UserID:       userID,
			TournamentID: tournamentID,
			Division:     division,
			NumGames:     sql.NullInt32{Int32: games, Valid: true},
			StartRating:  sql.NullInt32{Int32: startRating, Valid: true},
			RatingPoints: sql.NullInt32{Int32: ratingPoints, Valid: true},
			// TODO : We need to fill these in later once we have results
			// - OppRatingsSum
			// - NumWins
			IsLocked: isLocked,
			XID:      r.ID,
		}

		// Rating
		_, err := db.CreateRating(context.Background(), arg)
		if err != nil {
			fmt.Println(r)
			log.Fatal("failed to save rating: ", err)
		}
		// cnt++
	}

	// report number of ratings imported
	arg := models.ListRatingsParams{
		Limit:  1000000,
		Offset: 0,
	}
	ratings, err := db.ListRatings(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%d found, ", len(ratings))

	fmt.Println("... ", time.Since(now))
}

func truncateRatings() {
	db.TruncateRatings(context.Background())
}
