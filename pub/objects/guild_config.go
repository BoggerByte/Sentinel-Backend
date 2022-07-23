package objects

type GuildConfigPermissions struct {
	Edit int64 `json:"edit"`
	Read int64 `json:"read"`
}

type GuildConfigData struct {
	UseConfig bool `json:"use_config"`
}

type GuildConfig struct {
	Permissions GuildConfigPermissions `json:"permissions"`
	Data        GuildConfigData        `json:"data"`
	Preset      string                 `json:"preset"`
}

var DefaultGuildConfig = GuildConfig{
	Permissions: GuildConfigPermissions{
		Edit: 40,
		Read: 0xfffffffffff, // @everyone permissions
	},
	Data: GuildConfigData{
		UseConfig: false,
	},
	Preset: "default",
}
