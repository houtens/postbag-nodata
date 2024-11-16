package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/houtens/postbag/internal/models"
)

func importMembershipTypes() {
	filename := "seed/data/membership_types.csv"
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	type MembershipTypesCSV struct {
		Name   string `csv:"name"`
		Code   string `csv:"code"`
		Years  string `csv:"years"`
		Junior string `csv:"junior"`
		Post   string `csv:"post"`
		Life   string `csv:"life"`
	}

	rows := []MembershipTypesCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	fmt.Print("Import membership_types... ")
	now := time.Now()
	for _, r := range rows {
		years, _ := strconv.Atoi(r.Years)
		junior, _ := strconv.ParseBool(r.Junior)
		post, _ := strconv.ParseBool(r.Post)
		life, _ := strconv.ParseBool(r.Life)

		arg := models.CreateMembershipTypeParams{
			Name:     r.Name,
			Code:     r.Code,
			NumYears: int32(years),
			IsJunior: junior,
			IsPost:   post,
			IsLife:   life,
		}
		_, err := db.CreateMembershipType(context.Background(), arg)
		if err != nil {
			// debug
			fmt.Println(r)
			log.Fatal(err)
		}
	}
	fmt.Println(time.Since(now))
}

func truncateMembershipTypes() {
	db.TruncateMembershipTypes(context.Background())
}
