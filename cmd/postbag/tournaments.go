package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

// listTournaments handles rendering the list of tournaments
func (a app) listTournaments(c *gin.Context) {
	tournaments, err := a.service.ListTournaments()
	if err != nil {
		c.HTML(http.StatusNotFound, "server error", gin.H{})
	}

	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "tournaments/list", gin.H{
		"Title":       "Tournaments",
		"IsAdmin":     isAdmin,
		"Tournaments": tournaments,
	})
}

// showTournament handles rendering the details of a tournament
func (a app) showTournament(c *gin.Context) {
	userID := c.Param("id")

	tournament, err := a.service.FetchTournament(userID)
	if err != nil {
		c.Redirect(http.StatusNotFound, "/tournaments")
		c.Abort()
		return
	}

	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "tournaments/detail", gin.H{
		"Title":      "Tournament detail",
		"IsAdmin":    isAdmin,
		"Tournament": tournament,
	})
}

// TODO: editTournament handles editing an existing tournament
func (a app) editTournament(c *gin.Context) {
	// Parse the tournament ID from request path as UUID
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, "Not found")
	}
	_ = id

	c.String(http.StatusOK, "edit tournament by id")
}

// createTournament renders the form to create a tournament
func (a app) createTournament(c *gin.Context) {
	c.HTML(http.StatusOK, "tournaments/create", gin.H{})
}

func validateName(name string) error {
	if len(name) < 3 {
		return errors.New("name is too short")
	}
	return nil
}

// NewTournament
type NewTournament struct {
	Name             string    `form:"name" binding:"required"`
	StartDate        time.Time `form:"start_date" time_format:"2006-01-02" binding:"required"`
	EndDate          time.Time `form:"end_date" time_format:"2006-01-02" binding:"required"`
	Rounds           int       `form:"rounds"`
	Divisions        int       `form:"divisions"`
	ABSPRated        bool      `form:"absp"`
	WESPARated       bool      `form:"wespa"`
	Invitational     bool      `form:"invitational"`
	FreeChallenge    bool      `form:"free"`
	PenaltyChallenge bool      `form:"penalty"`
	// TODO: add postcode for the venue
	// needs Organiser, Director, Computer Operator
}

// saveTournament parses the create form and saves to the database
func (a app) saveTournament(c *gin.Context) {

	var input NewTournament
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// fmt.Printf("%+v\n", input)

	c.Redirect(http.StatusFound, "/tournaments")
}
