package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// listClubs renders the full list of clubs
func (a app) listClubs(c *gin.Context) {
	clubs, err := a.service.GetClubsList()
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "clubs/list", gin.H{
		"Title":   "Clubs",
		"IsAdmin": isAdmin,
		"Clubs":   clubs,
	})
}

// showClub renders the detail for an individual club
func (a app) showClub(c *gin.Context) {
	// read the club id from the request path
	clubID := c.Param("id")

	club, err := a.service.GetClubDetail(clubID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	members, err := a.service.GetClubMembers(clubID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "clubs/detail", gin.H{
		"Title":   "Clubs detail",
		"IsAdmin": isAdmin,
		"Club":    club,
		"Members": members,
	})
}
