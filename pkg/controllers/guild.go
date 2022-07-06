package controllers

import (
	"database/sql"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GuildController struct {
	store db.Store
}

func NewGuildController(store db.Store) *GuildController {
	return &GuildController{store: store}
}

func (ctrl *GuildController) Get(c *gin.Context) {
	var uri forms.GetGuildURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	guild, err := ctrl.store.GetGuild(c, uri.DiscordID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, guild)
}

func (ctrl *GuildController) GetAll(c *gin.Context) {
	payload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	guilds, err := ctrl.store.GetUserGuilds(c, payload.UserDiscordID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, guilds)
}
