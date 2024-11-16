package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

func getCountryID(name, county string) uuid.NullUUID {
	if county == "Scottish League" ||
		county == "Aberdeenshire" ||
		county == "Borders" ||
		county == "Dum & Gal" ||
		county == "East Lothian" ||
		county == "Fife" ||
		county == "Grampian" ||
		county == "Highland" ||
		county == "Lothian" ||
		county == "North Lanarkshire" ||
		county == "Perthshire" ||
		county == "Shetland" ||
		county == "Strathclyde" ||
		county == "Tayside" ||
		name == "ScottishLeague" ||
		name == "Glasgow Scrabble Club" ||
		name == "Edinburgh Scrabble Club" ||
		name == "ScottishScrabbleAssoc." {
		return countryByName("Scotland")
	}

	if name == "AllIrelandScrabbleAssoc." ||
		county == "EIRE" {
		return countryByName("Ireland")
	}

	if name == "Belfast" {
		return countryByName("Northern Ireland")
	}

	if county == "Clwyd" ||
		county == "Gwent" ||
		county == "M Glamorgan" ||
		county == "S Glamorgan" ||
		county == "W Glamorgan" {
		return countryByName("Wales")

	}

	return countryByName("England")
}

func importClubs() {
	filename := "seed/data/export/clubs.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("could not open data file", err)
	}
	defer f.Close()

	// import type definiton
	type ClubsCSV struct {
		Country     string `csv:"country"`
		County      string `csv:"county"`
		Name        string `csv:"club_name"`
		Type        string `csv:"type"`
		Contact     string `csv:"contact"`
		ContactName string `csv:"contact_name"`
		Email       string `csv:"email"`
		Phone       string `csv:"phone"`
		Website     string `csv:"website"`
		Updated     string `csv:"updated"`
		XID         string `csv:"id"`
	}

	rows := []ClubsCSV{}
	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	fmt.Print("Import clubs... ")
	now := time.Now()
	for _, r := range rows {
		// Purge null values "\N" from imported values
		r.County = purgeNull(r.County)
		r.Phone = purgeNull(r.Phone)
		r.Email = purgeNull(r.Email)

		// Match name/county to a country
		countryID := getCountryID(r.Name, r.County)

		arg := models.CreateClubParams{
			Name:        r.Name,
			County:      sql.NullString{String: r.County, Valid: true},
			IsActive:    true,
			Phone:       sql.NullString{String: r.Phone, Valid: true},
			Email:       sql.NullString{String: r.Email, Valid: true},
			ContactName: sql.NullString{String: r.ContactName, Valid: true},
			CountryID:   countryID,
			XID:         r.XID,
		}

		// club
		_, err := db.CreateClub(context.Background(), arg)
		if err != nil {
			log.Fatal("failed to save club", err)
		}
		// fmt.Println(club)
	}

	// report the number added
	arg := models.ListClubsParams{
		Limit:  10000,
		Offset: 0,
	}
	clubs, err := db.ListClubs(context.Background(), arg)
	fmt.Printf("%d found, ", len(clubs))

	fmt.Println(time.Since(now))
}

func truncateClubs() {
	db.TruncateClubs(context.Background())
}
