package forms

type GetGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type GetUserGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}
