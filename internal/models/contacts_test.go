package models

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
)

func createRandomContact(t *testing.T) Contact {
	user := createRandomUser(t)
	country := createRandomCountry(t)

	arg := CreateContactParams{
		UserID:    user.ID,
		Address1:  sql.NullString{String: randomName(), Valid: true},
		Address2:  sql.NullString{String: randomName(), Valid: true},
		Address3:  sql.NullString{String: randomName(), Valid: true},
		Address4:  sql.NullString{String: randomName(), Valid: true},
		Postcode:  sql.NullString{String: randomString(12), Valid: true},
		CountryID: uuid.NullUUID{UUID: country.ID, Valid: true},
		Phone:     sql.NullString{String: randomName(), Valid: true},
		Mobile:    sql.NullString{String: randomName(), Valid: true},
		Notes:     sql.NullString{String: randomString(300), Valid: true},
	}

	contact, err := db.CreateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact)

	return contact
}

func TestCreateContact(t *testing.T) {
	createRandomContact(t)
}

func TestGetContact(t *testing.T) {
	contact1 := createRandomContact(t)

	contact2, err := db.GetContact(context.Background(), contact1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ID, contact2.ID)
	require.Equal(t, contact1.Address1, contact2.Address1)
	require.Equal(t, contact1.Address2, contact2.Address2)
	require.Equal(t, contact1.Address3, contact2.Address3)
	require.Equal(t, contact1.Address4, contact2.Address4)
	require.Equal(t, contact1.Postcode, contact2.Postcode)
	require.Equal(t, contact1.Phone, contact2.Phone)
	require.Equal(t, contact1.Mobile, contact2.Mobile)
	require.Equal(t, contact1.UserID, contact2.UserID)
	require.Equal(t, contact1.Notes, contact2.Notes)
}

func TestListContacts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomContact(t)
	}

	arg := ListContactsParams{
		Limit:  5,
		Offset: 5,
	}

	contacts, err := db.ListContacts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, contacts, 5)
}

func TestUpdateContact(t *testing.T) {
	contact1 := createRandomContact(t)
	country := createRandomCountry(t)

	arg := UpdateContactParams{
		ID:        contact1.ID,
		UserID:    contact1.UserID,
		Address1:  sql.NullString{String: randomName(), Valid: true},
		Address2:  sql.NullString{String: randomName(), Valid: true},
		Address3:  sql.NullString{String: randomName(), Valid: true},
		Address4:  sql.NullString{String: randomName(), Valid: true},
		Postcode:  sql.NullString{String: randomString(12), Valid: true},
		CountryID: uuid.NullUUID{UUID: country.ID, Valid: true},
		Phone:     sql.NullString{String: randomName(), Valid: true},
		Mobile:    sql.NullString{String: randomName(), Valid: true},
		Notes:     sql.NullString{String: randomString(1000), Valid: true},
	}

	contact2, err := db.UpdateContact(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, contact2)

	require.Equal(t, contact1.ID, contact2.ID)
	require.NotEqual(t, contact1.Address1, contact2.Address1)
	require.NotEqual(t, contact1.Address2, contact2.Address2)
	require.NotEqual(t, contact1.Address3, contact2.Address3)
	require.NotEqual(t, contact1.Address4, contact2.Address4)
	require.NotEqual(t, contact1.Postcode, contact2.Postcode)
	require.NotEqual(t, contact1.Phone, contact2.Phone)
	require.NotEqual(t, contact1.Mobile, contact2.Mobile)
	require.NotEqual(t, contact1.Notes, contact2.Notes)
}

func TestDeleteContact(t *testing.T) {
	contact1 := createRandomContact(t)
	err := db.DeleteContact(context.Background(), contact1.ID)
	require.NoError(t, err)

	contact2, err := db.GetContact(context.Background(), contact1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, contact2)

}
