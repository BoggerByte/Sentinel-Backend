package controllers

import (
	"database/sql"
	"errors"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/middlewares"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GuildController struct {
	store db.Store
}

type ResponseGuild struct {
	ID             int64  `json:"id"`
	DiscordID      string `json:"discord_id"`
	OwnerDiscordID string `json:"owner_discord_id"`
	Name           string `json:"name"`
	Icon           string `json:"icon"`
	CanReadConfig  bool   `json:"can_read_config"`
	CanEditConfig  bool   `json:"can_edit_config"`
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

func (ctrl *GuildController) GetUserOne(c *gin.Context) {
	var uri forms.GetUserGuildURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	payload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	guild, err := ctrl.store.GetUserGuild(c, db.GetUserGuildParams{
		AccountDiscordID: payload.UserDiscordID,
		GuildDiscordID:   uri.DiscordID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, ResponseGuild{
		ID:             guild.ID,
		DiscordID:      guild.DiscordID,
		OwnerDiscordID: guild.OwnerDiscordID,
		Name:           guild.Name,
		Icon:           guild.Icon,
		CanReadConfig:  utils.AnyOfPermissions(guild.Permissions, guild.ConfigRead),
		CanEditConfig:  utils.AnyOfPermissions(guild.Permissions, guild.ConfigEdit),
	})
}

func (ctrl *GuildController) GetUserAll(c *gin.Context) {
	payload := c.MustGet(middlewares.AuthorizationPayloadKey).(*token.Payload)

	guilds, err := ctrl.store.GetUserGuilds(c, payload.UserDiscordID)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	var rGuilds []ResponseGuild
	for _, guild := range guilds {
		rGuilds = append(rGuilds, ResponseGuild{
			ID:             guild.ID,
			DiscordID:      guild.DiscordID,
			OwnerDiscordID: guild.OwnerDiscordID,
			Name:           guild.Name,
			Icon:           guild.Icon,
			CanReadConfig:  utils.AnyOfPermissions(guild.Permissions, guild.ConfigRead),
			CanEditConfig:  utils.AnyOfPermissions(guild.Permissions, guild.ConfigEdit),
		})
	}

	c.JSON(http.StatusOK, rGuilds)
}
