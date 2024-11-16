package models

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomTitle(t *testing.T) Title {

	randomName := randomTitle()

	title, err := db.CreateTitle(context.Background(), randomName)
	require.NoError(t, err)
	require.NotEmpty(t, title)

	require.Equal(t, randomName, title.Name)

	require.NotZero(t, title.ID)

	return title
}

func TestCreateTitle(t *testing.T) {
	createRandomTitle(t)
}

func TestGetTitle(t *testing.T) {
	title1 := createRandomTitle(t)
	title2, err := db.GetTitle(context.Background(), title1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, title2)

	require.Equal(t, title1.ID, title2.ID)
	require.Equal(t, title1.Name, title2.Name)
	require.Equal(t, title1.CreatedAt, title2.CreatedAt)
	require.Equal(t, title1.UpdatedAt, title2.UpdatedAt)
}

func TestListTitles(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTitle(t)
	}

	arg := ListTitlesParams{
		Limit:  5,
		Offset: 5,
	}

	titles, err := db.ListTitles(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, titles, 5)

	for _, title := range titles {
		require.NotEmpty(t, title)
	}
}

func TestDeleteTitle(t *testing.T) {
	title1 := createRandomTitle(t)
	db.DeleteTitle(context.Background(), title1.ID)

	title2, err := db.GetTitle(context.Background(), title1.ID)

	// ensure error and that title does not exist
	require.Error(t, err)
	require.Empty(t, title2)

}
