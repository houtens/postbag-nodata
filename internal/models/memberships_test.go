package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomMembership(t *testing.T) Membership {
	user := createRandomUser(t)
	membershipType := createRandomMembershipType(t)
	paymentType := createRandomPaymentType(t)

	arg := CreateMembershipParams{
		UserID:           user.ID,
		Cost:             float32(randomInt(0, 50)),
		MembershipTypeID: membershipType.ID,
		PaymentTypeID:    paymentType.ID,
		ExpiresAt:        randomDate(),
	}

	membership, err := db.CreateMembership(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, membership)

	return membership
}

func TestCreateMembership(t *testing.T) {
	createRandomMembership(t)
}

func TestGetMembership(t *testing.T) {
	membership1 := createRandomMembership(t)
	membership2, err := db.GetMembership(context.Background(), membership1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, membership2)
}

func TestListMemberships(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomMembership(t)
	}

	arg := ListMembershipsParams{
		Limit:  5,
		Offset: 1,
	}

	memberships, err := db.ListMemberships(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, memberships, 5)
}

func TestUpdateMembership(t *testing.T) {
	membership1 := createRandomMembership(t)

	user := createRandomUser(t)
	membershipType := createRandomMembershipType(t)
	paymentType := createRandomPaymentType(t)

	arg := UpdateMembershipParams{
		ID:               membership1.ID,
		UserID:           user.ID,
		Cost:             float32(randomInt(100, 200)),
		MembershipTypeID: membershipType.ID,
		PaymentTypeID:    paymentType.ID,
		ExpiresAt:        randomDate(),
	}

	membership2, err := db.UpdateMembership(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, membership2)
	require.NotEqual(t, membership1.UserID, membership2.UserID)
	require.NotEqual(t, membership1.Cost, membership2.Cost)
	require.NotEqual(t, membership1.MembershipTypeID, membership2.MembershipTypeID)
	require.NotEqual(t, membership1.PaymentTypeID, membership2.PaymentTypeID)
	require.NotEqual(t, membership1.ExpiresAt, membership2.ExpiresAt)
}

func TestDeleteMembership(t *testing.T) {
	membership1 := createRandomMembership(t)
	db.DeleteMembership(context.Background(), membership1.ID)

	membership2, err := db.GetMembership(context.Background(), membership1.ID)

	require.Error(t, err)
	require.Empty(t, membership2)
}
