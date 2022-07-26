package forms

import db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"

type GetGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type GetUserGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type CreateOrUpdateGuildsJSON struct {
	Guilds []db.CreateOrUpdateGuildParams
}

type CreateOrUpdateGuildJSON struct {
	db.CreateOrUpdateGuildParams
}
