package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

const (
	baseURL     = "https://absp-database.org/pix/"
	downloadDir = "../../seed/data/pix"
	avatarDir   = "../../static/images/avatars"
)

var db *sql.DB

// getImage makes a get request to the image url and returns a response
func getImage(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch url: %w", err)
	}

	// We only care about 200s
	if res.StatusCode != 200 {
		return nil, errors.New("return a non-200 response")
	}

	return res, nil
}

// saveImage takes an http response adn saves the body to file
func saveImage(rs *http.Response, id int) error {
	body, err := io.ReadAll(rs.Body)
	if err != nil {
		return err
	}
	output := fmt.Sprintf("%s/%d.jpg", downloadDir, id)
	err = os.WriteFile(output, body, 0600)
	if err != nil {
		return err
	}

	return nil
}

// postgresConnect establishes a database connection
func postgresConnect() error {
	var err error
	dsn := "postgres://postgres:@localhost/postbag?sslmode=disable"
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("unable to connect to postgres: %w", err)
	}

	return nil
}

// processFetch fetches an images and saves them to disk
func processFetch() error {
	// ensure the download directory exists
	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		return fmt.Errorf("cannot find output directory: %w", err)
	}

	// establish a database connection
	err := postgresConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Read x_id of all users from database
	query := "select x_id from users"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("could not read users: %w", err)
	}
	defer rows.Close()

	// prepare slice for xids and read them from the database
	xids := make([]int, 0)
	for rows.Next() {
		var xid int
		err = rows.Scan(&xid)
		if err != nil {
			return fmt.Errorf("unable to scan row: %w", err)
		}
		xids = append(xids, xid)
	}

	var url string
	// try to fetch each image and save it
	for _, id := range xids {
		// convert int to request path to the image - don't bother with strconv
		url = fmt.Sprintf("%s%d.jpg", baseURL, id)
		fmt.Println(url)

		// fetch image, or skip to the next
		res, err := getImage(url)
		if err != nil {
			log.Println(err)
			continue
		}

		// save the image to the downlad directory
		err = saveImage(res, id)
		if err != nil {
			return fmt.Errorf("failed to save image: %w", err)
		}
	}

	return nil
}

// processSeed iterates over avatar images and saves them to the user model
func processSeed() error {
	// check the path to the downloads directory exists
	if _, err := os.Stat(downloadDir); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("cannot find output directory: %w", err)
	}

	// establish a database connection
	err := postgresConnect()
	if err != nil {
		return err
	}
	defer db.Close()

	// Read x_id of all users from the database
	query := "select x_id from users"
	rows, err := db.Query(query)
	if err != nil {
		return fmt.Errorf("could not read users: %w", err)
	}
	defer rows.Close()

	// Ensure the static/images directory exists
	if _, err := os.Stat(avatarDir); os.IsNotExist(err) {
		err := os.Mkdir(avatarDir, 0755)
		if err != nil {
			return fmt.Errorf("could not create destination directory: %w", err)
		}
	}

	// Iterate over xids
	for rows.Next() {
		var xid int
		err = rows.Scan(&xid)
		if err != nil {
			return fmt.Errorf("unable to scan row: %w", err)
		}

		// find all downloaded pix as source
		source := fmt.Sprintf("%s/%d.jpg", downloadDir, xid)
		if _, err := os.Stat(downloadDir); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("cannot find output directory: %w", err)
		}

		// open the source - if not found then continue
		src, err := os.Open(source)
		if err != nil {
			continue
		}
		defer src.Close()

		// Generate a new uuid
		avatarID := uuid.New()

		destination := fmt.Sprintf("%s/%s.jpg", avatarDir, avatarID)
		dst, err := os.Create(destination)
		if err != nil {
			return fmt.Errorf("cannot open the destination file: %w", err)
		}
		defer dst.Close()

		// Saving the link and moving the file should be atomic
		_, err = io.Copy(dst, src)
		if err != nil {
			return fmt.Errorf("failed to copy: %w", err)
		}

		// Save link to the database
		avatar := fmt.Sprintf("%s.jpg", avatarID.String())
		fmt.Println(avatar)
		query := `UPDATE users SET avatar = $1 where x_id = $2`
		// fmt.Println(query, avatar, xid)
		_, err = db.Exec(query, avatar, xid)
		if err != nil {
			return fmt.Errorf("failed to save avatar to the database: %w", err)
		}

	}

	return nil
}

func main() {
	var fetch = flag.Bool("fetch", false, "fetch all profile pix from the absp database")
	var seed = flag.Bool("seed", false, "seed the database users with profile pix")

	// Parse flags
	flag.Parse()

	var err error

	// Fetch profiles images from the absp-database
	if *fetch {
		err = processFetch()
	}
	if err != nil {
		log.Fatal(err)
	}

	// seed the database with profile pix
	if *seed {
		err = processSeed()
	}
	if err != nil {
		log.Fatal(err)
	}
}
