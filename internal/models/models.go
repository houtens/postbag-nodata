// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0

package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID        uuid.UUID `db:"id"`
	Event     int32     `db:"event"`
	UserID    uuid.UUID `db:"user_id"`
	RemoteIp  string    `db:"remote_ip"`
	UserAgent string    `db:"user_agent"`
	Message   string    `db:"message"`
	CreatedAt time.Time `db:"created_at"`
}

type AuthRole struct {
	ID                 uuid.UUID `db:"id"`
	Name               string    `db:"name"`
	CanLogin           bool      `db:"can_login"`
	IsGuest            bool      `db:"is_guest"`
	IsMembersAdmin     bool      `db:"is_members_admin"`
	IsClubsAdmin       bool      `db:"is_clubs_admin"`
	IsRatingsAdmin     bool      `db:"is_ratings_admin"`
	IsTournamentsAdmin bool      `db:"is_tournaments_admin"`
	IsSuperAdmin       bool      `db:"is_super_admin"`
	CreatedAt          time.Time `db:"created_at"`
	UpdatedAt          time.Time `db:"updated_at"`
}

type Club struct {
	ID          uuid.UUID      `db:"id"`
	Name        string         `db:"name"`
	County      sql.NullString `db:"county"`
	Website     sql.NullString `db:"website"`
	IsActive    bool           `db:"is_active"`
	Phone       sql.NullString `db:"phone"`
	Email       sql.NullString `db:"email"`
	ContactName sql.NullString `db:"contact_name"`
	CountryID   uuid.NullUUID  `db:"country_id"`
	XID         string         `db:"x_id"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
}

type Contact struct {
	ID        uuid.UUID      `db:"id"`
	Address1  sql.NullString `db:"address1"`
	Address2  sql.NullString `db:"address2"`
	Address3  sql.NullString `db:"address3"`
	Address4  sql.NullString `db:"address4"`
	Postcode  sql.NullString `db:"postcode"`
	CountryID uuid.NullUUID  `db:"country_id"`
	Phone     sql.NullString `db:"phone"`
	Mobile    sql.NullString `db:"mobile"`
	UserID    uuid.UUID      `db:"user_id"`
	Notes     sql.NullString `db:"notes"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type Country struct {
	ID       uuid.UUID      `db:"id"`
	Name     string         `db:"name"`
	Flag     sql.NullString `db:"flag"`
	Code     sql.NullString `db:"code"`
	Priority bool           `db:"priority"`
	XID      sql.NullString `db:"x_id"`
}

type Invoice struct {
	ID            uuid.UUID      `db:"id"`
	TournamentID  uuid.UUID      `db:"tournament_id"`
	NumPlayers    int32          `db:"num_players"`
	NumNonMembers int32          `db:"num_non_members"`
	NumGames      int32          `db:"num_games"`
	IsMultiday    bool           `db:"is_multiday"`
	IsOverseas    bool           `db:"is_overseas"`
	LevyCost      float32        `db:"levy_cost"`
	ExtrasCost    float32        `db:"extras_cost"`
	TotalCost     float32        `db:"total_cost"`
	IsPaid        bool           `db:"is_paid"`
	Description   sql.NullString `db:"description"`
	ExtrasComment sql.NullString `db:"extras_comment"`
	Comment       sql.NullString `db:"comment"`
	CreatedAt     time.Time      `db:"created_at"`
	UpdatedAt     time.Time      `db:"updated_at"`
}

type Membership struct {
	ID               uuid.UUID `db:"id"`
	UserID           uuid.UUID `db:"user_id"`
	Cost             float32   `db:"cost"`
	MembershipTypeID uuid.UUID `db:"membership_type_id"`
	PaymentTypeID    uuid.UUID `db:"payment_type_id"`
	ExpiresAt        time.Time `db:"expires_at"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type MembershipType struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	NumYears  int32     `db:"num_years"`
	IsJunior  bool      `db:"is_junior"`
	IsPost    bool      `db:"is_post"`
	IsLife    bool      `db:"is_life"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type PaymentType struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Rating struct {
	ID            uuid.UUID     `db:"id"`
	UserID        uuid.UUID     `db:"user_id"`
	TournamentID  uuid.UUID     `db:"tournament_id"`
	Division      int32         `db:"division"`
	NumGames      sql.NullInt32 `db:"num_games"`
	StartRating   sql.NullInt32 `db:"start_rating"`
	EndRating     sql.NullInt32 `db:"end_rating"`
	RatingPoints  sql.NullInt32 `db:"rating_points"`
	OppRatingsSum sql.NullInt32 `db:"opp_ratings_sum"`
	NumWins       float32       `db:"num_wins"`
	IsLocked      bool          `db:"is_locked"`
	XID           string        `db:"x_id"`
	CreatedAt     time.Time     `db:"created_at"`
	UpdatedAt     time.Time     `db:"updated_at"`
}

type Result struct {
	ID           uuid.UUID     `db:"id"`
	Player1ID    uuid.NullUUID `db:"player1_id"`
	Player2ID    uuid.NullUUID `db:"player2_id"`
	Score1       int32         `db:"score1"`
	Score2       int32         `db:"score2"`
	Spread       int32         `db:"spread"`
	TournamentID uuid.UUID     `db:"tournament_id"`
	RoundNum     int32         `db:"round_num"`
	Type         int32         `db:"type"`
	IsLocked     bool          `db:"is_locked"`
	XID          string        `db:"x_id"`
	CreatedAt    time.Time     `db:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"`
}

type Title struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type Tournament struct {
	ID             uuid.UUID      `db:"id"`
	Name           string         `db:"name"`
	ShortName      sql.NullString `db:"short_name"`
	StartDate      sql.NullTime   `db:"start_date"`
	EndDate        sql.NullTime   `db:"end_date"`
	State          uuid.UUID      `db:"state"`
	NumDivisions   sql.NullInt32  `db:"num_divisions"`
	NumRounds      sql.NullInt32  `db:"num_rounds"`
	NumEntries     sql.NullInt32  `db:"num_entries"`
	IsPc           bool           `db:"is_pc"`
	IsFc           bool           `db:"is_fc"`
	IsRr           bool           `db:"is_rr"`
	IsWespa        bool           `db:"is_wespa"`
	IsInvitational bool           `db:"is_invitational"`
	IsLocked       bool           `db:"is_locked"`
	CreatorID      uuid.NullUUID  `db:"creator_id"`
	OrganiserID    uuid.NullUUID  `db:"organiser_id"`
	DirectorID     uuid.NullUUID  `db:"director_id"`
	CoperatorID    uuid.NullUUID  `db:"coperator_id"`
	XID            string         `db:"x_id"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      time.Time      `db:"updated_at"`
}

type TournamentState struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type User struct {
	ID            uuid.UUID      `db:"id"`
	FirstName     string         `db:"first_name"`
	LastName      string         `db:"last_name"`
	AltName       sql.NullString `db:"alt_name"`
	Email         sql.NullString `db:"email"`
	PasswordHash  sql.NullString `db:"password_hash"`
	AbspNum       sql.NullInt32  `db:"absp_num"`
	ClubID        uuid.NullUUID  `db:"club_id"`
	TitleID       uuid.NullUUID  `db:"title_id"`
	RoleID        uuid.UUID      `db:"role_id"`
	XLife         bool           `db:"x_life"`
	XPost         bool           `db:"x_post"`
	XID           string         `db:"x_id"`
	IsDeceased    bool           `db:"is_deceased"`
	Avatar        sql.NullString `db:"avatar"`
	PwToken       sql.NullString `db:"pw_token"`
	PwTokenExpiry sql.NullTime   `db:"pw_token_expiry"`
}
