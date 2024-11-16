package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// redirectProfile handles redirection from /profile to the user's own profile
func (a app) redirectProfile(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	// Get userid from the session in the context
	userID := session.Values["userid"]
	// Redirect to the user's own profile
	c.Redirect(http.StatusFound, fmt.Sprintf("/profile/%s", userID))
}

// showProfile handlers displaying a user profile
func (a app) showProfile(c *gin.Context) {
	session := c.MustGet("session").(*sessions.Session)
	isAdmin := session.Values["isadmin"]

	// Read the userID from the request path
	userID := c.Param("id")

	// Get user profile data
	userProfile, err := a.service.GetUserProfile(userID)
	if err != nil {
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	// Get upcoming tournaments data
	upcomingTournaments, err := a.service.GetUpcomingTournaments(userID)
	if err != nil {
		fmt.Println("error/404")
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	// Get recent ratings data
	recentRatings, err := a.service.GetRecentRatings(userID)
	if err != nil {
		fmt.Println("error/404")
		c.HTML(http.StatusNotFound, "error/404", gin.H{})
	}

	// Render the page
	c.HTML(http.StatusOK, "profile/detail", gin.H{
		"Title":       "Profile",
		"IsAdmin":     isAdmin,
		"UserID":      userID,
		"Profile":     userProfile,
		"Tournaments": upcomingTournaments,
		"Ratings":     recentRatings,
	})
}

// editProfile renders a form to allow some profile details to be updated
func (a app) editProfile(c *gin.Context) {
	// userID from params
	id := c.Param("id")
	// userID from the session
	session := c.MustGet("session").(*sessions.Session)
	// Get userid from the session in the context
	userID := session.Values["userid"].(string)

	// If the userID does not match the one being edited we redirect to the profile
	if userID != id {
		c.Redirect(http.StatusFound, "/profile")
		c.Abort()
		return
	}

	// user details
	user, err := a.service.GetEditUserProfile(userID)
	if err != nil {
		a.logger.Error("unable to get the user profile")
		c.Redirect(http.StatusFound, "/profile")
		c.Abort()
		return
	}

	// contact details
	contact, err := a.service.GetContactDetails(userID)
	if err != nil {
		a.logger.Error("unable to get the contact details for user")
		c.Redirect(http.StatusFound, "/profile")
		c.Abort()
		return
	}

	// clubs list for dropdown
	clubs, err := a.service.GetClubsList()
	if err != nil {
		a.logger.Error("unable to get the clubs list")
		c.Redirect(http.StatusFound, "/profile")
		c.Abort()
		return
	}

	// avatar link

	// Do we display the admin menu?
	isAdmin := session.Values["isadmin"]

	c.HTML(http.StatusOK, "profile/edit", gin.H{
		"Title":   "Edit profile",
		"IsAdmin": isAdmin,
		"User":    user,
		"Contact": contact,
		"Clubs":   clubs,
	})

}
