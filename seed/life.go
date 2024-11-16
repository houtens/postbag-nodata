package main

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

func getFreePaymentType() uuid.UUID {
	paymentType, err := db.GetPaymentTypeByName(context.Background(), "Free")
	if err != nil {
		log.Fatal("could not get payment type free: ", err)
	}

	return paymentType.ID
}

// importLife adds an entry to memberships for life members
func setLifeMembership() {
	// select all users that are life members
	users, err := db.ListUsersByXLife(context.Background())
	if err != nil {
		log.Fatal("could not list users by xlife: ", err)
	}

	// list all membership types and store in map
	membershipTypes := readMembershipTypes()
	// get payment type id for free membership
	freePaymentType := getFreePaymentType()
	// expires in 100 years time
	expiresAt := time.Now().AddDate(100, 0, 0)

	var membershipTypeID uuid.UUID
	// for each user, lookup post or pdf life member in map to get membership type id
	for _, u := range users {
		if u.XPost {
			membershipTypeID = membershipTypes["life-post"]
		} else {
			membershipTypeID = membershipTypes["life-pdf"]
		}

		arg := models.CreateMembershipParams{
			UserID:           u.ID,
			Cost:             float32(0),
			MembershipTypeID: membershipTypeID,
			PaymentTypeID:    freePaymentType,
			ExpiresAt:        expiresAt,
		}
		_, err := db.CreateMembership(context.Background(), arg)
		if err != nil {
			log.Fatal("could not create membership: ", err)
		}
	}

	// do we have a life membership role for the user - if so update, or delete
}
