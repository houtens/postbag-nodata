package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/houtens/postbag/internal/service"
)

// login handles the login workflow; rendering the login form or redirecting authenticated users to the profile page
func (a app) login(c *gin.Context) {
	// If a valid session cookie is present then redirect to the profile detail
	session, err := store.Get(c.Request, SessionID)
	if err != nil {
		a.logger.Error("could not get a session from the store and this should always work")
		c.AbortWithError(500, err)
		return
	}

	// If already logged in then redirect profile
	if !session.IsNew {
		a.logger.Info("redirecting to profile as session exists")
		c.Redirect(http.StatusFound, "/profile")
		return
	}

	// Let profile do the parsing of userid and redirectign to the correct profile page
	c.HTML(http.StatusOK, "auth/login", gin.H{
		"Title": "Login",
	})
}

// validateLogin handles POST, checks credentials and establishes an authenticated user session
func (a app) validateLogin(c *gin.Context) {
	var input service.LoginFormInput

	// Read username and password from the form
	err := c.Bind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		return
	}

	// Authenticate the user/pass
	sessionData, err := a.service.Authenticate(&input)
	if err != nil {
		if errors.Is(err, service.ErrFailedAuthentication) {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		} else {
			// TODO: we should log this and handle it more gracefully
			c.String(http.StatusOK, "internal server error")
			c.Abort()
			return
		}
	}

	// TODO: thoroughly test
	// Create a new session
	session, err := store.Get(c.Request, SessionID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Save the session
	if session.IsNew {
		session.Values["userid"] = sessionData.UserID
		session.Values["isadmin"] = sessionData.IsAdmin
		session.Save(c.Request, c.Writer)
	}

	c.Set("session", session)

	// profile will redirect to the user profile
	c.Redirect(http.StatusFound, "/profile")
	c.Abort()
}

// logout handles logging a user out by clearing the session state
func (a app) logout(c *gin.Context) {
	// Logout
	session, err := store.Get(c.Request, SessionID)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	// Remove the session cookie by setting a negative maxage
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/auth/login")
	c.Abort()
}

// submitResetPassword handles resetting the user password
func (a app) submitResetPassword(c *gin.Context) {
	var input service.PasswordResetInput

	err := c.Bind(&input)
	if err != nil {
		// error binding data - no need for flash
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}

	// Validate that the form has matching password fields that meet the criteria
	err = a.service.ResetUserPassword(input.User, input.Token, input.Password, input.PasswordRepeat)
	if err != nil {
		// log and continue
		fmt.Println(err)
		c.Redirect(http.StatusFound, "/auth/login")
	}

	c.Redirect(http.StatusFound, "/auth/login")
	c.Abort()
}

// getResetPassword validates a reset token and renders a password reset form
func (a app) resetPassword(c *gin.Context) {

	token := c.Query("token")

	fmt.Println(len(token))

	// Validate the token is a valid token
	if !a.service.ValidateResetToken(token) {
		fmt.Println("invalid token")
	}

	// Lookup the user and ensure the token is valid
	user, err := a.service.FindUserByValidToken(token)
	if err != nil {
		c.Redirect(http.StatusFound, "/invalid")
		c.Abort()
		return
	}

	c.HTML(http.StatusOK, "auth/reset", gin.H{
		"Title": "Reset Password",
		"user":  user,
		"token": token,
	})
}

// forgotPassword renders the password form
func (a app) forgotPassword(c *gin.Context) {
	c.HTML(http.StatusOK, "auth/forgot", gin.H{
		"Title": "Forgot Password",
	})
}

// submitForgotPassword handles the issuing of token based password reset workflow
func (a app) submitForgotPassword(c *gin.Context) {
	var input service.EmailInput

	err := c.Bind(&input)
	if err != nil {
		c.Redirect(http.StatusFound, "/auth/login")
		c.Abort()
		return
	}

	fmt.Println(input)

	err = a.service.GenerateResetToken(input)
	if err != nil {
		// Do we just log this or redirect?
		fmt.Println(err)
	}

	// setFlash(c, "info:Check your email for a password reset link.")
	c.Redirect(http.StatusFound, "/auth/login")
	c.Abort()
}
