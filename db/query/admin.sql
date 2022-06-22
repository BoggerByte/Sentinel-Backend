-- name: CreateAdmin :one
INSERT INTO admin (discord_id, username, discriminator, verified, email,
                   avatar, flags, banner, accent_color, public_flags)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING *;

-- name: GetAdmin :one
SELECT *
FROM admin
WHERE id = $1
LIMIT 1;

-- name: ListAdmins :many
SELECT *
FROM admin
ORDER BY id
OFFSET $2 LIMIT $1;

-- name: DeleteAdmin :exec
DELETE
FROM admin
WHERE id = $1;