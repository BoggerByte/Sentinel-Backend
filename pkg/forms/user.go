package forms

type GetUserURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}
