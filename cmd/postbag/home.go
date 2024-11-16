package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// getHome simply redirects to the login page and the login handler determines whether we have a valid session
func (a *app) getHome(c *gin.Context) {
	c.Redirect(http.StatusFound, "/auth/login")
}
