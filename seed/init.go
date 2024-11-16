package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/houtens/postbag/config"
	"github.com/houtens/postbag/internal/models"
	_ "github.com/lib/pq"
)

func init() {

	cfg := config.Load()

	// Set database connection string
	dsn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	conn, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		log.Fatal("cannot connect to the database:", err)
	}

	db = models.New(conn)
}
