package models

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func randomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func randomBool() bool {
	return randomInt(0, 100)%2 == 0
}

func randomID() string {
	r := randomInt(0, 1000000)
	return strconv.FormatInt(r, 10)
}

func randomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func randomCountry() string {
	return randomString(8)
}

func randomFlag() string {
	return randomString(32)
}

func randomCountryCode() string {
	return strings.ToUpper(randomString(3))
}

func randomName() string {
	return randomString(6)
}

func randomTitle() string {
	return strings.ToUpper(randomString(3))
}

func randomDateString() string {
	y := int(randomInt(1970, 2023))
	m := int(randomInt(1, 12))
	d := int(randomInt(1, 28))

	return fmt.Sprintf("%d-%02d-%02d", y, m, d)
}

func randomDate() time.Time {
	date, err := time.Parse("2006-01-02", randomDateString())
	if err != nil {
		log.Fatal("unable to parse random date string")
	}
	return date
}

func randomUUID() uuid.UUID {
	return uuid.New()
}
