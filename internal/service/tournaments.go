package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

// Tournament struct returned to the handler and passed directly into the template
type Tournament struct {
	ID           string
	Name         string
	ShortName    string
	StartDate    time.Time
	EndDate      time.Time
	NumDivisions int
	NumRounds    int
	NumEntries   int
	IsLocked     bool
	// more we can add here; should join to state
}

// ListTournaments retunrs a list of Tournament for passing into the template
func (s Service) ListTournaments() ([]Tournament, error) {
	q := models.New(s.DB)

	// recent tournaments are the past 20 or over X years
	arg := models.ListRecentTournamentsParams{
		Limit:  100,
		Offset: 0,
	}

	// returns a models.Tournament slice
	tournaments, err := q.ListRecentTournaments(context.Background(), arg)
	if err != nil {
		fmt.Println(err)
	}

	// tournament slice
	results := []Tournament{}

	for _, t := range tournaments {
		results = append(
			results,
			Tournament{
				ID:           t.ID.String(),
				Name:         t.Name,
				ShortName:    t.ShortName.String,
				StartDate:    t.StartDate.Time,
				EndDate:      t.EndDate.Time,
				NumDivisions: int(t.NumDivisions.Int32),
				NumRounds:    int(t.NumRounds.Int32),
				NumEntries:   int(t.NumEntries.Int32),
				IsLocked:     t.IsLocked,
			},
		)
	}

	return results, nil
}

// FetchTournament queries the data and returns a Tournament struct
func (s Service) FetchTournament(id string) (Tournament, error) {
	q := models.New(s.DB)

	// Parse id to uuid type
	userID, err := uuid.Parse(id)
	if err != nil {
		return Tournament{}, err
	}

	t, err := q.GetTournament(context.Background(), userID)
	if err != nil {
		fmt.Println("could not retrieve tournament")
	}

	result := Tournament{
		ID:           t.ID.String(),
		Name:         t.Name,
		ShortName:    t.ShortName.String,
		StartDate:    t.StartDate.Time,
		EndDate:      t.EndDate.Time,
		NumDivisions: int(t.NumDivisions.Int32),
		NumRounds:    int(t.NumRounds.Int32),
		NumEntries:   int(t.NumEntries.Int32),
	}

	return result, nil
}

type ActiveUser struct {
	ID        string
	FirstName string
	LastName  string
}

func (s Service) FilterActiveUsers(query string) ([]ActiveUser, error) {
	q := models.New(s.DB)

	users, err := q.GetActiveUsers(context.Background())
	if err != nil {
		fmt.Println("failed to query active users")
		return []ActiveUser{}, err
	}

	au := []ActiveUser{}
	fmt.Printf("%+v\n", au)

	for _, u := range users {
		if strings.Contains(strings.ToLower(u.FirstName), strings.ToLower(query)) ||
			strings.Contains(strings.ToLower(u.LastName), strings.ToLower(query)) {

			fmt.Println("matched:", strings.ToLower(u.FirstName), strings.ToLower(query))

			au = append(au, ActiveUser{
				ID:        u.ID.String(),
				FirstName: u.FirstName,
				LastName:  u.LastName,
			})
		}
	}

	return au, nil
}
