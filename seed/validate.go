package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

func calculateRatingWins() {

	// for each ratings entry - user and tournament
	// find user/tournament in results
	// count wins for player - union all on player1 score1 vs player2 score2 where match

	fmt.Printf("Updating ratings with number of wins...")
	now := time.Now()

	// iterate over each rating entry
	arg := models.ListRatingsParams{
		Limit:  200000,
		Offset: 0,
	}
	ratings, err := db.ListRatings(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	for _, r := range ratings {
		// fmt.Println(r)
		// count wins for user in tournament - may be decimal
		arg := models.GetWinsByUserTournamentParams{
			TournamentID: r.TournamentID,
			Player1ID: uuid.NullUUID{
				UUID:  r.UserID,
				Valid: true,
			},
		}
		wins, err := db.GetWinsByUserTournament(context.Background(), arg)
		if err != nil {
			// log.Fatal(err)
			continue
		}
		// fmt.Println(r.UserID, r.TournamentID, wins)

		// Update ratings table num_wins
		arg0 := models.UpdateRatingsNumWinsParams{
			UserID:       r.UserID,
			TournamentID: r.TournamentID,
			NumWins:      float32(wins),
		}
		_, err = db.UpdateRatingsNumWins(context.Background(), arg0)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(time.Since(now))
}

func validateTournamentResults() {
	fmt.Printf("Updating tournament results with metadata...")
	now := time.Now()

	// validate data already stored
	arg := models.ListTournamentsParams{
		Limit:  100000,
		Offset: 0,
	}
	tournaments, err := db.ListTournaments(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	for _, t := range tournaments {
		// Only update locked tournaments
		if !t.IsLocked {
			continue
		}

		err := updateRoundsAndDivisions(t)
		if err != nil {
			log.Fatal(err)
		}

		err = updatePlayerCounts(t)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println(time.Since(now))
}

func updateRoundsAndDivisions(t models.Tournament) error {
	rt, err := db.GetRatingsTournamentMetadata(context.Background(), t.ID)
	if err != nil {
		return fmt.Errorf("unable to fetch ratings for tournament: %w", err)
	}

	// Handle no rows returned from rt
	if rt.Rounds == nil || rt.Divisions == nil {
		return nil
	}

	// Convert types for values of rounds and divisions found in the tournaments table
	// and those calculated from the ratings.
	foundRounds := int(t.NumRounds.Int32)
	calcRounds := int(rt.Rounds.(int64))
	foundDivisions := int(t.NumDivisions.Int32)
	calcDivisions := int(rt.Divisions.(int64) + 1)

	// Only update the tournament values if larger ones are calculated
	if foundRounds < calcRounds || foundDivisions < calcDivisions {
		// NOTE: debug
		// fmt.Printf("%s %s - rounds:%d = %d, division:%d = %d\n", t.ID, t.Name, foundRounds, calcRounds, foundDivisions, calcDivisions)

		// update this tournament and set rounds/divisions accordingly
		rounds := int32(rt.Rounds.(int64))
		divisions := int32(rt.Divisions.(int64) + 1)
		arg := models.UpdateTournamentRoundsDivisionsParams{
			ID: t.ID,
			NumRounds: sql.NullInt32{
				Int32: rounds,
				Valid: true,
			},
			NumDivisions: sql.NullInt32{
				Int32: divisions,
				Valid: true,
			},
		}
		_, err := db.UpdateTournamentRoundsDivisions(context.Background(), arg)
		if err != nil {
			return fmt.Errorf("unable to update tournament rounds and divisions: %w", err)
		}
	}
	return nil
}

func updatePlayerCounts(t models.Tournament) error {

	if !t.NumEntries.Valid {
		// Count players in each tournament and update tournaments.entrants
		numPlayers, err := db.GetCountPlayersInTournament(context.Background(), t.ID)
		if err != nil {
			return fmt.Errorf("unable to get count of players in tournament: %w", err)
		}

		// NOTE: debug
		// fmt.Printf("%s %s %d %T\n", t.ID, t.Name, numPlayers, numPlayers)
		np := int32(numPlayers)
		arg := models.UpdateTournamentEntriesParams{
			ID: t.ID,
			NumEntries: sql.NullInt32{
				Int32: np,
				Valid: true,
			},
		}
		_, err = db.UpdateTournamentEntries(context.Background(), arg)
		if err != nil {
			return fmt.Errorf("unable to update entrants count in tournament: %w", err)
		}
	}

	return nil
}
