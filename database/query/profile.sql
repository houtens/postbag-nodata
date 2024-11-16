-- name: GetProfile :one
SELECT
  u.id,
  u.first_name,
  u.last_name,
  u.alt_name,
  u.absp_num,
  u.avatar,
  c.name as club_name,
  c.id as club_id,
  c.county,
  t.name as title_name,
  r.name as role_name
FROM users u
LEFT JOIN clubs c ON u.club_id = c.id 
LEFT JOIN titles t ON u.title_id = t.id 
LEFT JOIN auth_roles r ON u.role_id = r.id
WHERE u.id = $1
LIMIT 1;
