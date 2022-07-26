package middlewares

import (
	"errors"
	"fmt"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	AuthorizationHeaderKey  = "authorization"
	AuthorizationTypeBearer = "bearer"
	AuthorizationPayloadKey = "authorization_payload"
)

func NewAuthMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			err := errors.New("authentication header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 {
			err := errors.New("invalid authentication header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		authType := strings.ToLower(fields[0])
		if authType != AuthorizationTypeBearer {
			err := fmt.Errorf("unsupported authentication type: %s", authType)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Set(AuthorizationPayloadKey, payload)
		c.Next()
	}
}

func NewDiscordBotAuthMiddleware(config utils.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		sentinelAPISecret := c.GetHeader(AuthorizationHeaderKey)
		if len(sentinelAPISecret) == 0 {
			err := errors.New("authentication header is not provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if sentinelAPISecret != config.DiscordBotSentinelAPISecret {
			err := errors.New("invalid credentials")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		c.Next()
	}
}
