package controllers

import "github.com/gin-gonic/gin"

type Account interface {
	Get(c *gin.Context)
}

type Auth interface {
	FinalizeLogin(c *gin.Context)
	RefreshToken(c *gin.Context)
}

type Guild interface {
	Get(c *gin.Context)
	GetAll(c *gin.Context)
}

type GuildConfig interface {
	Overwrite(c *gin.Context)
	Get(c *gin.Context)
}

type Oauth2 interface {
	NewURL(c *gin.Context)
	DiscordCallback(c *gin.Context)
}

type Controllers struct {
	Account
	Auth
	Guild
	GuildConfig
	Oauth2
}

func errorResponse(err error) gin.H {
	return gin.H{"message": err.Error()}
}
