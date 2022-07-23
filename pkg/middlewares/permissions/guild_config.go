package permissions

import (
	"database/sql"
	"encoding/json"
	"errors"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GuildConfigPermissions struct {
	store db.Store
}

func NewGuildConfigPermissions(store db.Store) *GuildConfigPermissions {
	return &GuildConfigPermissions{
		store: store,
	}
}

func (p *GuildConfigPermissions) Overwrite() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req forms.RequireDiscordIDRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		authPayload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

		userGuildRel, err := p.store.GetUserGuildRel(c, db.GetUserGuildRelParams{
			AccountDiscordID: authPayload.UserDiscordID,
			GuildDiscordID:   req.DiscordID,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := errors.New("no relations with guild")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		guildConfig, err := p.store.GetGuildConfig(c, req.DiscordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := errors.New("guild config not found")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		var guildConfigObj = objects.DefaultGuildConfig
		err = json.Unmarshal(guildConfig.Json, &guildConfigObj)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if !utils.AnyOfPermissions(userGuildRel.Permissions, guildConfigObj.Permissions.Edit) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.Next()
	}
}

func (p *GuildConfigPermissions) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req forms.RequireDiscordIDRequest
		if err := c.ShouldBindUri(&req); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}

		authPayload := c.MustGet("authorization_payload").(*token.Payload)

		userGuildRel, err := p.store.GetUserGuildRel(c, db.GetUserGuildRelParams{
			AccountDiscordID: authPayload.UserDiscordID,
			GuildDiscordID:   req.DiscordID,
		})
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := errors.New("no relations with guild")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		guildConfig, err := p.store.GetGuildConfig(c, req.DiscordID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				err := errors.New("guild config not found")
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": err.Error()})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		var guildConfigObj = objects.DefaultGuildConfig
		err = json.Unmarshal(guildConfig.Json, &guildConfigObj)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		if !utils.AnyOfPermissions(userGuildRel.Permissions, guildConfigObj.Permissions.Read) {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		c.Next()
	}
}
