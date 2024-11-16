package models

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func createRandomInvoice(t *testing.T) Invoice {
	tournament := createRandomTournament(t)

	arg := CreateInvoiceParams{
		TournamentID:  tournament.ID,
		NumPlayers:    0,
		NumNonMembers: 0,
		NumGames:      0,
		IsMultiday:    randomBool(),
		IsOverseas:    randomBool(),
		LevyCost:      float32(randomInt(5, 50)),
		ExtrasCost:    0,
		TotalCost:     float32(randomInt(5, 50)),
		IsPaid:        randomBool(),
	}

	invoice, err := db.CreateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoice)

	return invoice
}

func TestCreateInvoice(t *testing.T) {
	createRandomInvoice(t)
}

func TestGetInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)

	invoice2, err := db.GetInvoice(context.Background(), invoice1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, invoice1.ID, invoice2.ID)
	require.Equal(t, invoice1.NumPlayers, invoice2.NumPlayers)
	require.Equal(t, invoice1.IsMultiday, invoice2.IsMultiday)
	// etc
}

func TestListInvoices(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomInvoice(t)
	}

	arg := ListInvoicesParams{
		Limit:  5,
		Offset: 5,
	}

	invoices, err := db.ListInvoices(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, invoices, 5)
}

func TestUpdateInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)

	arg := UpdateInvoiceParams{
		ID:            invoice1.ID,
		TournamentID:  invoice1.TournamentID,
		NumPlayers:    int32(randomInt(1, 100)),
		NumNonMembers: int32(randomInt(1, 100)),
		NumGames:      int32(randomInt(1, 100)),
		IsMultiday:    !invoice1.IsMultiday,
		IsOverseas:    !invoice1.IsOverseas,
		LevyCost:      float32(randomInt(60, 100)),
		ExtrasCost:    float32(randomInt(1, 100)),
		TotalCost:     float32(randomInt(100, 200)),
		IsPaid:        !invoice1.IsPaid,
		Description:   sql.NullString{Valid: false},
		ExtrasComment: sql.NullString{Valid: false},
		Comment:       sql.NullString{Valid: false},
	}
	invoice2, err := db.UpdateInvoice(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, invoice2)

	require.Equal(t, invoice1.ID, invoice2.ID)
	require.NotEqual(t, invoice1.NumPlayers, invoice2.NumPlayers)
	require.NotEqual(t, invoice1.NumNonMembers, invoice2.NumNonMembers)
	require.NotEqual(t, invoice1.NumGames, invoice2.NumGames)
	require.NotEqual(t, invoice1.IsMultiday, invoice2.IsMultiday)
	require.NotEqual(t, invoice1.IsOverseas, invoice2.IsOverseas)
	// etc
}

func TestDeleteInvoice(t *testing.T) {
	invoice1 := createRandomInvoice(t)
	err := db.DeleteInvoice(context.Background(), invoice1.ID)
	require.NoError(t, err)

	invoice2, err := db.GetInvoice(context.Background(), invoice1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, invoice2)
}
