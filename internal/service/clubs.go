package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

type Club struct {
	ID          string
	Name        string
	County      string
	ContactName string
	Country     string
}

type ClubsList []Club

func (s Service) GetClubsList() (ClubsList, error) {
	queries := models.New(s.DB)

	clubsList := ClubsList{}

	arg := models.ListClubsWithCountryParams{
		Offset: 0,
		Limit:  1000,
	}
	clubs, err := queries.ListClubsWithCountry(context.Background(), arg)
	if err != nil {
		return ClubsList{}, err
	}

	for _, c := range clubs {
		clubsList = append(
			clubsList,
			Club{
				ID:      c.ID.String(),
				Name:    c.Name,
				County:  c.County.String,
				Country: c.Country,
			})
	}
	return clubsList, nil
}

func (s Service) FilterClubs(query string) (ClubsList, error) {
	clubs, err := s.GetClubsList()
	if err != nil {
		return ClubsList{}, err
	}

	filteredClubs := ClubsList{}

	// Include only results with query match for Name, Country or County
	for _, c := range clubs {
		found := false
		if strings.Contains(strings.ToLower(c.Name), query) {
			found = true
		}
		if strings.Contains(strings.ToLower(c.Country), query) {
			found = true
		}
		if strings.Contains(strings.ToLower(c.County), query) {
			found = true
		}

		// Append if we have a match
		if found {
			filteredClubs = append(filteredClubs, c)
		}
	}

	return filteredClubs, nil
}

type ClubDetail struct {
	ClubName    string
	County      string
	Website     string
	Phone       string
	Email       string
	ContactName string
	Flag        string
	Code        string
}

func (s Service) GetClubDetail(id string) (ClubDetail, error) {
	queries := models.New(s.DB)

	// Parse the club id
	clubID, err := uuid.Parse(id)
	if err != nil {
		return ClubDetail{}, nil
	}

	// Retrieve club details with join on country
	club, err := queries.GetClubWithCountry(context.Background(), clubID)
	if err != nil {
		return ClubDetail{}, nil
	}

	clubDetail := ClubDetail{
		ClubName:    club.Name,
		County:      club.County.String,
		Website:     club.Website.String,
		Phone:       club.Phone.String,
		Email:       club.Email.String,
		ContactName: club.ContactName.String,
		Flag:        club.Flag.String,
		Code:        club.Code.String,
	}

	return clubDetail, nil
}

type ClubMember struct {
	ID     string
	Name   string
	Rating int
}

type ClubMembers []ClubMember

func (s Service) GetClubMembers(id string) (ClubMembers, error) {
	queries := models.New(s.DB)

	clubID, err := uuid.Parse(id)
	if err != nil {
		return ClubMembers{}, nil
	}

	// Retrieve club members with ratings
	arg := models.ListRatingsForClubParams{
		ClubID: uuid.NullUUID{UUID: clubID, Valid: true},
		Limit:  200,
		Offset: 0,
	}

	members, err := queries.ListRatingsForClub(context.Background(), arg)
	if err != nil {
		return ClubMembers{}, nil
	}
	fmt.Printf("%+v\n", members)

	clubMembers := []ClubMember{}

	for _, m := range members {
		memberName := fullName(m.FirstName, m.LastName, m.Title)
		clubMembers = append(
			clubMembers,
			ClubMember{
				ID:     m.UserID.String(),
				Name:   memberName,
				Rating: int(m.EndRating.Int32),
			},
		)
	}

	return clubMembers, nil
}
