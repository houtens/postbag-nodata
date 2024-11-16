package models

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/houtens/postbag/config"
	_ "github.com/lib/pq"
)

var db *Queries

func TestMain(m *testing.M) {

	// Tests require database access so load the config
	cfg := config.Load()

	// postgres connection string for the test db
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s_test?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	conn, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		log.Fatal("cannot connect to the database", err)
	}

	db = New(conn)

	os.Exit(m.Run())
}
