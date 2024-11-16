package service

import (
	"fmt"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(storedHash string, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(providedPassword))
	return err == nil
}

// fullName returns a full name string with title if one exists
func fullName(firstName string, lastName string, title string) string {
	name := fmt.Sprintf("%s %s", firstName, lastName)
	if len(title) > 0 {
		name = fmt.Sprintf("%s (%s)", name, title)
	}
	return name
}

func calcWinRate(won, lost, drawn int64) string {
	var played = won + lost + drawn
	if played == 0 {
		return "0"
	}

	wins := float64(won) + 0.5*float64(drawn)
	return fmt.Sprintf("%.1f", 100*wins/float64(played))
}
