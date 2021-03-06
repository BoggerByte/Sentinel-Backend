package forms

import (
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
)

// RequireDiscordIDRequest was created to somehow get guild discord id
type RequireDiscordIDRequest struct {
	DiscordID string `uri:"discord_id" json:"discord_id" forms:"discord_id" binding:"required"`
}

type GetGuildConfigPresetURI struct {
	Preset string `uri:"preset" binding:"required,oneof=default"`
}

type OverwriteGuildConfigJSON struct {
	Permissions objects.GuildConfigPermissions `json:"permissions" binding:"required"`
	Data        objects.GuildConfigData        `json:"data" binding:"required"`
	Preset      string                         `json:"preset" binding:"required,oneof=default custom"`
}
