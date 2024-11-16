-- name: CreateTournamentState :one
INSERT INTO tournament_state (
    name,
    code
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetTournamentState :one
SELECT * FROM tournament_state WHERE id = $1 LIMIT 1;

-- name: ListTournamentStates :many
SELECT * FROM tournament_state ORDER BY id LIMIT $1 OFFSET $2;

-- name: UpdateTournamentState :one
UPDATE tournament_state 
SET 
    name = $2,
    code = $3
WHERE id = $1 RETURNING *;

-- name: DeleteTournamentState :exec
DELETE FROM tournament_state WHERE id = $1;

-- name: TruncateTournamentState :exec
truncate tournament_state cascade;

