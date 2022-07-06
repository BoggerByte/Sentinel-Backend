package forms

import (
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
)

// RequireDiscordIDRequest was created to somehow get guild discord id
type RequireDiscordIDRequest struct {
	DiscordID int64 `uri:"discord_id" json:"discord_id" forms:"discord_id" binding:"required,min=1"`
}

type OverwriteGuildConfigJSON struct {
	Permissions objects.GuildConfigPermissions `json:"permissions" binding:"required"`
	Data        objects.GuildConfigData        `json:"data" binding:"required"`
}
