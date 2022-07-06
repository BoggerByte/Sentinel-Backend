package controllers

import "github.com/gin-gonic/gin"

type Account interface {
	Get(c *gin.Context)
}

type Auth interface {
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
	GenerateURL(c *gin.Context)
	HandleRedirect(c *gin.Context)
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
