-- name: CreateOrUpdateGuildConfig :one
INSERT INTO guild_config (id, json)
VALUES ((SELECT id
         FROM guild
         WHERE discord_id = $1), $2)
ON CONFLICT (id) DO UPDATE
    SET json = $2
RETURNING *;

-- name: TryCreateGuildConfig :one
INSERT INTO guild_config (id, json)
VALUES ((SELECT id
         FROM guild
         WHERE discord_id = $1), $2)
ON CONFLICT (id) DO NOTHING
RETURNING *;

-- name: GetGuildConfig :one
SELECT c.*
FROM guild g
         JOIN guild_config c on g.id = c.id
WHERE g.discord_id = $1;

-- name: UpdateGuildConfig :exec
UPDATE guild_config c
SET json = $1
FROM guild g
WHERE c.id = g.id
  AND g.discord_id = $2;

-- name: DeleteGuildConfig :exec
DELETE
FROM guild_config c USING guild g
WHERE c.id = g.id
  AND g.discord_id = $1;