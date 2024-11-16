package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

// MemberCSV defines struct for parsing member csv data
type MemberCSV struct {
	XID          string `csv:"id"`
	ABSPNumber   string `csv:"absp_number"`
	AccessLevel  string `csv:"access_level"`
	FirstName    string `csv:"first_name"`
	LastName     string `csv:"last_name"`
	Status       string `csv:"status"`
	Address1     string `csv:"address1"`
	Address2     string `csv:"address2"`
	Address3     string `csv:"address3"`
	Address4     string `csv:"address4"`
	Postcode     string `csv:"postcode"`
	Phone        string `csv:"phone"`
	Email        string `csv:"email"`
	Club         string `csv:"club"`
	ClubID       string `csv:"club_id"`
	Country      string `csv:"country"`
	Year         string `csv:"year"`
	PaidWhen     string `csv:"paid_when"`
	PaidHow      string `csv:"paid_how"`
	LastPayment  string `csv:"last_payment"`
	Mobile       string `csv:"mobile"`
	Notes        string `csv:"notes"`
	Life         string `csv:"life"`
	OBPost       string `csv:"ob_post"`
	OBPDF        string `csv:"ob_pdf"`
	DOB          string `csv:"dob"`
	LastPlayed   string `csv:"last_played"`
	Rating       string `csv:"rating"`
	TotalGames   string `csv:"total_games"`
	Updated      string `csv:"updated"`
	PasswordHash string `csv:"password_hash"`
	FBID         string `csv:"fb_id"`
}

var DeceasedMemberText string = `deceased|Deceased|dec'd`

func setAuthRole(m MemberCSV, r map[string]uuid.UUID) uuid.UUID {
	var authRole uuid.UUID

	switch m.AccessLevel {
	case "0":
		// Non-members
		authRole = r["No Login"]
	case "1":
		// Ordinary member
		authRole = r["Member"]
	case "2":
		// Tournament organiser is no longer needed - make ordinary member
		authRole = r["Member"]
	case "4":
		// Unknown - Rick Blakeway?
		authRole = r["Member"]
	case "5":
		// Members admin - Nuala, Natalie
		authRole = r["Members Admin"]
	case "6":
		// Unknown - Peter Ashurst?
		authRole = r["Member"]
	case "10":
		// Super Admin - Stewart, Wayne, Ross -- also Brodie, Whiteoak, Rachel
		authRole = r["Super Admin"]
	case "11":
		// Non-user roles - uk scrabble forum and Ratings officer (for email)
		authRole = r["No Login"]
	default:
		// Needs checking - they all appear to be non-members
		authRole = r["No Login"]
	}

	return authRole
}

func cleanEmail(e string) string {
	e = strings.TrimSpace(e)

	// fix known typos
	if e == "stanley1967@talktalknet" {
		e = "stanley1967@talktalk.net"
	}
	if e == "pcb17673gmail.com" {
		e = "pcb17673@gmail.com"
	}

	e = strings.ToLower(e)

	return e
}

func setEmail(e string) sql.NullString {
	// trim whitespace and clean-up some typos
	e = cleanEmail(e)

	// Validate email address
	ok := govalidator.IsEmail(e)

	// Failures
	if !ok {
		return sql.NullString{Valid: false}
	}

	// Convert to lowercase before returning
	// e = strings.ToLower(e)
	return sql.NullString{String: e, Valid: true}
}

func setMobile(m string) string {
	rx := regexp.MustCompile(`[^0-9 -+]+`)
	failed := rx.MatchString(m)
	if failed {
		return ""
	}
	return m
}

// setNotes purges deceased notes and \N strings, then removes surrounding whitespace
func setNotes(n string) string {
	match, err := regexp.MatchString(`Deceased|dec'd|\\N`, n)
	if err != nil {
		log.Fatal(err)
	}
	// If pattern is found, return an empty string
	if match {
		return ""
	}
	return strings.TrimSpace(n)
}

func setPasswordHash(pw string) sql.NullString {
	// Hash must be exactly 60 characters long
	if len(pw) != 60 {
		return sql.NullString{Valid: false}
	}

	// Bcrypt with 10 rounds begins with $2y$10$
	match, err := regexp.MatchString(`^\$2y\$10\$`, pw)
	if err != nil {
		log.Fatal(err)
	}
	if !match {
		return sql.NullString{Valid: false}
	}

	// Return a bcrypt hashed password
	return sql.NullString{String: pw, Valid: true}
}

