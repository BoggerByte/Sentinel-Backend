package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	memdb "github.com/BoggerByte/Sentinel-backend.git/pkg/db/memory"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/db/sqlc"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/forms"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/modules/token"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/services"
	"github.com/BoggerByte/Sentinel-backend.git/pkg/utils"
	"github.com/BoggerByte/Sentinel-backend.git/pub/objects"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v9"
	"net/http"
)

type Oauth2Controller struct {
	store                db.Store
	memStore             memdb.Store
	config               utils.Config
	tokenMaker           token.Maker
	discordOauth2Service *services.DiscordOauth2Service
}

func NewOauth2Controller(
	store db.Store,
	memStore memdb.Store,
	config utils.Config,
	tokenMaker token.Maker,
	discordOauth2Service *services.DiscordOauth2Service,
) *Oauth2Controller {
	return &Oauth2Controller{
		store:                store,
		memStore:             memStore,
		config:               config,
		tokenMaker:           tokenMaker,
		discordOauth2Service: discordOauth2Service,
	}
}

func (ctrl *Oauth2Controller) NewURL(c *gin.Context) {
	state := utils.RandomString(32)

	err := ctrl.memStore.SetOauth2Flow(c, state, memdb.Oauth2Flow{
		Completed:     false,
		UserDiscordID: 0,
	}, ctrl.config.Oauth2FlowStateDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url":   ctrl.discordOauth2Service.NewURL(state),
		"state": state,
	})
}

func (ctrl *Oauth2Controller) DiscordCallback(c *gin.Context) {
	var form forms.Oauth2RedirectForm
	if err := c.ShouldBindQuery(&form); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	_, err := ctrl.memStore.GetOauth2Flow(c, form.State)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			err := errors.New("state not exists or expired")
			c.JSON(http.StatusMethodNotAllowed, errorResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	// obtaining user data using Discord oauth2 API
	dToken, err := ctrl.discordOauth2Service.Exchange(form.Code)
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

	defaultGuildConfigObj := objects.DefaultGuildConfig
	defaultGuildConfigJSON, _ := json.Marshal(defaultGuildConfigObj)

	// create user and his relations form obtained oauth2 data
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
					Json:      defaultGuildConfigJSON,
				})
				if err != nil && !errors.Is(err, sql.ErrNoRows) {
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
		"session_id":       session.ID,
		"access_token":     accessToken,
		"access_duration":  ctrl.config.AccessTokenDuration.Milliseconds(),
		"refresh_token":    refreshToken,
		"refresh_duration": ctrl.config.RefreshTokenDuration.Milliseconds(),
	})
}
