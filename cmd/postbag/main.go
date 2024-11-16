package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/houtens/postbag/config"
	"github.com/houtens/postbag/internal/service"

	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

const SessionID string = "POSTBAGSESSION"

// TODO: Add to the config
var store = sessions.NewCookieStore([]byte("something-super-secret-goes-here"))

type app struct {
	logger  *slog.Logger
	service *service.Service
}

func main() {
	// Load the config and connect to the database
	cfg := config.Load()
	db := postgresConnection(cfg)

	// Create a slog logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Initialise an app struct
	app := &app{
		logger: logger,
		service: &service.Service{
			DB: db,
		},
	}

	// Create gin router and register routes (and setup)
	r := gin.Default()

	// Set once for the session cookies
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   60 * 60 * 12, // Login session is valid for 12 hours
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}

	// Serve static files on the right path
	r.Use(static.Serve("/static/css", static.LocalFile("./static/css", false)))
	r.Use(static.Serve("/static/js", static.LocalFile("./static/js", false)))
	r.Use(static.Serve("/static/avatars", static.LocalFile("./static/images/avatars", false)))
	r.Use(static.Serve("/images", static.LocalFile("./static/images", false)))

	r.GET("/", app.getHome)

	// Unauthenticated routes - must not be logged in
	anon := r.Group("/")
	anon.Use(UnauthenticatedMiddleware())

	// Login
	anon.GET("/auth/login", app.login)
	anon.POST("/auth/login", app.validateLogin)

	// Forgot password and password reset
	anon.GET("/auth/forgot", app.forgotPassword)
	anon.POST("/auth/forgot", app.submitForgotPassword)
	anon.GET("/auth/reset", app.resetPassword)
	anon.POST("/auth/reset", app.submitResetPassword)

	// Authenticated routes - must be logged in with a valid session
	auth := r.Group("/")
	auth.Use(SessionMiddleware())

	auth.GET("/auth/logout", app.logout)

	// User profiles
	auth.GET("/profile", app.redirectProfile)
	auth.GET("/profile/:id", app.showProfile)
	auth.GET("/profile/edit/:id", app.editProfile)

	auth.GET("/tournaments", app.listTournaments)
	auth.GET("/tournaments/:id", app.showTournament)
	auth.GET("/tournaments/edit/:id", app.editTournament)
	auth.GET("/tournaments/create", app.createTournament)
	auth.POST("/tournaments/create", app.saveTournament)
	// auth.GET("/tournaments/add", app.addTournament)

	auth.GET("/clubs", app.listClubs)
	auth.GET("/clubs/:id", app.showClub)
	auth.GET("/api/clubs/filtered-list", app.filterClubs)

	auth.GET("/ratings", app.listRatings)
	auth.GET("/api/ratings/filtered-list", app.filterRatings)

	auth.GET("/api/organisers/filtered-list", app.filterOrganisers)

	// auth.GET("/admin", app.getAdminDashboard)
	// auth.GET("/admin/tournaments", app.getAdminTournamentsList)

	// Match all non-existent routes
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{})
	})

	// Use HTML template helpers
	r.SetFuncMap(setTemplateHelpers())

	// Load all html templates and partials
	r.LoadHTMLGlob("templates/**/*.html")

	// Listen and Serve
	if err := r.Run(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)); err != nil {
		log.Fatal(err)
	}
}

func postgresConnection(c config.Config) *sql.DB {
	// Set up the dsn connection string
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.DBUser,
		c.DBPass,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)

	// Connect to the database
	conn, err := sql.Open(c.DBDriver, dsn)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	return conn
}

// SessionMiddleware checks for an existing session and saves it to the request context. If none, redirects to login page
func SessionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request provides a session cookie and has an existing session
		session, err := store.Get(c.Request, SessionID)
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
		}

		// If session is new redirect to the login page and do not save the session
		if session.IsNew {
			c.Redirect(http.StatusFound, "/auth/login")
			c.Abort()
			return
		}

		// Save existing session to the context
		c.Set("session", session)

		// Continue with the next handler
		c.Next()
	}
}

func UnauthenticatedMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if the request provides a session cookie and has an existing session
		// Send to profile if they do, otherwise continue with the next handler func
		session, err := store.Get(c.Request, SessionID)
		if err != nil {
			// We might arrive here if the session key has been changed. Can we just expire the cookie and redirect to login?
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		// If session is not new then one already exists so redirect to the profile page
		if !session.IsNew {
			c.Redirect(http.StatusFound, "/profile")
			c.Abort()
			return
		}

		// Continue with the next handler
		c.Next()
	}
}
