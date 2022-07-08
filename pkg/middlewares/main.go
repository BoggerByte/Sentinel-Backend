package middlewares

import (
	"github.com/gin-gonic/gin"
)

type GuildConfig interface {
	Overwrite() gin.HandlerFunc
	Get() gin.HandlerFunc
}

type Permissions struct {
	GuildConfig
}

type Middlewares struct {
	CORS        gin.HandlerFunc
	Auth        gin.HandlerFunc
	Permissions Permissions
}
