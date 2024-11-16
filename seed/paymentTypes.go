package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

func importPaymentTypes() {
	filename := "seed/data/payment_types.csv"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	type PaymentTypesCSV struct {
		Name string `csv:"name"`
	}

	rows := []PaymentTypesCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Import payment_types... ")
	now := time.Now()
	for _, r := range rows {
		_, err := db.CreatePaymentType(context.Background(), r.Name)
		if err != nil {
			// debug
			fmt.Println(r)
			log.Fatal(err)
		}
	}
	fmt.Println(time.Since(now))
}

func truncatePaymentTypes() {
	db.TruncatePaymentTypes(context.Background())
}
