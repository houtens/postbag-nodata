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
	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

func importTournaments() {
	filename := "seed/data/export/tournaments.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open csv file %s\n", err)
	}
	defer f.Close()

	type TournamentCSV struct {
		ID        string `csv:"id"`
		ShortName string `csv:"short_name"`
		Title     string `csv:"title"`
		Date      string `csv:"date"`
		Entries   string `csv:"entries"`
		Division  string `csv:"division"`
		Rounds    string `csv:"rounds"`
		Results   string `csv:"results"`
		TeamEvent string `csv:"team_event"`
		Organiser string `csv:"organiser"`
		TD        string `csv:"td"`
		CO        string `csv:"co"`
		Creator   string `csv:"creator"`
		Updated   string `csv:"updated"`
		Comment   string `csv:"comment"`
		Invoice   string `csv:"invoice"`
		Locked    string `csv:"locked"`
		Datafile  string `csv:"datafile"`
	}

	rows := []TournamentCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	// Read the states into a map of name to uuid
	statesMap := readTournamentStates()

	fmt.Print("Import tournaments... ")
	now := time.Now()
	// cnt := 0
	for _, r := range rows {
		// if cnt > 0 && cnt%1000 == 0 {
		// 	fmt.Printf(" %d", cnt)
		// }

		// Set the state based on value of Locked. 0 = approved, 1 = rated
		var state uuid.UUID
		if r.Locked == "0" {
			state = statesMap["approved"]
		} else {
			state = statesMap["rated"]
		}

		// Parse the date into start and end times for the day/days
		startDate, _ := time.Parse("2006-01-02", r.Date)
		endDate, _ := time.Parse("2006-01-02", r.Date)
		// Find end of day by adding 1-day (to the date) and subtracting 1-second
		endDate = endDate.AddDate(0, 0, 1).Add(-1 * time.Second)

		numDivisions, _ := strconv.Atoi(r.Division)
		numRounds, _ := strconv.Atoi(r.Rounds)

		// Lookup users to assign
		creatorID := userByXID(r.Creator)
		organiserID := userByXID(r.Organiser)
		directorID := userByXID(r.TD)
		coperatorID := userByXID(r.CO)

		isLocked := setLocked(r.Locked)

		// Title should not be empty in future - use ShortName when Title is absent
		title := r.Title
		if len(title) == 0 {
			title = r.ShortName
		}

		arg := models.CreateTournamentParams{
			Name:         title,
			ShortName:    sql.NullString{String: r.ShortName, Valid: true},
			StartDate:    sql.NullTime{Time: startDate, Valid: true},
			EndDate:      sql.NullTime{Time: endDate, Valid: true},
			State:        state,
			NumDivisions: sql.NullInt32{Int32: int32(numDivisions), Valid: true},
			NumRounds:    sql.NullInt32{Int32: int32(numRounds), Valid: true},
			CreatorID:    creatorID,
			OrganiserID:  organiserID,
			DirectorID:   directorID,
			CoperatorID:  coperatorID,
			IsLocked:     isLocked,
			XID:          r.ID,
		}

		// tournament
		_, err := db.CreateTournament(context.Background(), arg)
		if err != nil {
			log.Fatal("failed to save model", err)
		}
		// fmt.Println(tournament)
		// cnt++
	}

	// report the number of tournaments added
	arg := models.ListTournamentsParams{
		Limit:  100000,
		Offset: 0,
	}
	tournaments, err := db.ListTournaments(context.Background(), arg)
	fmt.Printf("%d found, ", len(tournaments))

	fmt.Println("... ", time.Since(now))
}

func truncateTournaments() {
	db.TruncateTournaments(context.Background())
}
