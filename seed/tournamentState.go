package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/houtens/postbag/internal/models"
)

func importTournamentState() {
	filename := "seed/data/tournament_state.csv"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	type TournamentStateCSV struct {
		Name string `csv:"name"`
		Code string `csv:"code"`
	}

	rows := []TournamentStateCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Import tournament_state... ")
	now := time.Now()
	for _, r := range rows {
		arg := models.CreateTournamentStateParams{
			Name: r.Name,
			Code: r.Code,
		}
		_, err := db.CreateTournamentState(context.Background(), arg)
		if err != nil {
			// debug
			fmt.Println(r)
			log.Fatal(err)
		}
	}
	fmt.Println(time.Since(now))
}

func truncateTournamentState() {
	db.TruncateTournamentState(context.Background())
}
