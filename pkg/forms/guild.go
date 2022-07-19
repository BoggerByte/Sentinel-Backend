package forms

type GetGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type GetAllAccountGuildsURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}
