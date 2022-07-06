package forms

type GetUserURI struct {
	DiscordID int64 `uri:"discord_id" binding:"required,min=1"`
}
