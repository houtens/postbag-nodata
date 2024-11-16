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

type InvoiceCSV struct {
	ID         string `csv:"id"`
	Date       string `csv:"date"`
	Tournament string `csv:"tournament"`
	Players    string `csv:"players"`
	Games      string `csv:"games"`
	SuppMems   string `csv:"supp_mems"`
	Multiday   string `csv:"multiday"`
	Penalty    string `csv:"penalty"`
	Amount     string `csv:"amount"`
	Type       string `csv:"type"`
	Comment    string `csv:"comment"`
	Paid       string `csv:"paid"`
	Locked     string `csv:"locked"`
}

func setMultiday(s string) bool {
	return s != "0"
}

func setTotalCost(s string) float32 {
	sInt, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return float32(sInt) / 100
}

func importInvoices() {
	filename := "seed/data/export/invoices.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows := []InvoiceCSV{}
	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	// Read tournament id and xid into a map
	tournamentsMap := readTournaments()

	fmt.Print("Import invoices... ")

	now := time.Now()

	for _, r := range rows {
		numPlayers, _ := strconv.Atoi(r.Players)
		numGames, _ := strconv.Atoi(r.Games)
		multiday := setMultiday(r.Multiday)
		totalCost := setTotalCost(r.Amount)
		createdDate, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			log.Fatal(err)
		}
		tournamentID := tournamentsMap[r.Tournament]

		arg := models.SeedInvoiceParams{
			TournamentID:  tournamentID,
			NumPlayers:    int32(numPlayers),
			NumNonMembers: 0,
			NumGames:      int32(numGames),
			IsMultiday:    multiday,
			IsOverseas:    false,
			LevyCost:      0,
			ExtrasCost:    0,
			TotalCost:     totalCost,
			IsPaid:        true,
			Description:   sql.NullString{String: r.Comment, Valid: true},
			ExtrasComment: sql.NullString{Valid: false},
			Comment:       sql.NullString{Valid: false},
			CreatedAt:     createdDate,
			UpdatedAt:     createdDate,
		}

		_, err = db.SeedInvoice(context.Background(), arg)
		if err != nil {
			log.Fatal(err)
		}

	}
	fmt.Println(time.Since(now))
}

func truncateInvoices() {
	db.TruncateInvoices(context.Background())
}
