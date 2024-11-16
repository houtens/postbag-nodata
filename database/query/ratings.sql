-- #160 is_locked

-- name: CreateRating :one
INSERT INTO ratings (
    user_id,
    tournament_id,
    division,
    num_games,
    start_rating,
    end_rating,
    rating_points,
    opp_ratings_sum,
    num_wins,
    is_locked,
    x_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
) RETURNING *;

-- name: GetRating :one
SELECT * FROM ratings WHERE id = $1 LIMIT 1;

-- name: ListRatings :many
SELECT * FROM ratings
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateRating :one
UPDATE ratings
SET
user_id = $2,
tournament_id = $3,
division = $4,
num_games = $5,
start_rating = $6,
end_rating = $7,
rating_points = $8,
opp_ratings_sum = $9,
num_wins = $10,
is_locked = $11,
x_id = $12
WHERE id = $1
RETURNING *;

-- name: DeleteRating :exec
DELETE FROM ratings WHERE id = $1;

-- end CRUD

-- name: GetRatingByXID :one
SELECT * FROM ratings WHERE x_id = $1 LIMIT 1;

-- name: GetLatestRating :one
SELECT end_rating FROM ratings r join tournaments t on r.tournament_id = t.id 
WHERE r.user_id = $1 AND t.is_locked = true ORDER BY t.end_date DESC LIMIT 1;

-- name: ListRatingsTournaments :many
SELECT * FROM
    ratings r JOIN tournaments t
ON r.tournament_id = t.id
ORDER BY r.id
LIMIT $1
OFFSET $2;

-- name: ListRatingsByUser :many
SELECT * FROM 
    ratings r JOIN tournaments t 
ON r.tournament_id = t.id
WHERE r.user_id = $1 
AND r.is_locked = 't'
ORDER BY t.end_date;

-- name: ListRatingsByUserRev :many
SELECT * FROM 
    ratings r JOIN tournaments t 
ON r.tournament_id = t.id
WHERE r.user_id = $1 
AND r.is_locked = 't'
ORDER BY t.end_date desc;

-- name: ListTournamentRatings :many
SELECT * FROM ratings r JOIN users u
ON r.user_id = u.id
WHERE r.tournament_id = $1 ORDER BY start_rating DESC;


-- name: ListRecentRatingsForPlayer :many
SELECT t.id, t.name, t.start_date, r.num_games, r.start_rating, r.end_rating, r.num_wins
FROM ratings r JOIN tournaments t ON r.tournament_id = t.id 
WHERE t.is_locked = true 
AND r.is_locked 
AND user_id = $1 
ORDER BY t.end_date DESC 
LIMIT 10;


-- name: ListRatingsRank :many
SELECT row_number() over(ORDER BY end_rating DESC), first_name, last_name, titles.name as title, clubs.name as club, max_ratings.user_id, end_rating, end_date, countries.name
FROM (
  SELECT DISTINCT ON (r.user_id) user_id, r.end_rating, t.end_date
  FROM ratings r JOIN tournaments t ON r.tournament_id = t.id
  WHERE (t.is_locked = true and r.is_locked = true)
  ORDER BY user_id, t.end_date DESC
) max_ratings JOIN users ON max_ratings.user_id = users.id
JOIN titles on users.title_id = titles.id
LEFT OUTER JOIN contacts on users.id = contacts.user_id
LEFT OUTER JOIN countries on contacts.country_id = countries.id
LEFT OUTER JOIN clubs on users.club_id = clubs.id
WHERE end_date > '2022-01-01' ORDER BY end_rating DESC 
LIMIT $1
OFFSET $2;

-- name: ListRatingsForClub :many
SELECT u.id, u.first_name, u.last_name, tt.name as title, mr.user_id, mr.end_rating
FROM (
  SELECT DISTINCT ON (r.user_id) user_id, r.end_rating, t.end_date 
  FROM ratings r JOIN tournaments t ON r.tournament_id = t.id 
  WHERE (t.is_locked = true and r.is_locked = true) 
  ORDER BY user_id, t.end_date DESC
) mr JOIN users u ON mr.user_id = u.id 
JOIN titles tt on tt.id = u.title_id
WHERE club_id = $1
ORDER BY end_rating DESC
LIMIT $2
OFFSET $3;


-- name: GetRatingsTournamentMetadata :one
select max(coalesce(division, 0)) divisions, max(coalesce(num_games, 0)) rounds from ratings where tournament_id = $1 and is_locked = 't' limit 1;

-- name: GetCountPlayersInTournament :one
select count(distinct user_id) from ratings where tournament_id = $1 limit 1;

-- name: UpdateRatingsNumWins :one
update ratings set num_wins = $3 where user_id = $1 and tournament_id = $2 returning *;

-- name: TruncateRatings :exec
truncate ratings cascade;