func setAbspNumber(s string) sql.NullInt32 {
	value, err := strconv.Atoi(s)
	if err != nil {
		fmt.Println("unable to parse absp_number")
		log.Fatal(err)
	}

	if value <= 0 {
		return sql.NullInt32{Valid: false}
	}

	return sql.NullInt32{Int32: int32(value), Valid: true}
}

func setXLife(x string) bool {
	return x == "1"
}

func setXPost(x string) bool {
	return x == "1"
}

func importMembers() {
	filename := "seed/data/export/members.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	rows := []MemberCSV{}
	if err := gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal(err)
	}

	// Read titles and clubs into maps
	titlesMap := readTitles()
	clubsMap := readClubs()
	authRolesMap := readAuthRoles()

	fmt.Print("Import members/users... ")

	now := time.Now()

	for _, r := range rows {
		// Extract Alt names such as "Paul{Tranmere} Thomson"
		firstName, altName := splitAmbiguousNames(r.FirstName)

		// Parse title from CSV status fields
		titleID, ok := titlesMap[r.Status]
		if !ok {
			// These should all exist -- if not, look at clubsMap below
			log.Fatal("unable to lookup title")
		}

		// Parse clubs from CSV via XID
		clubID, ok := clubsMap[r.ClubID]
		if !ok {
			clubID = uuid.NullUUID{Valid: false}
		}

		// Parse and import email and passwordhash
		email := setEmail(r.Email)
		passwordHash := setPasswordHash(r.PasswordHash)

		// Parse and import absp numbers
		if r.ABSPNumber == "8" {
			continue
		}
		abspNum := setAbspNumber(r.ABSPNumber)
		xLife := setXLife(r.Life)
		xPost := setXPost(r.OBPost)

		// Determine the appropriate AuthRole based on the old access_level
		authRole := setAuthRole(r, authRolesMap)

		// If we have a country xid then try to match it to an xid in the country table
		var country models.Country
		if r.Country != "\\N" && r.Country != "0" {
			// fmt.Printf("%v %v %v\n", r.FirstName, r.LastName, r.Country)
			country, err = db.GetCountryByXID(context.Background(), sql.NullString{String: r.Country, Valid: true})
			if err != nil {
				fmt.Printf("failed to match country for %v %v, %v\n", r.FirstName, r.LastName, r.Country)
				log.Fatal(err)
			}
		}

		arg := models.CreateUserParams{
			FirstName:    firstName,
			LastName:     r.LastName,
			AltName:      sql.NullString{String: altName, Valid: true},
			Email:        email,
			PasswordHash: passwordHash,
			AbspNum:      abspNum,
			TitleID:      uuid.NullUUID{UUID: titleID, Valid: true},
			ClubID:       clubID,
			RoleID:       authRole,
			XPost:        xPost,
			XLife:        xLife,
			XID:          r.XID,
		}

		user, err := db.CreateUser(context.Background(), arg)
		if err != nil {
			fmt.Println(r)
			log.Fatal(err)
		}

		// Save contact details
		address1 := purgeNull(r.Address1)
		address2 := purgeNull(r.Address2)
		address3 := purgeNull(r.Address3)
		address4 := purgeNull(r.Address4)
		postcode := purgeNull(r.Postcode)
		phone := purgeNull(r.Phone)

		mobile := setMobile(r.Mobile)
		notes := setNotes(r.Notes)

		arg1 := models.CreateContactParams{
			UserID:    user.ID,
			Address1:  sql.NullString{String: address1, Valid: true},
			Address2:  sql.NullString{String: address2, Valid: true},
			Address3:  sql.NullString{String: address3, Valid: true},
			Address4:  sql.NullString{String: address4, Valid: true},
			Postcode:  sql.NullString{String: postcode, Valid: true},
			CountryID: uuid.NullUUID{UUID: country.ID, Valid: true},
			Phone:     sql.NullString{String: phone, Valid: true},
			Mobile:    sql.NullString{String: mobile, Valid: true},
			Notes:     sql.NullString{String: notes, Valid: true},
		}

		// Create contact details
		_, err = db.CreateContact(context.Background(), arg1)
		if err != nil {
			log.Fatal(err)
		}
	}

	// report how many were added
	arg := models.ListUsersParams{
		Limit:  10000,
		Offset: 0,
	}
	members, err := db.ListUsers(context.Background(), arg)
	if err != nil {
		log.Fatal()
	}
	fmt.Printf("%d found, ", len(members))
	fmt.Println(time.Since(now))
}

func truncateUsers() {
	db.TruncateUsers(context.Background())
}

func truncateContacts() {
	db.TruncateContacts(context.Background())
}
