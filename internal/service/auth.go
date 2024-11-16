package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

var (
	ErrFailedAuthentication = errors.New("failed authentication")
	ErrRoleNotFound         = errors.New("role not found")
	ErrPasswordMismatch     = errors.New("passwords do not match")
)

type LoginFormInput struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type SessionData struct {
	UserID  string
	IsAdmin bool
}

func isAdmin(r models.AuthRole) bool {
	return r.IsMembersAdmin || r.IsClubsAdmin || r.IsRatingsAdmin || r.IsTournamentsAdmin || r.IsSuperAdmin
}

// Login sercice authenticates a user/password combination
func (s Service) Authenticate(input *LoginFormInput) (SessionData, error) {
	queries := models.New(s.DB)

	email := sql.NullString{String: input.Email, Valid: true}
	user, err := queries.GetUserByEmail(context.Background(), email)
	var loginSuccess = true
	if err != nil {
		loginSuccess = false
	}

	// Password validation
	ok := verifyPassword(user.PasswordHash.String, input.Password)
	if !ok {
		loginSuccess = false
	}

	// Only deny access after checking both email and password to prevent timing attack
	if !loginSuccess {
		return SessionData{}, ErrFailedAuthentication
	}

	// Lookup the ACL role
	role, err := queries.GetAuthRole(context.Background(), user.RoleID)
	if err != nil {
		return SessionData{}, ErrFailedAuthentication
	}

	sessionData := SessionData{
		UserID:  user.ID.String(),
		IsAdmin: isAdmin(role),
	}

	return sessionData, nil
}

func (s Service) GetUserRole(models.User) (models.AuthRole, error) {
	return models.AuthRole{}, nil
}

type EmailInput struct {
	Email string `form:"email"`
}

type PasswordResetInput struct {
	Password       string `form:"password"`
	PasswordRepeat string `form:"password-repeat"`
	User           string `form:"user"`
	Token          string `form:"token"`
}

func (s Service) GenerateResetToken(input EmailInput) error {
	queries := models.New(s.DB)

	// if none, return with an error
	arg := sql.NullString{
		String: input.Email,
		Valid:  true,
	}

	// Sanitise the input

	// Lookup the user
	user, err := queries.GetUserByEmail(context.Background(), arg)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Generate a reset token
	token, err := generateToken(32)
	if err != nil {
		return fmt.Errorf("token generation failed: %w", err)
	}

	// TODO: move this into config
	path := "http://localhost:3000/auth/reset?token="

	// Debug - in the absence of email we need to log the reset token
	fmt.Printf("Token: %s%s\n", path, token)

	arg1 := models.UpdateUserSetTokenParams{
		ID: user.ID,
		PwToken: sql.NullString{
			String: token,
			Valid:  true,
		},
	}
	_, err = queries.UpdateUserSetToken(context.Background(), arg1)
	if err != nil {
		// log
		return fmt.Errorf("unable to update user with token: %w", err)
	}

	return nil
}

// FindUserByValidToken returns the userid with matchign pw_token providing the pw_token_expiry is in the future (valid)
func (s Service) FindUserByValidToken(token string) (string, error) {
	queries := models.New(s.DB)

	arg := sql.NullString{
		String: token,
		Valid:  true,
	}

	// Fetch a user
	user, err := queries.GetUserByValidToken(context.Background(), arg)
	if err != nil {
		return "", err
	}

	// Return the uuid string value from the user id
	return user.ID.String(), nil
}

// ExpirePasswordResetToken resets the expiry time to now to invalidate the token
func (s Service) ExpirePasswordResetToken(token string) {
	queries := models.New(s.DB)

	arg := sql.NullString{
		String: token,
		Valid:  true,
	}

	_, err := queries.UpdateUserExpireToken(context.Background(), arg)
	if err != nil {
		log.Fatal(err)
	}
}

// ValidateResetToken returns true for tokens of correct length and containing only valid characters; otherwise false
func (s Service) ValidateResetToken(token string) bool {
	// validate length
	if len(token) != 32 {
		return false
	}

	// validate charset
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

	// validate the charset is okay
	return verifyCharset(token, charset)
}

// ResetUserPassword()
func (s Service) ResetUserPassword(user string, token string, password string, passwordRepeat string) error {
	queries := models.New(s.DB)

	// Do we need to doubly validate the passwords match?
	if password != passwordRepeat {
		// Passwords do not match
		return ErrPasswordMismatch
	}

	// Convert user from string to uuid
	userID, err := uuid.Parse(user)
	if err != nil {
		return err
	}

	// Convert token string to sql NullString
	nullToken := sql.NullString{
		String: token,
		Valid:  true,
	}

	pwHash := hashPassword(password)

	// Hash the password and create NullString
	passwordHash := sql.NullString{
		String: pwHash,
		Valid:  true,
	}

	// Update the password
	arg := models.UpdateUserPasswordHashParams{
		ID:           userID,
		PwToken:      nullToken,
		PasswordHash: passwordHash,
	}

	queries.UpdateUserPasswordHash(context.Background(), arg)

	return nil
}

func generateToken(length int) (string, error) {
	buf := make([]byte, length)

	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf)[:length], nil
}

func verifyCharset(input string, charset string) bool {
	for _, char := range input {
		if !strings.ContainsRune(charset, char) {
			return false
		}
	}
	return true
}
