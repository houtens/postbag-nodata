package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// filterRatings returns an htmx partial ratings list filtered by the query param
func (a app) filterRatings(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))

	ratings, err := a.service.FilterRatings(query)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	c.HTML(http.StatusOK, "ratings/filtered-list", gin.H{
		"Ratings": ratings,
	})
}

// filterClubs renders a partial htmx response with a filtered list of clubs
func (a app) filterClubs(c *gin.Context) {
	query := strings.ToLower(c.Query("q"))

	clubs, err := a.service.FilterClubs(query)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	c.HTML(http.StatusOK, "clubs/filtered-list", gin.H{
		"Clubs": clubs,
	})

}

// filterOrganisers returns a list of active (not deceased) users with partial match for first, last or alt names
func (a app) filterOrganisers(c *gin.Context) {
	query := strings.ToLower(c.Query("organiser"))

	organisers, err := a.service.FilterActiveUsers(query)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	// organisers, directors and coperators are all rendered from the active-users template
	c.HTML(http.StatusOK, "active-users/filtered-list", gin.H{
		"Organisers": organisers,
	})

}
