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

type CountryCSV struct {
	Flag     string `csv:"flag"`
	Name     string `csv:"name"`
	Code     string `csv:"code"`
	Priority string `csv:"priority"`
}

type OldCountryCSV struct {
	ID       string `csv:"id"`
	Name     string `csv:"name"`
	Filename string `csv:"filename"`
	Tier     string `csv:"tier"`
}

func readCountries() []CountryCSV {
	filename := "seed/data/countries.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("could not open countries data", err)
	}
	defer f.Close()

	rows := []CountryCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	return rows
}

func readOldCountries() []OldCountryCSV {
	filename := "seed/data/export/countries.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("could not open countries data", err)
	}
	defer f.Close()

	rows := []OldCountryCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	return rows
}

func importCountries() {

	rows := readCountries()
	old := readOldCountries()

	fmt.Print("Import countries... ")
	now := time.Now()
	for _, r := range rows {
		priority, err := strconv.ParseBool(r.Priority)
		if err != nil {
			log.Fatal("unable to parse bool from country:", err)
		}

		var xid sql.NullString

		for _, o := range old {
			if o.Name == r.Name {
				xid.String = o.ID
				xid.Valid = true
				break
			}
			if o.Name == "United States of America" && r.Name == "United States" {
				xid.String = o.ID
				xid.Valid = true
				break
			}
			if o.Name == "United Kingdom of Great Britain and Northern Ireland" && r.Name == "United Kingdom" {
				xid.String = o.ID
				xid.Valid = true
				break
			}
			if o.Name == "Czech Republic" && r.Name == "Czechia" {
				xid.String = o.ID
				xid.Valid = true
			}

		}

		// fmt.Printf("%v, %s\n", xid, r.Name)

		arg := models.CreateCountryParams{
			Name:     r.Name,
			Flag:     sql.NullString{String: r.Flag, Valid: true},
			Code:     sql.NullString{String: r.Code, Valid: true},
			Priority: priority,
			XID:      xid,
		}

		// country
		_, err = db.CreateCountry(context.Background(), arg)
		if err != nil {
			log.Fatal("failed to save country", err)
		}
		// fmt.Println(country)

	}

	// report how many were added
	countries, err := db.ListCountries(
		context.Background(),
		models.ListCountriesParams{
			Limit:  1000,
			Offset: 0,
		},
	)
	if err != nil {
		log.Fatal("unable to read countries")
	}
	fmt.Printf("%d found, ", len(countries))
	fmt.Println(time.Since(now))
}

func truncateCountries() {
	db.TruncateCountries(context.Background())
}
