package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"
	"github.com/google/uuid"
	"github.com/houtens/postbag/internal/models"
)

var (
	ErrFailedCreateResult = errors.New("failed to create result in database")
)

type ResultCSV struct {
	ID         string `csv:"id"`
	Player1    string `csv:"player1"`
	Player2    string `csv:"player2"`
	Score1     string `csv:"score1"`
	Score2     string `csv:"score2"`
	Spread     string `csv:"spread"`
	Type       string `csv:"type"`
	Tournament string `csv:"tournament"`
	RoundNum   string `csv:"round"`
	Locked     string `csv:"locked"`
}

func parsePlayers(player1, player2 string) (uuid.NullUUID, uuid.NullUUID) {
	p1 := userByXID(player1)
	p2 := userByXID(player2)
	return p1, p2
}

func parseScores(score1, score2 string) (int32, int32) {
	s1, err := strconv.Atoi(score1)
	if err != nil {
		log.Fatal("could not convert score1 to int: ", err)
	}
	s2, err := strconv.Atoi(score2)
	if err != nil {
		log.Fatal("could not convert score2 to int: ", err)
	}
	return int32(s1), int32(s2)
}

func parseSpread(spread string) int32 {
	sp, err := strconv.Atoi(spread)
	if err != nil {
		log.Fatal("could not convert spread to int: ", spread)
	}
	return int32(sp)
}

func parseTournament(tournament string) uuid.NullUUID {
	return tournamentByXID(tournament)
}

func parseRound(round string) int32 {
	rd, err := strconv.Atoi(round)
	if err != nil {
		log.Fatal("unable to parse round: ", err)
	}
	return int32(rd)
}

func parseLocked(locked string) bool {
	isLocked, err := strconv.ParseBool(locked)
	if err != nil {
		log.Fatal("could not parse result locked to boolean: ", err)
	}
	return isLocked
}

func parseType(typ string) int32 {
	t, err := strconv.Atoi(typ)
	if err != nil {
		log.Fatal("could not parse result type: ", err)
	}
	return int32(t)
}

func importResults() {
	filename := "seed/data/export/results.csv"

	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("could not open csv file %s\n", err)
	}
	defer f.Close()

	rows := []ResultCSV{}
	if err = gocsv.UnmarshalFile(f, &rows); err != nil {
		log.Fatal("unable to unmarshal csv:", err)
	}

	fmt.Print("Importing results... ")
	now := time.Now()
	for _, r := range rows {
		// NOTE: types of results can be re-mapped
		// type 0 - either p1 or p2 is missing, no scores and no spread
		// type 1 - normal results with both players, score and spread
		// type 2 - players defined, no scores, but spread
		// 		could subtract the spread from player -- 0 -225
		// type 3 - players defined, scores and spread represent games and total games
		// 		Split these so we record a win for each game
		// type 4 - player1 or player 2 missing, walkover with 100-0 score and spread
		// ignore, no - needed for manual results

		p1, p2 := parsePlayers(r.Player1, r.Player2)
		s1, s2 := parseScores(r.Score1, r.Score2)

		sp := parseSpread(r.Spread)
		tn := parseTournament(r.Tournament)

		// Skip tournaments that do not exist
		if !tn.Valid {
			continue
		}

		rnd := parseRound(r.RoundNum)
		lck := parseLocked(r.Locked)
		typ := parseType(r.Type)

		switch r.Type {
		case "1":
			processTypeOne(p1, p2, s1, s2, sp, tn.UUID, typ, rnd, lck, r.ID)

		case "2":
			processTypeTwo(p1, p2, s1, s2, sp, tn.UUID, typ, rnd, lck, r.ID)

		case "3":
			processTypeThree(p1, p2, s1, s2, sp, tn.UUID, typ, rnd, lck, r.ID)

		case "4":
			break
			// byes with a missing player
			// arg, err := createResultParams(p1, p2, s1, s2, sp, tn.UUID, typ, rnd, lck, r.ID)
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// err = createResult(arg)
			// if err != nil {
			// 	log.Fatal(err)
			// }
		}
	}
	fmt.Println(time.Since(now))
}

func processTypeThree(p1 uuid.NullUUID, p2 uuid.NullUUID, s1 int32, s2 int32, sp int32, tid uuid.UUID, typ int32, rnd int32, lck bool, xid string) {
	// NOTE:Type 3 results record only games won and no scores or spread.
	// s1 and s2 are the number of games won by each player; sp is the sum of games

	// fmt.Println(s1, s2, sp, rnd, typ)

	// single game
	if sp == 1 {
		if s1 > s2 {
			sp = 1
		} else {
			sp = -1
		}
		// create results params and save result
		arg, err := createResultParams(p1, p2, s1, s2, sp, tid, int32(5), rnd, lck, xid)
		if err != nil {
			log.Fatal(err)
		}
		err = createResult(arg)
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	// compound games
	if sp > 1 {
		// buffer the scores as wins since they aren't scores
		var w1, w2 int32
		w1, w2 = s1, s2
		for x := 0; x < int(w1+w2); x++ {
			s1, s2, sp = 0, 1, -1
			if x < int(w1) {
				s1, s2, sp = 1, 0, 1
			}
			arg, err := createResultParams(p1, p2, s1, s2, sp, tid, int32(6), int32(x), lck, xid)
			if err != nil {
				log.Fatal(err)
			}
			err = createResult(arg)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func processTypeTwo(p1 uuid.NullUUID, p2 uuid.NullUUID, s1 int32, s2 int32, sp int32, tid uuid.UUID, typ int32, rnd int32, lck bool, xid string) {
	switch {
	case sp > 0:
		s1 = sp
	case sp < 0:
		s2 = -1 * sp
	default:
		break
	}

	// create results params and save result
	arg, err := createResultParams(p1, p2, s1, s2, sp, tid, typ, rnd, lck, xid)
	if err != nil {
		log.Fatal(err)
	}
	err = createResult(arg)
	if err != nil {
		log.Fatal(err)
	}
}

func processTypeOne(p1 uuid.NullUUID, p2 uuid.NullUUID, s1 int32, s2 int32, sp int32, tid uuid.UUID, typ int32, rnd int32, lck bool, xid string) {
	spread := s1 - s2
	if spread != sp {
		log.Fatal("spread and scores do not match")
	}

	// create results params and save result
	arg, err := createResultParams(p1, p2, s1, s2, sp, tid, typ, rnd, lck, xid)
	if err != nil {
		log.Fatal(err)
	}
	err = createResult(arg)
	if err != nil {
		log.Fatal(err)
	}
}

func createResultParams(
	p1 uuid.NullUUID,
	p2 uuid.NullUUID,
	s1 int32,
	s2 int32,
	sp int32,
	tid uuid.UUID,
	typ int32,
	rnd int32,
	lck bool,
	xid string,
) (models.CreateResultParams, error) {
	return models.CreateResultParams{
		Player1ID:    p1,
		Player2ID:    p2,
		Score1:       s1,
		Score2:       s2,
		Spread:       sp,
		TournamentID: tid,
		Type:         typ,
		RoundNum:     rnd,
		IsLocked:     lck,
		XID:          xid,
	}, nil
}

func createResult(mr models.CreateResultParams) error {
	_, err := db.CreateResult(context.Background(), mr)
	if err != nil {
		return ErrFailedCreateResult
	}
	return nil
}

func truncateResults() {
	db.TruncateResults(context.Background())
}
