-- name: CreateOrUpdateUser :one
INSERT INTO "user" (discord_id, username, discriminator, verified, email,
                    avatar, banner, accent_color)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
ON CONFLICT (discord_id) DO UPDATE
    SET username      = $2,
        discriminator = $3,
        verified      = $4,
        email         = $5,
        avatar        = $6,
        banner        = $8,
        accent_color  = $9
RETURNING *;

-- name: GetUser :one
SELECT *
FROM "user"
WHERE discord_id = $1
LIMIT 1;