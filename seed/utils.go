package main

import (
	"context"
	"log"
	"regexp"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

func purgeNull(s string) string {
	if s == `\N` {
		return ""
	}
	return s
}

func splitAmbiguousNames(s string) (string, string) {
	// Just return the original string if we cannot parse a "Paul{Tranmere}"
	match, err := regexp.MatchString(`^.*{.*}$`, s)
	if err != nil {
		return s, ""
	}
	if !match {
		return s, ""
	}

	rx := regexp.MustCompile(`^(.*){(.*)}$`)
	res := rx.FindStringSubmatch(s)

	var firstName, altName string
	if len(res) > 2 {
		firstName = res[1]
		altName = res[2]
	}

	return firstName, altName
}

func readClubs() map[string]uuid.NullUUID {
	arg := models.ListClubsParams{
		Limit:  1000,
		Offset: 0,
	}
	clubs, err := db.ListClubs(context.Background(), arg)
	if err != nil {
		log.Fatal("error retrieving clubs")
	}

	// ClubsMap
	cm := make(map[string]uuid.NullUUID)
	for _, c := range clubs {
		cm[c.XID] = uuid.NullUUID{UUID: c.ID, Valid: true}
	}
	return cm
}

func readTournaments() map[string]uuid.UUID {
	arg := models.ListTournamentsParams{
		Limit:  100000,
		Offset: 0,
	}
	tournaments, err := db.ListTournaments(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	// Save to map
	tm := make(map[string]uuid.UUID)
	for _, t := range tournaments {
		tm[t.XID] = t.ID
	}

	return tm
}

func readTournamentStates() map[string]uuid.UUID {
	arg := models.ListTournamentStatesParams{
		Limit:  100,
		Offset: 0,
	}
	states, err := db.ListTournamentStates(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	// States Map
	sm := make(map[string]uuid.UUID)
	for _, s := range states {
		sm[s.Code] = s.ID
	}

	return sm
}

// return all the AuthRoles
func readAuthRoles() map[string]uuid.UUID {
	arg := models.ListAuthRolesParams{
		Limit:  100,
		Offset: 0,
	}
	authRoles, err := db.ListAuthRoles(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	// Return roles as map of name to id
	rm := make(map[string]uuid.UUID)
	for _, r := range authRoles {
		rm[r.Name] = r.ID
	}
	return rm
}

// returns all the Titles
func readTitles() map[string]uuid.UUID {
	// Lookup titles and set according to the CSV status field
	arg := models.ListTitlesParams{
		Limit:  5,
		Offset: 0,
	}
	titles, err := db.ListTitles(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}

	// Return tiles as a map of Name to ID
	tm := make(map[string]uuid.UUID)
	for _, t := range titles {
		tm[t.Name] = t.ID
	}
	return tm
}

func readMembershipTypes() map[string]uuid.UUID {
	arg := models.ListMembershipTypesParams{
		Limit:  1000,
		Offset: 0,
	}
	membershipTypes, err := db.ListMembershipTypes(context.Background(), arg)
	if err != nil {
		log.Fatal("failed to list membership_types: ", err)
	}

	// membership map of shortname code to ID
	mm := make(map[string]uuid.UUID)
	for _, m := range membershipTypes {
		mm[m.Code] = m.ID
	}

	return mm
}

func readPaymentTypes() map[string]uuid.UUID {
	arg := models.ListPaymentTypesParams{
		Limit:  20,
		Offset: 0,
	}
	paymentTypes, err := db.ListPaymentTypes(context.Background(), arg)
	if err != nil {
		log.Fatal("failed to list payment_types", err)
	}

	pm := make(map[string]uuid.UUID)
	for _, p := range paymentTypes {
		pm[p.Name] = p.ID
	}
	return pm
}

func countryByName(name string) uuid.NullUUID {
	country, err := db.GetCountryByName(context.Background(), name)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}
	return uuid.NullUUID{UUID: country.ID, Valid: true}
}

func tournamentByXID(xid string) uuid.NullUUID {
	u, err := db.GetTournamentByXID(context.Background(), xid)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}

	return uuid.NullUUID{UUID: u.ID, Valid: true}
}

func userByXID(xid string) uuid.NullUUID {
	u, err := db.GetUserByXID(context.Background(), xid)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}

	return uuid.NullUUID{UUID: u.ID, Valid: true}
}

func setLocked(s string) bool {
	return s == "1"
}
