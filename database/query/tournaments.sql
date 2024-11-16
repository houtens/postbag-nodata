-- #160 is_locked

-- name: CreateTournament :one
INSERT INTO tournaments (
    name,
    short_name,
    start_date,
    end_date,
    state,
    num_divisions,
    num_rounds,
    num_entries,
    is_pc,
    is_fc,
    is_rr,
    is_wespa,
    is_invitational,
    is_locked,
    creator_id,
    organiser_id,
    director_id,
    coperator_id,
    x_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15,
    $16, $17, $18, $19
) RETURNING *;

-- name: GetTournament :one
SELECT * FROM tournaments
WHERE id = $1 LIMIT 1;

-- name: ListTournaments :many
SELECT * FROM tournaments
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateTournament :one
UPDATE tournaments
SET
    name = $2, 
    short_name = $3,
    start_date = $4,
    end_date = $5,
    state = $6,
    num_divisions = $7,
    num_rounds = $8,
    num_entries = $9,
    is_pc = $10,
    is_fc = $11,
    is_rr = $12,
    is_wespa = $13,
    is_invitational = $14,
    is_locked = $15,
    creator_id = $16,
    organiser_id = $17,
    director_id = $18,
    coperator_id = $19,
    x_id = $20
WHERE id = $1
RETURNING *;

-- name: DeleteTournament :exec
DELETE FROM tournaments
WHERE id = $1;

-- end crud

-- name: GetTournamentByXID :one
SELECT * FROM tournaments
WHERE x_id = $1 LIMIT 1;

-- name: ListUpcomingTournaments :many
SELECT * FROM tournaments
WHERE is_locked = false
ORDER BY start_date asc;

-- name: ListRecentTournaments :many
SELECT * FROM tournaments
WHERE is_locked = true
ORDER BY end_date desc
LIMIT $1
OFFSET $2;


-- name: ListUpcomingTournamentsForPlayer :many
SELECT DISTINCT(t.name), t.start_date, t.end_date, t.num_rounds
FROM ratings r JOIN tournaments t ON r.tournament_id = t.id
WHERE r.is_locked = false AND r.user_id = $1 ORDER BY t.end_date ASC;


-- name: ListRecentTournamentsForPlayer :many
SELECT DISTINCT(t.name), t.start_date, t.end_date, t.num_rounds
FROM results r JOIN tournaments t ON r.tournament_id = t.id
WHERE t.is_locked = true 
AND (r.player1_id = $1 OR r.player2_id = $1) 
ORDER BY t.end_date DESC 
LIMIT 10;


-- name: UpdateTournamentRoundsDivisions :one
update tournaments set num_rounds = $2, num_divisions = $3 where id = $1 returning *;

-- name: UpdateTournamentEntries :one
update tournaments set num_entries = $2 where id = $1 returning *;

-- name: TruncateTournaments :exec
truncate tournaments cascade;

