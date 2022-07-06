package forms

type GetGuildURI struct {
	DiscordID int64 `uri:"discord_id" binding:"required,min=1"`
}

type GetAllAccountGuildsURI struct {
	DiscordID int64 `uri:"discord_id" binding:"required,min=1"`
}
