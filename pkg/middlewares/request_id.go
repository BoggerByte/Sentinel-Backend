package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
)

func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}
