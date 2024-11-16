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

// TODO Role permissions need to be verified with tests

func importAuthRoles() {
	filename := "seed/data/auth_roles.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("could not open data file", err)
	}
	defer f.Close()

	type AuthRolesCSV struct {
		Name               string `csv:"name"`
		CanLogin           string `csv:"can_login"`
		IsGuest            string `csv:"is_guest"`
		IsMembersAdmin     string `csv:"is_members_admin"`
		IsRatingsAdmin     string `csv:"is_ratings_admin"`
		IsTournamentsAdmin string `csv:"is_tournaments_admin"`
		IsSuperAdmin       string `csv:"is_super_admin"`
	}

	rows := []AuthRolesCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	fmt.Print("Import auth_roles... ")
	// Parse each role for its permissions
	now := time.Now()
	for _, r := range rows {

		CanLogin, err := strconv.ParseBool(r.CanLogin)
		if err != nil {
			log.Fatal(err)
		}

		IsGuest, err := strconv.ParseBool(r.IsGuest)
		if err != nil {
			log.Fatal(err)
		}

		IsMembersAdmin, err := strconv.ParseBool(r.IsMembersAdmin)
		if err != nil {
			log.Fatal(err)
		}

		IsRatingsAdmin, err := strconv.ParseBool(r.IsRatingsAdmin)
		if err != nil {
			log.Fatal(err)
		}

		IsTournamentsAdmin, err := strconv.ParseBool(r.IsTournamentsAdmin)
		if err != nil {
			log.Fatal(err)
		}

		IsSuperAdmin, err := strconv.ParseBool(r.IsSuperAdmin)
		if err != nil {
			log.Fatal(err)
		}

		arg := models.CreateAuthRoleParams{
			Name:               r.Name,
			CanLogin:           CanLogin,
			IsGuest:            IsGuest,
			IsMembersAdmin:     IsMembersAdmin,
			IsRatingsAdmin:     IsRatingsAdmin,
			IsTournamentsAdmin: IsTournamentsAdmin,
			IsSuperAdmin:       IsSuperAdmin,
		}

		// authRole
		_, err = db.CreateAuthRole(context.Background(), arg)
		if err != nil {
			log.Fatal("failed to save country", err)
		}
		// fmt.Println(authRole)
	}
	fmt.Println(time.Since(now))
}

func truncateAuthRole() {
	db.TruncateAuthRole(context.Background())
}
