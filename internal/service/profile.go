package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

// UpcomingTournaments are yet to be played
type UpcomingTournament struct {
	Name      string
	StartDate time.Time
	NumRounds int
}

type UpcomingTournamentsOutput []UpcomingTournament

func (s Service) GetUpcomingTournaments(id string) (UpcomingTournamentsOutput, error) {
	queries := models.New(s.DB)

	// Parse id to uuid type
	userID, err := uuid.Parse(id)
	if err != nil {
		return UpcomingTournamentsOutput{}, err
	}

	upcomingTournaments := UpcomingTournamentsOutput{}

	upcoming, err := queries.ListUpcomingTournamentsForPlayer(context.Background(), userID)
	if err != nil {
		return UpcomingTournamentsOutput{}, err
	}

	for _, u := range upcoming {
		upcomingTournaments = append(
			upcomingTournaments,
			UpcomingTournament{
				Name:      u.Name,
				StartDate: u.StartDate.Time,
				NumRounds: int(u.NumRounds.Int32),
			})
	}
	return upcomingTournaments, nil
}

// RecentResults are completed tournaments
type RecentRatings struct {
	ID          uuid.UUID
	Tournament  string
	StartDate   time.Time
	NumGames    int
	NumWins     float32
	StartRating int
	EndRating   int
}

type RecentRatingsOutput []RecentRatings

func (s Service) GetRecentRatings(id string) (RecentRatingsOutput, error) {
	queries := models.New(s.DB)
	userID, err := uuid.Parse(id)
	if err != nil {
		return RecentRatingsOutput{}, fmt.Errorf("parse id error: %w", err)
	}

	recentRatings := RecentRatingsOutput{}
	recent, err := queries.ListRecentRatingsForPlayer(context.Background(), userID)
	if err != nil {
		return RecentRatingsOutput{}, fmt.Errorf("fetch recent results error: %w", err)
	}
	_ = recent

	for _, r := range recent {
		recentRatings = append(
			recentRatings,
			RecentRatings{
				ID:         r.ID,
				Tournament: r.Name,
				StartDate:  r.StartDate.Time,
				NumGames:   int(r.NumGames.Int32),
				// TODO: NumWins needs to be seeded
				NumWins:     r.NumWins,
				StartRating: int(r.StartRating.Int32),
				EndRating:   int(r.EndRating.Int32),
			})
	}

	return recentRatings, nil
}

type UserProfileOutput struct {
	FullName     string
	AltName      string
	ABSPNum      string
	Avatar       string
	ClubName     string
	ClubID       uuid.UUID
	ValidMember  int
	Rating       int
	Won          int
	Lost         int
	Drawn        int
	WinRate      string
	ScoreFor     string
	ScoreAgainst string
}

func (s Service) GetUserProfile(id string) (UserProfileOutput, error) {
	queries := models.New(s.DB)
	userID, err := uuid.Parse(id)
	if err != nil {
		return UserProfileOutput{}, err
	}

	userProfile := UserProfileOutput{}

	profile, err := queries.GetProfile(context.Background(), userID)
	if err != nil {
		return UserProfileOutput{}, err
	}

	// set TitleName with or without "(GM)" or "(Exp)"
	userProfile.FullName = fmt.Sprintf("%s %s", profile.FirstName, profile.LastName)
	if profile.TitleName.String != "" {
		userProfile.FullName = fmt.Sprintf("%s %s (%s)", profile.FirstName, profile.LastName, profile.TitleName.String)
	}
	userProfile.Avatar = profile.Avatar.String
	userProfile.ABSPNum = fmt.Sprintf("%d", profile.AbspNum.Int32)
	userProfile.ClubName = profile.ClubName.String
	userProfile.ClubID = profile.ClubID.UUID

	rating, err := queries.GetLatestRating(context.Background(), userID)
	if err != nil {
		return UserProfileOutput{}, err
	}
	userProfile.Rating = int(rating.Int32)

	nullUserID := uuid.NullUUID{
		UUID:  userID,
		Valid: true,
	}

	record, err := queries.GetWinLossDraw(context.Background(), nullUserID)
	if err != nil {
		return UserProfileOutput{}, err
	}
	userProfile.Won = int(record.Won)
	userProfile.Lost = int(record.Lost)
	userProfile.Drawn = int(record.Drawn)

	// Win percentage
	winRate := calcWinRate(record.Won, record.Lost, record.Drawn)
	userProfile.WinRate = winRate

	averages, err := queries.GetAverageScores(context.Background(), nullUserID)
	if err != nil {
		return UserProfileOutput{}, err
	}
	userProfile.ScoreFor = fmt.Sprintf("%.0f", averages.For)
	userProfile.ScoreAgainst = fmt.Sprintf("%.0f", averages.Against)

	// recent
	_, err = queries.ListRecentRatingsForPlayer(context.Background(), userID)
	if err != nil {
		return UserProfileOutput{}, err
	}

	// TODO: this isn't calculated
	userProfile.ValidMember = 1
	// validMember, err := queries.GetValidMembership(context.Background(), userID)
	// if err != nil {
	// 	validMember = -1
	// }

	return userProfile, nil
}

type EditUserProfile struct {
	ID        string
	FirstName string
	LastName  string
	AltName   string
	Email     string
	ClubID    string
	TitleID   string
	AbspNum   int32
}

func (s Service) GetEditUserProfile(id string) (EditUserProfile, error) {

	// Get uuid from the path
	uuid, err := uuid.Parse(id)
	if err != nil {
		return EditUserProfile{}, errors.New("cannot parse id as uuid")
	}

	queries := models.New(s.DB)

	user, err := queries.GetUser(context.Background(), uuid)
	if err != nil {
		return EditUserProfile{}, errors.New("user not found")
	}
	// fmt.Printf("%+v\n", user)

	editUserProfile := EditUserProfile{
		ID:        id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AltName:   user.AltName.String,
		Email:     user.Email.String,
		AbspNum:   user.AbspNum.Int32,
		ClubID:    user.ClubID.UUID.String(),
		TitleID:   user.TitleID.UUID.String(),
	}

	return editUserProfile, nil
}

type ContactOutput struct {
	Address1 string
	Address2 string
	Address3 string
	Address4 string
	Postcode string
	Phone    string
	Mobile   string
}

func (s Service) GetContactDetails(id string) (ContactOutput, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return ContactOutput{}, nil
	}
	fmt.Println(uuid)

	queries := models.New(s.DB)

	c, err := queries.GetContactByUserID(context.Background(), uuid)
	if err != nil {
		fmt.Println("error", err)
		return ContactOutput{}, nil
	}

	contactOutput := ContactOutput{
		Address1: c.Address1.String,
		Address2: c.Address2.String,
		Address3: c.Address3.String,
		Address4: c.Address4.String,
		Postcode: c.Postcode.String,
		Phone:    c.Phone.String,
		Mobile:   c.Mobile.String,
	}

	fmt.Printf("%+v\n", contactOutput)
	return contactOutput, nil
}
