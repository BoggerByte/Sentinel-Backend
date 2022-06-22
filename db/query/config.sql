-- name: CreateConfig :one
INSERT INTO config (id, version, filename, created_at, guild_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetConfig :one
SELECT *
FROM config
WHERE id = $1
LIMIT 1;

-- name: ListConfigs :many
SELECT *
FROM config
ORDER BY id
OFFSET $1 LIMIT $2;

-- name: DeleteConfig :exec
DELETE
FROM config
WHERE id = $1;