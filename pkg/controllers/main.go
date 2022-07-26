package controllers

import "github.com/gin-gonic/gin"

type User interface {
	GetUser(c *gin.Context)
}

type Auth interface {
	RefreshToken(c *gin.Context)
}

type Guild interface {
	GetGuild(c *gin.Context)
	GetUserGuild(c *gin.Context)
	GetUserGuilds(c *gin.Context)
	CreateOrUpdateGuilds(c *gin.Context)
	CreateOrUpdateGuild(c *gin.Context)
}

type GuildConfig interface {
	OverwriteGuildConfig(c *gin.Context)
	GetGuildsConfigs(c *gin.Context)
	GetGuildConfig(c *gin.Context)
	GetGuildConfigPreset(c *gin.Context)
}

type Oauth2 interface {
	GetNewOauth2URL(c *gin.Context)
	GetNewInviteBotURL(c *gin.Context)
	HandleDiscordCallback(c *gin.Context)
}

type Controllers struct {
	User
	Auth
	Guild
	GuildConfig
	Oauth2
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
