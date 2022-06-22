-- name: CreateGuild :one
INSERT INTO guild (discord_id, name, icon, owner_id)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetGuild :one
SELECT *
FROM guild
WHERE id = $1
LIMIT 1;

-- name: ListGuilds :many
SELECT *
FROM guild
ORDER BY id
OFFSET $1 LIMIT $2;

-- name: DeleteGuild :exec
DELETE
FROM guild
WHERE id = $1;