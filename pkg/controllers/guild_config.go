package controllers

import (
	"encoding/json"
	db "github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
	"github.com/gin-gonic/gin"
	"net/http"
)

type GuildConfigController struct {
	store db.Store
}

func NewGuildConfigController(store db.Store) *GuildConfigController {
	return &GuildConfigController{
		store: store,
	}
}

func (ctrl *GuildConfigController) GetGuildConfigPreset(c *gin.Context) {
	var uri forms.GetGuildConfigPresetURI
	if err := c.ShouldBindUri(&uri); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var guildConfig objects.GuildConfig
	switch uri.Preset {
	case "default":
		guildConfig = objects.DefaultGuildConfig
	}

	c.JSON(http.StatusOK, guildConfig)
}

func (ctrl *GuildConfigController) OverwriteGuildConfig(c *gin.Context) {
	var uri forms.RequireDiscordIDRequest
	_ = c.ShouldBindUri(&uri)
	var newGuildConfig forms.OverwriteGuildConfigJSON
	if err := c.ShouldBindJSON(&newGuildConfig); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	newGuildConfigJSON, err := json.Marshal(newGuildConfig)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	_, err = ctrl.store.CreateOrUpdateGuildConfig(c, db.CreateOrUpdateGuildConfigParams{
		DiscordID: uri.DiscordID,
		Json:      newGuildConfigJSON,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}

func (ctrl *GuildConfigController) GetGuildConfig(c *gin.Context) {
	var uri forms.RequireDiscordIDRequest
	_ = c.ShouldBindUri(&uri)

	guildConfig, err := ctrl.store.GetGuildConfig(c, uri.DiscordID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, guildConfig)
}

func (ctrl *GuildConfigController) GetGuildsConfigs(c *gin.Context) {
	guildsConfigs, err := ctrl.store.GetGuildsConfigs(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, guildsConfigs)
}
