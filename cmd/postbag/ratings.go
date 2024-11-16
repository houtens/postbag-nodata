package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// listRatings handles redirection from /profile
func (a app) listRatings(c *gin.Context) {
	ratings, err := a.service.GetRatingsList()
	if err != nil {
		fmt.Println("error/404")
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "ratings/list", gin.H{
		"Title":   "Ratings",
		"IsAdmin": isAdmin,
		"Ratings": ratings,
	})
}
