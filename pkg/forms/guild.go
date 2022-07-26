package forms

type GetGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type GetUserGuildURI struct {
	DiscordID string `uri:"discord_id" binding:"required"`
}

type CreateOrUpdateGuildJSON struct {
	DiscordID      string `json:"discord_id" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Icon           string `json:"icon" binding:"required"`
	OwnerDiscordID string `json:"owner_discord_id" binding:"required"`
}

type CreateOrUpdateGuildsJSON struct {
	Guilds []CreateOrUpdateGuildJSON `json:"guilds" binding:"required"`
}
