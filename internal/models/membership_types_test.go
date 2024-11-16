package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomMembershipType(t *testing.T) MembershipType {
	arg := CreateMembershipTypeParams{
		Name:     randomString(12),
		NumYears: int32(randomInt(1, 5)),
		IsJunior: randomBool(),
		IsPost:   randomBool(),
	}

	membershipType, err := db.CreateMembershipType(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, membershipType)

	return membershipType
}

func TestCreateMembershipType(t *testing.T) {
	createRandomMembershipType(t)
}

func TestGetMembershipType(t *testing.T) {
	membershipType1 := createRandomMembershipType(t)
	membershipType2, err := db.GetMembershipType(context.Background(), membershipType1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, membershipType2)
}

func TestListMembershipTypes(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomMembershipType(t)
	}

	arg := ListMembershipTypesParams{
		Limit:  5,
		Offset: 5,
	}

	users, err := db.ListMembershipTypes(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, users, 5)
}

func TestUpdateMembershipType(t *testing.T) {
	membershipType1 := createRandomMembershipType(t)

	arg := UpdateMembershipTypeParams{
		ID:       membershipType1.ID,
		Name:     randomString(12),
		NumYears: membershipType1.NumYears + 1,
		IsJunior: !membershipType1.IsJunior,
		IsPost:   !membershipType1.IsPost,
	}

	membershipType2, err := db.UpdateMembershipType(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, membershipType2)
	require.NotEqual(t, membershipType1.Name, membershipType2.Name)
	require.NotEqual(t, membershipType1.NumYears, membershipType2.NumYears)
	require.NotEqual(t, membershipType1.IsJunior, membershipType2.IsJunior)
	require.NotEqual(t, membershipType1.IsPost, membershipType2.IsPost)
}

func TestDeleteMembershipType(t *testing.T) {
	membershipType1 := createRandomMembershipType(t)
	// Delete membership type
	err := db.DeleteMembershipType(context.Background(), membershipType1.ID)
	require.NoError(t, err)
	// Try to read deleted membership type - should fail
	membershipType2, err := db.GetMembershipType(context.Background(), membershipType1.ID)
	require.Error(t, err)
	require.Empty(t, membershipType2)
}
