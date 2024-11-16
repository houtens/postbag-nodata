package main

import (
	"flag"

	"github.com/houtens/postbag/internal/models"
)

var db *models.Queries

func main() {
	var fixtures = flag.Bool("fixtures", false, "seed fixtures data")
	var members = flag.Bool("members", false, "seed members data")
	var tournaments = flag.Bool("tournaments", false, "seed tournaments data")
	var invoices = flag.Bool("invoices", false, "seed invoices data")
	var ratings = flag.Bool("ratings", false, "seed ratings data and set end ratings")
	var results = flag.Bool("results", false, "seed results data")
	var fixwins = flag.Bool("fixwins", false, "fix wins data")
	var all = flag.Bool("all", false, "seed all data")

	flag.Parse()

	if *fixtures || *all {
		truncateCountries()
		truncateTitles()
		truncateAuthRole()
		truncateTournamentState()
		truncatePaymentTypes()
		truncateMembershipTypes()

		importCountries()
		importTitles()
		importAuthRoles()
		importClubs()
		importTournamentState()
		importMembershipTypes()
		importPaymentTypes()
	}

	if *members || *all {
		truncateContacts()
		truncateUsers()
		truncateMemberships()

		importMembers()
		importMemberships()
		setLifeMembership()
	}

	if *tournaments || *all {
		truncateTournaments()
		importTournaments()
	}

	if *invoices || *all {
		truncateInvoices()
		importInvoices()
	}

	if *ratings || *all {
		truncateRatings()
		importRatings()
		endRatings()
	}

	if *results || *all {
		truncateResults()
		importResults()
	}

	// These are idempotent
	if *fixwins || *all {
		validateTournamentResults()
		calculateRatingWins()
	}
}
