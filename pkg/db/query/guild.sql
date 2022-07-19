-- name: CreateOrUpdateGuild :one
INSERT INTO guild (discord_id, name, icon, owner_discord_id)
VALUES ($1, $2, $3, $4)
ON CONFLICT (discord_id) DO UPDATE
    SET name=$2,
        icon=$3,
        owner_discord_id=$4
RETURNING *;

-- name: GetGuild :one
SELECT *
FROM guild
WHERE discord_id = $1
LIMIT 1;

-- name: DeleteGuild :exec
DELETE
FROM guild
WHERE discord_id = $1;

-- name: GetUserGuilds :many
SELECT coalesce(g.id, 0),
       ug.guild_discord_id AS discord_id,
       ug.permissions,
       coalesce(g.owner_discord_id, '0'),
       coalesce(g.icon, '#'),
       coalesce(g.name, '')
FROM user_guild ug
         LEFT OUTER JOIN guild g ON g.discord_id = ug.guild_discord_id
WHERE ug.account_discord_id = $1;