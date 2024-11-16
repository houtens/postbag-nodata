package models

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCountry(t *testing.T) Country {
	arg := CreateCountryParams{
		Name:     randomCountry(),
		Flag:     sql.NullString{String: randomFlag(), Valid: true},
		Code:     sql.NullString{String: randomCountryCode(), Valid: true},
		Priority: false,
	}

	country, err := db.CreateCountry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, country)

	require.Equal(t, arg.Name, country.Name)
	require.Equal(t, arg.Flag, country.Flag)
	require.Equal(t, arg.Code, country.Code)
	require.Equal(t, arg.Priority, country.Priority)

	require.NotZero(t, country.ID)

	return country
}

func TestCreateCountry(t *testing.T) {
	createRandomCountry(t)
}

func TestGetCountry(t *testing.T) {
	country1 := createRandomCountry(t)
	country2, err := db.GetCountry(context.Background(), country1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, country2)

	require.Equal(t, country1.ID, country2.ID)
	require.Equal(t, country1.Name, country2.Name)
	require.Equal(t, country1.Flag, country2.Flag)
	require.Equal(t, country1.Code, country2.Code)
	require.Equal(t, country1.Priority, country2.Priority)
}

func TestListCountries(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCountry(t)
	}

	arg := ListCountriesParams{
		Limit:  5,
		Offset: 5,
	}

	countries, err := db.ListCountries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, countries, 5)

	for _, country := range countries {
		require.NotEmpty(t, country)
	}
}

func TestUpdateCountry(t *testing.T) {
	country1 := createRandomCountry(t)

	arg := UpdateCountryParams{
		ID:       country1.ID,
		Flag:     sql.NullString{String: randomFlag(), Valid: true},
		Code:     sql.NullString{String: randomCountryCode(), Valid: true},
		Priority: !country1.Priority,
	}

	country2, err := db.UpdateCountry(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, country2)

	require.Equal(t, country1.ID, country2.ID)
	require.Equal(t, arg.Flag, country2.Flag)
	require.Equal(t, arg.Code, country2.Code)
	require.Equal(t, arg.Priority, country2.Priority)
}

func TestDeleteCountry(t *testing.T) {
	country1 := createRandomCountry(t)
	err := db.DeleteCountry(context.Background(), country1.ID)
	require.NoError(t, err)

	country2, err := db.GetCountry(context.Background(), country1.ID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, country2)
}
