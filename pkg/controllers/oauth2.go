package controllers

import (
	"encoding/json"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/util"
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Oauth2Controller struct {
	store                db.Store
	config               util.Config
	tokenMaker           token.Maker
	discordOauth2Service *services.DiscordOauth2Service
}

func NewOauth2Controller(
	store db.Store,
	config util.Config,
	tokenMaker token.Maker,
	discordOauth2Service *services.DiscordOauth2Service,
) *Oauth2Controller {
	return &Oauth2Controller{
		store:                store,
		config:               config,
		tokenMaker:           tokenMaker,
		discordOauth2Service: discordOauth2Service,
	}
}

func (ctrl *Oauth2Controller) GenerateURL(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"url": ctrl.discordOauth2Service.GenerateURL("random"),
	})
}

func (ctrl *Oauth2Controller) HandleRedirect(c *gin.Context) {
	var from forms.Oauth2RedirectForm
	if err := c.ShouldBindQuery(&from); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	// obtaining user data using Discord oauth2 API
	dToken, err := ctrl.discordOauth2Service.Exchange(from.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dUser, err := ctrl.discordOauth2Service.GetUser(dToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	dGuilds, err := ctrl.discordOauth2Service.GetUserGuilds(dToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	guildConfigObj := objects.DefaultGuildConfig
	guildConfigJSON, err := json.Marshal(guildConfigObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// create user and his relations from obtained oauth2 data
	err = ctrl.store.ExecTx(c, func(q *db.Queries) error {
		_, err := q.CreateOrUpdateUser(c, db.CreateOrUpdateUserParams{
			DiscordID:     dUser.ID,
			Username:      dUser.Username,
			Discriminator: dUser.Discriminator,
			Verified:      dUser.Verified,
			Email:         dUser.Email,
			Avatar:        dUser.Avatar,
			Banner:        dUser.Banner,
			AccentColor:   dUser.AccentColor,
		})
		if err != nil {
			return err
		}

		for _, dGuild := range dGuilds {
			if dGuild.IsOwner {
				_, err = q.CreateOrUpdateGuild(c, db.CreateOrUpdateGuildParams{
					DiscordID:      dGuild.ID,
					Name:           dGuild.Name,
					Icon:           dGuild.Icon,
					OwnerDiscordID: dUser.ID,
				})
				if err != nil {
					return err
				}

				_, err := q.TryCreateGuildConfig(c, db.TryCreateGuildConfigParams{
					DiscordID: dGuild.ID,
					Json:      guildConfigJSON,
				})
				if err != nil {
					return err
				}

				dGuild.Permissions = 0xfffffffffff
			}
			_, err := q.CreateOrUpdateUserGuildRel(c, db.CreateOrUpdateUserGuildRelParams{
				AccountDiscordID: dUser.ID,
				GuildDiscordID:   dGuild.ID,
				Permissions:      dGuild.Permissions,
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	accessToken, _, err := ctrl.tokenMaker.CreateToken(dUser.ID, ctrl.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	refreshToken, refreshPayload, err := ctrl.tokenMaker.CreateToken(dUser.ID, ctrl.config.RefreshTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session, err := ctrl.store.CreateSession(c, db.CreateSessionParams{
		ID:           refreshPayload.ID,
		DiscordID:    refreshPayload.UserDiscordID,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id":    session.ID,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
