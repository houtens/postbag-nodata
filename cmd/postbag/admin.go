package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// TODO: count pending tournaments
// TODO: count scheduled tournaments
// TODO: count tournaments ready to be rated
// TODO: count users, members, active players etc
func (a app) getAdminDashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/dashboard", gin.H{})
}

func (a app) getAdminTournamentsList(c *gin.Context) {
	c.HTML(http.StatusOK, "admin/tournaments-list", gin.H{})
}

// TODO: list, new, edit

// func (a app) getAdminUsersList(c *gin.Context) {
// 	// users, err := a.session.GetUsers()
// 	c.String(http.StatusOK, "list admin users")
// }

// func (a app) getAdminUserNew(c *gin.Context) {
// 	c.String(http.StatusOK, "list admin users")
// }

// func (a app) postAdminUserNew(c *gin.Context) {
// 	c.String(http.StatusOK, "admin create")
// }

// func (a app) getAdminUserEdit(c *gin.Context) {
// 	c.String(http.StatusOK, "admin - edit user")
// }
