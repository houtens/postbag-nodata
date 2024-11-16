package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

// MembershipCSV is used to parse the csv data (table payments)
type MembershipCSV struct {
	XID              string `csv:"id"`
	Date             string `csv:"date"`
	Member           string `csv:"member"`
	Amount           string `csv:"amount"`
	Method           string `csv:"method"`
	SubscriptionYear string `csv:"subscription_year"`
}

func getUserByXID(xid string) (models.User, error) {
	// query to lookup the user by XID
	user, err := db.GetUserByXID(context.Background(), xid)
	return user, err
}

func getMembershipType(u models.User, m map[string]uuid.UUID) uuid.UUID {
	// Ignore life memberships here - any life members that have a membership
	// were paid for before becoming life members.
	if u.XPost {
		return m["legacy-post"]
	}
	return m["legacy-pdf"]
}

// func getMembership
func getExpiresDate(expiry string) time.Time {
	// Parse the year string into a time.Time
	date, _ := time.Parse("2006", expiry)
	// Add one year and subtract one second
	date = date.AddDate(1, 0, 0).Add(time.Second * -1)
	return date
}

func getPaymentType(p string, m map[string]uuid.UUID) uuid.UUID {
	switch p {
	case "paypal":
		return m["PayPal"]
	case "bacs":
		return m["BACS"]
	case "BACS":
		return m["BACS"]
	case "cash":
		return m["Cash"]
	case "Cash":
		return m["Cash"]
	case "cheque":
		return m["Cheque"]
	case "standing_order":
		return m["Standing Order"]
	case "dd":
		return m["Direct Debit"]
	case "free":
		return m["Free"]
	}
	return uuid.UUID{}
}

func importMemberships() {
	filename := "seed/data/export/payments.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows := []MembershipCSV{}
	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	paymentTypesMap := readPaymentTypes()
	membershipTypesMap := readMembershipTypes()

	fmt.Print("Import memberships... ")
	now := time.Now()

	for _, r := range rows {
		// fmt.Printf("%+v\n", r)
		// Special cases to ignore
		if r.Method == "" || r.SubscriptionYear == "0" {
			continue
		}
		user, err := getUserByXID(r.Member)
		if err != nil {
			// Member not found for the associated payment
			// fmt.Printf("missing member %v - ignore payment\n", r.Member)
			continue
		}

		// We are processing paid memberships - not life members.
		membershipTypeID := getMembershipType(user, membershipTypesMap)

		cost, _ := strconv.ParseFloat(r.Amount, 32)
		expiresAt := getExpiresDate(r.SubscriptionYear)
		paymentType := getPaymentType(r.Method, paymentTypesMap)
		createdAt, err := time.Parse("2006-01-02", r.Date)
		if err != nil {
			log.Fatal(err)
		}

		arg := models.SeedMembershipParams{
			UserID:           user.ID,
			Cost:             float32(cost),
			MembershipTypeID: membershipTypeID,
			PaymentTypeID:    paymentType,
			ExpiresAt:        expiresAt,
			CreatedAt:        createdAt,
			UpdatedAt:        createdAt,
		}

		membership, err := db.SeedMembership(context.Background(), arg)
		if err != nil {
			log.Fatal(err)
		}
		_ = membership
	}

	// report how many were added
	arg := models.ListMembershipsParams{
		Limit:  10000,
		Offset: 0,
	}
	memberships, err := db.ListMemberships(context.Background(), arg)
	fmt.Printf("%d found, ", len(memberships))

	fmt.Println(time.Since(now))

}

func truncateMemberships() {
	db.TruncateMemberships(context.Background())
}
