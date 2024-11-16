-- #160 is_locked

-- name: CreateResult :one
INSERT INTO results (
    player1_id,
    player2_id,
    score1,
    score2,
    spread,
    tournament_id,
    type,
    round_num,
    is_locked,
    x_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetResult :one
SELECT * FROM results
WHERE id = $1 LIMIT 1;

-- name: ListResults :many
SELECT * FROM results
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateResult :one
UPDATE results
SET
player1_id = $2,
player2_id = $3,
score1 = $4,
score2 = $5,
spread = $6,
tournament_id = $7,
type = $8,
round_num = $9,
is_locked = $10,
x_id = $11
WHERE id = $1
RETURNING *;

-- name: DeleteResult :exec
DELETE FROM results
where id = $1;

-- name: GetResultByXID :one
SELECT * FROM results
WHERE x_id = $1 LIMIT 1;

-- name: GetWinLossDraw :one
SELECT
  SUM(CASE WHEN r.score1 > r.score2 THEN 1 ELSE 0 END) as won,
  SUM(CASE WHEN r.score1 < r.score2 THEN 1 ELSE 0 END) as lost,
  SUM(CASE WHEN r.score1 = r.score2 THEN 1 ELSE 0 END) as drawn
FROM (
    SELECT player1_id, score1, score2, type FROM results WHERE results.player1_id = $1 AND is_locked = true
    UNION ALL
    SELECT player2_id, score2, score1, type FROM results WHERE player2_id = $1 AND is_locked = true
) AS r LIMIT 1;

-- name: GetWinsByUserTournament :one
SELECT SUM(CASE WHEN score1 > score2 THEN 1 WHEN score1 = score2 THEN 0.5 ELSE 0 END)::float wins FROM (
    SELECT score1, score2 FROM results WHERE results.tournament_id = $1 AND results.player1_id = $2
    UNION ALL
    SELECT score2, score1 FROM results WHERE tournament_id = $1 AND player2_id = $2) AS x LIMIT 1;

-- name: GetAverageScores :one
SELECT ROUND(AVG(s.score1)) AS for, ROUND(AVG(s.score2)) AS against FROM
(SELECT score1, score2
  FROM results
  WHERE results.player1_id = $1 AND type = 1 AND is_locked = true
  UNION ALL
  SELECT score2, score1 FROM results
  WHERE player2_id = $1 AND type = 1 AND is_locked = true
) AS s;

-- name: TruncateResults :exec
truncate results cascade;
