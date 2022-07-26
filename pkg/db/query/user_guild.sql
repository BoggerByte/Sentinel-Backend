-- name: CreateOrUpdateUserGuildRel :one
INSERT INTO user_guild (guild_discord_id, account_discord_id, permissions)
VALUES ($1, $2, $3)
ON CONFLICT (guild_discord_id, account_discord_id) DO UPDATE
    SET permissions = $3
RETURNING *;

-- name: CreateUserGuildRel :one
INSERT INTO user_guild (guild_discord_id, account_discord_id, permissions)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserGuildRel :one
SELECT *
FROM user_guild
WHERE account_discord_id = $1
  AND guild_discord_id = $2
LIMIT 1;