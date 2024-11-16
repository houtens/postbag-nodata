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

func importTitles() {
	filename := "seed/data/titles.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("could not open countries data", err)
	}
	defer f.Close()

	type TitleCSV struct {
		Name string `csv:"name"`
	}

	rows := []TitleCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	fmt.Print("Import titles... ")
	now := time.Now()
	for _, r := range rows {
		// title
		_, err := db.CreateTitle(context.Background(), r.Name)
		if err != nil {
			log.Fatal("failed to save title", err)
		}
	}

	// report how many were added
	arg := models.ListTitlesParams{
		Limit:  1000,
		Offset: 0,
	}
	titles, err := db.ListTitles(context.Background(), arg)
	fmt.Printf("%d found, ", len(titles))
	fmt.Println(time.Since(now))
}

func truncateTitles() {
	db.TruncateTitles(context.Background())
}
