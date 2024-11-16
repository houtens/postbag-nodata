package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomPaymentType(t *testing.T) PaymentType {
	name := randomString(16)
	paymentType, err := db.CreatePaymentType(context.Background(), name)
	require.NoError(t, err)
	require.NotEmpty(t, paymentType)

	return paymentType
}

func TestCreatePaymentType(t *testing.T) {
	createRandomPaymentType(t)
}

func TestGetPaymentType(t *testing.T) {
	paymentType1 := createRandomPaymentType(t)
	paymentType2, err := db.GetPaymentType(context.Background(), paymentType1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, paymentType2)
}

func TestListPaymentTypes(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomPaymentType(t)
	}

	arg := ListPaymentTypesParams{
		Limit:  5,
		Offset: 1,
	}

	paymentTypes, err := db.ListPaymentTypes(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, paymentTypes, 5)
}

func TestUpdatePaymentType(t *testing.T) {

	paymentType1 := createRandomPaymentType(t)
	arg := UpdatePaymentTypeParams{
		ID:   paymentType1.ID,
		Name: randomString(16),
	}

	paymentType2, err := db.UpdatePaymentType(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, paymentType2)
	require.NotEqual(t, paymentType1.Name, paymentType2.Name)
}

func TestDeletePaymentType(t *testing.T) {
	paymentType1 := createRandomPaymentType(t)
	db.DeletePaymentType(context.Background(), paymentType1.ID)

	paymentType2, err := db.GetPaymentType(context.Background(), paymentType1.ID)

	require.Error(t, err)
	require.Empty(t, paymentType2)
}
